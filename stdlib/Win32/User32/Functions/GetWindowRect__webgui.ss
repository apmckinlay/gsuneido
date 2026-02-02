// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (hwnd)
	{
	GetWindowPlacement(hwnd, place = Object())
	return place.rcNormalPosition
	}