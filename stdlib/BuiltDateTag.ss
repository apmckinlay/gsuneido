// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	tag: "^// BuiltDate ([<>]) (\d+)"

	// see also: CheckLibrary.BuiltDate_skip?
	Find(text)
		{
		if not text.Has?("// BuiltDate ")
			return false

		dir = false isnt text.Extract("^// BuiltDate > (\d+)") ? '>' : '<'
		lines = text.Lines()
		for i in lines.Members()
			{
			if dir isnt lines[i].Extract(.tag, 1)
				continue

			return Object(:dir, line: i, date: Date(lines[i].Extract(.tag, 2)))
			}

		return false
		}

	Format(dir, date)
		{
		return "// BuiltDate " $ dir $ ' ' $ date.Format("yyyyMMdd")
		}

	InsertLine(text)
		{
		lines = text.Lines()
		for i in lines.Members()
			if not lines[i].Prefix?("//")
				return i

		return lines.Size()
		}

	Edit(libview)
		{
		if false is editor = libview.Editor
			return false

		if false is choice = ToolDialog(libview.Window.Hwnd,
			[BuiltDateTagControl, .Find(editor.Get())],
			"BuiltDate Tag", closeButton?: false)
			return false

		.Apply(editor, choice)
		return true
		}

	Apply(editor, choice)
		{
		text = editor.Get()
		if false isnt cur = .Find(text)
			editor.SelectLine(cur.line)
		else
			editor.SetSelect(editor.PositionFromLine(.InsertLine(text)))

		editor.ReplaceSel(
			choice is #remove ? "" : .Format(choice.dir, choice.date) $ "\r\n")
		editor.SetFocus()
		}
	}
