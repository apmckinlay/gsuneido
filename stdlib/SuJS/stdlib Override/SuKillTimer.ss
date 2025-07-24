// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
function (hwnd, id)
	{
	if hwnd isnt 0 and hwnd is id // KillTimer from server
		if SuRender().Timers.Member?(id)
			id = SuRender().Timers[id]
	SuUI.GetCurrentWindow().ClearInterval(id)
	}
