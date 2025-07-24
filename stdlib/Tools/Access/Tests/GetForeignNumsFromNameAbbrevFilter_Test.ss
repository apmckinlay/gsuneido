// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	SetupForeignTable(testClass = false)
		{
		if testClass is false
			testClass = this

		if false isnt cache = Suneido.GetDefault(#ForeignKeyTables, false)
			cache.Reset()

		table1 = testClass.MakeTable('(table1_num, table1_name, table1_abbrev)
			key (table1_num)',
			[table1_num: 1, table1_name: 'AAA', table1_abbrev: 'aaa'],
			[table1_num: 2, table1_name: 'BBB', table1_abbrev: 'bbb'],
			[table1_num: 3, table1_name: 'CCC', table1_abbrev: 'ccc'])
		table2 = testClass.MakeTable('(table2_num, table1_num, table2_name)
			key (table2_num)
			index (table1_num) in ' $ table1)

		testClass.MakeLibraryRecord([
			name: "Field_table1_name_test",
			text: `Field_string { Prompt: "Name 1" }`])
		testClass.MakeLibraryRecord([
			name: "Field_table1_abbrev_test",
			text: `Field_string { Prompt: "Abbrev 1" }`])
		testClass.MakeLibraryRecord([
			name: "Field_table1_num_test",
			text: `Field_num { Prompt: "Num 1 test" }`])
		testClass.MakeLibraryRecord(
			[name: "Field_table2_name", text: `Field_string { Prompt: "Name 2" }`])
		testClass.MakeLibraryRecord(
			[name: "Field_table2_num", text: `Field_num { Prompt: "Num 2" }`])
		return [:table1, :table2]
		}
	Test_main()
		{
		.SetupForeignTable()
		sf = SelectFields(#(table2_num, table2_name, table1_num_test))

		fn = GetForeignNumsFromNameAbbrevFilter
			{
			GetForeignNumsFromNameAbbrevFilter_limit: 3
			}
		// operation is ''
		Assert(fn('table1_abbrev_test', sf, '') is: false)
		// not a foreign name/abbrev field
		Assert(fn('table2_name', sf, 'equals', 'test') is: false)
		Assert(fn('test_field', sf, 'equals', 'test') is: false)
		// invalid operation
		msg = "ERROR: GetParamsWhere.CallClass: invalid operation: invalid"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg, msg, msg))
		Assert(fn('table1_abbrev_test', sf, 'invalid', 'test') is: false)

		Assert(fn('table1_abbrev_test', sf, 'equals', 'test')
			is: #(numField: 'table1_num_test', nums: ()))
		Assert(fn('table1_abbrev_test', sf, 'equals', 'bbb')
			is: #(numField: 'table1_num_test', nums: (2)))

		res = fn('table1_abbrev_test', sf, 'greater than', 'aaa')
		Assert(res.numField is: 'table1_num_test')
		Assert(res.nums equalsSet: #(2, 3))

		res = fn('table1_abbrev_test', sf, 'does not contain', 'aaa')
		Assert(res.numField is: 'table1_num_test')
		Assert(res.nums equalsSet: #(2, 3, ""))

		// over limit
		Assert(fn('table1_abbrev_test', sf, 'greater than', 'a') is: false)
		}
	}
