// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Controller
	{
	Title: 'Undelete'
	New(lib, data)
		{
		super(.controls(lib, data))
		.list = .Vert.List
		defaultColWidth = 150
		.list.SetColWidth(0, defaultColWidth)
		.list.SetColWidth(1, defaultColWidth)
		.data = data
		}

	controls(lib, data)
		{
		.Title $= ' - ' $ lib
		return Object('Vert',
			Object('List', #(name, tran_asof), data)
			'Skip'
			#(Horz Fill (Button Undelete xmin: 75) Skip (Button Cancel xmin: 75))
			)
		}

	List_AllowCellEdit(col /*unused*/, row /*unused*/)
		{
		return false
		}

	List_DoubleClick(row, col /*unused*/)
		{
		return row isnt false ? .Window.Result(.data[row]) : true
		}

	On_Undelete()
		{
		sel = .list.GetSelection()
		if sel.Size() isnt 1
			Beep()
		else
			.Window.Result(.list.GetRow(sel[0]))
		}
	}
