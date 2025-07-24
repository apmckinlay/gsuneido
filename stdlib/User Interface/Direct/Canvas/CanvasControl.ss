// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: Canvas
	Title: Canvas
	Xstretch: 1
	Ystretch: 1
	New()
		{
		.CreateWindow("SuWhiteArrow", "", WS.VISIBLE)
		.SubClass()
		.items = Object()
		.selected = Object()
		.offset = ScaleWithDpiFactor(15) /*= dpi factor offset*/
		.paste_offset = 0
		// Get id of Canvas clipboard format
		.clipformat = RegisterClipboardFormat("Suneido_CANVAS")
		.formatting = DrawFormatting()
		}
	AddItem(item)
		{
		.items.Add(item)
		.Send('CanvasChanged')
		InvalidateRect(.Hwnd, .RectConversion(item.BoundingRect()), true)
		}
	AddItemAndSelect(item)
		{
		.items.Add(item)
		.Select(.items.Find(item))
		.Send('CanvasChanged')
		InvalidateRect(.Hwnd, .RectConversion(item.BoundingRect()), true)
		}
	MoveToBack(item)
		{
		.items.Remove(item)
		.items.Add(item, at: 0)
		.Select(0)
		.Send('CanvasChanged')
		InvalidateRect(.Hwnd, .RectConversion(item.BoundingRect()), true)
		}
	MoveToFront(item)
		{
		.items.Remove(item)
		.items.Add(item)
		.Select(.items.Find(item))
		.Send('CanvasChanged')
		InvalidateRect(.Hwnd, .RectConversion(item.BoundingRect()), false)
		}
	RemoveItem(item)
		{
		.items.Remove(item)
		.selected.Remove(item)
		.Send('CanvasChanged')
		InvalidateRect(.Hwnd, .RectConversion(item.BoundingRect()), true)
		}
	ResetSize(item)
		{
		// check if item is grouped
		if item.Grouped?
			{
			Alert('You cannot resize grouped items.', title: 'Error',
				flags: MB.ICONERROR)
			return
			}
		oriRect = .RectConversion(item.BoundingRect())
		item.ResetSize()
		.Select(.items.Find(item))
		.Send('CanvasChanged')
		InvalidateRect(.Hwnd, oriRect, true)
		InvalidateRect(.Hwnd, .RectConversion(item.BoundingRect()), true)
		}
	DeleteAll()
		{
		.items = Object()
		.selected = Object()
		.Send('CanvasChanged')
		.InvalidateClientArea()
		}
	RectConversion(rect)
		{
		rc = Object()
		rc.top = -CanvasItem.HandleBoundary() + (rect.y1 <= rect.y2 ? rect.y1 : rect.y2)
		rc.left = -CanvasItem.HandleBoundary() + (rect.x1 <= rect.x2 ? rect.x1 : rect.x2)
		rc.bottom = CanvasItem.HandleBoundary() + (rect.y1 > rect.y2 ? rect.y1 : rect.y2)
		rc.right = CanvasItem.HandleBoundary() + (rect.x1 > rect.x2 ? rect.x1 : rect.x2)
		return rc
		}
	GetSelected()
		{
		return .selected
		}
	GetAllItems()
		{
		return .items
		}
	SelectPoint(x, y)
		{
		.ClearSelect()
		if false isnt i = .ItemAtPoint(x, y)
			.Select(i)
		}
	ItemAtPoint(x, y)
		{
		for item in .selected
			if item.Contains(x, y)
				return .items.Find(item)
		for (i = .items.Size() - 1; i >= 0; --i)
			if .items[i].Contains(x, y)
				return i
		return false
		}
	SelectRect(x1, y1, x2, y2)
		{
		.MaybeClearSelect()
		for item in .items
			if item.Overlaps?(x1, y1, x2, y2)
				{
				.selected.Add(item)
				InvalidateRect(.Hwnd, .RectConversion(item.BoundingRect()), true)
				}
		}
	MaybeClearSelect()
		{
		if not KeyPressed?(VK.CONTROL) and not KeyPressed?(VK.SHIFT)
			{
			.ClearSelect()
			return true
			}
		else
			return false
		}
	SelectAll()
		{
		.ClearSelect()
		for item in .items
			.selected.Add(item)
		.InvalidateClientArea()
		}
	ClearSelect()
		{
		for item in .selected
			InvalidateRect(.Hwnd, .RectConversion(item.BoundingRect()), true)
		.selected = Object()
		}
	Select(i)
		{
		.selected.Add(.items[i])
		InvalidateRect(.Hwnd, .RectConversion(.items[i].BoundingRect()), true)
		}
	UnSelect(i)
		{
		.selected.Remove(.items[i])
		InvalidateRect(.Hwnd, .RectConversion(.items[i].BoundingRect()), true)
		}
	MoveSelected(dx, dy)
		{
		if .selected.Empty?()
			return
		rects = Object()
		for item in .selected
			{
			r = item.BoundingRect()
			rect = Object(left: r.x1, right: r.x2, top: r.y1, bottom: r.y2, :item)
			rects.Add(rect)
			}
		move = DrawSelectTracker_CalcNextMove(dx, dy, rects, this)
		for item in .selected
			item.Move(move.x, move.y)
		.InvalidateClientArea()
		.Send('CanvasChanged')
		}
	DeleteSelected()
		{
		for item in .selected
			{
			.items.Remove(item)
			InvalidateRect(.Hwnd, .RectConversion(item.BoundingRect()), true)
			item.Destroy()
			}
		.selected = Object()
		.Send('CanvasChanged')
		}
	CopyItems()
		{
		if .selected.Empty?()
			return
		.paste_offset = .offset
		ClipboardWriteString(Base64.Encode(Pack(.Get(.selected))), .clipformat)
		}

	Get(items = false)
		{
		if items is false
			items = .items
		return items.Map({ it.Get() }).Copy()
		}

	FormatColor(color)
		{
		".SetColor(0x" $ color.Hex() $ ")"
		}
	FormatLineColor(color)
		{
		".SetLineColor(0x" $ color.Hex() $ ")"
		}
	PasteItems()
		{
		if not IsClipboardFormatAvailable(.clipformat)
			return
		.ClearSelect()
		if '' is s = ClipboardReadString(.clipformat)
			return

		items = Unpack(Base64.Decode(s))
		_canvas = this
		for itemDef in items
			{
			item = Construct(itemDef).SetupScale()
			item.Move(.paste_offset, .paste_offset)
			.AddItemAndSelect(item)
			}
		.paste_offset += .offset
		.Send('CanvasChanged')
		.InvalidateClientArea()
		}

	color: 16777215
	SetColor(color)
		{
		.color = color
		if not .selected.Empty?()
			{
			for item in .selected
				{
				item.SetColor(color)
				InvalidateRect(.Hwnd, .RectConversion(item.BoundingRect()), true)
				}
			.Send('CanvasChanged')
			}
		}
	GetColor()
		{
		return .color
		}
	lin_color: 0
	SetLineColor(color)
		{
		.lin_color = color
		if not .selected.Empty?()
			{
			for item in .selected
				{
				item.SetLineColor(color)
				InvalidateRect(.Hwnd, .RectConversion(item.BoundingRect()), true)
				}
			.Send('CanvasChanged')
			}
		}
	GetLineColor()
		{
		return .lin_color
		}
	hfont: false
	SetFont(font, size = 11, weight = 500, underline = false)
		{
		.destroy_hfont()
		hdc = GetDC(.Hwnd)
		charset = GetLanguage().charset
		lf = Object(
			lfFaceName: font,
			lfHeight: size * GetDeviceCaps(hdc, GDC.LOGPIXELSY) / PointsPerInch,
			lfWeight: weight,
			lfUnderline: underline,
			lfEscapement: 0,
			lfOrientation: 0,
			lfCharSet: CHARSET[charset])
		.hfont = CreateFontIndirect(lf)
		.Send('CanvasChanged')
		}
	Resize(x, y, .w, .h)
		{
		super.Resize(x, y, w, h)
		}
	w: false
	h: false
	GetWidth()
		{
		return .w
		}
	GetHeight()
		{
		return .h
		}
	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		.formatting.SetDC(hdc)
		_canvas = this
		WithBkMode(hdc, TRANSPARENT)
			{
			.formatting.Paint(.items, ps.rcPaint)
			invalidItems = Object()
			for item in .selected
				{
				item.PaintHandles(hdc)
				if false is item.GetDefault('Valid', true)
					invalidItems.Add(item)
				}
			for item in invalidItems
				.RemoveItem(item)
			}
		EndPaint(.Hwnd, ps)
		return 0
		}
	DoWithReport(block)
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		.formatting.SetDC(hdc)
		_report = .formatting
		return Finally(block, { EndPaint(.Hwnd, ps) })
		}
	destroy_hfont()
		{
		if .hfont isnt false
			DeleteObject(.hfont)
		.hfont = false
		}
	SetReadOnly(readOnly)
		{
		.ClearSelect()
		SendMessage(.Hwnd, EM.SETREADONLY, readOnly, 0)
		super.SetReadOnly(readOnly)
		}
	// On_Copy is redirected to here (focus)
	// other keyboard shortcuts are handled in DrawCanvasControl
	On_Copy()
		{ .CopyItems() }

	SetXminYmin(.Xmin, .Ymin)
		{
		.WindowRefresh()
		}
	SyncItem(@unused) { }
	DESTROY()
		{
		for item in .items
			item.Destroy()
		.destroy_hfont()
		.formatting.Destroy()
		return 0
		}
	}
