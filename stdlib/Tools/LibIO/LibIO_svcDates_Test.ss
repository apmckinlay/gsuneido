// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
LibIOTestBase
	{
	Test_exportimport_svc_dates()
		{
		// test new libary and book records, not committed, useSVCDates = false
		.VerifyExportImport('VerifyExportImportLibraries',
			lib_committed: '', useSVCDates: false,
			lib_committed_expcted: '')

		// test new libary and book records, committed, useSVCDates = false
		.VerifyExportImport('VerifyExportImportLibraries2',
			lib_committed: #20140101.1000, useSVCDates: false
			lib_committed_expcted: '')

		// test new libary and book records, committed, useSVCDates = true
		.VerifyExportImport('VerifyExportImportLibraries3'
			lib_committed: #20140101.1000, useSVCDates:
			lib_committed_expcted: #20140101.1000, modified_updated?: false)

		// test import to existing libary and book records with no committed date,
		// imported record is not committed, useSVCDates = false
		.VerifyExportImport('VerifyExportImportLibraries',
			lib_committed: '', useSVCDates: false
			lib_committed_expcted: '')

		// test import to existing libary and book records with no committed date,
		// imported record is committed, useSVCDates = false
		.VerifyExportImport('VerifyExportImportLibraries2',
			lib_committed: #20140101.1000, useSVCDates: false
			lib_committed_expcted: '')

		// test import to existing libary and book records with committed date,
		// imported record is committed, useSVCDates = false
		.VerifyExportImport('VerifyExportImportLibraries3',
			lib_committed: #20140301.1000, useSVCDates: false
			lib_committed_expcted: #20140101.1000)

		// test import to existing libary and book records with committed date,
		// imported record is committed, useSVCDates = true
		.VerifyExportImport('VerifyExportImportLibraries3',
			lib_committed: #20140301.1000, useSVCDates:,
			lib_committed_expcted: #20140301.1000, modified_updated?: false)

		// test import to existing libary and book records with committed date,
		// imported record is NOT committed, useSVCDates = true
		.VerifyExportImport('VerifyExportImportLibraries3',
			lib_committed: '', useSVCDates:,
			lib_committed_expcted: #20140301.1000, modified_updated?: false)
		}
	}