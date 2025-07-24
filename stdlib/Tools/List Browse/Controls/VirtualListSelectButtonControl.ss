// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
VirtualListThumbImageButtonControl
	{
	New()
		{
		super(@.layout())
		.Pushed?(true)
		}

	layout()
		{
		.sysWidth = GetSystemMetrics(SM.CXVSCROLL)
		return Object(text: "Select", command: 'VirtualListThumb_ArrowSelect',
			image: 'zoom.emf', mouseEffect:, imagePadding: 0.15,
			enlargeOnHover: #(imagePadding: .2, direction: 'top-left'),
			buttonWidth: .sysWidth, buttonHeight: .sysWidth)
		}

	GetWidth()
		{
		return .sysWidth
		}

	GetHeight()
		{
		return .sysWidth
		}

	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		SetWindowPos(.Hwnd, HWND.TOP, 0, 0, 0, 0, SWP.NOSIZE | SWP.NOMOVE)
		}
	}
