// Copyright (C) 2025 Axon Development Corporation All rights reserved worldwide.
function (hwnd)
	{
	windowManager = SuRenderBackend().WindowManager
	return windowManager.ActivateWindow(windowManager.GetWindow(hwnd))
	}