// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
CanvasControl
	{
	SetTracker(tracker, item)
		{
		if .tracker isnt false
			.tracker.Release()
		.tracker = tracker(.Hwnd, item, canvas: this)
		}
	tracker: false
	dragging: false
	resizing?: false
	lParamDown: lParam
	CONTEXTMENU(lParam)
		{
		if .tracker is false or .dragging is true
			return 0
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		// select the item right clicked on
		ScreenToClient(.Hwnd, pt = Object(:x, :y))
		.tracker.MouseDown(pt.x, pt.y)
		.mouseUp(pt.x, pt.y)
		menu = DrawControl.Menu[0].Copy().Map({ it.Tr('&') })
		ContextMenu(menu).ShowCall(this, x, y)
		}
	Default(@args)
		{
		method = args[0]
		if method.Prefix?('On_Context_') and args.Member?('item')
			{
			.Send('On_' $ ToIdentifier(args.item))
			return 0
			}
		}
	RBUTTONDOWN(lParam /*unused*/)
		{
		.Send('WhichDrawCanvasClicked', .Name)
		return 0
		}
	LBUTTONDOWN(lParam)
		{
		.Send('WhichDrawCanvasClicked', .Name)
		if .tracker is false
			return 0
		SetFocus(.Hwnd)
		.dragging = true
		SetCapture(.Hwnd)
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		// check if user clicked on a handle
		// and using a select tracker
		// if they did, they are resizing
		if .tracker.Base?(DrawRectTracker)
			{
			for item in .GetAllItems()
				if item.IsHandle?(x, y)
					{
					.resizing? = true
					.prevtracker = .tracker
					.origx = x
					.origy = y
					.tracker.ResizeDown(item, x, y)
					return 0
					}
			}
		.lParamDown = lParam
		.tracker.MouseDown(x, y)
		return 0
		}
	MOUSEMOVE(lParam)
		{
		if not .dragging
			return 0
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		if .resizing? and (item = .GetSelected()).Size() isnt 0 and
			.tracker.Base?(DrawRectTracker)
			.tracker.ResizeMove(item[0], x, y)
		else
			.tracker.MouseMove(x, y)
		// prevent clear selected objects when dragging back to the original spot
		.lParamDown = -1
		return 0
		}
	LBUTTONUP(lParam)
		{
		if not .dragging
			return 0
		.dragging = false
		ReleaseCapture()
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		.selectPoint(lParam, x, y)
		if .resizing? and (item = .GetSelected()).Size() isnt 0 and
			.tracker.Base?(DrawRectTracker)
			{
			.resizing? = false
			.tracker.ResizeUp(item[0], x, y)
			.tracker = .prevtracker
			if item.Size() > 1
				{
				Alert('You must select one item', title: 'Error',
					flags: MB.ICONERROR)
				return 0
				}
			oriRect = .RectConversion(item[0].BoundingRect())
			item[0].Resize(.origx, .origy, x, y)
			InvalidateRect(.Hwnd, oriRect, true)
			InvalidateRect(.Hwnd, .RectConversion(item[0].BoundingRect()), false)
			.Send('CanvasChanged')
			return 0
			}
		if false isnt item = .mouseUp(x, y)
			{
			.ClearSelect()
			if not Object?(item)
				item = Object(item)
			item.Each({
				.AddItemAndSelect(it.SetColor(.GetColor()).SetLineColor(.GetLineColor()))
				})
			.resizing? = false
			.Send('Canvas_LButtonUp')
			}
		.Send('CanvasChanged')
		return 0
		}

	mouseUp(x, y)
		{
		.DoWithReport({ .tracker.MouseUp(:x, :y) })
		}

	selectPoint(lParam, x, y)
		{
		if .mouseSamePosition(lParam) and .MaybeClearSelect() and
			false isnt i = .ItemAtPoint(x, y)
			.Select(i)
		}

	mouseSamePosition(lParam)
		{
		return .lParamDown is lParam
		}
	LBUTTONDBLCLK()
		{
		.EditItem()
		return 0
		}
	GETDLGCODE(wParam, lParam)
		{
		keys = .Send('List_GetDlgCode', :wParam, :lParam)
		return keys is 0 ? DLGC.WANTCHARS | DLGC.WANTARROWS : keys
		}
	KEYDOWN(wParam)
		{
		if .arrowKeys(wParam)
			return 0
		if wParam is VK.DELETE
			.DeleteSelected()
		else if KeyPressed?(VK.CONTROL)
			.ctrlKeys(wParam)
		return 0
		}

	arrowKeys(wParam)
		{
		if not Object(VK.LEFT, VK.RIGHT, VK.UP, VK.DOWN).Has?(wParam)
			return false
		x = wParam.Even?() ? 0 : wParam is VK.LEFT ? -1 : 1
		y = wParam.Odd?()  ? 0 : wParam is VK.UP   ? -1 : 1
		.MoveSelected(x, y)
		return true
		}

	ctrlKeys(wParam)
		{
		// CTRL + A, X, V, Z do not work inside a Book unless we handle them here
		// Ctrl+C / On_Copy will be redirected to CanvasControl
		// so we don't handle it here
		switch wParam
			{
			case VK.A :
				.Send(#On_Select_All)
			case VK.X :
				.Send(#On_Cut)
			case VK.V :
				.Send(#On_Paste)
			case VK.Z:
				.Send(#On_Undo)
			case VK.Y:
				.Send(#On_Redo)
			default :
			}
		}

	EditItem()
		{
		item = .GetSelected()
		if item.Size() isnt 1
			return
		prevRect = item[0].BoundingRect()
		item[0].Edit(canvas: this)
		.Send('CanvasChanged')
		InvalidateRect(.Hwnd, .RectConversion(prevRect), true)
		InvalidateRect(.Hwnd, .RectConversion(item[0].BoundingRect()), true)
		}
	CutItems()
		{
		.CopyItems()
		.DeleteSelected()
		}
	SelectAll()
		{
		super.SelectAll()
		.tracker = DrawSelectTracker(.Hwnd, false, canvas: this)
		SetFocus(.Hwnd)
		}
	Destroy()
		{
		if .tracker isnt false
			.tracker.Release()
		super.Destroy()
		}
	}
