// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	// NOTE: assumes white background (i.e. WndPane)
	Name: ReporterColumns
	Xmin: 700
	New()
		{
		.hdr = .Vert.Header
		.colheads = .Vert.ReporterColHeads

		.char = .Xmin / Reporter.LandscapeChars
		}
	Controls()
		{
		return Object('Vert'
			'ReporterColHeads'
			Object('Header', style: HDS.BUTTONS | HDS.DRAGDROP,
				headerSelectPrompt: 'no_prompts')
			)
		}

	itemClicked: false
	HeaderClick(item, button /*unused*/)
		{
		.itemClicked = item
		// have to delay the Clicked message in order to allow the item tracking process
		// to end before the Controller does anything with the Clicked message
		// (like create a dialog). Creating a dialog in the middle of the header
		// tracking process messes up the focus and the header
		// tracking gets stuck like the user is dragging the column.
		.Defer(.sendItemClickedMessage)
		}
	sendItemClickedMessage()
		{
		.Send('Click', .hdr.GetItem(.itemClicked).text)
		}
	HeaderReorder(item, newpos)
		{
		col = .cols[item]
		.cols.Delete(item)
		.cols.Add(col, at: newpos)
		.columnsChanged()
		}
	HeaderResize(i, width /*unused*/)
		{
		wdth = (.hdr.GetItem(i).width / .char)
		.cols[i].width = wdth.Round(0)
		.columnsChanged()
		}
	columnsChanged()
		{
		.colheads.SetColumns(.getHeadings(), .getWidths())
		}
	getHeadings()
		{
		// TODO: headings
		return .cols.Map({ it.text })
		}
	getWidths()
		{
		return .cols.Map({ (it.width * .char).Round(0) })
		}

	ContextMenu(x, y)
		{
		ScreenToClient(.hdr.Hwnd, pt = Object(:x, :y))
		.cur = .hdr.HitTest(pt.x, pt.y).iItem
		menu = Object('Add/Remove Columns')
		if .cur isnt -1
			menu.Add('Remove Column', '', 'Properties')
		ContextMenu(menu).ShowCall(this, x, y)
		}
	On_Context_AddRemove_Columns()
		{
		.Send('On_AddRemove_Columns')
		}
	On_Context_Remove_Column()
		{
		.Send('ClearProperties', .hdr.GetItem(.cur).text)
		.hdr.DeleteItem(.cur)
		.cols.Delete(.cur)
		.columnsChanged()
		}
	On_Context_Properties()
		{
		.Send('Properties', .hdr.GetItem(.cur).text)
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
			width = (col.width * .char)
			.hdr.AddItem(col.text, width.Round(0),
				tip: col.text $ ' (click to modify)')
			}
		.columnsChanged()
		}
	}
