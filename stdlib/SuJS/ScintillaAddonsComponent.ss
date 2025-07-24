// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
ScintillaComponent
	{
	New(readonly, .height, tabthrough, .width, .argXmin, .argYmin)
		{
		super(readonly, height, tabthrough)
		.AddEventListenerToCM('inputRead', .onInput)
		}

	padding: 10 // (4px padding + 1px border) * 2
	Recalc()
		{
		super.Recalc()
		if .argYmin isnt false
			.Ymin = .argYmin
		else
			.Ymin = .CM.DefaultTextHeight() *
				(.height is false ? 7/*=default height*/ : .height) + .padding
		if .argXmin isnt false
			.Xmin = .argXmin
		else
			.Xmin = SuRender().GetTextMetrics(.El,
				'M'.Repeat(.width is false ? 60/*=default width*/ : .width)).width
		.SetMinSize()
		}

	ScrollToBottom(noFocus? = false)
		{
		.CM.ScrollIntoView(Object(ch: 0, line: .CM.LastLine()))
		if noFocus? is false
			.CM.Focus()
		}

	addonCommands: #()
	SetAddonCommands(commands)
		{
		.addonCommands = commands.Map({ it.Replace('ctrl', 'control') })
		}

	OnKeyDown(cm, event)
		{
		pressed = Object(
			control: event.ctrlKey, shift: event.shiftKey, alt: event.altKey)
		if event.key is "Enter"
			.Event(#Enter_Pressed, :pressed)
		if event.key is "Backspace"
			.Event(#Backspace_Pressed)
		EditorKeyDownComponentHandler(this, event, pressed, extraCommands: .addonCommands)

		super.OnKeyDown(cm, event)
		}

	onInput(unused, changeObj)
		{
		try
			{
			if changeObj.origin is "+input"
				.Event(#CHARADDED, changeObj.text.Join('\r\n'))
			}
		}
	}
