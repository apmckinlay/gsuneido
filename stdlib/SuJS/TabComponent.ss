// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: 'Tab'
	ContextMenu: true
	styles: '
		.su-tab-control {
			position: relative;
		}
		.su-tab-container {
			display: flex;
			position: absolute;
			top: 0px;
			left: 0px;
			width: 100%;
			height: 100%;
		}
		.su-tab-container.su-tab-top {
			border-bottom: 1px solid lightgrey;
		}
		.su-tab-container.su-tab-bottom {
			border-top: 1px solid lightgrey;
		}
		.su-tab {
			border-right: 1px solid lightgrey;
			outline: none;
			cursor: pointer;
			padding: 5px;
			transition: 0.3s;
			display: flex;
			align-items: baseline;
			overflow: hidden;
			font-weight: normal;
		}
		.su-tab-top .su-tab {
			border-top: 1px solid lightgrey;
			border-top-left-radius: 0.5em;
			border-top-right-radius: 0.5em;
		}
		.su-tab-bottom .su-tab {
			border-bottom: 1px solid lightgrey;
			border-bottom-left-radius: 0.5em;
			border-bottom-right-radius: 0.5em;
		}
		.su-tab-text {
			text-overflow: ellipsis;
			white-space: nowrap;
			overflow: hidden;
			user-select: none;
			text-align: center;
		}
		.su-tab:first-child {
			border-left: 1px solid lightgrey;
		}
		.su-tab:hover {
			background-color: lightblue;
		}
		.su-tab.selected {
			background-color: white;
			font-weight: bold;
			color: blue;
		}
		.su-tab-button {
			font-family: suneido;
			font-style: normal;
			font-weight: normal;
			padding: 5px;
			align-self: center;
			cursor: default;
			user-select: none;
		}
		.su-tab-button:hover {
			outline: 1px solid black;
			outline-offset: -4px;
		}
		.su-tab-extra {
			order: 999;
			flex-grow: 1;
		}'
	tabButton: false
	dropButton: false
	showingDropButton: false
	extraControl: false
	New(.close_button = false, orientation = 'top', tabButton = false,
		extraControl = false, .staticTabs = #())
		{
		LoadCssStyles('su_tabs.css', .styles)
		.tabs = Object()
		.closeImage = close_button isnt false
			? [char: IconFontHelper.GetCode(#close).Chr(),
				font: 'suneido2', color: 'darkgrey']
			: -1

		.CreateElement('div', className: 'su-tab-control')
		.TargetEl = .tabEl = CreateElement('div', .El, 'su-tab-container')
		.tabEl.classList.Add(orientation is 'bottom' ? 'su-tab-bottom' : 'su-tab-top')
		.initHiddenTab()
		if .vertical = orientation in ('left', 'right')
			{
			.SetStyles(#('flex-direction': 'column', 'order': '100'), .tabEl)

			.Ystretch = 1
			.Xstretch = 0
			}
		else
			.Xstretch = 1

		if tabButton isnt false
			{
			.tabButton = CreateElement('div', .tabEl, className: 'su-tab-button')
			.tabButton.SetAttribute('translate', 'no')
			.tabButton.title = tabButton
			.tabButton.textContent = IconFontHelper.GetCode('expand.emf').Chr()
			.tabButton.AddEventListener('click', .onTabButton)
			}

		if extraControl isnt false
			{
			.extraControl = .Construct(extraControl)
			.extraControl.El.classList.Add('su-tab-extra')
			}


		.initSize()
		.initResizeObserver()
		}

	initHiddenTab()
		{
		// for calculate tab width
		.hiddenTab = CreateElement('div', .tabEl, className: 'su-tab selected')
		.hiddenTab.SetStyle('display', 'none')
		}

	initSize()
		{
		metrics = SuRender().GetTextMetrics(.tabEl, 'M')
		.Xmin = .Ymin = Max(metrics.height + 12 /*=padding + border*/,
			.extraControl is false ? 0 : .extraControl.Ymin)
		.SetMinSize()
		.updateWidth()
		}

	initResizeObserver()
		{
		.resizeObserver = SuUI.MakeWebObject('ResizeObserver', .onResize)
		.resizeObserver.Observe(.tabEl)
		}

	onTabButton(event)
		{
		.EventWithOverlay('ButtonClicked', event.target.title,
			Object(x: event.clientX, y: event.clientY))
		}

	onResize(@unused)
		{
		width = .vertical
			? SuRender.GetClientRect(.tabEl).height
			: SuRender.GetClientRect(.tabEl).width
		if .w > width
			{
			if .showingDropButton is true
				return
			if .dropButton is false
				{
				.dropButton = CreateElement('div', className: 'su-tab-button')
				.dropButton.SetAttribute('translate', 'no')
				.dropButton.title = 'Go to Tab'
				.dropButton.innerHTML = '&nbsp;'/*arrow_down in suneido font*/
				.dropButton.AddEventListener('click', .onTabButton)
				}
			.tabEl.AppendChild(.dropButton)
			.showingDropButton = true
			}
		else
			{
			if .showingDropButton is false
				return
			.dropButton.Remove()
			.showingDropButton = false
			}
		}

	OnContextMenu(event)
		{
		i = .tabs.FindIf({ it.textEl is event.target })
		.RunWhenNotFrozen({
			.EventWithFreeze('ContextMenu', event.clientX, event.clientY, i) })
		event.StopPropagation()
		event.PreventDefault()
		}

	Insert(i, text, data, image, id)
		{
		el = CreateElement('div', className: 'su-tab')
		textEl = CreateElement('div', el, className: 'su-tab-text')
		textEl.textContent = text
		if .vertical
			{
			el.SetStyle('writing-mode', 'vertical-lr')
			el.SetStyle('text-orientation', 'upfront')
			}
		item = Object(:text, :data, image: .baseImage(text, image), :id, :el, :textEl)
		if i is .tabs.Size()
			{
			if .tabButton is false
				.tabEl.AppendChild(el)
			else
				.tabEl.InsertBefore(el, .tabButton)
			.tabs.Add(item)
			}
		else
			{
			.tabEl.InsertBefore(el, .tabs[i].el)
			.tabs.Add(item, at: i)
			}

		.addImageEl(item)
		el.AddEventListener('mouseenter', .eventFactory(.onMouseEnter, item))
		el.AddEventListener('mouseleave', .eventFactory(.onMouseLeave, item))
		el.AddEventListener('click', .eventFactory(.click, item, #Click))
		.updateWidth()
		}

	baseImage(text, image)
		{
		baseImage = image isnt -1 or .staticTabs.Has?(text)
			? image
			: .closeImage
		mouseOverImage = .closeImage isnt -1 and not .staticTabs.Has?(text)
			? .closeImage
			: image
		return Object(:baseImage, :mouseOverImage)
		}

	addImageEl(item)
		{
		if -1 is imageOb = item.image.baseImage
			return
		el = CreateElement('div', item.el, at: 0)
		el.SetAttribute('translate', 'no')
		el.textContent = imageOb.char
		.SetStyles(.imageStyle(imageOb), el)
		if item.image.Any?({ it is .closeImage })
			el.AddEventListener(#click, .eventFactory(.click, item, #Tab_Close))
		item.imageEl = el
		}

	imageStyle(imageOb)
		{
		return Object(
			'font-family': imageOb.font,
			'font-style': 'normal',
			'font-weight': 'normal',
			'margin-right': '5px',
			'user-select': 'none',
			'color': ToCssColor(imageOb.GetDefault(#color, #inherit)))
		}

	onMouseEnter(item)
		{
		if false isnt imageEl = item.GetDefault(#imageEl, false)
			.updateImage(imageEl, item.image.mouseOverImage)

		target = item.textEl
		if target.offsetWidth < target.scrollWidth
			target.SetAttribute('title', item.text)
		else
			target.RemoveAttribute('title')
		}

	onMouseLeave(item)
		{
		if false isnt imageEl = item.GetDefault(#imageEl, false)
			.updateImage(imageEl, item.image.baseImage)
		}

	updateImage(imageEl, imageOb)
		{
		.SetStyles(.imageStyle(imageOb), imageEl)
		imageEl.textContent = imageOb.char
		}

	Remove(i)
		{
		if .tabs[i] is .selected
			{
			.selected = false
			}
		.tabs[i].el.Remove()
		.tabs.Delete(i)
		.updateWidth()
		}

	selected: false
	Select(i)
		{
		if .selected isnt false
			.selected.el.className = "su-tab"
		.selected = .tabs[i]
		.selected.el.className = "su-tab selected"
		}

	eventFactory(@args)
		{
		fn = args[0]
		return { |event|
			args.event = event
			fn(@+1args)
			}
		}

	click(item, endPoint, event)
		{
		i = .tabs.Find(item)
		.RunWhenNotFrozen({ .EventWithOverlay(endPoint, i) })
		event.StopPropagation()
		}

	DoTabChange(i, change?)
		{
		if change? is false
			return
		.Select(i)
		}

	SetText(i, text)
		{
		.tabs[i].text = text
		.tabs[i].textEl.textContent = text
		.updateWidth(i)
		}

	SetImage(i, image)
		{
		tab = .tabs[i]
		imageOb = .baseImage(tab.text, image is -1 ? -1 : image)
		if tab.image is imageOb
			return
		if tab.Member?(#imageEl)
			{
			tab.imageEl.Remove()
			tab.Delete(#imageEl)
			}
		tab.image = imageOb
		if imageOb.baseImage isnt -1
			.addImageEl(tab)
		.updateWidth(i)
		}

	updateWidth(tabChanged = false)
		{
		if tabChanged isnt false
			.tabs[tabChanged].Delete(#width)
		w = 1
		for tab in .tabs
			{
			if tab.Member?(#width)
				w += tab.width
			else
				{
				textWidth = Max(
					SuRender().GetTextMetrics(.hiddenTab, tab.textEl.textContent).width,
					SuRender().GetTextMetrics(tab.textEl, tab.textEl.textContent).width)
				tab.textEl.SetStyle(.vertical ? 'height' : 'width', textWidth $ 'px')
				width = textWidth + 10/*=padding left and right*/ + 1/*=border*/
				if tab.Member?(#imageEl)
					width += SuRender().
						GetTextMetrics(tab.imageEl, tab.imageEl.textContent).width +
						5/*=padding right*/
				tab.width = width
				w += width
				}
			}
		if .tabButton isnt false
			w += SuRender().GetTextMetrics(.tabButton, .tabButton.textContent).width +
				10/*=padding*/
		.w = w
		}

	Destroy()
		{
		.resizeObserver.Unobserve(.tabEl)
		.resizeObserver = false
		super.Destroy()
		}
	}
