// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// base class for text controls (e.g. P, Head, Pre)
WndProc
	{
	Xstretch: 1
	Font: ""
	Size: "0"
	Weight: "normal"
	Flags: 0
	WndClass: "Html"
	Underline: false
	Italic: false
	Bold: false

	New(text)
		{
		.CreateWindow(.WndClass, "", WS.VISIBLE)
		.SubClass()
		.text = text
		.hfont = CreateFontIndirect(
			.LogFont(.Font, .Size, .Weight, .Underline, .Italic))
		}
	CalcSize()
		{
		r = .CalcRect()
		.Xmin = r.right
		.Ymin = r.bottom + 6 /*= offset */
		}
	CalcRect(w = 0)
		{
		r = Object(right: w)
		.WithSelectObject(.hfont)
			{|hdc|
			DrawText(hdc, .text, -1, r, DT.CALCRECT + DT.NOPREFIX + .Flags)
			}
		return r
		}
	ERASEBKGND()
		{ return 1 }
	PAINT()
		{
		dc = BeginPaint(.Hwnd, ps = Object())
		GetClientRect(.Hwnd, r = Object())
		FillRect(dc, r, GetSysColorBrush(COLOR.BTNFACE))

		.PreDraw(dc)
		WithHdcSettings(dc, [.hfont, SetBkMode: TRANSPARENT],
			{ DrawText(dc, .text, -1, r, .DrawFlags()) })
		.PostDraw(dc, r)
		EndPaint(.Hwnd, ps)
		return 0
		}
	DrawFlags()
		{
		return DT.NOPREFIX + .Flags
		}
	PreDraw(dc/*unused*/)
		{ }
	PostDraw(dc/*unused*/, r/*unused*/)
		{ }
	DESTROY()
		{
		DeleteObject(.hfont)
		return 0
		}
	}
