// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (hwnd, msg, wParam, lParam)
	{
	Assert(msg isString:)
	SuRenderBackend().RecordAction(hwnd, msg, [wParam, lParam])
	}