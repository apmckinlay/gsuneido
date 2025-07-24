// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'ChooseManyAsObjectList'
	Xmin: 400
	Ymin: 200
	list: false

	New(list, .columns, .saveColName, .editableColumns = false,
		.selectCol = "choosemany_select", .noOkayCancel = false)
		{
		super(.layout(list, .columns))
		.orig_list = list
		.list = .Vert.List
		if .editableColumns is false
			.list.SetReadOnly(readOnly: true, grayOut: false)
		}
	layout(data, columns)
		{
		return Object('Vert'
			Object('List', :data, :columns, noDragDrop:, noHeaderButtons:, resetColumns:,
				columnsSaveName: .saveColName, checkBoxColumn: .selectCol),
			.noOkayCancel
				? #(HorzEqual, #(Button All) 'Skip' #(Button None))
				: #(AllNoneOkCancel))
		}

	List_SingleClick(row, col)
		{
		if col isnt 0 or row is false
			return 0

		.list.GetRow(row)[.selectCol] = .list.GetRow(row)[.selectCol] isnt true
		.list.RepaintRow(row)
		return 0
		}

	List_AllowCellEdit(col, row /*unused*/)
		{
		if Object?(.editableColumns) and .editableColumns.Has?(col)
			return true
		return false
		}

	List_DeleteRecord(@unused)
		{
		return false
		}

	List_WantNewRow(@unused)
		{
		return false
		}

	On_All()
		{
		.set_selected(true)
		}

	On_None()
		{
		.set_selected(false)
		}

	set_selected(value)
		{
		for row in .list.Get()
			row[.selectCol] = value
		.list.Repaint()
		}

	On_OK()
		{
		.Send('On_OK')
		}

	OK()
		{
		return .list.Get()
		}

	Destroy()
		{
		if .list isnt false
			.list.ClearSelectFocus()
		super.Destroy()
		}
	}
