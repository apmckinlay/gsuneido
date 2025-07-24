// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
DiffBase
	{
	Name: "MergeDiff"
	CallClass(title, listMaster, listCurrent,
		titleCurrent, titleMaster, titleMerge = 'MERGED',
		gotoButton = false, comment = '', commentBgColor = false, newOnRight? = false,
		extraControls = #(Skip))
		{
		Window(
			Object(this, listMaster, listCurrent,
				titleCurrent, titleMerge, titleMaster, :gotoButton, :comment,
				:commentBgColor, :newOnRight?, :extraControls),
			:title, keep_placement: 'Diff2Control')
		}
	New(listMaster, listCurrent, lib, recname,
		titleCurrent, titleMaster, .titleMerge = 'MERGED',
		.base = false,
		.tags = "12", .gotoButton = false, .comment = false, .refresh = false,
		.commentBgColor = false, .newOnRight? = false, extraControls = #(Skip))
		{
		super(listMaster, listCurrent, lib, recname, titleCurrent, titleMaster,
			base, tags, gotoButton, comment, refresh, commentBgColor, newOnRight?,
			:extraControls)
		}

	SetRedirs()
		{
		.Redir('On_Use_Master')
		.Redir('On_Merge')
		}

	GetControls()
		{
		.listCurrent = .FindControl(#listCurrent)
		.listMerge = .FindControl(#listMerge)
		.listMaster = .FindControl(#listMaster)
		.titleCurrentCtrl = .FindControl(#titleCurrent)
		.titleMergeCtrl = .FindControl(#titleMerge)
		.titleMasterCtrl = .FindControl(#titleMaster)
		}

	SetTitles(titleCurrent, titleMaster)
		{
		.titleCurrentCtrl.Set(titleCurrent)
		.titleMergeCtrl.Set(.titleMerge)
		.titleMasterCtrl.Set(titleMaster)
		}

	SetWarn(msg)
		{
		.FindControl(#mergeWarn).Set(msg)
		}

	OvBars()
		{
		return Object(.FindControl(#ovbar1),
			.FindControl(#ovbar2),
			.FindControl(#ovbar3))
		}
	ListControls()
		{
		return Object(.listCurrent, .listMerge, .listMaster)
		}
	SetLists(diffs)
		{
		.listCurrent.Set(diffs.Map({ .getDiff(it, '2') }).Join('\r\n'))
		.listMaster.Set(diffs.Map({ .getDiff(it, '1')  }).Join('\r\n'))
		mergedCode = diffs.Map({ it[2] }).Join('\r\n')
		.listMerge.Set(mergedCode)
		}

	getDiff(diff, idx)
		{
		// if + or # with no idx, then exists in both local and master, but not original
		// (i.e. both added the same line, or changed the same line in the same way)
		if .addModify(diff, idx)
			return diff[2]

		otherIdx = idx is '1' ? '2' : '1'
		if diff[1] is '-'$otherIdx or diff[1] is '#'$otherIdx
			return diff[0]

		return ''
		}
	addModify(diff, idx)
		{
		indic = diff[1].Extract('\d$')
		return diff[1] is "" or
			((diff[1].Prefix?('+') or diff[1].Prefix?('#')) and
				(indic is idx or indic is false))
		}

	Getter_CodeToCheck()
		{
		return .listMerge.Get()
		}

	SetProcs()
		{
		.prev1proc = SetWindowProc(.listCurrent.Hwnd, GWL.WNDPROC, .ListCurproc)
		.prev2proc = SetWindowProc(.listMerge.Hwnd, GWL.WNDPROC, .ListMergeproc)
		.prev3proc = SetWindowProc(.listMaster.Hwnd, GWL.WNDPROC, .ListMasterproc)

		.listMerge.SetupMargin()
		.listMerge.SetFocus()
		}

	AddMarker(indic, row, lists, model)
		{
		tag = " " $ indic[1..]
		text = tag.Tr('12', .tags)
		if indic is '<'
			lists.AddMarker(row, type: 'Remove')
		else if indic.Prefix?('-')
			{
			lists.AddMarker(row, type: 'Remove')
			.listMerge.AddMarginText(row, :text)
			}
		else if indic is '>'
			.listCurrent.AddMarker(row, type: 'Add')
		else if indic.Prefix?('+')
			.addListMarker(indic, row, type: 'Add', :text)
		else
			{
			lists.AddMarker(row, type: 'Modify')
			list = indic.Suffix?('1') ? .listMaster : .listCurrent
			list.AddMarginText(row, text: tag.Tr('12', .tags))

			super.PaintIndics(.listCurrent, row, model.GetRowIndics(row).NewToOld)
			super.PaintIndics(.listMaster, row, model.GetRowIndics(row).OldToNew)
			}
		}

	addListMarker(indic, row, type, text)
		{
		.listMerge.AddMarker(row, :type)
		.listMerge.AddMarginText(row, :text)
		if indic[-1] is indic //no suffix
			{
			.listMaster.AddMarker(row, :type)
			.listCurrent.AddMarker(row, :type)
			}
		else
			{
			list = indic.Suffix?('1') ? .listMaster : .listCurrent
			list.AddMarker(row, :type)
			}
		}

	ListCurproc(hwnd, msg, wparam, lparam)
		{
		idx = .ListControls().FindIf({ it is .listCurrent })
		super.ListProc(hwnd, msg, wparam, lparam, idx, .prev1proc)
		}
	ListMergeproc(hwnd, msg, wparam, lparam)
		{
		idx = .ListControls().FindIf({ it is .listMerge })
		super.ListProc(hwnd, msg, wparam, lparam, idx, .prev2proc)
		}
	ListMasterproc(hwnd, msg, wparam, lparam)
		{
		idx = .ListControls().FindIf({ it is .listMaster })
		super.ListProc(hwnd, msg, wparam, lparam, idx, .prev3proc)
		}

	DiffPanes()
		{
		return Object('Vert'
			Object('HorzEqual'
				Object('Vert'
				Object('Static' size: '+2', xstretch: 1,
					name: 'titleCurrent'),
				super.DiffWithBar('listCurrent', 'ovbar1')),
				Object('Vert'
				Object('Horz'
					Object('Static' size: '+2', xstretch: 1, name: 'titleMerge')
					Object('Static' size: '+2', xstretch: 1, name: 'mergeWarn',
						weight: 'bold', color: CLR.RED)),
				super.DiffWithBar('listMerge', 'ovbar2'))
				Object('Vert'
				Object('Static' size: '+2', xstretch: 1,
					name: 'titleMaster'),
				super.DiffWithBar('listMaster', 'ovbar3'))
				)
			)
		}

	On_Previous_Change()
		{
		if .focusedDiff isnt false
			super.PreviousChange(.focusedDiff.LineFromPosition())
		}

	On_Next_Change()
		{
		if .focusedDiff isnt false
			super.NextChange(.focusedDiff.LineFromPosition())
		}

	GetDiffDesc(diff)
		{
		// for diff record: 0 means before change, 2 means after change
		// for diff[1]: 1 means master made change, 2 means local made change,
		// 		false means either both made same change, or neither made change
		if false is ind = diff[1].Extract("\d$")
			{
			// if neither, or both
			ind = .newOnRight? ? 0 : 2
			}
		else
			{
			// if local made the change then show before to master
			// if master made the change then show after to master
			if ind is '2'
				ind = .newOnRight? ? 0 : 2
			else
				ind = .newOnRight? ? 2 : 0
			}
		otherInd = ind is 0 ? 2 : 0
		return Object(left: diff[ind], right: diff[otherInd])
		}

	UpdateSelectedIndics(prevline, curline, diffs, model)
		{
		prevDiff = diffs[prevline][1]
		curDiff = diffs[curline][1]
		if prevDiff.Prefix?("#")
			{
			prevRowIndics = model.GetRowIndics(prevline)
			super.PaintIndics(.listCurrent, prevline, prevRowIndics.NewToOld)
			super.PaintIndics(.listMaster, prevline, prevRowIndics.OldToNew)
			}

		if curDiff.Prefix?("#")
			{
			curRowIndics = model.GetRowIndics(curline)
			super.PaintIndics(.listCurrent, curline, curRowIndics.NewToOld, selected?:)
			super.PaintIndics(.listMaster, curline, curRowIndics.OldToNew, selected?:)
			}
		}

	UpdateSelectMarker(line, type, lists)
		{
		// The old list is always compared to the new list,
		// and since we want everything displayed in relation to changes in the new list,
		// the Diff output needs to be reversed
		if type is '<' or type.Prefix?('-')
			{
			.listCurrent.AddMarker(line, '', selected?:)
			.listMaster.AddMarker(line, 'Remove', selected?:)
			.listMerge.AddMarker(line, 'Remove', selected?:)
			}
		else if type is '>' or type.Prefix?('+')
			{
			.listCurrent.AddMarker(line, 'Add', selected?:)
			.listMaster.AddMarker(line, '', selected?:)
			.listMerge.AddMarker(line, 'Add', selected?:)
			}
		else if type.Prefix?('#')
			lists.AddMarker(line, 'Modify', selected?:)
		else
			lists.AddMarker(line, '', selected?:)
		}

	focusedDiff: false
	Scintilla_SetFocus(source)
		{
		.focusedDiff = source
		}

	ClearCallBacks()
		{
		SetWindowProc(.listCurrent.Hwnd, GWL.WNDPROC, .prev1proc)
		ClearCallback(.ListCurproc)
		SetWindowProc(.listMerge.Hwnd, GWL.WNDPROC, .prev2proc)
		ClearCallback(.ListMergeproc)
		SetWindowProc(.listMaster.Hwnd, GWL.WNDPROC, .prev3proc)
		ClearCallback(.ListMasterproc)
		}
	}
