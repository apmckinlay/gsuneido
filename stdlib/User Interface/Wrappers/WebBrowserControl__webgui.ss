// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: "WebBrowser"
	Xstretch: 1
	Ystretch: 1
	ComponentName: "WebBrowser"
	ComponentArgs: #()
	comObject: class
		{
		New(.browser) { }
		Navigate(url)
			{
			fragment = url.AfterLast('#')
			url = url.RemoveSuffix('#' $ fragment)
			.LocationURL = url
			if url is "about:blank"
				return
			if url.Prefix?('http://') or url.Prefix?('https://')
				{
				.browser.Act('Navigate', url)
				return
				}
			res = JsSuneidoAPP(url)
			if res isnt ''
				{
				.Load(res)
				if fragment isnt ''
					.browser.Act('Locate', fragment)
				}
			}
		Load(res)
			{
			.browser.Act('Load', res)
			}
		Print()
			{
			.browser.Act('Print')
			}
		Default(@args)
			{
			SuServerPrint("COM", Display(args))
			}
		LocationURL: ''
		}

	New(what)
		{
		.webview = (.comObject)(this)
		.Load(what)
		}

	Load(what)
		{
		if what is false
			return
		if what.Prefix?('MSHTML:')
			.webview.Load(JsSuneidoAPP.Convert(what.RemovePrefix('MSHTML:')))
		else
			.webview.Navigate(what)
		}

	SetCssStyle(style)
		{
		.InsertAdjacentHTML(false, 'beforeend', '<style>' $
			JsSuneidoAPP.Convert(style.Tr('\r\n')) $ '</style>')
		}

	Getter_LocationURL()
		{
		return .webview.LocationURL
		}

	LinkGoTo(url, query)
		{
		.webview.Navigate('suneido:' $
			XmlEntityDecode(Base64.Decode(url)) $ Url.BuildQuery(query))
		}

	DoCopy()
		{
		.Act('Copy')
		}

	DoPrint()
		{
		.webview.Print()
		}

	TriggerKeyDown(key)
		{
		.Act('TriggerKeyDown', key)
		}

	// if id is false, insert to body
	InsertAdjacentHTML(id, position, text)
		{
		.Act('InsertAdjacentHTML', id, position, text)
		}

	ScrollIntoView(id, alignToTop)
		{
		.Act('ScrollIntoView', id, alignToTop)
		}

	OnNavComplete(block)
		{
		block()
		}
	}
