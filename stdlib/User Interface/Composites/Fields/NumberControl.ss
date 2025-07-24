// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
/* e.g.
Window(#(Vert
	(Field width: 5)
	(Number mask: '-###,###,###.##')
	(Number mask: false, width: 10)
	(Number mask: '##.##', width: 10)
	))
*/
HandleEnterControl
	{
	Name: 'Number'
	DefaultWidth: 15

	New(mask = '-###,###,###', .readonly = false,
		.rangefrom = false, .rangeto = false, width = false,
		set = false, .mandatory = false, status = "", justify = "RIGHT",
		font = "", size = "", weight = "", underline = false,
		hidden = false, tabover = false)
		{
		super(@.options(:mask, :width, :readonly, :justify, :status, :font, :size,
			:weight, :underline, :hidden, :tabover))
		.SubClass()
		.prev = Object(text: "", begin: Object(x: 0), end: Object(x: 0))
		if set isnt false
			{
			.Set(set)
			.Send("NewValue", .Get())
			}
		}
	options(@args)
		{
		.mask = args.mask = .RetrieveMask(args.mask)
		return args
		}
	RetrieveMask(mask)
		{
		if String?(mask) and mask =~ "^-?[A-Z][a-zA-Z_0-9]*$"
			mask = mask.Extract('-?') $ Global(mask.Replace('-', ''))
		return mask
		}
	SetFontAndSize(font, size, weight, underline, width, height/*unused*/)
		{ // overrides EditControl
		// NOTE: using '9' assumes digits are all the same width (normal)
		// NOTE: ignores width if mask is set
		if .mask isnt false
			super.SetFontAndSize(font, size, weight, underline, 1, 1,
				text: .mask.Tr('#', '9'))
		else
			super.SetFontAndSize(font, size, weight, underline, width, 1, text: '9')
		}

	selEnd: 999
	EN_SETFOCUS(skipSetText? = false)
		{
		if skipSetText? is false
			{
			SetWindowText(.Hwnd, GetWindowText(.Hwnd).Tr(','))
			.SendMessage(EM.SETSEL, 0, .selEnd)
			}
		return super.EN_SETFOCUS()
		}
	NumberFormat: 'Format'
	KillFocus()
		{
		s = .GetUnvalidated()
		if s isnt '' and .Valid?()
			{
			dirty? = .Dirty?()
			if .mask isnt false and s.Number?() // Valid? doesn't check this when readonly
				s = Number(s)[.NumberFormat](.mask)
			SetWindowText(.Hwnd, s)
			.SendMessage(EM.SETSEL, 0, .selEnd)
			.Dirty?(dirty?)
			}
		}
	Valid?()
		{
		s = .GetUnvalidated()
		if .readonly and s isnt "#" and s isnt "-"
			return true
		return .validCheck?(s, .mandatory, .rangefrom, .rangeto)
		}
	validCheck?(data, mandatory, rangefrom, rangeto)
		{
		valid = true
		if data is ''
			valid = not mandatory
		else if not Number?(data) and not data.Number?()
			valid =  false
		else
			{
			try
				n = String?(data) ? Number(data) : data
			catch (unused, "can't convert String to number")
				return false
			if .notInRange(n, rangefrom, rangeto) or IsInf?(n)
				valid = false
			}
		return valid
		}
	notInRange(n, rangefrom, rangeto)
		{
		return (rangefrom isnt false and n < rangefrom) or
			(rangeto isnt false and n > rangeto)
		}
	ValidData?(@args)
		{
		value = args[0]
		return .validCheck?(value, args.GetDefault('mandatory', false),
			args.GetDefault('rangefrom', false), args.GetDefault('rangeto', false))
		}

	Get()
		{
		s = .GetUnvalidated()
		return s isnt '' and .Valid?() ? Number(s) : ""
		// return "" for invalid instead of s so that rules work
		}
	GetUnvalidated()
		{
		s = .GetConvertedText()
		if .isSimpleMathExpression?(s)
			try s = String(s.Eval()) // relying on condition to keep this Eval safe
		return s
		}
	isSimpleMathExpression?(s)
		{
		return not s.Blank?() and s.Tr(" ()0-9eE.+*/-") is "" and s.Has1of?("-+/*")
		}
	GetConvertedText()
		{
		return super.Get().Tr(',').Trim()
		}
	Set(value)
		{
		super.Set(.ConvertValue(value, { it.Tr(',') }))
		}

	SetRange(low, high)
		{
		.rangefrom = low
		.rangeto = high
		}

	ConvertValue(value, convert)
		{
		if Number?(value) or (String?(value) and value.Number?())
			{
			value = .mask is false ? String(value) : Number(value)[.NumberFormat](.mask)
			if GetFocus() is .Hwnd
				value = convert(value)
			}
		return value
		}
	}
