// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_isCustomizedMandatory()
		{
		isCustomizedMandatory = RecordControl.RecordControl_isCustomizedMandatory
		Assert(isCustomizedMandatory(#()) is: false)
		Assert(isCustomizedMandatory(#(Custom: false)) is: false)
		Assert(isCustomizedMandatory(#(Custom: #(), Name: 'test_field')) is: false)
		Assert(isCustomizedMandatory(#(Custom: #(test_field: #()), Name: 'test_field'))
			is: false)
		Assert(isCustomizedMandatory(
			#(Custom: #(test_field: #(first_focus:)), Name: 'test_field'))
			is: false)
		Assert(isCustomizedMandatory(
			#(Custom: #(test_field: #(mandatory:)), Name: 'test_field'))
			is: true)
		Assert(isCustomizedMandatory(
			#(Custom: #(test_field: #(mandatory:, first_focus:)), Name: 'test_field'))
			is: true)
		}
	}