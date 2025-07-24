// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
// TODO: num/number/date fields do not show tootip
class
	{
	New(.grid)
		{
		// FIXME: This control does not apparently destroy its tooltip.
		grid.Map = Object()
		grid.Map[TTN.SHOW] = 'TTN_SHOW'
		.tip = grid.Construct(ToolTipControl)
		.tip.SendMessage(TTM.SETMAXTIPWIDTH, 0, 400) /*= tip max width */
		.tip.SendMessage(TTM.SETDELAYTIME, TTDT.AUTOPOP, 30000) /*= tip delay time*/
		.tip.Activate(false)
		.tip.AddTool(grid.Hwnd, LPSTR_TEXTCALLBACK)
		.tip.SetFont(StdFonts.Mono())
		grid.SetRelay(.tip.RelayEvent)
		}

	last_row: false
	last_col: false
	UpdateToolTip(row_num, col, model)
		{
		if row_num is .last_row and col is .last_col
			return

		.tip.Activate(false)
		rec = model.GetRecord(row_num)
		if rec isnt false and col isnt false
			.showTooltip(col, rec)
		.last_row = row_num
		.last_col = col
		}

	showTooltip(col, rec)
		{
		if rec.GetDefault('vl_full_display', #()).Member?(col)
			{
			.tip.UpdateTipText(.grid.Hwnd, rec.vl_full_display[col])
			.tip.Activate(true)
			}
		}

	Activate(status)
		{
		.tip.Activate(status)
		}

	TTN_SHOW(lParam, getCellRect)
		{
		if .last_row is false or .last_col is false
			return true

		r = getCellRect(.last_row, .last_col)
		.tip.AdjustRect(false, r)
		ClientToScreen(.grid.Hwnd, p = [x: r.left, y: r.top])
		SetWindowPos(NMHDR(lParam).hwndFrom, 0,
			p.x, p.y, 0, 0, // rect
			SWP.NOACTIVATE | SWP.NOSIZE | SWP.NOZORDER)
		return true
		}

	ResetLast()
		{
		.last_row = false
		.last_col = false
		}

	Destroy()
		{
		.tip.Destroy()
		}
	}
