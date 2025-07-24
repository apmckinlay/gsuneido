// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
// used by QueryApply and QueryApplyMulti
function (query)
	{
	keys = QueryKeys(query)
	if keys is #("")
		return query
	if "" isnt sort = QueryGetSort(query).RemovePrefix("reverse ")
		{
		s = sort.Split(',')
		for key in keys
			if #() is key.Split(',').Difference(s)
				return query
		return false
		}
	return query $ ' sort ' $ ShortestKey(keys)
	}