// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// Based on contributions from Roberto Artigas Jr. (rartiga1@midsouth.rr.com)
Controller
	{
	Title: "Language Translations"
	New()
		{
		super(.layout())
		}
	layout()
		{
		return Object("Browse"
			'translatelanguage'
			title: 'Language Translations'
			columns: QueryColumns('translatelanguage'))
		}
	}
