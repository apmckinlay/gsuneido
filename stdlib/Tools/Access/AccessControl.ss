// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
CommandParent
	{
	Name: "Access"
	last_button: false
	select_button: false
	new_button: false
	data: false
	New(@args)
		{
		super(.makecontrols(args))
		.data = .Vert.Scroll.Border.Data
		.data.SetProtectField(.protectField)
		.data.AddObserver(.Access_RecordChange)

		for ctrl in .Vert.Horz.HorzEven.GetChildren()
			if ctrl.Name.Has?('New')
				.new_button = ctrl

		.status = .Vert.GetDefault('Status', FakeObject(Set:, GetValid:))

		.Setup()
		.lock = AccessLock(this, false)
		.loopedAddons = AccessLoopAddonManager(this)

		.saveOnlyLinked = args.Member?('saveOnlyLinked') and
			args.saveOnlyLinked is true
		.protect = args.GetDefault(#protect, false) is true

		.Window.AddValidationItem(this)

		.model.Plugins_Init()

		if args.GetDefault('fromMultiView', false) is false
			{
			.initDefaultSelect(args)
			.Load_initial_record(args)
			.SetDefaultStatus()
			}
		keys = .model.GetKeyField().Split(",")
		.attachmentsManager = AttachmentsManager(.model.GetQuery(), keys)
		}

	initDefaultSelect(args)
		{
		selName = AccessSelectMgr.DefaultSelectName(this, .title, .query)
		.SetSelectMgr(selName, args)
		.setInitialWhere(fromNew?:)
		}

	setInitialWhere(fromNew? = false)
		{
		if .ApplySelects(fromNew?) is false
			return
		if not fromNew?
			{
			.Select_vals.Each({ it.check = false })
			AlertDelayed('No records found that ' $
				'match the current select.\r\n\r\nSelection has been reset.',
				.title, flags: MB.ICONINFORMATION)
			}
		}

	SetDefaultStatus()
		{
		.status.SetDefaultMsg(.selectMgr.UsingDefaultFilter?()
			? 'initial Select applied - for more details click on "Select"'
			: '')
		}

	accessGoTo: false
	ApplySelects(fromNew? = false, _accessGoTo? = false)
		{
		.accessGoTo = accessGoTo?
		where = SelectRepeatControl.BuildWhere(.sf, .Select_vals)
		whereStr = .sf.Joins(where.joinflds) $ where.where
		found = .SetWhere(whereStr, quiet:)
		if not found and .selectMgr.HasSavedDefault? and not accessGoTo?
			{
			.Select_vals.Each({ it.check = false })
			.SetWhere('', quiet:)
			.Send('SelectControl_Changed')
			selectType = fromNew? ? 'default' : 'current'
			AlertDelayed('No records found that ' $
				'match the ' $ selectType $ ' select.\r\n\r\nSelection has been cleared.',
				.title, flags: MB.ICONINFORMATION)
			return false
			}
		return true
		}

	Startup()
		{
		super.Startup()
		.setFocusToLocate()
		}
	Load_initial_record(args)
		{
		.start_last? = false
		if args.Member?('startLast') and args.startLast is true
			{
			.start_last? = true
			.On_Last()
			}
		else if args.Member?('startNew') and args.startNew is true
			.On_New()
		else
			.On_First()

		// the following is necessary for when the table is empty
		if .protect and .record is false or ReadOnlyAccess(this) is true
			.data.SetReadOnly(true)
		}
	locate_status: false
	setFocusToLocate()
		{
		if .HasLocate?() and not .newrecord?
			{
			.locate.SelectAll()
			SetFocus(.locate.EditHwnd())
			}
		}
	Commands:
		(
		("New",		"Ctrl+N")
		("Edit",	"Alt+E")
		("First",	"Alt+F")
		("Prev",	"Alt+P")
		("Next",	"Alt+N")
		("Last",	"Alt+L")
		("Select",	"Alt+S")
		("Locate",	"Ctrl+L")
		("NextTab",	"Ctrl+Tab")
		("PrevTab",	"Shift+Ctrl+Tab")
		)
	Setup()
		{
		.locate = .Vert.Horz.Locate
		.select = false
		.edit_button = class
			{
			Pushed?(unused) { return false }
			Grayed(unused) { return false }
			GetEnabled() { return false }
			SetEnabled(unused) { return false }
			}
		if not .linked?
			{
			.edit_button = .Vert.Horz.HorzEven.Edit
			.select_button = .Vert.Horz.HorzEven.Select
			}
		.first_button = .Vert.Horz.HorzEven.First
		.last_button = .Vert.Horz.HorzEven.Last
		.prev_button = .Vert.Horz.HorzEven.Prev
		.next_button = .Vert.Horz.HorzEven.Next
		}
	makecontrols(args)
		{
		.types = AccessTypes(args)
		.model = AccessModel(args, .types.DynamicTypes)
		.validField = args.GetDefault("validField", false)
		.warningField = args.GetDefault("warningField", false)
		.protectField = args.GetDefault("protectField", false)
		.historyFields = args.GetDefault("historyFields", false)
		.Addons = AddonManager(this, args.GetDefault('addons', #()))
		// DEPRECATING: defaultNewValues is to be removed under: 33976
		.defaultNewValues = args.GetDefault("defaultNewValues", false)
		.nextNumber = AccessNextNum(this, args.GetDefault("nextNum", false))
		return .SetupControls(args)
		}
	SetupControls(args)
		{
		.setTitles(args)
		.option = args.GetDefault('option', 'Access')
		.menus = RecordMenuManager(.protectField, .option, .historyFields, ctrl: this,
			warningField: .warningField)
		.set_customization()
		custom_screen = not args.GetDefault('excludeCustomize?', false)
			? .customIcon
			: false
		return Object('Vert',
			.title isnt ""
				? Object('CenterTitle', .title,
					titleLeftCtrl: args.GetDefault('titleLeftCtrl', false),
					:custom_screen)
				: #(Skip 0),
			Object('Scroll', Object('Border', Object('Record', .control_ob(args),
				custom: .customFields, accessControl:), border: 5)),
			Object('Horz',
				.get_buttons(.new_button_control(args), .menus.Current, .menus.Global,
					args),
				.model.GetLocateLayout()
				)
			#(Skip 2)
			'Status'
			)
		}
	customKey: false
	GetAccessCustomKey()	// message send from BrowseControl
		{
		return .customKey
		}
	set_customization()
		{
		.customKey = .BuildCustomKey(.option, .title)
		.customFields = Customizable.GetCustomizedFields(.customKey)
		.customIcon = .customKey isnt false
		}
	BuildCustomKey(option, title)
		{
		return option isnt 'Access' and option isnt '' and title isnt ''
			? title $ ' ~ ' $ option
			: false
		}
	setTitles(args)
		{
		.title = args.Member?("title") ? args.title : .query
		// need to do the following to ensure the width isnt governed by the
		// length of the query string (Static control for title).
		if args.Member?("dynamicTypes")
			.title = args.title
		.Title = .title $ " - Access"
		}
	GetTitle()
		{
		return .title
		}
	control_ob(args)
		{
		switch (args.Size(list:))
			{
		case 1 : // just query
			control = Object("Vert")
			if .types.DynamicTypes is false
				control.Add(@.fields)
		case 2 : // single control
			control = args[1]
		default : // form
			control = Object("Form").Add(@args.Values(list:)[1..])
			}
		return control
		}
	getter_fields()
		{
		return .model.GetFields()
		}
	getter_sf()
		{
		return .model.GetSelectFields()
		}
	Getter_Option() // used by Select presets
		{ return .option }

	new_button_control(args)
		{
		newbutton = #(Button, "&New", tip: 'Ctrl+N', xstretch: 1)
		if args.Member?('newOptions')
			types = args.newOptions
		else
			types = .types.DynamicTypeList(omitTypes:)
		if Object?(types) and types.Size() > 1
			newbutton = Object('MenuButton', 'New', types, name: 'New')
		return newbutton
		}

	linked?: false
	get_buttons(newbutton, current_menu, global_menu, args)
		{
		.linked? = args.Member?("linked?") and args.linked? is true
		buttons = Object('HorzEven'
			newbutton
			#(EnhancedButton, text: "First", tip: 'Alt+F', buttonStyle:, mouseEffect:)
			#(EnhancedButton, text: "Prev", tip: 'Alt+P', buttonStyle:, mouseEffect:)
			#(EnhancedButton, text: "Next", tip: 'Alt+N', buttonStyle:, mouseEffect:)
			#(EnhancedButton, text: "Last", tip: 'Alt+L', buttonStyle:, mouseEffect:)
			Object('MenuButton', 'Current', current_menu, sendParents?:))
		if .linked?
			return buttons
		buttons = buttons.Add(#(EnhancedButton, text: "Edit", tip: 'Alt+E',
			buttonStyle:, mouseEffect:), at: 2)
		buttons = buttons.Add(#(EnhancedButton, text: "Select...", tip: 'Alt+S',
			buttonStyle:, mouseEffect:), at: 7)
		return buttons.Add(Object('MenuButton', 'Global', global_menu))
		}
	record: false
	RecordSet?() // false when setdata is not called, like on an empty readonly table
		{
		return .record isnt false
		}
	original_record: false
	BeforeRecord(x) // default sends message
		{
		.Send("Access_BeforeRecord", x)
		.Addons.Send('BeforeRecord', x)
		}

	lastSetDataTime: false
	setdata(x, newrec = false)
		{
		.Send("Access_BeforeLoadingRecord", x)
		.types.DetectTypeChange(newrec, x, .change_type)
		.model.NotifyObservers('before_setdata')
		.nextNumber.SetData(x, newrec)
		.BeforeRecord(x)
		.record = x
		.model.SetKeyQuery(.record)
		.original_record = x.Copy()
		if .protectField isnt false
			.record[.protectField]
		if newrec
			Customizable.SetRecordDefaultValues(.customKey, .record, .protectField)
		CustomizeField.SetFormulas(.customKey, .record, .protectField)
		.data.Set(.record)
		.record_change_members = false
		.Send("Access_SetRecord", :x)
		.model.Plugins_Execute(access: this, member: 'setdata',
			pluginType: 'AccessObservers')
		.newrecord? = newrec
		if .new_button isnt false and .new_button.Method?("Pushed?")
			.new_button.Pushed?(newrec)
		if (.protect)
			.data.SetReadOnly(true)
		.setWarnings()
		.model.NotifyObservers('setdata')
		.lastSetDataTime = Date()
		if not newrec
			.view_mode()
		.loopedAddons.Start()
		}

	change_type(typename, control)
		{
		.title = typename
		.Vert.CenterTitle.Set(.title)
		.set_customization()
		.Vert.Remove(1)
		.Vert.Insert(1, Object('Scroll',
			Object('Border', Object('Record', control, custom: .customFields),
				border: 5)))
		.data = .Vert.Scroll.Border.Data
		.data.SetProtectField(.protectField)
		.data.AddObserver(.Access_RecordChange)
		.ResetCommands()
		}
	newrecord?: false
	// we need this new_setdata? flag for cases where we want setdata to get done
	// in the On_New method, like delete and restore
	new_setdata?: false
	On_New(@args)
		{
		rec = Record()
		if (ReadOnlyAccess(this) is true or not .Save() or .protect or
			false is AllowInsertRecord?(rec, .protectField))
			return
		sameType? = .types.SetCurrentType(args)
		// return if already on new non-dirty record AND type hasn't changed
		if .newNonDirtyRecord?(sameType?)
			return

		.c.Rewind()
		.setdata(rec, newrec:)
		.data.SetReadOnly(false)
		.types.ApplyStickyValues(.record)
		.set_position_button_state('?')
		.Send("Access_NewRecord", data: .data.Get())
		.applyDefaultNewValues()
		.data.Dirty?(rec.forceDirty is true)
		.FocusFirst(.Vert.Scroll.Hwnd, custom: .customFields)
		.edit_button.Pushed?(true)
		.Send("Access_AfterNewRecord", .data)
		.loopedAddons.Stop()
		}

	newNonDirtyRecord?(sameType?)
		{
		return sameType? and .newrecord? and not .data.Dirty?() and	not .new_setdata?
		}

	// vvvvvvvvvvvvvvvvvvvvvvvvv DEPRECATING: vvvvvvvvvvvvvvvvvvvvvvvvv
	// - defaultNewValues and SetDefaultNewValues is to be removed under: 33976
	defaultNewValues: false
	SetDefaultNewValues(values)
		{
		.new_setdata? = true
		.defaultNewValues = values
		}
	applyDefaultNewValues()
		{
		if .defaultNewValues is false
			return
		data = .GetData()
		loopFields = .defaultNewValues.Member?('fillinSequence')
			? .defaultNewValues.fillinSequence
			: .defaultNewValues.Members()
		.applyEachDefaultNewValue(loopFields, data)
		data.forceDirty = .defaultNewValues.GetDefault('forceDirty', false)
		}
	applyEachDefaultNewValue(loopFields, data)
		{
		for field in loopFields
			{
			if .defaultNewValues[field] is ""
				continue
			data[field] = .defaultNewValues[field]
			if false is fieldCtrl = .GetControl(field)
				if false isnt tabs = .FindControl('Tabs')
					{
					tabs.ConstructAllTabs()
					fieldCtrl = .GetControl(field)
					}
			if fieldCtrl isnt false and fieldCtrl.Base?(KeyControl)
				fieldCtrl.Field.Process_newvalue()
			.Record_NewValue(field, data[field])
			}
		}
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	On_Edit()
		{
		if ReadOnlyAccess(this) is true or .protect or .edit_button.GetEnabled() is false
			return
		.edit_button.SetEnabled(false)
		if .protectField isnt false
			{
			protect = .data.GetField(.protectField)
			if String?(protect) and protect isnt ""
				{
				Alert(protect, 'Reason Protected', .Window.Hwnd, MB.ICONINFORMATION)
				.edit_button.SetEnabled(true)
				return
				}
			}
		.toggleEdit()

		// could be entering or leaving edit mode
		// can get called multiple times
		// e.g. when you try to save an invalid record
		.Send('Access_FinishOnEdit')
		.edit_button.SetEnabled(AccessProtect.AllowEdit?(.data, .protect, .protectField))
		}
	toggleEdit()
		{
		if .data.GetReadOnly() is true
			{
			.edit_mode()
			if not .data.HasFocus?()
				.FocusFirst(.Vert.Scroll.Hwnd, custom: .customFields)
			}
		else
			{
			.Save() // does view_mode
			if .data.GetReadOnly() is false // save failed
				.edit_button.Pushed?(true)
			else
				{
				.new_button.Pushed?(false)
				.loopedAddons.Start()
				}
			}
		}
	view_mode()
		{
		if .linked?
			return
		.data.SetReadOnly(true)
		.edit_button.Pushed?(false)
		.edit_button.SetEnabled(AccessProtect.AllowEdit?(.data, .protect, .protectField))
		.lock.Unlock()
		.kill_valid_timer()
		}
	edit_mode(noReload = false)
		{
		if .linked?
			return
		// reload record in case record has been changed by another user
		// while this user viewing record (must be done BEFORE lock)
		if false is noReload and not .Reload()
			{
			.AlertInfo(.title, "The current record has been deleted.")
			return
			}

		if not .lock.Trylock()
			return

		if .validField isnt false
			.record[.validField]

		.edit_button.Pushed?(true)
		.data.SetReadOnly(false)
		.loopedAddons.Stop()
		}
	GetLockKey()
		{
		return .model.GetLockKey(.record)
		}
	On_First()
		{
		.get(#Next, #On_New, 'first')
		}
	On_Last()
		{
		.get(#Prev, #On_New, 'last')
		}
	On_Next()
		{
		.get(#Next, #On_Last)
		// TODO: don't really need to do On_Last
		// since if Next fails you're already on last
		}
	On_Prev()
		{
		.get(#Prev, #On_First)
		// TODO: don't really need to do On_First
		// since if Prev fails you're already on first
		}
	get(dir, onfail, firstlast = '?')
		{
		if not .Save()
			return
		if firstlast is 'first' or firstlast is 'last'
			.c.Rewind()
		Transaction(read:)
			{|t| x = .c[dir](t) }
		if x is false
			{
			.beep()
			// need to check if table is empty without select
			// in this case, we don't want to give the message
			if .select is true and not .model.TableEmpty?() and .firstRead? is true
				{
				if 'deleted' isnt .CheckDeleted(quiet:)
					.noRecordFound(dir)
				else // current record is deleted and last in select,
					// remove select and go to last
					.recordDeleted()
				return
				}
			this[onfail]()
			}
		else
			{
			.firstRead? = false
			.setdata(x)
			.set_position_button_state(firstlast)
			}
		}
	beep() // so tests can override
		{
		Beep()
		}
	noRecordFound(dir)
		{
		position = dir is 'Prev'
			? 'first'
			: (dir is 'Next' ? 'last' : '?')
		.set_position_button_state(position)
		.Select_vals.Each({ it.check = false })
		Alert('No records found that ' $
			'match the current select.\r\n\r\nSelection has been reset.',
			.title, flags: MB.ICONINFORMATION)
		.SetWhere("")
		.Defer(uniqueID: 'reopen_select_dialog') // need the orig select to close first
			{
			.On_Select()
			}
		}
	recordDeleted()
		{
		.setInitialWhere()
		if not QueryEmpty?(.query)
			.On_Last()
		else
			.On_New()
		}

	set_position_button_state(position)
		{
		if .last_button is false
			return
		.last_button.Grayed(position is 'last')
		.next_button.Grayed(position is 'last')
		.first_button.Grayed(position is 'first')
		.prev_button.Grayed(position is 'first')
		}
	On_Go()
		{
		if (not .HasLocate?() or
			not .Save() or
			(locate = .locate.Get().locate) is "" or
			(by = .locate.Get().locateby) is "" or
			not .model.GetLocateKeys().Member?(by))
			return
		.AccessGoto(.model.GetLocateKey(by), locate)
		.setFocusToLocate()
		}
	AccessGoto(field, value)
		{
		if ((false is field = .model.FindGotoField(field)) or not .Save())
			return

		.set_position_button_state('?')
		x = false
		value = DatadictEncode(field, value)
		// 10281 - handle 'invalid query' exception only, see catch.
		// This is needed for queries that have sort on rule field or
		// non-indexed field (cursor throws exception)
		try
			{
			.c.Seek(field, value)
			Transaction(read:)
				{|t| x = .c.Next(t) }
			if x isnt false and x[field] is value
				{
				.setdata(x)
				.firstRead? = false
				return
				}
			else if .invalidLocate?(value)
				{
				.showTipForInvalidLocate(x, field, value)
				return
				}
			}
		catch (err /*unused*/, "invalid query")
			{ } // fall through
		.gotoWithoutSelect(field, value)
		}
	invalidLocate?(value)
		{
		return .HasLocate?() and value isnt '' and not .locate.Valid?()
		}
	showTipForInvalidLocate(x, field, value)
		{
		if x isnt false and
			String(x[field]).Prefix?(String(value))
			.locateTip("Matches more than one record")
		else
			.locateTip("No matching records")
		}
	gotoWithoutSelect(field, value)
		{
		if false isnt y = .model.LookupRecord(field, value)
			.setdata(y)
		else
			.On_Last()
		}
	locateTip(msg)
		{
		.locate.BalloonTip(msg)
		}
	getter_keyquery()
		{
		return .model.GetKeyQuery()
		}
	On_Current_Save()
		{
		.Save()
		}

	// IMPORTANT: there was an issue with linked browses getting out of sync due
	// to the browse being cleared BEFORE the access had successfully deleted the
	// header record. In order to prevent these types of issues, this method should
	// either:
	// - successfully delete the record and load the next/new one OR
	// - reload the record to ensure consistency and prevent potentially multiple
	//		"bad" states after an attempted delete
	On_Current_Delete(@unused)
		{
		if .notAbleToDelete()
			{
			.failedDelete()
			return
			}
		.data.Dirty?(false)
		.invalid? = false
		errMsg = ''
		.Send('AccessBeforeDeleting')
		curData = .GetData()
		Transaction(update:)
			{|t|
			// errMsg can be:
			// '' 				- no error, delete success
			// false 			- key exception error, has been alerted
			// non-empty string	- other errors
			errMsg = .deleteCurrentRecord(t)
			}
		if errMsg isnt ''
			.failedDelete(errMsg)
		else
			.deleteRecordAttachments(curData)
		}

	deleteRecordAttachments(rec)
		{
		.attachmentsManager.QueueDeleteRecordFiles(rec)
		.deleteOldAttachments()
		}

	failedDelete(msg = 'Delete failed')
		{
		if msg isnt false
			.AlertError('Current Delete', msg)
		if not ReadOnlyAccess(this)
			.Reload()
		}

	RecordConflict?(x, quiet? = false)
		{
		return RecordConflict?(.original_record, x, .fields, .Window.Hwnd, :quiet?)
		}
	notAbleToDelete()
		{
		if .readOnlyAccess() is true or false is .Send("Access_AllowDelete") or
			not .RecordSet?()
			return true
		if not .EditMode?()
			{
			if not .lock.Trylock()
				return true
			// TODO: this should be around the deleting
			// otherwise someone could lock it between our unlock and our delete
			// although they'd have to do it at exactly the wrong time
			.lock.Unlock()
			}
		if false is AccessProtect.AllowDelete?(
			this, .record, .protect, .protectField, .newrecord?)
			return true
		return .model.NotifyObservers("delete") is false
		}
	readOnlyAccess() // for tests
		{
		return ReadOnlyAccess(this)
		}
	deleteCurrentRecord(t)
		{
		errMsg = ''
		KeyException.TryCatch(
			block:
				{
				deleted? = false
				if not .newrecord?
					{
					if ((x = .CheckDeleted(t, quiet:)) is 'deleted')
						deleted? = true
					else if .RecordConflict?(x)
						errMsg = 'Another user has modified this record'
					}
				if errMsg is ""
					errMsg = .deleteCurrentRecordFn(t, deleted?)
				}
			catch_block:
				{|e|
				errMsg = .getDeleteErrMsg(e)
				if not t.Ended?()
					t.Rollback()
				}
			)
		return errMsg
		}
	deleteCurrentRecordFn(t, deleted?)
		{
		errMsg = ''
		.Send('AccessBeforeDelete', :t)
		.model.Before_Delete(t, x: .GetData())
		if .newrecord? or deleted?
			{
			t.Complete()
			.new_setdata? = true
			if .start_last?
				.On_Last()
			else
				.On_First()
			.new_setdata? = false
			}
		else
			{
			if .saveOnlyLinked isnt true and
				0 is t.QueryDo("delete " $ .keyquery)
				errMsg = "Delete failed"
			.Send('AccessAfterDelete', :t)
			t.Complete()
			.Send('AccessAfterDeleting')
			.On_Next()
			}
		return errMsg
		}
	getDeleteErrMsg(e)
		{
		if e.Has?('blocked by foreign key') and
			"" isnt (msg = .showForeignKeyUsage())
			errMsg = 'This record can not be deleted. ' $ msg
		else
			{
			errMsg = false // KeyException generates its own alert
			KeyException(e, action: 'delete')
			}
		return errMsg
		}

	showForeignKeyUsage()
		{
		return .model.ForeignKeyUsage(.data.Get())
		}

	On_Current_Restore()
		{
		if .record is false or false is .Send("Access_AllowRestore")
			return
		.model.NotifyObservers("restore")
		.RestoreAttachmentFiles()
		if .newrecord?
			{
			// have to simulate New button click.  If dynamic types are
			// used then we must translate " " to "_" in the menu choice
			// in order to simulate the MenuButton click
			.data.Dirty?(false)
			.new_setdata? = true
			.On_New(.types.GetTypeName())
			.new_setdata? = false
			.Send("Access_Restore", newrecord: true)
			.kill_valid_timer()
			return
			}
		if not .Reload()
			{
			Alert("The current record has been deleted.", title: 'Current Restore',
				flags: MB.ICONERROR)
			return
			}
		.kill_valid_timer()
		.data.Dirty?(false)
		.data.SetAllValid()
		.Send("Access_Restore", newrecord: false)
		}
	Reload(forceViewMode = false)
		{
		if .newrecord?
			return true
		edit? = .EditMode?()
		x = Query1(.keyquery)
		if x isnt false
			{
			.setdata(x)
			if edit? and not forceViewMode
				.edit_mode(noReload:)
			}
		return x isnt false
		}

	On_Current_Print()
		{
		if not Suneido.GetDefault('user_roles', #()).Has?('admin')
			{
			.AlertInfo('Current Print',
				'You must be in the admin role to use Current > Print')
			return
			}

		if .record is false
			return
		CurrentPrint(.record, .Window.Hwnd, .base_query, Display(.Parent)[.. -2],
			excludeFields: .GetExcludeSelectFields())
		}
	PrintReport(reportClass, data = false, permission = false, additionalArgs = false)
		{
		if permission isnt false and AccessPermissions(permission) isnt true
			{
			.AlertInfo(.title, 'You do not have permission for this option')
			return
			}

		if not .SaveFor('Print ' $ reportClass.AccessTitle $ ' on')
			return

		args = Object(data is false ? .GetData() : data)
		if additionalArgs isnt false
			args.Merge(additionalArgs)
		ToolDialog(.Window.Hwnd, reportClass(@args))
		.Reload()
		}

	On_Current_Reason_Protected()
		{
		AccessProtect.ReasonProtected(this, .protectField, .showForeignKeyUsage())
		}

	On_Current_View_Warnings()
		{
		if '' isnt warnings = .checkFieldRule(true, .warningField)
			.AlertWarn(.Title $ ' Warnings', warnings.Replace('; ', ';\r\n'))
		else
			.AlertInfo(.Title, 'No warnings provided')
		}

	set_select_button_state()
		{
		if .select_button isnt false
			.select_button.Pushed?(.select)
		}
	On_Select()
		{
		if not .Save()
			return
		BookLog('Access Select Start')
		.set_select_button_state()
		SelectControl(
			this, .selectMgr.Name(), defaultButton: .start_last? ? 'Last' : 'First'
			noUserDefaultSelects?: .accessGoTo)
		BookLog('Access Select End')
		}

	firstRead?: false
	SetWhere(where, quiet = false, hwnd = false, extraMsg = '') // called by Select
		{
		if hwnd is false
			hwnd = .Window.Hwnd
		.select = where > ""
		.Defer(.set_select_button_state, uniqueID: 'select_button_state')

		preQuery = .GetQuery()
		ret = .model.AddMoreToQuery(where)
		if preQuery isnt .GetQuery()
			.types.ClearStickyFieldValues()
		.SetDefaultStatus()
		if String?(ret) and not quiet
			Alert(ret $ extraMsg, title: 'Select', :hwnd, flags: MB.ICONINFORMATION)
		.firstRead? = true
		return ret is true
		}
	ModifyWhere(where, hwnd)
		{
		k = .getKey()
		if not .SetWhere(where, :hwnd, extraMsg: '\r\n\r\nSelect will be cleared')
			{
			.select = false
			.Select_vals.Each({ it.check = false })
			}
		.firstRead? = true
		.AccessGoto(@k)
		if .firstRead? is true // no records matched
			{
			.noRecordFound('?')
			return
			}
		// so that the Select... button will give up it's blue rectangle correctly when
		// closing the select dialog
		.Defer(uniqueID: 'select_button_state')
			{
			SetFocus(.select_button.Hwnd)
			.set_select_button_state()
			}
		}
	getKey()
		{
		field = .model.GetKeyField()
		// apm - what if field is composite?
		return Object(field, .data.GetField(field))
		}
	getter_c()
		{
		return .model.GetCursor()
		}
	ChangeQuery(query)
		{
		.model.SetQuery(query)
		if .start_last?
			.On_Last()
		else
			.On_First()
		}
	getter_query()
		{
		return .model.GetQuery()
		}

	invalid?: false
	valid_timer: false
	record_change_members: false
	Access_RecordChange(member)
		{
		if .valid_timer is false
			{
			.record_change_members = Object(member)
			.valid_timer = Defer(.record_change)
			}
		else if Object?(.record_change_members)
			.record_change_members.AddUnique(member)
		else
			.record_change_members = Object(member)

		.model.Plugins_Execute(access: this, :member, pluginType: 'AccessObservers')
		data = .data.Get(excludeHandleFocus:)
		.model.Plugins_Execute(:data, :member, hwnd: .Window.Hwnd, query: .keyquery,
			pluginType: 'Observers')
		}
	record_change()
		{
		// handle case where AccessControl destroyed and Delayed function still called
		if .record_change_members is false
			{
			.kill_valid_timer()
			return
			}
		if .invalid?
			.check_valid()
		else
			.setWarnings()
		.kill_valid_timer()
		.Send("Access_RecordChange", .record_change_members)
		.record_change_members = false
		}

	Record_NewValue(field, value)
		{
		if .fields.Has?(field)
			{
			.model.Plugins_Execute(data: .data is false ? false : .data.Get(),
				:field,	hwnd: .Window.Hwnd, query: .keyquery, pluginType: 'AfterField')
			.Send("Access_AfterField", field, value)
			}
		}

	allowCustomTabs?: false
	Tabs_BeforeConstruct(controls)
		{
		if false is idx = controls.FindIf({ it.Tab is 'Custom' })
			return
		.allowCustomTabs? = true
		if false is cl = OptContribution('CustomTabPermissions', false)
			return
		if not String?(tableName = .Send("GetCustomizableName"))
			tableName = .model.GetTableName()

		controls.Delete(idx)
		cl.WithPermissableTabs(tableName)
			{ |tab|
			controls.Add(Object('Customizable' Tab: tab, tabName: tab))
			}
		}

	AllowCustomTabs?()
		{
		return .allowCustomTabs?
		}

	Valid?()
		{
		return .check_valid()
		}
	check_valid(evalRule? = false)
		{
		if not .data.Dirty?()
			return true

		// Use RecordControl's Valid to poll controls
		if ((invalid_fields = .data.Valid()) isnt true)
			return .setInvalidFields(invalid_fields)

		if "" isnt errCustom = CustomizeField.CheckCustomFields(
			.customFields, .data, .protectField)
			return .setInvalidFields(errCustom)

		//check valid field to determine if record is valid.
		if '' isnt stat = .checkFieldRule(evalRule?, .validField)
			return .setInvalidFields(stat)

		.invalid? = false
		.setWarnings()
		return true
		}

	checkFieldRule(evalRule?, field)
		{
		if field is false
			return ''
		return evalRule? is true
			? .data.Get().Eval(Global('Rule_' $ field))
			: .data.GetField(field)
		}

	setInvalidFields(msg)
		{
		if not .invalid?
			{
			.invalid? = true
			.beep()
			}
		.invalid_status(msg)
		return false
		}

	invalid_status(msg)
		{ .status.Set(msg, invalid:) }

	setWarnings()
		{
		stat = .checkFieldRule(true, .warningField)
		.status.Set(stat, normal: stat is '', warn: stat isnt '')
		}

	On_Save()
		{ .Save() }

	SaveForAndToggleEdit(action, trackActionResult = false)
		{
		addonVals = Object()
		.Addons.Send('Before_SaveForAndToggleEdit', addonVals)
		if false is .SaveFor(action)
			return trackActionResult ? #(save: false, edit: false) : false

		if false is .SetEditMode()
			return trackActionResult ? #(save:, edit: false) : false

		.Addons.Send('After_SaveForAndToggleEdit', addonVals)
		return trackActionResult ? #(save:, edit:) : true
		}

	SaveFor(action)
		{
		if .NewRecord?() and not .data.Dirty?()
			{
			.AlertInfo(action.Replace(' (on|to)$', ''),
				"Can't " $ action $ " an empty new record")
			return false
			}
		return .Save()
		}

	Save()
		{
		if ReadOnlyAccess(this)
			return true

		if false is .Send('Access_SavePreCheck')
			return false
		// if no record has loaded yet, we don't want to be calling view_mode
		// which calls protect rule which can potentially trigger next num rules
		// Also, no need to leave edit mode if on new record
		// and user hasn't done anything
		if (.record is false or (.newrecord? and not .data.Dirty?()))
			return true
		if true is result = .save()
			{
			.view_mode()
			.nextNumber.Confirm()
			}
		return result
		}

	lastSaveTime: false
	save()
		{
		if not .data.Dirty?()
			return true

		.record_change_without_delay()

		if false is .check_valid(evalRule?:)
			{
			.model.NotifyObservers('accessInvalid')
			return false
			}
		if true isnt msg = .nextNumber.Renew() // in case reservation has expired
			{
			.AlertWarn('Next Number', msg)
			return false
			}

		// update orig and new record members in temp members in case
		// an error occurs during the save, only update the real members
		// at the end of this method when the transaction has completed
		.temp_original_record = 0
		.temp_newrecord? = .newrecord?

		if not .EditMode?()
			{
			if false isnt debug = Suneido.GetDefault("RecordNotInEditMode", false)
				SuneidoLog('INFO: Record Not In Edit Mode',
					params: [firstChangeMember: debug.member],
					calls: debug.callstack)
			ProgrammerError("Access dirty, not in edit mode (" $ .title $ ")")
			Suneido.Delete(#RecordNotInEditMode)
			}

		if false is .Send('AccessBeforeSaving')
			return false
		if false is .saveWithKeyExceptionTran()
			return false

		return .afterSave()
		}

	saveWithKeyExceptionTran()
		{
		errMsgOb = Object()
		try
			return KeyExceptionTransaction()
				{|t|
				data = .GetData()
				.menus.UpdateHistory(data, .newrecord?)
				result = .Send('AccessBeforeSave', :t)
				.Addons.Collect('SaveValid', data, .GetOriginal(), t).Each()
					{ if it > "" errMsgOb.Add(it) }
				if result isnt false
					result = .do_save(t)
				if result is false or result is 'deleted' or errMsgOb.NotEmpty?()
					{
					if not t.Ended?()
						t.Rollback()
					if errMsgOb.NotEmpty?()
						.AlertError('', errMsgOb.Join('\n'))
					return false
					}
				if result is true
					.Send('AccessAfterSave', :t)
				if t.Ended?()
					return false
				.lastSaveTime = Date()
				}
		catch (unused, "interrupt: KeyException")
			return false
		}
	do_save(t)
		{
		result = true
		if .saveOnlyLinked isnt true
			result = .save_record(t)
		else if not .newrecord? // still check if record has been deleted
			result = Record?(.CheckDeleted(t)) ? true : 'deleted'
		else if .saveOnlyLinked is true and .newrecord?
			.temp_newrecord? = false
		return .notifyAfterSave(result, t)
		}
	save_record(t)
		{
		result = true
		if .newrecord?
			.save_output(t)
		else
			result = .save_update(t)
		return result
		}
	notifyAfterSave(result, t)
		{
		if result is true
			{
			if false is .model.NotifyObservers('save', t)
				result = false
			else if false is (.temp_original_record = t.Query1(.keyquery))
				SuneidoLog("ERROR: AccessControl.do_save - can't get original record:" $
					"\r\n" $ "keyquery: " $ .keyquery, calls:)
			}
		if result is true
			.temp_newrecord? = false
		return result
		}
	afterSave()
		{
		// use temp original record because we only want to set original record
		// if save successfully makes it through the transaction
		if .temp_original_record isnt 0
			.original_record = .temp_original_record
		.newrecord? = .temp_newrecord?
		.model.NotifyObservers('after_save')
		.data.Dirty?(false)
		.Send('AccessAfterSaving')
		.deleteOldAttachments()
		return true
		}
	deleteOldAttachments()
		{
		.attachmentsManager.ProcessQueue()
		}
	RestoreAttachmentFiles()
		{
		.attachmentsManager.ProcessQueue(restore?:)
		for browse in .linkedBrowses
			if Instance?(browse) and not browse.Destroyed?() and
				browse.Method?('RestoreAttachmentFiles')
				browse.RestoreAttachmentFiles()
		}
	QueueDeleteAttachmentFile(newFile, oldFile, name, action)
		{
		return .attachmentsManager.QueueDeleteFile(
			newFile, oldFile, .GetData(), name, action)
		}
	linkedBrowses: ()
	RegisterLinkedBrowse(browseCtrl, name)
		{
		if .linkedBrowses.Readonly?()
			.linkedBrowses = Object()
		.linkedBrowses[name] = browseCtrl
		}
	Access_Save()
		{
		return .Save()
		}
	// returns 'deleted' if deleted, otherwise the record
	CheckDeleted(t = false, quiet = false)
		{
		DoWithTran(t)
			{|tran|
			x = tran.Query1(.keyquery)
			}

		if x is false
			{
			if not quiet
				Alert("Access: can't get record to update", title: 'Error',
					flags: MB.ICONERROR)
			// return deleted so that the user can get off the screen but the
			// access does not try to save any browses or explorer list views
			return 'deleted'
			}
		return x
		}
	save_output(t)
		{
		t.QueryOutput(.base_query, .record)
		.types.SetStickyValues(.record)
		.model.SetKeyQuery(.record)
		}
	save_update(t)
		{
		if ((x = .CheckDeleted(t)) is 'deleted')
			return 'deleted'
		// NOTE: original_record should only be false
		// if the lookup failed earlier (an assertion failure)
		if .original_record is false
			SuneidoLog("AccessControl.save_update - original_record is false")
		else if .RecordConflict?(x)
			return false
		x.Update(.record)
		.model.SetKeyQuery(.record)
		return true
		}
	Status(status)
		{
		if .status.GetValid()
			.status.Set(status)
		}
	GetFields()
		{
		return .fields
		}
	GetKeys()
		{
		return .model.GetKeys()
		}
	GetExcludeSelectFields()
		{
		return .model.GetExcludeSelectFields()
		}
	NewRecord?()
		{
		return .newrecord?
		}
	GetOriginal() // used by KeyFieldControl
		{
		return .original_record
		}
	GetQuery()
		{
		return .query
		}
	GetKeyQuery()
		{
		return .keyquery
		}
	GetTransQuery()
		{
		return .query
		}
	getter_base_query()
		{
		return .model.GetBaseQuery()
		}
	GetBaseQuery()
		{
		return .base_query
		}
	Locate?(available = true)
		{
		.locate_status = available
		}
	HasLocate?()
		{
		return .locate_status isnt false
		}
	GetData()
		{
		return .data.Get()
		}
	GetCurrentSelectedData()
		{
		return .GetData()
		}
	GetControl(field)
		{
		return .data.GetControl(field)
		}
	GetRecordControl()
		{
		return .data
		}
	AccessObserver(fn, at = false)
		{
		.model.AddObserver(fn, at)
		}
	RemoveAccessObserver(fn)
		{
		.model.RemoveObserver(fn)
		}

	SetVisible(visible?)
		{
		super.SetVisible(visible?)
		.loopedAddons.Subscriber(enabled?: visible?)
		}

	SetReadOnly(readOnly)
		// pre:		readOnly is a Boolean value
		// post:	this is readOnly iff readOnly is true
		{
		Assert(Boolean?(readOnly), "Access SetReadOnly: not boolean")
		.protect = readOnly
		.data.SetReadOnly(readOnly)
		}

	SetSelectMgr(selectName, args)
		{
		// set initial select if any
		.selectMgr = AccessSelectMgr(args.GetDefault('select', #()),
			name: selectName)
		if .select_button isnt false
			.selectMgr.LoadSelects(this)
		}
	Getter_Select_vals()
		{
		return .selectMgr.Select_vals()
		}
	SetSelectVals(select_vals)
		{
		.selectMgr.SetSelectVals(select_vals, .GetSelectFields())
		}
	SetSelect(selects)
		{
		.selectMgr.Reset(selects)
		.setInitialWhere(fromNew?:)
		.Send('SelectControl_Changed')
		}
	GetSelectFields()
		{
		return .sf
		}
	// so MultiView > Access > embedded/linked virtual list use its own select mgr name
	OverrideSelectManager?()
		{
		return false
		}

	On_Global(option)
		{
		BookLog("Access Global " $ option)
		.menus.On_Global(option, .Window.Hwnd)
		}

	On_Current(option)
		{
		if option is 'Delete'
			return
		BookLog("Access Current " $ option)
		.menus.On_Current(option, .GetData(), this)
		}

	On_Locate()
		{
		if .HasLocate?()
			SetFocus(.locate.EditHwnd())
		}

	Dirty?(dirty = '')
		{
		.data.Dirty?(dirty)
		}
	EditMode?()
		{
		return not .data.GetReadOnly()
		}

	SetEditMode()
		{
		if not .EditMode?()
			.On_Edit()
		return .EditMode?()
		}
	SetMainRecordField(field, value)
		{
		.data.SetField(field, value)
		}
	ConfirmDestroy()
		{
		if ReadOnlyAccess(this)
			return true

		.Send('Access_ConfirmDestroy')

		// this will be used in CloseWindowConfirmation
		// to log if the changes are discarded
		Suneido.AccessRecordDestroyed = .getKey()

		// This must be done before valid checking
		.record_change_without_delay(destroying:)

		if false is .check_valid()
			.model.NotifyObservers('accessInvalid')

		// Save was being done in Destroy.  Moved here so user could fix
		// duplicate keys before the Destroy is done.
		return ((not .data.Dirty?()) or (.invalid? is false)) and .Save()
		}
	record_change_without_delay(destroying = false)
		{
		// make sure the current field gets into record_change_members
		hwnd = .getFocus()
		.ClearFocus()
		// if destroying don't need to keep focus, and it causes problems
		if not destroying
			.setFocus(hwnd)

		// clear valid timer first to ensure delayed record_change does not get done
		.kill_valid_timer()

		// because of delay in observer, make sure the recordchange message is sent
		if .record_change_members isnt false
			{
			.Send("Access_RecordChange", .record_change_members)
			.record_change_members = false
			}
		}
	getFocus()
		{
		return GetFocus()
		}
	setFocus(hwnd)
		{
		SetFocus(hwnd)
		}
	old_accels: false
	kill_valid_timer()
		{
		if .valid_timer isnt false
			{
			.valid_timer.Kill()
			.valid_timer = false
			}
		}

	CloseWindowConfirmation()
		{
		if Boolean?(result = .Send("Access_CloseWindowConfirmation"))
			return result
		return true
		}

	Destroy()
		{
		.loopedAddons.Subscriber()
		Suneido.AccessRecordDestroyed  = ''
		.lock.Unlock()
		.c.Close()
		if .select_button isnt false
			.selectMgr.SaveSelects()
		.Window.RemoveValidationItem(this)
		.kill_valid_timer()
		.nextNumber.PutBack()
		super.Destroy()
		}
	}