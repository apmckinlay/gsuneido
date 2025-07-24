// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(query, cols, title = '', hwnd = 0)
		{
		ToolDialog(hwnd, Object('Params',
			Object(.print_rep, Query: query, Cols: cols, :title),
			:title, name: title))
		}

	print_rep: QueryFormat
		{
		Query()
			{
			return .Query
			}
		Output()
			{
			ob = Object('Row')
			for col in .Cols
				ob.Add(col)
			return ob
			}
		}
	}
