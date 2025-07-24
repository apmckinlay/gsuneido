// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'LinkList'
	Title: 'Pick a Path'
	CallClass(hwnd, list)
		{
		return ToolDialog(hwnd, Object(this, list))
		}

	New(.list)
		{
		.listbox = .Data.ListBox
		margin = 2 * GetSystemMetrics(SM.CXVSCROLL) // to avoid horizontal scroll bar
		.listbox.Xmin = .listbox.GetHorizontalExtent() + margin
		}

	Controls()
		{
		return Object('Record', Object('ListBox', .list.Map({ it.path $ '/' $ it.name})))
		}

	ListBoxSelect(i /*unused*/)
		{
		selected = .Get()
		selected.path = selected.path[.. selected.path.FindLast('/')]
		.Window.Result(selected)
		}

	Get()
		{
		name = ''
		path = .listbox.GetText(.listbox.GetCurSel())
		for item in .list
			if path.Prefix?(item.path) and path.Suffix?(item.name)
				name = item.name
		return false isnt path ? Object(:path, :name) : false
		}
	}
