// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
class
	{
	GetQuery(source)
		{
		// if there is no book open (from running demo data), do not do the authorize
		 // or if there is no "name" member i.e. stdlib only query based
		if not Suneido.Member?('browser') or not source.Member?(#name)
			query = source.query
		else
			query = false is ReporterDataSource.Authorized?(source.name)
				? ''
				: source.query
		if query.BeforeFirst('.').GlobalName?() and not Uninit?(query.BeforeFirst('.'))
			query = Global(query)()
		return QueryStripSort(query)
		}

	GetKeys(query)
		{
		return QueryKeys(query)
		}

	GetColumns(query)
		{
		return QueryColumns(query)
		}

	HasColumn?(query, field)
		{
		return query isnt "" and .GetColumns(query).Has?(field)
		}

	BuildExclude(source)
		{
		exclude = #()
		if Object?(source.exclude)
			exclude = source.exclude
		else if Function?(source.exclude)
			exclude = (source.exclude)()
		else if String?(source.exclude) and source.exclude isnt ''
			exclude = (source.exclude.Eval())() // needs Eval
		exclude = exclude.Copy()
		if false isnt cl = OptContribution('CustomTabPermissions', false)
			cl.ReporterExcludes(source, exclude)
		return exclude
		}

	CheckQuery(query)
		{
		try
			QueryStrategy(query)
		catch (err)
			{
			Print(ERROR: err)
			Print(QUERY: query)
			if err.Has?("invalid query")
				throw "Reporter: Unable to generate report. Sort field(s) may be invalid."
			else if err.Has?("invalid column(s) in expressions")
				throw "Reporter: Unable to generate report. " $
					"Please check formulas. Formulas must be defined " $
					"prior to other formulas that use them."
			else
				throw "Reporter: Error in report: " $ err
			}
		}
	}