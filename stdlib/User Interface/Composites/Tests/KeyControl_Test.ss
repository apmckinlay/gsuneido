// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Key_BuildQuery()
		{
		// override Send for test class
		idClass = KeyControl { Send(@args/*unused*/) { return "" } }

		// make sure query has sort and is handled properly (needs to be at end)
		query = "test_query sort testsortfield"

		// no restrictions
		restrictions = invalidRestrictions = false
		Assert(idClass.Key_BuildQuery(query, restrictions, invalidRestrictions)
			is: "test_query sort testsortfield")

		query = function () { return "test_query sort testsortfield" }
		// no restrictions
		restrictions = invalidRestrictions = false
		Assert(idClass.Key_BuildQuery(query, restrictions, invalidRestrictions)
			is: "test_query sort testsortfield")

		// just restrictions
		restrictions = "testfield1 is 'testval1' and testfield2 is 'testval2'"
		Assert(idClass.Key_BuildQuery(query, restrictions, invalidRestrictions)
			sameText: BuiltDate() < #20250422
				? "test_query where testfield1 is 'testval1' and
					testfield2 is 'testval2' sort testsortfield"
				: "test_query sort testsortfield
					where testfield1 is 'testval1' and testfield2 is 'testval2'")

		// both restrictions and invalidRestrictions
		invalidRestrictions = "testfield3 is 'testval3'"
		Assert(idClass.Key_BuildQuery(query, restrictions, invalidRestrictions)
			sameText: BuiltDate() < #20250422
				? "test_query where testfield1 is 'testval1' and
					testfield2 is 'testval2' where testfield3 is 'testval3'
					sort testsortfield"
				: "test_query sort testsortfield
					where testfield1 is 'testval1' and testfield2 is 'testval2'
					where testfield3 is 'testval3'")

		// just invalidRestrictions
		restrictions = false
		Assert(idClass.Key_BuildQuery(query, restrictions, invalidRestrictions)
			sameText: BuiltDate() < #20250422
				? "test_query where testfield3 is 'testval3' sort testsortfield"
				: "test_query sort testsortfield where testfield3 is 'testval3'")

		// just optionalRestrictions (ChooseDate)
		.MakeLibraryRecord([name: 'Field_test_field'
			text: `Field_date
				{
				Prompt: 'Terminate Date'
				Control: (ChooseDate)
				}`])
		optRestrictions = #(test_field)
		Assert(idClass.Key_BuildQuery(query, false, false, :optRestrictions)
			sameText: BuiltDate() < #20250422
				? 'test_query  where (test_field is "" or test_field >= ' $
					Display(Date().NoTime()) $ ') sort testsortfield'
				: 'test_query sort testsortfield where (test_field is "" or
					test_field >= ' $ Display(Date().NoTime()) $ ')')

		// just optionalRestrictions (ChooseList)
		.MakeLibraryRecord([name: 'Field_test_field2'
			text: `Field_string
				{
				Prompt: 'Status'
				Control: (ChooseList #("active", "inactive"))
				}`])
		optRestrictions = #(test_field2)
		Assert(idClass.Key_BuildQuery(query, false, false, :optRestrictions)
			sameText: BuiltDate() < #20250422
				? 'test_query where test_field2 is "active" sort testsortfield'
				: 'test_query sort testsortfield where test_field2 is "active"')

		// both optionalRestrictions
		optRestrictions = #(test_field, test_field2)
		Assert(idClass.Key_BuildQuery(query, false, false, :optRestrictions)
			sameText: BuiltDate() < #20250422
				? 'test_query where (test_field is "" or test_field >= ' $
					Display(Date().NoTime()) $
					')  where test_field2 is "active" sort testsortfield'
				: 'test_query sort testsortfield where (test_field is "" or
					test_field >= ' $ Display(Date().NoTime()) $ ')
					where test_field2 is "active"')

		// all restrictions
		invalidRestrictions = 'testfield3 is "testval3"'
		restrictions = 'testfield1 is "testval1" and testfield2 is "testval2"'
		Assert(idClass.Key_BuildQuery(query, restrictions, invalidRestrictions,
			:optRestrictions)
			sameText: BuiltDate() < #20250422
				? 'test_query  where testfield1 is "testval1" and ' $
					'testfield2 is "testval2"  where testfield3 is "testval3"  ' $
					'where (test_field is "" or test_field >= ' $
					Display(Date().NoTime()) $
					')  where test_field2 is "active" sort testsortfield'
				: 'test_query sort testsortfield where testfield1 is "testval1" and
					testfield2 is "testval2" where testfield3 is "testval3"
					where (test_field is "" or
					test_field >= ' $ Display(Date().NoTime()) $ ')
					where test_field2 is "active"')
		}
	}