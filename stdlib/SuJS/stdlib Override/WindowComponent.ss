// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
WindowBaseComponent
	{
	styles: `
		.su-window-anchor {
			padding: 0;
			display: inline-block;
			border: none;
			border-spacing: 0px;
		}

		.su-window-anchor::backdrop {
			background-color: rgba(256,256,256,0.5);
		}

		.su-window-container {
			position: fixed;
			background-color: var(--su-color-buttonface);
			display: inline-block;
			border: 1px solid lightblue;
			min-width: 0px;
			height: auto;
			box-shadow: 3px 3px 3px grey;
		}

		.su-window-content {
			display: flex;
			overflow: auto;
		}

		.su-window-menubar {
			display: flex;
			flex-direction: row;
			background-color: white;
			user-select: none;
			height: 1.5em;
		}

		.su-window-menu-item {
			padding: 0px 5px;
		}

		.su-window-menu-item:hover {
			background-color: lightblue;
		}

		.su-window-header {
			display: flex;
			flex-direction: row;
			align-items: center;
			height: 1.5em;
			justify-content: space-between;
			background-color: var(--su-color-windowheader);
		}
		.su-window-title {
			display: inline;
			user-select: none;
			padding-left: 5px;
			line-height: 1.5em;
			cursor: default;
		}
		.su-window-buttons {
		}
		.su-window-button,
		.su-window-close {
			padding: 0 1em;
			display: inline-block;
			color:	#aaa;
			font-weight: bold;
			line-height: 1.5em;
			font-family: suneido;
			font-style: normal;
			font-weight: normal;
			user-select: none;
			cursor: default;
		}
		.su-window-button:hover,
		.su-window-button:focus {
			color: black;
			background-color: darkgrey;
			text-decoration: none;
		}
		.su-window-close:hover,
		.su-window-close:focus {
			color: white;
			background-color: red;
			text-decoration: none;
		}
		.su-window-left-resize {
			position: absolute;
			left: -5px;
			top: 0px;
			width: 7px;
			height: 100%;
			cursor: ew-resize;
		}
		.su-window-right-resize {
			position: absolute;
			right: -5px;
			top: 0px;
			width: 7px;
			height: 100%;
			cursor: ew-resize;
		}
		.su-window-top-resize {
			position: absolute;
			top: -5px;
			left: 0px;
			height: 7px;
			width: 100%;
			cursor: ns-resize;
		}
		.su-window-bottom-resize {
			position: absolute;
			bottom: -5px;
			left: 0px;
			height: 7px;
			width: 100%;
			cursor: ns-resize;
		}
		.su-window-left-top-resize {
			position: absolute;
			top: -5px;
			left: -5px;
			width: 7px;
			height: 7px;
			cursor: nwse-resize;
		}
		.su-window-right-top-resize {
			position: absolute;
			top: -5px;
			right: -5px;
			width: 7px;
			height: 7px;
			cursor: nesw-resize;
		}
		.su-window-right-bottom-resize {
			position: absolute;
			bottom: -5px;
			right: -5px;
			width: 7px;
			height: 7px;
			cursor: nwse-resize;
		}
		.su-window-left-bottom-resize {
			position: absolute;
			bottom: -5px;
			left: -5px;
			width: 7px;
			height: 7px;
			cursor: nesw-resize;
		}`

	CallClass(@args)
		{
		_ctrlspec = args
		new this(@args)
		}

	AnchorElement: 'div'
	New(control, title = false,
		x/*unused*/ = false, y/*unused*/ = false,
		w/*unused*/ = false, h/*unused*/ = false,
		exStyle/*unused*/ = 0, style = 0, wndclass/*unused*/ = "SuBtnfaceArrow",
		.show = true, .newset = false, .exitOnClose = false,
		.keep_placement = false, .border = 0, .onDestroy/*unused*/ = false,
		parentHwnd = 0, useDefaultSize/*unused*/ = false, menubar = #())
		{
		.xmin0 = .Xmin
		.ymin0 = .Ymin
		LoadCssStyles('su_window.css', .styles)


		.ParentEl = parentHwnd is 0
			? SuUI.GetCurrentDocument().body
			: SuRender().GetRegisteredComponent(parentHwnd).El

		.Window = .Controller = this
		.HwndMap = SuUI.HtmlElMap(this) // Element map
		.AnchorEl = .setupEl(CreateElement(.AnchorElement, .ParentEl,
			className: "su-window-anchor"))
		if style is 0
			style = WS.OVERLAPPEDWINDOW
		.containerEl = .setupEl(CreateElement('div', .AnchorEl,
			className: "su-window-container"))
		.setupHeader(title, style)
		.setupMenubar(menubar)
		.El = .setupEl(CreateElement('div', .containerEl))
		.SetStyles(Object('padding': border $ 'px', 'box-sizing': 'border-box'))


		.Ctrl = .Construct(control)

		.setupResize(style)
		.Recalc()
		.InitUniqueId()
		.RunPendingRefresh()
		.RegisterActiveWindow()

		DoStartup(.Ctrl)
		}

	setupEl(el)
		{
		el.Control(this)
		el.Window(this)
		return el
		}

	Startup()
		{
		// so the initial Dialog/ModalWindow still saves proper size info
		// when closing without been resized
		.Event(#WINDOWRESIZE, SuRender.GetClientRect(.El))
		}

	defaultButton: false
	SetDefaultButton(uniqueId)
		{
		.defaultButton = SuRender().GetRegisteredComponent(uniqueId)
		.HighlightDefaultButton(true)
		}
	CallDefaultButton()
		{
		if .defaultButton is false or .defaultButton.GetEnabled() is false
			return
		.defaultButton.CLICKED()
		}
	HighlightDefaultButton(highlight?)
		{
		if .defaultButton is false or .defaultButton.GetEnabled() is false
			return
		.defaultButton.Highlight(highlight?)
		}

	headerEl: false
	titleEl: false
	buttonsEl: false
	minimizeEl: false
	maximizeEl: false
	closeEl: false
	headerHeight: 32
	DisableMinimize?: false
	NonClientHeight: 0
	setupHeader(.title, style)
		{
		if ((style & WS.POPUP) isnt 0)
			return

		.headerEl = CreateElement('div', .containerEl, className: "su-window-header")
		.titleEl = CreateElement('div', .headerEl, className:  "su-window-title")
		if .title isnt false
			.titleEl.innerHTML = .title

		.buttonsEl = CreateElement('div', .headerEl, className: 'su-window-buttons')
		.buttonsEl.SetAttribute('translate', 'no')
		if ((style & WS.MINIMIZEBOX) isnt 0 and .DisableMinimize? is false)
			{
			.minimizeEl = CreateElement('div', .buttonsEl, className: "su-window-button")
			.minimizeEl.textContent = IconFontHelper.GetCode('collapse').Chr()
			.minimizeEl.AddEventListener('click', .MINIMIZE)
			}
		if ((style & WS.MAXIMIZEBOX) isnt 0)
			{
			.maximizeEl = CreateElement('div', .buttonsEl, className: "su-window-button")
			.maximizeEl.textContent = IconFontHelper.GetCode('square').Chr()
			.maximizeEl.AddEventListener('click', .MAXIMIZE)
			.headerEl.AddEventListener('dblclick', .doubleClick)
			}
		if ((style & WS.SYSMENU) isnt 0)
			{
			.closeEl = CreateElement('div', .buttonsEl, className: "su-window-close")
			.closeEl.textContent = IconFontHelper.GetCode('delete').Chr()
			.closeEl.AddEventListener('click', .CLOSE)
			}
		.headerEl.Window(this)
		.headerEl.AddEventListener('mousedown', .moveMouseDown)
		.headerEl.AddEventListener('contextmenu', .headerContextMenu)
		.NonClientHeight += 1.5 /*=hearder height*/
		}

	style: 0
	headerContextMenu(event)
		{
		status = Object(
			restore: .State isnt WindowPlacement.normal,
			maximize: .maximizeEl isnt false and .State isnt WindowPlacement.maximized,
			close: .closeEl isnt false)
		extra = .snapMgr.ContextMenu
		.RunWhenNotFrozen(
			{ .EventWithOverlay('HeaderContextMenu', status, extra,
				event.clientX, event.clientY) })
		event.StopPropagation()
		event.PreventDefault()
		}

	menuEl: false
	setupMenubar(menubar)
		{
		if menubar.Empty?()
			return

		.menuEl = CreateElement('div', .containerEl, className: 'su-window-menubar')
		for menu in menubar
			{
			el = CreateElement('div', .menuEl, className: 'su-window-menu-item')
			el.textContent = menu
			el.AddEventListener('click', .eventFactory(el, menu, .onMenubar))
			el.AddEventListener('focus', .menuFocus)
			el.tabIndex = "-1"
			}
		.menuEl.Window(this)
		.NonClientHeight += 1.5 /*=menubar height*/
		}

	menuFocus(event)
		{
		el = false
		try
			el = event.relatedTarget
		if false isnt control = .GetControlFromEl(el)
			.Event('SyncPrevMenuFocus', control.UniqueId)
		}

	onMenubar(el, menu, event/*unused*/)
		{
		r = SuRender.GetClientRect(el)
		.RunWhenNotFrozen({ .EventWithOverlay('MenuBar', menu, r.left, r.bottom) })
		}

	ExtraContextCall(args)
		{
		.snapMgr.ContextCall(args, this)
		}

	moveMouseDown(event)
		{
		if event.target is .maximizeEl or
			event.target is .minimizeEl or
			event.target is .closeEl or
			event.button isnt 0
			return

		.origPos = .curPos = Object(x: event.x, y: event.y)
		.StartMouseTracking(.moveMouseUp, .moveMouseMove)
		}

	moveMouseMove(event)
		{
		viewportRect = SuRender.GetClientRect()
		if .State isnt WindowPlacement.normal
			{
			left = event.x - .curPos.x / viewportRect.width * .windowRect.width
			.restore(:left, top: 0)
			}
		containerRect = SuRender.GetClientRect(.containerEl)
		.moveHorz(event, containerRect, viewportRect)
		.moveVert(event, containerRect, viewportRect)
		.curPos = Object(x: event.x, y: event.y)
		event.PreventDefault()
		}

	moveHorz(event, containerRect, viewportRect)
		{
		if event.x < viewportRect.left or event.x > viewportRect.right
			return
		diff = event.x - .curPos.x
		.containerEl.SetStyle(#left, (containerRect.left + diff) $ 'px')
		}

	moveVert(event, containerRect, viewportRect)
		{
		if event.y < viewportRect.top or event.y > viewportRect.bottom
			return
		diff = event.y - .curPos.y
		.containerEl.SetStyle(#top, (containerRect.top + diff) $ 'px')
		}

	moveMouseUp(event/*unused*/)
		{
		if .curPos isnt .origPos
			{
			.updateWindowRect()
			.syncWindowPlacement()
			}
		.curPos = false
		.StopMouseTracking()
		}

	resizes: #()
	setupResize(style)
		{
		.resizes = Object()
		if ((style & WS.SIZEBOX) is 0)
			return

		if .Ctrl.Xstretch > 0
			{
			.makeResize([#left])
			.makeResize([#right])
			}
		if .Ctrl.Ystretch > 0
			{
			.makeResize([#top])
			.makeResize([#bottom])
			}
		if .Ctrl.Xstretch > 0 and .Ctrl.Ystretch > 0
			{
			.makeResize([#left, #top])
			.makeResize([#right, #top])
			.makeResize([#right, #bottom])
			.makeResize([#left, #bottom])
			}
		.El.className = "su-window-content"
		.Ctrl.SetStyles(#(width: '100%', height: '100%'))
		}

	makeResize(sides)
		{
		el = CreateElement('div', .containerEl,
			className: 'su-window-' $ sides.Join('-') $ '-resize')
		el.AddEventListener('mousedown', .eventFactory(el, sides, .resizeMouseDown))
		.resizes.Add(el)
		}

	eventFactory(el, sides, handler)
		{
		return { |event| handler(el, sides, event) }
		}

	getter_snapMgr()
		{
		return SuRender().SnapManager
		}
	GetResizes()
		{
		return .resizes
		}

	curPos: false
	dirs: (
		left: 	[#x, -1],
		right:	[#x, 1],
		top:	[#y, -1],
		bottom:	[#y, 1])
	resizeMouseDown(el, sides, event)
		{
		if event.target isnt el or event.button isnt 0
			return
		.curPos = Object()
		for side in sides
			.curPos[.dirs[side][0]] = event[.dirs[side][0]]
		.StartMouseTracking(.resizeMouseUp, .eventFactory(el, sides, .resizeMouseMove))
		}

	resizeMouseMove(el/*unused*/, sides, event)
		{
		viewportRect = SuRender.GetClientRect()
		newPos = .curPos.Copy()
		for side in sides
			{
			pos = event[.dirs[side][0]]
			prePos = .curPos[.dirs[side][0]]
			if .resize(side, pos, prePos, viewportRect)
				newPos[.dirs[side][0]] = pos
			}
		.curPos = newPos
		event.PreventDefault()
		}

	resize(side, pos, prePos, viewportRect)
		{
		if ((viewportRect[side] - pos) * .dirs[side][1] < 0)
			return false

		return side in (#left, #right)
			? .resizeHorz(side, pos, prePos)
			: .resizeVert(side, pos, prePos)
		}

	resizeHorz(side, pos, prePos)
		{
		offset = pos - prePos
		rect = .El.GetBoundingClientRect()
		newWidth = rect.width + .dirs[side][1] * offset
		if newWidth < Max(.Ctrl.Xmin + .border * 2, .getHeaderXmin())
			return false

		if .snapMgr.HorzResize(newWidth, this) isnt true
			{
			.containerEl.SetStyle('width', '')
			.El.SetStyle('width', newWidth $ 'px')
			if side is #left
				.containerEl.SetStyle('left', pos $ 'px')
			}
		return true
		}

	resizeVert(side, pos, prePos)
		{
		offset = pos - prePos
		rect = .El.GetBoundingClientRect()
		newHeight = rect.height + .dirs[side][1] * offset
		if newHeight < .Ctrl.Ymin
			return false

		.El.SetStyle('height', newHeight $ 'px')
		if side is #top
			.containerEl.SetStyle('top', pos $ 'px')
		return true
		}

	headerXmin: false
	getHeaderXmin()
		{
		if .headerXmin isnt false
			return .headerXmin

		width = 0
		if .titleEl isnt false and .title isnt false
			width += SuRender().GetTextMetrics(.titleEl, .title).width
		if .buttonsEl isnt false
			width += .buttonsEl.clientWidth
		return .headerXmin = width
		}

	resizeMouseUp(event/*unused*/)
		{
		.curPos = false
		.StopMouseTracking()

		if .State is WindowPlacement.normal
			{
			.updateWindowRect()
			.syncWindowPlacement()
			.Event(#WINDOWRESIZE, SuRender.GetClientRect(.El)) // to WindwoBase
			}
		.TopDownWindowResize()
		}

	Recalc()
		{
		if .Ctrl is false
			return
		padding = .border * 2
		.Xmin = Max(.xmin0, .Ctrl.Xmin) + padding
		.Ymin = Max(.ymin0, .Ctrl.Ymin) + padding
		.SetMinSize()
		.updateWindowRect()
		}

	SetTitle(.title)
		{
		.titleEl.innerHTML = .title
		.headerXmin = false
		}

	GetContainerEl()
		{
		return .containerEl
		}

	CLOSE()
		{
		.RunWhenNotFrozen({ .EventWithFreeze('CLOSE') })
		}

	doubleClick(event)
		{
		if event.target is .maximizeEl or
			event.target is .minimizeEl or
			event.target is .closeEl
			return
		.MAXIMIZE()
		}

	State: 0 // = WindowPlacement.normal
	previousState: 0
	windowRect: false
	MAXIMIZE(forceRestore? = false)
		{
		if .State isnt WindowPlacement.normal or forceRestore? is true
			.restore()
		else
			.maximize()
		.syncWindowPlacement()
		.TopDownWindowResize()
		}

	maximize()
		{
		.SetState(WindowPlacement.maximized)
		.snapMgr.Remove(this)
		.SetStyles(
			Object(left: '0px', top: '0px', width: '100%',
				height: 'calc(100% - (' $ SuRender().Taskbar.GetTaskbarHeight() $ '))'),
			.containerEl)
		.SetStyles(
			Object(width: '100%', height: 'calc(100% - ' $ .NonClientHeight $ 'em)'))
		.resizes.Each({ it.SetStyle('display', 'none') })
		}

	MINIMIZE()
		{
		if .minimizeEl is false
			return
		.SetState(WindowPlacement.minimized)
		.snapMgr.Remove(this)
		.SetStyles(#(display: 'none'), .containerEl)
		.syncWindowPlacement()
		.EventWithOverlay('AfterMinimized')
		}

	UpdateMaximize()
		{
		if .State is WindowPlacement.normal
			return
		.SetStyles(
			Object(
				height: 'calc(100% - (' $ SuRender().Taskbar.GetTaskbarHeight() $ '))'),
			.containerEl)
		}

	restore(left = false, top = false)
		{
		.SetState(WindowPlacement.normal)
		.snapMgr.Remove(this)
		viewPort = .getViewPort()
		width = Max(Min(.windowRect.width, viewPort.width - 2), .Ctrl.Xmin + .border * 2)
		height = Max(Min(.windowRect.height, viewPort.height - 2),
			.Ctrl.Ymin + .border * 2)
		if top is false
			top = Max(Min(.windowRect.top, viewPort.bottom - height), 0)
		if left is false
			left = Max(Min(.windowRect.left, viewPort.right - width), 0)
		.SetStyles(Object(left: left $ 'px', top: top $ 'px', right: '',
			width: '', height: ''), .containerEl)
		.SetStyles(Object(width: width $ 'px',
			height: 'calc(' $ height $ 'px - ' $ .NonClientHeight $ 'em)'))
		.resizes.Each({ it.SetStyle('display', 'initial') })
		}

	getViewPort()
		{
		viewPort = SuRender.GetClientRect()
		taskBar = SuRender().Taskbar.GetDimension()
		viewPort.height -= taskBar.height
		viewPort.bottom -= taskBar.height
		return viewPort
		}

	SetState(state)
		{
		.previousState = .State
		.State = state
		if .maximizeEl isnt false
			.maximizeEl.textContent = IconFontHelper.GetCode(
				state isnt WindowPlacement.normal ? 'restore' : 'square').Chr()
		}

	updateWindowRect()
		{
		if .State isnt WindowPlacement.normal
			return
		.windowRect = .getWindowRect()
		}

	getWindowRect()
		{
		windowRect = SuRender.GetClientRect(.containerEl)
		windowRect.width -= 2 /*=border*/
		windowRect.height -= 2 /*=border*/
		return windowRect
		}

	syncWindowPlacement()
		{
		.Event(#SYNCWINDOWPLACEMENT, rect: .windowRect,
			maximized: .State is WindowPlacement.maximized,
			minimized: .State is WindowPlacement.minimized)
		}

	SetWindowPlacement(place)
		{
		.windowRect = place.rcNormalPosition
		if place.showCmd is SW.SHOWMAXIMIZED
			.maximize()
		else
			.restore()
		}

	Center()
		{
		super.Center()
		.updateWindowRect()
		.syncWindowPlacement()
		}

	PlaceActive()
		{
		if .State is WindowPlacement.minimized
			{
			.SetState(.previousState)
			.SetStyles(#(display: 'initial'), .containerEl)
			.syncWindowPlacement()
			}
		if .State isnt WindowPlacement.normal
			return
		oldWindowRect = .windowRect
		.restore() // force window into viewport
		.updateWindowRect()
		if .windowRectChanged?(oldWindowRect, .windowRect)
			.syncWindowPlacement()
		}

	// Due to the border width difference with DPI scale > 100%,
	// the real window size could have a small difference from the set size even if
	// users don't move/resize the window. Allow 2px diff to avoid unnecessary saves
	windowRectChanged?(old, cur)
		{
		for m in #(left, right, top, bottom, width, height)
			{
			oldValue = old.GetDefault(m, -999/*=an impossible number*/)
			curValue = cur.GetDefault(m, -999/*=an impossible number*/)
			if ((oldValue - curValue).Abs() > 2)
				return true
			}
		return false
		}

	Destroy()
		{
		.snapMgr.Remove(this)
		.defaultButton = false
		if .Member?(#Ctrl)
			.Ctrl.Destroy()
		.El = .AnchorEl
		if .headerEl isnt false
			.headerEl.Window(false)
		if .menuEl isnt false
			.menuEl.Window(false)
		super.Destroy()
		}
	}
