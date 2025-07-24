// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
LibIOTestBase
	{
	Test_import_record_with_pending_delete()
		{
		.testPendingDeletes(
			exportCommitted:	#20131212.1000,
			expectedCommitted: 	#20131111.1000,
			useSVCDates: 		false)
		.testPendingDeletes(
			exportCommitted: 	#20131212.1000,
			expectedCommitted: 	#20131212.1000,
			useSVCDates: 		true)
		.testPendingDeletes(
			exportCommitted: 	'',
			expectedCommitted: 	#20131111.1000,
			useSVCDates: 		false)
		.testPendingDeletes(
			exportCommitted: 	'',
			expectedCommitted: 	'',
			useSVCDates: 		true)
		}

	testPendingDeletes(exportCommitted, expectedCommitted, useSVCDates)
		{
		name =  .TempTableName() // unique record name each time
		modified_updated? = not useSVCDates
		rec = [lib_committed: expectedCommitted, :name]
		.SvcLibrary.Output(rec.Copy(), deleted:, committed:)
		.SvcBook.Output(rec.Copy(), deleted:, committed:)

		Assert(.SvcLibrary.Get(name, deleted:) isnt: false)
		Assert(.SvcBook.Get('/' $ name, deleted:) isnt: false)

		.VerifyExportImport(name,
			lib_committed: exportCommitted, :useSVCDates,
			lib_committed_expcted: expectedCommitted, :modified_updated?)

		Assert(.SvcLibrary.Get(name, deleted:) is: false)
		Assert(.SvcBook.Get('/' $ name, deleted:) is: false)
		}
	}
