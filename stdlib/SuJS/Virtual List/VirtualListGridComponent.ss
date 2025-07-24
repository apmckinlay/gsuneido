// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: 'VirtualListGrid'
	Xstretch: 1
	Ystretch: 1
	styles: `
		.su-vlist-container {
			position: relative;
			overflow: auto;
			border: 1px solid black;
			box-sizing: border-box;
		}
		.su-vlist-container:focus {
			outline: none;
		}
		.su-vlist-table {
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

	New()
		{
		LoadCssStyles('vlist-control.css', .styles)

		.CreateElement('div', className: 'su-vlist-container')
		.SetMinSize()
		.SetStyles([ 'background-color': 'white' ])
		.El.tabIndex = "0"

		.init()

		.El.AddEventListener('keydown', .keydown)
		.El.AddEventListener('scroll', .scroll)
		}

	OverrideHeight(height/*unused*/)
		{
		.SetStyles([ 'overflow': 'hidden' ])
		}

	OnFocus(event)
		{
		super.OnFocus()
		if event.target isnt .El or
			SuRender().Overlay.Status is #Closing and
				.keydownOverlayId is SuRender().OverlayId
			return
		.Event(#SETFOCUS)
		}

	OnBlur(event)
		{
		super.OnBlur()
		if event.target isnt .El or SuRender().Overlay.Status is #Opening
			return

		id = false
		if event isnt false
			try
				{
				focusEl = event.relatedTarget
				id = focusEl.Control().UniqueId
				}
		.Event(#KILLFOCUS, wParam: id)
		}

	init()
		{
		.table = CreateElement('table', .El, 'su-vlist-table')
		.table.SetAttribute('translate', 'no')
		.header = .Construct('ListHeader', .table)
		.body = .Construct('VirtualListGridBody', .table, .header)

		.header.Reset()
		.body.Reset()
		}

	UpdateHead(headCols, markCol = false, showSortIndicator = false)
		{
		.header.Update(headCols, markCol, virtualList?:, :showSortIndicator)
		.body.HeaderChanged()
		}

	keyMap: false
	keydownOverlayId: false
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
			" ": 		VK.SPACE,
//			PageUp: 	VK.PRIOR,
//			PageDown: 	VK.NEXT,
			Home:		VK.HOME,
			End:		VK.END,
			ArrowUp:	VK.UP,
			ArrowDown:	VK.DOWN,
			Enter:		VK.RETURN,
			Escape:		VK.ESCAPE,
			"+":		VK.ADD,
			"-":		VK.SUBTRACT,
			Delete:		Object(key: VK.DELETE, overlay?:),
			Insert:		Object(key: VK.INSERT, overlay?:)
			)

		if not .keyMap.Member?(event.key)
			return

		key = Object?(.keyMap[event.key]) ? .keyMap[event.key].key : .keyMap[event.key]
		overlay? = Object?(.keyMap[event.key]) and .keyMap[event.key].overlay?
		if overlay?
			{
			.EventWithOverlay(#KEYDOWN, key, 0,
				ctrl: event.ctrlKey, shift: event.shiftKey)
			.keydownOverlayId = SuRender().OverlayId
			}
		else
			.RunWhenNotFrozen({
				.EventWithFreeze(#KEYDOWN, key, 0,
					ctrl: event.ctrlKey, shift: event.shiftKey) })

		event.PreventDefault()
		event.StopPropagation()
		}

	lastScrollTop: 0
	scroll(event /*unused*/)
		{
		if .El isnt false and .lastScrollTop isnt .El.scrollTop
			.lastScrollTop = .El.scrollTop
		}

	SetReadOnly(readOnly)
		{
		.El.SetStyle('background-color',
			ToCssColor(readOnly ? CLR.ButtonFace : CLR.WHITE))
		}

	Recalc()
		{
		.body.Recalc()
		}

	SetColWidth(col, width)
		{
		.header.SetColWidth(col, width)
		}

	Default(@args)
		{
		event = args[0]
		.body[event](@+1 args)
		}

	BeforeResize()
		{
		.body.BeforeResize()
		}

	AfterResize()
		{
		.body.AfterResize()
		}

	Getter_(member)
		{
		return .body[member]
		}

	UpdateScroll()
		{
		if .El.scrollTop isnt .lastScrollTop
			.El.scrollTop = .lastScrollTop
		}

	Destroy()
		{
		.body.Destroy()
		.header.Destroy()
		super.Destroy()
		}
	}
