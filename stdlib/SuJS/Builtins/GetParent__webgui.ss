// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
function (hwnd)
	{
	if false is control = SuRenderBackend().GetRegisteredControl(hwnd)
		return NULL
	if not control.Member?(#Parent)
		return NULL
	return control.Parent.Hwnd
	}