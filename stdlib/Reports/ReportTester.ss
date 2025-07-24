// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Report
	{
	// NOTE: trackAccesses only tracks Access Points that are explicitly added
	// by the report (will not track passivly added ones (i.e. because of datadicts)
	New(@args)
		{
		super(@args)
		.trackAccesses = args.GetDefault('trackAccesses', false)
		}
	Test()
		{
		return .Run(ReportTestDriver(.trackAccesses), quiet?:)
		}

	GetFont()
		{
		return Document_Builder.GetDefaultFont()
		}

	VerifyLines(report, size, lines)
		{
		report.Map!({ it.SplitCSV() })
		Assert(report isSize: size)
		// Using idx so only the lines we want to verify can be passed in
		// i.e. lines can have
		/*[ [0: #(column list)]
		[3: #(first line item)]
		[4: #(second line item)]
		[5: #(third line item)]
		[8: #(grand total)] ]
		*/
		// if we only want to test the columns, line items and total and
		// 1,2 are headers and 6,7 are subtotals
		for idx in lines.Members()
			Assert(report[idx] is: lines[idx])
		}
	}
