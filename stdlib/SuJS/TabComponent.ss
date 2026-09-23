// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
/* Element hierarchy
.su-tab-control              (root)
	.su-tab-container        (.tabEl) — flex row; .su-tab-top/bottom/left/right/scroll
	|-	.su-tab-viewport     (.viewportEl) — scroll window
		|-	.su-tab-row      (.viewportRowEl) — absolute flex strip that actually scrolls
			|-	.hiddenTab   (display:none; width probe)
			|-	.su-tab × N
				|-	imageEl  (optional close/custom icon)
				|-	.su-tab-text  (.textEl) — clipped label
			|-	.su-tab-button   (expand/tab-list dropdown)
	|-	.su-tab-extra        (optional user-supplied control at the end)
*/
Component
	{
	Name: #Tab
	ContextMenu: true
	styles:       "
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
			box-sizing: border-box;
		}
		.su-tab-viewport {
			position: relative;
			overflow: hidden;
			flex-grow: 1;
		}
		.su-tab-row {
			position: absolute;
			top: 0px;
			display: flex;
			width: 100%;
			height: 100%;
		}
		.su-tab-container.su-tab-top {
			border-bottom: 1px solid lightgrey;
		}
		.su-tab-container.su-tab-bottom {
			border-top: 1px solid lightgrey;
		}
		.su-tab-container.su-tab-left {
			border-right: 1px solid lightgrey;
		}
		.su-tab-container.su-tab-right {
			border-left: 1px solid lightgrey;
		}
		.su-tab {
			border-right: 1px solid lightgrey;
			outline: none;
			cursor: pointer;
			padding: 5px;
			transition: background-color 0.3s, font-weight 0.3s, color 0.3s;;
			display: flex;
			align-items: baseline;
			overflow: hidden;
			font-weight: normal;
		}
		.su-tab-scroll .su-tab {
			overflow: visible;
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
		.su-tab-right .su-tab {
			border-right: 1px solid lightgrey;
			border-top-right-radius: 0.5em;
			border-bottom-right-radius: 0.5em;
		}
		.su-tab-left .su-tab {
			border-left: 1px solid lightgrey;
			border-top-left-radius: 0.5em;
			border-bottom-left-radius: 0.5em;
		}
		.su-tab-top .su-tab:last-child
		.su-tab-bottom .su-tab:last-child {
			border-right: 1px solid lightgrey;
		}
		.su-tab-right .su-tab:last-child,
		.su-tab-left .su-tab:last-child {
			border-bottom: 1px solid lightgrey;
		}
		.su-tab-top .su-tab:first-child
		.su-tab-bottom .su-tab:first-child {
			border-left: 1px solid lightgrey;
		}
		.su-tab-left .su-tab:first-child
		.su-tab-right .su-tab:first-child {
			border-top: 1px solid lightgrey;
		}
		.su-tab-text {
			text-overflow: ellipsis;
			white-space: nowrap;
			overflow: hidden;
			user-select: none;
			text-align: center;
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
		}"
	tabButton:    false
	extraControl: false
	New(.close_button = false, orientation = #top, tabButton = false,
		extraControl = false, .staticTabs = #(), .scrollTabs = false)
		{
		LoadCssStyles("su_tabs.css", .styles)
		.tabs = Object()
		.buttons = Object()
		.closeImage = close_button isnt false
			? [char: IconFontHelper.GetCode(#close).Chr(),
				font: #suneido2, color: #darkgrey]
			: -1

		.CreateElement(#div, className: "su-tab-control")
		.TargetEl = .tabEl = CreateElement(#div, .El, "su-tab-container")

		if .scrollTabs is true
			.tabEl.classList.Add("su-tab-scroll")

		.viewportEl = CreateElement(#div, .tabEl, "su-tab-viewport")
		.viewportRowEl = CreateElement(#div, .viewportEl, "su-tab-row")
		.initHiddenTab()
		if orientation in (#left, #right)
			{
			.vertical = true
			.tabEl.classList.Add(orientation is #right ? "su-tab-right" : "su-tab-left")
			.SetStyles(#("flex-direction": column, order: "100"), .viewportRowEl)

			.Ystretch = 1
			.Xstretch = 0
			}
		else
			{
			.vertical = false
			.tabEl.classList.Add(orientation is #bottom ? "su-tab-bottom" : "su-tab-top")
			.Xstretch = 1
			}

		.initTabButton(tabButton)

		if extraControl isnt false
			{
			.extraControl = .Construct(extraControl)
			.extraControl.El.classList.Add("su-tab-extra")
			}

		.initSize()
		.initResizeObserver()
		}

	initHiddenTab()
		{
		// for calculate tab width
		.hiddenTab = CreateElement(#div, .viewportRowEl, className: "su-tab selected")
		.hiddenTab.SetStyle(#display, #none)
		}

	initSize()
		{
		metrics = SuRender().GetTextMetrics(.tabEl, 'M')
		.Xmin = .Ymin = Max(metrics.height + 12/*=padding + border*/,
			.extraControl is false ? 0 : .extraControl.Ymin)
		.SetMinSize()
		.updateWidth()
		}

	initResizeObserver()
		{
		.resizeObserver = SuUI.MakeWebObject(#ResizeObserver, .onResize)
		.resizeObserver.Observe(.tabEl)
		}

	initTabButton(tabButton)
		{
		if tabButton isnt false
			{
			.tabButton = CreateElement(#div, .viewportRowEl, className: "su-tab-button")
			.tabButton.SetAttribute(#translate, #no)
			.tabButton.title = tabButton
			.tabButton.textContent = IconFontHelper.GetCode("expand.emf").Chr()
			.tabButton.AddEventListener(#click, .onTabButton)
			}
		}

	onTabButton(event)
		{
		.EventWithOverlay(#ButtonClicked, event.target.title,
			Object(x: event.clientX, y: event.clientY))
		}

	onScrollPrevious()
		{
		if false is leftTab = .findLeftTab()
			return

		if leftTab is 0
			return

		prev = .tabs[leftTab-1]
		.viewportRowEl.SetStyle(#left, -prev.el.offsetLeft $ "px")
		}

	onScrollNext()
		{
		if false is leftTab = .findLeftTab()
			return

		width = .vertical
			? SuRender.GetClientRect(.viewportEl).height
			: SuRender.GetClientRect(.viewportEl).width

		if .w + .viewportRowEl.offsetLeft <= width
			return // already displaying all tabs from cur to the end

		if false is next = .tabs.GetDefault(leftTab + 1, false)
			return

		.viewportRowEl.SetStyle(#left, -next.el.offsetLeft $ "px")
		}

	restoreScroll()
		{
		if .scrollTabs is true
			.viewportRowEl.SetStyle(#left, "0px")
		}

	findLeftTab()
		{
		offset = .viewportRowEl.offsetLeft
		for (i = 0; i < .tabs.Size(); i += 1)
			if offset + .tabs[i].el.offsetLeft >= 0
				return i
		return false
		}

	onResize(@unused)
		{
		width = .vertical
			? SuRender.GetClientRect(.tabEl).height
			: SuRender.GetClientRect(.tabEl).width

		if .w > width
			.showDropScrollButton()
		else
			.hideDropScrollButton()
		.restoreScroll()
		}

	showDropScrollButton()
		{
		if .buttons.NotEmpty?()
			return

		if .scrollTabs is false
			{
			dropButton = .createButton("Go to Tab",
				"&nbsp;" /*arrow_down in suneido font*/,
				.onTabButton)
			.tabEl.AppendChild(dropButton)
			.buttons.Add(dropButton)
			}
		else
			{
			prevButton = .createButton(#Previous, '.', .onScrollPrevious)
			.tabEl.AppendChild(prevButton)
			.buttons.Add(prevButton)
			nextButton = .createButton(#Next, '/', .onScrollNext)
			.tabEl.AppendChild(nextButton)
			.buttons.Add(nextButton)
			}
		}

	hideDropScrollButton()
		{
		if .buttons.Empty?()
			return

		.buttons.Each({ it.Remove() })
		.buttons.Delete(all:)
		}

	createButton(title, icon, cb)
		{
		button = CreateElement(#div, className: "su-tab-button")
		button.SetAttribute(#translate, #no)
		button.title = title
		button.innerHTML = icon
		button.AddEventListener(#click, cb)
		return button
		}

	OnContextMenu(event)
		{
		i = .tabs.FindIf({ it.textEl is event.target })
		.RunWhenNotFrozen(
			{
			.EventWithFreeze(#ContextMenu, event.clientX, event.clientY, i)
			})
		event.StopPropagation()
		event.PreventDefault()
		}

	Insert(i, text, data, image, id)
		{
		el = CreateElement(#div, className: "su-tab")
		textEl = CreateElement(#div, el, className: "su-tab-text")
		textEl.textContent = text
		if .vertical
			{
			el.SetStyle("writing-mode", "vertical-lr")
			el.SetStyle("text-orientation", #upfront)
			}
		item = Object(:text, :data, image: .baseImage(text, image), :id, :el, :textEl)
		if i is .tabs.Size()
			{
			if .tabButton is false
				.viewportRowEl.AppendChild(el)
			else
				.viewportRowEl.InsertBefore(el, .tabButton)
			.tabs.Add(item)
			}
		else
			{
			.viewportRowEl.InsertBefore(el, .tabs[i].el)
			.tabs.Add(item, at: i)
			}

		.addImageEl(item)
		el.AddEventListener(#mouseenter, .eventFactory(.onMouseEnter, item))
		el.AddEventListener(#mouseleave, .eventFactory(.onMouseLeave, item))
		el.AddEventListener(#click, .eventFactory(.click, item, #Click))
		.updateWidth()
		}

	baseImage(text, image)
		{
		baseImage = image isnt -1 or .staticTabs.Has?(text) ? image : .closeImage
		mouseOverImage = .closeImage isnt -1 and not .staticTabs.Has?(text)
			? .closeImage
			: image
		return Object(:baseImage, :mouseOverImage)
		}

	addImageEl(item)
		{
		if -1 is imageOb = item.image.baseImage
			return
		el = CreateElement(#div, item.el, at: 0)
		el.SetAttribute(#translate, #no)
		el.textContent = imageOb.char
		.SetStyles(.imageStyle(imageOb), el)
		if item.image.Any?({ it is .closeImage })
			el.AddEventListener(#click, .eventFactory(.click, item, #Tab_Close))
		item.imageEl = el
		}

	imageStyle(imageOb)
		{
		return Object(
			"font-family": imageOb.font,
			"font-style": #normal,
			"font-weight": #normal,
			"margin-right": "5px",
			"user-select": #none,
			color: ToCssColor(imageOb.GetDefault(#color, #inherit)))
		}

	onMouseEnter(item)
		{
		if false isnt imageEl = item.GetDefault(#imageEl, false)
			.updateImage(imageEl, item.image.mouseOverImage)

		target = item.textEl
		if "" isnt tooltip = item.data.GetDefault(#tooltip, "")
			target.SetAttribute(#title, tooltip)
		else if target.offsetWidth < target.scrollWidth
			target.SetAttribute(#title, item.text)
		else
			target.RemoveAttribute(#title)
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
			.selected = false
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
		return {|event|
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

	w: -1
	updateWidth(tabChanged = false)
		{
		if tabChanged isnt false
			.tabs[tabChanged].Delete(#width)
		w = 1
		for tab in .tabs
			if tab.Member?(#width)
				w += tab.width
			else
				{
				textWidth = Max(
					SuRender().GetTextMetrics(.hiddenTab, tab.textEl.textContent).width,
					SuRender().GetTextMetrics(tab.textEl, tab.textEl.textContent).width)
				tab.textEl.SetStyle(.vertical ? #height : #width, textWidth $ "px")
				width = textWidth + 10/*=padding left and right*/ + 1/*=border*/
				if tab.Member?(#imageEl)
					width += SuRender().
							GetTextMetrics(tab.imageEl, tab.imageEl.textContent).width +
						5/*=padding right*/
				tab.width = width
				w += width
				}
		if .tabButton isnt false
			w += SuRender().GetTextMetrics(.tabButton, .tabButton.textContent).width +
				10/*=padding*/
		if .w isnt w
			{
			.w = w
			.onResize()
			}
		}

	Destroy()
		{
		.resizeObserver.Unobserve(.tabEl)
		.resizeObserver = false
		super.Destroy()
		}
	}
