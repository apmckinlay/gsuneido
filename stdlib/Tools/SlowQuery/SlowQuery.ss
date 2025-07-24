// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
class
	{
	ensure()
		{
		Database('ensure slow_queries
			(sq_num, sq_hash, sq_query, sq_time, sq_last_used)
			key (sq_num)
			key (sq_hash)')
		}

	LogIfTooSlow(t, query, _slowQueryLog = false, block = false)
		{
		if t > 5 /*= seconds */ and not query.Has?('=~') and
			not query.Has?('/*SLOWQUERY SUPPRESS*/')
			{
			if slowQueryLog isnt false
				{
				if slowQueryLog.GetDefault('logged', false) is true or
					slowQueryLog.GetDefault('suppressed', false) is true
					return
				slowQueryLog.logged = true
				}
			if block isnt false
				block(.log(query, t))
			}
		}

	log(query, time)
		{
		.ensure()
		query = query.Trim()
		sha1 = Sha1(query)
		hash = sha1.ToHex()

		RetryTransaction()
			{ |tran|
			tran.QueryDo('delete slow_queries
				where sq_hash is ' $ Display(hash))
			t = Timestamp()
			tran.QueryOutput('slow_queries',
				[sq_num: t, sq_last_used: t,
				sq_hash: hash, sq_query: query, sq_time: time])
			}
		return Base64.Encode(sha1) // for SuneidoLog, it checks first 30 characters
		}

	slow?(query)
		{
		if not TableExists?('slow_queries')
			return false
		hash = Sha1(query.Trim()).ToHex()
		QueryApply1('slow_queries', sq_hash: hash)
			{
			it.sq_last_used = Timestamp()
			it.Update()
			return true
			}
		return false
		}

	Count()
		{
		if not TableExists?('slow_queries')
			return 0
		return QueryCount('slow_queries')
		}

	Purge()
		{
		if not TableExists?('slow_queries')
			return
		QueryApplyMulti('slow_queries
			where sq_last_used < ' $ Display(Date().Plus(months: -3)), update:)
			{
			it.Delete()
			}
		}

	Validate(query, allCols, after, queryState, indexes = false)
		{
		if not .slow?(query)
			return true

		if indexes is false
			indexes = .SelectableIndexes(query, allCols)
		if indexes.Empty?() // no suggestion
			return true

		queryState.indexes = indexes
		type = queryState.sortCol is false ? 'Select' : 'Sort'
		queryState.msg = 'Your current ' $ type $ ' could be slow. ' $
			'Adding additional criteria can make it faster.'

		.SuggestionWindow(queryState, after)
		return false
		}

	SelectableIndexes(query, allCols)
		{
		return (.getIndexes)(query).Intersect(allCols)
		}

	getIndexes: Memoize
		{
		Func(query)
			{
			if '' is table = QueryGetTable(query, nothrow:)
				return #()
			return QueryList('indexes where table is ' $ Display(table), 'columns').
				Map({ it.BeforeFirst(',') }).
				UniqueValues().
				RemoveIf(Internal?).
				Map({ f = query.Extract(it $ '\sto\s(\w+)'); f is false ? it : f }).
				SortWith!({ Object(it !~ '^date|_date', SelectPrompt(it)) })
			}
		}

	SuggestionWindow(queryState, after)
		{
		args = Object(['SlowQuerySuggestions', queryState],
			closeButton?: false, keep_size: false, onDestroy: after)
		if after is false
			Dialog(@(args.Add(0, at: 0)))
		else
			ModalWindow(@args)
		}

	AddIndexedFilter(filter, topFilters)
		{
		newFilters = topFilters.Get()
		conditions = newFilters.conditions
		conditions.RemoveIf({ it.condition_field is filter.condition_field and
			(it.check isnt true or it[it.condition_field].operation in ('', 'all')) })
		conditions.Add(filter)
		topFilters.Set(newFilters)
		topFilters.HighlightLastRow(conditions.Size() - 1)
		topFilters.SetSelectApplied(false)
		}

	AddParamsIndexedFilter(filter, recordControl, filterField)
		{
		filterCtrl = recordControl.FindControl(filterField)
		conditions = filterCtrl.Get()

		conditions.RemoveIf({ it.condition_field is filter.condition_field and
			it[it.condition_field].operation in ('', 'all') })
		conditions.Add(filter)

		recordControl.SetField(filterField, conditions)
		lastRow = conditions.Size() - 1
		if filterCtrl.Method?('FocusRow')
			filterCtrl.FocusRow(lastRow, value?:)
		}
	}