// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	AddWhere(query, where)
		{
		return QueryAddWhere(query, where)
		}

	AvailableColumns(query)
		{
		return QueryAvailableColumns(query)
		}

	GetTable(query, nothrow)
		{
		return QueryGetTable(query, :nothrow)
		}

	StripSort(query)
		{
		return QueryStripSort(query)
		}

	ExtendColumns(query, sf, fields)
		{
		joinob = sf.JoinsOb(fields, withDetails?:)
		query = BuildQueryJoin(query, joinob)
		missingRules = fields.Difference(QueryColumns(query))
		return query $ Opt(' extend ', missingRules.Join(', '), ' ')
		}
	}