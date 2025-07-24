// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: 'ColorRect'
	Xmin: 50
	Ymin: 50
	New(color = 0)
		{
		.CreateWindow("SuBtnfaceArrow", "", WS.VISIBLE, WS_EX.STATICEDGE)
		.SubClass()
		if Object?(color)
		   color = RGB(color[0], color[1], color[2])
		.color = color
		.Send('Data')
		}
	ERASEBKGND()
		{
		return 1
		}
	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		GetClientRect(.Hwnd, r = Object())
		brush = CreateSolidBrush(.color)
		FillRect(hdc, r, brush)
		DeleteObject(brush)
		EndPaint(.Hwnd, ps)
		return 0
		}
	LBUTTONDBLCLK()
		{
		.Send("ColorRect_DoubleClick", .color)
		return 0
		}
	Get()
		{
		return .color
		}
	Set(color)
		{
		.color = color
		.Repaint()
		}
	Dirty?(unused)
		{
		return false
		}
	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}