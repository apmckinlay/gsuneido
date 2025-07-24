// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function (name)
	{
	name = name.Lower()
	if false isnt i = ProvinceNames.FindIf({ it.Lower() is name })
		return ProvinceCodes[i]
	if false isnt i = StateNames.FindIf({ it.Lower() is name })
		return StateCodes[i]
	return false
	}
