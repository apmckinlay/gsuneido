// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_children()
		{
		m = SchemaModel.SchemaModel_children

		Assert(m(#(), 100) is: #())

		Assert(m(#([column: field_a, field: 0]), 100)
			is: #([name: field_a, num: 100, group: false]))

		cols = [
			[column: #field_a, field: 0],
			[column: #field_b, field: 1],
			[column: #field_c, field: 2],
			[// Uppercase rule columns:
				column: #rule_a, field: -1],
			[column: #rule_b, field: -1]]
		expectedCols = [
			[name: #field_a, num: 100, group: false],
			[name: #field_b, num: 101, group: false],
			[name: #field_c, num: 102, group: false],
			[name: #rule_a, num: 103, group: false],
			[name: #rule_b, num: 104, group: false]]
		Assert(m(cols, 100) is: expectedCols)

		Assert(m(#([column: rule_a, field: -1]), 100)
			is: #([name: rule_a, num: 101, group: false]))
		}
	}
