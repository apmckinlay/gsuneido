// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xstretch: 1
	list: false
	New(.query, .columns, .linkField, .validField = false,
		.protectField = false, .headerFields = #(), .dataMember = "lineitem_data",
		.expandLayout = false, .extraMenu = #(), .extraFmts = false,
		.hdrCornerCtrl = false, .keyField = false, .stretchColumn = false,
		.customDelete = #(), .preventCustomExpand? = false, .switchToForm = false)
		{
		super(.createControls())
		.allDataMember = 'all_' $ .dataMember
		.list = .FindControl('VirtualList_LineItem')
		.list.Commands = #()
		.Send('Data')
		.Send("AccessObserver", .AccessChanged)
		.Send("AddSetObserver", .recordSetHandler)
		.setButtonRedirects()
		}

	createControls()
		{
		return Object('VirtualList',
			validField: .validField,
			protectField: .protectField,
			menu: .createMenu(),
			preventCustomExpand?: .preventCustomExpand?,
			enableDeleteBar:,
			disableSelectFilter:,
			loadAll?:,
			extraFmts: .extraFmts,
			hdrCornerCtrl: .hdrCornerCtrl,
			stretchColumn: .stretchColumn,
			linked?:,
			switchToForm: .switchToForm
			name: 'VirtualList_LineItem')
		}

	createMenu()
		{
		menu = Object()
		if not .customDelete.Empty?()
			menu.Add(@.customDelete)
		else
			menu.Add('Delete/Undelete')

		menu.Add('', 'Reason Protected')
		if not .extraMenu.Empty?()
			menu.Add('')
		menu.Add(@.extraMenu)
		return menu
		}

	setButtonRedirects()
		{
		buttonCmds = Object()
		if .expandLayout isnt false
			.findButtons(.expandLayout.ctrl, buttonCmds)
		if .hdrCornerCtrl isnt false
			.findButtons(.hdrCornerCtrl, buttonCmds)
		for cmd in buttonCmds
			.list.Redir(cmd, this)
		}

	findButtons(layout, cmds)
		{
		if String?(layout[0]) and layout[0].Suffix?("Button")
			{
			cmds.Add("On_" $ ToIdentifier(layout.GetDefault('command', { layout[1] })))
			return
			}
		for m in layout
			if Object?(m)
				.findButtons(m, cmds)
		}

	VirtualList_NeedValidation?()
		{
		return false // AccessControl handles all control validation
		}

	Recv(@args)
		{
		if args[0].Prefix?('On_')
			{
			args = args.Copy()
			.list.SelectExpandRecord(args.source)
			args.rec = .list.GetSelectedRecord()
			.Send(@args)
			}
		return 0
		}

	AccessChanged(@args)
		{
		// the control is destroyed in the middle of TabsControl before_setdata
		if .list is false or .list.GetModel() is false
			return false

		args = args.Copy()
		args.list = this
		args.grid = .list.GetGrid()
		args.model = .list.GetModel()
		return .excludeForSave(args)
			{
			LineItemSave(@args)
			}
		}

	excludeForSave(args, block)
		{
		if args[0] isnt 'save' or .exclude is false
			return block()

		for item in .list.GetLoadedData()
			if ((.exclude)(item))
				{
				if Object?(item.vl_origin)
					item.vl_deleted = true // temporarily set flag
				item.vl_excluded = true
				}
		result = block()
		for item in .list.GetLoadedData() // restore deleted flag
			if item.vl_excluded is true
				{
				item.vl_deleted = false
				item.Delete('vl_excluded')
				}
		return result
		}

	getter_exclude()
		{
		if 0 is .exclude = .Send("LineItem_ExcludeItemFn") // once only
			.exclude = false
		return .exclude
		}

	GetAllLineItems(includeAll? = false)
		{
		return .list.GetLoadedData().Filter({
			it.vl_expand? isnt true and .filterItem(it, includeAll?) })
		}

	filterItem(item, includeAll?)
		{
		if includeAll? is true or .exclude is false
			return true
		return not (.exclude)(item)
		}

	GetLineItems()
		{
		return .GetAllLineItems().Filter({ it.vl_deleted isnt true })
		}

	GetDeleted()
		{
		return .GetAllLineItems().Filter({ it.vl_deleted is true })
		}

	AddRecord(rec, pos = 'current', preset? = false)
		{
		.list.AddRecord(rec, pos)
		.Send("SetField", .allDataMember, .GetAllLineItems(includeAll?:))
		if not preset?
			{
			.Send("SetField", .dataMember, .GetLineItems())
			.Send("InvalidateFields", Object(.dataMember, .allDataMember))
			}
		}

	DeleteRow(rec)
		{
		.On_Context_DeleteUndelete(rec, force:)
		}

	VirtualList_RepaintDisabled()
		{
		return .Loading?()
		}

	Loading?()
		{
		return .setting?
		}

	GetExpanded()
		{
		return .list.GetModel().ExpandModel.GetExpanded()
		}

	GetExpandedControl(rec)
		{
		return .list.GetExpandedControl(rec)
		}

	linkValue: false
	headerData: false
	setting?: false
	Set(value)
		{
		.setting? = true
		state = .linkValue is value ? .getState() : false
		if .list.GetModel() isnt false
			.list.CollapseAll()

		.linkValue = value
		if not .headerFields.Empty?()
			{
			.headerData = .Send("GetData")
			.headerData.Observer(.observer_HeaderData)
			}
		query = QueryAddWhere(.query,
			" where " $ .linkField $ " is " $ Display(.linkValue))

		accessCustomKey = .Send('GetAccessCustomKey')
		customKey = .buildCustomKey(accessCustomKey, .Name)

		.list.SetQuery(query, .columns, :customKey)
		.Send('DoWithoutDirty',
			{
			.Send("SetField", .dataMember, .GetAllLineItems())
			.Send("SetField", .allDataMember, .GetAllLineItems(includeAll?:))
			.Send('LineItem_SetWithoutDirty', .headerData, :state)
			})
		.setState(state)
		.Send('LineItem_AfterSet')
		.setting? = false
		}

	recordSetHandler()
		{
		if .Name is ''
			return

		// the following prevents expanding line-item records when this ctrl isnt on
		// the current tab (for performance reasons)
		curTabCtrl = .Send('TabsControl_GetCurrentControl')
		if curTabCtrl is 0 or false is curTabCtrl.FindControl(.Name)
			return

		.Send('DoWithoutDirty',
			{
			if 0 isnt types = .Send('GetSavedExpandedTypes')
				.ExpandByField(@types.Merge(#(keepPos?:)))
			})
		}

	buildCustomKey(accessCustomKey, ctrlName)
		{
		return String?(accessCustomKey)
			? accessCustomKey $ ' | ' $ ctrlName
			: false
		}

	getState()
		{
		state = Object(expanded: #(), offset: 0, selected: false)
		if false is model = .list.GetModel()
			return state

		if false isnt rec = .list.GetSelectedRecord()
			state.selected = rec[.keyField]
		state.offset = model.Offset
		if model.ExpandModel isnt false
			state.expanded = .GetExpanded().Map({ it[.keyField] }).Instantiate()
		return state
		}

	setState(state)
		{
		if state is false
			return
		if state.selected isnt false
			.list.SelectRecordByKeyPair(state.selected, .keyField)
		.list.ExpandByField(state.expanded, .keyField, keepPos?:)
		.list.VertScroll(state.offset - .list.GetModel().Offset)
		}

	VirtualList_AutoSave?()
		{
		return false
		}

	Get() { return .linkValue }

	observer_HeaderData(member)
		{
		if .Destroyed?()
			return
		if .changeHeaderData(member)
			.list.RepaintGrid()
		}

	changeHeaderData(member)
		{
		if not .headerFields.Has?(member)
			return false
		changed = false
		for row in .GetAllLineItems(includeAll?:)
			{
			if row[member] isnt .headerData[member]
				{
				changed = true
				row[member] = .headerData[member]
				}
			}
		return changed
		}

	VirtualList_ExtraSetupRecordFn()
		{
		return .handleBeforeRecord
		}

	handleBeforeRecord(record)
		{
		if .headerData is false
			return
		if record.New?()
			record[.linkField] = .linkValue
		// reasons for not using PreSet:
		// 	a.  observer is added after this
		//  b.  rules that have header field deps won't kick in (including protect rule)
		//  	all depends on when protect rule gets referenced, if
		// 		CustomizeField.SetFormulas is called it will reference protect rule
		//		so after header fields are set protect field does not get kicked in
		//		cause deps aren't seen as changing
		// 		eg.  Invoiced Order, pickdrop_type is not unprotected
		for field in .headerFields
			record[field] = .headerData[field]

		.Send('LineItem_BeforeRecord', record, headerData: .headerData)
		}

	VirtualList_NewRowAdded(rec)
		{
		.Send('LineItem_NewRowAdded', rec)
		if rec.Member?('CustomizableSetDefaultValues')
			{
			.Send("SetField", .dataMember, .GetLineItems())
			.Send("SetField", .allDataMember, .GetAllLineItems(includeAll?:))
			.Send("InvalidateFields", Object(.dataMember, .allDataMember))
			}
		}

	VirtualList_InvalidDataChanged(rec)
		{
		rec.lineitem_dirty? = true
		.Send("SetField", .dataMember, .GetLineItems())
		.Send("SetField", .allDataMember, .GetAllLineItems(includeAll?:))
		.Send("InvalidateFields", Object(.dataMember, .allDataMember))
		}

	VirtualList_RecordChange(member, record)
		{
		.setFieldIfChanged(member, record)

		.Send('LineItem_RecordChange', member, record)

		// there could be gui operation in previous lines, especailly suneido.js
		// through internet, needs to check the current control is destroyed again
		if .Destroyed?()
			return

		.Send("InvalidateFields", Object(.dataMember, .allDataMember))

		// the above Sends can make other header fields change
		// but the observer will NOT be called recursively
		// only trigger "SetField" once per row
		if record.lineitem_dirty? isnt true
			.setFieldIfChanged(member, record)
		}

	setFieldIfChanged(member, record)
		{
		if (.list.GetModel().Columns().Has?(member) and
			(not Record?(record.vl_origin) or
			record.vl_origin[member] isnt record[member]))
			{
			record.lineitem_dirty? = true
			.Send("SetField", .dataMember, .GetLineItems())
			.Send("SetField", .allDataMember, .GetAllLineItems(includeAll?:))
			}
		}

	VirtualList_AfterField(field, value, record)
		{
		record.lineitem_dirty? = true
		.Send('LineItem_AfterField', field, value, record)
		.Send("SetField", .dataMember, .GetLineItems())
		.Send("SetField", .allDataMember, .GetAllLineItems(includeAll?:))
		.Send("InvalidateFields", Object(.dataMember, .allDataMember))
		}

	VirtualList_RecordMenu?()
		{
		return false
		}

	VirtualList_GetSubTitle()
		{
		return ListCustomize.SubTitle(this, linkField: .linkField)
		}

	On_Context_DeleteUndelete(rec, force = false)
		{
		if rec is false or .GetReadOnly() is true
			return false

		if not force and false is ProtectRuleAllowsDelete?(rec, .protectField)
			return false

		if not .appAllowDelete(rec)
			return false

		if rec.Member?('vl_origin')
			{
			rec.vl_deleted = rec.vl_deleted isnt true
			rowIndex = .list.GetModel().GetRecordRowNum(rec)
			.list.VirtualListThumb_Expand(rowIndex, expand: false)
			}
		else
			.list.DeleteRow(rec)
		.list.AfterDelete()

		.Send("LineItem_DeleteRecord", rec)
		data = .GetLineItems()
		.Send("SetField", .dataMember, data)
		.Send("SetField", .allDataMember, .GetAllLineItems(includeAll?:))
		.Send("InvalidateFields", Object(.dataMember, .allDataMember))
		return true
		}

	appAllowDelete(rec)
		{
		// always allow "undelete"
		if rec.vl_deleted isnt true
			{
			result = .Send("LineItem_AllowDelete", rec, tranFromSave: false)
			if String?(result)
				{
				.AlertInfo('Delete Record', result)
				return false
				}
			if result is false
				return false
			}
		return true
		}

	LineItemsDirty?()
		{
		if .list.GetModel().EditModel.HasChanges?()
			return true

		for rec in .GetAllLineItems()
			if rec.vl_deleted is true
				return true

		return false
		}

	Valid?()
		{
		if false is .Send("LineItem_ExtraValid")
			return false

		return .AccessChanged('valid')
		}

	VirtualList_SwitchToForm()
		{
		.Send("LineItem_SwitchToForm")
		}

	GetSelectedRecord()
		{
		return .list.GetSelectedRecord()
		}

	GetCustomFields()
		{
		return .list.GetCustomFields()
		}

	VirtualList_Expand(unused)
		{
		return .expandLayout is false ? 0 : .expandLayout.Copy()
		}

	VirtualList_DisableSort?()
		{
		return true
		}

	VirtualList_BeforeSave_PreTran()
		{
		return false
		}

	VirtualList_SaveOutstandingChanges?()
		{
		return false
		}

	VirtualList_AllowNextRowWithoutSave()
		{
		return true
		}

	VirtualList_AllowMove(rec)
		{
		return false isnt .Send('LineItem_AllowMove', :rec)
		}

	VirtualList_Move()
		{
		return .Send('LineItem_Move')
		}

	VirtualList_ShowEditButton?()
		{
		return false
		}

	VirtualList_ReadOnly?()
		{
		return .GetReadOnly()
		}

	VirtualList_SelectRowOnExpandFocus()
		{
		return true
		}

	VirtualList_AddGlobalMenu?()
		{
		return false
		}

	Repaint()
		{
		.list.Repaint()
		}

	ExpandByField(@args)
		{
		.list.ExpandByField(@args)
		}

	VirtualList_AfterExpand(rec, ctrl)
		{
		// Idealy code in the handler for this should save the values
		// returned by GetSavedExpandedKeys
		.Send('LineItem_AfterExpand', :rec, :ctrl)
		}

	ScrollToTop()
		{
		.list.On_VirtualListThumb_ArrowHome()
		}

	SelectRecord(rec)
		{
		.list.SelectRecord(rec)
		}

	SelectRecordByFocus(focus)
		{
		.list.SelectRecordByFocus(focus)
		}

	VirtualList_ItemSelected(rec)
		{
		.Send('LineItem_ItemSelected', rec)
		}

	VirtualList_DefaultExpandLayout()
		{
		.Send('LineItem_DefaultExpandLayout')
		}

	VirtualList_MouseWheel(wParam)
		{
		.WndProc.Callsuper(.GetGridHwnd(), WM.MOUSEWHEEL, :wParam, lParam: 0)
		}

	VirtualList_AfterSave(data, t)
		{
		.Send('LineItem_AfterSave', data, t)
		}

	GetHdrCornerControl()
		{
		return .list.GetHdrCornerControl()
		}

	GetGridHwnd()
		{
		return .list.GetGridHwnd()
		}

	SetHighlighted(items, clr)
		{
		.list.VirtualListHighlightRecords(items, clr)
		}

	RegisterLinkedBrowse(ctrl, name)
		{
		.Send('RegisterLinkedBrowse', ctrl, name)
		}

	Destroy()
		{
		.Send("RemoveAccessObserver", .AccessChanged)
		.Send("RemoveSetObserver", .recordSetHandler)
		.Send('NoData')
		super.Destroy()
		}
	}