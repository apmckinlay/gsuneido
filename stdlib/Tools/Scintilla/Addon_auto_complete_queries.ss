// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Addon_auto_complete
	{
	AutoComplete(word)
		{
		if word.Size() >= 2
			.AutoShow(word, .matchingWords(word))
		}

	// NOTE: This list must be kept alphabetically sorted to work with .Matches
	keywords: (alter, average, cascade, count, create, delete,
		drop, ensure, extend, false, index, insert, intersect, join,
		leftjoin, list, lower, max, min, minus, project, remove, rename,
		reverse, semijoin, sort, summarize, times, total, true, union, unique,
		update, view, where)
	matchingWords(word)
		{
		matches = []
		.Matches(.tables, word, matches)
		.Matches(.fields, word, matches)
		.Matches(.keywords, word, matches)
		return matches.Sort!().Unique!()
		}

	tables: ()
	fields: ()
	IdleAfterChange()
		{
		.fields = .collectFields(.tables = .collectTables())
		}

	collectTables()
		{
		return TablesList(includeViews?:)
		}

	collectFields(tables)
		{
		fields = Object()
		for token in QueryScanner(.Get())
			if tables.BinarySearch?(token)
				fields.Add(@QueryColumns(token))
		return fields.Sort!().Unique!()
		}
	}