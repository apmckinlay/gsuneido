// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// asof:
	// - Date:	Queries only one asof period. Useful for getting a record at a set time
	// - []: 	Queries from compact date to most recent transaction
	// - [from: <date>, to: <date>]: Queries a specific range
	CallClass(query, asof = [])
		{
		if Date?(asof)
			asof = [from: asof, to: asof]
		return ServerEval('QueryHistory.Collect', query, asof)
		}

	GUI(query, parent = 0)
		{
		if false is asof = QueryHistoryAsofControl(parent)
			return false
		Working('Collecting History...', { results = this(query, asof) })
		return results
		}

	Collect(query, asof)
		{
		results = Object()
		Transaction(read:)
			{ .queryAsof(it, query, asof, results) }
		return results
		}

	queryAsof(t, query, asof, results)
		{
		from = t.Asof(asof.GetDefault(#from, Date.Begin()))
		tAsof = t.Asof(asof.GetDefault(#to, Date.End()))
		key = ShortestKey(query)
		lastPack = Object().Set_default([])
		while tAsof isnt false and tAsof >= from
			{
			try
				t.QueryApply(query)
					{
					if lastPack[it[key]] is curPack = Pack(it)
						continue
					lastPack[it[key]] = curPack
					it.tran_asof
					results.Add(it)
					}
			catch (unused, 'nonexistent table')
				{ }
			tAsof = t.Asof(-1)
			}
		results.Sort!(By(#tran_asof))
		}
	}
