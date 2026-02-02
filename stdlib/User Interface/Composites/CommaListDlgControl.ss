// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xmin: 325
	Ymin: 150
	Xstretch: 1
	Ystretch: 1
	Name: "CommaListDlg"
	Title: "Item List"
	New(list, max_items, .protectLineFunc = false)
		{
		super(.layout(max_items))
		.Vert.List.SetColWidth(0, 250) /*= default width */
		list_data = Object()
		for item in list.Split(',')
			list_data.Add(Object(commalist_value: item.Trim()))
		.Vert.List.Set(list_data)
		}
	layout(max_items)
		{
		footer = ''
		if max_items isnt false
			footer = Object('Horz'
				Object('Static' 'Max ' $ Display(max_items) $ ' items allowed'
					size: '+2', weight: 'bold'))

		return Object('Vert'
			Object('List' columns: #(commalist_value))
			#(Skip 5)
			footer
			)
		}
	OK()
		{
		ob = Object()
		data = .Vert.List.Get()
		for item in data
			if item.commalist_value isnt ''
				ob.Add(item.commalist_value)
		return ob.Join(',')
		}
	List_AllowCellEdit(col, row)
		{
		if col is -1
			return false

		if .protectLineFunc is false
			return true

		columns = .Vert.List.GetColumns()
		data = .Vert.List.Get()
		result = (.protectLineFunc)(data[row][columns[col]])
		return result is false
		}
	List_DeleteRecord(rec)
		{
		if .protectLineFunc is false
			return true

		result = (.protectLineFunc)(rec.commalist_value)
		return result is false
		}
	}
