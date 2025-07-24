// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	HorzMargin: 	5
	model:			false
	New(.grid, .rowHeight)
		{
		.createSystemObjects()
		}

	createSystemObjects()
		{
		.background = GetSysColorBrush(COLOR.WINDOW)
		.focused = GetSysColorBrush(COLOR.HIGHLIGHT)
		.alwaysHighlightSelected? = .grid.AlwaysHighlightSelected?()
		.invalid = CreateSolidBrush(CLR.LIGHTRED)
		.shaded = CreateSolidBrush(CLR.azure) /*= shaded */
		.warning = CreateSolidBrush(CLR.WarnColor)

		.expandBackground = GetSysColorBrush(COLOR.BTNFACE)
		.expandLine = ExtCreatePen(PS.GEOMETRIC | PS.DOT, 1,
			Object(lbColor: 0x00000000, lbStyle: BS.SOLID), 0, 0)
		}

	SetModel(.model)
		{
		.colModel = .model.ColModel
		}

	PAINT()
		{
		hdc = BeginPaint(.grid.Hwnd, ps = Object())
		if .model is false
			FillRect(hdc, ps.rcPaint, .background)
		else
			WithHdcSettings(hdc, [.grid.GetFont(), SetBkMode: TRANSPARENT],
				{ .paintGrid(hdc, ps.rcPaint) })
		EndPaint(.grid.Hwnd, ps)
		return 0
		}

	paintGrid(hdc, rc)
		{
		numCols = .colModel.GetSize()
		// find the start of column that need painting
		col = 0
		left = -.colModel.Offset
		while(col < numCols)
			{
			next_left = left + .colModel.GetColWidth(col)
			if next_left > rc.left
				break
			left = next_left
			++col
			}

		top = .paint(hdc, rc, col, left, numCols)

		if (top < rc.bottom)
			{
			rcSel = Object(left: rc.left, :top, right: rc.right, bottom: rc.bottom)
			FillRect(hdc, rcSel, .background)
			}
		}

	paint(hdc, rc, col, left, numCols)
		{
		// for each row needing painting
		.colModel.SetDC(hdc)
		row_num = .grid.GetRows(rc.top)
		top = row_num * .rowHeight
		prevExpand? = false
		while(row_num < .model.VisibleRows and top < rc.bottom)
			{
			if false is rec = .model.GetRecord(row_num)
				break

			if rec.vl_expand? is true
				{
				.paintExpandLines(hdc, top)
				top += .rowHeight
				prevExpand? = true
				++row_num
				continue
				}
			rcSel = Object(left: rc.left, :top, right: rc.right,
				bottom: top + .rowHeight)
			.fillRow(hdc, rcSel, row_num + .model.Offset, rec, prevExpand?)
			prevExpand? = false

			c = col
			x = left + .HorzMargin
			while(c < numCols and x < rc.right)
				{
				width = .paintCell(c, rec, hdc, top, x)
				x += width
				++c
				}

			top += .rowHeight
			++row_num
			}
		return top
		}

	paintCell(c, rec, hdc, top, x)
		{
		width = .colModel.GetColWidth(c)
		orig = false
		if width > 2 * .HorzMargin
			{
			if .model.EditModel.ColumnInvalid?(rec, .colModel.Get(c)) and
				.grid.GetReadOnly() isnt true
				{
				FillRect(hdc, Object(:top, left: x - .HorzMargin,
					right: x + width - .HorzMargin, bottom: top + .rowHeight), .invalid)
				orig = .model.ColModel.SetBackgroundBrush(.invalid, selected: '')
				}
			rect = Rect(x, top + 2, width - 2 * .HorzMargin, .rowHeight - 2)
			.colModel.PaintCell(c, rect, rec)
			}
		if orig isnt false
			.model.ColModel.SetBackgroundBrush(@orig)
		return width
		}

	paintExpandLines(hdc, top)
		{
		if .colModel.Offset >= .rowHeight
			return
		FillRect(hdc, Object(:top, left: 0,
			right: .rowHeight, bottom: top + .rowHeight), .expandBackground)
		}

	fillRow(hdc, rcSel, row, rowRec, prevExpand?)
		{
		bgBrush = .getBrush(row, rowRec, prevExpand?)
		FillRect(hdc, rcSel, bgBrush)
		selected = .model.Selection.HasSelectedRow?(rowRec)
		focused = bgBrush is .focused
		.model.ColModel.SetBackgroundBrush(bgBrush, selected: focused and selected)

		if selected
			{
			rcSel.left = -2 // only draw top and bottom dotted lines for focus
			rcSel.right = 20000
			DrawFocusRect(hdc, rcSel)
			}

		textColor = focused and selected ? COLOR.HIGHLIGHTTEXT : COLOR.WINDOWTEXT
		SetTextColor(hdc, GetSysColor(textColor))
		}

	brushMgr: false
	getBrush(row, rowRec, prevExpand? = false)
		{
		if .model.Selection.HasSelectedRow?(rowRec) and
			(.alwaysHighlightSelected? or .grid.HasFocus?())
			return  .focused

		if .model.EditModel.GetWarningMsg(rowRec) isnt ""
			return .warning

		if .brushMgr isnt false
			if false isnt brush = .brushMgr.GetBrush(rowRec)
				return brush

		return prevExpand? or row % 2 is 0 ? .shaded : .background
		}

	SetBackground(color)
		{
		DeleteObject(.background)
		.background = GetSysColorBrush(color)
		}

	HighlightValues(member, values, color)
		{
		if .brushMgr is false
			.brushMgr = VirtualListBrushes()
		.brushMgr.HighlightValues(member, values, color)
		}

	HighlightRecords(recs, color)
		{
		if .brushMgr is false
			.brushMgr = VirtualListBrushes()
		.brushMgr.HighlightRecords(recs, color)
		}

	ClearHighlightRecord(rec)
		{
		if .brushMgr isnt false
			.brushMgr.ClearHighlightRecord(rec)
		}

	ClearHighlight()
		{
		if .brushMgr isnt false
			.brushMgr.Destroy()
		.brushMgr = false
		}

	destroyObject(ob)
		{
		if ob isnt false
			DeleteObject(ob)
		}

	Destroy()
		{
		.destroyObject(.shaded)
		.destroyObject(.invalid)
		.destroyObject(.expandBackground)
		.destroyObject(.expandLine)
		.destroyObject(.warning)

		if .brushMgr isnt false
			.brushMgr.Destroy()
		}
	}
