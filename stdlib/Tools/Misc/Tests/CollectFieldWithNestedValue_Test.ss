// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	TestFunc()
		{
		table = .MakeTable('(ky, cols) key (ky)', [cols: 'noField', ky: 1])
		expected = 'wantedField'
		Assert(CollectFieldWithNestedValue(table, 'cols', expected, 'ky') is: Object())
		QueryOutput(table, [cols: 'noField, anotherField, wantedField', ky: 2])
		Assert(CollectFieldWithNestedValue(table, 'cols', expected, 'ky') is: #(2))

		table2 = .MakeTable('(ky,  columnA, columnB) key (ky)',
			[columnA: 'noField', columnB: 'noField', ky: 1])
		Assert(CollectFieldWithNestedValue(table2, #('columnA', 'columnB'),
			expected, 'ky') is: Object())
		QueryOutput(table2, [columnA: 'noField, anotherField',
			columnB: 'wrongField, alsoWrongField' ky: 2])
		Assert(CollectFieldWithNestedValue(table2, #('columnA', 'columnB'),
			expected, 'ky') is: Object())

		QueryDo('update ' $ table2  $ ' where ky is 2
			 set columnB = "wantedField, anotherField, moreFields"')
		Assert(CollectFieldWithNestedValue(table2, #('columnA', 'columnB'),
			expected, 'ky') is: #(2))
		Assert(CollectFieldWithNestedValue(table2, #('columnA'), expected, 'ky')
			is: Object())
		Assert(CollectFieldWithNestedValue(table2, #('columnB'), expected, 'ky') is: #(2))

		QueryOutput(table2, [columnA: 'wantedField, anotherField',
			columnB: 'wrongField, alsoWrongField' ky: 3])
		QueryOutput(table2, [columnA: 'anotherField',
			columnB: 'wantedField' ky: 4])
		Assert(CollectFieldWithNestedValue(table2, #('columnA'), expected, 'ky') is: #(3))
		Assert(CollectFieldWithNestedValue(table2, #('columnA', 'columnB'),
			expected, 'ky') equalsSet: #(2, 3, 4))
		Assert(CollectFieldWithNestedValue(table2, #('columnB'), expected, 'ky')
			equalsSet: #(2, 4))
		}
	}