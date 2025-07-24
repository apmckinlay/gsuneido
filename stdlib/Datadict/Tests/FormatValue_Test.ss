// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(FormatValue('', 'city') is: '')
		Assert(FormatValue('Saskatoon', 'city') is: 'Saskatoon')
		Assert(FormatValue(false, 'boolean_yesno') is: 'no')
		Assert(FormatValue(true, 'boolean_yesno') is: 'yes')
		Assert(FormatValue('', 'boolean_yesno') is: 'no')
		ob = #(1 2 3)
		Assert(FormatValue(ob, 'fred') is: Display(ob))
		Assert(FormatValue(#20170101, 'fred') is: #20170101.ShortDate())
		Assert(FormatValue(#20170101.1256, 'fred') is: #20170101.1256.ShortDateTime())
		Assert(FormatValue(0, 'dollars') is: '$0.00')
		Assert(FormatValue(12.5, 'dollars') is: '$12.50')
		Assert(FormatValue(12.5, 'dollars') is: '$12.50')
		Assert(FormatValue('<span style=><b>abc</b></span>', 'scintilla_rich') is: 'abc')
		infoText = "Email:joe@abc.com" $ InfoControl.LabelDelimiter $ 'SAMPLE LABEL'
		Assert(FormatValue(infoText, 'info') is: 'Email:joe@abc.com')
		Assert(FormatValue(123, 'number') is: 123)

		field = .TempTableName()
		rateMask = .TempName()
		.MakeLibraryRecord([name: 'Field_' $ field,
			text: `Field_dollars
				{
				Format: (OptionalNumber mask: '-` $ rateMask $ `')
				}`])
		.MakeLibraryRecord([name: rateMask, text: `'#,###,###.##'`])
		Assert(FormatValue(.8, field) is: '$.80')
		Assert(FormatValue(1.8, field) is: '$1.80')
		Assert(FormatValue(1000.8, field) is: '$1,000.80')

		field2 = .TempTableName()
		.MakeLibraryRecord([name: 'Field_' $ field2,
			text: `Field_dollars
				{
				Format: (OptionalNumber mask: '-#,###,###.##')
				}`])
		Assert(FormatValue(.8, field2) is: '$.80')
		Assert(FormatValue(1.8, field2) is: '$1.80')
		Assert(FormatValue(1000.8, field2) is: '$1,000.80')
		}
	}
