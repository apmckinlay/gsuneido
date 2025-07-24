// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// primarily for use by Split
// a container that uses this must have a Movesplit(pos) function
/* This is a variation on the normal splitter control
   that can display a 'handle' that is used for 'opening' or 'closing' the splitter
   (This feature was designed for BookControl) */
SplitterControl
	{
	SplitName: "Splitter"
	tips: false
	New(tips = true)
		{
		.tips = tips ? .Construct("ToolTip") : tips
		.d = ScaleWithDpiFactor(4 /*= scale factor*/)
		.d2 = 2 * .d
		.h1 = ScaleWithDpiFactor(13 /*= scale factor*/)
		.h2 = 2 * (.h1 + .d)
		.createBitmaps()
		.handleUnderMouse? = false
		.SetRelay(.tips.RelayEvent)
		}

	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		rc = .getHandleRect()
		rc.bottom += rc.top
		rc.right += rc.left
		.tips.RemoveTool(.Hwnd, 0)
		.tips.AddTool(.Hwnd,
			(.Parent.Open? ? "Close" : "Open") $ " " $ .SplitName, rect: rc)
		}

	hdcs: ()
	bmps: ()
	createBitmaps()
		{
		// Create arrow bitmaps
		hdc = GetDC(.Hwnd)
		hpen_normal = CreatePen(PS.SOLID, 0, GetSysColor(COLOR.TRIDSHADOW))
		.hdcs = Object()
		.bmps = Object()
		.bmph = .createArrowBitmap(hdc, hpen_normal, .d, .d2) // Horizontal arrow
		.bmpv = .createArrowBitmap(hdc, hpen_normal, .d2, .d) // Vertical arrow
		DeleteObject(hpen_normal)
		ReleaseDC(.Hwnd, hdc)
		}

	createArrowBitmap(hdcthis, hpen, width, height)
		{
		/* CreateCompatibleDC creates a new DC in memory.
		As a result, the SelectObject is limited to the scope of said DC.
		This means we do not need to restore these SelectObjects.
		Furthermore, later code relies on the result of CreateCompatibleBitmap being
		in the DC's SelectObject. */
		.hdcs.Add(hdc = CreateCompatibleDC(hdcthis))
		.bmps.Add(hbmp = CreateCompatibleBitmap(hdcthis, width, height))
		SelectObject(hdc, hbmp)
		PatBlt(hdc, 0, 0, width, height, ROP.WHITENESS)
		SelectObject(hdc, hpen)
		if width < height
			.horizontalBitMap(hdc width, height)
		else
			.verticalBitMap(hdc width, height)
		return Object(:hdc, :hbmp)
		}

	horizontalBitMap(hdc width, height)
		{
		.iterateOverXY(width - 1, (height / 2).Int() - 1)
			{|x, y|
			MoveTo(hdc, x, y)
			LineTo(hdc,	x, (height - y))
			}
		}

	iterateOverXY(x, y, block)
		{
		for (; x >= 0 and y >= 0; x--, y--)
			block(x, y)
		}

	verticalBitMap(hdc, width, height)
		{
		.iterateOverXY((width / 2).Int() - 1, height - 1)
			{|x, y|
			MoveTo(hdc, x, y)
			LineTo(hdc, (width - x), y)
			}
		}

	MouseMove(wParam, lParam, x, y)
		{
		newHandleUnderMouse? = .mouseOnHandle?(x, y)
		if newHandleUnderMouse? isnt .handleUnderMouse?
			.repaintHandle()
		.handleUnderMouse? = newHandleUnderMouse?
		if not .SplitterControl_dragging and .handleUnderMouse?
			SetCursor(LoadCursor(ResourceModule(), .getCursor()))
		else if .Parent.CanDrag?
			return super.MouseMove(:wParam, :lParam, :x, :y)
		else
			SetCursor(LoadCursor(NULL, IDC.ARROW))
		return 0
		}

	MouseLeave()
		{
		if not .handleUnderMouse?
			return
		.handleUnderMouse? = false
		.repaintHandle()
		}

	LBUTTONDOWN(lParam)
		{
		if .handleUnderMouse? or not .Parent.CanDrag?
			return 0
		else
			return super.LBUTTONDOWN(lParam)
		}

	LBUTTONUP(lParam)
		{
		if not .SplitterControl_dragging and .handleUnderMouse?
			{
			ClientToScreen(.Hwnd, pt = Object(x: LOSWORD(lParam) y: HISWORD(lParam)))
			.LBUTTONDBLCLK()
			ScreenToClient(.Hwnd, pt)
			SendMessage(.Hwnd, WM.MOUSEMOVE, 0, pt.x | pt.y << 16)
			}
		else
			super.LBUTTONUP(lParam)
		return 0
		}

	LBUTTONDBLCLK()
		{
		if .Parent.Open?
			.Parent.Close()
		else
			.Parent.Open()
		return 0
		}

	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		GetClientRect(.Hwnd, r = Object())
		FillRect(hdc, r, GetSysColorBrush(COLOR.BTNFACE))

		if .handleUnderMouse? // Highlight entire rectangle when hovered
			WithHdcSettings(hdc, [brush: GetSysColor(COLOR.HIGHLIGHT)])
				{
				rc = .getHandleRect()
				PatBlt(hdc, rc.left, rc.top, rc.right, rc.bottom, ROP.PATCOPY)
				}

		data = .getArrowData()
		rop = .handleUnderMouse? ? ROP.MERGEPAINT : ROP.SRCAND
		if .Parent.Open?
			.paintArrowOpened(hdc, data.pt, data.bmp.hdc, rop)
		else
			.paintArrowClosed(hdc, data.pt, data.bmp.hdc, rop)
		EndPaint(.Hwnd, ps)
		return 0
		}

	paintArrowOpened(hdc, pt, bdc, rop)
		{
		switch .Parent.Associate
			{
		case "east":
			BitBlt(hdc, pt.x, pt.y, .d, .d2, bdc, 0, 0, rop)
		case "west":
			StretchBlt(hdc, pt.x + .d, pt.y, -.d, .d2, bdc, 0, 0, .d, .d2, rop)
		case "north":
			StretchBlt(hdc, pt.x, pt.y + .d, .d2, -.d, bdc, 0, 0, .d2, .d, rop)
		case "south":
			BitBlt(hdc, pt.x, pt.y, .d2, .d, bdc, 0, 0, rop)
			}
		}

	paintArrowClosed(hdc, pt, bdc, rop)
		{
		switch .Parent.Associate
			{
		case "east":
			StretchBlt(hdc, pt.x + .d, pt.y, -.d, .d2, bdc, 0, 0, .d, .d2, rop)
		case "west":
			BitBlt(hdc, pt.x, pt.y, .d, .d2, bdc, 0, 0, rop)
		case "north":
			BitBlt(hdc, pt.x, pt.y, .d2, .d, bdc, 0, 0, rop)
		case "south":
			StretchBlt(hdc, pt.x, pt.y + .d, .d2, -.d, bdc, 0, 0, .d2, .d, rop)
			}
		}

	repaintHandle()
		{
		rc = .getHandleRect()
		if .Dir is "vert"
			rc.right += rc.left
		else
			rc.bottom += rc.top
		InvalidateRect(.Hwnd, rc, true)
		}

	getHandleRect()
		{
		GetClientRect(.Hwnd, rcHandle = Object())
		pt = .getArrowPt()
		if .Parent.Dir is "vert"
			{
			rcHandle.left = pt.x - .h1
			rcHandle.right = .h2
			}
		else
			{
			rcHandle.top = pt.y - .h1
			rcHandle.bottom = .h2
			}
		return rcHandle
		}

	getCursor()
		{
		if false is cursorOb = .cursors.GetDefault(associate = .Parent.Associate, false)
			throw "invalid splitter Associate: " $ associate
		return .Parent.Open?
			? cursorOb.opened
			: cursorOb.closed
		}

	getter_cursors()
		{
		return .cursors = Object(
			east: 	Object(opened: IDC.HANDE, closed: IDC.HANDW),
			west: 	Object(opened: IDC.HANDW, closed: IDC.HANDE),
			south: 	Object(opened: IDC.HANDS, closed: IDC.HANDN),
			north: 	Object(opened: IDC.HANDN, closed: IDC.HANDS)).Set_readonly()
		}

	mouseOnHandle?(x, y)
		{
		rcHandle = .getHandleRect()
		if .Parent.Dir is "vert"
			rcHandle.right += rcHandle.left
		else
			rcHandle.bottom += rcHandle.top
		return PtInRect(rcHandle, Object(:x, :y))
		}

	getArrowPt()
		{
		GetClientRect(.Hwnd, rc = Object())
		pt = Object()
		if .Parent.Dir is "vert"
			{
			pt.x = rc.left + ((rc.right - rc.left - .d2) / 2).Int()
			pt.y = rc.top + ((rc.bottom - rc.top - .d) / 2).Int()
			}
		else
			{
			pt.x = rc.left + ((rc.right - rc.left - .d) / 2).Int()
			pt.y = rc.top + ((rc.bottom - rc.top - .d2) / 2).Int()
			}
		return pt
		}

	getArrowData()
		{
		data = Object(pt: .getArrowPt())
		if .Parent.Dir is "vert"
			data.bmp = .bmpv
		else
			data.bmp = .bmph
		return data
		}

	Destroy()
		{
		if .tips isnt false
			.tips.Destroy()
		.hdcs.Each(DeleteDC)
		.bmps.Each(DeleteObject)
		super.Destroy()
		}
	}
