// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
CommandParent
	{
	Xmin: 300
	Ymin: 200
	Name:		'VirtualList'
	ComponentName: 'VirtualList'
	model: 		false
	view: 		false
	Commands:	()

	New(query = false, .columns = #(), .columnsSaveName = false,
		.headerSelectPrompt = false, menu = false, .startLast = false,
		headerMenu = false, readonly = false, filterBy = false,
		.hideCustomColumns? = false, .mandatoryFields = #(), .enableMultiSelect = false,
		.checkBoxColumn = false, .checkBoxAmountField = false, .lockFields = #(),
		.disableSelectFilter = false, .sortSaveName = false,
		thinBorder = false, .protectField = false, .validField = false,
		.title = false, .disableCheckSortLimit? = false,
		.option = false, historyFields = false, titleLeftCtrl = false, .nextNum = false,
		.excludeSelectFields = #(), .loadAll? = false, .extraFmts = false,
		hdrCornerCtrl = false, expandExcludeFields = #(), addons = #(), .keyField = false,
		.stickyFields = false, .enableUserDefaultSelect = false, .stretchColumn = false,
		.filtersOnTop = false, select = #(), .saveQuery = false,
		.hideColumnsNotSaved? = false, enableDeleteBar = false, .linked? = false,
		.preventCustomExpand? = false, switchToForm = false, .defaultColumns = false,
		.asof = false, .useQuery = 'auto', excludeCustomize? = false,
		.warningField = false)
		{
		super(Object('VirtualListView', menu, .headerSelectPrompt, headerMenu,
			readonly, filterBy, .checkBoxColumn, checkBoxAmountField, disableSelectFilter,
			thinBorder, .protectField, .validField, .title,
			.option, historyFields, titleLeftCtrl, :hdrCornerCtrl, :expandExcludeFields,
			:addons, filtersOnTop: .filtersOnTop, :select, :enableDeleteBar,
			linked?: .linked?, :preventCustomExpand?, :switchToForm, :excludeCustomize?))
		.view = .VirtualListView
		.baseQuery = query

		if query isnt false and query isnt ''
			.SetQuery(query, columns, asof)

		if .protectField isnt false and false isnt .Send("VirtualList_NeedValidation?")
			{
			.needValidation? = true
			.Window.AddValidationItem(this)
			}
		if .protectField isnt false // editable
			.Commands = #(
				("New",		"Ctrl+N")
				("Edit",	"Alt+E")
				("Select",	"Alt+S"))
		}

	selectName: ''
	SetQuery(query, columns = #(), filters = false, customKey = false,
		selectName = '', .asof = false, .useQuery = 'auto')
		{
		where = .VirtualListView.SetFilter(filters)
		recycled = .recycleExpands()
		if .model isnt false
			.model.Destroy()

		extraSetupRecord = .VirtualListView.BeforeRecord

		customKey = .getCustomKey(customKey, query)
		.handleSorts(customKey)

		useExpandModel? = not .preventCustomExpand? or
			.VirtualListView.VirtualListGrid_Expand([]) isnt 0
		.model = VirtualListModel(query, .startLast, columns, .columnsSaveName,
			.headerSelectPrompt, where, .hideCustomColumns?, .mandatoryFields,
			.checkBoxColumn, .checkBoxAmountField, sortSaveName: .sortSaveName,
			lockFields: .lockFields, protectField: .protectField,
			:extraSetupRecord, observerList: .view, validField: .validField,
			nextNum: .nextNum, disableCheckSortLimit?: .disableCheckSortLimit?,
			excludeSelectFields: .excludeSelectFields, loadAll?: .loadAll?,
			extraFmts: .extraFmts, :customKey, keyField: .keyField,
			stickyFields: .stickyFields, enableUserDefaultSelect:.enableUserDefaultSelect,
			disableSelectFilter: .disableSelectFilter, stretchColumn: .stretchColumn,
			enableMultiSelect: .enableMultiSelect, saveQuery: .saveQuery,
			hideColumnsNotSaved?: .hideColumnsNotSaved?, linked?: .linked?,
			:useExpandModel?, option: .option, defaultColumns: .defaultColumns, :asof,
			useQuery: .useQuery, warningField: .warningField)
		.Send('RegisterLinkedBrowse', this, .Name)
		if selectName isnt ''
			.selectName = selectName

		.restoreRecycledExpands(recycled)
		if .model.EditModel.Editable?() and
			.Send('VirtualList_CustomizeColumnAllowHideMandatory?') isnt true
			.model.ColModel.AddMissingMandatoryCols()
		.view.SetModel(.model, .selectName)
		}

	handleSorts(customKey)
		{
		if .sortableList?()
			{
			if false isnt .Send('VirtualList_SaveSort?') and .sortSaveName is false and
				customKey isnt false
				.sortSaveName = customKey $ ' Sort'
			}
		else
			.disableCheckSortLimit? = true
		}

	sortableList?()
		{
		return .filtersOnTop is true and .Send('VirtualList_DisableSort?') isnt true
		}

	recycleExpands()
		{
		// assuming all the expands are using the same layout
		recycled = Object()
		if .model isnt false and .model.ExpandModel isnt false
			recycled = .model.ExpandModel.RecycleExpands()
		return recycled
		}

	getCustomKey(customKey, query)
		{
		if customKey isnt false
			return customKey

		customKey = .Send('VirtualList_GetCustomKey')
		if customKey not in (0, false)
			return customKey

		return ListCustomize.BuildCustomKeyFromQueryTitle(query, .title)
		}

	restoreRecycledExpands(recycled)
		{
		if .model.ExpandModel isnt false
			.model.ExpandModel.SetRecycledExpands(recycled)
		}

	Getter_Select_vals()
		{
		return .model.ColModel.GetSelectVals()
		}
	SetSelectVals(select_vals)
		{
		.view.SetSelectVals(select_vals)
		}
	SetDefaultSelect(defaultSel)
		{
		.model.ColModel.GetSelectMgr().PrependInitialSelect(defaultSel)
		sf = .GetSelectFields()
		where = SelectRepeatControl.BuildWhere(sf, .Select_vals)
		.SetWhere(sf.Joins(where.joinflds) $ where.where, quiet:)
		}
	ApplySelects()
		{
		sf = .GetSelectFields()
		where = SelectRepeatControl.BuildWhere(sf, .Select_vals)
		whereStr = sf.Joins(where.joinflds) $ where.where
		result = .SetWhere(whereStr, quiet:)
		.view.SelectControl_Changed()
		return result
		}

	Refresh()
		{
		if not .SaveFirst()
			return
		if .model.ExpandModel isnt false
			.model.ExpandModel.CollapseAll()
		.model.RefreshData()
		.view.Repaint(keepPos?:)
		}

	Msg(args)
		{
		args.source = this
		super.Msg(args)
		}

	Valid?()
		{
		if false is .Send("VirtualList_ExtraValid")
			return false

		return .view.Valid?()
		}

	On_Edit(source, force = false) // from shortcut
		{
		return .view.On_Edit(source, force)
		}

	On_New() // from shortcut
		{
		if .model isnt false
			.view.On_New()
		}

	On_Select()
		{
		if .model isnt false and .disableSelectFilter isnt true
			.view.On_VirtualListThumb_ArrowSelect()
		}

	Redir(msg, ctrl = 'focus')
		{
		.view.Redir(msg, ctrl)
		}

	// Sent from SelectControl, required to redirect to VirtualListViewControl
	SelectControl_Changed()
		{
		.view.SelectControl_Changed()
		}

	UsingCursor?()
		{
		.model.UsingCursor?()
		}

	Default(@args)
		{
		event = args[0]
		.view[event](@+1 args)
		}

	ConfirmDestroy()
		{
		if .model is false
			return true

		SetFocus(NULL)
		.view.SaveOutstandingChanges()
		return .model.EditModel.AllowLeaving?()
		}

	QueueDeleteAttachmentFile(newFile, oldFile, name, action)
		{
		.view.QueueDeleteAttachmentFile(newFile, oldFile, name, action)
		}

	RestoreAttachmentFiles()
		{
		.view.RestoreAttachmentFiles()
		}

	needValidation?: false
	Destroy()
		{
		if .protectField isnt false and .needValidation?
			.Window.RemoveValidationItem(this)
		if .model isnt false
			{
			if .model.ExpandModel isnt false
				.model.ExpandModel.DestroyAll() //destroy constructed and recycled expands
			.model.Destroy()
			}
		super.Destroy()
		}
	}
