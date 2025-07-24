// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide
// e.g. enlargeOnHover: #(imagePadding: .15, x: 60, y: 15)
ButtonControl
	{
	New(.text = false, command = false,
		tabover = false, .defaultButton = false, style = 0, tip = false, pad = false,
		font = "", size = "", weight = "", textColor = false,
		.width = false, .buttonWidth = false, .buttonHeight = false,
		italic = false, underline = false, strikeout = false,
		image = false, mouseOverImage = false, mouseDownImage = false,
		.imageColor = false, .mouseOverImageColor = false,
		.book = 'imagebook', .mouseEffect = false, .imagePadding = 0,
		.buttonStyle = false, classic = false, .noBgnd = false, .enlargeOnHover = false,
		.alignTop = false, hidden = false)
		{
		super(.text is false ? 'EnhancedButton' : .text, command,
			tabover: tabover or not .buttonStyle, :defaultButton,
			style: style | BS.OWNERDRAW, :tip, :pad,
			:font, :size, :weight, color: .textColor = TranslateColor(textColor),
			:italic, :underline, :strikeout, :hidden)
		.buttonName = command
		.SetImage(image, mouseOverImage, mouseDownImage)
		.skip = .ButtonControl_pad / 2
		.calcSize()
		.setBrush()
		.setDraw(classic)
		}
	setDraw(classic)
		{
		.draw = (classic is false) ? .drawThemed : .drawUnthemed
		}
	image: false
	imageObj: false
	imageW: false
	imageH: false
	mouseOverImage: false
	mouseOverImageObj: false
	mouseDownImage: false
	mouseDownImageObj: false
	GETDLGCODE()
		{
		if .buttonStyle is false
			return 'callsuper'
		return DLGC.BUTTON | (.defaultButton ? DLGC.DEFPUSHBUTTON : DLGC.UNDEFPUSHBUTTON)
		}
	BM_SETSTYLE(wParam)
		{
		if .buttonStyle is false
			return 'callsuper'

		if wParam is BS.DEFPUSHBUTTON
			.setDefaultButton(true)
		else if wParam is BS.PUSHBUTTON
			.setDefaultButton(false)
		return 0
		}
	setDefaultButton(defaultButton)
		{
		.defaultButton = defaultButton
		InvalidateRect(.Hwnd, NULL, true)
		}
	GetImage()
		{
		return .image
		}
	SetImage(image, mouseOverImage = false, mouseDownImage = false)
		{
		.imageObj = .getImage(.image, image, .imageObj)
		.image = image
		if .imageObj isnt false
			.WithDC(.Hwnd)
				{ |hdc|
				.imageH = .imageObj.Height(hdc)
				.imageW = .imageObj.Width(hdc)
				}
		.mouseOverImageObj =
			.getImage(.mouseOverImage, mouseOverImage, .mouseOverImageObj)
		if mouseOverImage isnt false
			.mouseOverImage = mouseOverImage
		.mouseDownImageObj =
			.getImage(.mouseDownImage, mouseDownImage, .mouseDownImageObj)
		if mouseDownImage isnt false
			.mouseDownImage = mouseDownImage
		InvalidateRect(.Hwnd, NULL, false)
		}
	getImage(oldImage, newImage, oldImageObj)
		{
		if newImage is false or newImage is oldImage
			return oldImageObj
		if oldImage isnt false
			oldImageObj.Close()
		if Object?(newImage)
			return TwoImages(@newImage)
		return ImageResource(newImage, book: .book)
		}
	alignTop: false
	calcSize()
		{
		if .buttonStyle is false and not .alignTop // to align with other field in the row
			.Top = 0
		// The ButtonControl calculates the width of the text and its pad
		// during initialization
		.tw = .text is false ? 0 : .Xmin
		w = .buttonWidth isnt false
			? .buttonWidth
			: .width isnt false ? .CalcWidth(.width) : 0
		h = .buttonHeight isnt false ? .buttonHeight : .Ymin
		.updateSize(w, h)
		.Ymin = h
		.Xmin = Max(.ix + .iw + .tw, w)
		}
	updateSize(w, h)
		{
		.ih = h
		.iy = 0
		.iw = .imageObj is false ? 0 : (.ih * .imageW / .imageH).Round(0)
		ixFactor = .7
		txFactor = .4
		if w isnt 0 and .iw > w
			.iw = w
		noImageOrStyle = .imageObj is false or not .buttonStyle
		.ix = Max(noImageOrStyle or .text is false ? 0 : .skip * ixFactor,
			(w - .iw - .tw) / 2)
		.tx = .ix + .iw - (noImageOrStyle ? 0 : .skip * txFactor)
		}
	Resize(x, y, w, h)
		{
		.updateSize(w, h)
		return super.Resize(x, y, w, h)
		}
	imageBrush: false
	mouseOverImageBrush: false
	disabledImageBrush: false
	setBrush()
		{
		.disabledImageBrush = CreateSolidBrush(CLR.GRAY)
		.imageBrush = .imageColor isnt false
			? CreateSolidBrush(.imageColor)
			: .buttonStyle isnt false
				? CreateSolidBrush(CLR.BLACK)
				: CreateSolidBrush(CLR.EnhancedButtonFace)
		.mouseOverImageBrush = .mouseOverImageColor isnt false
			? CreateSolidBrush(.mouseOverImageColor)
			: CreateSolidBrush(CLR.BLACK)
		}
	deleteBrush()
		{
		if .imageBrush isnt false
			DeleteObject(.imageBrush)
		if .mouseOverImageBrush isnt false
			DeleteObject(.mouseOverImageBrush)
		if .disabledImageBrush isnt false
			DeleteObject(.disabledImageBrush)
		}
	SetImageColor(imageColor = false, mouseOverImageColor = false)
		{
		if .imageBrush isnt false
			DeleteObject(.imageBrush)
		.imageColor = imageColor
		.imageBrush = CreateSolidBrush(imageColor)

		if mouseOverImageColor isnt false
			{
			if .mouseOverImageBrush isnt false
				DeleteObject(.mouseOverImageBrush)
			.mouseOverImageColor = mouseOverImageColor
			.mouseOverImageBrush = CreateSolidBrush(mouseOverImageColor)
			}
		.Repaint()
		}
	GetImageColor()
		{
		return .imageColor
		}
	SetEnabled(enabled)
		{
		.enabled = enabled
		super.SetEnabled(enabled)
		}
	prevFocusHwnd: 0
	SETFOCUS(wParam)
		{
		.prevFocusHwnd = wParam
		return 'callsuper'
		}
	CLICKED()
		{
		if .prevFocusHwnd isnt 0 and .buttonStyle is false
			{
			SetFocus(.prevFocusHwnd)
			}
		return super.CLICKED()
		}
	focusRectOffSet: -3
	DRAWITEM(dis)
		{
		hdc = dis.hDC
		rect = dis.rcItem
		textRect = rect.Copy()
		textRect.left += .tx
		textRect.right = textRect.left + .tw
		disabled = (dis.itemState & ODS.DISABLED) isnt 0
		focused = (dis.itemState & ODS.FOCUS) isnt 0

		.DrawBackground(hdc, rect, disabled, focused)
		.drawText(hdc, textRect)
		.drawImage(hdc, rect)

		if .buttonStyle and (dis.itemState & ODS.FOCUS) isnt 0
			{
			InflateRect(rect, .focusRectOffSet, .focusRectOffSet)
			DrawFocusRect(hdc, rect)
			}
		.DrawExtra(hdc, :rect)
		return true
		}
	DrawBackground(hdc, rect, disabled, focused)
		{
		if .mouseEffect and (.buttonStyle or .mousedown or .mouseover or .pushed)
			(.draw)(hdc, rect, disabled, focused)
		else
			FillRect(hdc, rect, GetSysColorBrush(COLOR.BTNFACE))
		}
	drawImage(hdc, rect)
		{
		if ((.pushed or .mousedown) and .mouseDownImageObj isnt false)
			curImage = .mouseDownImageObj
		else if (.mouseover and .mouseOverImageObj isnt false)
			curImage = .mouseOverImageObj
		else
			curImage = .imageObj
		if curImage isnt false
			if curImage.IsRasterImage()
				curImage.Draw(hdc, .ix, .iy, .iw, .ih)
			else
				.drawVectorImage(hdc, curImage, rect)
		}
	DrawExtra(hdc /*unused*/, rect /*unused*/) { }
	imageOffSet: 0.5
	drawVectorImage(hdc, curImage, rect)
		{
		bgBrush = false
		if not curImage.Base?(ImageFont)
			{
			if false is samplePoint = .FindVisibleSamplePoint(hdc)
				return
			bgBrush = CreateSolidBrush(GetPixel(hdc, samplePoint.x, samplePoint.y))
			}
		pressed = .mousedown or .pushed
		imageBrush = .getImageBrush(pressed)

		offset = .mouseEffect and pressed and .text is false
			? -ScaleWithDpiFactor(.imageOffSet)
			: 0
		paddingW = (.iw * .imagePadding).Round(0)
		paddingH = (.ih * .imagePadding).Round(0)

		hrgn = CreateRectRgn(rect.left + 2, rect.top + 2, rect.right - 2, rect.bottom - 2)
		SelectClipRgn(hdc, hrgn)
		w = .iw - paddingW * 2
		h = .ih - paddingH * 2
		x = .ix + paddingW - offset
		y = .iy + paddingH - offset
		curImage.Draw(hdc, x, y, w, h, imageBrush, bgBrush)
		SelectClipRgn(hdc, NULL)
		DeleteObject(hrgn)

		if bgBrush isnt false
			DeleteObject(bgBrush)
		}
	getImageBrush(pressed)
		{
		return not .enabled
			? .disabledImageBrush
			: not pressed and not .mouseover
				? .imageBrush
				: .mouseOverImageBrush
		}
	FindVisibleSamplePoint(hdc)
		{
		sampleOffset = -3 // in order to avoid sampling the border color
		res = GetClipBox(hdc, rc = Object())
		InflateRect(rc, sampleOffset, sampleOffset)
		if res is REGIONTYPE.COMPLEXREGION or res is REGIONTYPE.SIMPLEREGION
			{
			if Abs(rc.bottom - rc.top) < Abs(rc.right - rc.left)
				{
				direction = 'y'
				start = Min(rc.top, rc.bottom)
				end = Max(rc.top, rc.bottom)
				}
			else
				{
				direction = 'x'
				start = Min(rc.left, rc.right)
				end = Max(rc.left, rc.right)
				}
			return .scan(direction, rc, start, hdc, end)
			}
		return false
		}
	scan(direction, rc, start, hdc, end)
		{
		do
			{
			x = direction is 'y' ? rc.left : start
			y = direction is 'y' ? start : rc.top
			if PtVisible(hdc, x, y) is true
				return Object(:x, :y)
			x = direction is 'y' ? rc.left : end
			y = direction is 'y' ? end : rc.top
			if PtVisible(hdc, x, y) is true
				return Object(:x, :y)
			start++
			end--
			}
		while(start <= end)
		return false
		}
	drawUnthemed(dc, rect, disabled, unused)
		{
		state = .mousedown or .pushed ? DFCS.PUSHED :
				.mouseover ? DFCS.HOT :
				disabled or .grayed ? DFCS.INACTIVE : 0
		DrawFrameControl(dc, rect, DFC.BUTTON, DFCS.BUTTONPUSH | state)
		}
	drawThemed(dc, rect, disabled, focused)
		{
		state = disabled or .grayed
			? THEME.PBS.DISABLED
			: .mousedown or .pushed
				? THEME.PBS.PRESSED
				: .mouseover
					? THEME.PBS.HOT
					: .defaultButton or focused
						? THEME.PBS.DEFAULTED
						: THEME.PBS.NORMAL
		hTheme = GetWindowTheme(.Hwnd)
		DrawThemeBackground(hTheme, dc, THEME.BP.PUSHBUTTON, state, rect, NULL)
		}
	drawText(dc, rect)
		{
		if .text isnt false
			WithHdcSettings(dc, [SetBkMode: TRANSPARENT, SetTextColor: .textColor])
				{
				DrawText(dc, GetWindowText(.Hwnd), -1, rect,
					DT.CENTER | DT.VCENTER | DT.SINGLELINE)
				}
		}

	//mousedown is for LEFTBUTTONDOWN event
	//pushed is to record long last pushed status (e.g. in DrawPaletteControl)
	pushed: false
	mousedown: false
	mouseover: false
	enabled: true
	grayed: false
	Pushed?(state = -1)
		{
		if state isnt -1 and state isnt .pushed
			{
			.pushed = state
			InvalidateRect(.Hwnd, NULL, true)
			}
		return .pushed
		}
	gray: 0xaaaaaa
	Grayed(state = -1)
		{
		if state isnt -1
			.SetTextColor(state is true ? .gray : false)
		return .textColor is .gray
		}
	SetText(text) // needed?
		{
		.Set(text)
		}
	Set(.text)
		{
		super.Set(.text isnt false ? .text : '')
		.calcSize()
		.WindowRefresh()
		}
	SetTextColor(color)
		{
		.textColor = color is false ? GetSysColor(COLOR.BTNTEXT) : TranslateColor(color)
		InvalidateRect(.Hwnd, NULL, false)
		}
	SetSize(.buttonWidth, .buttonHeight)
		{
		.calcSize()
		}
	SetMouseEffect(mouseEffect)
		{
		.mouseEffect = mouseEffect
		}
	GetTextColor()
		{
		return .textColor
		}
	GetButtonName()
		{
		return .buttonName
		}
	SetPressed?(pressed?)
		{
		.mousedown = pressed?
		}
	GetPressed?()
		{
		return .mousedown
		}
	MouseOver?()
		{
		return .mouseover
		}
	ERASEBKGND(wParam)
		{
		if not .noBgnd
			{
			hdc = wParam
			GetClientRect(.Hwnd, rect = Object())
			FillRect(hdc, rect, GetSysColorBrush(COLOR.BTNFACE))
			}
		return 1
		}
	RBUTTONDOWN()
		{
		rc = GetWindowRect(.Hwnd)
		.Send('EnhancedRButtonDown', rc)
		return 0
		}
	LBUTTONDOWN()
		{
		if false isnt .Send('EnhancedButtonAllowPush')
			.mousedown = true
		return 'callsuper'
		}
	LBUTTONUP()
		{
		.mousedown = false
		return 'callsuper'
		}
	prevSize: false
	MOUSEMOVE()
		{
		TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
			dwFlags: TME.LEAVE hwndTrack: .Hwnd))
		.enlarge()
		if .mouseover or .mousedown or .pushed
			return 'callsuper'
		.mouseover = true
		InvalidateRect(.Hwnd, NULL, false)
		return 'callsuper'
		}
	enlarge()
		{
		if .enlargeOnHover is false or .prevSize isnt false
			return
		rect = .GetClientRect()
		posRect = .GetRect().ToWindowsRect()
		.prevSize = Object(width: rect.GetWidth(), height: rect.GetHeight(),
			imagePadding: .imagePadding, left: posRect.left, top: posRect.top)
		.mouseEffect = true
		.buttonStyle = true
		.imagePadding = .enlargeOnHover.imagePadding
		.WithSelectObject(.GetFont())
			{|hdc|
			GetTextExtentPoint32(hdc, 'M', 1, ex = Object())
			}
		textSize = .text is false ? 0 : .text.Size()
		expandX = ex.x * textSize + ScaleWithDpiFactor(15) 		/*= x padding */
		expandY = ex.y + ScaleWithDpiFactor(10)					/*= y padding */
		.updateSize(expandX, expandY)
		if .enlargeOnHover.GetDefault('direction', false) is 'top-left'
			.Resize(posRect.left - expandX + .buttonWidth, posRect.bottom - expandY,
				expandX, expandY)
		else
			.Resize(posRect.left, posRect.top, .prevSize.width + expandX, expandY)
		SetWindowPos(.Hwnd, HWND.TOP, 0, 0, 0, 0,
			SWP.NOCOPYBITS | SWP.NOACTIVATE | SWP.NOMOVE | SWP.NOREDRAW | SWP.NOSIZE)
		.Defer(.Repaint)
		}
	MOUSELEAVE()
		{
		.mouseover = false
		if .prevSize isnt false
			{
			.mouseEffect = false
			.buttonStyle = false
			.imagePadding = .prevSize.imagePadding
			.updateSize(.prevSize.width, .prevSize.height)
			.Resize(.prevSize.left, .prevSize.top, .prevSize.width, .prevSize.height)
			.prevSize = false
			}
		InvalidateRect(.Hwnd, NULL, true)
		return 'callsuper'
		}
	closeImages()
		{
		if .imageObj isnt false
			{
			.imageObj.Close()
			.image = .imageObj = false
			}
		if .mouseOverImageObj isnt false
			{
			.mouseOverImageObj.Close()
			.mouseOverImage = .mouseOverImageObj = false
			}
		if .mouseDownImageObj isnt false
			{
			.mouseDownImageObj.Close()
			.mouseDownImage = .mouseDownImageObj = false
			}
		}
	ContextMenu(x, y)
		{
		.Send('EnhancedButton_ContextMenu', :x, :y)
		return super.ContextMenu(x, y)
		}
	Destroy()
		{
		.closeImages()
		.deleteBrush()
		super.Destroy()
		}
	}
