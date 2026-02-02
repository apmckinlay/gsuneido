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
		if not WebView2.Available?()
			{
			sid = LoginSessionAddress()
			.alert(terminalUser?: sid.Prefix?('wts'))
			ExitClient(true)
			}
		else
			.webview = WebView2(this)

		.Load(what)
		}

	alert(terminalUser? = false)
		{
		msg = "Microsoft WebView2 is not installed on this computer."
		msg $= "\n\nPlease " $
			(terminalUser? ? "contact your administrator to " : "") $
				"install WebView2 using the link below:" $
				"\n\nhttps://go.microsoft.com/fwlink/p/?LinkId=2124703"
		Alert(msg, title: "Microsoft WebView2 Required")
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
