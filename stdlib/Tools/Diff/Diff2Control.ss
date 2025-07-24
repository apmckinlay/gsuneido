// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
// used by LibDiff and SvcControl
// TODO if we passed in functions for listNew and listOld
// then refresh could be implemented here (e.g. instead of LibDiff)
DiffBase
	{
	Name: "Diff"
	CallClass(title, listNew, listOld, titleNew, titleOld,
		gotoButton = false, comment = false, commentBgColor = false, newOnRight? = false,
		extraControls = #(Skip), titleNewRight = "", titleOldRight = "")
		{
		Window(
			Object(this, listNew, listOld, "", "", titleNew, titleOld, :gotoButton,
				:comment, :commentBgColor, :newOnRight?, :extraControls, :titleNewRight,
				:titleOldRight),
			:title, keep_placement: 'Diff2Control')
		}

	New(listNew, listOld, lib, recname, titleNew = "ListNew", titleOld = "ListOld",
		.base = false,
		.tags = "12", .gotoButton = false, .comment = false, .refresh = false,
		.commentBgColor = false, .newOnRight? = false, extraControls = #(Skip),
		titleNewRight = "", titleOldRight = "")
		{
		super(listNew, listOld, lib, recname, titleNew, titleOld,
			base, tags, gotoButton, comment, refresh, commentBgColor, newOnRight?,
			:extraControls, :titleNewRight, :titleOldRight)
		}

	SetRedirs()
		{
		// do nothing
		}

	GetControls()
		{
		.listNew = .FindControl(#listNew)
		.listOld = .FindControl(#listOld)
		.titleNew = .FindControl(#titleNew)
		.titleOld = .FindControl(#titleOld)
		.titleNewRight = .FindControl(#titleNewRight)
		.titleOldRight = .FindControl(#titleOldRight)
		}

	SetTitles(titleNew, titleOld, titleNewRight = '', titleOldRight = '')
		{
		.titleNew.Set(titleNew)
		.titleOld.Set(titleOld)
		.titleNewRight.Set(titleNewRight)
		.titleOldRight.Set(titleOldRight)
		}

	OvBars()
		{
		return Object(.FindControl(#ovbar1), .FindControl(#ovbar2))
		}
	ListControls()
		{
		return Object(.listNew, .listOld)
		}
	SetLists(diffs)
		{
		if false is .FindControl('listNew') or false is .FindControl('listOld')
			return

		newCode = diffs.Map({ it[2] }).Join('\r\n')
		oldCode = diffs.Map({ it[0] }).Join('\r\n')
		.listNew.Set(newCode)
		.listOld.Set(oldCode)
		}
	UpdateList(listOld, listNew, base = false)
		{
		if not Instance?(.Lists) or not Instance?(.Ovbar)
			return

		d = .GetDiffs(listOld, listNew, base)
		.SetLists(d.diffs)
		.AddMarkers(d.diffs, .Lists, d.model, .Ovbar)
		.ToFirstDiff()
		}
	Getter_CodeToCheck()
		{
		return .newOnRight? ? .listNew.Get() : .listOld.Get()
		}
	SetProcs()
		{
		.prev1proc = SetWindowProc(.listNew.Hwnd, GWL.WNDPROC, .ListNewproc)
		.prev2proc = SetWindowProc(.listOld.Hwnd, GWL.WNDPROC, .ListOldproc)

		if .base isnt false
			.listNew.SetupMargin()
		}

	AddMarker(indic, row, lists, model)
		{
		tag = " " $ indic[1..]
		if indic is '<'
			lists.AddMarker(row, type: 'Remove')
		else if indic.Prefix?('-')
			{
			lists.AddMarker(row, type: 'Remove')
			.listNew.AddMarginText(row, text: tag.Tr('12', .tags))
			}
		else if indic is '>'
			.listNew.AddMarker(row, type: 'Add')
		else if indic.Prefix?('+')
			{
			.listNew.AddMarker(row, type: 'Add')
			.listNew.AddMarginText(row, text: tag.Tr('12', .tags))
			}
		else
			{
			lists.AddMarker(row, type: 'Modify')
			.listNew.AddMarginText(row, text: tag.Tr('12', .tags))

			super.PaintIndics(.listNew, row, model.GetRowIndics(row).NewToOld)
			super.PaintIndics(.listOld, row, model.GetRowIndics(row).OldToNew)
			}
		}

	ListNewproc(hwnd, msg, wparam, lparam)
		{
		idx = .ListControls().FindIf({ it is .listNew })
		super.ListProc(hwnd, msg, wparam, lparam, idx, .prev1proc)
		}
	ListOldproc(hwnd, msg, wparam, lparam)
		{
		idx = .ListControls().FindIf({ it is .listOld })
		super.ListProc(hwnd, msg, wparam, lparam, idx, .prev2proc)
		}

	DiffPanes()
		{
		return 	Object('Vert'
			Object('Horz'
				Object('Vert'
					Object('Horz'
					Object('Static' size: '+2',
						name: .newOnRight? ? 'titleOld' : 'titleNew'),
					'Fill',
						Object('Static',
							name: .newOnRight? ? 'titleOldRight' : 'titleNewRight',
							justify: 'RIGHT')
					Object('Skip')),
					// Set the position of listNew and listOld based on .newOnRight?
					super.DiffWithBar(.newOnRight? ? 'listOld' : 'listNew', 'ovbar1')
				),
				Object('Vert'
					Object('Horz'
						Object('Static' size: '+2',
							name: .newOnRight? ? 'titleNew' : 'titleOld'),
						'Fill',
						Object('Skip'),
						Object('Static',
							name: .newOnRight? ? 'titleNewRight' : 'titleOldRight',
							justify: 'RIGHT')
					Object('Skip')),
					super.DiffWithBar(.newOnRight? ? 'listNew' : 'listOld', 'ovbar2')
					)
				)
			)
		}

	On_Previous_Change()
		{
		super.PreviousChange(.listNew.LineFromPosition())
		}

	On_Next_Change()
		{
		super.NextChange(.listNew.LineFromPosition())
		}

	GetDiffDesc(diff)
		{
		return Object(
			left: diff[.newOnRight? ? 0 : 2],
			right: diff[.newOnRight? ? 2 : 0])
		}

	UpdateSelectedIndics(prevline, curline, diffs, model)
		{
		prevDiff = diffs[prevline][1]
		curDiff = diffs[curline][1]

		if prevDiff.Prefix?("#")
			{
			prevRowIndics = model.GetRowIndics(prevline)
			super.PaintIndics(.listNew, prevline, prevRowIndics.NewToOld)
			super.PaintIndics(.listOld, prevline, prevRowIndics.OldToNew)
			}

		if curDiff.Prefix?("#")
			{
			curRowIndics = model.GetRowIndics(curline)
			super.PaintIndics(.listNew, curline, curRowIndics.NewToOld, selected?:)
			super.PaintIndics(.listOld, curline, curRowIndics.OldToNew, selected?:)
			}
		}

	UpdateSelectMarker(line, type, lists)
		{
		// The old list is always compared to the new list,
		// and since we want everything displayed in relation to changes in the new list,
		// the Diff output needs to be reversed
		if type is '<' or type.Prefix?('-')
			{
			.listNew.AddMarker(line, '', selected?:)
			.listOld.AddMarker(line, 'Remove', selected?:)
			}
		else if type is '>' or type.Prefix?('+')
			{
			.listNew.AddMarker(line, 'Add', selected?:)
			.listOld.AddMarker(line, '', selected?:)
			}
		else if type.Prefix?('#')
			lists.AddMarker(line, 'Modify', selected?:)
		else
			lists.AddMarker(line, '', selected?:)
		}

	ClearCallBacks()
		{
		SetWindowProc(.listNew.Hwnd, GWL.WNDPROC, .prev1proc)
		ClearCallback(.ListNewproc)
		SetWindowProc(.listOld.Hwnd, GWL.WNDPROC, .prev2proc)
		ClearCallback(.ListOldproc)
		}
	}