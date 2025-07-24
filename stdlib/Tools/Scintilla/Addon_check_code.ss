// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonForThreadTasks
	{
	AddonName: 		"CheckCode"
	warningLevel:	98
	errorLevel:		99
	Init()
		{
		super.Init()
		.marker_error = .MarkerIdx(level: .errorLevel)
		.indic_error = .IndicatorIdx(level: .errorLevel)

		.marker_warning = .MarkerIdx(level: .warningLevel)
		.indic_warning = .IndicatorIdx(level: .warningLevel)

		.oldQcInfo = false
		.oldLineWarnings = false
		.oldCheckCodeResult = false
		}

	Styling()
		{
		return [
			[level: .warningLevel,
				marker: [SC.MARK_ROUNDRECT, back: CLR.DARKGRAY],
				indicator: [INDIC.SQUIGGLE, fore: CLR.DARKGRAY]],
			[level: .errorLevel,
				marker: [SC.MARK_ROUNDRECT, back: CLR.RED],
				indicator: [INDIC.SQUIGGLE, fore: CLR.RED]]]
		}

	Set()
		{
		//TODO: Used named threads when exe is released to limit # of threads to 2
		.fastWarnings = false
		.extraWarnings = false
		.name = .Send(#CurrentName)
		.table = .Send(#CurrentTable)
		.oldQcInfo = false
		.oldLineWarnings = false
		.oldCheckCodeResult = false
		.Send(#CheckCode_QualityChanged, [warningText: .MsgWhenChecking, rating: false])
		}
	prev_result: ''
	prev_warnings: ()
	prev_line_warnings: ()

	Invalidate()
		{
		.IdleAfterChange()
		}

	PreThread()
		{
		.code = .Get()
		if 0 is .name = .Send(#CurrentName)
			.name = false
		if 0 is .table = .Send(#CurrentTable)
			.table = false
		return .name isnt false and .table isnt false
		}

	ThreadFn(startTime)
		{
		.fastWarnings = .extraWarnings = false
		.checkCodeResult = CheckCode(.code, .name, .table, .checkCodeWarnings = Object())
		if not QcIsEnabled()
			{
			.Defer({ .displayWarnings(.checkCodeWarnings) },
				uniqueID: 'CheckCodeDisplayWarns')
			return
			}

		if false isnt (.fastWarnings = .checkWarning(startTime)) and
			false isnt (.extraWarnings = .checkWarning(startTime, extraChecks:))
			.Defer({ .updateLibView(startTime) },  uniqueID: 'CheckCodeUpdateLibView')
		}

	checkWarning(startTime, extraChecks = false)
		{
		if .IsOutdatedRecord(startTime)
			return false
		try
			warningOb = Qc_Main(.table, .name, .code, minimizeOutput?:,
				:extraChecks)
		catch(e /*unused*/, "Outdated")
			return false

		return warningOb
		}

	MsgWhenChecking: "Checking Code..."

	updateLibView(startTime)
		{
		if .IsOutdatedRecord(startTime)
			return

		lineWarnings = Object()
		curWarnings = Object()
		lineWarnings.Append(.checkCodeWarnings)
		if .fastWarnings isnt false
			{
			lineWarnings.Append(.fastWarnings.lineWarnings)
			curWarnings.Append(.fastWarnings)
			}
		if .extraWarnings isnt false
			{
			lineWarnings.Append(.extraWarnings.lineWarnings)
			curWarnings.Append(.extraWarnings)
			}

		curWarningsText = QcContinuousWarningsOutput(curWarnings)
		checkCodeText = .createCheckCodeText(.checkCodeWarnings, .name, .table)
		qcInfo = Object()
		qcInfo.warningText = Opt(checkCodeText, '\n') $ curWarningsText
		if .extraWarnings is false or .fastWarnings is false
			qcInfo.warningText = .MsgWhenChecking $ qcInfo.warningText
		qcInfo.rating = Qc_Main.CalcRatings(curWarnings)

		if lineWarnings isnt .oldLineWarnings
			{
			.displayWarnings(lineWarnings)
			.oldLineWarnings = lineWarnings
			.oldCheckCodeResult = .checkCodeResult
			}
		if qcInfo isnt .oldQcInfo
			{
			.updateAnnotations(qcInfo.warningText)
			.Send("CheckCode_QualityChanged", qcInfo)
			.oldQcInfo = qcInfo
			}
		}

	updateAnnotations(allWarningsText)
		{
		if .Empty?() or .Destroyed?()
			return
		.SendToAddons("ClearAllAnnotations")
		warnings = Object()
		for line in allWarningsText.Lines()
			{
			if line.Prefix?(.table $ ':' $ .name $ ':')
				{
				afterPrefix = line.AfterFirst(.table $ ':' $ .name $ ':')
				lineNum = Number(afterPrefix.BeforeFirst(' ')) - 1
				warningText = afterPrefix.AfterFirst(" - ")
				warnings[lineNum] = warnings.Member?(lineNum)
					? warnings[lineNum] $ '\r\n' $ warningText
					: warningText
				}
			}
		for lineNum in warnings.Members()
			.SendToAddons("AddAnnotation", lineNum, warnings[lineNum])
		}

	createCheckCodeText(warnings, .name, lib)
		{
		outputText = ""
		for warning in warnings
			if not warning.GetDefault('noOutput', false)
				outputText $= '\n' $ lib $ ':' $ .name $ ':' $ (warning.line + 1) $ ' ' $
				warning.msg
		return outputText
		}

	displayWarnings(warnings)
		{
		if .Destroyed?() or warnings is false or .Empty?()
			return
		.clear()
		.processWarnings(warnings)
		.setStatus()
		.MarkersChanged() // e.g. update overview bar
		.prev_result = .checkCodeResult
		.prev_warnings = warnings
		}

	clear()
		{
		.ClearIndicator(.indic_error)
		.ClearIndicator(.indic_warning)
		.MarkerDeleteAll(.marker_error)
		.MarkerDeleteAll(.marker_warning)
		}

	processWarnings(warnings)
		{
		for err in warnings
			{
			if err.Member?(#pos) and err.Member?(#len) and err.Member?(#msg)
				.squiggle(err.pos, err.len, warning: err.msg.Prefix?("WARNING"))
			else if err.Size(list:) >= 2 // old style qc
				.squiggle(err[0], err[1], warning: err.GetDefault(#warning, false))
			else
				.AddWarning(err[0] - 1) //Convert to 0 based numbering
			}
		}

	squiggle(pos, len, warning = false)
		{
		.SetIndicator(warning ? .indic_warning : .indic_error, pos, len)
		.MarkerAdd(.LineFromPosition(pos), warning ? .marker_warning : .marker_error)
		}

	AddWarning(line)
		{
		.MarkerAdd(line, .marker_warning)
		}

	setStatus()
		{
		if .checkCodeResult is false
			.Send('Status', '\terrors', invalid:)
		else
			.Send('Status', '\t ', normal:)
		}
	}
