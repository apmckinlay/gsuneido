// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Code Recovery'
	CallClass(exit? = false, block = false)
		{
		// If no block is provided, it is an explicit call
		if block is false
			.ctrl(false, exit?)
		else
			try
				block()
			catch (error)
				.ctrl(error, exit?)
		}

	ctrl(error, exit?)
		{
		ToolDialog(0, [this, error, exit?], title: .Title, keep_size: false)
		}

	libList: false
	modifiedList: false
	New(.error, .exit?)
		{
		.libList = .FindControl(#libList)
		.modifiedList = .FindControl(#modifiedList)
		.seperator = LibIO.Seperator() $ '\r\n'
		}

	Controls()
		{
		controls = [#Vert,
			#(Static, 'This control is for emergencies.\r\n\r\n' $
				'As such, Right-Click > Restore is designed to pull directly ' $
				'from the local database.\r\nUse Svc > Compare to ' $
				'ensure database integrity'),
			#Skip
			]
		if .error isnt false
			controls.Add(@.ErrorDisplay(.error))
		controls.Add(
			#Skip,
			[#Horz,
				[#ChooseList, .libLists(), name: #libList]
				],
			#(ListBox, (), xstretch: 1, ystretch: 1, readonly:, multiSelect:,
				font: '@mono', size: '+1', name: modifiedList)
			)
		return controls
		}

	ErrorDisplay(error)
		{
		return [
			[#Static, 'ERROR: ' $ error, size: '+2', weight: #bold, color: CLR.RED]
			#Skip,
			[#Static, FormatCallStack(error.Callstack()), color: CLR.RED]]
		}

	internalLibs: (Test_lib, configlib)
	libLists()
		{
		libs = []
		try
			libs = .libsWithModifications().Remove(@.internalLibs)
		catch (e)
			.print(ERROR: e)
		return libs
		}

	print(@args)
		{
		try
			ServerPrint(@args)
		catch (e)
			{
			errOb = Object('PRINT ERROR: ' $ e, '', 'Attempted to print:')
			for m, v in args
				errOb.Add('\t' $ m $ ': ' $ v)
			.AlertError(.Title, errOb.Join('\r\n'))
			}
		}

	libsWithModifications()
		{
		libs = []
		LibraryTables().Each()
			{
			try
				if not QueryEmpty?(it $
					' where lib_modified isnt "" and group in (-1, -2)')
					libs.Add(it)
			catch (e)
				.print(ERROR: e)
			}
		return libs
		}

	NewValue(lib)
		{
		try
			.newValue(lib)
		catch (e)
			.print(ERROR: e)
		}

	newValue(lib)
		{
		if .libList isnt false and .modifiedList isnt false
			.buildList(lib)
		}

	buildList(lib)
		{
		.modifiedList.DeleteAll()
		QueryApply(lib $ ' where group in (-1, -2) and lib_modified isnt "" sort name')
			{
			prefix = it.lib_committed is '' ? '+ ' : '  '
			prefix = it.group is -1
				? it.lib_committed is ''
					? '+ '
					: '  '
				: '- '
			.modifiedList.AddItem(prefix $ it.name)
			}
		}

	ListBox_ContextMenu(x, y)
		{
		if -1 isnt .modifiedList.GetCurSel()
			try
				ContextMenu(#(Export, Restore, (Restore))).ShowCall(this, x, y)
			catch (e)
				.print(ERROR: e)
		}

	On_Context_Restore()
		{
		lib = .libList.Get()
		for item in .modifiedList.GetAllSelected()
			try
				.restoreItem(lib, item)
			catch (e)
				.print(ERROR: e)
		.NewValue(lib)
		}

	restoreItem(lib, sel)
		{
		prefix = sel.BeforeFirst(' ').Trim()
		name = sel.AfterFirst(' ').Trim()
		if prefix is '+'
			.delete(lib, name)
		else
			.restore(lib, name, deleted: prefix is '-')
		}

	restore(lib, name, deleted)
		{
		QueryApply1(lib, group: deleted ? -2 : -1, :name)
			{
			it.group = -1
			it.text = it.lib_before_text
			it.lib_invalid_text = it.lib_before_text = it.lib_modified = ''
			it.Update()
			}
		.print(RESTORED: lib $ ':' $ name)
		}

	delete(lib, name)
		{
		QueryApply1(lib, group: -1, :name)
			{
			it.Delete()
			}
		.print(REMOVED: lib $ ':' $ name)
		}

	On_Context_Export()
		{
		svcTable = SvcTable(.libList.Get())
		if '' isnt filename = .filename('Export (append) to')
			.modifiedList.GetAllSelected().Each()
				{
				.export(svcTable, it.AfterFirst(' ').Trim(), filename)
				}
		}

	filename(subTitle)
		{
		return SaveFileName(title: .Title $ ' - ' $ subTitle)
		}

	export(svcTable, name, filename)
		{
		lib = svcTable.Table()
		if false is rec = svcTable.Get(name)
			throw 'Export stopped, cannot find record: ' $ lib $ ':' $ name

		content = Object(name, .exportInfo(rec, lib), rec.lib_current_text, .seperator)
		File(filename, 'a')
			{ |f|
			f.Write(content.Join('\r\n'))
			}
		.print(EXPORTED: lib $ ':' $  name)
		}

	exportInfo(rec, lib)
		{
		ob = Object(lib_committed: rec.lib_committed, orig_lib: lib)
		LibIO.SetParent(ob, rec, lib)
		return 'librec_info: ' $ Display(ob)
		}

	On_Import()
		{
		if false is info = LibViewImportRecordControl(.WindowHwnd(), '')
			return
		LibIO.Import(info.fileName, info.lib)
		.print('IMPORTED: ' $ info.fileName)
		}

	Destroy()
		{
		exit? = .exit?
		super.Destroy()
		if exit?
			Exit(true)
		else
			Unload()
		}
	}
