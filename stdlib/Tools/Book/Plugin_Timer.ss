// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
#(
ExtensionPoints:
	(
	('timerfunction')
	)
Contributions:
	(
	(Timer, timerfunction,
		func: function () { EventManager() },
		name: 'EventManager')

	(Timer, timerfunction,
		func: function () { IM_MessengerManager.MessengerThread() },
		name: 'StartMessengerThread')
	)
)