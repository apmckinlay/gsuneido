// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Quality Checker'
	Commands: ((RunCheck, "Ctrl+Q", "Run Quality Checker"))
	Controls:
		(Vert
			(Horz
				(Skip medium:)
				(Vert
					(Skip small:)
					(StarRating))
				Fill
				RefreshButton
				(Skip small:))
			(Editor name: 'qualityCtrl' readonly:)
		)

	New(.libview)
		{
		.threadRunning = false
		.starRating = .FindControl('StarRating')
		.qualityCtrl = .FindControl('qualityCtrl')
		.lib = .libview.CurrentTable()
		.classToCheck = .libview.CurrentName()
		.On_Refresh()
		}

	On_Refresh()
		{
		if .libview.Empty?()
			return
		.qualityCtrl.Set("Running quality checker on " $ .classToCheck)
		if .threadRunning is false
			Thread({ .checkFunc(.lib, .classToCheck) })
		}

	checkFunc(lib, classToCheck)
		{
		Thread.Name('QualityChecker-thread')
		.threadRunning = true
		outputStr = ''
		if lib is ''
			outputStr = "Error: No library selected"
		else if classToCheck is ''
			outputStr = "Error: No class/function record selected"

		code = .getLibText(lib, classToCheck)
		if code is false
			outputStr = "Error: Please select a proper class or function"
		if outputStr isnt ''
			{
			.delayedDisplayOutput(outputStr)
			return
			}
		warningTextAndRating = .createWarningsText(lib, classToCheck, code)

		outputStr = "Class or function checked: " $ classToCheck $ "\n\n" $
			warningTextAndRating.warningText

		.delayedDisplayOutput(outputStr, warningTextAndRating.rating)
		if .Destroyed?()
			return
		.threadRunning = false
		}

	delayedDisplayOutput(text, rating = false)
		{
		.Defer({ .displayOutput(text, rating) })
		}

	displayOutput(text, rating)
		{
		if .Destroyed?()
			return
		.qualityCtrl.Set(text)
		if rating isnt false
			.starRating.SetRating(rating)
		}

	getLibText(lib, name)
		{
		if lib is '' or false is queryResult = Query1(lib, :name, group: -1)
			return false
		return queryResult.lib_current_text
		}

	createWarningsText(lib, classToCheck, code)
		{
		warningsContinuousMethods = Qc_Main.CheckWithExtra(
			lib, classToCheck, code, minimizeOutput?: false)
		continuousText = QcContinuousWarningsOutput(warningsContinuousMethods)

		warningsSlowMethods = Qc_Main.SlowChecks(lib, classToCheck, code).Values(list:)
		slowText = .slowWarningsText(warningsSlowMethods)

		return Record(warningText: continuousText $ "\n" $ slowText,
			rating: Qc_Main.CalcRatings(warningsContinuousMethods))
		}

	slowWarningsText(warningsSlowMethods)
		{
		warningText = "The following quality checks do not contribute to" $
			" the code quality rating.\n\n"
		for individualMethodWarnings in warningsSlowMethods
			{
			warningText $= individualMethodWarnings.desc $ ': \n'
			for warning in individualMethodWarnings.warnings
				warningText $= Opt(warning.name, '\t') $ '\n'
			warningText $= '\n'
			}
		return warningText
		}
 }