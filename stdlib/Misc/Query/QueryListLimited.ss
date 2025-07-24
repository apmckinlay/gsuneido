// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function (query, field, limit, returnDetailIfOverLimit? = false, tran = false)
	{
	query = QueryStripSort(query)
	DoWithTran(tran)
		{ |t|
		q = query $ ' project ' $ field
		count = t.QueryCount(q)
		if count > limit
			return returnDetailIfOverLimit? is true
				? Object(:count, list: t.QueryAll(q, :limit).Map({ it[field] }))
				: false
		return t.QueryList(query, field)
		}
	}
