// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	return Qc_Main.CheckWithExtra(.table, .name, .text, minimizeOutput?:, extraChecks:).
		Filter({ Object?(it) and it.GetDefault('warnings', #()).NotEmpty?() }).
		Map!({ it.desc }).Join('\r\n')
	}