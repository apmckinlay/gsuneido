// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 'ReporterColumns'
	ComponentName: 'ReporterColumns'

	New()
		{
		.hdr = SuJsListHeader(headerSelectPrompt: 'no_prompts')
		.char = 700 / Reporter.LandscapeChars
		.ComponentArgs = #()
		}

	cols: ()
	Get()
		{
		return .cols
		}

	Set(cols)
		{
		.cols = cols
		.hdr.Clear()
		for col in cols
			{
			width = col.width * .char
			.hdr.AddItem(col.text, width.Round(0),
				tip: col.text $ ' (click to modify)')
			}
		.Act(#UpdateHead, .hdr.Get())
		}

	HeaderResize(col, width)
		{
		.hdr.SetItemWidth(col, width)
		.Act("SetColWidth", col, width)
		wdth = width / .char
		.cols[col].width = wdth.Round(0)
		}

	HeaderReorder(item, newpos)
		{
		col = .cols[item]
		.cols.Delete(item)
		.cols.Add(col, at: newpos)
		.hdr.Reorder(item, newpos)
		.Act(#UpdateHead, .hdr.Get())
		}

	HeaderClick(col)
		{
		.Send('Click', .hdr.GetItem(col).text)
		}

	CONTEXTMENU_HEADER(x, y, col)
		{
		.cur = col
		hdrSize = .hdr.GetItemCount()
		menu = Object('Add/Remove Columns')
		if hdrSize > 0 and .cur < hdrSize
			menu.Add('Remove Column', '', 'Properties')
		ContextMenu(menu).ShowCall(this, x, y)
		}

	ContextMenu(x, y)
		{
		menu = Object('Add/Remove Columns')
		ContextMenu(menu).ShowCall(this, x, y)
		}

	On_Context_AddRemove_Columns()
		{
		.Send('On_AddRemove_Columns')
		}

	On_Context_Remove_Column()
		{
		if .hdr.GetItemCount() <= 0
			return
		.Send('ClearProperties', .hdr.GetItem(.cur).text)
		.hdr.DeleteItem(.cur)
		.cols.Delete(.cur)
		.Act(#UpdateHead, .hdr.Get())
		}

	On_Context_Properties()
		{
		if .hdr.GetItemCount() <= 0
			return
		.Send('Properties', .hdr.GetItem(.cur).text)
		}
	}