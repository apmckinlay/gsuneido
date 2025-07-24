// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	Name: 'CapitalizeWords'
	KillFocus()
		{
		if not .Dirty?()
			return
		s = super.Get()
		if s is ""
			return
		SetWindowText(.Hwnd, s.CapitalizeWords(lower: false))
		.Dirty?(true)
		.SelectAll()
		}
	}