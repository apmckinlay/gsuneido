// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (wParam)
	{
	clicks = (HISWORD(wParam) / WHEEL_DELTA).Int()
	wsl = SPI_GetWheelScrollLines()
	lines = clicks * wsl
	return Object(:lines, :clicks, page?: wsl is -1, down?: clicks < 0)
	}