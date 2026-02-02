// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (hwnd)
	{
	if false is window = SuRenderBackend().WindowManager.GetWindow(hwnd)
		return NULL
	return SuRenderBackend().WindowManager.ActivateWindow(window)
	}