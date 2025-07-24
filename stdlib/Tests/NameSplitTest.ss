// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(NameSplit('Bob').first is: 'Bob')
		Assert(NameSplit('Bob').last is: '')

		Assert(NameSplit('Bobby Smith').first is: 'Bobby')
		Assert(NameSplit('Bob Smith ').last is: 'Smith')

		Assert(NameSplit('Bob P. Smith').first is: 'Bob P.')
		Assert(NameSplit('Bob P. Smith').last is: 'Smith')

		Assert(NameSplit('Smith, Sue').first is: 'Sue')
		Assert(NameSplit('Smith, Sue Z.').first is: 'Sue Z.')
		Assert(NameSplit('Jones, Bob').last is: 'Jones')

		Assert(NameSplit('Mary Ann', split_on: ',').first is: 'Mary Ann')
		Assert(NameSplit('Smith, Mary Ann', split_on: ',').first is: 'Mary Ann')
		Assert(NameSplit('Smith, Mary Ann', split_on: ',').last is: 'Smith')
		}
	}