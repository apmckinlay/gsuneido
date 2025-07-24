// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	errorMsg: 'Please ensure the scanner is turned on and connected.\r\n'
	// returnScanOb? for debugging purposes
	Scan(filename, returnScanOb? = false)
		{
		BookLog('Scan Attachment Start - Clscan')
		if false is .handleFile(filename)
			return 'Failed to overwrite existing file'

		cmd = .handleSettings(filename)
		if cmd.Has?('User Cancelled') or cmd.Has?('Problem')
			return cmd

		scanResults = .execute(cmd)
		if returnScanOb?
			return .makeOptionOb(.formatResults(scanResults))

		return .handleScanResults(filename, scanResults)
		}

	handleFile(filename)
		{
		// Previously confirmed override from user
		if false is FileExists?(filename)
			return true
		return DeleteFile(filename)
		}

	handleSettings(filename)
		{
		baseCmd = '/SetFileName "' $ filename $ '" '
		scanSettings = OptContribution('GetClscanScannerSettings',
			function () { return #() })()

		if false is scanner = .checkAndSetScanner(UserSettings.Get(
			'Clscan - Scanner Source'))
			return 'User Cancelled'

		if not scanner.Prefix?('WIA')
			baseCmd $= '/SetPageSize "USLETTER" '

		return .checkSettings(baseCmd, scanner, scanSettings) $ scanner
		}

	checkAndSetScanner(source)
		{
		// uses default scanner
		if source is false
			return ''

		if false is .scannerList.Has?(source)
			{
			if false is YesNo('Scanner saved is not detected.  If you want to continue ' $
				'it will use the default scanner.  Do you want to continue?',
				'Scanner Source')
				return false
			return ''
			}
		return ' /SetScanner "' $ source $ '"'
		}

	checkSettings(cmd, scanner, scanSettings)
		{
		checkCmd = scanner
		scanCmd = cmd
		for option in scanSettings.Copy().Delete('setupLocation')
			{
			if option.check is ''
				continue
			checkCmd $= ' ' $ option.check
			scanCmd $= ' ' $ .FormatArgs(option.scan, option.value)
			}
		if checkCmd is scanner
			return cmd

		if UserSettings.Get('Clscan - Scanner Settings', #()) isnt scanSettings
			{
			if false is settingsResult = .getScannerSettingsOb(.execute(checkCmd))
				return 'Problem detecting scanner'

			if '' isnt msg = .compareSettings(settingsResult, scanSettings)
				{
				// SuneidoLog to help us find different responses from scanners
				SuneidoLog('INFO: Company settings do not match scanners',
					params: settingsResult)
				if false is YesNo(msg $ '\r\nContinue?', 'Scan Attachment')
					return 'User Cancelled'
				}
			UserSettings.Put('Clscan - Scanner Settings', scanSettings)
			}
		return scanCmd
		}

	compareSettings(result, scanSettings)
		{
		msg = ''
		if result.Members().FindIf({ it.Prefix?( 'Untested Windows') }) isnt false
			return msg // attempt to scan to wake scanner
		setupLocation = scanSettings.GetDefault('setupLocation', '')
		if scanSettings.Member?('duplex') and scanSettings.duplex.value is 'Y'
			msg $= .checkDuplex(result.GetDefault('duplex', false))
		if scanSettings.Member?('resolutions')
			msg $= .checkResolution(result.GetDefault('resolutions', #()),
				scanSettings.resolutions.value, setupLocation)
		if scanSettings.Member?('color')
			msg $= .checkColor(result.GetDefault('color', #()), scanSettings.color.value,
				setupLocation)

		return msg
		}

	checkDuplex(duplex)
		{
		return duplex is false
			? 'Scanner does not support changing double-sided options.\r\n' $
				'Depending on scanner settings, both or one side may be ' $
				'scanned.\r\n'
			: ''
		}

	checkResolution(resolution, companySetResolution, setupLocation)
		{
		minMaxResult? = resolution.Member?(0) and resolution[0].Lower().Has?('min')
			? .checkMinMaxResolution(resolution[0].Lower(), companySetResolution)
			: false

		return (resolution.Has?(String(companySetResolution)) or minMaxResult?)
			? ''
			: .notSupportedMsg('resolution', setupLocation)
		}

	checkMinMaxResolution(resolution, companySetResolution)
		{
		min = Number(resolution.Extract('min value=(\d*)'))
		max = Number(resolution.Extract('max value=(\d*)'))
		step = Number(resolution.Extract('step=(\d*)'))
		if companySetResolution < min or max < companySetResolution
			return false

		base = companySetResolution - min
		return base % step is 0
		}

	notSupportedMsg(type, setupLocation)
		{
		return 'The ' $ type $ ' specified' $  Opt(' in ', setupLocation) $
			' is not supported by your scanner.\r\n'
		}

	checkColor(color, companySetColor, setupLocation)
		{
		return false is color.Has?(companySetColor)
			? .notSupportedMsg('color', setupLocation)
			: ''
		}

	getScannerSettingsOb(scannersResults)
		{
		type = ''
		allowedSettings = Object()
		if scannersResults.Prefix?('Problem while opening')
			return false
		resultLines = scannersResults.Lines()
		startIdx = resultLines.FindIf({ it.Has?('Selected scanner') })
		for line in resultLines[startIdx + 1 ..] //skip scanner line
			{
			if line.Blank?()
				continue
			if line.Has?(":")
				{
				type = line.AfterFirst(' ').BeforeFirst(':')
				if type isnt 'duplex'
					allowedSettings[type] = Object()
				}
			else
				{
				if type is ''
					continue
				if type isnt 'duplex'
					allowedSettings[type].Add(line)
				else
					allowedSettings[type] = line is 'True'
				}
			}
		return allowedSettings
		}


	makeOptionOb(results)
		{
		scanOb = Object()
		resultPos = ''
		for item in results
			{
			if item.Blank?()
				{
				resultPos = results.FindIf({ it is item})
				break
				}
			mem = item.BeforeFirst(' =')
			val = item.AfterFirst('= ')
			scanOb[mem] = val
			}
		scanOb.results = results[resultPos+1 ..].Join(' ')
		return scanOb
		}

	formatResults(results, parsingMessage = ':\r\n')
		{
		return results.AfterLast(parsingMessage).Split('\r\n')
		}

	execute(cmd)
		{
		if false is prog = ExternalApp('clscan')
			return ''
		try result = RunPipedOutput('"' $ prog $ '" ' $ cmd)
		catch (err)
			return '\nProblem while opening the scanner: ' $ err
			// need \n for handleScanResults parsing
		return result
		}

	FormatArgs(arg, value)
		{
		return Opt(arg, ' ', '"', value, '"')
		}

	GetAvailableSources()
		{
		sources = .execute('/GetScanners')
		return .formatResults(sources)
		}

	handleScanResults(filename, scanResults)
		{
		selPrefix = 'Selected scanner '
		if scanResults.Has?(selPrefix)
			{
			scanner = scanResults.AfterFirst(selPrefix).BeforeFirst('\n').Trim()
			BookLog('Scan Attachment - Clscan Scanner: ' $ scanner)
			}
		if scanResults.Has?(filename $ ' is saved.')
			{
			BookLog('Scan Attachment End - Clscan')
			return true
			}
		lastLine = scanResults.AfterLast('\n')
		if lastLine is "Unable to scan." or
			lastLine.Prefix?('Problem while opening the scanner')
			return .errorMsg $ lastLine

		BookLog('Scan Attachment End - Clscan')
		return true
		}

	SelectScanner(window)
		{
		scannerOb = .scannerList
		if false isnt selectedScanner = UserSettings.Get('Clscan - Scanner Source')
			if scannerOb.Has?(selectedScanner)
				scannerOb.Remove(selectedScanner).Add(selectedScanner at: 0)
		ClscanSelectSourceControl(window, scannerOb)
		}

	getter_scannerList()
		{
		list = #()
		try
			list = .GetAvailableSources()
		catch (e)
			SuneidoLog('ERRATIC: (CAUGHT) Init Scanner List: ' $ e,
				caughtMsg: 'no message given to user')
		return .scannerList = list
		}

	ScanningAllowed?()
		{
		return not OptContribution('Hosted?', function() { return false })() and
			not Sys.SuneidoJs?()
		}

	ScannerAvailable?()
		{
		return not .scannerList.Empty?()
		}

	// This is run from Accountinglib_PostLogin
	PreGetScannerSources()
		{
		Thread()
			{
			Thread.Name('Clscan-thread')
			Database.SessionId('(Clscan-thread)')
			.scannerList
			}
		}

	SearchLogForUnfinished(logFile)
		{
		unfinished = Object()
		finished = 0
		last = false
		File(logFile, 'r')
			{ |f|
			while false isnt l = f.Readline()
				{
				if not l.Suffix?('gS')
					continue
				if Object?(last) and last.Size() < 4 /*= user's next book logging*/
					{
					lastUser = last[0].AfterFirst(',').Split('\t')[2]
					if l.Has?(lastUser)
						last.Add(l)
					}
				if l.Has?('Scan Attachment Start')
					last = unfinished[.getUserId(l)] = Object(l)
				else if l.Has?('Scan Attachment End')
					{
					unfinished.Delete(.getUserId(l))
					last = false
					finished++
					}
				}
			}
		for m, v in unfinished
			{
			Print(m)
			v.Each(Print)
			}
		Print("Total: ", :finished, unfinished: unfinished.Members().Size())
		return
		}

	getUserId(l)
		{
		cust = l.BeforeFirst(',')
		lineOb = l.AfterFirst(',').Split('\t')
		user = lineOb[2]
		return cust $ ',' $ user
		}
	}