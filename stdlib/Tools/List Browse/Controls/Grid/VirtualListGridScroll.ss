// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(grid, model, rowHeight, dx, dy, fromRowIndex = false, fromLeft = false,
		fromRight = false)
		{
		gridRect = grid.GetClientRect().ToWindowsRect()
		scroll = .scrollInfo(gridRect, fromLeft, fromRight, dx)
		lastRowBottom = gridRect.bottom - gridRect.bottom % rowHeight
		childRcs = .getChildRects(model)
		gridRc = GetWindowRect(grid.Hwnd)
		prevTop = fromRowIndex isnt false ? fromRowIndex * rowHeight : 0
		hwdp = BeginDeferWindowPos(childRcs.Size())
		for childRc in childRcs.Sort!(By('top'))
			{
			childRectTop = childRc.top - gridRc.top
			if .expandAboveIndex(fromRowIndex, childRectTop, rowHeight)
				{
				prevTop = Max(childRc.bottom - gridRc.top, prevTop)
				continue
				}
			subRect = Object(left: scroll.left, top: prevTop, right: scroll.right,
				bottom: Min(childRectTop, lastRowBottom))
			prevTop = childRc.bottom - gridRc.top
			if subRect.top < gridRect.bottom
				ScrollWindowEx(
					grid.Hwnd, scroll.dx, dy, subRect, subRect, NULL, NULL, SW.INVALIDATE)
			if .horzMove?(dy, model, rowHeight)
				{
				indentRect = Object(left: scroll.left, top: childRectTop,
					bottom: Min(prevTop, lastRowBottom), right: rowHeight)
				ScrollWindowEx(grid.Hwnd, 0, dy,
					indentRect, indentRect, NULL, NULL, SW.INVALIDATE)
				}
			if fromLeft is false
				hwdp = .moveExpand(grid, childRc, hwdp, scroll, dy)
			}
		.moveRemainder([:fromRowIndex, :prevTop, :lastRowBottom],
			rowHeight, scroll, grid, dy)
		EndDeferWindowPos(hwdp)
		}

	scrollInfo(gridRect, fromLeft, fromRight, dx)
		{
		return fromRight isnt false
			? Object(left: gridRect.left, right: fromRight, dx: -dx)
			: Object(left: fromLeft is false ? gridRect.left : fromLeft
				right: gridRect.right, :dx)
		}

	getChildRects(model)
		{
		childRcs = Object()
		if model is false or model.ExpandModel is false
			return childRcs
		expands = model.ExpandModel.GetControls()
		for c in expands
			{
			childHwnd = c.Hwnd
			childRc = GetWindowRect(childHwnd)
			childRc.hwnd = childHwnd
			childRcs.Add(childRc)
			}
		return childRcs
		}

	expandAboveIndex(fromRowIndex, childRectTop, rowHeight)
		{
		return fromRowIndex isnt false and childRectTop / rowHeight <= fromRowIndex
		}

	horzMove?(dy, model, rowHeight)
		{
		return dy isnt 0 and model.ColModel.Offset < rowHeight
		}

	moveExpand(grid, childRc, hwdp, scroll, dy)
		{
		ScreenToClient(grid.Hwnd, pt = Object(x: childRc.left, y: childRc.top))
		hwdp = DeferWindowPos(hwdp, childRc.hwnd, 0,
			pt.x + scroll.dx, pt.y + dy, 0, 0,
			SWP.NOSIZE | SWP.NOZORDER | SWP.NOACTIVATE)
		return hwdp
		}

	moveRemainder(remainder, rowHeight, scroll, grid, dy)
		{
		prevTop = remainder.prevTop
		fromRowIndex = remainder.fromRowIndex
		if fromRowIndex isnt false
			prevTop = Max(remainder.prevTop, fromRowIndex * rowHeight)
		if prevTop < remainder.lastRowBottom
			{
			botRect = Object(left: scroll.left, right: scroll.right,
				top: prevTop, bottom: remainder.lastRowBottom)
			ScrollWindowEx(
				grid.Hwnd, scroll.dx, dy, botRect, botRect, NULL, NULL, SW.INVALIDATE)
			}
		}

	MoveExpandControls(grid, model, rowHeight, toRowIndex, rows)
		{
		childRcs = .getChildRects(model)
		gridRc = GetWindowRect(grid.Hwnd)
		hwdp = BeginDeferWindowPos(childRcs.Size())
		for childRc in childRcs.Sort!(By('top'))
			{
			childRectTop = childRc.top - gridRc.top
			if toRowIndex * rowHeight > childRectTop
				{
				ScreenToClient(grid.Hwnd, pt = Object(x: childRc.left, y: childRc.top))
				hwdp = DeferWindowPos(hwdp, childRc.hwnd, 0,
					pt.x, pt.y + rows * rowHeight, 0, 0,
					SWP.NOSIZE | SWP.NOZORDER | SWP.NOACTIVATE)
				}
			}
		EndDeferWindowPos(hwdp)
		}

	HSCROLL(grid, model, colModel, wParam, rowHeight)
		{
		SetFocus(grid.Hwnd)
		dx = .horzScroll(grid, colModel, wParam)
		if dx is 0
			return 0

		.UpdateHorzScrollPos(grid, colModel)
		grid.Send('VirtualListGrid_HorzScroll', dx)
		.CallClass(grid, model, rowHeight, dx, 0)
		return 0
		}

	horzScroll(grid, colModel, wParam)
		{
		switch LOWORD(wParam)
			{
		case SB.LEFT:
			return .calcMovePixel(-colModel.Offset, grid, colModel)
		case SB.RIGHT:
			return .calcMovePixel(colModel.GetTotalWidths(), grid, colModel)
		case SB.LINELEFT:
			return .calcMovePixel(-10 /*= move left */, grid, colModel)
		case SB.LINERIGHT:
			return .calcMovePixel(10 /*= move right */, grid, colModel)
		default:
			return .horzScroll2(grid, colModel, wParam)
			}
		}

	horzScroll2(grid, colModel, wParam)
		{
		switch LOWORD(wParam)
			{
		case SB.PAGELEFT:
			return .calcMovePixel(-grid.GetClientRect().GetWidth(), grid, colModel)
		case SB.PAGERIGHT:
			return .calcMovePixel(grid.GetClientRect().GetWidth(), grid, colModel)
		case SB.THUMBTRACK:
			return .calcMovePixel(HIWORD(wParam) - colModel.Offset, grid, colModel)
		default:
			return 0
			}
		}

	calcMovePixel(movePix, grid, colModel)
		{
		newOffset = Max(0, Min(colModel.Offset + movePix,
			colModel.GetTotalWidths() - grid.GetClientRect().GetWidth()))
		movePix = colModel.Offset - newOffset
		colModel.Offset = newOffset
		return movePix
		}

	UpdateHorzScrollPos(grid, colModel)
		{
		nPage = grid.GetClientRect().GetWidth()
		totalWidth = colModel isnt false and colModel.StretchColumn is false
			? colModel.GetTotalWidths()
			: nPage - 1
		sif = Object(
			cbSize:	SCROLLINFO.Size(),
			fMask:	SIF.RANGE | SIF.POS | SIF.PAGE,
			:nPage,
			nMin:	0,
			nMax:	colModel is false ? 0 : totalWidth,
			nPos:	colModel is false ? 0 : colModel.Offset
			)
		SetScrollInfo(grid.Hwnd, SB.HORZ, sif, redraw:)
		}
	}
