// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	errorLogByteLimit: 15000
	GetLog(files, logHeading = false, archiveDirectory = false, archiveFunction = false)
		{
		errors = ''
		for file in files.Members()
			errors $= .getLogFile(files[file], file, logHeading)

		.ArchiveLogFiles(files, archiveDirectory, archiveFunction)
		return errors
		}
	GetLogAsObject(files, archiveDirectory = false, archiveFunction = false)
		{
		retOb = Object()
		for file in files.Members()
			retOb[file] = .getLogFile(files[file], file)
		.ArchiveLogFiles(files, archiveDirectory, archiveFunction)
		return retOb
		}

	getLogFile(file, logName, logHeading = false)
		{
		if false is errOb = .getFileWithLimit(file)
			return ''
		errors = errOb.str
		if errors.Blank?()
			return ''

		if errors.Size() >= .errorLogByteLimit
			{
			errors $= '...\nTOO MANY ERRORS: LOG HAS BEEN TRIMMED'
			if .truncatedErrors?(file)
				errors $= '\nERROR/WARNING found in trimmed section\n'
			}
		if errOb.suppressedErrors > 0
			errors $= '\nSUPPRESSIONS: ' $ errOb.suppressedErrors $
				' Non-Suneido errors were suppressed.\n'
		heading = ''
		if logHeading isnt false
			heading $= '\n\n' $ logName $ logHeading $ ':\n'
		return heading $ errors
		}

	// broken out for test
	getFileWithLimit(file)
		{
		try
			{
			File(file)
				{|f|
				result = .readFileWithSuppressions(f, .errorLogByteLimit)
				}
			}
		catch
			return false

		return result
		}

	readFileWithSuppressions(f, errorLogByteLimit)
		{
		str = ''
		suppressedErrors = 0
		while false isnt line = f.Readline()
			{
			if .SuppressLine?(line)
				{
				suppressedErrors++
				continue
				}
			if .skippedLine?(line)
				continue
			str $= line $ '\r\n'
			if str.Size() >= errorLogByteLimit
				break
			}
		return Object(:suppressedErrors, :str)
		}

	// These are for suppressing / ignoring non-suneido errors ONLY
	// Do not add to this list unless suppressing errors from other unrelated applications
	suppressionRegexes: (
		` \d\d\ds\d\d\dms\d\d\dus `, // ' 293s572ms754us '
		` PREV(IOUS)?: \(\d+\) `, // ' (11) '
		` PREV(IOUS)?:  ?- `, // ' PREVIOUS  - ' (usually PCmiler errors)
		`(i?)SQLite`,
		`ERROR crashpad_client_win`, // google drive crashpad error
		`Failed to start crash handler process`,
		`CRASHPAD_HANDLE_START_ERROR`,
		`GoogleDriveFSPipe`,
		`GetProtoFromRegistryValue Opening registry key`,
		`ctxmenu.cc:213:GenerateContextMenu`,
		`ALK Technologies\\PCMILER30`, // pcmiler errors
		`Code: 8, Message: 'attempt to write a readonly database'`,
		`Code: 14, Message: 'unable to open database file'`,
		`ERROR in rule.*SHOW`,
		`ERROR: in rule.*SHOW`,
		// WebView2 error when exit due to lost connection
		`Failed to unregister class Chrome_WidgetWin_0`
		// not needed after BuiltDate 2025-08-21
		`WARNING: Query1 slow:.*views`
		)
	SuppressLine?(line)
		{
		return .suppressionRegexes.Any?({ line =~ it })
		}

	// these are axon logs that we still want in the file but do not need to see on the
	// main nightly checks page like user connectin timeouts
	skippedRegexes: (
		`dbms server.*closing idle connection`,
		`PREV.*FATAL: lost connection: EOF`
		)
	skippedLine?(line)
		{
		return .skippedRegexes.Any?({ line =~ it })
		}

	truncatedErrors?(file)
		{
		File(file)
			{ |f|
			// Need to handle if the ByteLimit ended up in the middle of the words ERROR
			// or WARNING - adjust backwards by the size of WARNING (bigger of the two)
			return .searchForError(f, .errorLogByteLimit - "WARNING".Size())
			}
		}

	// broken out for test
	searchForError(f, startbyte)
		{
		f.Seek(startbyte)
		while false isnt str = f.Readline()
			{
			if str =~ '\<ERROR|WARNING\>' and not .SuppressLine?(str)
				return true
			}
		return false
		}

	ArchiveLogFiles(files, archiveDirectory, archiveFunction)
		{
		if archiveDirectory is false or archiveFunction is false
			return
		EnsureDir(archiveDirectory)

		suffix = Display(Timestamp()).Tr('#.')
		for file in files
			{
			archiveFunction(file, archiveDirectory $ '/',
				Paths.Basename(file).Replace('.log', '_server' $ suffix $ '.log'),
				Paths.Basename(file).Replace('.log', '_server*.log'))
			}
		}
	}
