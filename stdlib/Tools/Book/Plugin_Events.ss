// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
#(
ExtensionPoints:
	(
	('eventfunction')
	)
Contributions:
	(
	(Events, eventfunction, name: 'NewButtonMessages',
		serverfunc: function (user)
			{ return BookNotification.NewEvents(user) },
		clientfunc: function (result)
			{
			result = BookNotification.HandleNewEvents(result)
			if result is true
				PubSub.Publish('book notification')
			})

	(Events, eventfunction, name: 'CustomizationNotification',
		serverfunc: function ()
			{ return ServerEval('CustomizableOnServer.GetChangesOnServer') },
		clientfunc: function (result)
			{ Customizable.HandleChangesOnClient(result) })

	(Events, eventfunction, name: 'CheckJsDownloadTasks',
		serverfunc: function (user)
			{ return JsDownload.CheckTask(user) },
		clientfunc: function (result)
			{ JsDownload.WarnIfOutstanding(result) })
	)
)
