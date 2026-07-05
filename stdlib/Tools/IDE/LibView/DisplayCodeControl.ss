// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonsControl
	{
	IDE:,
	New(@args)
		{
		super(@.processArgs(args))
		}

	processArgs(args)
		{
		args.Add(true, at: #readonly)
		.SetupHighlightStyle(args, this)
		return args
		}

	Addon_brace_match:,
	Addon_highlight_occurrences:,
	Addon_flag:,
	Addon_show_margin:,
	Addon_show_whitespace:,
	Addon_show_line_numbers:,
	Addon_go_to_line:,
	Addon_highlight_cursor_line:,
	Addon_indent_guides:,
	Addon_show_references:,
	Addon_scroll_zoom:,

	// static
	SetupHighlightStyle(args, control)
		{
		type = args.GetDefault(#type, #code)
		if type is #code
			{
			control.Addon_suneido_style = true
			control.Addon_calltips = true
			control.Addon_go_to_definition = true
			}
		else if type is #html
			control.Addon_html = true
		else if type is #md
			control.Addon_md = true
		}
	}
