// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (hwnd)
	{
	window = SuRenderBackend().WindowManager.GetWindow(hwnd)
	Assert(window isnt: false)
	window.CLOSE()
	}