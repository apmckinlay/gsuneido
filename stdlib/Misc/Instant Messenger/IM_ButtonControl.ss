// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New()
		{
		super(.layout())

		if IM_Available?()
			{
			.messengerButton = .FindControl('messengerButton')
			.subs = [
				PubSub.Subscribe('messengerMessages', .messengerMessages)
				PubSub.Subscribe('messengerOffline', .messengerOffline)
				PubSub.Subscribe('messengerOnline', .messengerOnline)
				]
			IM_MessengerManager.MessengerThread()
			.initStatus()
			}
		}

	layout()
		{
		return IM_Available?()
			? #(Horz
				#(Vert Fill (EnhancedButton, command: 'IM', image: 'im.emf', mouseEffect:,
					imageColor: 0x009933, mouseOverImageColor: 0x006903, imagePadding: .1,
					tip: 'Open Instant Messenger', name: 'messengerButton') Fill)
				#Skip)
			: #(Vert)
		}

	initStatus()
		{
		if Suneido.GetDefault(#IM_Status, '') is ''
			return
		Suneido.IM_Status is 'online'
			? .messengerOnline()
			: .messengerOffline()
		}

	green: 0x009933
	red: 0x0000ff
	grey: 0xc0c0c0
	darkGreen: 0x003300
	darkRed: 0x0000aa
	messengerButton: false
	messengerMessages(notify = true, setRead = false)
		{
		if IM_MessengerControl.Opened?()
			return

		if .messengerButton is false //or .messengerButton.GetImageColor() is .red
			return

		if notify
			.messengerButton.SetImageColor(.red, .darkRed)

		if setRead
			.messengerButton.SetImageColor(.green, .darkGreen)

		.flashWindow(setRead, notify)
		}

	flashWindow(setRead, notify)
		{
		if super.HasFocus?() is false and setRead
			FlashWindowEx(Object(cbSize: FLASHWINFO.Size(), hwnd: .Window.Hwnd,
				dwFlags: FLASHW.STOP, uCount: 0))

		if super.HasFocus?() is false and notify
			FlashWindowEx(Object(cbSize: FLASHWINFO.Size(), hwnd: .Window.Hwnd,
				dwFlags: FLASHW.ALL|FLASHW.TIMERNOFG, uCount: 15))
		}

	messengerOffline()
		{
		if .messengerButton isnt false and
			.messengerButton.GetImageColor() isnt .grey
			.messengerButton.SetImageColor(.grey, .grey)
		}

	messengerOnline()
		{
		if .messengerButton isnt false and
			.messengerButton.GetImageColor() is .grey
			.messengerButton.SetImageColor(.green, .darkGreen)
		}

	On_IM()
		{
		if .messengerButton is false
			return

		if true isnt TestHttpServer()
			{
			AlertError('Unable to connect to Instant Messenger server.')
			return
			}
		IM_MessengerManager.MessengerThread()
		.messengerButton.SetImageColor(.green, .darkGreen)
		IM_MessengerControl()
		}

	Destroy()
		{
		if IM_Available?()
			{
			.subs.Each(#Unsubscribe)
			IM_MessengerManager.StopThread()
			IM_MessengerManager.Request("/IM_ContactsUpdate")
			}
		super.Destroy()
		}
	}
