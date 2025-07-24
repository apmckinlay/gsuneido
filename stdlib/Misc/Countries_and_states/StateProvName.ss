// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
function (state_prov)
	{
	if false isnt i = ProvinceCodes.Find(state_prov)
		return ProvinceNames[i]
	if false isnt i = StateCodes.Find(state_prov)
		return StateNames[i]
	return state_prov
	}