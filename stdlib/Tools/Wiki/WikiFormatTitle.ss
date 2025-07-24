// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function (title)
	{
	return title.
		Replace('([^A-Z])([A-Z])', `\1 \2`).
		Replace('([^0-9])([0-9])', `\1 \2`).
		Replace("([A-Z])([A-Z][^A-Z ])", `\1 \2`).
		Trim()
	}