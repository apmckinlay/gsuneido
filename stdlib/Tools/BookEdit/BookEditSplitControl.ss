// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
CodeViewControl
	{
	addons: #(
		Addon_html:,
		Addon_html_edit:,
		Addon_html_display:,
		Addon_show_whitespace:,
		Addon_go_to_help_page:,
		Addon_svc:,
		Addon_multiple_selection:,
		Addon_move_lines:,
		Addon_show_modified_lines:,
		Addon_show_references:,
		Addon_suneido_style: false,
		Addon_show_margin: false,
		Addon_brace_match: false,
		Addon_class_outline: false,
		Addon_speller: (ignore: (dl, dt, dd, ul, ol, li, href, br, pre,
			suneido, Suneido, gSuneido, gsuneido, jSuneido, jsuneido,
			https, stdlib, builtin, Builtin)))
	New(data, readonly = false)
		{
		super(:data, addons: .addons, :readonly)
		}

	InitialSet(data)
		{
		table = data.table
		// Need to lookup the current data on construct to ensure it is up to date
		data = data.num is 0 // Tree Book root "folder"
			? [name: data.table, plugin:] // Fake record for the "root" folder of the tree
			: Query1(data.table, num: data.num)
		data.table = table
		super.InitialSet(data)
		}

	Set(data)
		{
		data.name = SvcBook.MakeName([name: data.name, path: data.path])
		.plugin? = data.plugin is true
		if BookResource?(data.name, readOnly?:)
			{
			.Editor.Ymin = 0
			.Editor.SetVisible(false)
			.GetSplit().SetSplit([0, 1])
			// Required for book resources as they do NOT go through the parent's .Set
			.Table = data.GetDefault(#table, .Table)
			.RecName = data.GetDefault(#name, .RecName)
			}
		else
			super.Set(data)
		}

	MenuSelect(tip)
		{ .Status(tip) }

	Goto(name /*unused*/, source /*unused*/)
		{ return 0 /* need this because BookEditControl has a Goto */ }

	Deletable?()
		{
		return .plugin?
			? 'Cannot delete, item is from a pluigin'
			: ''
		}
	}
