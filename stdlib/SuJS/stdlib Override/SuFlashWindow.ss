// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(windowId, show?, text = 'New message(s)')
		{
		if show?
			SuRender().Notification.Notify(text, windowId)
		else
			SuRender().Notification.OnWindowActivated(windowId)
		}
	}
