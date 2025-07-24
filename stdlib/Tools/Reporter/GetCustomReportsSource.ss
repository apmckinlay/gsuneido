// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
function (sourceName)
	{
	if Customizable.AccessToDataSource?(sourceName)
		return Customizable.GetPermissableDataSources()[sourceName]
	contribs = false
	Plugins().ForeachContribution('Reporter', 'queries')
		{ |x|
		if x.name is sourceName
			return x
		contribs = true
		}
	return contribs ? false : [query: sourceName]
	}