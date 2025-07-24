// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_invalid_operator?()
		{
		// test string op with non String fields
		dd = Field_date {}
		op = #('contains', '=~', pre: '(?i)(?q)', suf: '')
		Assert(Select2.Invalid_operator?(op, dd) is: true)
		dd = Field_number {}
		Assert(Select2.Invalid_operator?(op, dd) is: true)

		// test operator with Image (only empty and not empty are valid)
		dd = Field_image {}
		op = #('equals', 'is', pre: '', suf: '')
		Assert(Select2.Invalid_operator?(op, dd) is: true)
		op = #('empty', 'is', pre: '', suf: '')
		Assert(Select2.Invalid_operator?(op, dd) is: false)
		op = #('not empty', 'isnt', pre: '', suf: '')
		Assert(Select2.Invalid_operator?(op, dd) is: false)

		// test valid operator
		dd = Field_string {}
		op = #('contains', '=~', pre: '(?i)(?q)', suf: '')
		Assert(Select2.Invalid_operator?(op, dd) is: false)
		}

	Test_formatSummarizeCalcField()
		{
		Assert(Select2.Select2_formatSummarizeCalcField('testselect2field',
			#(calc0: Field_number), 'Field_string') is: 'Field_string')
		Assert(Select2.Select2_formatSummarizeCalcField('testselect2_total_calc0',
			#(calc0: Field_number), 'Field_string') is: 'Field_string')

		Assert(Select2.Select2_formatSummarizeCalcField('total_calc0',
			#(calc0: Field_number), 'Field_string') is: 'Field_number')
		Assert(Select2.Select2_formatSummarizeCalcField('max_calc1',
			#(calc0: Field_number, calc1: Field_date), 'Field_string') is: 'Field_date')
		}

	Test_Empty_field()
		{
		fn = Select2.Empty_field

		sf = SelectFields(#('test_field', 'test_field_num'),
			headerSelectPrompt: 'no_prompts')
		fld = 'test_field'
		op = Object('equals')
		Assert(fn(fld, op, sf) is: 'test_field')

		op[0] = 'empty'
		Assert(fn(fld, op, sf) is: 'test_field')

		fld = 'not_in_sf'
		Assert(fn(fld, op, sf) is: 'not_in_sf')

		op[0] = 'equals'
		fld = 'test_field_abbrev'
		Assert(fn(fld, op, sf) is: 'test_field_abbrev')

		op[0] = 'suffix is empty'
		Assert(fn(fld, op, sf) is: 'test_field_num')
		fld = 'test_field_name'
		Assert(fn(fld, op, sf) is: 'test_field_num')

		sf = SelectFields(#('test_field', 'test_field_name'),
			headerSelectPrompt: 'no_prompts')
		Assert(fn(fld, op, sf) is: 'test_field_name')
		}

	Test_value()
		{
		fn = { |val|
			Select2.Select2_value(#(), 'string',
				Select2.Ops[12/*= matches*/], [val0: val], 0)
			}
		Assert(fn('(abc)') is: Object(errs: '', val: Display('(abc)'),
			selectFunction: false))
		Assert(fn('abc)') is: Object(errs: 'Invalid value ' $ Display('abc)') $
				' for operator ' $ Display("matches"), val: false))

		fn = { |val|
			Select2.Select2_value(#(), 'string',
				Select2.Ops[8/*= contains*/], [val0: val], 0)
			}
		Assert(fn('(abc)') is: Object(errs: '', val: Display('(?i)(?q)(abc)'),
			selectFunction: false))
		Assert(fn('abc)') is: Object(errs: '', val: Display('(?i)(?q)abc)'),
			selectFunction: false))

		fn = { |val|
			Select2.Select2_value(#(), 'string',
				Select2.Ops[4/*= equals*/], [val0: val], 0)
			}
		Assert(fn('(abc)') is: Object(errs: '', val: Display('(abc)'),
			selectFunction: false))
		Assert(fn('abc)') is: Object(errs: '', val: Display('abc)'),
			selectFunction: false))
		}

	Test_Where()
		{
		GetForeignNumsFromNameAbbrevFilter_Test.SetupForeignTable(this)
		sf = SelectFields(#(table2_num, table2_name, table1_num_test))
		select2 = Select2(sf)
		where = select2.Where([
			menu_option0:, print0:, checkbox0:, fieldlist0: "Name 2"
			menu_option1:, print1:, checkbox1:, fieldlist1: "Num 1 test Name"])
		Assert(where is: #(where: '', errs: '', joinflds: #()))

		where = select2.Where([
			menu_option0:, print0:, checkbox0:, fieldlist0: "Name 2", oplist0: "equals",
			val0: 'abc',
			menu_option1:, print1:, checkbox1:, fieldlist1: "Num 1 test Name"])
		Assert(where is: #(where: ' where table2_name is "abc"', errs: '',
			joinflds: #(table2_name)))

		where = select2.Where([
			menu_option0:, print0:, checkbox0:, fieldlist0: "Name 2",
			oplist0: "not equal to", val0: 'abc',
			menu_option1:, print1:, checkbox1:, fieldlist1: "Num 1 test Name"])
		Assert(where is: #(where: ' where table2_name isnt "abc"', errs: '',
			joinflds: #(table2_name)))

		where = select2.Where([
			menu_option0:, print0:, checkbox0:,
			fieldlist0: "Name 2", oplist0: "equals", val0: 'abc',
			menu_option1:, print1:, checkbox1:,
			fieldlist1: "Num 1 test Name", oplist1: "equals", val1: 'abc'])
		Assert(where
			is: #(where: ' where table2_name is "abc" where table1_num_test in ()',
				errs: '', joinflds: #(table2_name, table1_num_test)))

		where = select2.Where([
			menu_option0:, print0:, checkbox0:,
			fieldlist0: "Name 2", oplist0: "equals", val0: 'abc',
			menu_option1:, print1:, checkbox1:,
			fieldlist1: "Num 1 test Name", oplist1: "equals", val1: 'BBB'])
		Assert(where
			is: #(where: ' where table2_name is "abc" where table1_num_test in (2)',
				errs: '', joinflds: #(table2_name, table1_num_test)))
		}
	}
