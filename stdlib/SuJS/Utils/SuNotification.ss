// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	unavailable: false
	New(.taskbar)
		{
		try
			SuUI.GetCurrentWindow().Notification
		catch
			{
			Print("This browser does not support desktop notification")
			.unavailable = true
			return
			}
		if not .hasPermission?()
			SuUI.GetCurrentWindow().Notification.RequestPermission()

		.notifications = Object().Set_default(Object())
		}

	Notify(msg, windowId)
		{
		.taskbar.SetFlashing(windowId, true)
		if .unavailable is true or not .hasPermission?()
			return

		.cleanup()
		if .notifications[windowId].Member?(msg)
			.notifications[windowId][msg].close()

		try
			{
			notification = SuUI.MakeWebObject('Notification', msg, #(
				requireInteraction:,
				icon: '/favicon.ico'))
			}
		catch
			{
			.unavailable = true
			return
			}
		notification.AddEventListener('close',
			{ |event/*unused*/| .notifications[windowId].Remove(notification) })
		notification.AddEventListener('click',
			{ |event/*unused*/|
			notification.close()
			SuRender().ActivateWindow(windowId)
			})
		.notifications[windowId][msg] = notification
		}

	cleanup()
		{
		for id in .notifications.Members().Copy()
			if .notifications[id].Empty?()
				.notifications.Delete(id)
		}

	hasPermission?()
		{
		return SuUI.GetCurrentWindow().Notification['permission'] is "granted"
		}

	OnWindowActivated(windowId)
		{
		.taskbar.SetFlashing(windowId, false)
		if .unavailable is true or not .notifications.Member?(windowId)
			return

		for notification in .notifications[windowId]
			notification.close()
		}
	}
