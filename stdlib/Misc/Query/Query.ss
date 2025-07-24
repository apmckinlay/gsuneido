// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
// TEMPORARY: remove after BuiltDate #20250422
class
	{
	StripSort(query)
		{
		// Most of the time there is no sort (e.g. 95% of the time)
		// so it is better to do a fast Has? first
		if not query.Has?("sort")
			return query
		if Sys.Client?()
			return ServerEval("Query.StripSort", query)
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
	GetSort(query)
		{
		if not query.Has?('sort')
			return ""
		sort = ""
		keep = false
		scnr = QueryScanner(query)
		for token in scnr
			if token is "sort"
				keep = true
			else if keep and scnr.Type() not in (#COMMENT, #WHITESPACE, #NEWLINE)
				if token is "reverse"
					sort $= "reverse "
				else
					sort $= token // column or comma
		return sort
		}
	}