// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Browse All Notes'
	CallClass(hwnd)
		{
		return ToolDialog(hwnd, this, border: 0)
		}
	New()
		{
		.list = .FindControl('List')
		.list.SetReadOnly(true, false) // to prevent F2 from activating cell-edit
		.load()
		}
	Controls: (Vert,
		(List, columns: #(name, text, path), columnsSaveName: 'BrowseAllNotes',
			resetColumns:, noShading:),
		(Horz, (Skip, small:), (Static, 'Double-click to view')))

	List_DoubleClick(row, col)
		{
		if row is false or col is false
			return 0

		data = .list.GetRow(row)
		NotesEditControl(.Window.Hwnd, data.name,
			UserNotes.Query(data.name, data.path), book_option: data.path,
			nestedCtrl:)
		return true
		}

	load()
		{
		.list.DeleteAll()
		QueryApply('user_notes')
			{ |x|
			if x.usernote_permission? isnt true
				continue
			.list.AddRow(x)
			}
		}
	}