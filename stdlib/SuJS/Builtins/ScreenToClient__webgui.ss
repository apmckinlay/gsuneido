// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (hwnd, pt/*unused*/)
	{
	if false is control = SuRenderBackend().GetRegisteredControl(hwnd)
		throw "ScreenToClient: unexpected control - " $ Display(hwnd)
	Assert(control base: ListControl)
	}