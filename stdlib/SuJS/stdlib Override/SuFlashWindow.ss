// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(windowId, show?)
		{
		if show?
			SuRender().Notification.Notify('New message(s)', windowId)
		else
			SuRender().Notification.OnWindowActivated(windowId)
		}
	}
