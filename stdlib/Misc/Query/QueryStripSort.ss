// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// DEPRECATED: use Query.StripSort (built-in after #20250422)
function (query)
	{
	if BuiltDate() >= #20250422
		return Query.StripSort(query)

	// Most of the time there is no sort (e.g. 95% of the time)
	// so it is better to do a fast Has? first
	if not query.Has?("sort")
		return query
	if Sys.Client?()
		return ServerEval("QueryStripSort", query)
	tableHint = query.Has?("/* tableHint:")
		? "/* tableHint: " $ QueryGetTable(query, orview:) $ " */ "
		: ""
	q = Suneido.ParseQuery(query)
	if q.type isnt "sort"
		return query
	// keep checkquery suppressions
	sups = QueryGetSuppressions(query)
	keep = sups.Empty?() ? '' : ' ' $ sups.Join(' ')
	return tableHint $ q.source.String $ keep
	}
