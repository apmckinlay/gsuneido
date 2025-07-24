// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Diff to Original'
	CallClass(lib, name)
		{
		if false isnt local_rec = .local_rec(lib = lib.Tr('()'), name)
			return ToolDialog(0, [this, lib, name, local_rec])
		return false
		}

	local_rec(lib, name)
		{
		if false is local_rec = Query1(lib, :name, group: -1)
			{
			.AlertInfo(.Title, lib $ ':' $ name $ ' - cannot find local record')
			return false
			}
		if local_rec.lib_modified is ''
			{
			.AlertInfo(.Title, lib $ ':' $ name $ ' - has no local changes')
			return false
			}
		return local_rec
		}

	New(.lib, .name, local_rec)
		{ super(.layout(local_rec)) }

	layout(local_rec)
		{
		extraControls = Object(#Horz)
		comment = .lib $ ':' $ .name
		commentBgColor = false
		extraControls.Add('Skip', #(MenuButton, 'Restore', 'Restore'))
		return Object('Vert',
			Object('Diff2', local_rec.lib_current_text, local_rec.lib_before_text,
				.lib, .name, 'Modified', 'Original', :comment, :commentBgColor,
				extraControls: #(MenuButton, 'Restore', 'Restore')))
		}

	On_Restore(unused)
		{ .Window.Result('Restore') }
	}