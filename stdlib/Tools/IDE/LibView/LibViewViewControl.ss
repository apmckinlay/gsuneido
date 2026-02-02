// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
CodeViewControl
	{
	Name: LibViewView
	New(data, readonly = false)
		{
		super(:data, addons: .Addons(), :readonly)
		}

	addons: #(
		Addon_status: #(addon: Addon_editor_status), // Replacing base status bar
		Addon_libview_todo:,
		Addon_auto_complete_code:,
		Addon_auto_delimit:,
		Addon_auto_indent:,
		Addon_check_ascii:,
		Addon_check_code:,
		// vvvvvvvvvvvvvvvvv TEMPORARY: remove when 32381 is completed vvvvvvvvvvvvvvvvv
		Addon_check_lib_invalid_text:,
		// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
		Addon_comment:,
		Addon_star_rating:,
		Addon_show_whitespace:,
		Addon_unwrap:,
		Addon_debugger:,
		Addon_stepping_debugger:,
		Addon_multiple_selection:,
		Addon_move_lines:,
		Addon_svc:,
		Addon_show_modified_lines:,
		Addon_show_references:,
		Addon_inspect:,
		Addon_annotation:,
		Addon_coverage:)
	Addons()
		{
		addons = .addons.Copy()
		if addons.GetDefault('Addon_check_code', false) isnt false and
			addons.GetDefault('Addon_libview_todo', false) isnt false
			addons.Addon_libview_todo = #(init_qctext: 'Checking Code...')
		return addons
		}

	InitialSet(data)
		{
		table = data.table
		num = data.GetDefault('keyNum', 0)
		// Need to lookup the current data on construct to ensure it is up to date
		data = num is 0 // Tree Library root "folder"
			? [group: 0, keyNum: num] // Fake record for the "root" folder of the tree
			: Query1(table, :num)
		data.group = data.group isnt -1
		data.table = table
		data.keyNum = num
		super.InitialSet(data)
		}

	// Group and Num are used by: stdlib:Addon_libview_todo
	Group:	false
	Num: 	false
	Set(data)
		{
		.Group = data.group
		.Num = data.GetDefault('keyNum', 0)
		data.text = data.lib_current_text
		super.Set(data)
		}
	}
