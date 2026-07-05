// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: VirtualListExpandBar
	New(.preventCustomExpand?, .enableDeleteBar = false)
		{
		}

	GetLayout()
		{
		return false
		}

	enableExpand: false
	SetInfo(.model, .rowHeight, .headerYmin, .expandBtns)
		{
		.grid = .Send('GetGrid')
		.showExpandBar()
		showExpandButton = .ShowExpand?() and .model.ExpandModel isnt false
		.grid.Act('SetExpandButtons', .expandBtns.GetExpandButtons(showExpandButton))
		}

	showExpandBar()
		{
		if .Send('VirtualListGrid_Expand', []) isnt 0
			.enableExpand = true
		else if not .preventCustomExpand?
			{
			if 0 is expandInfo = .Send('Customizable_ExpandInfo')
				expandInfo = Object(availableFields: false, defaultLayout: '')
			if 0 is customKey = .Send('GetAccessCustomKey')
				customKey = ''
			table = .model.GetTableName()
			c = Customizable(table, defaultLayout: expandInfo.defaultLayout,
				user: Suneido.User, :customKey)
			.enableExpand = c.LayoutExists?(CustomizeExpandControl.LayoutName)
			}
		}

	ShowExpand?()
		{
		return .enableExpand or .enableDeleteBar
		}

	ShowEditButtons()
		{
		.RefreshEditState()
		}

	RefreshEditState()
		{
		if .model is false or not .model.EditModel.Editable?() or
			.model.ExpandModel is false
			return
		if false is .Controller.Send('VirtualList_ShowEditButton?')
			return
		updates = Object()
		for rec in .model.ExpandModel.GetExpanded()
			{
			row = .model.GetRecordRowNum(rec) + .model.Offset
			updates[row] = not .needEditBtn?(rec)
				? 'hidden'
				: .model.EditModel.RecordLocked?(rec)
			}
		if updates.NotEmpty?()
			{
			grid = .Parent.FindControl(#VirtualListGrid)
			grid.Act('VirtualListExpandEditPushed', updates)
			}
		}

	needEditBtn?(rec)
		{
		// control data is only what is in the expand control (not the columns)
		if false isnt recordControl = .model.ExpandModel.GetRecordControl(rec)
			recordControl.GetControlData().Members().Each()
				{
				if not FieldProtected?(it, rec, .model.EditModel.ProtectField)
					return true
				}
		return false
		}

	Default(@unused) { }
	}
