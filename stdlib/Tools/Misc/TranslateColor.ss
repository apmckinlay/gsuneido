// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function (color)
	{
	return String?(color) and CLR.Member?(color) ? CLR[color] : color
	}
