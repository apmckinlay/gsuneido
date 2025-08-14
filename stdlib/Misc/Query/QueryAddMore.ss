// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// similar to QueryAddWhere but handles more than where
// NOTE: use QueryAddWhere if just adding a where
function (query, more)
	{
	// use \n in case query/more ends with //comment
	if BuiltDate() > #20250422
		return Query.StripSort(query) $ Opt("\n", more) $
			Opt("\nsort ", Query.GetSort(query))

	if not query.Has?("sort")
		return query $ Opt("\n", more)
	if Sys.Client?()
		return ServerEval("QueryAddMore", query, more)
	tableHint = query.Has?("/* tableHint:")
		? "/* tableHint: " $ QueryGetTable(query, orview:) $ " */ "
		: ""
	q = Suneido.ParseQuery(query)
	if q.type isnt "sort"
		return query $ Opt("\n", more)
	// keep checkquery suppressions
	sups = QueryGetSuppressions(query)
	keep = sups.Empty?() ? '' : ' ' $ sups.Join(' ')
	return tableHint $ q.source.String $ keep $ Opt("\n", more) $ Opt("\n", q.string)
	}
