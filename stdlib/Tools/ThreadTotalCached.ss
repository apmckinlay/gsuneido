// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
class
	{
	cacheSize: 200
	CallClass(calcOnserverFn, query, screenForTotals, filters)
		{
		if .client?()
			return ServerEval('ThreadTotalCached', calcOnserverFn, query, screenForTotals,
				filters)
		try
			{
			.processFilters(filters)
			screenTotals = 'ThreadTotal_' $ screenForTotals
			_queryTotals = QueryStripSort(query)
			if not Suneido.Member?(screenTotals)
				Suneido[screenTotals] = LruCache(.getFunc, .cacheSize)
			return Suneido[screenTotals].Get(calcOnserverFn, filters)
			}
		catch (e)
			{
			SuneidoLog('ERROR: ' $ e, params: filters)
			ErrorLog('filter: ' $ Json.Encode(filters))
			throw e
			}
		}

	// extracted for testing
	client?()
		{
		return Sys.Client?()
		}

	processFilters(filters)
		{
		filters.RemoveIf({ it[it.condition_field].operation is '' }).Sort!({|x,y|
			x.condition_field < y.condition_field })
		}

	getFunc(calcOnserverFn, filters)
		{
		q = _queryTotals
		filters = filters // just to avoid unused argument
		return Global(calcOnserverFn)(q)
		}

	Reset(screenTotals)
		{
		if .client?()
			return ServerEval('ThreadTotalCached.Reset', screenTotals)

		if false isnt cache = Suneido.GetDefault(screenTotals, false)
			cache.Reset()
		}
	}
