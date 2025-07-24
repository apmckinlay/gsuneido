// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: ChooseButton
	New(.text, .list, width = false)
		{
		super([MenuButtonControl, text, list, left:, :width, name: 'ChooseButton'])
		.prevproc = SetWindowProc(.ChooseButton.Hwnd, GWL.WNDPROC, .Fieldproc)
		.thisproc = GetWindowLongPtr(.ChooseButton.Hwnd, GWL.WNDPROC)
		.ChooseButton.Ymin = (.Ymin -= 2)
		.Send(#Data)
		}
	// if you use this class instead of MenuButtonControl in New
	// then it will only pull down the list if you click on the arrow
	menubutton: MenuButtonControl
		{
		LBUTTONDOWN(lParam)
			{
			GetClientRect(.Hwnd, r = [])
			x = LOSWORD(lParam)
			if x > r.right - 15 /* = width of the arrow */
				super.LBUTTONDOWN(:lParam)
			else
				SetFocus(.Hwnd)
			return 0
			}
		}
	Fieldproc(hwnd, msg, wparam, lparam)
		{
		_hwnd = .WindowHwnd()
		if msg is WM.GETDLGCODE and
			false isnt (m = MSG(lparam)) and
			m.message is WM.CHAR
			return DLGC.WANTALLKEYS
		if msg is WM.CHAR
			return .char(wparam.Chr())
		return CallWindowProc(.prevproc, hwnd, msg, wparam, lparam)
		}
	char(c)
		{
		if .GetReadOnly() is true
			return 0
		list = .list.Assocs().Filter({ it[1].Prefix?(c) })
		if list.Empty?()
			return 0
		i = list.FindIf({ it[1] is .ChooseButton.Get() })
		i = i is false ? 0 : (i + 1) % list.Size()
		return .On_ChooseButton(list[i][1], list[i][0])
		}

	Set(value)
		{
		.ChooseButton.Set(value is "" ? .text : value)
		}
	Get()
		{
		return .ChooseButton.Get()
		}
	On_ChooseButton(value, index)
		{
		.Set(value)
		.Send(#NewValue, value)
		.Send("On_" $ .Name, value, index)
		return 0
		}
	GetReadOnly()
		{
		return .ChooseButton.Disabled?()
		}
	SetReadOnly(readonly)
		{
		.ChooseButton.Disable(readonly)
		.ChooseButton.Grayed?(readonly or .Get() is '')
		}
	SetList(.list)
		{
		.ChooseButton.SetMenu(list)
		}
	prevproc: false
	Destroy()
		{
		.Send(#NoData)
		if .prevproc isnt false
			{
			if GetWindowLongPtr(.ChooseButton.Hwnd, GWL.WNDPROC) is .thisproc
				SetWindowProc(.ChooseButton.Hwnd, GWL.WNDPROC, .prevproc)
			ClearCallback(.Fieldproc)
			.prevproc = false
			}
		super.Destroy()
		}
	}
