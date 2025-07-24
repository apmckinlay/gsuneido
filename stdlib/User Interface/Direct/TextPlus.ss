// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
/*
Base class for CheckBoxControl and RadioButtonControl
Not intended to be used directly

Testing:
- minimum box size matches Windows
	(this means it'll get cut off for extremely small font sizes)
- width is correct
	Dialog(0, #(Border (Horz (ColorRect) (CheckBox "ahoy") (ColorRect))))
- focus rectangle, just around text (should just touch 'y')
- no focus rectangle if text is ""
- blue box for mouse over or focus
- readonly gray box, no blue box
- checkmark centered in box
- different font sizes (box size should match font size)
	Dialog(0, #(Border (Vert
		(CheckBox Howdy size: 6)
		(CheckBox Howdy)
		(CheckBox Howdy size: 20)
		(CheckBox Howdy size: 30)
		(CheckBox Howdy size: 40)
		)))
- baseline alignment
	Dialog(0, #(Border (Horz (Static hello) (CheckBox hello))))
- overlap in FormControl
	Dialog(0, #(Border (Form
		(date group: 1) nl
		(etapay_readyto_pay? group: 1) nl
		(date group: 1)
		)))
- can click on box or text to toggle
- spacebar toggles
- can't modify if readonly
- works in HighlightControl and MainFieldControl
- all of the above work with lefttext
- from list/browse make sure double click to open control doesn't also toggle
*/
WndProc
	{
	valid: true
	New(.text = "", .lefttext = false, .font = "", .size = "", .weight = "",
		.readonly = false, .tip = "", hidden = false, tabover = false)
		{
		style = 0
		.SetHidden(hidden)
		if hidden is false
			style |= WS.VISIBLE
		if tabover is false
			style |= WS.TABSTOP
		.CreateWindow("SuBtnfaceArrow", "", style)
		.SubClass()
		.makeObs()
		.Recalc()
		}
	makeObs()
		{
		penThickness = .Antialias ? Max(4, (Number(.size) / 6)) : 1
		.obs = Object(
			grayPen: CreatePen(PS.INSIDEFRAME, penThickness, CLR.GRAY)
			bluePen: CreatePen(PS.INSIDEFRAME, penThickness, RGB(51, 153, 255))
			darkGrayPen: CreatePen(PS.INSIDEFRAME, penThickness, CLR.DARKGRAY),
			errBrush: CreateSolidBrush(CLR.ErrorColor))
		}
	Recalc()
		{
		super.SetFont(.font, .size, .weight, .text is "" ? "X" : .text)
		box = .imageSize()
		gap = Max(3, .TextExtent(' ').x)
		if .text is ""
			{
			gap = 0
			.Xmin = 0
			}
		.wtext = .Xmin
		if .lefttext
			{
			.xtext = 1 // allow for focus rect
			.xbox = .Xmin + gap
			}
		else
			{
			.xbox = 0
			.xtext = box + gap
			}
		.Xmin += box + gap

		focusRectBorder = 2
		.Xmin += focusRectBorder
		.Ymin += focusRectBorder
		formControlOverlap = 3
		.Ymin += formControlOverlap
		++.Top // so baseline is correct
		} // Exit()
	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		GetClientRect(.Hwnd, r = Object())
		r.Set_readonly()
		WithBkMode(hdc, TRANSPARENT)
			{
			.PaintImage(hdc, r)
			.paintText(hdc, r)
			if .HasFocus?() and .text isnt ""
				.paintFocusRect(hdc, r)
			}
		EndPaint(.Hwnd, ps)
		return 0
		}
	imageSize()
		{
		return Max(12, (.Top * .9).Round(0))
		}
	Antialias: true
	PaintImage(hdc, r)
		{
		h = .imageSize()
		top = Max(0, r.top + .Top - h + 1)
		rect = Rect(r.left + .xbox, top, h, h)
		if .Antialias
			Image.PaintWithAntialias(hdc, h, h, rect)
				{|hdcBmp, w, h|
				.paint(hdcBmp, Rect(0, 0, 0, Max(w, h)))
				}
		else
			.paint(hdc, rect)
		}
	paint(hdc, rect)
		{
		h = rect.GetHeight()
		brush = GetSysColorBrush(COLOR.BTNFACE)
		FillRect(hdc, [right: h, bottom: h], brush)

		DoWithHdcObjects(hdc, .outerObs())
			{
			.PaintOuter(hdc, rect)
			}
		if .Get() is true
			{
			imageBrush = GetStockObject(.enabled ? SO.BLACK_BRUSH : SO.GRAY_BRUSH)
			DoWithHdcObjects(hdc, [imageBrush])
				{
				bgBrush = .GetReadOnly() or not .enabled ? brush : false
				.PaintInner(hdc, rect, :imageBrush, :bgBrush)
				}
			}
		}
	outerObs()
		{
		return .GetReadOnly() or not .enabled
			? [GetSysColorBrush(COLOR.BTNFACE), .obs.grayPen]
			: .valid
				? [.MouseOver?() ? .obs.bluePen : .obs.darkGrayPen]
				: [.obs.errBrush, .obs.grayPen]
		}

	PaintOuter(hdc /*unused*/, rect /*unused*/)
		{
		throw 'PaintOuter method must be defined by derived class'
		}

	PaintInner(hdc /*unused*/, rect /*unused*/, imageBrush /*unused*/, bgBrush /*unused*/)
		{
		throw 'PaintInner method must be defined by derived class'
		}

	paintText(hdc, r)
		{
		WithHdcSettings(hdc, [.GetFont(), SetTextColor: .enabled ? .color : CLR.GRAY])
			{
			TextOut(hdc, r.left + .xtext, r.top, .text, .text.Size())
			}
		}

	paintFocusRect(hdc, r)
		{
		rf = [left: r.left + .xtext - 1, top: r.top,
			right: r.left + .xtext + .wtext + 1, bottom: r.bottom]
		DrawFocusRect(hdc, rf)
		}

	readonly: false
	set_readonly: false
	SetReadOnly(ro)
		{
		.set_readonly = ro
		.Repaint()
		}
	enabled: true
	SetEnabled(.enabled)
		{
		super.SetEnabled(enabled)
		.Repaint()
		}
	GetReadOnly()
		{
		return .readonly or .set_readonly
		}

	SETFOCUS()
		{
		.Repaint()
		return 0
		}
	KILLFOCUS()
		{
		.Repaint()
		return 0
		}
	LBUTTONDOWN()
		{
		.SetFocus()
		.gotMouseDown = true
		return 0
		}
	// if another window (e.g. list) creates this control on double-click
	// then we will receive an LBUTTONUP that we need to ignore
	gotMouseDown: false
	LBUTTONUP()
		{
		if .gotMouseDown
			.Toggle()
		return 0
		}
	value: ""
	Toggle()
		{
		if not .GetReadOnly()
			{
			.value = .value isnt true
			.ToolTip(.tip)
			.dirty? = true
			.Repaint()
			}
		}
	Get()
		{
		return .value is true
		}
	Set(value)
		{
		if value is .value
			return

		if not BooleanOrEmpty?(value)
			.ToolTip('invalid value: ' $ Display(value) $ '. Value must be true or false')
		else
			.ToolTip(.tip)
		.value = value
		.dirty? = false
		.Repaint()
		}

	Valid?()
		{
		.valid = BooleanOrEmpty?(.value)
		.Repaint()
		return .valid
		}

	GetUnvalidated()
		{
		return .value
		}

	ValidData?(value)
		{
		return BooleanOrEmpty?(value)
		}

	dirty?: false
	Dirty?(dirty = "")
		{
		Assert(BooleanOrEmpty?(dirty))
		if (dirty isnt "")
			.dirty? = dirty
		return .dirty?
		}

	GETDLGCODE()
		{
		return DLGC.WANTCHARS
		}
	CHAR(wParam)
		{
		if wParam is VK.SPACE
			.Toggle()
		return 0
		}

	mouseover?: false
	MOUSEMOVE()
		{
		if .mouseover? or .GetReadOnly()
			return 0
		TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
			dwFlags: TME.LEAVE hwndTrack: .Hwnd))
		.mouseover? = true
		.Repaint()
		return 0
		}
	MOUSELEAVE()
		{
		.mouseover? = false
		.Repaint()
		return 0
		}
	MouseOver?()
		{
		return .mouseover?
		}

	// for HighlightControl
	FindControl(name)
		{
		return name is 'Static' ? this : false
		}
	SetFont(.font = "", .size = "", .weight = "")
		{
		super.SetFont(font, size, weight)
		.WindowRefresh() // will call Recalc which will handle size change
		}
	color: 0
	SetColor(color)
		{
		.color = TranslateColor(color)
		}
	GetColor()
		{
		return .color
		}

	GetText()
		{
		return .text
		}

	ContextMenu(x, y) {	.DoDevContextMenu(x, y)	}

	Resize(x, y, w, h)
		{
		// adjust to handle overlap in FormControl
		super.Resize(x, y + 1, w, h - 2)
		}

	Destroy()
		{
		for ob in .obs
			DeleteObject(ob)
		.obs = #()
		super.Destroy()
		}
	}
