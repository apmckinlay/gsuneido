// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(table)
		{
		Window(Object(this, table))
		}

	New(.table)
		{
		.Title = 'Check Records - ' $ .table
		}

	Controls()
		{
		extends = Object('table = ' $ Display(.table), 'check_record_code')
		columns = Object('name', 'check_record_code')
		select = Object(#(check_record_code))
		if QcIsEnabled()
			{
			extends.Add('check_record_qc')
			columns.Add('check_record_qc')
			select.Add(#(check_record_qc))
			}
		return Object('VirtualList',
			query: .table $ ' where group is -1 extend ' $ extends.Join(', '),
			columnsSaveName: 'CheckRecords',
			:columns, :select, filtersOnTop:)
		}

	VirtualList_DoubleClick(rec, col)
		{
		if col isnt 'name'
			return 0
		GotoLibView(rec.name, libs: Object(.table))
		return false
		}
	}