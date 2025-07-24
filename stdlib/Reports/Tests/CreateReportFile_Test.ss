// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_checkReportOptions()
		{
		checkReportOptions = CreateReportFile.CreateReportFile_checkReportOptions
		rpt = Mock()
		rpt.When.GetSelectFields().Return(sf = Mock())
		rpt.When.Menu_params_fields().Return(#())
		sf.Fields = #(a, b, d_param, e_name, e_abbrev, boolfield?)
		sf.When.GetConverted().Return(#())

		checkReportOptions(rpt, #())
		checkReportOptions(rpt, #(a:, b:))
		checkReportOptions(rpt, #(a_param:, b:))
		checkReportOptions(rpt, #(a_param:, d_param:))
		checkReportOptions(rpt, #(e_num_param:, b:, e_name_param:, e_abbrev_param:))
		checkReportOptions(rpt, #(boolfield_param:))
		Assert( { checkReportOptions(rpt, #('boolfield?_param':)) }
			throws: 'Invalid Report Options')
		Assert( { checkReportOptions(rpt, #(a_param:, b:, c:)) }
			throws: 'Invalid Report Options')

		// handle when base report has menu options modified
		rpt = Mock()
		rpt.When.GetSelectFields().Return(sf)
		rpt.When.Menu_params_fields().Return(#('a','b', 'newField'))
		params = Object(a_param:, b_param:)
		checkReportOptions(rpt, params)
		Assert(params is: #(newField_param: #(operation: ""), a_param:, b_param:))
		}
	}
