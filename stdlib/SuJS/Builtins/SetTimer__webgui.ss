// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
function (hwnd/*unused*/, id/*unused*/, ms, f)
	{
	return SuRenderBackend().TimerManager.SetTimer(ms, f)
	}