// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
#(
ExtensionPoints:
	(
	('Current'), ('Global')
	)
Contributions:
	(
	(AccessMenus, Current, Inspect, devel:,
		function (data, hwnd) { Inspect(data, :hwnd) })
	(AccessMenus, Current, 'Go To QueryView', devel:,
		function (data/*unused*/, hwnd/*unused*/, access)
			{
			GotoQueryView(access.NewRecord?()
				? access.AccessControl_query
				: access.AccessControl_keyquery)
			})
	(AccessMenus, Current, History,
		view: function (data, hwnd, historyFields)
			{
			ViewHistory(data, hwnd, historyFields)
			}
		update: function(data, historyFields)
			{
			UpdateHistory(data, historyFields)
			})
	)
)
