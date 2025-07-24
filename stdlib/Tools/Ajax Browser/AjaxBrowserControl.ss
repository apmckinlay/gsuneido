// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'the page'
	New()
		{
		super(.layout())
		.browser = .FindControl('WebBrowser')
		if .Window.Member?('Ctrl')
			{
			.curAccels = .Window.SetupAccels(.Commands)
			// make accelerators work inside book
			.Window.Ctrl.Redir('On_Refresh', this)
			.Window.Ctrl.Redir('On_Home', this)
			.Window.Ctrl.Redir('On_End', this)
			.Window.Ctrl.Redir('On_Up', this)
			.Window.Ctrl.Redir('On_Down', this)
			.Window.Ctrl.Redir('On_Left', this)
			.Window.Ctrl.Redir('On_Right', this)
			.Window.Ctrl.Redir('On_PgDn', this)
			.Window.Ctrl.Redir('On_PgUp', this)
			.Window.Ctrl.Redir('On_Escape', this)
			.Window.Ctrl.Redir('On_Copy', this)
			}
		}
	layout()
		{
		if true isnt result = .checkForHttp()
			return result
		return Object('WebBrowser', .UrlAddress())
		}
	checkForHttp()
		{
		if Sys.Client?() and true isnt TestHttpServer()
			return Object('Center'
				Object('Vert'
					Object('Static' 'Unable to access ' $ .Title, justify: 'CENTER',
						xstretch: 1, textStyle: 'main')
					'Skip'
					Object('Static' "Cannot connect to HTTP service (server IP: " $
						ServerIP() $ ", port: " $ Display(HttpPort()) $ ")",
						justify: 'CENTER', xstretch: 1)
					)
				)
		// sending requests during the book check causes transaction check failures
		if BookCheck.CheckingBook?()
			return Object('Static', 'Checking Book')
		return true
		}
	uuid: false
	UrlAddress()
		{
		.uuid = ServerEval('WebSession.Register', Suneido.User)
		return .Server() $ .Endpoint $ '?sessionid=' $ .uuid
		}
	UrlPage(@args /*unused*/)
		{
		return ''
		}
	Server()
		{
		if false isnt host = Suneido.GetDefault(#JsConnectionHost, false)
			return (Suneido.GetDefault(#isHttps?, false) ? 'https://' : 'http://') $
				host $ '/'

		if '' is serverip = ServerIP()
			serverip = '127.0.0.1'
		return 'http://' $ serverip $ ':' $ HttpPort() $ '/'
		}
	Commands: (
		(Refresh,"F5")
		(Home,	"Home"),
		(End,	"End"),
		(Up,	"Up"),
		(Down,	"Down"),
		(Right,	"Right"),
		(Left,	"Left"),
		(PgUp,	"PRIOR"),
		(PgDn,	"NEXT"),
		(Escape,"Escape"),
		(Copy, "Ctrl+C"))

	On_Refresh()
		{
		if .uuid isnt false
			.CloseSession(.uuid)
		if .browser isnt false
			.browser.Load(.UrlAddress())
		}
	sendCommand(key)
		{
		if .browser isnt false
			.browser.TriggerKeyDown(key)

		}
	On_Home()
		{
		.sendCommand(36) /*= home */
		}
	On_End()
		{
		.sendCommand(35) /*= end */
		}
	On_Up()
		{
		.sendCommand(38) /*= up */
		}
	On_Down()
		{
		.sendCommand(40) /*= down */
		}
	On_Left()
		{
		.sendCommand(37) /*= left */
		}
	On_Right()
		{
		.sendCommand(39) /*= right */
		}
	On_PgUp()
		{
		.sendCommand(33) /*= page up */
		}
	On_PgDn()
		{
		.sendCommand(34) /*= page down */
		}
	On_Escape()
		{
		.sendCommand(27) /*= escape */
		}
	On_Copy()
		{
		if .browser isnt false
			.browser.DoCopy()
		}
	Navigate(url)
		{
		if .browser isnt false
			.browser.Load(url)
		}
	curAccels: false
	CloseSession(sessionid)
		{
		ServerEval('WebSession.Close', sessionid)
		return true
		}
	GetSessionId()
		{
		return .uuid
		}
	Destroy()
		{
		if .uuid isnt false
			.CloseSession(.uuid)

		if .curAccels isnt false
			{
			.Window.Ctrl.RemoveRedir(this)
			.Window.RestoreAccels(.curAccels)
			.Window.Ctrl.Redir('On_Copy')
			}
		if .browser isnt false
			.browser.Destroy()
		super.Destroy()
		}
	}
