// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: 'Status'
	Xstretch: 1
	border: 2
	New(text = "")
		{
		.CreateWindow("SuBtnfaceArrow", text, WS.VISIBLE)
		.SubClass()
		.SetFont(text: "M")
		.Ymin *= 1.5 /* = ymin factor */
		.Ymin += .border
		.untranslated = text
		.text = TranslateLanguage(text)
		}
	PAINT()
		{
		DoWithHdcObjects(hdc = BeginPaint(.Hwnd, ps = Object()), [.GetFont()])
			{
			GetClientRect(.Hwnd, r = Object())
			r.top += ScaleWithDpiFactor(3) /* = top border */
			DrawStatusText(hdc, r, .text, 0)
			}
		EndPaint(.Hwnd, ps)
		return 0
		}
	Set(text)
		{
		.untranslated = text
		.text = TranslateLanguage(text)
		}
	Get()
		{
		return .untranslated
		}
	GetReadOnly() // read-only not applicable to statustext
		{ return true }
	}
