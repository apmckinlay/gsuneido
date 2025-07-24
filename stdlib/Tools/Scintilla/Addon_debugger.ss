// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// requires On_ methods be forwarded to Scintilla as LibView does
ScintillaAddon
	{
	insertCodes: 	false
	breakLevel: 	56
	watchLevel: 	57
	New(@args)
		{
		super(@args)
		.insertCodes = Object().Set_default(Object())
		.prevInserted? = Object().Set_default(false)
		}
	Init()
		{
		.breakFlag = .MarkerIdx(level: .breakLevel)
		.breakIndicator = .IndicatorIdx(level: .breakLevel)

		.watchFlag = .MarkerIdx(level: .watchLevel)
		.watchIndicator = .IndicatorIdx(level: .watchLevel)

		.insertTemplates = Object()
		.insertTemplates[.watchFlag] =  .watchTemplateCode
		.insertTemplates[.breakFlag] =  '/*DEBUG*/Debugger.Dialog(' $
			'GetActiveWindow(), "Break Point", GetCallStack())/*DEBUG*/'
		}

	Styling()
		{
		return [
			[level: .breakLevel,
				marker: [SC.MARK_SHORTARROW, back: CLR.RED],
				indicator: [INDIC.BOX, fore: CLR.RED]],
			[level: .watchLevel,
				marker: [SC.MARK_SHORTARROW, back: CLR.amber]
				indicator: [INDIC.ROUNDBOX, fore: CLR.amber]]]
		}

	watchTemplateCode(selText)
		{
		if selText is ''
			return '/*DEBUG*/Print(WATCH: Locals(0))/*DEBUG*/'
		else
			{
			printLabel = selText.Tr('\r\n').Ellipsis(30) /*= max number of characters */
			return '/*DEBUG*/Print(`WATCH - ' $ printLabel $
				'`: ' $ selText $ ')/*DEBUG*/'
			}
		}

	ContextMenu()
		{
		return #('Watch Point\tCtrl+W', 'Break Point\tCtrl+B', 'Clear All Points')
		}

	On_Clear_All_Points()
		{
		.MarkerDeleteAll(.watchFlag)
		.MarkerDeleteAll(.breakFlag)
		.ClearIndicator(.watchIndicator)
		.ClearIndicator(.breakIndicator)

		.insertCodes = Object().Set_default(Object())
		LibraryOverrideClear()
		}

	On_Watch_Point()
		{
		.addDebugger(.watchFlag)
		}

	On_Break_Point()
		{
		.addDebugger(.breakFlag)
		}

	addDebugger(marker)
		{
		lineNumber = .LineFromPosition()
		state = .MarkerGet(lineNumber)

		lib = .Send("CurrentTable")
		name = .Send("CurrentName")

		existFlag = false
		for markerFlag in .insertTemplates.Members()
			if 0 isnt (state & (1 << markerFlag))
				existFlag = markerFlag

		if existFlag isnt false // removing
			{
			.remove(lineNumber, existFlag, lib, name)
			return
			}

		sel = .GetSelText()
		if not .isPointValid?([:lib, :name, :lineNumber], marker, sel)
			{
			InfoWindowControl('Cannot add debugging point' $ Opt(' with "', sel, '"'),
				titleSize: 0, marginSize: 7, autoClose: 1)
			return
			}

		if marker isnt existFlag
			{
			.MarkerAdd(lineNumber, marker)
			.addIndicator(marker, sel)
			.insertCodes[lib $ ':' $ name][lineNumber] = [:marker, selection: sel]
			}

		.MarkersChanged()

		.ensureDebuggingPoints(lib, name)
		}

	remove(lineNumber, existFlag, lib, name)
		{
		.MarkerDelete(lineNumber, existFlag)
		// need to delete twice after switching tabs in lib view
		.MarkerDelete(lineNumber, existFlag)

		lineSize = .GetLine().Size()
		lineStart = .Get()[.. .GetCurrentPos()].FindLast('\n')
		.ClearIndicator(.watchIndicator, lineStart, lineSize)
		.ClearIndicator(.breakIndicator, lineStart, lineSize)
		.insertCodes[lib $ ':' $ name].Delete(lineNumber)
		.MarkersChanged()
		.ensureDebuggingPoints(lib, name)
		}

	addIndicator(marker, selText)
		{
		text = .Get()
		lineStart = text.FindLast('\n', .GetCurrentPos())
		if marker is .watchFlag and selText isnt ''
			{
			line = .GetLine()
			selStart = lineStart + line.Find(selText) + 1
			.SetIndicator(.watchIndicator, selStart, selText.Size())
			}
		else
			{
			len = .GetLine().Size()
			indic = marker is .watchFlag ? .watchIndicator : .breakIndicator
			.SetIndicator(indic, lineStart, len)
			}
		}

	isPointValid?(point, marker, sel)
		{
		currentCode = .Get()
		savedCodes = .insertCodes[point.lib $ ':' $ point.name].Copy()
		savedCodes[point.lineNumber] = [:marker, selection: sel]
		inserts = .collectInserts(savedCodes, point.lineNumber, marker)
		newCode = .buildDebugCode(currentCode, inserts)
		return Compilable?(newCode)
		}

	ensureDebuggingPoints(lib, name)
		{
		currentCode = .Get()

		savedCodes = .insertCodes[lib $ ':' $ name]
		inserts = .collectInserts(savedCodes)
		if not inserts.Empty?()
			{
			LibraryOverride(lib, name, .buildDebugCode(currentCode, inserts))
			.prevInserted?[lib $ ':' $ name] = true
			}
		else if .prevInserted?[lib $ ':' $ name]
			{
			LibraryOverride(lib, name)
			.prevInserted?[lib $ ':' $ name] = false
			}
		}
	collectInserts(savedCodes, testLineNumber = false, testMarker = false)
		{
		inserts = Object()
		index = 0
		lines = .GetMarkerLines()
		if testLineNumber isnt false
			lines.AddUnique(testLineNumber).Sort!()
		for lineNumber in lines
			{
			marker = .MarkerGet(lineNumber)
			isDebugPoint? = false
			for markerFlag in .insertTemplates.Members()
				{
				if ((testLineNumber is lineNumber and markerFlag is testMarker) or
					(0 isnt (marker & (1 << markerFlag))))
					{
					.updateSavedCodes(savedCodes, index, lineNumber)
					.addInsertCode(markerFlag, savedCodes, lineNumber, inserts)
					isDebugPoint? = true
					}
				}
			if isDebugPoint?
				++index
			}
		return inserts
		}

	updateSavedCodes(savedCodes, index, lineNumber)
		{
		savedLines = savedCodes.Members().Sort!()
		savedLineNumber = savedLines[index]
		if lineNumber isnt savedLineNumber
			{
			ins = savedCodes[savedLineNumber].Copy()
			savedCodes.Delete(savedLineNumber)
			savedCodes[lineNumber] = ins
			}
		}

	addInsertCode(markerFlag, savedCodes, lineNumber, inserts)
		{
		insertCode = .insertTemplates[markerFlag]
		if Function?(insertCode)
			insertCode = insertCode(savedCodes[lineNumber].selection)
		inserts[lineNumber] = insertCode
		}

	buildDebugCode(currentCode, inserts)
		{
		codeLines = currentCode.Lines()
		i = 0
		for lineNumber in inserts.Members().Sort!()
			{
			line = .GetLine(lineNumber).Trim()
			if line =~ '^(if|for|while|return)'
				{
				codeLines.Add(inserts[lineNumber], at: lineNumber + i)
				++i // insert offset
				}
			else
				{
				codeLines.Add('if(true){' $ inserts[lineNumber], at: lineNumber + i)
				codeLines.Add('}', at: lineNumber + i + 2)
				i += 2 // insert offset
				}
			}
		return codeLines.Join('\r\n')
		}

	IdleAfterChange()
		{
		lib = .Send("CurrentTable")
		name = .Send("CurrentName")
		.ensureDebuggingPoints(lib, name)
		}

	Destroy()
		{
		if .prevInserted?.NotEmpty?()
			LibraryOverrideClear()
		}
	}