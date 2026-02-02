// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonForThreadTasks
	{
	AddonName: 		'ShowModifiedLines'
	markerType: 	'diff'
	addLevel: 		28
	modifyLevel: 	29
	deleteLevel: 	27
	Init()
		{
		super.Init()
		.init()
		.marker_add = .MarkerIdx(level: .addLevel, type: .markerType)
		.marker_modify = .MarkerIdx(level: .modifyLevel, type: .markerType)
		.marker_delete = .MarkerIdx(level: .deleteLevel, type: .markerType)
		}

	init()
		{
		.halfHeight = (.TextHeight() / 2).Round(0)
		.width = .GetMarginWidthN(1)
		}

	Styling()
		{
		.init()
		return [
			.buildMarker('add', level: .addLevel),
			.buildMarker('modify', level: .modifyLevel),
			.buildMarker('delete', top: '000', level: .deleteLevel)
			]
		}

	buildMarker(type, top = '001', bottom = '001', level = 0)
		{
		rep = (.width / 3).Int()  /*= modified line uses 1/3 of the marker bar */
		top = top.Map({ it.Repeat(rep) }).RightFill(.width, top[-1])
		bottom = bottom.Map({ it.Repeat(rep) }).RightFill(.width, bottom[-1])
		marker = `/* XPM */
			static char * XFACE[] = {
			/* <Values> */
			/* <width/columns> <height/rows> <colors> <chars per pixel>*/
			"` $ .width $ ` ` $ (.halfHeight * 2) $ ` 2 1",
			/* <Colors> */
			"0 c none",
			"1 c #` $ .getRGBString(.GetSchemeColor(type)) $ `",
			/* <Pixels> */`
		marker $= ('"' $ top $ '",\r\n').Repeat(.halfHeight)
		marker $= ('"' $ bottom $ '",\r\n').Repeat(.halfHeight)
		marker = marker.RemoveSuffix(',\r\n') $ `};`
		return [:level, marker: [marker, back: .GetSchemeColor(type), type: .markerType]]
		}

	getRGBString(color)
		{
		hex = color.Hex().LeftFill(6/*= rgb hex string*/, '0')
		return hex[-2..] $ hex[2::2] $ hex[::2]
		}

	Set() // reseting members when viewing new records
		{
		.diffs = false
		.rec = false
		}

	setRec()
		{
		if not .validateParams()
			.rec = false
		else
			.rec = Query1(SvcTable(.table).NameQuery(.name))

		if .rec isnt false and .rec.lib_before_text is ''
			.rec.lib_before_text = .rec.text
		return not .invalid?(.rec)
		}

	invalid?(rec)
		{
		return rec is false or not rec.Member?('lib_committed') or
			rec.lib_committed is "" or not rec.Member?('lib_before_text') or
			(rec.lib_modified isnt "" and rec.lib_before_text is "")
		}

	PreThread()
		{
		.text = .Get()
		return .setRec()
		}

	validateParams()
		{
		try
			{
			.name = .Send('CurrentName')
			.table = .Send('CurrentTable')
			}
		catch(unused, '*socket connection timeout')
			return false
		return .name not in (0,'') and .table not in (0,'')
		}

	prevChangedLines: #()
	ThreadFn(startTime)
		{
		changedLines = .getChangedLines(.rec.lib_before_text, .text)
		if .IsOutdatedRecord(startTime) or .prevChangedLines.EqualSet?(changedLines)
			return

		.Defer(uniqueID: 'modified_lines')
			{
			if not .Destroyed?() and not .IsOutdatedRecord(startTime)
				{
				.addDiffMarkers(.prevChangedLines = changedLines)
				.MarkersChanged()
				}
			}
		}

	getChangedLines(beforeText, text)
		{
		.diffs = Diff.SideBySide(beforeText.Lines(), text.Lines())
		i = 0
		deleteCount = 0
		while i < .diffs.Size()
			{
			if .diffs[i][1] is "<"
				deleteCount++
			.diffs[i].Add(i - deleteCount)
			i++
			}
		return .diffs.Filter({ it[1] isnt "" })
		}

	clearMarkers()
		{
		.MarkerDeleteAll(.marker_delete)
		.MarkerDeleteAll(.marker_add)
		.MarkerDeleteAll(.marker_modify)
		}

	addDiffMarkers(changedLines)
		{
		.clearMarkers()
		lineIndex = 3
		for line in changedLines
			{
			if line[1] is "#" and line[0] is line[2]
				continue
			if line[1] is "<"
				.addMarker(line[lineIndex], .marker_delete)
			else if line[1] is ">"
				.addMarker(line[lineIndex], .marker_add)
			else
				.addMarker(line[lineIndex], .marker_modify)
			}
		}

	addMarker(line, marker)
		{
		.MarkerAdd(line, marker)
		}

	Invalidate()
		{
		super.IdleAfterChange() // Do not clear: .prevChangedLines on Invalidate
		}

	IdleAfterChange() // Explicit modification, next ThreadFn call should run
		{
		.prevChangedLines = #(false)
		super.IdleAfterChange()
		}

	ContextMenu()
		{
		return #('Show Original Lines\tCtrl+O')
		}

	On_Show_Original_Lines()
		{
		if not .setRec()
			return
		result = .findPrevLines()
		if result is false
			return
		ToolDialog(0, Object(.popup, result.prevLines, result.localLineHeight, this),
			keep_size: false, border: 0)
		}

	findPrevLines()
		{
		if .diffs is false
			return false
		sel = .GetSelect()
		first = .LineFromPosition(sel.cpMin)
		last = .LineFromPosition(sel.cpMax)
		localLineHeight = last - first + 1
		prevLines = ''
		lineIndex = 3
		for (i = first; i < .diffs.Size(); i++)
			{
			localLineNumber = .diffs[i][lineIndex]
			if localLineNumber > last
				break
			if localLineNumber >= first and .diffs[i][1] isnt '>'
				prevLines $= .diffs[i][0] $ '\n'
			}
		first = .PositionFromLine(first)
		last = .PositionFromLine(last + 1)
		.SetSelect(first, last - first)
		return Object(:localLineHeight, :prevLines)
		}

	popup: Controller
		{
		Title: "Original Lines"
		New(.prevLines, localLineHeight, .scintillaAddon)
			{
			super(Object('Vert' Object('DisplayCode',
				set: prevLines,
				height: Min(30/*=min*/,
					Max(prevLines.Lines().Size(), localLineHeight)),
				width: 100)
				#(Horz Fill (MenuButton 'Restore' ('Restore')) Fill)))
			}

		On_Restore_Restore()
			{
			.scintillaAddon.Paste(.prevLines)
			.Window.Result(true)
			}
		}
	}
