// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	baseName: #bizTestNextNameTest
	Test_main()
		{
		.TearDownIfTablesNotExist(#params)
		Params.Ensure()

		name = .baseName
		QueryOutput(#params,
			[
				params: [report_name: name, Source: "Recurring Payable Invoices"],
				report: "Reporter - " $ name])

		fn = ReporterActions.ReporterActions_renameReport
		renamed_name = #bizTestNextNameTest_Renamed
		fn(name, renamed_name, #simple, skipUpdateBookMenu:)
		rec = Query1(`params`, report: "Reporter - " $ renamed_name)
		Assert(rec isnt false)
		Assert(Query1(`params`, report: "Reporter - " $ name) is: false)

		ReporterActions.Delete(renamed_name)
		Assert(Query1(`params`, report: "Reporter - " $ renamed_name) is: false)
		}

	Teardown()
		{
		QueryDo('delete params
			where report >= ' $ Display('Reporter - ' $ .baseName) $
			' and report < ' $ Display('Reporter - ' $ .baseName $ `~`))
		super.Teardown()
		}
	}
