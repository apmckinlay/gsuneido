// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	New(.parent, .model)
		{
		.editModel = .model.EditModel
		.Hwnd = .parent.Hwnd
		.Window = Object(Hwnd: .parent.Window.Hwnd)
		.UniqueId = .parent.UniqueId
		.Controller = this
		}

	editor: false
	Editing: false
	EditCell(rec, col)
		{
		msg = .allowEdit(rec, col)
		if msg isnt '' and msg isnt true
			{
			if String?(msg)
				.parent.AlertInfo('Reason Protected', msg)
			return false
			}
		return .editCell(col)
		}

	editCell(col)
		{
		.col = col
		customFields = .model.ColModel.GetCustomFields()
		custom = customFields isnt false
			? customFields.GetDefault(col, false)
			: false
		readonly = FieldProtected?(col, .rec, .editModel.ProtectField)
		.Editing = true
		row = .model.GetRecordRowNum(.rec) + .model.Offset
		.editor = new ListEditWindow(
			0,	readonly, .model.ColModel.FindCol(.col), row, this,
			:custom, :customFields)
		return .rec
		}

	allowEdit(rec, col)
		{
		freshRec = .model.ReloadRecord(rec)
		if String?(freshRec)
			return freshRec
		.rec = freshRec
		allow? = .parent.Controller.Send(
			"VirtualList_AllowCellEdit", col, rec: freshRec, oldrec: rec)
		if allow? isnt '' and allow? isnt 0
			return allow?
		if true isnt msg = .model.LockRecord(freshRec)
			return msg
		return ''
		}

	ChildOf?(hwnd)
		{
		if .editor is false
			return false
		return .editor.ChildOf?(hwnd)
		}

	ClosingListEdit()
		{
		if .editor is false
			return false
		return .editor.ClosingListEdit
		}

	GetEditorHwnd()
		{
		if .editor is false
			return 0
		return .editor.Hwnd
		}

	GetClientRect()
		{
		return .parent.GetClientRect()
		}

	GetRow(unused)
		{
		return .rec
		}

	GetCol(unused)
		{
		return .col
		}

	Default(@args)
		{
		if not .parent.Member?('Controller')
			return 0
		if not .parent.Controller.Method?(args[0])
			return 0
		.parent.Controller[args[0]](@+1args)
		}

	ListEditWindow_Commit(col /*unused*/, row /*unused*/, dir, data, valid?,
		unvalidated_val = '', readonly = false, dirty? = false)
		{
		.parent.Controller.Send('VirtualList_PreCommitCellValue', .rec)
		.Editing = false
		.editor = false
		if dirty? and not readonly and .valueChanged?(valid?, unvalidated_val, data)
			.commitChange(data, valid?, unvalidated_val)
		if dir isnt 0
			.editNextCell(.col, dir)
		}

	valueChanged?(valid?, unvalidated_val, data)
		{
		prevInvalid = ListControl.GetInvalidFieldData(.rec, .col)
		if prevInvalid isnt '' and valid? is true // from invalid -> valid
			invalidValChg? = true
		else
			invalidValChg? = not valid? and	unvalidated_val isnt prevInvalid
		return invalidValChg? or data isnt .rec[.col]
		}

	commitChange(data, valid?, unvalidated_val = '')
		{
		_committing = .col
		.CommitChange(.parent, .rec, .col, data, valid?, :unvalidated_val)
		}

	// also called by expand editing
	CommitChange(grid, rec, col, data, valid?, unvalidated_val = '')
		{
		model = grid.Controller.GetModel()
		grid.Send('RecordDirty?', true)
		model.NextNum.CheckPutBackNextNum(rec, col, data)
		rec[col] = data // triggers observer first
		if valid?
			{
			if true is result = grid.Controller.Send('VirtualList_CommitCellValue',
				rec, col, data)
				model.EditModel.ClearChanges(rec)
			else if Object?(result)
				.editModel.AddChanges(rec, col, data, invalidCols: result.Copy())
			else if result is 0 // handle saving automatically
				model.EditModel.RemoveInvalidCol(rec, col)
			}
		else
			{
			model.EditModel.AddInvalidCol(rec, col)
			}
		if ListControl.SetInvalidFieldData(rec, col, unvalidated_val)
			.Send('VirtualList_InvalidDataChanged', rec)
		grid.Controller.Send('VirtualList_AfterField', col, data, rec)
		if model.CheckBoxColModel isnt false
			{
			model.CheckBoxColModel.AutoSelectByAmount(col, data, rec)
			grid.Controller.UpdateTotalSelected(recalc:)
			}
		model.ColModel.Plugins_Execute(data: rec, field: col, hwnd: grid.Hwnd,
			query: model.GetQuery(), pluginType: 'AfterField')

		grid.Controller.Send('VirtualList_AfterChanged', record: rec, saved: false)
		grid.RepaintSelectedRows()
		}

	editNextCell(col, dir)
		{
		numCols = .model.ColModel.GetColumns().Size()
		colIndex = .model.ColModel.FindCol(col)
		forever
			{
			colIndex += dir
			rowChange = 0
			if colIndex >= numCols
				{
				rowChange = 1
				colIndex = numCols - colIndex
				}

			if colIndex < 0
				{
				rowChange = -1
				colIndex = numCols + colIndex
				}

			nextCol = .model.ColModel.Get(colIndex)
			if rowChange isnt 0 and not .rowChanged(rowChange)
				return
			if false is record = .parent.GetSelectedRecord()
				return
			if .nextColEditable(nextCol, colIndex, record)
				{
				.parent.ScrollColToView(colIndex)
				.editCell(nextCol)
				return
				}
			}
		}
	rowChanged(rowChange)
		{
		if .model.AutoSave? is true and false is
			.parent.Send('VirtualListGrid_SaveRecord', .parent.GetSelectedRecord())
			return false
		return .moveNextRow?(rowChange)
		}
	moveNextRow?(rowChange)
		{
		if false is .parent.SelectNextRow(rowChange)
			{
			if rowChange is 1
				.parent.InsertRow(pos: 'end')
			return false
			}
		return true
		}
	minWidth: 5
	nextColEditable(nextCol, colIndex, record)
		{
		if record.vl_deleted is true
			return false

		if .model.ColModel.GetColWidth(colIndex) <= .minWidth or
			.model.EditModel.ProtectField is false
			return false

		customFields = .model.ColModel.GetCustomFields()
		if true is ListCustomize.GetControlOption(customFields, nextCol, 'tabover')
			return false

		result = .allowEdit(record, nextCol)
		return result is true or result is ''
		}

	Return()
		{
		if .editor isnt false
			.editor.Return()
		.Editing = false
		.editor = false
		}
	}