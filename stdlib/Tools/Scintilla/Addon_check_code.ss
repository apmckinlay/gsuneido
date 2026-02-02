// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonForThreadTasks
	{
	AddonName:		'CheckCode'
	warningLevel:	98
	errorLevel:		99
	Init()
		{
		super.Init()
		.marker_error = .MarkerIdx(level: .errorLevel)
		.indic_error = .IndicatorIdx(level: .errorLevel)

		.marker_warning = .MarkerIdx(level: .warningLevel)
		.indic_warning = .IndicatorIdx(level: .warningLevel)
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

	prevWarnings: #(lines: (), qc: (warningText: 'init', rating: false))
	ThreadFn(startTime)
		{
		checkCodeResult = CheckCode(.code, .name, .table, checkCodeWarnings = Object())
		if .IsOutdatedRecord(startTime)
			return

		qcWarnings = .qcWarnings()
		if .IsOutdatedRecord(startTime)
			return

		warnings = .buildWarnings(checkCodeWarnings, qcWarnings)
		if .updateLibView?(warnings, .prevWarnings) and not .IsOutdatedRecord(startTime)
			.Defer({ .updateLibView(startTime, warnings, checkCodeResult) },
				uniqueID: #Addon_check_code)
		}

	qcWarnings()
		{
		try
			if QcIsEnabled()
				return Qc_Main.CheckWithExtra(.table, .name, .code, minimizeOutput?:)
		catch (e /*unused*/, 'Outdated')
			{ }
		return Object(lines: Object(), all: Object())
		}

	buildWarnings(checkCodeWarnings, qcWarnings)
		{
		warnings = Object(lines: Object(), qc: Object())
		warnings.lines.Append(checkCodeWarnings)
		warnings.lines.Append(qcWarnings.GetDefault('lineWarnings', #()))
		warnings.qc.warningText =
			Opt(.createCheckCodeText(checkCodeWarnings, .name, .table), '\n') $
			QcContinuousWarningsOutput(qcWarnings)
		warnings.qc.rating = Qc_Main.CalcRatings(qcWarnings)
		return warnings
		}

	updateLibView?(warnings, prevWarnings)
		{
		return warnings.qc.rating isnt prevWarnings.qc.rating or
			warnings.qc.warningText isnt prevWarnings.qc.warningText or
			not warnings.lines.EqualSet?(prevWarnings.lines)
		}

	updateLibView(startTime, warnings, checkCodeResult)
		{
		if .IsOutdatedRecord(startTime)
			return

		if warnings.lines isnt .prevWarnings.lines
			.displayWarnings(warnings.lines, checkCodeResult)

		if warnings.qc isnt .prevWarnings.qc
			{
			.updateAnnotations(warnings.qc.warningText)
			.Send(#CheckCode_QualityChanged, warnings.qc)
			}
		.prevWarnings = warnings
		}

	updateAnnotations(allWarningsText)
		{
		if .Empty?() or .Destroyed?()
			return
		.SendToAddons(#ClearAllAnnotations)
		warnings = Object()
		for line in allWarningsText.Lines()
			{
			if line.Prefix?(.table $ ':' $ .name $ ':')
				{
				afterPrefix = line.AfterFirst(.table $ ':' $ .name $ ':')
				lineNum = Number(afterPrefix.BeforeFirst(' ')) - 1
				warningText = afterPrefix.AfterFirst(' - ')
				warnings[lineNum] = warnings.Member?(lineNum)
					? warnings[lineNum] $ '\r\n' $ warningText
					: warningText
				}
			}
		for lineNum in warnings.Members()
			.SendToAddons(#AddAnnotation, lineNum, warnings[lineNum])
		}

	createCheckCodeText(warnings, name, lib)
		{
		outputText = ''
		for warning in warnings
			if not warning.GetDefault('noOutput', false)
				outputText $= '\n' $ lib $ ':' $ name $ ':' $ (warning.line + 1) $ ' ' $
				warning.msg
		return outputText
		}

	displayWarnings(warnings, checkCodeResult)
		{
		if .Destroyed?() or warnings is false or .Empty?()
			return
		.clear()
		.processWarnings(warnings)
		if checkCodeResult is false
			.Send(#Status, '\terrors', invalid:)
		else
			.Send(#Status, '\t ', normal:)
		.MarkersChanged() // e.g. update overview bar
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
				.squiggle(err.pos, err.len, warning: err.msg.Prefix?('WARNING'))
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
	}
