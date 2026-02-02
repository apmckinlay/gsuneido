// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 'Static'
	ComponentName: 'Static'
	Unsortable: true
	Hasfocus?:	false

	New(.text = "", font = "", size = "", weight = "",
		justify = "LEFT", underline = false, color = ""
		/*, whitebgnd = false, status = "",
		sunken = false*/, tip = false, tabstop = false,
		bgndcolor = "", hidden = false, textStyle = false)
		{
		.SuSetHidden(hidden)
		.ComponentArgs = Object(text, font, size, weight, justify, underline, color,
			tip, tabstop, bgndcolor, hidden, textStyle)
		}

	Get()
		{
		return .text
		}

	Set(text, logfont = false, refreshRequired? = false)
		{
		text = String(text)
		if .text is text and logfont is false
			return
		.text = text
		.Act('Set', text, logfont, refreshRequired?)
		}

	SetColor(color)
		{
		.Act(#SetColor, color)
		}

	LBUTTONDOWN()
		{
		if 0 isnt .Send('Static_Click')
			return 0
		return 'callsuper' // for selecting text
		}

	On_Select_All()
		{
		.Act('SelectAll')
		}

	On_Context_Select_All()
		{
		.On_Select_All()
		}

	On_Context_Copy()
		{
		.Act('On_Copy')
		}

	ContextMenu(x, y)
		{
		if 0 isnt .Send('Static_ContextMenu', :x, :y)
			return

		ContextMenu(#('Copy', 'Select &All')).ShowCall(this, x, y)
		}
	}