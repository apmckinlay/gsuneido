// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
function (query, block)
	{
	key = ShortestKey(query)
	QueryApply(query)
		{|x|
		where = ''
		for field in key.Split(',')
			where $= ' where ' $ field $ ' is ' $ Display(x[field])
		lookupQuery = QueryStripSort(query) $ where
		RetryTransaction()
			{|t|
			y = t.Query1(lookupQuery)
			// this could fail if another user modified the key
			// after we started the QueryApply
			if y isnt false
				{
				try
					block(t, y)
				catch (ex, "block:")
					if ex is "block:break"
						break
					// else block:continue ... so continue
				}
			else
				SuneidoLog('INFO: QueryApplyEach skipped ' $ lookupQuery,
					params: x)
			}
		}
	}