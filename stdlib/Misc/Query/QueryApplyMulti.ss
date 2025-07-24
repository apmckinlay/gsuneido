// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
// QueryApplyMulti chooses between two strategies based on their estimated cost:
// - withCursor uses a cursor with multiple update transactions (the old way)
//	 This is better to process a large percentage of the records in the table.
// - withLookup uses a read-only QueryApply and then lookups in update transactions
//	 This is better to process a small percentage of the records in the table.
class
	{
	updatesPerTran: 100
	CallClass(query, block, update)
		{
		if not update
			{
			QueryApply(query, block)
			return
			}

		lookupCost = .lookupCost(query)
		if lookupCost.nrecs > 0 // use lookup for 0 nrecs, possibly avoiding update tran
			{
			cursorCost = .cursorCost(query)
			if cursorCost.cost < lookupCost.cost
				{
				.withCursor(cursorCost.query, block)
				return
				}
			}
		.withLookup(query, block)
		}
	cursorCost(query)
		{
		Cursor(query) {|c| keys = c.Keys()}
		if "" isnt sort = QueryGetSort(query).RemovePrefix("reverse ")
			{
			s = sort.Split(',')
			for key in keys
				if #() is key.Split(',').Difference(s)
					return Object(:query, cost: CursorCost(query).cost)
			throw "non-key sort is not allowed"
			}
		// determine the best (lowest cost) key sort
		best_cost = 999999999999
		best_query = false
		for key in keys
			{
			c = CursorCost(q = query $ Opt(" sort ", key))
			if c.cost < best_cost
				{
				best_cost = c.cost
				best_query = q
				}
			}
		return Object(query: best_query, cost: best_cost)
		}
	withCursor(query, block)
		{
		Cursor(query)
			{|c|
			x = false
			do
				// can't use RetryTransaction because of c.Next
				Transaction(update:)
					{ |t|
					for (i = 0; i < .updatesPerTran and false isnt x = c.Next(t); ++i)
						{
						try
							block(x)
						catch (ex, "block:")
							if ex is "block:break"
								{
								x = false // to get out of do-while
								break
								}
							// else block:continue ... so continue
						}
					}
				while x isnt false
			}
		}
	lookupCost(query)
		{
		cost = QueryCost(query)
		cost.cost += cost.nrecs * 1000 /*= estimated lookup cost */
		return cost
		}
	withLookup(query, block)
		{
		keys = ShortestKey(query).Split(',')
		lookupQuery = QueryStripSort(query)
		i = 0
		t = false
		Finally(
			{
			QueryApply(query) // read-only
				{|x|
				if t is false
					t = Transaction(update:)
				if i++ % .updatesPerTran is 0
					{
					t.Complete()
					t = Transaction(update:)
					}
				args = .project(x, keys).Add(lookupQuery)
				if false isnt y = t.Query1(@args)
					{
					try
						block(y)
					catch (ex, "block:")
						if ex is "block:break"
							return
						// else block:continue ... so continue
					}
				}
			},
			{
			if t isnt false
				t.Complete()
			})
		}
	project(rec, keys)
		{
		// can't use Records.Project because it omits missing fields
		ob = Object()
		for f in keys
			ob[f] = rec[f]
		return ob
		}
	}