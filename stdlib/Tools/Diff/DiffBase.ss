// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
// used by LibDiff and SvcControl
// TODO if we passed in functions for listNew and listOld
// then refresh could be implemented here (e.g. instead of LibDiff)
PassthruController
	{
	New(list1, list2, .lib, .recName, title1, title2, .base = false,
		.tags = "12", .gotoButton = false, .comment = false, .refresh = false,
		.commentBgColor = false, .newOnRight? = false, .extraControls = #(Skip),
		titleNewRight = "", titleOldRight = "")
		{
		.SetRedirs()
		result = .GetDiffs(list1, list2, base)
		.model = result.model
		.diffs = result.diffs
		.getControls()

		.SetTitles(title1, title2, :titleNewRight, :titleOldRight)
		.setupChangePrefix(title1, title2)

		.SettupLists()
		.ToFirstDiff()
		}

	CurrentTable()
		{
		return .lib
		}

	ToFirstDiff()
		{
		if not .diffs.Empty?() and .diffs[0][1] is "" // not already on a change
			.On_Next_Change()
		}

	SetRedirs()
		{
		throw 'Must Implement'
		}

	SetTitles(@unused)
		{
		throw 'Must Implement'
		}

	On_Next_Change()
		{
		throw 'Must Implement'
		}

	GetDiffs(list1, list2, base)
		{
		if String?(list1)
			list1 = .normalizeLineEnds(list1).Lines()
		if String?(list2)
			list2 = .normalizeLineEnds(list2).Lines()

		if String?(base)
			base = base.Lines()

		model = Diff2Model(list2, list1, base)
		diffs = model.Diffs

		return Object(:model, :diffs)
		}

	// relace single '\r' with '\r\n'
	normalizeLineEnds(s)
		{
		start = i = 0
		normalized = ''
		while s.Size() isnt i = s.Find('\r', i)
			{
			if i + 1 < s.Size() and s[i + 1] isnt '\n' or i + 1 is s.Size()
				{
				normalized $= s[start..i] $ '\r\n'
				start = i + 1
				}
			i++
			}
		return normalized $= s[start..i]
		}

	getControls()
		{
		.GetControls()

		.mergeWarn = .FindControl(#mergeWarn)
		.showLine = .FindControl(#show_line_difference)
		.lineDiff = .FindControl(#lineDiffVert)
		.setShowHideLineDiff()
		.text = .FindControl(#text)
		}

	GetControls()
		{
		throw 'Must Implement'
		}

	setShowHideLineDiff()
		{
		if 'novalue' is value = UserSettings.Get('VersionControl-LineDiff', 'novalue')
			UserSettings.Put('VersionControl-LineDiff', value = true)
		if value is true
			.showLineDiff()
		.showLine.Set(value)
		}

	showLineDiff()
		{
		.lineDiff.Insert(0, Object('WorkSpaceCode', name: 'lineDiff', height: 2,
			readonly:, Addon_adjust_height:))
		}

	removeLineDiff()
		{
		.lineDiff.Remove(0)
		}

	SettupLists()
		{
		.lists = .lists_class(@.ListControls())

		.ovbar = new .lists_class(@.OvBars())
		.ovbar.SetNumRows(.diffs.Size())
		.ovbar.SetMaxRowHeight(.lists.Get(0).TextHeight(0), scaled?:)
		.ovbar.SetTopMargin(GetSystemMetrics(SM.CXHSCROLL))

		selColor = CLR.GRAY
		.lists.SetSelBack(true, selColor)

		.SetLists(.diffs)

		if .base isnt false
			{
			.warnings = .checkCode(.CodeToCheck)
			if not .warnings.Empty?()
				.SetWarn(.warnings[0])
			}

		.Defer({ .AddMarkers(.diffs, .lists, .model, .ovbar) })
		.SetProcs()
		}

	SetWarn(unused) { }

	SetLists(@unused)
		{
		throw 'Must Implement'
		}

	SetProcs()
		{
		throw 'Must Implement'
		}

	checkCode(mergedCode)
		{
		results = Object()
		CheckCode(mergedCode, :results)
		return results.Map!({ it.msg.Trim('- ') $ ' - line: ' $ it.line })
		}

	Static_DoubleClick(source)
		{
		if source isnt .mergeWarn
			return 0
		// Want to use InfoWindow here but it is not calculating text widths correctly
//		InfoWindowControl( .warnings.Join('\r\n'), 'Code Errors')
		.AlertInfo('Code Errors', .warnings.Join('\r\n'))
		}

	setupChangePrefix(titleNew, titleOld)
		{
		size = Max(titleNew.Find(' '), titleOld.Find(' '))
		titleLeft = .newOnRight? ? titleOld : titleNew
		titleRight = .newOnRight? ? titleNew : titleOld
		.leftChangePre = titleLeft[..size].Trim().LeftFill(size, ' ')
		if ((titleCut = titleRight[..size]) is "AS OF 20")
			.rightChangePre = "SELECTED".LeftFill(size, ' ')
		else
			.rightChangePre	= titleCut.Trim().LeftFill(size, ' ')
		}

	lists_class: class
		{
		New(@lists)
			{
			.lists = lists
			}
		Get(i)
			{
			return .lists[i]
			}
		Size()
			{
			return .lists.Size()
			}
		forList(block)
			{
			for idx in .lists.Members()
				(block)(.lists[idx], idx)
			}
		Default(@args)
			{
			.forList()
				{ |list, idx /*unused*/|
				if list.Method?(args[0])
					list[args[0]](@+1args)
				}
			}
		// wrappers for ListControl Methods
		SetScroll(main)
			{
			mainScroll = GetScrollPos(.Get(main).Hwnd, SB.VERT)
			.forList()
				{ |list, idx|
				if idx isnt main
					{
					scroll = GetScrollPos(list.Hwnd, SB.VERT)
					list.LineScroll(0, (mainScroll - scroll))
					}
				}
			}
		SetXOffSet(main)
			{
			mainScroll = GetScrollPos(.Get(main).Hwnd, SB.HORZ)
			.forList()
				{ |list, idx|
				if idx isnt main
					list.SETXOFFSET(mainScroll)
				}
			}
		SetFirstVisibleLine(main)
			{
			line = .Get(main).GetFirstVisibleLine()
			.forList()
				{ |list, idx|
				if idx isnt main
					list.SetFirstVisibleLine(line)
				}
			return .Get(main).LineFromPosition()
			}
		}

	AddMarkers(diffs, lists, model, ovbar)
		{
		if not Instance?(lists) or not Instance?(model) or not Instance?(ovbar)
			return

		rowAdjust = VertScrollBar?(.lists.Get(0).Hwnd) ? 0 : -1
		colors = .colors()
		for row, diff in diffs
			{
			indic = diff[1]
			if indic is ''
				continue
			.AddMarker(indic, row, lists, model)
			ovbar.AddMark(row + rowAdjust, colors[diffs[row][1][0]])
			}
		.Defer(ovbar.Repaint)
		}

	AddMarker(@unused)
		{
		throw 'Must Implement'
		}

	PaintIndics(list, row, posInfo, selected? = false)
		{
		linePos = list.PositionFromLine(row)
		nextLinePos = list.PositionFromLine(row + 1)
		posInfo.Each()
			{
			list.AddIndic(linePos + it.pos, it.length,
				linePos, nextLinePos - linePos, :selected?)
			}
		}

	colors()
		{
		green = 0x0000c800
		return Object(
			'>': green
			'+': green
			'<': CLR.RED,
			'-': CLR.RED,
			'#': CLR.BLUE)
		}

	ListProc(hwnd, msg, wparam, lparam, main, prevproc)
		{
		return .listproc(hwnd, msg, wparam, lparam, main, .lists, prevproc)
		}

	listproc(hwnd, msg, wparam, lparam, main, lists, prevproc)
		{
		_hwnd = .WindowHwnd()
		if .preventScroll(msg)
			return 0
		result = CallWindowProc(prevproc, hwnd, msg, wparam, lparam)
		if msg is WM.VSCROLL or msg is WM.MOUSEWHEEL
			{
			lists.SetScroll(main)
			.Ovbar.Scintilla_VScroll()
			}
		else if msg is WM.HSCROLL
			lists.SetXOffSet(main)
		else if msg is WM.LBUTTONDOWN
			.updateSelectedLine(lists.Get(main).LineFromPosition())
		else if msg is WM.KEYDOWN and wparam in (VK.DOWN, VK.UP, VK.LEFT, VK.RIGHT)
			.updateSelectedLine(lists.SetFirstVisibleLine(main))
		return result
		}
	preventScroll(msg)
		{
		return msg is WM.MOUSEWHEEL and KeyPressed?(VK.CONTROL)
		}

	buttons()
		{
		buttons = ['Horz' 'Skip'
			#(Button "Previous Change") 'Skip' #(Button "Next Change")]
		if .gotoButton isnt false
			buttons.Add(#(Skip 50), #(Button "Go To Definition"))
		if .refresh
			buttons.Add(#(Skip 50), #(Button "Refresh"))
		return buttons
		}

	Controls()
		{
		controls = Object('VertSplit', .DiffPanes(),
			Object('Vert'
				Object('Horz', .buttons() #(Skip 20) 'show_line_difference',
					'Fill', .extraControls, #(Skip))
				#(Skip medium:)
				name: 'lineDiffVert'
				)
			)
		if .comment isnt false
			controls[1].Add(.Comment(.comment), at: 1)
		return controls
		}

	DiffWithBar(diffname, barname, readonly = true)
		{
		showMargin = TableExists?(.lib)
			? QueryColumns(.lib).Has?(#group)
			: false
		return 	Object('Horz'
			Object('ScintillaDiff', xstretch: 1, ystretch: 200, fontSize: 11,
				name: diffname, :readonly, :showMargin),
			Object('OverviewBar', name: barname))
		}

	commentBgColor: false
	Comment(comment)
		{
		bgColor = .commentBgColor is false ? CLR.azure : .commentBgColor
		return Object('Editor', set: comment, size: '+1', ystretch: 0, readonly:,
			readOnlyBgndColor: bgColor,
			height: comment.Has?('\n') ? 2 : 1)
		}

	curline: 0
	prevline: 0

	PreviousChange(selected)
		{
		line = .curline
		selected = .findPrevChange(selected, .diffs)
		// loop again to the beginning of the block of changes, do ensurevisible on
		// the first item, then on the last item so that as much of the block
		// of changes as possible is visible
		block_begin = selected
		for (; block_begin > 0 and .diffs[block_begin][1] isnt ""; --block_begin)
			{}
		if (selected >= 0)
			{
			.lists.EnsureVisible(Max(0, block_begin + 1))
			.lists.GotoLine(selected)
			.lists.EnsureVisible(selected)
			line = selected
			}
		else
			Beep()

		.updateSelectedLine(line)
		}
	findPrevChange(selected, diffs)
		{
		if (selected < 0)
			{ Beep(); return }
		// backup untill we find the start of the current change
		for (; selected >= 0 and diffs[selected][1] isnt ""; --selected)
			{}
		// backup untill we find the end of the prev change
		for (; selected >= 0 and diffs[selected][1] is ""; --selected)
			{}
		return selected
		}

	NextChange(selected)
		{
		line = .curline
		selected = .findNextChange(selected, .diffs)
		// loop again to the end of the block of changes, do ensurevisible on
		// the last item, then on the first item so that as much of the block
		// of changes as possible is visible
		block_end = selected
		for (; block_end < .diffs.Size() and .diffs[block_end][1] isnt ""; ++block_end)
			{}
		if (selected < .diffs.Size())
			{
			.lists.EnsureVisible(Min(.diffs.Size() - 1, block_end))
			.lists.GotoLine(selected)
			.lists.EnsureVisible(selected)
			line = selected
			}
		else
			Beep()

		.updateSelectedLine(line)
		}

	findNextChange(selected, diffs)
		{
		// find the end of the current change
		for (; selected < diffs.Size() and diffs[selected][1] isnt ""; ++selected)
			{}
		// find the start of the next change
		for (; selected < diffs.Size() and diffs[selected][1] is ""; ++selected)
			{}
		return selected
		}

	NewValue(value, source)
		{
		if source.Name is 'show_line_difference'
			{
			.toggleShowLineDiff(value)
			.Send('ShowHideLineDiff', .showLine.Get())
			}
		}

	ToggleShowLineDiff(value)
		{
		.toggleShowLineDiff(value)
		.showLine.Set(value)
		}

	toggleShowLineDiff(value)
		{
		if value is true and not .lineDiffExists?()
			.showLineDiff()
		else if value is false and .lineDiffExists?()
			.removeLineDiff()

		.updateSelectedLine(.curline)
		}

	lineDiffExists?()
		{
		return .FindControl('lineDiff') isnt false
		}

	updateSelectedLine(line)
		{
		.curline = Min(line, .diffs.Size() - 1)
		if false isnt lineDiffCtrl = .FindControl('lineDiff')
			{
			if .curline is false
				{
				lineDiffCtrl.Set('')
				return
				}
			diff = .diffs[.curline]
			if diff.Size() is 3/*=diff item size*/ and diff[0] isnt diff[2]
				{
				diffDesc = .GetDiffDesc(diff)
				lineDiffCtrl.Set(
					Opt(.leftChangePre,  ': ') $ diffDesc.left $ '\n' $
					Opt(.rightChangePre, ': ') $ diffDesc.right)
				}
			else
				lineDiffCtrl.Set('')
			}
		else
			{
			if .curline in (false, -1)
				return
			diff = .diffs[.curline]
			}

		.UpdateSelectedIndics(.prevline, .curline, .diffs, .model)

		.lists.RemoveMarker(.prevline)
		.UpdateSelectMarker(.curline, diff[1], .lists)
		.prevline = .curline
		}

	GetDiffDesc(@unused)
		{
		throw 'Must Implement'
		}

	UpdateSelectedIndics(@unused)
		{
		throw 'Must Implement'
		}

	UpdateSelectMarker(@unused)
		{
		throw 'Must Implement'
		}

	Overview_Click(row)
		{
		.lists.GotoLine(row)
		.lists.EnsureVisible(row)
		.updateSelectedLine(row)
		}

	On_Go_To_Definition()
		{
		line = .curline - .LineOffset()
		GoToDefinition(.recName, .lib, line)
		.Window.Result(true)
		}

	LineOffset()
		{
		signs = Object(.newOnRight? ? '>' : '<', '+1', '-2')
		return .diffs[.. .curline].CountIf({ signs.Has?(it[1]) }) - 1
		}

	Goto(name, source) // sent by Addon_go_to_definition
		{
		code = source.Get()
		if name.Prefix?('.') and not ClassHelp.Class?(code)
			{
			GotoLibView(name[1 ..]) // e.g. for rules
			return false
			}
		else
			return .gotoMethod?(name, code, source)
		}

	gotoMethod?(name, code, source)
		{
		if not name.Prefix?('.')
			return 0
		name = name[1 ..]
		if false isnt pos = ClassHelp.FindMethod(code, name)
			.Overview_Click(source.LineFromPosition(pos))
		else if false isnt pos = ClassHelp.FindDotDeclarations(code, name)
			.Overview_Click(source.LineFromPosition(pos))
		else if name[0].Upper?() and
			false isnt x = ClassHelp.FindBaseMethod(.lib, code, name)
			GotoLibView(x.name $ '.' $ name, libs: Object(x.lib))
		else
			return 0
		return true
		}

	Getter_CurLine()
		{
		return .curline
		}

	Getter_Diffs()
		{
		return .diffs
		}

	Getter_Goto()
		{
		return .gotoButton is true ? .lib $ ':' $ .recName : false
		}

	Getter_Lists()
		{
		if not .Member?('DiffBase_lists')
			return #()

		return .lists
		}

	Getter_Ovbar()
		{
		if not .Member?('DiffBase_ovbar')
			return #()

		return .ovbar
		}

	Resize(x, y, w, h)
		{
		selected = .lists.Get(0).LineFromPosition()
		super.Resize(x, y, w, h)
		.lists.GotoLine(selected)
		}

	ClearCallBacks()
		{
		throw 'Must Implement'
		}

	Destroy()
		{
		// BUG doesn't handle if something else subclasses after us
		// TODO handle like ChooseControl.Destroy
		.ClearCallBacks()
		UserSettings.Put('VersionControl-LineDiff', .showLine.Get(), Suneido.User)
		super.Destroy()
		}
	}