// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	result = ""
	QueryApply('views')
		{ |x|
		if "" isnt s = CheckQuery(x.view_definition)
			result $= x.view_name $ " " $ s $ "\n"
		}
	return result
	}