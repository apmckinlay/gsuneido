// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		foreignRec = [frannt_num: num = Timestamp(),
				frannt_name: 'Big Old Test', frannt_abbrev: 'bot',
				frannt_name_lower!: 'big old test']
		_findForeignRecTableTest = .MakeTable('
			(frannt_num, frannt_name, frannt_abbrev, frannt_name_lower!)
			key (frannt_num)
			key (frannt_name)
			key (frannt_name_lower!)
			index unique (frannt_abbrev)',
			foreignRec)

		foreignRec.Delete('frannt_name_lower!')
		func = FindForeignRecWithAbbrevNameOrNum
			{
			FindForeignRecWithAbbrevNameOrNum_getNumTable(basename)
				{
				return basename is 'frannt'
					? _findForeignRecTableTest
					: false
				}
			FindForeignRecWithAbbrevNameOrNum_queryKeys(unused)
				{
				return #('frannt_num', 'frannt_name', 'frannt_name_lower!')
				}
			FindForeignRecWithAbbrevNameOrNum_uniqueIndexes(unused)
				{
				return #('frannt_abbrev')
				}
			}

		Assert(func([], 'notarealtable_num') is: false)
		Assert(func([], 'badtable_name') is: false)
		Assert(func([], 'frannt_name') is: false)

		Assert(func([frannt_num: num], 'frannt_name') is: false)
		Assert(func([frannt_name: 'big old test'], 'frannt_name') is: foreignRec)
		Assert(func([frannt_num: num], 'frannt_name', useNum?:) is: foreignRec)
		Assert(func([frannt_num_renamed: num], 'frannt_name', useNum?:) is: foreignRec)
		Assert(func([different_num_renamed: num], 'frannt_name', useNum?:) is: false)
		}
	}