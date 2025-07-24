// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_table_not_exist()
		{
		table = 'tests' $ Display(Timestamp()).Tr('#.', '_')
		Assert({ GetNextNumber(table, 'nextnum') } throws: 'nonexistent table')
		}

	Test_table_empty()
		{
		table = .MakeTable('(nextnum) key()')
		Assert({ GetNextNumber(table, 'nextnum') } throws: 'no records')
		}

	Test_table_valid()
		{
		table = .MakeTable('(nextnum) key()', [nextnum: 1000])
		Assert(GetNextNumber(table, 'nextnum') is: 1000)
		Assert(GetNextNumber(table, 'nextnum') is: 1001)
		}

	Test_conflict()
		{
		table = .MakeTable('(nextnum) key()', [nextnum: 1000])
		Suneido.forceTooManyRetryTransaction = true
		Assert({ GetNextNumber(table, 'nextnum') } throws: 'too many retries')
		Suneido.forceTooManyRetryTransaction = false
		}
	}
