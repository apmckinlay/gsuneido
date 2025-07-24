// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
function (hwnd, id, ms, f, once? = false)
	{
	if hwnd isnt 0 and hwnd is id // SetTimer from server
		{
		timer = SuUI.GetCurrentWindow().
			SetInterval({ |a/*unused*/, b/*unused*/, c/*unused*/, d/*unused*/|
				SuRender().TimeOut(id)
				if once? is true
					SuUI.GetCurrentWindow().ClearInterval(SuRender().Timers[id]) },
			ms, 0, 0, 0, 0)
		SuRender().Timers[id] = timer
		return
		}
	return SuUI.GetCurrentWindow().SetInterval(f, ms, 0, 0, 0, 0)
	}
