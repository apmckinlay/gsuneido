// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
// contributions by Claudio Mascioni
// common parent for ChooseList and ChooseField
PassthruController
	{
	New(field, buttonBefore = false, .allowReadOnlyDropDown = false)
		{
		hidden = Object?(field) and field.GetDefault('hidden', false)
		.Button = .Construct(DropDownButtonControl, :hidden, :allowReadOnlyDropDown)
		.Field = .Construct(field)
		.prevproc = SetWindowProc(.Field.Hwnd, GWL.WNDPROC, .callback = .Fieldproc)
		.thisproc = GetWindowLongPtr(.Field.Hwnd, GWL.WNDPROC)
		.adjustForButton(buttonBefore)
		.Xmin = .Field.Xmin + .buttonWidth // allow for button
		.Ymin = .Field.Ymin
		.Top = .Field.Top
		.children = Object(.Field, .Button)
		.buttonBefore = buttonBefore is true
		.Send("Data")
		}
	adjustForButton(buttonBefore)
		{
		ex = .Field.TextExtent('M')
		bh = ex.y + ScaleWithDpiFactor(4) /* = 4 for edit ctrl margins */
		bw = bh - 2

		margins = .Field.GetMargins()
		if buttonBefore
			margins.left += bw
		else
			margins.right += bw
		.Field.SetMargins(margins)

		.buttonWidth = bw
		}
	fieldProcOverrideCheck: false
	SetFieldProcOverrideCheck(fn)
		{
		.fieldProcOverrideCheck = fn
		}
	allowOverride?()
		{
		if .fieldProcOverrideCheck is false
			return true
		return (.fieldProcOverrideCheck)()
		}
	Fieldproc(hwnd, msg, wparam, lparam)
		{
		if not .allowOverride?()
			return .defaultHandle(hwnd, msg, wparam, lparam)

		_hwnd = .WindowHwnd()
		if msg is WM.GETDLGCODE
			return .dlgCode(lparam)
		if .dropDown?(msg, wparam, lparam)
			{
			.On_DropDown()
			return 0
			}
		if msg is WM.SYSCHAR and .Alt_z?(wparam, lparam)
			return 0

		.handleFocus(msg)
		if msg is WM.CHAR and wparam is VK.TAB
			return 0
		if msg is WM.CHAR and wparam is VK.RETURN
			{
			.FieldReturn()
			return 0
			}
		if msg is WM.KEYUP and wparam is VK.ESCAPE
			{
			.Send('FieldEscape')
			return 0
			}

		return .defaultHandle(hwnd, msg, wparam, lparam)
		}
	dlgCode(lparam)
		{
		// topwindow is used when calling this method directly on
		// the class (there is no instance, therefore no .Window)
		topwindow = .Window
		if (false isnt (m = MSG(lparam)) and
			(m.message is WM.CHAR or m.message is WM.KEYDOWN) and
			(m.wParam is VK.UP or m.wParam is VK.DOWN or
				(m.wParam is VK.RETURN and not topwindow.Base?(Dialog))))
			return DLGC.WANTALLKEYS
		if (false isnt (m = MSG(lparam)) and
			m.message is WM.SYSCHAR and
			m.wParam.Chr() is 'z')
			return DLGC.WANTALLKEYS
		// need this so tab works
		return DLGC.WANTCHARS | DLGC.WANTARROWS | DLGC.HASSETSEL
		}
	dropDown?(msg, wparam, lparam)
		{
		if (((msg is WM.KEYDOWN or msg is WM.SYSKEYDOWN) and
			(wparam is VK.UP or wparam is VK.DOWN)) or
			(msg is WM.SYSKEYDOWN and .Alt_z?(wparam, lparam)))
			return true
		return false
		}
	handleFocus(msg)
		{
		if msg is WM.SETFOCUS
			.FieldSetFocus()
		if msg is WM.KILLFOCUS
			.FieldKillFocus()
		}
	defaultHandle(hwnd, msg, wparam, lparam)
		{
		if msg is WM.MOUSEWHEEL and not .HasFocus?()
			return DefWindowProc(hwnd, msg, wparam, lparam)

		return CallWindowProc(.prevproc, hwnd, msg, wparam, lparam)
		}

	On_DropDown()
		{
		throw "must be defined by derived class"
		}

	Alt_z?(wparam, lparam)
		{ return wparam.Chr().Lower() is 'z' and (((lparam >> 29) & 1) is 1) }
	FieldSetFocus()
		{
		.Send('Field_SetFocus')
		}
	FieldKillFocus()
		{
		}
	FieldReturn()
		{
		}

	children: false
	GetChildren()
		{
		return .children isnt false ? .children : Object()
		}

	Data()
		{
		// Block the Data message from the Field control
		}

	NoData()
		{
		}

	EditHwnd()
		{
		return .Field.Hwnd
		}

	Edit_ParentValid?()
		{
		return .Valid?()
		}

	SetFocus()
		{
		.Field.SetFocus()
		}

	InitDropDown()
		{
		if .Destroyed?() or .dropDownReadOnly()
			return false
		SetFocus(.Field.Hwnd)
		// Focus change could cause this field to be destroyed or protected
		if .Destroyed?() or .dropDownReadOnly()
			return false
		return GetWindowRect(.Field.Hwnd)
		}
	dropDownReadOnly()
		{
		if .allowReadOnlyDropDown is true
			return false
		return .GetReadOnly()
		}

	Set(value)
		{
		.Field.Set(value)
		}
	Get()
		{
		return .Field.Get()
		}

	Dirty?(dirty = "")
		{
		return .Field.Dirty?(dirty)
		}

	NewValue(value /*unused*/)
		{
		// resend newvalue so that this controller becomes the source
		// and thus responsible for Get method.
		.Field.Dirty?(true)
		.Send("NewValue", .Get())
		}

	SetValid(valid)
		{
		.Field.SetValid(valid)
		}

	RemoveButton() // used by KeyControl
		{
		.children.Remove(.Button)
		.Button.Destroy()
		.Button = false
		}

	THEMECHANGED()
		{
		.Resize(.x, .y, .w, .h)
		return 0
		}
	Resize(.x, .y, .w, .h)
		{
		.Field.Resize(x, y, w, h)
		if (.Button is false)
			return
		GetClientRect(.Field.Hwnd, r = Object())
		bh = r.bottom - r.top
		bw = .buttonWidth
		p = Object(x: .buttonBefore ? r.left : r.right - bw, y: 0)
		ClientToScreen(.Field.Hwnd, p)
		ScreenToClient(GetParent(.Field.Hwnd), p)
		if .Controller.Base?(ListEditWindow)
			{
			p.x += 2
			p.y -= 2
			}
		else
			{
			--p.x
			bw += 2
			}
		// leave extra gap above and bottom for windows 11 underline
		.Button.Resize(p.x, p.y, bw, bh)
		SetWindowPos(.Button.Hwnd, HWND.TOPMOST, 0, 0, 0, 0, SWP.NOSIZE | SWP.NOMOVE)
		}

	SetFont(font, size)
		{
		.Field.SetFont(font, size)
		}
	SetStatus(status)
		{
		.Field.SetStatus(status)
		}
	SetBgndColor(color)
		{
		.Field.SetBgndColor(color)
		}
	SetTextColor(color)
		{
		.Field.SetTextColor(color)
		}

	Update()
		{
		.Field.Update()
		.Button.Update()
		super.Update()
		}

	prevproc: false
	Destroy()
		{
		.Send("NoData")
		if .prevproc isnt false
			{
			if GetWindowLongPtr(.Field.Hwnd, GWL.WNDPROC) is .thisproc
				SetWindowProc(.Field.Hwnd, GWL.WNDPROC, .prevproc)
			ClearCallback(.callback)
			.prevproc = false
			}
		super.Destroy()
		}
	}
