// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
LibIOTestBase
	{
	Test_one()
		{
		.VerifyExportImportForDelete("rec1", #20190402, useSVCDates: false)
		.VerifyExportImportForDelete("rec2", #20190402, useSVCDates: true)
		.VerifyExportImportForDelete("rec3", #20190402, hasDeleted:, useSVCDates: false)
		.VerifyExportImportForDelete("rec4", #20190402, hasDeleted:, useSVCDates: true)
		}
	}
