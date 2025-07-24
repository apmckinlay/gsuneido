// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
CodeViewControl
	{
	addons: #(
		Addon_suneido_style: (query:),
		Addon_wrap:,
		Addon_brace_match: false,
		Addon_calltips: false,
		Addon_class_outline: false,
		Addon_folding: false,
		Addon_go_to_line: false,
		Addon_highlight_cursor_line: false,
		Addon_indent_guides: false,
		Addon_show_line_numbers: false,
		Addon_show_margin: false,
		Addon_status: false
		)
	New()
		{
		super(addons: .addons, readonly:)
		.Editor.SetWrap(true)
		}

	Scintilla_DoubleClick()
		{
		line = .Editor.GetLine()
		if line =~ '^\w+:[[:upper:]][_[:alnum:]]*[?!]?\>'
			{
			_hwnd = .Window.Hwnd
			.Editor.Home()
			.Editor.LineEndExtend()
			.Editor.On_Context_Go_To_Definition()
			}
		}
	}