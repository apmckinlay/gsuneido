// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	styleLevel: 59
	New(@args)
		{
		super(@args)
		.prevInserted = Object()
		}
	Init()
		{
		.breakFlag = .MarkerIdx(level: .styleLevel)
		.breakIndicator = .IndicatorIdx(level: .styleLevel)
		}

	Styling()
		{
		return [
			[level: .styleLevel,
				marker: [marker: SC.MARK_SHORTARROW, back: CLR.RED],
				indicator: [INDIC.ROUNDBOX, fore: CLR.RED]]]
		}

	ContextMenu()
		{
		return #('Stepping Break Point\tCtrl+S', 'Clear All Stepping Points')
		}

	On_Stepping_Break_Point()
		{
		if false is info = .getInfo()
			return

		src = .Get()
		pos = .GetSelectionStart()

		if false is pos = .findClosestNonBlankChar(src, pos)
			return .alert()

		if false is
			SteppingDebuggerManager().ToggleBreakPoint(info.lib, info.name, src, pos)
			return .alert()

		.updateMarkerAndIndicator(info)
		.prevInserted.Add(info)
		}

	getInfo()
		{
		lib = .Send("CurrentTable")
		name = .Send("CurrentName")
		if not String?(lib) or not String?(name)
			return false

		return Object(:lib, :name)
		}

	findClosestNonBlankChar(src, pos)
		{
		if not src[pos].Blank?()
			return pos

		if pos >= src.Size()
			return false

		left = pos - 1
		right = pos
		while (left >= 0 and src[left] in ('\t', ' '))
			left--
		while (right < src.Size() and src[right] in ('\t', ' '))
			right++
		return .pick(src, left, right)
		}

	pick(src, left, right)
		{
		if right < src.Size() and src[right] not in ('\r', '\n')
			return right
		if left >= 0 and src[left] not in ('\r', '\n')
			return left
		return false
		}

	alert()
		{
		InfoWindowControl('Cannot add stepping debugging point', titleSize: 0,
			marginSize: 7, autoClose: 1)
		return
		}

	updateMarkerAndIndicator(info)
		{
		.MarkerDeleteAll(.breakFlag)
		.ClearIndicator(.breakIndicator)

		for range in SteppingDebuggerManager().GetBreakPointRanges(info.lib, info.name)
			{
			lineNumber = .LineFromPosition(range.i)
			.MarkerAdd(lineNumber, .breakFlag)
			.SetIndicator(.breakIndicator, range.i, range.n)
			}
		.MarkersChanged()
		}

	On_Clear_All_Stepping_Points(clear = false)
		{
		if false is info = .getInfo()
			return

		.MarkerDeleteAll(.breakFlag)
		.ClearIndicator(.breakIndicator)

		clear is false
			? SteppingDebuggerManager().RemoveAllBreakPoints(@info)
			: SteppingDebuggerManager().ClearSource(@info)
		.prevInserted.Remove(info)
		}

	curInfo: false
	IdleAfterChange()
		{
		if false is info = .getInfo()
			return

		if info isnt .curInfo
			{
			.curInfo = info
			.updateMarkerAndIndicator(info)
			}
		else if not .GetReadOnly()
			{
			.On_Clear_All_Stepping_Points(clear:)
			}
		}

	SetFocus()
		{
		.update()
		}

	Set()
		{
		.update()
		}

	update()
		{
		if false is info = .getInfo()
			return

		.updateMarkerAndIndicator(info)
		}

	Destroy()
		{
		if .Send('RemoveAllBreakPointsOnDestroy?') is false
			return

		.prevInserted.Each({ SteppingDebuggerManager().RemoveAllBreakPoints(@it) })
		}
	}