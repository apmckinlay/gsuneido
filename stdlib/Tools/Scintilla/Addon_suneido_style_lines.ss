// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// style lines independently
// e.g. mismatched quotes will not affect following lines
// used by WorkSpaceOutputControl which is also used for WorkSpace Find
Addon_suneido_style
	{
	Style(from, to)
		{
		line = .LineFromPosition(from)
		from = .PositionFromLine(line)
		toLine = .LineFromPosition(to)
		for (; line <= toLine; ++line)
			{
			to = .PositionFromLine(line + 1)
			ScintillaStyle.Style(.Hwnd, from, line, to)
			from = to
			}
		}
	}