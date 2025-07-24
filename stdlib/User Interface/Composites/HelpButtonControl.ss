// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: HelpButton
	New(imageColor = '', tip = 'Open Help (F1)', .noAccels? = false,
		size = "")
		{
		super(Object('EnhancedButton', command: 'Help', image: 'questionMark_black.emf',
			imageColor: imageColor is '' ? CLR.Highlight : imageColor,
			mouseOverImageColor: CLR.Highlight
			mouseEffect:, imagePadding: .1, :tip, :size))
		.set_help_page()
		if not .noAccels?
			.Defer(.setup)
		}
	set_help_page()
		{
		if 0 is page = .Send('HelpButton_HelpPage')
			page = .Send('AccessGoTo_CurrentBookOption')
		.help_page = String?(page) ? page : false
		}
	setup()
		{
		// even though we are using .Delayed, still need check for Ctrl since there
		// are some cases where Delayed call is being called too early. We think this
		// may be caused by things like dll calls during the construct process which
		// may be doing things with the windows message queue (PC*Miler dlls for example)
		if .Member?(#Window) and .Window.Member?('Ctrl') and not .noAccels?
			{
			.curAccels = .Window.SetupAccels(#(('Help', 'F1')))
			.Window.Ctrl.Redir('On_Help', this)
			if .help_page is false and .Window.Member?(#HelpPage)
				.help_page = .Window.HelpPage
			}
		}
	On_Help()
		{
		path = .help_page is false
			? Suneido.GetDefault(#CurrentBookOption, "")
			: .help_page
		helpBook = Suneido.CurrentBook $ 'Help'
		if path.Suffix?("&goto") or path.Has?("#")
			OpenBook(helpBook, :path)
		else
			{
			pageinfo = Object(path: path.BeforeLast('/'), name: path.AfterLast('/'))
			OpenBook(helpBook, pageinfo)
			}
		}
	curAccels: false
	Destroy()
		{
		if .curAccels isnt false
			{
			.Window.Ctrl.RemoveRedir(this)
			.Window.RestoreAccels(.curAccels)
			}
		super.Destroy()
		}
	}
