// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
// See also: FindReferencesControl
Controller
	{
	// Title needs to always be the same i.e. not include table name
	// so keep_placement: works
	Title: 'History'
	Name: 'History'
	CallClass(table, name)
		{
		if SvcSettings() is false or table is "" or name is ""
			return false
		IdeTabbedView('VersionHistory', name, table)
		}
	New(.name, .table)
		{
		settings = SvcSettings()
		SvcSocketClient().RetryState()
		.svc = Svc(server: settings.svc_server, local?: settings.svc_local?)
		.svcTable = SvcTable(.table)
		.more = .Vert.Horz.More
		.branches = #()
		data = .get(.name, Date.End())
		.list = .Vert.List
		.list.Set(data)
		if not data.Empty?()
			.list.SetSelection(0)
		.list.VSCROLL(SB.TOP)
		.list.SetMultiSelect(true)
		.diffWindow = #()
		.viewWindow = #()
		if .branches.NotEmpty?()
			.alert("Branching with rename has been detected: " $ .branches.Join(", "))
		}
	Controls:
		(Vert
			(ListStretch, #(When, Who, Type, Name, Comment), columnsSaveName: 'History',
				defWidth: 100, headerSelectPrompt: 'no_prompts')
			(Skip 5)
			(Horz
				Fill
				(Button 'More' pad: 30)
				Skip
				(Button 'View' pad: 30)
				Skip
				(Button 'Diff' defaultButton:, tip: '(same as double click)' pad: 30)
				Skip
				(Button 'Diff to Current')
				Fill)
			)
	List_WantNewRow()
		{ return false }
	List_WantEditField(col)
		{
		return .colName(col) is 'Comment'
			? Object('Editor', xmin: .list.GetColWidth(2), height: 15, readonly:)
			: false
		}
	colName(col)
		{
		return .list.GetColumns()[col]
		}
	get(name, when)
		{
		list = .svc.Get10Before(.table, name, when).Map!({|x|
			Object(When: x.lib_committed, Who: x.id, Comment: x.comment, Type: x.type,
				Name: x.name isnt .name ? x.name : '') })
		if not list.Empty?()
			.when = list.Last().When
		if list.Size() < 10 /* = starting list size*/
			{
			renames = .svc.SearchForRename(.table, name)
			if renames.NotEmpty?()
				{
				if renames.Size() > 1 // should only happen with hash collisions
					{
					.branches = renames
					.more.SetEnabled(false)
					}
				else
					list.Append(.get(renames[0], .when))
				}
			else
				.more.SetEnabled(false)
			}
		return list
		}
	On_More()
		{
		if .branches.NotEmpty?()
			.alert("Branching with rename has been detected: " $ .branches.Join(", "))
		list = Object()
		list.Append(.list.Get())
		list.Append(.get(.getName(), .when))
		.list.Set(list)
		.list.SetSelection(list.Size() - 1)
		}

	On_Diff()
		{
		sel = .list.GetSelection()
		if sel.Size() is 1
			.diffToPrev(sel[0])
		else if sel.Size() is 2
			{
			data = .list.Get()
			sel.Sort!({|x,y| data[x].When < data[y].When })
			prev = data[sel[0]]
			next = data[sel[1]]
			if prev.Type is '-' or next.Type is '-'
				return .alert("Cannot view diff with a deleted item")
			.showdiff2(.getName(prev), prev.When, .getName(next), next.When)
			}
		else
			.alert("Please select one or two rows\n\n" $
				"Selecting one row will compare to the previous version\n" $
				"Selecting two rows will compare those two versions")
		}
	List_DoubleClick(row, col)
		{
		if row is false or .colName(col) is 'Comment'
			return 0
		.diffToPrev(row)
		return true
		}
	List_GetDlgCode(lParam)
		{
		if false isnt (m = MSG(lParam)) and
			(m.message is WM.CHAR or m.message is WM.KEYDOWN) and
			m.wParam is VK.RETURN
			return DLGC.WANTALLKEYS
		return 0
		}
	List_KeyDown(wParam)
		{
		if wParam is VK.RETURN
			.On_Diff()
		}
	diffToPrev(row)
		{
		data = .list.Get()
		x = data[row]
		if row + 1 < data.Size()
			{
			prev = data[row + 1]
			if prev.Type is '-' or x.Type is '-'
				{
				.view(prev.Type is '-' ? x.When : prev.When)
				return
				}
			.showdiff2(.getName(prev), prev.When, .getName(x), x.When,
				x.When.ShortDateTime() $ ' ' $ x.Who $ ' - ' $ x.Comment)
			}
		else
			.view(x.When)
		}
	getversion(name, when)
		{
		when = Date(when)
		x = .svc.GetOld(.table, name, when)
		return x isnt false
			? x
			: []
		}
	showdiff2(name1, when1, name2, when2, comment = '')
		{
		rec1 = .getversion(name1, when1)
		rec2 = .getversion(name2, when2)
		.showdiff(
			'As of ' $ when1.ShortDateTime() $ Opt(' - ', rec1.path), rec1.text,
			'As of ' $ when2.ShortDateTime() $ Opt(' - ', rec2.path), rec2.text,
			comment)
		}
	showdiff(desc1, text1, desc2, text2, comment = "")
		{
		if not .diffWindow.Empty?()
			.diffWindow.Destroy()

		.diffWindow = BookResource?(.name, imagesOnly?:)
			? DiffImageInMemoryControl(text1, text2, desc1, desc2)
			: Diff2CodeControl("Diff: " $ .name, text2, text1, desc2, desc1,
				.table, .name, :comment, newOnRight?:)
		}

	On_View()
		{
		if not Object?(res = .get1())
			return
		if res.Type is '-'
			return .alert("Cannot view diff with a deleted item.")
		.view(res.When)
		}
	view(when)
		{
		if not .viewWindow.Empty?()
			.viewWindow.Destroy()
		name = .getName()
		rec = .getversion(name, when)

		title = Paths.Basename(name) $ ' as of ' $ when.ShortDateTime() $ ' - ' $ rec.path
		viewer = BookResource?(name, imagesOnly?:)
			? ImageViewer
			: CodeViewer
		.viewWindow = viewer(rec.text, table: .table, :name, :title)
		}

	getName(sel = false)
		{
		rowName = sel isnt false ? sel.Name : .get1().Name
		return rowName isnt '' ? rowName : .name
		}

	On_Diff_to_Current()
		{
		if not Object?(x = .get1())
			return
		if x.Type is '-'
			return .alert("Cannot view diff with a deleted item.")
		cur = .svcTable.Get(.name)
		rec = .getversion(.getName(), x.When)
		.showdiff(
			x.When.ShortDateTime() $ ' - ' $ rec.path, rec.text,
			'Current - ' $ rec.path, cur.lib_current_text)
		}

	get1()
		{
		sel = .list.GetSelection()
		if (sel.Size() isnt 1)
			return .alert("Please select one row")
		data = .list.Get()
		return data[sel[0]]
		}

	FindLineHistory(line)
		{
		listData = .list.Get()
		if listData.Size() is 0
			return

		.loadUntilFirstNotFound(line)

		.locateFirstTimeExist(listData, line)
		}

	halfSecondInMs: 500
	loadUntilFirstNotFound(line)
		{
		forever
			{
			listData = .list.Get()
			lastItem = listData.Last()
			oldestRec = .svc.GetOld(.table, .getName(lastItem), lastItem.When)
			if not .hasLine?(oldestRec.text, line)
				break

			if not .more.GetEnabled()
				break

			Thread.Sleep(.halfSecondInMs)
			.On_More()
			.list.SetSelection(.list.GetNumRows() - 1)
			}
		}

	locateFirstTimeExist(listData, line)
		{
		listData = .list.Get()
		last = listData.Size() - 1
		searchIncrement = 10
		for(i = last; i >= last - searchIncrement; i--)
			{
			if i < 0
				{
				.alert('Cannot find the selected line from Version History')
				return
				}
			lineRec = listData[i]
			if lineRec.Type is '-'
				continue
			historyRec = .svc.GetOld(.table, .getName(lineRec), lineRec.When)
			if .hasLine?(historyRec.text, line)
				{
				Thread.Sleep(.halfSecondInMs)
				.list.SetSelection(i)
				return
				}
			}
		.list.SetSelection(last)
		}

	hasLine?(text, line)
		{
		return text =~ '^(?q)' $ line $ '(?-q)$'
		}

	alert(msg)
		{
		Alert(msg, .Title, .Window.Hwnd, MB.ICONINFORMATION)
		return false
		}

	List_ContextMenu(x, y)
		{
		sel = .list.GetSelection()
		if sel is #()
			return
		sel = .list.GetRow(sel[0])

		menufuncs = Object()
		options = Object()
		i = 1
		for cont in GetContributions("VersionHistoryContextMenu")
			{
			options.Add(cont.menuoption)
			menufuncs.Add(cont.menufunc, at: i++)
			}
		if options.Empty?()
			return

		index = ContextMenu(options).Show(.Window.Hwnd, x, y)
		if menufuncs.Member?(index)
			(menufuncs[index])(:sel, source: this)
		}

	Destroy()
		{
		if not .viewWindow.Empty?()
			.viewWindow.Destroy()
		if not .diffWindow.Empty?()
			.diffWindow.Destroy()
		super.Destroy()
		}
	}