// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// NOTE: see VirtualListMethodsMap for more public methods
Controller
	{
	Ymin:			0
	Xmin:			0
	Ystretch:		1
	Xstretch:		1
	Name:			"VirtualListView"
	recMenu: 		false

	New(menu = false, .headerSelectPrompt = false, headerMenu = false, readonly = false,
		filterBy = false, .checkBoxColumn = false, .checkBoxAmountField = false,
		.disableSelectFilter = false, .thinBorder = false,
		.protectField = false, .validField = false, .title = false,
		option = false, historyFields = false, titleLeftCtrl = false,
		.hdrCornerCtrl = '', .expandExcludeFields = #(), addons = #(),
		.filtersOnTop = false, .select = #(), .enableDeleteBar = false, .linked? = false,
		.preventCustomExpand? = false, .switchToForm = false, .excludeCustomize? = false)
		{
		super(.layout(addons, filterBy, titleLeftCtrl))

		.header = .FindControl('VirtualListHeader')
		.grid = .FindControl('VirtualListGrid')
		.thumb = .FindControl('VirtualListThumbBar')
		.scroll = .FindControl('VirtualListScroll')

		addCurrentMenu? = .Send('VirtualList_RecordMenu?') isnt false
		if addCurrentMenu? or .addGlobalMenu?
			.recMenu = RecordMenuManager(protectField, option, historyFields, this)
		.contextMenu = VirtualListContextMenu(
			menu, headerMenu, .recMenu, addCurrentMenu?, .addGlobalMenu?)
		.expandBar = .FindControl('VirtualListExpandBar')
		.expandBtns = .FindControl('VirtualListExpandButtons')
		.statusBar = .FindControl('Status')
		.readonly = readonly
		if .readonly or ReadOnlyAccess(this)
			.SetReadOnly(true)
		}

	Startup()
		{
		VirtualListViewExtra.Startup(this)
		}

	getter_topFilters()
		{
		return .FindControl('SelectRepeat')
		}

	getter_addGlobalMenu?()
		{
		if 0 is result = .Send('VirtualList_AddGlobalMenu?')
			return .addGlobalMenu? = true
		return .addGlobarMenu? = result
		}

	layout(addons, filterBy, titleLeftCtrl)
		{
		addons = addons.Copy().Append(GetContributions('FormatContextMenuItems'))
		addons.Addon_VirtualListTopFilters = [filtersOnTop: .filtersOnTop,
			checkBoxAmountField: .checkBoxAmountField]
		addons.Addon_VirtualListView_Edit = true
		addons.Addon_VirtualListView_SetStatus = true
		addons.Addon_VirtualListView_Thumb = true
		addons.Addon_VirtualListView_Header = true
		addons.Addon_VirtualListView_Grid = true
		addons.Addon_VirtualListView_ContextMenu = true
		addons.Addon_VirtualListFilter = [:filterBy]
		addons.Addon_VirtualListView_Expand = true
		.Addons = AddonManager(this, addons)
		return VirtualListViewLayout(this, titleLeftCtrl)
		}

	firstTime: true
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		if .firstTime
			{
			.firstTime = false
			row = false
			if .linked? is false and .model isnt false
				if false isnt row = .model.SetFirstSelection()
					.grid.SetFocusedRow(row)
			if .filtersOnTop
				.Addons.Send('LoadSavedFilters', this)

			// have to do after LoadSavedFilters otherwise extra
			// buttons are not constructed yet
			if row isnt false
				if false isnt rec = .model.GetRecord(row - .model.Offset)
					{
					.Send('VirtualList_ItemSelected', rec)
					.RefreshValid(rec)
					}
			}
		}

	ToggleFilter()
		{
		.Addons.Send('ToggleFilter', this)
		}

	UpdateTopFilters(vals = false)
		{
		if .topFilters isnt false
			.topFilters.Set([conditions: vals is false ? .Select_vals : vals])
		}
	Select_ExtraButtons()
		{
		.Send('VirtualList_Select_ExtraButtons')
		}
	Select_ExtraLayout()
		{
		extraLayout = Object(buttons: .Send('VirtualList_Select_ExtraButtons')
			checkBoxAmountField: .checkBoxAmountField)
		if not Object?(extraLayout.buttons)
			extraLayout.buttons = Object()
		return extraLayout
		}

	GetDefaultSelect()
		{
		return .select
		}

	Record_NewValue(field, value, source)
		{
		.Field_SetFocus(source)
		.Addons.SendToOneAddon('Record_NewValue', field, value, source)
		}

	SetFilter(filters = false)
		{
		return .Addons.Collect('SetFilter', filters).Join(' ')
		}

	model: 	false
	SetModel(model, selectName = '', updateWhere = false)
		{
		.model = model
		.expandBar.SetInfo(.model, .grid.RowHeight, .header.Ymin, .expandBtns)
		.contextMenu.BuildMenu(.model, this, .preventCustomExpand?, .excludeCustomize?)
		.SetStatusBar('')
		.Send('VirtualList_ModelChanged')
		.setModelWhere(selectName, updateWhere)
		if .switchToForm isnt true
			.SetDefaultStatus()
		.Repaint()
		}

	setModelWhere(selectName, updateWhere)
		{
		whereStr = .model.ColModel.GetSelectWhere(
			selectName, this, .model.AllAvailableColumns) $
			.Addons.Collect('ExtraWhere').Join(' ')
		if not whereStr.Blank?()
			.model.SetWhere(whereStr)
		else if updateWhere isnt false
			.model.SetWhere(updateWhere)
		}

	GetModel()
		{
		return .model
		}

	GetTableName()
		{
		return .model.GetTableName()
		}

	Seek(field, prefix)
		{
		row_num = .model.Seek(field, prefix)
		.Repaint()
		.grid.SelectRow(row_num)
		}

	SelectRecordByKeyPair(key, field)
		{
		if false isnt rec = .model.GetRecordByKeyPair(key, field)
			.grid.SelectRecord(rec)
		}

	SelectRecord(rec)
		{
		if rec isnt false
			.grid.SelectRecord(rec)
		}

	VirtualListScroll_Resize(width, height)
		{
		if .model is false or height is 0
			return 'none'

		height -= .header.Ymin + .horzScrollSize(width)

		.model.UpdateVisibleRows(.grid.GetRows(height))

		.thumb.SetThumbPosition(.model.GetPosition())

		return not .model.Begin?() or not .model.End?()
		}

	VirtualListScroll_FloatingPosition(width)
		{
		return .model.ColModel.HasSelectedVals?() and not .filtersOnTop
			? .horzScrollSize(width)
			: false
		}

	horzScrollSize(width)
		{
		return .model.ColModel.GetTotalWidths() > width
			? GetSystemMetrics(SM.CXHSCROLL) : 0
		}

	GetReadOnly()
		{
		return .grid.GetReadOnly()
		}

	RepaintExpandBar()
		{
		.expandBar.ShowEditButtons()
		.grid.ShowExpandButton()
		.expandBar.Repaint()
		}

	// from inner virtual list
	VirtualList_ItemSelected(rec /*unused*/, source)
		{
		.ClearSelect()
		if .model.ExpandModel isnt false
			.model.ExpandModel.ClearAllSelections(source)
		}

	RefreshValid(rec = false)
		{
		.Addons.SendToOneAddon('RefreshValid', this, .statusBar, rec)
		}
	SetStatusBar(msg, normal = false, warn = false, invalid = false)
		{
		.Addons.SendToOneAddon('SetStatusBar', .statusBar, msg, :normal, :warn, :invalid)
		}

	VirtualListHighlightRecords(recs, clr)
		{
		.grid.HighlightRecords(recs, clr)
		}
	RecordDirty?(dirty?)
		{
		.Send('RecordDirty?', dirty?)
		}
	VirtualListDirty?()
		{
		return not .model.EditModel.AllowLeaving?()
		}
	RecordLocked?(rec)
		{
		return .model.EditModel.RecordLocked?(rec)
		}
	Default(@args)
		{
		VirtualListMethodsMap(this, args)
		}

	Recv(@args)
		{
		.Addons.SendToOneAddon(@args)
		}

	// used to get specific value from selected record when using context menu
	GetContextColumn()
		{
		return .contextMenu.ContextCol
		}

	GetContextMenu()
		{
		return .contextMenu
		}

	GetContextRecordMenu()
		{
		return .recMenu
		}

	Save() // for Drill Down
		{
		return .SaveOutstandingChanges()
		}

	GetTitle()
		{
		if .title isnt false
			return .title
		if false isnt key = .GetAccessCustomKey()
			return key
		return .Option
		}

	AddRecord(record, pos = 'current') // pos: current, start, end
		{
		.grid.InsertRow(:record, :pos, force:)
		}

	GetMandatoryCols()
		{
		return .model.EditModel.Editable?() ? .model.ColModel.GetMandatoryFields() : #()
		}

	Customizable_ExpandInfo()
		{
		if 0 is defaultLayout = .Send('VirtualList_DefaultExpandLayout')
			defaultLayout = ''
		availableFields = .ExpandColumns()
		return Object(:defaultLayout, :availableFields)
		}

	ExpandColumns()
		{
		fields = .model.ColModel.GetAvailableColumns(.model.GetQuery())
		removeCols = Customizable.GetNonPermissableFields(.model.GetQuery())
		return fields.Difference(removeCols).Difference(.expandExcludeFields)
		}

	GetPrimarySort()
		{
		return .model.GetPrimarySort()
		}

	SaveSort?()
		{
		return .model.SaveSort?()
		}

	/* Used by Customize Column Dialog */
	GetCheckBoxField()
		{
		return .checkBoxColumn
		}

	MoveColumnToFront(col)
		{
		return .header.MoveColumnToFront(col)
		}

	Repaint(keepPos? = false)
		{
		.header.SetColModel(.model.ColModel, .GetPrimarySort(),
			showExpandBar?: .expandBar.ShowExpand?())
		.grid.SetModel(.model)
		if .Send('VirtualList_RepaintDisabled') isnt true // in the middle of setting
			{
			.grid.Repaint(:keepPos?)
			.RepaintExpandBar()
			.scroll.ResizeWindow()
			.thumb.SetSelectPressed(.model.ColModel.HasSelectedVals?())
			}
		}
	/* end of Used by Customize Column dialog */

	// to avoid repainting other elements, and keep the selection
	RepaintGrid(checkTotal = false)
		{
		.grid.Repaint()
		if checkTotal and
			.grid.GetClientRect().GetWidth() > .model.ColModel.GetTotalWidths()
			.grid.ScrollToLeft()
		}

	VertScroll(rows)
		{
		.grid.VertScroll(rows)
		}

	GetAlertTitle()
		{
		return .title is false ? "List" : .title // once only
		}

	// from inner virtual list or expand bar
	VirtualList_MouseWheel(wParam)
		{
		.grid.MOUSEWHEEL(wParam)
		}

	ClearSelect()
		{
		.grid.ClearSelect()
		}

	SelectControl_Changed()
		{
		.invalidateSelect()
		.UpdateTopFilters()
		}

	invalidateSelect()
		{
		.ClearSelect()
		if false isnt row = .model.SetFirstSelection()
			.grid.SetFocusedRow(row)
		}

	GetCheckedRecords()
		{
		return .model.GetCheckedRecords()
		}

	CheckRecordByKeys(keys)
		{
		for key in keys
			.model.CheckRecordByKey(key, forceCheck:)
		.UpdateTotalSelected(recalc:)
		.Repaint()
		}

	CheckAll()
		{
		.model.CheckAll()
		.Repaint()
		}

	UncheckAll()
		{
		.model.UncheckAll()
		.Repaint()
		}

	GetSelectedRecords()
		{
		return .grid.GetSelectedRecords()
		}

	SetRecordsToTop(field, values)
		{
		if .model is false
			return

		.collapseAllWhenSetRecsTop(values, field)
		if false is .model.SetRecordsToTop(field, values)
			return
		.model.UpdateOffset(-.model.Offset)
		.Repaint()
		.grid.SelectTopRecords(field, values)
		}

	collapseAllWhenSetRecsTop(values, field)
		{
		if .model.ExpandModel is false
			return
		recs = .GetLoadedData()
		recOnTop = recs.Any?({ values.Has?(it[field]) 		})
		expanded = recs.Any?({ it.vl_expanded_rows isnt '' 	})
		if recOnTop and expanded
			.CollapseAll(keepPos?:)
		}

	HighlightValues(member, values, color)
		{
		return .grid.HighlightValues(member, values, color)
		}

	ClearHighlight()
		{
		return .grid.ClearHighlight()
		}

	/* Used by Select dialog */
	GetSelectedRecord()
		{
		.grid.GetSelectedRecord()
		}

	GetCurrentSelectedData()
		{
		.grid.GetSelectedRecord()
		}

	GetFields()
		{
		return .model.ColModel.GetOriginalColumns()
		}

	Getter_Option()
		{
		if 0 isnt opt = .Send('VirtualList_GetOption')
			return opt
		return .model.ColModel.GetColumnsSaveName()
		}

	Getter_Select_vals()
		{
		return .model.ColModel.GetSelectVals()
		}

	SetSelectVals(select_vals)
		{
		.model.ColModel.SetSelectVals(select_vals)
		FilterButtonControl.UpdateStatus(this, .model.ColModel.HasSelectedVals?())
		}

	GetCurrentSelectWhere()
		{
		return .selectWhere
		}

	selectWhere: ''
	SetWhere(where)
		{
		.selectWhere = .model.ColModel.ExtendWithWhere(where, .model.AllAvailableColumns)
		.ResetWhere()
		.Send('VirtualList_SetWhere')
		.expandBar.ShowEditButtons()
		.SetDefaultStatus()
		return true
		}

	SetDefaultStatus()
		{
		if .statusBar is false
			return

		topClosed? = .topFilters is false
		.statusBar.SetDefaultMsg(
			.model.ColModel.GetSelectMgr().UsingDefaultFilter?() and topClosed?
			? 'initial Select applied'
			: '')
		}

	AfterTopFilter(type)
		{
		if type is "open" and .statusBar isnt false
			.statusBar.SetDefaultMsg("")
		if type is "close"
			.SetDefaultStatus()
		}

	GetSelectFields()
		{
		// option for extra selectable fields but not avaialble on columns
		if 0 is extra = .Send('VirtualList_ExtraSelectFields')
			extra = #()
		return .model.ColModel.GetSelectFields(extra)
		}
	/* end of Used by Select dialog */

	ResetWhere()
		{
		updateWhere = .selectWhere $ .Addons.Collect('ExtraWhere').Join(' ')
		.SetModel(.model, :updateWhere)
		}

	GetCustomFields()
		{
		if .model is false
			return #()
		return .model.ColModel.GetCustomFields()
		}

	SetReadOnly(readOnly)
		{
		if readOnly isnt true and .readonly
			return
		.grid.SetReadOnly(readOnly)
		}

	Editable?()
		{
		return .model.EditModel.Editable?() and not .GetReadOnly()
		}

	ReplaceRecord(oldRec, newRec)
		{
		.model.ReplaceRecord(oldRec, newRec)
		}

	ReloadRecord(rec, discard = false)
		{
		if discard is true and .model.OwnLock?(rec)
			{
			.model.NextNum.CheckPutBackNextNum(rec)
			.model.UnlockRecord(rec) // also clears changes
			return .ForceEditMode(rec)
			}
		if false isnt rec = .model.ReloadRecord(rec, force:)
			.grid.RepaintRecord(rec)
		return Record?(rec)
		}

	ReloadRecordByKeyPair(key, field)
		{
		if false isnt rec = .model.GetRecordByKeyPair(key, field)
			.ReloadRecord(rec)
		}

	extraSetupRecord: false
	BeforeRecord(x)
		{
		if .extraSetupRecord is false
			.extraSetupRecord = .Send('VirtualList_ExtraSetupRecordFn')
		if .extraSetupRecord isnt 0
			(.extraSetupRecord)(x, source: this)
		.Addons.Send('BeforeRecord', x)
		}

	ForceSetExtraSetupRecordFn()
		{
		.extraSetupRecord = .Send('VirtualList_ExtraSetupRecordFn')
		}

	SetSelectApplied(applied)
		{
		if false isnt .topFilters
			.topFilters.SetSelectApplied(applied)
		}

	RepaintSelectedRows()
		{
		.grid.RepaintSelectedRows()
		}

	GetLoadedData()
		{
		return .model.GetLoadedData().Copy()
		}
	GetDeleted()
		{
		return .model.GetLoadedData().Filter({ it.vl_deleted isnt true })
		}
	Field_SetFocus(source)
		{
		if true isnt .Send('VirtualList_SelectRowOnExpandFocus')
			return
		if .model.ExpandModel is false or
			false is ctrl = .model.ExpandModel.GetControl(source)
			return
		rec = ctrl.GetControl().Get()
		if ((false isnt cur = .GetSelectedRecord()) and cur is rec)
			return

		.grid.SelectRecord(rec)
		}
	SelectRecordByFocus(curFocus)
		{
		if curFocus isnt false and .model.ExpandModel isnt false and
			false isnt rec = .model.ExpandModel.GetCurrentFocusedRecord(curFocus)
			.grid.SelectRecord(rec)
		}

	SelectExpandRecord(source)
		{
		if false is result = .GetExpandCtrlAndRecord(source)
			return false

		.grid.SelectRecord(result.rec)
		return true
		}

	GetData(source)
		{
		if false isnt result = .GetExpandCtrlAndRecord(source)
			return result.rec
		return false
		}

	On_VirtualList_Button(source)
		{
		if false is ctrl = .model.ExpandModel.GetControl(source)
			return 0
		rec = ctrl.GetControl().Get()
		.Send("On_" $ source.Get(), :rec, :ctrl)
		}

	/* called from KeyControl messages */
	ExpandedIndex: false
	GetField(field)
		{
		rec = .ExpandedIndex is false
			? .grid.GetSelectedRecord()
			: .model.GetRecord(.ExpandedIndex)
		return rec[field]
		}
	SetField(field, value, idx /*unused*/ = false, invalidate /*unused*/ = false,
		source = false)
		{
		rec = source isnt false and not source.Window.Base?(ListEditWindow)
			? .GetExpandCtrlAndRecord(source).rec
			: .grid.GetSelectedRecord()
		Assert(rec isnt: false, msg: "VirtualList - cannot find source line")
		rec[field] = value
		}
	SetMainRecordField(field, value, source = false)
		{
		.SetField(field, value, :source)
		}
	GetTransQuery()
		{
		return .model.GetQuery()
		}
	/* called from KeyControl messages */

	GetGrid()
		{
		return .grid
		}

	GetGridHwnd()
		{
		return .grid.Hwnd
		}

	GetViewControls()
		{
		return Object(header: .header, thumb: .thumb, expandBar: .expandBar,
			expandBtns: .expandBtns, scroll: .scroll)
		}

	GetOutstandingChanges(all? = false)
		{
		return .model.EditModel.GetOutstandingChanges(:all?)
		}

	GetHdrCornerControl()
		{
		return .scroll.GetHdrCornerCtrl()
		}

	GetExpandBarWidth()
		{
		return .expandBar.Xmin
		}

	GetQuery()
		{
		return .model.GetQuery()
		}

	GetAccessCustomKey()
		{
		return .model.ColModel.GetCustomKey()
		}

	Map_GetAddress(source) // built-in on Address 1
		{
		return .model.ExpandModel.Map_GetAddress(source)
		}

	FieldPrompt_GetSelectFields(source = false)
		{
		.Controller.Send('FieldPrompt_GetSelectFields',
			rec: .model.ExpandModel.GetExpandRecord(source), :source)
		}

	SelectDisabled()
		{
		return .disableSelectFilter is true
		}

	FiltersOnTop?()
		{
		return .filtersOnTop
		}

	IsLinked?()
		{
		return .linked?
		}

	QueueDeleteAttachmentFile(newFile, oldFile, name, action)
		{
		curRec = .GetCurrentSelectedData()
		return .model.QueueDeleteAttachmentFile(newFile, oldFile, curRec, name, action)
		}

	RestoreAttachmentFiles()
		{
		.model.CleanupAttachments(true)
		}

	Valid?()
		{
		if .readonly or .GetReadOnly()
			return true
		save = VirtualListSave(.grid, .model, .contextMenu.UpdateHistory)
		return save.Valid?()
		}
	}
