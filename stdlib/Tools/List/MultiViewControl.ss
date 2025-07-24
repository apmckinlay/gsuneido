// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
CommandParent
	{
	Name: 'MultiView'

	constructed?: false
	New(.args, .accessArgs, .listArgs, .accessGoTo? = false, defaultMode = 0,
		.embedded = false)
		{
		super(.layout(args))
		.modeSaveName = args.option $ '_ListMode'
		.query = args[0]
		.access = .FindControl('Access')
		.list = .FindControl('VirtualList')
		.flip = .FindControl('Flip')
		.total = .FindControl('ThreadTotal')
		.locate = .access.FindControl('Locate')
		.access.Redir('On_Flip_To_List', this)
		.list.Redir('On_Flip_To_Access', this)
		.constructed? = true

		userSetting = UserSettings.Get(.modeSaveName, 'none') // false is accessmode
		if .embedded
			{
			.list.Commands = #()
			.access.Commands = #()
			userSetting = true
			}

		.initFlip(userSetting, defaultMode)
		.initSelects(args)
		if false isnt highlight = args.GetDefault('highlightrows', false)
			.list.HighlightValues(@highlight)
		}

	initFlip(userSetting, defaultMode, _accessGoTo? = false)
		{
		if accessGoTo?
			.accessGoTo? = true
		if (not .accessGoTo? and
			(userSetting is true or (userSetting is 'none' and defaultMode is .listMode)))
			.flipToList(fromNew?:)
		}

	initSelects(args)
		{
		selectName = AccessSelectMgr.DefaultSelectName(
			this, args.GetDefault('title', ''), .query)
		.access.SetSelectMgr(selectName, .accessArgs)
		if .embedded
			return
		.list.SetQuery(.query, .listArgs.columns, :selectName)

		if .accessMode?()
			.accessSelectsSetup()
		else
			.listSelectsSetup()
		}

	accessInitialRecordLoaded: false
	accessSelectsSetup()
		{
		.access.ApplySelects()
		if .accessInitialRecordLoaded
			return
		.access.Load_initial_record(.accessArgs)
		.access.SetDefaultStatus()
		.accessInitialRecordLoaded = true
		}

	listDefaultStatusSet: false
	listSelectsSetup()
		{
		.list.ApplySelects()
		if .listDefaultStatusSet
			return
		.list.SetDefaultStatus()
		.listDefaultStatusSet = true
		}

	OverrideSelectManager?()
		{
		return false
		}

	Msg(args)
		{
		// Distinguishes between messages sent from the MultiView Virtual List and a VL
		// inside an AccessControl. Assumes that any VL in an access will be named
		if args.source.Base?(VirtualListControl) and args.source.Name isnt 'VirtualList'
			return .Send(@args)
		return super.Msg(args)
		}

	PrintReport(reportClass, source, data = false, permission = false,
		additionalArgs = false)
		{
		if permission isnt false and AccessPermissions(permission) isnt true
			{
			.AlertInfo(.access.GetTitle(), 'You do not have permission for this option')
			return
			}

		if not .SaveFor('Print ' $ reportClass.AccessTitle $ ' on')
			return

		args = Object(data is false ? .GetData(source) : data)
		if additionalArgs isnt false
			args.Merge(additionalArgs)
		ToolDialog(.Window.Hwnd, reportClass(@args))
		if .AccessMode?()
			.Refresh()
		}

	Commands: #(
		#('Flip', 'Alt+A')
		)

	layout(args, _accessGoTo? = false)
		{
		.accessArgs = args.Copy().Merge(.accessArgs)
		.accessArgs.fromMultiView = true
		.listArgs = args.Copy().Merge(.listArgs)
		.controlCustomization()

		access = .accessArgs.Copy().Add('Access', at: 0)
		list = .listArgs.Copy().Add('VirtualList', at: 0)
		// setting up VL with empty query / columns so we only do 1 query when setting up
		list[1] = ''
		list.columns = #()
		list.enableUserDefaultSelect = not accessGoTo?
		if .embedded
			list.enableUserDefaultSelect = false

		list.titleLeftCtrl = .flipButton(highlighted: 0, command: 'Flip_To_Access')
		list.switchToForm = true
		if not .listArgs.Member?('filtersOnTop')
			list.filtersOnTop = true
		access.titleLeftCtrl = .flipButton(highlighted: 1, command: 'Flip_To_List')
		.ensureDefaultCols()
		totalLayout = .listArgs.Extract('totalLayout', #())
		if Object?(totalLayout) and not totalLayout.Empty?()
			list = Object('Vert', list, totalLayout)

		return Object('Flip', access, list)
		}

	controlCustomization()
		{
		// excludeCustomize? can be set in args, accessArgs, or listArgs, (setting the
		// controls accordingly). If .accessArgs has "excludeCustomize?: true",
		// then Customizable is NOT required.
		if .accessArgs.GetDefault('excludeCustomize?', false)
			return
		if .listArgs.GetDefault('protectField', false) isnt false and
			false is Display(.accessArgs[1]).Has?('Customizable')
			ProgrammerError('Multiview Access is missing Customizable Layout')
		}

	flipButton(highlighted, command)
		{
		tip = not .embedded
			? 'switch between list and form view (Alt+A)'
			: 'access this record in form view (Alt+A)'
		return Object('Border',
			Object('EnhancedButton',
				image: Object('view_list.emf', 'view_form.emf', :highlighted),
				mouseEffect:, imagePadding: .1, :command, :tip),
			border: 6)
		}
	ensureDefaultCols()
		{
		Assert(.listArgs hasMember: 'option')
		Assert(.listArgs hasMember: 'title')
		Assert(.listArgs hasMember: 'defaultColumns')
		if not .listArgs.Member?('columns')
			.listArgs.columns = QuerySelectColumns(.listArgs[0]).
				RemoveIf(Internal?).
				Difference(.listArgs.GetDefault('excludeSelectFields', #()))
		}

	AccessGoto(field, value, wrapper = false)
		{
		.access.AccessGoto(field, value)
		if wrapper isnt false
			wrapper.Redir('On_Flip', this)
		}

	GotoRecord(field, value)
		{
		if .AccessMode?()
			.AccessGoto(field, value)
		else
			{
			if false is .list.Save() // to mimic access
				return
			.list.Refresh()
			.list.SelectRecordByKeyPair(value, field)
			}
		}

	VirtualList_GetCustomKey()
		{
		return AccessControl.BuildCustomKey(.args.option, .args.title)
		}

	VirtualList_GetOption()
		{
		return .args.option
		}

	// START - Methods to handle common operations on both list and access
	Access_NewRecord(data)
		{
		return .Send('MultiView_NewRecord', data)
		}

	VirtualList_NewRowAdded(rec)
		{
		return .Send('MultiView_NewRecord', rec)
		}

	SetEditMode(source = false)
		{
		if .accessMode?()
			return .access.SetEditMode()
		return source isnt false
			? .list.On_Edit(source, force:)
			: false isnt .list.ForceEditMode(.list.GetSelectedRecord())
		}

	GetMultiFilter()
		{
		if .accessMode?()
			return false
		return .list.Select_vals
		}

	GetData(source)
		{
		return .accessMode?() ? .access.GetData() : .list.GetData(source)
		}

	GetCurrentSelectedData()
		{
		return .accessMode?() ? .access.GetData() : .list.GetSelectedRecord()
		}

	GetSelectedRecords()
		{
		return .accessMode?() ? Object(.access.GetData()) : .list.GetSelectedRecords()
		}

	Refresh()
		{
		if .accessMode?()
			{
			.access.Reload()
			return
			}
		.calculateTotal(true, false)
		.list.Refresh()
		}

	Refresh1(rec)
		{
		// Refresh1 discards changes.
		// Does not handle .calculateTotal, but could be added in future
		if .accessMode?()
			return .access.Reload()

		result = .list.ReloadRecord(rec, discard:)
		key = ShortestKey(.access.GetKeys())
		.list.SelectRecordByKeyPair(rec[key], key)
		return result
		}

	Valid?()
		{
		return .accessMode?() ? .access.Valid?() : .list.Valid?()
		}

	GetCtrl(data)
		{
		return .accessMode?() ? .access : .list.GetExpandedControl(data)
		}

	GetControl(field)
		{
		return .accessMode?() ? .access.GetControl(field) : .list
		}

	AccessMode?()
		{ return .accessMode?() }

	GetCurrentModeCtrl()
		{
		return .accessMode?() ? .access : .list
		}

	EditMode?(source)
		{
		if .accessMode?()
			return .access.EditMode?()

		if false is ctrl = .GetCtrl(.GetData(source))
			return false
		return not ctrl.GetReadOnly()
		}

	Save()
		{
		return .accessMode?() ? .access.Save() : .list.Save()
		}

	SaveFor(action)
		{
		return .accessMode?() ? .access.SaveFor(action) : .list.Save()
		}

	SaveForAndEnsureEdit(action, source = false)
		{
		if .accessMode?()
			return .access.SaveForAndToggleEdit(action)

		if false is .list.Save()
			return false
		result = source is false
			? .list.ForceEditMode(.list.GetSelectedRecord())
			: .list.On_Edit(source)
		return false isnt result
		}

	GetOrigin(data)
		{
		return .accessMode?() ? .access.GetOriginal() : data.vl_origin
		}

	NewRecord?(data)
		{
		// not sure if this is correct
		return .accessMode?() ? .access.NewRecord?() : data.New?()
		}

	Goto(field, value)
		{
		if .accessMode?()
			.access.AccessGoto(field, value)
		else
			{
			.On_Flip_To_Access()
			if .accessMode?()
				.access.AccessGoto(field, value)
			}
		}

	accessMode?()
		{
		if .embedded and .flip.GetCurrent() is .accessMode
			ProgrammerError('Embedded MultiView should never be in Access Mode')
		return  .flip.GetCurrent() is .accessMode
		}
	// END - Methods to handle common operations on both list and access - End


	On_Flip()
		{
		if false is .Send('AllowFlip?')
			return

		if .accessMode?()
			.On_Flip_To_List()
		else
			.On_Flip_To_Access()
		}

	On_Flip_To_List()
		{
		if .access.Save()
			.flipToList()
		}

	listMode: 1
	newRecordsSinceFlip: false
	modifiedRecordsSinceFlip: false
	flipToList(fromNew? = false)
		{
		if .needToApplySelects?(fromNew?)
			.listSelectsSetup()
		else
			.reloadModifiedRecords(fromNew?)
		.flip.SetCurrent(.listMode)
		.Defer(uniqueID: 'listSetFocus')
			{
			// ensure the default focus does not go to Access Locate,
			// which would cause the screen to freeze when pressing up/down arrow keys
			if not .Destroyed?() and not .accessMode?()
				SetFocus(.list.GetGridHwnd())
			}
		.newRecordsSinceFlip = false
		if not fromNew?
			{
			.Send('MultiView_EnterListMode', list: .list)
			.redirToList()
			}
		NotesControl.UpdateStatus(.list)
		}

	CommandParent_SkipInitCommands?()
		{
		return true
		}

	Startup()
		{
		super.Startup()
		if .accessMode?()
			.redirToAccess()
		else
			.redirToList()
		}

	redirToList()
		{
		if .Destroyed?() or .accessMode?() or .embedded
			return

		.access.RemoveCommands()
		.list.AddCommands()
		}

	needToApplySelects?(fromNew?)
		{
		if fromNew?
			return false
		return .access.Select_vals isnt .prevSelVal or .newRecordsSinceFlip
		}

	reloadModifiedRecords(fromNew?)
		{
		if fromNew?
			return

		if .modifiedRecordsSinceFlip is false
			return

		.modifiedRecordsSinceFlip = false
		.calculateTotal(true, false)
		}

	calculateTotal(saved?, listDirty?)
		{
		if .total is false
			return

		.total.AfterChanged(saved?, .list.GetQuery(), listDirty?, .list.Select_vals)
		}

	VirtualList_SwitchToForm()
		{
		.On_Flip_To_Access()
		}

	On_Flip_To_Access()
		{
		empty? = .list.GetLoadedData().Empty?()
		if not .list.SaveOutstandingChanges()
			{
			.AlertInfo(.access.GetTitle(), 'Please correct the highlighted line.')
			return
			}
		rec = false
		if not empty? and false is rec = .list.GetSelectedRecord()
			{
			.AlertInfo(.access.GetTitle(), 'Please select a line')
			return
			}
		goToField = .getGoToField()
		if .embedded
			{
			.flip_AccessGoTo(empty?, goToField, rec)
			return
			}
		.flipToAccess()
		if rec isnt false
			{
			.access.AccessGoto(goToField, rec[goToField])
			.focusLocate()
			}
		else if .access.RecordSet?() // false on an empty readonly table
			.access.Reload() // fix stale data between list vs access
		NotesControl.UpdateStatus(.access)
		}

	getGoToField()
		{
		sort = QueryGetSort(.query).RemovePrefix("reverse ")
		keys = .access.GetKeys()
		if sort.Split(',').Size() is 1 and keys.Has?(sort)
			return sort

		return ShortestKey(keys)
		}

	flip_AccessGoTo(empty?, key, rec)
		{
		if empty? or rec is false or
			0 is accessInfo = .Send('MultiView_EmbeddedAccessInfo', rec)
			return

		AccessGoTo(accessInfo.accessScreen, key, rec[key], .Window.Hwnd,
			defaultSelect: accessInfo.defaultSelect, onDestroy:	.list.Refresh)
		return
		}

	accessMode: 0
	prevSelVal: false
	flipToAccess()
		{
		.prevSelVal = .list.Select_vals.DeepCopy()
//Print('.prevSelVal' : .prevSelVal)
		.flip.SetCurrent(.accessMode)
		.redirToAccess()
		.accessSelectsSetup()
		}

	focusLocate()
		{
		if not .access.HasLocate?()
			return

		.locate.SelectAll()
		SetFocus(.locate.EditHwnd())
		}

	redirToAccess()
		{
		if .embedded
			return
		.list.RemoveCommands()
		.access.AddCommands()
		}

	// START - Common Message Handlers
	VirtualList_BeforeSave(@args)
		{
		return .Send('MultiView_BeforeSave', args)
		}

	AccessBeforeSave(@args)
		{
		args.data = .access.GetData()
		return .Send('MultiView_BeforeSave', args)
		}

	VirtualList_RecordChange(member, record)
		{
		.Send('MultiView_RecordChange', member, record)
		}

	Access_RecordChange(members)
		{
		record = .access.GetData()
		// NOTE: Careful with converting Access to MultiView.
		// In Access, members are all in one Object, in multiview it will be a string.
		for member in members
			.Send('MultiView_RecordChange', member, record)
		}

	Access_ConfirmDestroy()
		{
		.Send('MultiView_ConfirmDestroy')
		}

	VirtualList_AfterField(field, value, data)
		{
		.Send('MultiView_AfterField', field, value, data)
		}

	Access_AfterField(field, value)
		{
		if not .constructed?
			return

		data = .access.GetData()
		.Send('MultiView_AfterField', field, value, data)
		}

	AccessBeforeSaving()
		{
		data = .access.GetData()
		.Send('MultiView_BeforeSaving', data)
		}

	VirtualList_BeforeSave_PreTran(rec)
		{
		.Send('MultiView_BeforeSaving', rec)
		}

	Access_AllowDelete()
		{
		return .Send('MultiView_AllowDelete', .access.GetData())
		}

	VirtualList_AllowDelete(rec)
		{
		if .accessMode?()
			return false
		return .Send('MultiView_AllowDelete', rec)
		}

	VirtualList_AllowInsert()
		{
		return not .accessMode?()
		}

	VirtualList_AllowNewRecord()
		{
		if .accessMode?()
			return false

		return false isnt .Send('MultiView_AllowNewRecord')
		}

	AccessBeforeDelete(t)
		{
		.Send('MultiView_BeforeDelete', .access.GetData(), t)
		}

	VirtualList_BeforeDelete(data, t)
		{
		.Send('MultiView_BeforeDelete', data, t)
		}

	AccessAfterSave(t)
		{
		.Send('MultiView_AfterSave', .access.GetData(), t)
		.orig = .access.GetOriginal().Copy()
		}

	VirtualList_AfterSave(data, t)
		{
		.Send('MultiView_AfterSave', data, t)
		}

	VirtualList_AfterSaving(record)
		{
		.Send('MultiView_AfterSaving', record)
		}

	SelectControl_Changed()
		{
		if .accessMode?()
			.list.SetSelectVals(.access.Select_vals.Copy())
		else
			.access.SetSelectVals(.list.Select_vals.Copy())
		}

	VirtualList_BeforeApplySelect(@unused)
		{
		.access.SetSelectVals(.list.Select_vals.Copy())
		return true
		}

	Access_AfterNewRecord(@unused)
		{
		.newRecordsSinceFlip = true
		}

	AccessAfterSaving()
		{
		.Send('MultiView_AfterSaving', newRec = .access.GetData())
		.list.ReplaceRecord(.orig, newRec) // reload if loaded; also handles key update
		.modifiedRecordsSinceFlip = true
		}

	AccessAfterDeleting()
		{
		key = ShortestKey(.access.GetKeys())
		.list.RemoveRowByKeyPair(.access.GetOriginal()[key], key)
		}

	Access_SetRecord(x, source)
		{
		.Send('MultiView_ItemSelected', x, ctrl: source)
		}

	VirtualList_ItemSelected(rec, source)
		{
		.Send('MultiView_ItemSelected', rec, ctrl: source)
		}

	VirtualList_AfterChanged(saved)
		{
		.calculateTotal(saved, .list.VirtualListDirty?())
		}

	VirtualList_SetWhere()
		{
		.calculateTotal(true, false)
		}

	Access_BeforeRecord(rec)
		{
		.Send('MultiView_BeforeRecord', rec)
		}

	VirtualList_ExtraSetupRecordFn()
		{
		return .HandleBeforeRecord
		}

	HandleBeforeRecord(rec)
		{
		.Send('MultiView_BeforeRecord', rec)
		}
	// END - Common Message Handlers

	Destroy()
		{
		if not .accessGoTo? and not .embedded
			UserSettings.Put(.modeSaveName, .flip.GetCurrent() is .listMode)
		super.Destroy()
		}
	}
