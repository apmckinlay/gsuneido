// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	title: "Onsite Continuous Tests"
	Bucket: 'continuoustests'
	Run(testType, bookCheckOnly = false)
		{
		Database.SessionId('(continuous_tests)')
		testGroup = Sys.Client?() ? 'Client/Server' : 'Standalone'
		if false is testInfoOb = .getTestInfo(testType)
			{
			.SendResults(.title $ " - FAILED invalid Go File - Invalid testType",
				Display(testType))
			return
			}

		resultFileName = 'stillRunningResults.txt'
		if not FileExists?(resultFileName)
			.putFile(resultFileName,
				testGroup $ " Tests Started at " $ Display(date = Date()) $
				' (' $ date.UTC() $ ')\r\n')

		skipTags = Object('win32', 'installedSystem')
		if not Sys.Windows?()
			skipTags.Add('windows')

		if testInfoOb.Member?('additionalSkip')
			skipTags.Add(@testInfoOb.additionalSkip)

		AddFile(resultFileName, "FAILED - " $ testInfoOb.type $
			" started but did not finish\r\n")
		results = .runTests(testInfoOb, skipTags, bookCheckOnly)

		.addResultsToFile(testInfoOb.type, results, resultFileName)
		}

	getTestInfo(testType)
		{
		if testType is 'ServerContinuousTestInfo'
			return ServerSuneido.Get('ServerContinuousTestInfo')
		continuousTestTypes = GetContributions('ContinuousTestTypes')
		return continuousTestTypes.FindOne({ it.type is testType })
		}

	runTests(testInfoOb, skipTags, bookCheckOnly)
		{
		.DisableSuneidoVariables(testInfoOb)

		results = ''
		try
			{
			// use tables - only stdlib, axonlib used at this point
			// loop through testInfoOb.libs and use them.
			if not Sys.Client?()
				{
				if Libraries().Has?('configlib')
					Unuse('configlib')
				for lib in testInfoOb.libs
					Use(lib)
				LibraryTags.Reset()
				}

			additionalServerTests? =
				testInfoOb.GetDefault('testGroup', '') is 'Client/Server'

			results = ContinuousTests(noData?: testInfoOb.GetDefault('noData?', false),
				currency: testInfoOb.GetDefault('currency', 'USD'),
				:skipTags,
				dept?: testInfoOb.GetDefault('dept?', true),
				:additionalServerTests?, :bookCheckOnly)
			}
		catch (err)
			results = 'ERROR: ' $ err

		results $= '\r\n' $ ContinuousTests_ErrorLogs()
		if not bookCheckOnly
			results = .checkLibs(testInfoOb, results)
		return .CheckBook(results, testInfoOb)
		}

	DisableSuneidoVariables(params)
		{
		if params.GetDefault('disableCompanyTableChecking?', false) is true
			Suneido.DisableCompanyTableChecking = ""
		}

	checkLibs(testInfoOb, results)
		{
		if testInfoOb.GetDefault('skipLibCheck?', false) is false
			{
			libs = Libraries()
			if testInfoOb.Member?('additionalSkip')
				{
				libs = libs.Difference(testInfoOb.additionalSkip)
				results $= '\r\nChecking Libraries: ' $ libs.Join(',')
				}
			results $= '\r\n'
			try
				results $= Continuous_CheckLibraries(libs)
			catch (err)
				results $= err
			}
		return results
		}

	CheckBook(results, testInfoOb)
		{
		if testInfoOb.GetDefault('skipBookCheck?', false) is false
			try
				{
				Plugins().ForeachContribution('ContinuousTests', 'bookCheck', showErrors:)
					{ |x|
					results $= (x.bookCheck)()
					}
				}
			catch (err)
				results $= '\r\n\r\n' $ err

		return results
		}

	addResultsToFile(testType, results, resultFileName)
		{
		prevResults = .getFile(resultFileName)
		prevResults = prevResults is false ? ''
			: prevResults.Replace("FAILED - " $ testType $
				" started but did not finish", "")

		.putFile(resultFileName, prevResults $ .FormatResultContent(results, testType))
		}

	ErrorSpan: `<span class='error'>ERROR</span>`
	WarnSpan: `<span class='warning'>WARNING</span>`

	FormatResultContent(results, testSetName)
		{
		status = ''
		if results =~ .ErrorRegex
			status = .ErrorSpan
		else if results.Has?("WARNING")
			status = .WarnSpan

		summary = Xml('summary', status $ .Padding $ testSetName $ .Padding)
		text = Xml('p', results.Trim())
		return Xml('details', summary $ text )
		}

	// called from local continuous test script
	SendEmail(testGroup, testType, newExes? = false)
		{
		testGroupLabel = testGroup $ (newExes? ? ' Latest Exe' : '')
		for subfolder in #('timeout_gsport', 'timeout_gdev',
			'sujsweb_gdev', 'sujsweb_gsport')
			{
			if DirExists?('../' $ subfolder)
				{
				extra = ''
				if false is result =
					.getFile('../' $ subfolder $ '/stillRunningResults.txt')
					{
					extra = .snapshotTimeoutTesterFolder(subfolder)
					result =	'ERROR: Missing Timeout/SuJsWeb Tester result from ' $
						subfolder $ extra
					}
				else if result =~ .ErrorRegex
					extra = .snapshotTimeoutTesterFolder(subfolder)
				AddFile('stillRunningResults.txt', '\r\n' $ result $ extra)
				}
			}

		AddFile('stillRunningResults.txt',
			'\r\n' $ testGroupLabel $ ' Tests Ended at ' $ Display(Date()))
		fileName = .buildFileName(testType, newExes?)
		fileOb = .buildFormatedResultsFile(.getFile('stillRunningResults.txt'),
			fileName, testGroupLabel)

		message = ''
		Plugins().ForeachContribution('ContinuousTests', 'buildMessage')
			{ |x|
			message $= (x.buildMessage)(.title, fileOb.status, fileOb.filename)
			}

		.deleteFile('stillRunningResults.txt')
		.SendResults(.title $ ' - ' $ fileOb.title, message)
		}

	buildFileName(testType, newExes?)
		{
		suffix = '_' $ testType
		if newExes?
			suffix $= '_latestExe'
		return 'results' $ suffix $ '.txt'
		}

	snapshotTimeoutTesterFolder(timeoutFolder)
		{
		failedDir = timeoutFolder $ '_failedbak'
		if DirExists?(failedDir)
			DeleteDir(failedDir)
		EnsureDir(failedDir)
		try
			failed = MirrorDir('../' $ timeoutFolder, failedDir)
		catch (err)
			failed = Object(err)
		return '\r\nThe ' $ timeoutFolder $ ' is copied to ' $
			GetCurrentDirectory() $ '/' $ failedDir $
			Opt(', (Except ', failed.Join(','), ')')
		}

	buildFormatedResultsFile(message, filename, testGroupLabel)
		{
		results = .htmlFormatResults(message)
		results.filename = filename.Replace('.txt$', '') $
			(results.status is 'FAILED' ? Display(Date()).Replace('#', '') : '') $ '.html'
		results.title = testGroupLabel $ ' ' $ results.status
		headTag = '<html><head><meta charset="ascii"><title>' $
			XmlEntityEncode(results.title) $ '</title>' $ .Style() $ '</head>'
		results.text = headTag $ '<body><pre>' $ results.text $ '</pre></body></html>'
		.putFile(results.filename, results.text)
		return results
		}

	Style()
		{
		error = `.error { color: red; margin: 0 5px; }`
		warning = `.warning {margin: 0 5px; color: orange; }`
		summary = `summary { display: flex; align-items: center;}` $
			`summary::before { content: "\25B8"; font-size: 20px; }` $
			`summary:hover { cursor: pointer; }`
		return Xml('style', summary $ error $ warning )
		}

	minTestResultSize: 8
	ErrorRegex: "\<(ERROR|FAIL|FATAL)"
	Padding: "===================="
	htmlFormatResults(message, heading = false)
		{
		text = heading isnt false
			? "<h3>" $ heading $ "</h3>"
			: "<h3>Continuous Tests Results</h3>"
		text $= Opt("<pre>",.GetChangeFileText().Trim(),"</pre>")
		status = 'PASSED'
		messageOb = message.Lines()
		for line in messageOb
			{
			if .CheckForError(line)
				{
				line = '<span style="color: red">' $ line $ '</span>'
				status = 'FAILED'
				}
			else if .CheckForWarning(line)
				{
				line = '<span style="color: orange">' $ line $ '</span>'
				if status isnt 'FAILED'
					status = 'WARNING'
				}
			text $= line $ '\r\n'
			}
		if messageOb.RemoveIf(#Blank?).Size() < .minTestResultSize
			{
			text $= '<span style="color: red">Missing too much content, something went ' $
				'wrong with the tests</span>'
			status = 'FAILED'
			}
		return Object(:status, :text)
		}

	CheckForError(line)
		{
		return line =~ .ErrorRegex and line !~ .Padding
		}

	CheckForWarning(line)
		{
		return line.Has?("WARNING") and line !~ .Padding
		}

	GetChangeFileText()
		{
		if false is change = .getFile(`../` $ .ChangesResultFile)
			change = ''
		else
			change = 'Continuous Triggered by ' $ change
		return change
		}

	SendResults(subject, message)
		{
		if false is send = OptContribution('ContinuousTestSendResults', false)
			return
		groups = Object('programmers')
		if subject.Suffix?('FAILED') or message.Has?('ERROR:')
			groups.Add('failed_continuous_tests')
		send(groups, subject, message)
		}

	checkForLatestExe()
		{
		compareResults = .compareExes()
		errors = compareResults.errors
		copyMap = compareResults.copyMap
		for mem in copyMap.Copy().Members()
			.copyExes(mem, copyMap[mem], errors, copyMap)

		compareErrors = Object()
		errors.Each({ compareErrors.Add('\t' $ it.Join('\r\n\t')) })
		if compareErrors.RemoveIf({ it.Blank?() }).NotEmpty?()
			.SendResults(.continuousTestType() $ ' - Check Latest Exes FAILED',
				'Errors Encountered:\r\n' $ compareErrors.Join('\t\r\n'))

		return copyMap.Any?({ it.NotEmpty?() })
			? copyMap.Map!({ it.message }).Values().Join('\r\n\r\n')
			: false
		}

	/* Expected Format
	result = [
		copyMap: [
			version1: [
				serverExeStd:  'standard server exe name', 	// IE: gsport.exe
				serverExeTest: 'test server exe name',		// IE: gsport_version1.exe
				clientExeStd:  'standard client exe name',	// IE: gsuneido.exe
				clientExeTest: 'test client exe name',		// IE: gsuneido_version1.exe
				stdExePath:    'path to the standard exes',	// IE: C:/path/<serverExeStd>
				message:       'string sent with test results' // IE: "Copied <x>, to <y>"
				]
			version2: [
				...
				// same structure as exeVersion1, but with different values
				... ]
			]
		errors: [
			version1: #()
			version2: #()
			// If there are errors copying with a specific map, it will be removed from
			// the copyMap member, and it's errors will be included in this object
			version3: #('Example: Failed to get file date for comparison, (<file>)')
			]
		]
	*/
	compareExes()
		{
		return (OptContribution('ContinuousTest_CompareExes',
			{ return [copyMap: #(), errors: #()] }))()
		}

	copyExes(exeVersion, exeMap, errors, copyMap)
		{
		#(client, server).Each()
			{|exeType|
			path = Paths.Combine(exeMap.stdExePath, exeMap[exeType $ 'ExeStd'])
			try
				.copyFile(path, './' $ exeMap[exeType $ 'ExeTest'])
			catch (error)
				{
				errors[exeVersion].
					Add('Failed to copy exe: ' $ path $ ', (' $ error $ ')')
				copyMap.Delete(exeVersion)
				}
			}
		}

	//Extracted for tests
	copyFile(from, to)
		{
		// need to use a System command to copy the file to keep the permissions set
		// correctly. Using CopyFile changes the permissions on the exe files which will
		// make the copy fail the next time
		from = `"` $ Paths.ToLocal(from) $ `"`
		to = `"` $ Paths.ToLocal(to) $ `"`
		copyCmd = Sys.Windows?() ? 'copy /Y ' : 'cp -p '
		if 0 isnt result = System(copyCmd $ from $ ' ' $ to)
			throw 'system command failed with result: ' $ result
		}

	ChangesResultFile: 'getChangesSucceeded.txt'
	ForceRunFile: 'forceRun.txt'
	importTable: 'continuous_test_import_history'
	importDelimiter: '~~~'
	GetVersionControlChanges(forcedRun? = false)
		{
		.EnsureImportHistoryTable()
		EnsureDir('exeFolder')
		.deleteFile(.ChangesResultFile)
		instance = .continuousTestType()
		lastImport = ''
		try
			{
			forcedRun? = .checkForceRunFile(forcedRun?)
			latestSuneido = .checkForLatestExe()
			if false is changes = .getMasterChanges()
				return

			if changes.Has?(.importDelimiter)
				{
				lastImport = changes.AfterFirst(.importDelimiter)
				changes = changes.BeforeFirst(.importDelimiter)
				}

			if changes is '' and latestSuneido is false and not forcedRun?
				return
			}
		catch (err)
			{
			try
				.SendResults(instance $
					" - Import Version Control Changes FAILED", 'Instance: ' $
					instance $ '\r\n' $ err)
			return
			}

		.updateImportHistory(lastImport)

		.putFile(.ChangesResultFile, .buildMsg(changes, latestSuneido, forcedRun?))
		}

	//exported for tests
	updateImportHistory(lastImport)
		{
		lastImportTime = lastImport is ''
			? Display(Timestamp())
			: '#' $ lastImport
		QueryDo('update ' $ .importTable $
			' set ct_import_time = ' $ lastImportTime)
		}

	EnsureImportHistoryTable()
		{
		if false is TableExists?(.importTable)
			{
			Database('ensure ' $ .importTable $ '  (ct_import_time) key ()')
			QueryOutput(.importTable, [ct_import_time: Timestamp().Minus(hours:12)])
			}
		}

	checkForceRunFile(forcedRun?)
		{
		if not forcedRun?
			forcedRun? = FileExists?(.ForceRunFile)
		.deleteFile(.ForceRunFile)
		return forcedRun?
		}

	exportDir: 'SvcChangesExported'
	getMasterChanges()
		{
		EnsureDir(.exportDir)

		lastDate = QueryLast(.importTable $ ' sort ct_import_time').ct_import_time

		if lastDate < Date().Minus(days: 14)
			{
			.SendResults('Out of date test instance', 'Last recorded import of code ' $
				'was older than 14 days.  The following instance will need to be ' $
				'manually fixed.\r\nFrom: ' $ GetCurrentDirectory())
			return false
			}

		changes = ''
		checksumOb = Object()
		fileList = .getFileList(lastDate)
		for file in fileList
			{
			.getZipFile(file)

			Spawn(P.WAIT, 'unzip', '-o', 'SvcExports.zip', '-d', .exportDir)

			changes $= .getFile(.exportDir $ '/Changes.txt')
			.deleteFile(.exportDir $ '/Changes.txt')

			checksumOb = .collectChecksums()
			for file in Dir(.exportDir $ '/*')
				{
				if not TableExists?(file)
					LibTreeModel.Create(file)
				LibIO.Import(.exportDir $ '/' $ file, file, useSVCDates:)
				.deleteFile(ExeDir() $ '/' $ .exportDir $ '/' $ file)
				}
			}

		if checksumOb.NotEmpty?() and '' isnt msg = .compareChecksums(checksumOb)
			{
			.SendResults(.continuousTestType() $ " - Verify Library Checksums FAILED",
				'From: ' $ GetCurrentDirectory() $ '\r\n' $ msg)
			}

		if not fileList.Empty?()
			changes = changes $ .importDelimiter $ fileList.Last().Extract('\d+\.?\d+')

		return changes
		}

	continuousTestType()
		{
		return Paths.Basename(GetCurrentDirectory())
		}

	getFileList(lastDate)
		{
		fileList = Dir(.svcExportPath() $ '/*')
		fileList = fileList.Filter(
			{
			date = it.Extract('\d+\.?\d+')
			Date(lastDate) < Date(date)
			}).Sort!()
		return fileList
		}

	svcExportPath()
		{
		return OptContribution('ContinuousTest_LocalExportPath', '')
		}

	getZipFile(file)
		{
		from = Paths.Combine(.svcExportPath(), file)
		to = Paths.Combine(ExeDir(), 'SvcExports.zip')
		if FileExists?(to)
			.deleteFile(to)
		Retry(maxRetries: 3)
			{
			CopyFile(from, to, false)
			}
		}

	buildMsg(masterChanges, latestSuneido, forcedRun?)
		{
		msgs = Object()
		if latestSuneido isnt false
			msgs.Add(latestSuneido)
		if masterChanges isnt ''
			msgs.Add('Version Control Changes:\r\n' $ masterChanges)
		else if forcedRun?
			msgs.Add('Forcing test run - No Changes found')

		return msgs.Join('\r\n\r\n')
		}

	collectChecksums()
		{
		if false is s = .getFile(.exportDir $ '/checksums.json')
			return #()

		.deleteFile(.exportDir $ '/checksums.json')
		return Json.Decode(s)
		}

	compareChecksums(checksumOb)
		{
		libsToCheck = .libsToCheck()
		localChecksums = .checksumLibraries(libsToCheck, books: #())
		msgOb = Object()
		compareOb = Object()
		for lib in libsToCheck
			{
			localSum = .findLibChecksum(localChecksums, lib)
			svcSum = .findLibChecksum(checksumOb, lib)
			if localSum isnt svcSum
				compareOb.Add(Object(local: localSum.cksum, svc: svcSum.cksum, :lib))
			}
		libFill = compareOb.NotEmpty?()
			? compareOb.MaxWith({ it.lib.Size() }).lib.Size()
			: 0
		chksumLen = 10
		compareOb.Each()
			{
			msgOb.Add(it.lib.RightFill(libFill, ' ') $
				' - svc ' $ it.svc.LeftFill(chksumLen, ' ')  $
				'; local ' $ it.local)
			}
		return msgOb.Empty?()
			? ''
			: 'Checksum mismatches:\r\n\t' $ msgOb.Join('\r\n\t')
		}

	libsToCheck()
		{
		libs = Object()
		continuousTestTypes = GetContributions('ContinuousTestTypes')
		for testSet in continuousTestTypes
			libs.MergeUnion(testSet.libs)
		libs.Remove('configlib').Add('stdlib', at: 0)
		return libs
		}

	findLibChecksum(cksumOb, lib)
		{
		cksum = cksumOb.FindOne({ it.lib is lib })
		if cksum is false
			return Object(:lib, cksum: 'missing')
		return cksum
		}

	// methods needed for tests
	putFile(file, text)
		{
		PutFile(file, text)
		}
	getFile(filename)
		{
		return GetFile(filename)
		}
	deleteFile(file)
		{
		DeleteFile(file)
		}
	checksumLibraries(libs = false, record_results = false, books = false)
		{
		return ChecksumLibraries(libs, record_results, books)
		}
	}
