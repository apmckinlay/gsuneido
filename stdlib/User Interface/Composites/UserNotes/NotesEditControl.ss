// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(hwnd, title, query, book_option, nestedCtrl = false)
		{
		return ToolDialog(hwnd, [this, :title, :query, :book_option, :nestedCtrl],
			closeButton?: nestedCtrl, border: 0, keep_size: 'Notes',
			title: 'Notes for ' $ title)
		}
	New(.title, .query, .book_option, nestedCtrl = false)
		{
		.editor = .Vert.Editor
		.static = .Vert.Border.Horz.Static
		if false isnt x = Query1(.query)
			{
			.editor.Set(x.text)
			.static.Set("Last modified " $ x.last_modified.ShortDateTime() $
				" by " $ x.last_modified_by)
			}
		if AccessPermissions(.book_option) is 'readOnly'
			.editor.SetReadOnly(true)
		if nestedCtrl isnt false
			{
			.FindControl('Browse_All_Notes').SetVisible(false)
			.editor.SetReadOnly(true)
			.FindControl('Save').SetVisible(false)
			.FindControl('Cancel').Set('Close')
			}
		}
	Controls: (Vert
		(Skip 4)
		(Static ' These notes are attached to the screen in general, ' $
			'NOT to a particular record. ')
		(Skip 4)
		(ScintillaAddonsEditor, fontSize: 11)
		(Border (Horz
			(Static '') Fill Skip
			(LinkButton 'Browse All Notes') Skip
			(Button Save width: 8) Skip (Button Cancel width: 8))
			5))
	On_Browse_All_Notes()
		{
		if .editor.Dirty?() and YesNo("Save changes?", "Notes", .Window.Hwnd)
			.save()
		ToolDialog(.Window.Hwnd, 'AllNotes', border: 0)
		}
	On_Save()
		{
		.save()
		.Window.Result(true)
		}
	save()
		{
		text = .editor.Get().RightTrim()
		RetryTransaction()
			{ |t|
			if text is ''
				t.QueryDo('delete ' $ .query)
			else if false isnt (x = t.Query1(.query))
				{
				x.text = text
				x.last_modified = Date()
				x.last_modified_by = Suneido.User
				x.path = .book_option
				x.Update()
				}
			else
				t.QueryOutput(.query,
					[name: .title,
						:text,
						last_modified: Date(),
						last_modified_by: Suneido.User,
						path: .book_option])
			}
		.editor.Dirty?(false)
		}
	}
