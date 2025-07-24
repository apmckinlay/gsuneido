// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title1: 'Local'
	Title2: 'Master'
	New(image1, image2, style1 = false, style2 = false)
		{
		super(.controls(image1, image2, style1, style2))
		.localscroll = .FindControl(#localscroll)
		.masterscroll = .FindControl(#masterscroll)
		.localprevproc = SetWindowProc(.localscroll.Hwnd, GWL.WNDPROC, .Localproc)
		.masterprevproc = SetWindowProc(.masterscroll.Hwnd, GWL.WNDPROC, .Masterproc)
		}
	controls(image1, image2, style1, style2)
		{
		image1Ob = Object()
		image2Ob = Object()
		if image1.Suffix?('.emf')
			{
			image1Ob = Object('Image' image1, name: 'image1')
			image2Ob = Object('Image' image2, name: 'image2')
			}
		else
			{
			image1Ob = Object('Mshtml' image1, style: style1, name: 'image1')
			image2Ob = Object('Mshtml' image2, style: style2, name: 'image2')
			}
		return Object('Horz'
			#Skip
			Object('Vert'
				Object('Static', .Title1)
				#(Skip 5)
				Object('Scroll'
					image1Ob, name: 'localscroll', noEdge:)
				name: 'Vert1')
			#Skip
			#EtchedVertLine
			#Skip
			Object('Vert'
				Object('Static', .Title2)
				#(Skip 5)
				Object('Scroll'
					image2Ob, name: 'masterscroll', noEdge:)
				name: 'Vert2')
			#Skip
			)
		}
	Localproc(hwnd, msg, wparam, lparam)
		{
		.scrollproc(hwnd, msg, wparam, lparam, .localscroll, .masterscroll,
			.localprevproc)
		}
	Masterproc(hwnd, msg, wparam, lparam)
		{
		.scrollproc(hwnd, msg, wparam, lparam, .masterscroll, .localscroll,
			.masterprevproc)
		}
	scrollproc(hwnd, msg, wparam, lparam, scroll1, scroll2, prevproc)
		{
		_hwnd = .WindowHwnd()
		result = CallWindowProc(prevproc, hwnd, msg, wparam, lparam)
		if (msg is WM.VSCROLL)
			{
			scrollpos1 = .getScrollPos(scroll1.Hwnd, wparam, SB.VERT)
			scrollpos2 = .getScrollPos(scroll2.Hwnd, wparam, SB.VERT)
			scroll2.Scroll(0, scrollpos2 - scrollpos1)
			}
		else if (msg is WM.HSCROLL)
			{
			scrollpos1 = .getScrollPos(scroll1.Hwnd, wparam, SB.HORZ)
			scrollpos2 = .getScrollPos(scroll2.Hwnd, wparam, SB.HORZ)
			scroll2.Scroll(scrollpos2 - scrollpos1, 0)
			}
		return result
		}
	getScrollPos(hwnd, wparam, which)
		{
		if (LOWORD(wparam) is SB.THUMBTRACK)
			{
			GetScrollInfo(hwnd, which, info = Object(cbSize: SCROLLINFO.Size()
				fMask: SIF.TRACKPOS))
			return info.nTrackPos
			}
		else
			return GetScrollPos(hwnd, which)
		}
	Destroy()
		{
		SetWindowProc(.localscroll.Hwnd, GWL.WNDPROC, .localprevproc)
		ClearCallback(.Localproc)
		SetWindowProc(.masterscroll.Hwnd, GWL.WNDPROC, .masterprevproc)
		ClearCallback(.Masterproc)
		super.Destroy()
		}
	}
