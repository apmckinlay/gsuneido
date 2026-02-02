// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (hwnd, wpPlace)
	{
	window = SuRenderBackend().GetRegisteredControl(hwnd)
	window.GetWindowPlacement(wpPlace)
	}