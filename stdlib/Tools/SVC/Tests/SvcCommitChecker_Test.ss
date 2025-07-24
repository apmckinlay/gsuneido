// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
SvcTests
	{
	Test_last_library_change()
		{
		fn = SvcCommitChecker.SvcCommitChecker_last_library_change
		lastget = Date()
		prevChange = lastget.Plus(minutes: 1)
		changeAsOf = lastget.Plus(minutes: 2)
		exclude_lib = .MakeLibrary(
			[name: 'excluded_test', lib_modified: lastget.Plus(minutes: 4)])
		modifiedTable1 = .MakeLibrary([name: 'modified_test', lib_modified: prevChange])
		modifiedTable2 = .MakeLibrary(
			[name: 'prev_modified', lib_modified: lastget],
			[name: 'modified_test', lib_modified: changeAsOf])
		Suneido.LibraryTablesOverride = [modifiedTable1, modifiedTable2, exclude_lib]
		.SpyOn(SvcTable.Getter_ExcludedTables).Return([exclude_lib])

		// Most recent change was a record modification
		lastLocalChange = [when: lastget, what: 'N/A', how: 'SVC']
		Assert(fn(lastLocalChange) is: changeAsOf)
		Assert(lastLocalChange.what is: modifiedTable2 $ ': modified_test')
		Assert(lastLocalChange.how is: 'modified')

		// Most recent change was a record deletion
		// One of the deleted records is in excluded libraries, should be ignored
		validDeleteTime = lastget.Plus(minutes: 6)
		QueryOutput(modifiedTable1,
			[name: name = .TempName(), lib_modified: validDeleteTime, group: -2])
		QueryOutput(exclude_lib,
			[name: 'ExcludeMe', lib_committed: lastget.Plus(minutes: 5), group: -2])
		Assert(fn(lastLocalChange) is: validDeleteTime)
		Assert(lastLocalChange.what is: modifiedTable1 $ ': ' $ name)
		Assert(lastLocalChange.how is: 'deleted')

		// Most recent change was from getting Version Control changes
		lib_committed = lastget.Plus(minutes: 7)
		.MakeLibraryRecord([name: .TempName(), text: 'text', :lib_committed],
			table: modifiedTable2)
		Assert(fn(lastLocalChange) is: validDeleteTime)
		Assert(lastLocalChange.what is: modifiedTable1 $ ': ' $ name)
		Assert(lastLocalChange.how is: 'deleted')
		}

	Test_libraryChangeMessage()
		{
		fn = SvcCommitChecker.SvcCommitChecker_libraryChangeMessage
		lastTestRun = #20220131.0750
		lastLocalChange = [
			how:	'modified'
			what:	'lib:record'
			when:	#20220131.0800
			]
		message = fn(lastLocalChange, lastTestRun)
		.assertMessage(message, #modified, testRun:)

		lastLocalChange.how = #deleted
		lastTestRun = Date.Begin()
		message = fn(lastLocalChange, lastTestRun)
		.assertMessage(message, #deleted)

		lastLocalChange.how = #SVC
		message = fn(lastLocalChange, lastTestRun)
		.assertMessage(message, #SVC)

		lastTestRun = #20220131.0750
		message = fn(lastLocalChange, lastTestRun)
		.assertMessage(message, #SVC, testRun:)
		}

	assertMessage(message, how, testRun = false)
		{
		message = message.Lines().Filter({ it isnt '' })
		i = 0
		Assert(message[i++]
			is: 'You must run all the tests successfully before sending changes')
		if testRun
			Assert(message[i++] startsWith: 'Last successful test:')
		if how in (#modified, #deleted)
			{
			Assert(message[i++] is: 'Last record ' $ how $ ':')
			Assert(message[i++] is: '\t- lib:record')
			Assert(message[i++] startsWith: '\t- Date:')
			}
		Assert(message isSize: i)
		}

	Test_delete_and_referenced?()
		{
		delAndRefd? = SvcCommitChecker.SvcCommitChecker_delete_and_referenced?
		Suneido.LibraryTablesOverride = Object('Test_lib')

		// not a deletion so warning not applicable
		recName = .TempName()
		Assert(delAndRefd?(recName, '+', 'Test_lib') is: false)

		// record being deleted is referenced but has another definition
		recName = .TempName()
		.MakeLibraryRecord([name: recName, text: '#()'])
		.MakeLibraryRecord([name: 'OtherName', text: 'function() {' $ recName $ '}'])
		Assert(delAndRefd?(recName, '-', 'Test_lib') is: false)

		// record has no definitions left and is still referenced
		recName = .TempName()
		.MakeLibraryRecord([name: 'OtherName2', text: 'function() {' $ recName $ '}'])
		Assert(delAndRefd?(recName, '-', 'Test_lib'))
		Assert(delAndRefd?(recName, '-', 'Contrib') is: false)

		// record has no definitions left and is not referenced
		recName = .TempName()
		Assert(delAndRefd?(recName, '-', 'Test_lib') is: false)

		// record has a definition left and is not referenced
		recName = .TempName()
		.MakeLibraryRecord([name: recName, text: '#()'])
		Assert(delAndRefd?(recName, '-', 'Test_lib') is: false)
		}

	Test_verifyCodeQuality()
		{
		lib = .MakeLibrary()
		rec = Object(name: 'rec', :lib, type: '+', text: 'changes')
		.CommitAdd(svc = .Svc(), .SvcTable(lib), rec.name, '', 'add')
		master = [master: svc.Get(lib, rec.name), missingTestOld?:]

		// getCodeQualityRating runs twice per test. Except for the final two tests
		// [Local star rating, Master star rating]
		.SpyOn(Qc_Main.Qc_Main_getCodeQualityRating).Return(
			10, 10, 	// No quality change
			9, false, 	// No Master quality, compare to min value, above min
			7, false, 	// No Master quality, compare to min value, below min
			9, 	10, 	// Quality dropped one star
			8, 	9, 		// Quality dropped one star
			0, 	9, 		// Dropped to 0 quality
			9, 			// New record, no master rec, compare to min value, above min
			7, 			// New record, no master rec, compare to min value, above min
			)
		method = SvcCommitChecker.SvcCommitChecker_verifyCodeQuality
		Assert(method(rec, rec, master) is: '')
		Assert(method(rec, rec, master) is: '')
		Assert(method(rec, rec, master)
			is: lib $ ':rec rating is: 3.5,  maintain/exceed: 4\n')
		Assert(method(rec, rec, master)
			is: lib $ ':rec rating is: 4.5,  maintain/exceed: 5\n')
		Assert(method(rec, rec, master)
			is: lib $ ':rec rating is: 4,  maintain/exceed: 4.5\n')
		Assert(method(rec, rec, master)
			is: lib $ ':rec rating is: 0,  maintain/exceed: 4.5\n')

		master.master = false
		Assert(method(rec, rec, master) is: '')
		Assert(method(rec, rec, master)
			is: '+' $ lib $':rec rating is: 3.5,  min: 4\n')
		}

	Test_checkLineEnds()
		{
		scc = SvcCommitChecker
			{
			SvcCommitChecker_maxAllowedRecs: 3
			SvcCommitChecker_libraryCheck?(model/*unused*/, table/*unused*/)
				{ return true }
			}
		// only valid record
		changes = Record(
			[lib:'Test_lib', name: .TempName(), type: '+', text: '#()'],
			[lib:'Test_lib', name: .TempName(), type: '+', text: '#("HelloWorld!")'],
			[lib:'Test_lib', name: .TempName(), type: '+', text: '#(101)'])
		for change in changes
			.MakeLibraryRecord(change)
		Assert(scc.SvcCommitChecker_checkLineEnds('', changes, '') is: '')

		// invalid records
		name1 = .TempName()
		name2 = .TempName()
		changes = Record(
			[lib:'Test_lib', name: name1, type: '+', text: '#(101)\n'],
			[lib:'Test_lib', name: name2, type: '', text: '#(102)\n'])
		for change in changes
			.MakeLibraryRecord(change)
		msg = "Following record(s) uses non-standard line ending characters:\n\t" $
			"Test_lib:" $ name1 $ "\n\tTest_lib:" $ name2
		Assert(scc.SvcCommitChecker_checkLineEnds('', changes, '') is: msg)

		// should show only maxAllowedRecs invalid
		name1 = .TempName()
		name2 = .TempName()
		name3 = .TempName()
		name4 = .TempName()
		changes = Record(
			[lib:'Test_lib', name: .TempName(), type: '+', text: '#(good)'],
			[lib:'Test_lib', name: name1, type: '+', text: '#(101)\n'],
			[lib:'Test_lib', name: name2, type: '', text: '#(102)\n'],
			[lib:'Test_lib', name: name3, type: '', text: '#(103)\n'],
			[lib:'Test_lib', name: name4, type: '', text: '#(104)\n'],)
		for change in changes
			.MakeLibraryRecord(change)
		msg = "Following record(s) uses non-standard line ending characters:\n\t" $
			"Test_lib:" $ name1 $ "\n\tTest_lib:" $ name2 $ "\n\tTest_lib:" $ name3 $
			"\n\tToo many record to display"
		Assert(scc.SvcCommitChecker_checkLineEnds('', changes, '') is: msg)
		// invalid record getting removed
		name1 = .TempName()
		changes = Record(
			[lib:'Test_lib', name: name1, type: '-', text: '#()\n'])
		for change in changes
			.MakeLibraryRecord(change)
		Assert(scc.SvcCommitChecker_checkLineEnds('', changes, '') is: '')
		}

	Test_errors_in_local_changes?()
		{
		cl = SvcCommitChecker
			{
			SvcCommitChecker_checkRecord(change)
				{
				errors = Object()
				if '' isnt error = change.GetDefault('error', '')
					errors.Add(error)
				warnings = Object()
				if '' isnt warning = change.GetDefault('warnings', '')
					warnings.Add(warning)
				return Object(:errors , :warnings)
				}
			}
		fn = cl.SvcCommitChecker_errors_in_local_changes?

		table = lib = .TestLibName()
		deleted = [name: 'Deleted_NotChecked', type: '-', :lib]
		Assert(fn([deleted], table) is: '')

		valid = [name: 'Valid', type: ' ', :lib]
		Assert(fn([deleted, valid], table) is: '')

		suppressed = [name: 'Suppressed', type: ' ', :lib, error: 'suppressed']
		.assertCodeErrorMessage(fn, [suppressed], table,
			[suppressedStr = .errorString(suppressed)])

		invalid = [name: 'Invalid', type: ' ', :lib, error: 'invalid']
		.assertCodeErrorMessage(fn, [invalid], table,
			[invalidStr = .errorString(invalid)])

		otherErrors = [name: 'OtherErrors', type: '+', :lib,
			error: 'error message from the record that will be displayed on the ' $
				'pop up but what to test the ellipsis here 116 characters']
		.assertCodeErrorMessage(fn, [otherErrors], table,
			[otherIssuesStr = .errorString(otherErrors,
				'error message from the record that will be displayed on the ' $
				'pop up but what to test the ellipsis her...')])

		changes = [valid, suppressed, deleted, invalid, otherErrors]
		expectedErrors = [suppressedStr, invalidStr, otherIssuesStr]
		.assertCodeErrorMessage(fn, changes, table, expectedErrors)
		}

	assertCodeErrorMessage(fn, changes, table, expectedErrors)
		{
		results = fn(changes, table).Split('\r\n\t- ')
		Assert(results isSize: expectedErrors.Size() + 1)
		Assert(results[0] is: 'Unable to send selected changes:')
		expectedErrors.Each({ Assert(results has: it) })
		}

	errorString(change, error = false)
		{
		if error is false
			error = change.error
		return change.lib $ ':' $ change.name $ ' (' $ error $ ')'
		}

	Test_formatMsg()
		{
		m = SvcCommitChecker.SvcCommitChecker_formatMsg

		checkRecord = Object(errors: Object(), warnings: Object())
		Assert(m(false, '', checkRecord) is: '')

		Assert(m(true, '', checkRecord) is: 'Record is over size limit (100 kb)')

		msg = m(true, 'Quality Checker Text', checkRecord)
		Assert(msg has: 'Record is over size limit (100 kb)')
		Assert(msg has: 'Quality Checker Text')

		checkRecord.errors.Add('error 1', 'error 2')
		Assert(m(false, '', checkRecord)
			is: 'Record has syntax error(s):\r\n- error 1\r\n- error 2')

		checkRecord.warnings.Add('warn 1', 'warn 2')
		msg = m(false, '', checkRecord)
		Assert(msg has: 'Record has syntax error(s):\r\n- error 1\r\n- error 2')
		Assert(msg has: 'Record has the following warning(s):\r\n- warn 1\r\n- warn 2')

		msg = m(true, 'Quality Checker Text', checkRecord)
		Assert(msg has: 'Record is over size limit (100 kb)')
		Assert(msg has: 'Quality Checker Text')
		Assert(msg has: 'Record has syntax error(s):\r\n- error 1\r\n- error 2')
		Assert(msg has: 'Record has the following warning(s):\r\n- warn 1\r\n- warn 2')
		}

	Teardown()
		{
		Suneido.Delete(#LibraryTablesOverride)
		super.Teardown()
		}
	}
