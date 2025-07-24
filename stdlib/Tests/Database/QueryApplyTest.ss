// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		QueryApply('stdlib', update:) {|unused| break }
		Assert({ QueryApply('stdlib sort group', update:) {|unused| } }
			throws: 'QueryEnsureKeySort: query has non-key sort')
		QueryApply('stdlib where name = "QueryApplyTest"', update:)  {|unused| break }
		QueryApply('stdlib sort name,group', update:) {|unused| break }
		QueryApply('stdlib sort reverse name,group', update:) {|unused| break }
		table = .MakeTable('(num, name) key()')
		QueryApply(table, update:) {|unused| }
		}
	}
