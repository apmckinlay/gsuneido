// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: "List"
	Xmin: 100
	Ymin: 100
	Xstretch: 1
	Ystretch: 1
	styles: `
		.su-list-container {
			position: relative;
			overflow: auto;
			border: 1px solid black;
		}
		.su-list-container:focus {
			outline: none;
		}
		.su-list-table {
			position: absolute;
			top: 0px;
			left: 0px;
			width: 100%;
			height: 100%;
			table-layout: fixed;
			user-select: none;
			border: none;
			border-spacing: 0;
		}
		`

	New(noShading = false, indicateHovered = false, noDragDrop = false,
		noHeaderButtons = false, .stretch = false, noHeader = false)
		{
		LoadCssStyles('list-control.css', .styles)

		.CreateElement('div', className: 'su-list-container')
		.SetMinSize()
		.SetStyles([ 'background-color': 'white' ])
		.El.tabIndex = "0"

		.init(noShading, indicateHovered, noDragDrop, noHeaderButtons, noHeader)

		.El.AddEventListener('keydown', .keydown)
		.El.AddEventListener('focus', .focus)
		.El.AddEventListener('blur', .blur)
		}

	focus(event)
		{
		if event.target isnt .El
			return
		.Event(#SETFOCUS)
		}

	blur(event)
		{
		if event.target isnt .El
			return
		.body.KillFocus()
		.EventWithFreeze('LIST_KILLFOCUS')
		return
		}

	init(noShading, indicateHovered, noDragDrop, noHeaderButtons, noHeader)
		{
		.table = CreateElement('table', .El, 'su-list-table')
		.table.SetAttribute('translate', 'no')
		.header = .Construct('ListHeader',
			.table, :noDragDrop, :noHeaderButtons, stretch: .stretch, :noHeader)
		.body = .Construct('ListBody', .table, .header, noShading, indicateHovered)

		.header.Reset()
		.body.Reset()
		}

	UpdateHead(headCols, markCol = false, data = false, showSortIndicator = false)
		{
		.header.Update(headCols, markCol, :showSortIndicator)
		if data isnt 'skip'
			.body.HeaderChanged(data)
		}

	Default(@args)
		{
		event = args[0]
		.body[event](@+1 args)
		}

	Getter_(member)
		{
		return .body[member]
		}

	SetMaxWidth(col)
		{
		.header.SetMaxWidth(col)
		}

	SetColWidth(col, width)
		{
		.header.SetColWidth(col, width)
		}

	readOnly: false
	grayOut: true
	SetReadOnly(.readOnly, .grayOut = true)
		{
		.El.SetStyle('background-color',
			ToCssColor(.readOnly and .grayOut or .enabled is false
				? CLR.ButtonFace
				: CLR.WHITE))
		}

	enabled: true
	SetEnabled(.enabled)
		{
		.El.SetStyle('pointer-events', enabled is false ? 'none' : '')
		.SetReadOnly(.readOnly, .grayOut)
		}

	GetEnabled()
		{
		return .enabled
		}

	keyMap: false
	keydown(event)
		{
		// input in ListEditWindow
		if event.target isnt .El or SuRender().Frozen?()
			return

		if .keyMap is false
			.keyMap = Object(
			F2: 		VK.F2,
			F5: 		VK.F5,
			F8: 		VK.F8,
			Insert: 	VK.INSERT,
			" ": 		VK.SPACE,
//			PageUp: 	VK.PRIOR,
//			PageDown: 	VK.NEXT,
			Home:		VK.HOME,
			End:		VK.END,
			ArrowUp:	VK.UP,
			ArrowDown:	VK.DOWN,
			Delete:		VK.DELETE)

		if not .keyMap.Member?(event.key)
			return
		.EventWithFreeze(#KEYDOWN, .keyMap[event.key], 0,
			ctrl: event.ctrlKey, shift: event.shiftKey)
		event.PreventDefault()
		}

	Destroy()
		{
		.body.Destroy()
		.header.Destroy()
		super.Destroy()
		}
	}
