// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	Name: 'Number'
	New(.mask = '-###.###.###,##', .readonly = false,
		.rangefrom = false, .rangeto = false, width = false,
		set = false, .mandatory = false, status = "",
		font = "", size = "", weight = "")
		{
		super(width: .getWidth(width, mask),
			:readonly, justify: "RIGHT", style: ES.MULTILINE,
			:status, :font, :size, :weight)
		// have to have MULTILINE to get right justification on Windows 95
		.SubClass()
		.prev = Object(text: "", begin: Object(x: 0), end: Object(x: 0))
		if set isnt false
			{
			.Set(set)
			.Send("NewValue", .Get())
			}
		}

	defaultWidth: 15
	getWidth(width, mask)
		{
		return width is false ? mask is false ? .defaultWidth : mask.Size() : width
		}

	GETDLGCODE() // need this for tab to work with multiline
		{
		return DLGC.WANTCHARS | DLGC.WANTARROWS
		}
	EN_SETFOCUS()
		{
		SetWindowText(.Hwnd, GetWindowText(.Hwnd).Tr('.'))
		.SendMessage(EM.SETSEL, 0, 999 /* = to select all */)
		return super.EN_SETFOCUS()
		}
	KillFocus()
		{
		if .mask is false
			return
		s = super.Get().Tr('.').Tr(',', '.')
		if s isnt '' and .Valid?()
			{
			dirty? = .Dirty?()
			SetWindowText(.Hwnd, Number(s).EuroFormat(.mask))
			.Dirty?(dirty?)
			}
		}
	Valid?()
		{
		if .readonly
			return true
		valid = true
		s = super.Get().Tr('.').Tr(',', '.')
		if s is ''
			valid = not .mandatory
		else if not s.Number?()
			valid =  false
		else
			{
			n = Number(s)
			if ((.rangefrom isnt false and n < .rangefrom) or
				(.rangeto isnt false and n > .rangeto))
				valid = false
			}
		return valid
		}
	Get()
		{
		s = super.Get().Tr('.').Tr(',', '.')
		return s isnt '' and .Valid?() ? Number(s) : ""
		// return "" for invalid instead of s so that rules work
		}
	Set(value)
		{
		if (Number?(value) or (String?(value) and value.Number?()))
			{
			value = .mask is false ? String(value) : Number(value).EuroFormat(.mask)
			if GetFocus() is .Hwnd
				value = value.Tr('.')
			}
		super.Set(value)
		}
	}