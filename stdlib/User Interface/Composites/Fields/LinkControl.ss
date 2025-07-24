// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	Name: 'Link'
	New(width = 25, set = false, mandatory = false, tabover = false, hidden = false,
		readonly = false)
		{
		super(underline:, :width, :set, :mandatory, :tabover, :hidden, :readonly)
		.SubClass()
		}
	CTLCOLOREDIT(wParam)
		{
		SetTextColor(wParam, .ValidLink?(.Get()) ? CLR.BLUE : CLR.RED)
		return super.CTLCOLOREDIT(wParam)
		}
	ValidLink?(unused)
		{
		return true
		}
	Prefix: ''
	LBUTTONDBLCLK()
		{
		.GoToLink()
		}
	GoToLink()
		{
		adr = .Get()
		if adr isnt '' and .ValidLink?(adr)
			ShellExecute(.WindowHwnd(), 'open', .MergePrefix(adr))
		return 0
		}
	MergePrefix(adr)
		{
		return .Prefix $ adr
		}
	}
