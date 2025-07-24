// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	QueryApply(@args)
		{
		NameArgs(args, #(query, block, dir), #(Next))
		dir = args.Extract(#dir)
		block = args.block
		Assert(args.Extract(#update, "") is "",
			msg: "t.QueryApply should not have update: argument")
		readonly = args.Extract(#readonly, false)
		if .Update?() and not readonly
			{
			if false isnt (newquery = QueryAddKeySort(args.query))
				args.query = newquery
			else if Suneido.User is 'default'
				SuneidoLog("t.QueryApply failed to add key sort, " $
					"consider removing sort or specifying readonly:", calls:)
			}
		args.block = {|q|
			while false isnt x = q[dir]()
				try
					block(x)
				catch (e, "block:")
					if e is "block:break"
						break
					// else block:continue ... so continue
			}
		.Query(@args)
		}
	QueryApply1(@args)
		{
		Assert(args.Extract(#update, "") is "",
			msg: "t.QueryApply1 should not have update: argument")
		Assert(.Update?(),
			msg: "t.QueryApply1 should only be used on update Transactions")
		NameArgs(args, #(query, block))
		block = args.Extract(#block)
		if false isnt x = .Query1(@args)
			block(x)
		}
	QueryAccum(query, accum, block) // DEPRECATED
		{
		.QueryApply(query)
			{ |x|
			accum = block(accum, x)
			}
		return accum
		}
	QueryOutput(query, record)
		{
		.Query(query)
			{ |q|
			q.Output(record)
			}
		}
	QueryDo(@args)
		{
		result = .Query(@args)
		if (not Number?(result))
			throw "QueryDo: not a request"
		return result
		}
	SeekQuery(query)
		{
		return SeekQuery(this, query)
		}
	QueryMax(query, field, default = false)
		{
		x = .Query1(QueryStripSort(query) $ ' summarize max ' $ field)
		return x is false ? default : x['max_' $ field]
		}
	QueryMin(query, field, default = false)
		{
		x = .Query1(QueryStripSort(query) $ ' summarize min ' $ field)
		return x is false ? default : x['min_' $ field]
		}
	QueryCount(query)
		{
		x = .Query1(QueryStripSort(query) $ '\nsummarize count')
		return x is false ? 0 : x.count
		}
	QueryTotal(query, field)
		{
		x = .Query1(QueryStripSort(query) $ '\nsummarize total ' $ field)
		return x is false ? 0 : x['total_' $ field]
		}
	QueryAverage(query, field)
		{
		x = .Query1(QueryStripSort(query) $ '\nsummarize average ' $ field)
		return x is false ? false : x['average_' $ field]
		}
	QueryList(query, field)
		{
		x = .Query1(QueryStripSort(query) $ '\nsummarize list ' $ field)
		return x is false ? [] : x['list_' $ field]
		}
	Query1Cached(@args)
		{
		cache = .Data().GetInit(#Query1Cache,
			{ LruCache({ .Query1(@it) }, 100) }) /*= cache size */
		return cache.Get(args)
		}
	QueryAll(@args)
		{
		NameArgs(args, #(query, limit), #(false))
		limit = args.Extract(#limit)
		Assert(limit is false or Number?(limit),
			'QueryAll limit must be a number, got: ' $ limit)
		sort = #()
		reverse = false
		if limit is false
			{
			// sort in memory
			sort = QueryGetSort(args.query, keepReverse:)
			if reverse = sort.Prefix?("reverse ")
				sort = sort.RemovePrefix("reverse ")
			sort = sort.Split(',')
			args.query = QueryStripSort(args.query)
			}
		list = Object()
		args.block =
			{|q|
			while false isnt x = q.Next()
				if Number?(limit) and list.Size() >= limit
					break
				else
					list.Add(x)
			}
		.Query(@args)
		if not sort.Empty?()
			list.Sort!(By(@sort))
		if reverse
			list.Reverse!()
		return list
		}

	QueryRange(query, field)
		{
		query = QueryStripSort(query) $ ' summarize min ' $ field $ ', max ' $ field
		x = .Query1(query)
		return x is false ? false :
			Object(min: x['min_' $ field], max: x['max_' $ field])
		}

	QueryEmpty?(@args)
		{
		if BuiltDate() < #20250424
			{
			args[0] = QueryStripSort(args[0]) $
				" /* CHECKQUERY SUPPRESS: SORT REQUIRED */"
			return false is .QueryFirst(@args)
			}
		return not .QueryExists?(@args)
		}

	QueryAny1(@args) // query [, fields...]
		{
		query = args[0]
		if QueryStripSort(query) isnt query
			ProgrammerError('QueryAny1 does not take a query with a sort')

		fields = args.Delete(0)
		// will check that the values are the same on up to 10 records
		if Suneido.GetDefault('ValidateQueryAny1?', false) is true
			{
			count = 0
			orig = false
			.QueryApply(query)
				{
				if orig is false
					orig = it.Project(fields)
				else if orig isnt it.Project(fields)
					throw "QueryAny1: all recs do not match"
				if ++count > 10 /*= max records to check to limit affect on performance */
					break
				}
			}

		rec = .QueryFirst(query $ " /* CHECKQUERY SUPPRESS: SORT REQUIRED */")
		return rec is false ? false : rec.Project(fields)
		}
	}