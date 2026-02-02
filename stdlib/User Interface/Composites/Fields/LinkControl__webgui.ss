// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
FieldControl
	{
	Name: 'Link'
	ComponentName: 'Link'
	New(width = 25, set = false, mandatory = false, tabover = false, hidden = false,
		readonly = false)
		{
		super(underline:, :width, :set, :mandatory, :tabover, :hidden, :readonly)
		.updateTextColor()
		}

	EN_CHANGE(text)
		{
		super.EN_CHANGE(text)
		.updateTextColor()
		}
	SetReadOnly(readonly)
		{
		super.SetReadOnly(readonly)
		.updateTextColor()
		}

	SetEnabled(enabled)
		{
		super.SetEnabled(enabled)
		.updateTextColor()
		}

	curColor: false
	updateTextColor()
		{
		color = .GetReadOnly() is true or .GetEnabled() isnt true
			? false
			: .ValidLink?(.Get()) ? CLR.BLUE : CLR.RED
		if color is .curColor
			return
		.Act('SetTextColor', .curColor = color)
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