// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WebBrowserControl
	{
	Name:		'Browser'
	Xstretch:	1
	Ystretch:	1
	hidden: 	false

	subHwnd: false
	New(url = false, .asField = false)
		{
		super(.setup(url))
		.subs = Object()
		if asField
			.Send("Data")
		else
			.subs.Add(
				PubSub.Subscribe(
					'BrowserRedir_' $ .subHwnd = .WindowHwnd(),
					{|@args| this[args[0]](@+1args) }))
		}
	setup(url)
		{
		.history = Object()
		.hi = -1
		.hn = 0 // .history[.hi] is always current page
		return url
		}
	Goto(url, checkPageFrozen? = true)
		{
		// don't allow navigation if page is frozen
		if checkPageFrozen? and .Send("PageFrozen?") is true
			return

		try
			.Load(Url.Decode(url))
		catch (err, "COMobject call failed Navigate 80020009")
			{
			// Users open a help page's Properties dialog which will disable the Help Window,
			// then clicking the Help button will re-enable the Window (OpenBook).
			// Now the Help Window is free to navigate but the Preperties dialog prevents it
			if url.Has?('Help')
				SuneidoLog('Warning: (CAUGHT) ' $ err, params: Object(:url),
					caughtMsg: 'help will focus without navigating.')
			else
				throw err
			}
		}

	Load(what)
		{
		if what isnt false and .subHwnd isnt false
			{
			browserLoads = Suneido.GetInit('BrowserRedir_Loads', Object)
			browserLoads[.subHwnd] = browserLoads.GetDefault(.subHwnd, Object())
			browserLoads[.subHwnd][what] = Date()
			}
		super.Load(what)
		}

	Current()
		{ return Url.Decode(.LocationURL) }
	GoBack(i = 1)
		{
		// don't allow navigation if page is frozen
		if .Send("PageFrozen?") is true
			return
		if .Current().Prefix?("suneido:")
			{
			if (.hi is 0)
				{ Beep(); return }
			.Goto(.history[.hi = .hi - i], checkPageFrozen?: false)
			}
		else
			try .DoGoBack()
		}
	GoForward(i = 1)
		{
		// don't allow navigation if page is frozen
		if .Send("PageFrozen?") is true
			return
		if .Current().Prefix?("suneido:")
			{
			if .hi >= .hn
				{ Beep(); return }
			.Goto(.history[.hi += i], checkPageFrozen?: false)
			}
		else
			try .DoGoForward()
		}
	// returns true/false
	// this is for SuneidoAPP -
	// if this return false, then SuneidoAPP will still run .app after calling "Going"
	// if this returns true, then SuneidoApp will NOT call .app, and will return
	// after calling "Going"
	Going(url)
		{
		if .hidden
			{
			.hidden = false
			.Parent.Close()
			.OnNavComplete()
				{
				if .hidden is false
					// show the browser
					.SetVisible(true)
				}
			}
		if (not url.Prefix?("suneido:"))
			return 0
		x = 0 isnt .Send("Going", url)
		if .history.Member?(.hi) and url is .history[.hi]
			return x
		.history[++.hi] = url
		.hn = .hi
		return x
		}
	BackList()
		{ return .history[.. .hi].Reverse!() }
	ForwardList()
		{ return .history[.hi + 1 :: .hn - .hi] }
	CanBack?()
		{ return true }
	CanForward?()
		{ return true }

	Print()
		{
		if .hidden
			Beep()
		else if Sys.SuneidoJs?()
			.DoPrint()
		else
			{
			if not .LocationURL.Suffix?('?print')
				.Load(.LocationURL $ '?print')
			.DoPrint()
			}
		}
	PrintPreview()
		{
		if .hidden
			Beep()
		else
			{
			if not .LocationURL.Suffix?('?print')
				.Load(.LocationURL $ '?print')
			.DoPrintPreview()
			}
		}
	PageSetup()
		{
		if .hidden
			Beep()
		else
			.DoPageSetup()
		}

	Set(control) // program page
		{
		if .asField and String?(control) and control =~ 'https?://'
			{
			.Load(control)
			return
			}
		// hide the browser
		.SetVisible(false)
		.hidden = true
		// create the program page
		.Send("ProgramPage")
		.Parent.Open(control)
		}
	On_Copy()
		{
		if .hidden
			Beep()
		else
			.DoCopy()
		}
	Paste()
		{
		.DoPaste()
		}
	On_Find()
		{
		.DoFind()
		}
	Refresh()
		{
		.DoRefresh()
		}
	SetEnabled(unused)
		{
		}
	Destroy()
		{
		if .asField
			.Send("NoData")
		.subs.Each('Unsubscribe')
		if .subHwnd isnt false
			Suneido.GetDefault('BrowserRedir_Loads', Object()).Delete(.subHwnd)
		super.Destroy()
		}
	}
