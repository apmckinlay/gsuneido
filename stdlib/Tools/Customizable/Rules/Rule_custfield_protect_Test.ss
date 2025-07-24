// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		user = Suneido.User
		Suneido.User = 'default'

		.buildTestCustFields()
		rec = Record(custfield_field: 'testfield_unprotected')
		protect = rec.custfield_protect
		Assert(protect hasntMember: 'custfield_mandatory')
		Assert(protect hasntMember: 'custfield_hidden')
		Assert(protect hasntMember: 'custfield_tabover')
		Assert(protect hasntMember: 'custfield_readonly')
		Assert(protect hasMember: 'allowDelete')

		rec = Record(custfield_field: 'testfield_protected')
		protect = rec.custfield_protect
		Assert(protect hasMember: 'custfield_mandatory')
		Assert(protect hasMember: 'custfield_tabover')
		Assert(protect hasMember: 'custfield_readonly')
		Assert(protect hasMember: 'custfield_only_fillin_from')
		Assert(protect hasntMember: 'custfield_hidden')

		rec = Record(custfield_field: 'custom_000001')
		protect = rec.custfield_protect
		Assert(protect hasntMember: 'custfield_only_fillin_from')

		.MakeLibraryRecord([name: 'Rule_test_field', text: 'function() {return true}'])
		.MakeLibraryRecord(
			[name: 'Rule_test_field_a__protect', text: 'function() {return true}'])

		rec = Record(custfield_field: 'test_field')
		protect = rec.custfield_protect
		Assert(protect hasMember: 'custfield_formula')

		rec = Record(custfield_field: 'test_field_a')
		protect = rec.custfield_protect
		Assert(protect hasMember: 'custfield_formula')

		rec = Record(custfield_field: 'test_field_b')
		protect = rec.custfield_protect
		Assert(protect hasntMember: 'custfield_formula')

		rec = Record(custfield_field: 'testfield_AllowOff')
		protect = rec.custfield_protect
		Assert(protect hasMember: 'custfield_hidden')
		Suneido.User = user
		}

	buildTestCustFields()
		{
		.MakeLibraryRecord([name: 'CustFieldUnProtectedControl'
			text: 'FieldControl
				{
				New(mandatory = false, hidden = false, tabover = false, readonly = false)
					{
					super()
					}
				}'])
		.MakeLibraryRecord([name: 'Field_testfield_unprotected',
			text: 'Field_string
				{
				Prompt: "Test Protected"
				Control: (CustFieldUnProtected)
				 }'])

		.MakeLibraryRecord([name: 'CustFieldProtectedControl'
			text: 'FieldControl
				{
				New(hidden = false)
					{
					super()
					}
				}'])
		.MakeLibraryRecord([name: 'Field_testfield_protected',
			text: 'Field_string
				{
				Prompt: "Test Protected"
				Control: (CustFieldProtected)
				 }'])
		.MakeLibraryRecord([name: 'Field_testfield_AllowOff',
			text: 'Field_string
				{
				Prompt: "Test Protected"
				Control: (CustFieldProtected)
				AllowCustomizableOptions: (hidden: false)
				 }'])
		.MakeLibraryRecord([name: 'Field_custom_000001',
			text: 'Field_string
				{
				Prompt: "Test custom"
				 }'])
		}

	Test_override()
		{
		text = 'Field_string
			{
			Prompt: "Test 1"
			Control: (ScintillaAddonsEditor xmin: 250)
			}'
		.MakeLibraryRecord([name: 'Field_test1_field' text: text])

		rec = Record(custfield_field: 'test1_field')
		protect = rec.custfield_protect
		Assert(protect hasMember: 'custfield_mandatory')
		Assert(protect hasMember: 'custfield_readonly')

		text2 = 'Field_string
			{
			Prompt: "Test 2"
			Control: (ScintillaAddonsEditor xmin: 250)
			AllowCustomizableOptions: (readonly:, mandatory:)
			}'
		.MakeLibraryRecord([name: 'Field_test2_field' text: text2])

		rec = Record(custfield_field: 'test2_field')
		protect = rec.custfield_protect
		Assert(protect hasntMember: 'custfield_mandatory')
		Assert(protect hasntMember: 'custfield_readonly')
		}
	}