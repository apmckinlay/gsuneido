// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'New Messages'

	CallClass()
		{
		BookLog('View New Message')
		if false isnt hwnd = Suneido.GetDefault(#ShowNewEntries, false)
			DestroyWindow(hwnd.hwnd)
		w = Window(Object(this),
			style: WS.SYSMENU | WS.CAPTION | WS.SIZEBOX | WS.TILEDWINDOW,
			exStyle: WS_EX.DLGMODALFRAME,
			keep_placement:)
		Suneido.ShowNewEntries.hwnd = w.Hwnd
		return // no return value
		}

	from: 0
	max_read: 0
	Max_msg: 50
	curAccels: false
	New()
		{
		super(.layout())
		.browser = .FindControl('Browser')
		.prev = .FindControl('Previous')
		.next = .FindControl('Next')
		.size = BookNotification.Count()
		.setButtons()
		.curAccels = .Window.SetupAccels(#(("Copy", "Ctrl+C"), ("Escape", "Escape")))
		.Redir('On_Escape', this)
		}

	layout()
		{
		return Object('Vert',
			Object('Scroll' Object('Browser', .readRecords(), name: 'Browser'))
			#Skip
			Object('HorzEqual',
				'Fill'
				Object('Button', 'Previous', xmin: 90,
					tip: 'Previous ' $ .Max_msg $ ' messages')
				#Skip
				Object('Button', 'Next', xmin: 90,
					tip: 'Next ' $ .Max_msg $ ' messages')
				#Skip)
			#Skip
			)
		}

	readRecords()
		{
		currentMax = .from + .Max_msg
		if currentMax > .max_read
			.max_read = currentMax

		if Suneido.Member?('ShowNewEntries')
			Suneido.ShowNewEntries.status = false
		else
			Suneido.ShowNewEntries = Object(status:, time: Timestamp())
		return 'suneido:/from?NotificationsHtml(' $	.from $ "," $ (currentMax) $ ')'
		}

	setButtons(forceDisable? = false)
		{
		end = Min(.size, .from + .Max_msg)

		.next.SetEnabled(forceDisable? ? false : .loadOkay?() and end isnt .size)
		.prev.SetEnabled(forceDisable? ? false : .loadOkay?() and .from > 0)
		}

	On_Previous()
		{
		if .from > 0
			{
			.from -= .Max_msg
			.refresh()
			}
		}

	On_Escape()
		{
		.Window.Destroy()
		}

	On_Next()
		{
		.from = Min(.size, .from + .Max_msg)
		.refresh()
		}

	refresh()
		{
		if not .loadOkay?()
			{
			.setButtons(true)
			return
			}

		.setButtons()
		.browser.Goto(.readRecords())
		}

	loadOkay?()
		{
		return Suneido.ShowNewEntries.GetDefault('status', false)
		}

	Destroy()
		{
		.Window.RestoreAccels(.curAccels)
		if .loadOkay?()
			BookNotification.ClearNewEntries(0, .max_read)
		Suneido.Delete(#ShowNewEntries)
		super.Destroy()
		}
	}
