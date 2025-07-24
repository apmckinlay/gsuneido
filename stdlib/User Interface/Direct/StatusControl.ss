// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: 'Status'
	Xstretch: 1
	New(text = "")
		{
		.CreateWindow("SuBtnfaceArrow", '', WS.VISIBLE, WS_EX.STATICEDGE)
		.SubClass()
		.color = .normal = GetSysColor(COLOR.BTNFACE)
		.SetFont(text: "M")
		.vOffset = (.Ymin / 4 /*= offset factor*/).Round(0)
		.Ymin += .vOffset * 2
		.untranslated = text
		.text = TranslateLanguage(text)
		.err_brush = CreateSolidBrush(CLR.ErrorColor)
		.warn_brush = CreateSolidBrush(CLR.WarnColor)
		}
	ERASEBKGND()
		{
		return 1
		}
	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		GetClientRect(.Hwnd, r = Object())
		FillRect(hdc, r, .getBrush())
		text = .text isnt "" ? .text : .defaultMsg
		if text isnt ""
			WithHdcSettings(hdc, [.Hwnd_hfont, SetBkMode: TRANSPARENT],
				{ TextOut(hdc, 2, .vOffset, text, text.Size()) })
		EndPaint(.Hwnd, ps)
		return 0
		}
	getBrush()
		{
		return .color is CLR.ErrorColor
			? .err_brush
			: .color is CLR.WarnColor
				? .warn_brush
				: GetSysColorBrush(COLOR.BTNFACE)
		}
	statusLimit: 500
	Set(text, normal = false, warn = false, invalid = false)
		{
		.setColor(normal, warn, invalid)
		text = String(text)
		.untranslated = text
		.text = TranslateLanguage(text).Ellipsis(.statusLimit, atEnd:)
		.Repaint()
		}
	setColor(normal, warn, invalid)
		{
		if invalid
			.color = CLR.ErrorColor
		else if warn
			.color =  CLR.WarnColor
		else if normal
			.color = .normal
		}

	defaultMsg: ""
	SetDefaultMsg(text)
		{
		.defaultMsg = TranslateLanguage(text).Ellipsis(.statusLimit, atEnd:)
		.Repaint()
		}

	SetValid(valid = true)
		{
		.color = valid is true ? .normal : CLR.ErrorColor
		.Repaint()
		}
	SetWarning(warn = true)
		{
		.color = warn is true ? CLR.WarnColor : .normal
		.Repaint()
		}
	GetValid()
		{
		return .color is .normal
		}
	Get()
		{
		return .untranslated
		}
	GetReadOnly() // read-only not applicable to status
		{
		return true
		}
	Destroy()
		{
		DeleteObject(.err_brush)
		DeleteObject(.warn_brush)
		super.Destroy()
		}
	}
