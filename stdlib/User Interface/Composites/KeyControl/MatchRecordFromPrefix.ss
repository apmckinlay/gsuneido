// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function (baseQuery, field, val)
	{
	query = QueryStripSort(baseQuery) $
		" where " $ field $ " >= " $ Display(val) $
		" sort " $ field
	WithQuery(query)
		{|q|
		if false isnt (rec = q.Next()) and String(rec[field]).Prefix?(val)
			{
			if rec[field] is val
				return rec
			next = q.Next()
			if next is false or not String(next[field]).Prefix?(val)
				return rec
			}
		}
	return false
	}