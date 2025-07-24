// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	title: Compare
	CallClass(svc, table)
		{
		ToolDialog(_hwnd, [this, svc, table],
			title: .title $ ' ' $ table, keep_size: 'SvcCompare', closeButton?: false)
		}

	Controls: #(Vert,
		(GridControl, (
			((Static, '#'),
				(Static, 'record exists both locally and in version control, ' $
					'but they are different')),
			((Static, 'M'),
				(Static, 'record only exists in version control Master, not locally')),
			((Static, 'L'),
				(Static, 'record only exists Locally, not in version control')))
			overlap:)
		Skip,
		(ListBox, (), xstretch: 1, ystretch: 1, readonly:, multiSelect:,
			font: '@mono', size: '+1'),
		(Skip, medium:),
		(Horz,
			(Static 'Double click to view, or right click for more options')
			Fill,
			(Button, Close)
		))
	New(.svc, .table)
		{
		.list = .FindControl('ListBox')
		.updateList()
		}

	updateList()
		{
		if .list.GetCount() > 0
			.list.DeleteAll()
		.svc.Compare(.svcTable = SvcTable(.table)).Each(.list.AddItem)
		}

	ListBoxDoubleClick(i)
		{
		s = .list.GetText(i)
		name = s[2 ..]
		switch s[0]
			{
		case '#' :
			.viewDiff(name)
		case 'M' :
			.viewMaster(name)
		case 'L' :
			.goToLocal(name)
			}
		}

	viewDiff(name)
		{
		local = .svcTable.Get(name)
		if false is local = .getLocal(name)
			.AlertWarn('View Local', `Can't find ` $ name)
		else if false is master = .svc.Get(.table, name)
			.AlertWarn('View Master', `Can't get ` $ name)
		else
			{
			titles = .diffTitles(.svcTable.Type, local, master, name)
			Diff2CodeControl('Diff: ' $ .table $ ' : ' $ name,
				local.text, master.text, titles.local, titles.master, .table, name,
				keep_size: 'SvcCompare diff')
			}
		}

	getLocal(name)
		{
		if false is local = .svcTable.Get(name)
			return false
		local.text = .svcTable.Type is 'lib'
			? local.lib_current_text
			: local.text
		return local
		}

	diffTitles(type, local, master, name)
		{
		localTitle = type is 'lib'
			? Opt(local.path, '/') $ name
			: name
		masterTitle = type is 'lib'
			? Opt(master.path, '/') $ name
			: master.name
		return [local: 'LOCAL ' $ localTitle, master: 'MASTER ' $ masterTitle]
		}

	viewMaster(name)
		{
		if false is x = .svc.Get(.table, name)
			.AlertWarn('View Master', `Can't get ` $ name)
		else
			CodeViewer(x.text, table: .table, :name,
				title: 'MASTER ' $ .table $ ' : ' $ name)
		}

	goToLocal(name)
		{
		if .svcTable.Type is 'lib'
			GotoLibView(.table $ ':' $ name)
		else
			OpenBook(.table, .table $ name, bookedit?:)
		}

	ListBox_ContextMenu(x, y)
		{
		menu = Object()
		if '' isnt s = .list.Get()
			.addItemContextMenuOptions(menu, s)
		if menu.NotEmpty?()
			menu.Add('')
		menu.Add('Sync')
		ContextMenu(menu).ShowCall(this, x, y)
		}

	addItemContextMenuOptions(menu, s)
		{
		if s[0] is '#'
			menu.Add('Diff')
		if s[0] in ('M', '#')
			menu.Add('View Master')
		if s[0] in ('L', '#')
			menu.Add('Go To Local Definition')
		if s[0] in ('M', '#') or .list.GetSelected(all?:) > 1
			menu.Add('', 'Overwrite')
		}

	On_Context_Diff()
		{
		.viewDiff(.curname())
		}

	curname()
		{
		return .list.Get()[2 ..]
		}

	On_Context_View_Master()
		{
		.viewMaster(.curname())
		}

	On_Context_Go_To_Local_Definition()
		{
		.goToLocal(.curname())
		}

	On_Context_Overwrite()
		{
		selected = .list.GetAllSelected()
		if selected.NotEmpty?()
			.overwriteSelected(selected)
		return 0
		}

	overwriteSelected(selected)
		{
		if .svc.Outstanding?([.table])
			.AlertError(.title,
				'Unable to overwrite record: Local database is not up to date.\r\n' $
				'\r\nPlease update ' $ .table $
				' prior to overwriting individiual records.')
		else if selected.Size() is 1 or .overwrite?('All selected records')
			{
			.svc.Overwrite(.discrepancies(selected))
			.updateList()
			.modified = true
			}
		}

	discrepancies(overwrite)
		{
		discrepancies = Object()
		overwrite.Each({ discrepancies.Add([table: .table, name: it[2 ..]]) })
		return discrepancies.Sort!(By(#name))
		}

	overwrite?(prefix)
		{
		return YesNo(prefix $ ' will be overwritten with their\r\n' $
			'most recent version from SVC.\r\n\r\nContinue?', .title,
			flags: MB.ICONWARNING)
		}

	On_Context_Sync()
		{
		if not .overwrite?('All records')
			return
		overwrite = Object()
		for i in .. .list.GetCount()
			overwrite.Add(.list.GetText(i))
		SvcSyncTables.Sync(.table, .svc, overwrite)
		.updateList()
		.modified = true
		}

	modified: false
	On_Close()
		{
		.Window.Result(.modified)
		}
	}
