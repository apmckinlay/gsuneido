// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonsControl
	{
	IDE:,
	New(@args)
		{
		super(@args.Add(true, at: #readonly))
		}
	Addon_brace_match:,
	Addon_calltips:,
	Addon_highlight_occurrences:,
	Addon_suneido_style:,
	Addon_flag:,
	Addon_show_margin:,
	Addon_show_whitespace:,
	Addon_show_line_numbers:,
	Addon_go_to_definition:,
	Addon_go_to_line:,
	Addon_highlight_cursor_line:,
	Addon_indent_guides:,
	Addon_show_references:,
	Addon_scroll_zoom:
	}
