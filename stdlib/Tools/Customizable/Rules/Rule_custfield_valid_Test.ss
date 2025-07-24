// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert([].custfield_valid is: '')
		Assert([custfield_mandatory: true, custfield_readonly: true].custfield_valid
			is: "can not make mandatory field read-only/hidden.")
		Assert([custfield_mandatory: true, custfield_hidden: true].custfield_valid
			is: "can not make mandatory field read-only/hidden.")

		field1 = .TempName()
		.MakeLibraryRecord([name: "Field_" $ field1, text: `Field_string
			{
			Prompt: 'ABC'
			}`])
		field2 = .TempName()
		.MakeLibraryRecord([name: "Field_" $ field2, text: `Field_string
			{
			Prompt: 'DEF'
			}`])

		rec = [custfield_formula: 'abc = 123', custfield_fields_list: [field1, field2]]
		Assert(rec.custfield_valid is: "")
		rec.custfield_formula = "ABC = 123"
		Assert(rec.custfield_valid is: "Can not assign value to field in formula")
		rec.custfield_formula = "if true\n\tDEF = 123"
		Assert(rec.custfield_valid is: "Can not assign value to field in formula")
		rec.custfield_formula = "if true\n\tABCD = 123"
		Assert(rec.custfield_valid is: "")
		rec.custfield_formula = "return ABC"
		Assert(rec.custfield_valid is: "")
		rec.custfield_formula = "ABC"
		Assert(rec.custfield_valid is: "")
		rec.custfield_formula = "var = ABC"
		Assert(rec.custfield_valid is: "")
		}

	Test_customField_only_fillin_from()
		{
		rec = [custfield_field: 'efg', custfield_only_fillin_from: 'hij']
		Assert(rec.custfield_valid is: 'Only Fill-in From is not allowed for this field')

		rec = [custfield_field: 'efg', custfield_only_fillin_from: '']
		Assert(rec.custfield_valid is: '')

		rec = [custfield_field: 'custom_900001', custfield_only_fillin_from: 'hij']
		Assert(rec.custfield_valid is: '')

		rec = [custfield_field: 'custom_900001', custfield_only_fillin_from: '']
		Assert(rec.custfield_valid is: '')
		}
	}