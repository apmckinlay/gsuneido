// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'View Item as of'
	CallClass(lib, name, curItem, editor)
		{
		ToolDialog(0, Object(this, lib, name, curItem, editor))
		}

	New(.lib, .name, .curItem, .editor)
		{
		.list.SetReadOnly(true)
		if false is history = QueryHistory.GUI(.query, parent: .WindowHwnd())
			{
			.Window.Result(0)
			return
			}
		.list.AddRows(history.Reverse!())
		.list.SelectRow(0)
		.rb.Picked('Diff to Current')
		}

	getter_query()
		{ return .query = .lib $ ' where group in (-1, -2) and name is "' $ .name $ '"' }

	getter_list()
		{ return .list = .FindControl('List') }

	getter_rb()
		{ return .rb = .FindControl('RadioButtons') }

	getter_diffPane()
		{ return .diffPane = .FindControl('diffPane') }

	Controls()
		{
		return Object('Vert'
			#(HorzSplit,
				(Vert
					(ListStretch, #(tran_asof), columnsSaveName: 'ViewItemHistory',
						resetColumns:),
					ymin: 100, ystretch: 1, xstretch: 1, xmin: 100)
				(Horz (Pane Fill) name: 'diffPane' ystretch: 1, xstretch: 4),
				splitSaveName: 'ViewItemHistory')
			#EtchedLine
			#(Skip medium:)
			#(Horz
				Fill
				(RadioButtons, 'View', 'Diff to Previous', 'Diff to Current', horz:)
				(Skip medium:)
				(MenuButton 'Restore' ('Restore'))
				Fill))
		}

	selectedRec: false
	List_Selection(selection)
		{
		if selection is false
			return

		.update_diff(selection[0])
		}

	update_diff(selection)
		{
		.selectedRec = QueryHistory(.query, asof = .getRowAsof(selection)).Extract(0)
		listOld = .selectedRec.lib_current_text
		titleOld = .title(.selectedRec, asof)
		compareOb = .compareTo(selection)
		view_ob = compareOb.list is false
			? Object('Vert', Object('DisplayCode', set: listOld), xmin: .diffXmin)
			: Object('Diff2', compareOb.list, listOld, .lib, .name,
				:titleOld, titleNew: compareOb.title, newOnRight?:)

		.diffPane.RemoveAll()
		.diffPane.Insert(0, view_ob)
		}

	getRowAsof(idx)
		{ return .list.GetRow(idx).tran_asof }

	titleLength: 20 // Used to set titles to the same length to align properly
	title(rec, asof)
		{
		title = 'As of ' $ asof.ShortDateTime()
		if rec isnt false and rec.group is -2
			title $= ' (DELETED)'
		return title.RightFill(.titleLength, '\t')
		}

	compareTo(selection)
		{
		title = list = false
		if 'Diff to Current' is radioButton = .rb.Get()
			{
			title = 'Current'.RightFill(.titleLength, '\t')
			list = .curItem
			}
		else if radioButton is 'Diff to Previous' and ++selection < .list.GetNumRows()
			{
			asof = .getRowAsof(selection)
			if false isnt rec = QueryHistory(.query, :asof).Extract(0, false)
				list = rec.lib_current_text
			title = .title(rec, asof)
			}
		return [:title, :list]
		}

	/* Need the larger of 2 Xmins for Diff2 / DisplayCode (it's Diff's). It should come
	 * back as 980, but it is calculated dynamically based on the Width member of
	 * ScintillaAddonsControl (see the .calcXmin() method there).
	 * 483(ScintillaAddonsControl) + 7(OverviewBar) for 2 interior controls (total 980)
	 */
	getter_diffXmin()
		{
		if false is ctrl = .FindControl('Diff')
			return 980 /*= default size if Diff control is not constructed */
		return .diffXmin = ctrl.Xmin
		}

	GetField(unused)
		{ return 'Diff to Current' }

	NewValue(unused)
		{ .update_diff(.list.GetSelection()[0]) }

	On_Restore(unused)
		{
		if false is .selectedRec
			return
		if .selectedRec.group is -1
			{
			.editor.PasteOverAll(.selectedRec.lib_current_text)
			.Window.Result(0)
			}
		else
			.AlertWarn(.Title, 'Cannot restore record to deleted state')
		}
	}
