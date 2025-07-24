// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Hwnd
	{
	Name: "WebBrowser"
	Xstretch: 1
	Ystretch: 1
	webview: false

	New(what)
		{
		.CreateWindow('SuWhiteArrow', '', WS.VISIBLE)
		try
			{
			if not WebView2.Available?()
				{
				sid = LoginSessionAddress()
				if not sid.Prefix?('wts')
					{
					Alert("Microsoft WebView2 is not installed " $
						"on your system.\n\nPlease install it using the link below" $
						"\n\n\thttps://go.microsoft.com/fwlink/p/?LinkId=2124703",
						"WebView2 Required")
					ExitClient(true)
					}
				else
					.webview = WebView(this)
				}
			else
				.webview = WebView2(this)
			}
		catch (e, 'WebView:')
			{
			SuneidoLog("ERROR: (CAUGHT) Unable to create WebBrowser (" $
				e.RemovePrefix('WebView: ') $ ')', params: [:what], caughtMsg: 'fallback')
			.webview = WebView(this)
			}
		.Load(what)
		}

	WebView2?()
		{
		return .webview.Base?(WebView2)
		}

	Default(@args)
		{
		(.webview[args[0]])(@+1args)
		}

	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		.webview.Resize(w, h)
		UpdateWindow(.Hwnd)
		}

	SetFocus()
		{
		super.SetFocus()
		.webview.SetFocus()
		}

	SetVisible(visible)
		{
		super.SetVisible(visible)
		if visible
			UpdateWindow(.Hwnd)
		}

	Getter_LocationURL()
		{
		return .webview.LocationURL
		}

	InsertAdjacentHTML(id, position, text)
		{
		.webview.InsertAdjacentHTML(id, position, text)
		}

	ScrollIntoView(id, alignToTop)
		{
		.webview.ScrollIntoView(id, alignToTop)
		}

	Destroy()
		{
		if .webview isnt false
			.webview.Release()
		super.Destroy()
		}
	}
