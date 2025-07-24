// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
CursorWndProc
	{
	Name:		BookMarkContainer
	Xstretch:	1
	Ystretch:	1
	active:		false
	highlight:	false
	depressed:	false
	scroll:		0
	New()
		{
		super(.Name, WS.VISIBLE | WS.CLIPSIBLINGS)
		.SetFont()
		.WithSelectObject(.GetFont(), {|hdc| GetTextMetrics(hdc, tm = Object()) })
		.aveCharWidth = tm.AveCharWidth
		.marks = Object()
		.custColors = Object()
		.tips = .Construct(#ToolTip)
		.tools = .Construct(BookMarkToolControl)
		.yoffset = ScaleWithDpiFactor(6 /*= DPI factor*/) + .tools.Ymin
		.drawing = Object(
			shadbrush:	GetSysColorBrush(COLOR.TRIDSHADOW),
			borderPen: CreatePen(PS.SOLID, 0, GetSysColor(COLOR.TRIDSHADOW)),
			activeMarkFont: CreateFontIndirect(.LogFont(weight: FW.BOLD)),
			highlightMarkFont: CreateFontIndirect(.LogFont(underline:)))
		.tips.AddTool(.Hwnd, '')
		.SetRelay(.tips.RelayEvent)
		}

	GotoPath(path)
		{
		if .active is x = .findMark(path)
			return
		.repaintMark(.active) 		// Repaint the previous active bookmark
		.repaintMark(.active = x) 	// Repaint the new active bookmark
		if .active isnt false
			.ensureMarkVisible(.active)
		}

	findMark(path)
		{
		return .marks.FindIf({ it.path is path })
		}

	repaintMark(i)
		{
		if false isnt rc = .getMarkRect(i)
			InvalidateRect(.Hwnd, rc, true)
		}

	getMarkRect(index)
		{
		if not .getMarkRect?(index, .w, .marks.Size())
			return false
		GetClientRect(.Hwnd, rc = Object())
		top = .yoffset + (index - .scroll) * .markHeight
		return top > rc.bottom or top < .yoffset
			? false
			: Object(left: 1, right: rc.right, :top, bottom: top + .markHeight)
		}

	getMarkRect?(index, w, marks)
		{
		return index isnt false and w > 0 and index < marks
		}

	getter_markHeight()
		{
		return .markHeight = (.TextExtent(#M).y * 1.25).Round(0) + 2 /*= padding*/
		}

	AddMark(path)
		{
		if path is ''
			return
		if false is i = .findMark(path)
			{
			.marks.Add(Object(:path, color: CLR.YELLOW))
			.active = i = (nPos = .marks.Size()) - 1
			.updateScrollBar(nPos)
			.repaintMarks()
			}
		else
			.ensureMarkVisible(i)
		}

	repaintMarks()
		{
		GetClientRect(.Hwnd, rc = Object())
		rc.top += .yoffset
		InvalidateRect(.Hwnd, rc, true)
		}

	ensureMarkVisible(i)
		{
		if i < .scroll or i >= .maxVisibleMarks
			.setScroll(nPos: i)
		}

	RemoveMark(path)
		{
		if false is x = .findMark(path)
			.AlertInfo('Remove Mark', 'There is no bookmark currently selected.')
		else
			{
			if x is .active
				.active = false
			.marks.Delete(x)
			.updateScrollBar(nPos: .shiftScroll(x > .scroll))
			.repaintMarks()
			}
		return x
		}

	shiftScroll(down? = false)
		{
		return .scroll + (down? ? 1 : -1)
		}

	Resize(.x, .y, .w, .h)
		{
		super.Resize(x, y, w, h)
		.tools.Resize(x: 0, y: 3, :w, h: .tools.Ymin) // Move toolbar
		.updateScrollBar(nPos: .scroll)
		}

	maxVisibleMarks: 0
	updateScrollBar(nPos = false)
		{
		GetClientRect(.Hwnd, rc = Object())
		.maxVisibleMarks = ((rc.bottom - .yoffset) / .markHeight).RoundDown(0)
		.maxScroll = .marks.Size() - .maxVisibleMarks
		resize? = .maxScroll <= 0
			? .removeScrollBar()
			: .addScrollBar(nPos)
		if resize?
			.Resize(.x, .y, .w, .h)
		}

	removeScrollBar()
		{
		.scroll = 0
		if not .HasStyle?(WS.VSCROLL)
			return false
		.setScroll(fMask: SIF.RANGE | SIF.PAGE, nPage: 0, nMin: 0, nMax: 0)
		.RemStyle(WS.VSCROLL)
		return true
		}

	addScrollBar(nPos)
		{
		nPage = Max((.maxScroll / .markHeight).RoundDown(0), 1)
		sif = Object(fMask: SIF.RANGE | SIF.PAGE,
			:nPage, nMin: 0, nMax: .maxScroll + nPage - 1)
		if nPos isnt false
			sif.nPos = nPos
		.setScroll(@sif)
		if resize? = not .HasStyle?(WS.VSCROLL)
			.AddStyle(WS.VSCROLL)
		return resize?
		}

	setScroll(@sif)
		{
		repaint? = sif.Member?(#nPos)
			? .ensureScrollPos(sif)
			: false
		sif.cbSize = SCROLLINFO.Size()
		SetScrollInfo(.Hwnd, SB.VERT, sif, true)
		if repaint?
			.repaintMarks()
		return repaint?
		}

	ensureScrollPos(sif)
		{
		sif.fMask = sif.GetDefault(#fMask, 0) | SIF.POS
		sif.nPos = Max(Min(sif.nPos, .maxScroll), 0)
		repaint? = .scroll isnt sif.nPos
		.scroll = sif.nPos
		return repaint?
		}

	depressedMove: false
	MouseMove(x, y)
		{
		newHover = false
		.forSelectedMark(Object(:x, :y), { |mark, rc /*unused*/| newHover = mark })
		prevHighlight = .highlight
		if .depressed is false
			.mouseHoverNonSelected(newHover)
		else
			.mouseHoverSelected(.mouseDragScroll(y, newHover))
		if prevHighlight isnt .highlight
			.tips.UpdateTipText(.Hwnd, .toolTipText(.highlight))
		}

	mouseHoverNonSelected(newHover)
		{
		if newHover is .highlight
			return
		.repaintMark(.highlight)
		.repaintMark(.highlight = newHover)
		SetCursor(.highlight isnt false
			? LoadCursor(ResourceModule(), IDC.HAND) : LoadCursor(NULL, IDC.ARROW))
		}

	mouseDragScroll(y, newHover)
		{
		y -= .yoffset
		if y > 0 and y < .h - .yoffset
			return newHover
		.setScroll(nPos: .shiftScroll(down? = y >= .y))
		return Max(Min(.depressed + (down? ? 1 : -1), .marks.Size() - 1), 0)
		}

	mouseHoverSelected(newHover)
		{
		if .swapMark(newHover, .depressed) is true
			.depressedMove = true
		SetCursor(.depressed isnt false
			? LoadCursor(ResourceModule(), IDC.DRAG1) : LoadCursor(NULL, IDC.ARROW))
		}

	swapMark(index1, index2)
		{
		if .skipSwap?(index1, index2)
			return false
		.marks.Swap(index1, index2)
		.repaintMark(index1)
		.repaintMark(index2)
		.active = .swapIndex(.active, index1, index2)
		.depressed = .swapIndex(.depressed, index1, index2)
		return true
		}

	skipSwap?(index1, index2)
		{
		return index1 is index2 or index1 is false or index2 is false
		}

	swapIndex(curIndex, index1, index2)
		{
		return curIndex is index1
			? index2
			: curIndex is index2
				? index1
				: curIndex
		}

	toolTipText(i)
		{
		mark = .marks.GetDefault(i, [])
		text = mark.path
		if text is '' or RGBColors.BadContrast?(mark.color.ToRGB())
			return text
		GetClientRect(.Hwnd, rc = Object())
		textSize = (text.Size() * .aveCharWidth).RoundDown(0)
		return textSize > rc.right - rc.left
			? text
			: ''
		}

	MouseLeave()
		{
		.repaintMark(.highlight)
		.tips.UpdateTipText(.Hwnd, '')
		.highlight = false
		}

	PAINT()
		{
		BeginPaint(.Hwnd, ps = Object())

		start = Max(.scroll, 0)
		end = Min(.scroll + .maxVisibleMarks, .marks.Size())

		GetClientRect(.Hwnd, rc = Object())
		WithHdcSettings(ps.hdc, [.drawing.borderPen SetBkMode: TRANSPARENT],
			{
			.drawBookMarks(start, end, ps.hdc)
			MoveTo(ps.hdc, rc.left, rc.top + .yoffset)
			LineTo(ps.hdc, rc.left, rc.bottom)
			})
		EndPaint(.Hwnd, ps)
		return 0
		}

	drawBookMarks(start, end, hdc)
		{
		for i in start .. end
			if false isnt rc = .getMarkRect(i)
				{
				highlight? = .highlight?(i, .highlight, .depressed)
				depressed? = .depressed?(i, .highlight, .depressed)
				WithHdcSettings(hdc, .markSettings(i, highlight?))
					{ .drawBookMark(i, rc, hdc, highlight?, depressed?) }
				}
		}

	highlight?(i, highlight, depressed)
		{
		return (i is highlight and depressed is false) or i is depressed
		}

	depressed?(i, highlight, depressed)
		{
		return i is depressed and i is highlight
		}

	markSettings(i, highlight?)
		{
		font = i is .active
			? .drawing.activeMarkFont
			: highlight?
				? .drawing.highlightMarkFont
				: .Hwnd_hfont
		color = highlight?
			? COLOR.HIGHLIGHT
			: COLOR.BTNTEXT
		return Object(font, brush: .marks[i].color, SetTextColor: GetSysColor(color))
		}

	drawBookMark(i, rc, hdc, highlight?, depressed?)
		{
		text = .marks[i].path
		rc.left += i is .active or highlight? ? 0 : 9
		.shadowShading(hdc, rc, depressed?, highlight?)

		// Main surface
		offset = highlight? ? 1 : 0
		PatBlt(hdc, rc.left, rc.top,
			rc.right - rc.left - offset,
			rc.bottom - rc.top - offset,
			ROP.PATCOPY)

		rc.left += 2; rc.right -= 2; rc.top += 2; rc.bottom -= 2
		DrawTextEx(
			hdc,
			Paths.ToWindows(text), // PATH_ELLIPSIS requires backslashes
			text.Size(),
			rc,
			DT.PATH_ELLIPSIS | DT.VCENTER | DT.NOPREFIX | DT.NOCLIP,
			NULL
			)
		}

	shadowShading(hdc, rc, depressed?, highlight?)
		{
		FrameRect(hdc, rc, .drawing.shadbrush)
		if not (depressed? and highlight?)
			{
			rc.bottom--
			rc.right--
			}
		else
			{
			rc.top++
			rc.left++
			}
		}

	LBUTTONDBLCLK()
		{
		.AddMark(.Controller.CurrentPage())
		return 0
		}

	LBUTTONDOWN(lParam)
		{
		pt = Object(x: LOSWORD(lParam) y: HISWORD(lParam))
		.lButtonDown(pt)
		return 0
		}

	lButtonDown(pt)
		{
		.forSelectedMark(pt)
			{ |i, rc|
			.depressed = i
			InvalidateRect(.Hwnd, rc, false)
			}
		SetCapture(.Hwnd)
		}

	forSelectedMark(pt, block)
		{
		for (mark = .scroll; mark < .marks.Size(); ++mark)
			{
			if false is rc = .getMarkRect(mark)
				break
			else if PtInRect(rc, pt)
				{
				block(mark, rc)
				break
				}
			}
		}

	RBUTTONDOWN(lParam)
		{
		.forSelectedMark(Object(x: LOSWORD(lParam) y: HISWORD(lParam)))
			{ |i, rc|
			color = .marks[i].color
			result = ChooseColorWrapper(color, .Hwnd, custColors: .custColors)
			if result isnt false
				.marks[i].color = result
			InvalidateRect(.Hwnd, rc, true)
			}
		return 0
		}

	LBUTTONUP(lParam)
		{
		ReleaseCapture()
		if false isnt r = .getMarkRect(.depressed)
			.lbuttonUp(lParam, r)
		.depressedMove = .depressed = false // Unflag depressed button
		return 0
		}

	lbuttonUp(lParam, r)
		{
		bookmark = .marks[.depressed].path
		if false is .Send(#ValidBookmark?, bookmark)
			{
			.RemoveMark(bookmark)
			.AlertInfo('Bookmark Removed', 'Page ' $ bookmark $ ' not found.\n' $
				'Page may have been deleted or renamed.\nBookmark has been removed')
			return
			}
		if .depressedMove is false and
			PtInRect(r, Object(x: LOSWORD(lParam) y: HISWORD(lParam)))
			.Send(#Goto, bookmark)
		InvalidateRect(.Hwnd, r, true)	// Paint button in 'up' position
		}

	VSCROLL(wParam, lParam /*unused*/)
		{
		action = LOWORD(wParam)
		if action in (SB.LINEUP, SB.LINEDOWN, SB.THUMBTRACK)
			.setScroll(nPos: action is SB.THUMBTRACK
				? HIWORD(wParam)
				: .shiftScroll(action is SB.LINEDOWN))
		return 0
		}

	MOUSEWHEEL(wParam)
		{
		if .maxVisibleMarks >= .marks.Size()
			return 0
		.setScroll(nPos: .shiftScroll(GetWheelScrollInfo(wParam).down?))
		.tips.UpdateTipText(.Hwnd, '')
		return 0
		}

	GetState()
		{
		return .marks
		}

	SetState(stateobject)
		{
		.marks = Object()
		stateobject.Each()
			{
			if not String?(it)
				.marks.Add(it)
			else
				.AddMark(it)
			}
		.GotoPath(.Controller.CurrentPage())
		}

	Destroy()
		{
		.tips.Destroy()
		.tools.Destroy()
		.drawing.Each(DeleteObject)
		super.Destroy()
		}
	}
