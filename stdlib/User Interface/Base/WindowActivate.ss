// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function (hwnd)
	{
	GetWindowPlacement(hwnd, place = Object(length: WINDOWPLACEMENT.Size()))
	if place.showCmd is SW.SHOWMINIMIZED
		{
		place.showCmd = SW.SHOWNORMAL
		SetWindowPlacement(hwnd, place)
		}
	SetActiveWindow(hwnd)
	}