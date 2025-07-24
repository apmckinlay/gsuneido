// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cache: false
	Setup()
		{
		.cache = Suneido.GetDefault("ForeignKeyTables", false)
		Suneido.Delete('ForeignKeyTables')
		.sf = SelectFields(#())
		for f in .fields
			.sf.AddField(@f)
		}
	fields:
		(
		('amount', 'Amount')
		('tax', 'Tax')
		('totamount', 'Total Amount')
		('totamount', 'Total Amount Sold')
		('invnum', 'Invoice #')
		('printed', 'Printed?')
		('compercent', '% Commission')
		('ab', 'A / B')
		('ref1', 'Ref (1)')
		('test_date', 'Date')
		)
	Test_FormulaPromptsToFields()
		{
		for f in .fields
			Assert(f[0] is: .sf.FormulaPromptsToFields(f[1], Object()))
		for f in .fields
			for g in .fields
				Assert(f[0] $ '/' $ g[0]
					is: .sf.FormulaPromptsToFields(f[1] $ '/' $ g[1], Object()))
		Assert('(amount+tax).Abs()',
			is: .sf.FormulaPromptsToFields('(Amount+Tax).Abs()', Object()))
		Assert('Date()' is: .sf.FormulaPromptsToFields('Date()', Object()))
		}
	Test_FormulaFields()
		{
		Assert(#(amount, tax), is: .sf.FormulaFields('(Amount+Tax).Abs()'))
		}
	Test_PromptToField()
		{
		sf = SelectFields
			{
			SelectFields_field_ob: #(One: one, Two: two, Three: three)
			}
		Assert(sf.PromptToField("One") is: "one")
		Assert(sf.PromptToField("Two") is: "two")
		Assert(sf.PromptToField("Three") is: "three")
		Assert(sf.PromptToField("Four") is: false)
		}
	Test_CustomFieldDuplicatePrompts()
		{
		p1 = .TempTableName()
		uniqueField = "custom_" $ p1
		.MakeLibraryRecord([name: "Field_" $ uniqueField, text: .fakeCustomField(p1)])

		p2 = .TempTableName()
		dupFieldNoSelect = "custom_" $ p2
		.MakeLibraryRecord(
			[name: "Field_" $ dupFieldNoSelect, text: .fakeCustomField(p2)])

		dupFieldWithSelect = "custom_" $ p2 $ '_2'
		.MakeLibraryRecord(
			[name: "Field_" $ dupFieldWithSelect, text: .fakeCustomField(p2, true)])

		// reverse order of fields with/without select prompts
		p3 = .TempTableName()
		dupFieldWithSelect2 = "custom_" $ p3
		.MakeLibraryRecord(
			[name: "Field_" $ dupFieldWithSelect2, text: .fakeCustomField(p3, true)])

		dupFieldNoSelect2 = "custom_" $ p3 $ '_2'
		.MakeLibraryRecord(
			[name: "Field_" $ dupFieldNoSelect2, text: .fakeCustomField(p3)])

		fields = Object(uniqueField, dupFieldNoSelect, dupFieldWithSelect,
			dupFieldNoSelect2, dupFieldWithSelect2)
		//PromptToField should not depend on order of field processed
		for fieldOb in Object(fields, fields.Copy().Reverse!())
			{
			sf = new SelectFields(fieldOb, joins: false)
			Assert(sf.PromptToField(p1) is: uniqueField)
			Assert(sf.PromptToField(p2) is: dupFieldNoSelect)
			Assert(sf.PromptToField(p3) is: dupFieldNoSelect2)
			Assert(sf.PromptToField(p2 $ " ~ Other Table") is: dupFieldWithSelect)
			Assert(sf.PromptToField(p3 $ " ~ Other Table") is: dupFieldWithSelect2)
			}
		}

	fakeCustomField(prompt, selectPrompt = false)
		{
		text = `Field_string_custom
			{
			Prompt: ` $ Display(prompt)
		if selectPrompt
			text $= `
			SelectPrompt: ` $ Display(prompt $ " ~ Other Table")
		return text $= `
			}`
		}

	Test_exclude_fields()
		{
		excludeFn = SelectFields.SelectFields_exclude_fields
		mock = Mock()
		field_list = Object('string', 'internal', 'abc', 'number')
		Assert(mock.Eval(excludeFn, field_list, #('abc'), false) is:
			#('string', 'number'))
		}
	Test_ScanFormula()
		{
		.test_ScanFormula('', '')
		.test_ScanFormula('123 + Total Amount', '123 + totamount')
		.test_ScanFormula('Total Amount + Total Amount Sold', 'totamount + totamount')
		.test_ScanFormula("'Ref (1) is ' $ Ref (1)", "'Ref (1) is ' $ ref1")
		.test_ScanFormula('"Total Amount % Commission" $ Total Amount $ "Tax"',
			'"Total Amount % Commission" $ totamount $ "Tax"')
		}
	test_ScanFormula(src, expected)
		{
		s = ''
		block = { s $= it }
		.sf.ScanFormula(src, block, block)
		Assert(s, is: expected)
		}
	Test_ScanFields()
		{
		Assert(.sf.ScanFields('') is: #())

		Assert(.sf.ScanFields('123 + 456') is: #(
			#(pos: 0, end: 3)
			#(pos: 6, end: 9)))

		Assert(.sf.ScanFields('((123 + 456) + Tax)') is: #(
			#(pos: 2, end: 5)
			#(pos: 8, end: 11)
			#(pos: 15, end: 18)))

		Assert(.sf.ScanFields('123 + 456 + Total Amount Sold') is: #(
			#(pos: 0, end: 3)
			#(pos: 6, end: 9)
			#(pos: 12, end: 29)))

		Assert(.sf.ScanFields('123 + 456 + Total Amount Sold $ "123 Amount"') is: #(
			#(pos: 0, end: 3)
			#(pos: 6, end: 9)
			#(pos: 12, end: 29)
			#(pos: 32, end: 44)))
		}
	Test_text_match()
		{
		fn = SelectFields.SelectFields_text_match

		cases = #(``, `abc 123 ABC $%^`, `\"\'`, `\n\r`, `\\`)
		for first in #('"', "'", '`')
			cases.Each()
				{
				s = first $ it $ first
				Assert(fn(first, s) is: s)
				}

		Assert(fn('1', '123abc') is: '')
		Assert(fn('`', '`123') is: '')
		Assert(fn('"', '"123') is: '')
		Assert(fn('"', `"abc\"!@#'"123`) is: `"abc\"!@#'"`)
		}
	Test_number_match()
		{
		fn = SelectFields.SelectFields_number_match

		Assert(fn(false, '123') is "123")
		Assert(fn(false, '123.123') is "123.123")
		Assert(fn(false, '.123') is ".123")
		Assert(fn(false, '.123 * 100') is ".123")
		Assert(fn(false, '.123.456 * 100') is ".123")
		Assert(fn(false, '100.123.456 * 100') is "100.123")
		Assert(fn(false, 'aaa100.123.456 * 100') is "")
		Assert(fn(false, ' 100.123.456 * 100') is "")
		}
	Test_first()
		{
		m = SelectFields.SelectFields_first
		Assert(m('Test this') is: 'Test')
		Assert(m('Test + This') is: 'Test')
		Assert(m('Test_abc + This') is: 'Test_abc')
		Assert(m('_abc + This') is: '_abc')
		u = '_'.Repeat(300)
		Assert(m(u $ ' Test') is: u)
		Assert(m('(Amount+Tax).Abs()') is: '(')
		Assert(m('ABC') is: 'ABC')
		Assert(m('A') is: 'A')
		Assert(m('A+BC') is: 'A')
		Assert(m('+ABC') is: '+')
		Assert(m('/ABC') is: '/')
		Assert(m('++ABC') is: '+')
		Assert(m('((ABC+BC)*D)') is: '(')
		}

	Test_Joins()
		{
		mastertable1 = .MakeTable('(sftest_num, sftest_name, sftest_abbrev)
			key(sftest_num)')
		mastertable2 = .MakeTable('(sftest2_num, sftest2_name, sftest2_abbrev)
			key(sftest2_num)')
		.MakeTable('(sftest_num, sftrantest_num, sftrantest_desc, sftest2_num_billto)
			key(sftrantest_num)
			index(sftest_num) in ' $ mastertable1 $
			' index(sftest2_num_billto) in ' $ mastertable2 $ '(sftest2_num)')
		.MakeLibraryRecord([name: "Field_sftest_num",
			text: `Field_string { Prompt: "sftest_num" }`])
		.MakeLibraryRecord([name: "Field_sftest2_num",
			text: `Field_string { Prompt: "sftest2_num" }`])
		.MakeLibraryRecord([name: "Field_sftrantest_num", text: `Field_internal {} `])
		.MakeLibraryRecord([name: "Field_sftrantest_desc",
			text: `Field_string { Prompt: "sftestdesc" }`])
		fields = #("sftest_num", "sftrantest_desc", "sftest2_num_billto")
		sf = SelectFields(fields)
		Assert(sf.Fields hasSubset: #(
			"sftest2_num_billto Abbrev": "sftest2_abbrev_billto",
			"sftest2_num_billto Name": "sftest2_name_billto",
			"sftest_num Abbrev": "sftest_abbrev", "sftest_num Name": "sftest_name",
			sftestdesc: "sftrantest_desc"))
		joins = sf.Joins(#("sftest_num Abbrev", "sftest_num Name", "sftestdesc",
			"sftest2_num_billto Abbrev", "sftest2_num_billto Name"))
		q1 = " leftjoin by(sftest_num) (" $ mastertable1 $
			" project sftest_num, sftest_name, sftest_abbrev) "
		Assert(joins has: q1)
		q2 = " leftjoin by(sftest2_num_billto) (" $ mastertable2 $
			" project sftest2_num, sftest2_name, sftest2_abbrev " $
			"rename sftest2_num to sftest2_num_billto " $
			"rename sftest2_name to sftest2_name_billto " $
			"rename sftest2_abbrev to sftest2_abbrev_billto)"
		Assert(joins has: q2)

		joinob = sf.JoinsOb(#("sftest_num Abbrev", "sftest_num Name", "sftestdesc",
			"sftest2_num_billto Abbrev", "sftest2_num_billto Name"), withDetails?:)
		pos1 = joinob.FindIf({ it.str is q1 })
		Assert(pos1 isnt: false, msg: 'pos 1')
		Assert(joinob[pos1].fields equalsSet: #(sftest_abbrev, sftest_name))
		pos2 = joinob.FindIf({ it.str.Trim() is q2.Trim() })
		Assert(pos2 isnt: false, msg: 'pos 2')
		Assert(joinob[pos2].fields.Sort!()
			is: #(sftest2_abbrev_billto, sftest2_name_billto))
		}

	Test_nameAndAbbrevExist()
		{
		fields = Object('prefix_num', 'prefix_abbrev', 'prefix_name')
		Assert(SelectFields.SelectFields_nameAndAbbrevExist(fields, 'prefix',
			'prefix_num'))

		fields = Object('prefix_num_suffix', 'prefix_abbrev_suffix', 'prefix_name_suffix')
		Assert(SelectFields.SelectFields_nameAndAbbrevExist(fields, 'prefix',
			'prefix_num_suffix'))

		fields = Object('prefix_num', 'otherField_name', 'otherField_abbrev')
		Assert(SelectFields.SelectFields_nameAndAbbrevExist(fields, 'prefix',
			'prefix_num') is: false)
		}

	Teardown()
		{
		if .cache isnt false
			Suneido.ForeignKeyTables = .cache
		else
			Suneido.Delete('ForeignKeyTables')
		super.Teardown()
		}
	}
