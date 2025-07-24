// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Addon_auto_complete
	{
	New(@args)
		{
		super(@args)
		.LongIdle()
		}
	AutoComplete(word)
		{
		if word.Size() < 2
			return
		.AutoShow(word, .matching_words(word))
		}
	matching_words(word)
		{
		matches = []
		.Matches(.tables, word, matches)
		.Matches(.fields, word, matches)
		.Matches(.keywords, word, matches)
		return matches.Sort!().Unique!()
		}

	keywords: (alter average count create delete
		drop ensure extend false index insert intersect join
		leftjoin list lower minus project remove rename reverse
		sort summarize times total true union unique update view where)

	tables: ()
	LongIdle()
		{
		.tables = QueryList('tables', 'table')
		.tables.Add(@QueryList('views', 'view_name'))
		.tables.Sort!()
		}

	fields: ()
	IdleAfterChange()
		{
		.fields = []
		for token in Scanner(.Get())
			if .tables.BinarySearch?(token)
				.tableFields(token)
		.fields.Sort!().Unique!()
		}
	tableFields(table)
		{
		.fields.Add(@QueryColumns(table))
		}
	}
