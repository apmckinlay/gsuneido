// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	webBrowser: false
	New(.parent)
		{
		.webBrowser = WebBrowser(.parent.Hwnd)
		if Number?(.webBrowser)
			throw "WebView: Error Code (" $ .webBrowser $ ')'
		}

	Load(what)
		{
		if what is false
			return
		if what.Prefix?('MSHTML:')
			{
			.parent.SetEnabled(false)
			if '' is content = what.RemovePrefix('MSHTML:')
				content = ' ' // need to write something to make sure the body is created
			doc = .document()
			doc.Write(content)
			doc.Close()
			doc.Release()
			.parent.SetEnabled(true)
			}
		else
			.webBrowser.Navigate(what)
		}

	document()
		{
		if .webBrowser.ReadyState is 0 // Uninitialized
			.webBrowser.Navigate('about:blank') // Navigate to initialize .Document
		return .webBrowser.Document
		}

	Ready?()
		{
		return .webBrowser.ReadyState
		}

	Resize(w, h)
		{
		.webBrowser.Width = w
		.webBrowser.Height = h
		}

	OnNavComplete(block)
		{
		block()
		}

	SetFocus() {}

	SetCssStyle(style)
		{
		comOb = .webBrowser.Document
		styleElement = comOb.CreateElement('style')
		styleElement.Type = "text/css"
		styleElement.AppendChild(comOb.CreateTextNode(style))
		comOb.Body.AppendChild(styleElement)
		}

	Getter_LocationURL()
		{
		return .webBrowser.LocationURL
		}

	TriggerKeyDown(key)
		{
		hwnd = .getChildWindow()
		SendMessage(hwnd, WM.CHAR, key, 0)
		}

	getChildWindow()
		{
		i = 0
		hwnd = .parent.Hwnd
		do
			{
			prev = hwnd
			hwnd = GetWindow(hwnd, GW.CHILD)
			}
			while hwnd isnt 0 and ++i < 10 /*= max hwnd find attempts */
		return prev
		}

	DoFind()
		{
		.findDialog(.getChildWindow())
		}

	findDialog(hwnd)
		{
		idm_find = 67
		from_accel = 1
		SendMessage(hwnd, WM.COMMAND, MAKELONG(idm_find, from_accel), 0)
		}

	DoGoBack()
		{
		.webBrowser.GoBack()
		}

	DoGoForward()
		{
		.webBrowser.GoForward()
		}

	DoCopy()
		{
		try
			.webBrowser.ExecWB(OLECMDID.COPY, OLECMDEXECOPT.DODEFAULT, NULL, NULL)
		catch(err /*unused*/, "COMobject call failed ExecWB 80020009")
			{
			// NON FATAL ERROR from COM - caused by holding Ctl+C down
			}
		}

	DoPaste()
		{
		.webBrowser.ExecWB(OLECMDID.PASTE, OLECMDEXECOPT.DODEFAULT, NULL, NULL)
		}

	DoRefresh()
		{
		.webBrowser.ExecWB(OLECMDID.PREREFRESH, OLECMDEXECOPT.DODEFAULT, NULL, NULL)
		.webBrowser.ExecWB(OLECMDID.REFRESH, OLECMDEXECOPT.DODEFAULT, NULL, NULL)
		}

	DoPrint()
		{
		.webBrowser.ExecWB(OLECMDID.PRINT, OLECMDEXECOPT.DODEFAULT, NULL, NULL)
		}

	DoPrintPreview()
		{
		.webBrowser.ExecWB(OLECMDID.PRINTPREVIEW, OLECMDEXECOPT.DODEFAULT, NULL, NULL)
		}

	DoPageSetup()
		{
		.webBrowser.ExecWB(OLECMDID.PAGESETUP, OLECMDEXECOPT.DODEFAULT, NULL, NULL)
		}

	InsertAdjacentHTML(id, position, text)
		{
		doc = .document()
		parent = doc.GetElementById(id)
		parent.InsertAdjacentHTML(position, text)
		doc.Release()
		parent.Release()
		}

	ScrollIntoView(id, alignToTop)
		{
		doc = .document()
		el = doc.GetElementById(id)
		el.ScrollIntoView(alignToTop)
		doc.Release()
		el.Release()
		}

	Release()
		{
		if .webBrowser isnt false
			.webBrowser.Release()
		}
	}
