// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
ScintillaComponent
	{
	New(readonly, .height, tabthrough, .width, .argXmin, .argYmin, .componentAddons = #())
		{
		super(readonly, height, tabthrough)
		.AddEventListenerToCM('inputRead', .onInput)
		if .componentAddons.NotEmpty?()
			.AddEventListenerToCM('beforeChange', .onBeforeChange)
		.defers = Object()

		.addons = AddonManager(this, componentAddons)
		.addons.Send(#Init)
		}

	updateUIDelay: false
	UpdateUI()
		{
		super.UpdateUI()
		if .componentAddons.Empty?()
			return

		if .updateUIDelay isnt false
			.updateUIDelay.Kill()
		.updateUIDelay = SuDelayed(0)
			{
			.updateUIDelay = false
			.addons.Send(#UpdateUI)
			}
		}

	onBeforeChange(unused, changeObj)
		{
		if changeObj.GetDefault(#origin, false) is '+delete'
			{
			from = .CM.IndexFromPos(changeObj.from)
			to = .CM.IndexFromPos(changeObj.to)
			.addons.Send(#BeforeDelete, from, to - from)
			}
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
			.Xmin = SuRender().GetTextMetrics(.CMEl,
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
			control: event.GetDefault(#ctrlKey, false),
			shift: event.GetDefault(#shiftKey, false),
			alt: event.GetDefault(#altKey, false))
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
				{
				.Event(#CHARADDED, c = changeObj.text.Join('\r\n'))
				.addons.Send(#CharAdded, c)
				}
			}
		}

	// methods to support Addons running on browsers
	GetLength()
		{
		return .CM.GetValue().Size()
		}

	GetSelText()
		{
		return .CM.GetSelection()
		}

	GetSelect()
		{
		from = .CM.IndexFromPos(.CM.GetCursor('from'))
		to = .CM.IndexFromPos(.CM.GetCursor('to'))
		return Object(cpMin: from, cpMax: to)
		}

	GetAt(pos)
		{
		if pos < 0 or pos >= .GetLength()
			return '\x00'
		c = .GetRange(pos, pos + 1, '\n')
		return c is '' ? '\r' : c
		}

	GetRange(start, end, sep = '\r\n') // gets from start to end - 1
		{
		startPos = .CM.PosFromIndex(start)
		endPos = .CM.PosFromIndex(end)
		return .CM.GetRange(startPos, endPos, sep)
		}

	GetLine(line = false)
		{
		if line is false
			line = .LineFromPosition()
		return line >= .CM.LineCount() or line < 0
			? ''
			: .CM.GetLine(line)
		}

	LineFromPosition(pos = false)
		{
		if pos is false
			return .CM.GetCursor('from').line
		return .CM.PosFromIndex(pos).line
		}

	GetCurrentWord()
		{
		wordChars = .GetWordChars()
		org = end = .GetCurrentPos()
		while wordChars.Has?(.GetAt(org - 1))
			--org
		while wordChars.Has?(.GetAt(end))
			++end
		return org < end ? .GetRange(org, end) : ""
		}

	GetCurrentPos()
		{
		return .CM.IndexFromPos(.CM.GetCursor('head'))
		}

	SetSelect(i, n = 0)
		{
		.SetSel(i, i + n)
		}

	On_Delete()
		{
		.DELETE()
		}

	Defer(block)
		{
		defer = false
		defer = SuDelayed(0, { /*.defers.Remove1(defer); */block(); })
		.defers.Add(defer)
		}

	GetColumn(i)
		{
		pos = .CM.PosFromIndex(i)
		return pos.ch
		}

	BraceMatch(pos)
		{
		return BraceMatch(.CM.GetValue(), pos)
		}

	SetHighlightGuide(@unused) {}
	StyleToEnd() {}

	Destroy()
		{
		if .updateUIDelay isnt false
			{
			.updateUIDelay.Kill()
			.updateUIDelay = false
			}

		while not Same?(.defers, k = .defers.PopLast())
			k.Kill()

		super.Destroy()
		}
	}
