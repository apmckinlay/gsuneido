// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
NumberControl
	{
	Name: 'Dollar'
	EN_SETFOCUS()
		{
		SetWindowText(.Hwnd, .GetUnvalidated())
		return super.EN_SETFOCUS()
		}
	NumberFormat: 'DollarFormat'
	GetConvertedText()
		{
		return .convert(GetWindowText(.Hwnd))
		}

	convert(value)
		{
		value = value.Tr(',').Trim()
		value = value.Replace('^\$', '').Trim()
		if value.Prefix?('(') and value.Suffix?(')')
			value = value.Replace('^\(', '-').Replace('\)$', '')
		return value
		}

	Set(value)
		{
		super.Set(.ConvertValue(value, .convert))
		}
	}