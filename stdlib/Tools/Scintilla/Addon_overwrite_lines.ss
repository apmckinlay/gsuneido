// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	ContextMenu()
		{
		.goto = .Send('Getter_Goto')
		if .goto is false
			return #()
		.lib = .goto.BeforeFirst(':')
		.name = .goto.AfterFirst(':')
		if BookTable?(.lib)
			return #()
		rec = Query1(.lib, name: .name, group: -1)
		if rec is false or rec.lib_modified is ''
			return #()
		return #('Overwrite Highlighted Lines')
		}

	On_Overwrite_Highlighted_Lines()
		{
		model = .Send('GetModel')
		local = model.GetLocalRec(.lib, .name)
		master = model.GetMasterRec(.lib, .name)
		model.CheckCommitted(SvcTable(.lib), .name, local, master)

		.calculateLineChanges()
		.overwriteLines(Query1(.lib, name: .name, group: -1))

		Defer(.refresh)
		LibUnload(.name)
		}

	calculateLineChanges()
		{
		diffs = .Send('Getter_Diffs')
		if diffs is false or diffs is 0
			return
		sel = .GetSelect()
		.firstLine = .LineFromPosition(sel.cpMin)
		.lastLine = .LineFromPosition(sel.cpMax)

		.prevLines = ''
		for (i = .firstLine; i <= .lastLine; i++)
			if diffs[i][1] isnt '>' and not diffs[i][1].Prefix?('+')
				.prevLines $= diffs[i][0] $ '\r\n'

		first = .PositionFromLine(.firstLine)
		last = .PositionFromLine(.lastLine+1)-1
		.SetSelect(first, last - first)

		last = .lastLine
		for (i = 0; i <= last; i++)
			{
			if diffs[i][1] is '<' or diffs[i][1].Prefix?('-')
				{
				if i < .firstLine
					.firstLine -= 1
				.lastLine -= 1
				}
			}
		}

	overwriteLines(rec)
		{
		if rec is false
			return

		lines = rec.lib_current_text.Lines()
		rec.text = Opt(lines[.. .firstLine].Join('\r\n'), '\r\n') $
			.prevLines $
			lines[.lastLine+1..].Join('\r\n')

		if rec.text.Suffix?('\r\n') is false
			rec.text $= '\r\n'

		CodeState(.lib, rec)
		}

	refresh()
		{
		.Send('ResetSelection')
		}
	}
