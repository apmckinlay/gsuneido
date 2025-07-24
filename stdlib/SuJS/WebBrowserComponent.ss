// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: 		'WebBrowser'
	Xstretch:	1
	Ystretch:	1

	contentReady?: false
	New()
		{
		.CreateElement('iframe')
		.El.SetStyle('background-color', 'white')
		SuUI.GetCurrentWindow().AddEventListener('online', .online)
		.pendingInsert = Object()
		.El.addEventListener(#load,
			{ |unused|
			.contentReady? = true
			if (.El.src is '' or .sameOrigin?(.El))
				{
				contentWindow = SuUI.GetContentWindow(.El)
				contentWindow.suIframeSend = .goto
				SuRender.NotAllowDragAndDrop(contentWindow)
				.El.contentDocument.AddEventListener(#mousedown, .mousedown)
				.El.contentDocument.AddEventListener(#scroll, .scroll)
				if .fragment isnt ''
					{
					try
						{
						target = .El.contentDocument.QuerySelector(
							'[name="' $ .fragment $ '"]')
						.fragment = ''
						target.ScrollIntoView()
						}
					catch (e)
						Print(e)
					}
				if .pendingInsert.NotEmpty?()
					{
					for insert in .pendingInsert
						.InsertAdjacentHTML(@insert)
					.pendingInsert = Object()
					}
				}
			})
		SuRender().RegisterIframe(.El)
		.resizeObserver = SuUI.MakeWebObject('ResizeObserver', .onResize)
		.resizeObserver.Observe(.El)
		}

	online()
		{
		if Boolean?(.El) or '' is src = .El.GetDefault('src','')
			return
		.El.src = ''
		.El.src = src
		}

	sameOrigin?(iframe)
		{
		try
			{
			contentWindow = SuUI.GetContentWindow(iframe)
			contentWindow['location']['hostname']
			return true
			}
		return false
		}

	Load(res)
		{
		.pendingScroll = .lastScrollY = false
		.contentReady? = false
		.pendingInsert = Object()
		.El.srcdoc = res
		}

	Navigate(url)
		{
		.pendingScroll = .lastScrollY = false
		.contentReady? = false
		.pendingInsert = Object()
		.El.src = url
		}

	fragment: ''
	Locate(.fragment) {}

	goto(@args)
		{
		url = args[0]
		args = args[1..]
		query = Object()
		for (i = 0; i < args.Size(); i++)
			{
			query[args[i]] = args[++i]
			}
		.EventWithOverlay(#LinkGoTo, url, query)
		}

	mousedown(event)
		{
		.Window.DoActivate(event, target: .El)
		}

	lastScrollY: false
	scroll(event/*unused*/)
		{
		if .isHidden?()
			return
		.lastScrollY = SuUI.GetContentWindow(.El).scrollY
		}

	hidden?: false
	onResize(@unused)
		{
		hidden? = .isHidden?()
		if .hidden? is true and hidden? is false
			{
			if .pendingScroll isnt false
				{
				.ScrollIntoView(@.pendingScroll)
				.pendingScroll = false
				}
			else if .lastScrollY isnt false
				SuUI.GetContentWindow(.El).ScrollTo(Object(top: .lastScrollY))
			}

		.hidden? = hidden?
		}

	isHidden?()
		{
		return .El.offsetHeight is 0 and .El.offsetWidth is 0
		}

	Copy()
		{
		SuUI.GetCurrentDocument().ExecCommand('copy')
		}

	Print()
		{
		SuUI.GetContentWindow(.El).Print()
		}

	InsertAdjacentHTML(id, position, text)
		{
		if .contentReady? is true
			.El.contentDocument.GetElementById(id).InsertAdjacentHTML(position, text)
		else
			.pendingInsert.Add(Object(id, position, text))
		}

	pendingScroll: false
	ScrollIntoView(id, alignToTop)
		{
		if .isHidden?()
			.pendingScroll = [id, alignToTop]
		else
			try .El.contentDocument.GetElementById(id).ScrollIntoView(alignToTop)
		}

	TriggerKeyDown(key)
		{
		SuUI.GetContentWindow(.El).Eval("document.dispatchEvent(
			new KeyboardEvent('keydown', { keyCode: " $ key $ " }))")
		}

	Destroy()
		{
		SuRender().UnregisterIframe(.El)
		.resizeObserver.Unobserve(.El)
		.resizeObserver = false
		super.Destroy()
		}
	}
