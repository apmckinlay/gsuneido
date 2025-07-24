// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	TestFunc()
		{
		table = .MakeTable('(ky, cols) key (ky)', [cols: 'noField', ky: 1])
		expected = 'wantedField'
		Assert(TableHasNestedValue?(table, 'cols', expected) is: false)
		QueryOutput(table, [cols: 'noField, anotherField, wantedField', ky: 2])
		Assert(TableHasNestedValue?(table, 'cols', expected))

		table2 = .MakeTable('(ky,  columnA, columnB) key (ky)',
			[columnA: 'noField', columnB: 'noField', ky: 1])
		Assert(TableHasNestedValue?(table2, #('columnA', 'columnB'), expected) is: false)
		QueryOutput(table2, [columnA: 'noField, anotherField',
			columnB: 'wrongField, alsoWrongField' ky: 2])
		Assert(TableHasNestedValue?(table2, #('columnA', 'columnB'), expected) is: false)

		QueryDo('update ' $ table2  $ ' where ky is 2
			 set columnB = "wantedField, anotherField, moreFields"')
		Assert(TableHasNestedValue?(table2, #('columnA', 'columnB'), expected))
		Assert(TableHasNestedValue?(table2, #('columnA'), expected) is: false)
		Assert(TableHasNestedValue?(table2, #('columnB'), expected))
		}
	}