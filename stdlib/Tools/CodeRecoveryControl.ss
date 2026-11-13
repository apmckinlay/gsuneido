// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
/* README:
	This class is designed to function in cases of low level code issues.
	As a result, it is designed to rely on as little external code as possible.
	This reduces the risk of the Code Recovery control failing when an issue occurs.

	For example, this class could make use of: FormatCallStack.
	However, if one of the methods FormatCallStack relied on had issues, it could
	prevent the Code Recovery control from launching (preventing corrective measures).
*/
Controller
	{
	Title: "Code Recovery"
	CallClass(exit? = false, block = false)
		{
		// If no block is provided, it is an explicit call to launch the control
		if block is false
			block = { .launch(exit?, false) }
		try
			block()
		catch (errorBlock)
			try
				.launch(exit?, errorBlock)
			catch (errorCtrl)
				{
				ErrorLog(.Title $ " encountered an issue when attempting to launch" $
					.formatError("Block", errorBlock) $ .formatError(.Title, errorCtrl))
				.shutdown(exit?)
				}
		}

	launch(exit?, error)
		{
		Window([this, exit?, error], title: .Title, w: 750, h: 350,
			exStyle: WS_EX.TOOLWINDOW | WS_EX.STATICEDGE | WS_EX.OVERLAPPEDWINDOW,
			style: WS.CAPTION | WS.SIZEBOX | WS.SYSMENU, border: 5)
		}

	formatError(source, error)
		{
		return "\r\n============================================================\r\n" $
			source $ " Error: " $ error $ "\r\n\t" $ .callstack(error)
		}

	callstack(error, join = "\r\n\t")
		{
		return error.Callstack()[..10 /*= max levels*/].Map!({ it.fn }).Join(join)
		}

	shutdown(exit?)
		{
		if exit?
			Shutdown(alsoServer: Sys.Client?())
		}

	libList:      false
	modifiedList: false
	New(.exit?, .error)
		{
		.libList = .FindControl(#libList)
		.modifiedList = .FindControl(#modifiedList)
		.seperator = LibIO.Seperator() $ "\r\n"
		}

	Controls()
		{
		controls = [#Vert,
			#(Static, "This control is for emergencies."),
			#Skip,
			#(Static,
				'As such, Right-Click > Restore is designed to pull directly ' $
				'from the local database.'),
			#(Static,
				'Be sure to use Version Control > Compare to ensure database ' $
				'integrity after restoring.'),
			#Skip]
		if .error isnt false
			controls.Add(@.ErrorDisplay(.error))
		controls.Add(
			#Skip,
			[#Horz,
				[#ChooseList, .libLists(), name: #libList]],
			#(ListBox, (), xstretch: 1, ystretch: 1, readonly:, multiSelect:,
				font: "@mono", size: "+1", name: modifiedList))
		return controls
		}

	ErrorDisplay(error)
		{
		return [
			[#Static, "ERROR: " $ error, size: "+2", weight: #bold, color: CLR.RED],
			[#Static, .callstack(error, "\r\n"), color: CLR.RED]]
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
			errOb = ["PRINT ERROR: " $ e, "", "Attempted to print:"]
			for m, v in args
				errOb.Add('\t' $ m $ ": " $ v)
			.AlertError(.Title, errOb.Join("\r\n"))
			}
		}

	libsWithModifications()
		{
		libs = []
		LibraryTables().Each()
			{
			try
				if not QueryEmpty?(
					it $ ' where lib_modified isnt "" and group in (-1, -2)')
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
			prefix = it.lib_committed is "" ? "+ " : "  "
			prefix = it.group is -1
				? it.lib_committed is ""
					? "+ "
					: "  "
				: "- "
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
			it.lib_invalid_text = it.lib_before_text = it.lib_modified = ""
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
		if "" isnt filename = .filename("Export (append) to")
			.modifiedList.GetAllSelected().Each()
				{
				.export(svcTable, it.AfterFirst(' ').Trim(), filename)
				}
		}

	filename(subTitle)
		{
		return SaveFileName(title: .Title $ " - " $ subTitle)
		}

	export(svcTable, name, filename)
		{
		lib = svcTable.Table()
		if false is rec = svcTable.Get(name)
			throw "Export stopped, cannot find record: " $ lib $ ':' $ name

		content = [name, .exportInfo(rec, lib), rec.lib_current_text, .seperator]
		File(filename, 'a')
			{|f|
			f.Write(content.Join("\r\n"))
			}
		.print(EXPORTED: lib $ ':' $ name)
		}

	exportInfo(rec, lib)
		{
		ob = Object(lib_committed: rec.lib_committed, orig_lib: lib)
		LibIO.SetParent(ob, rec, lib)
		return "librec_info: " $ Display(ob)
		}

	On_Import()
		{
		if false is info = LibViewImportRecordControl(.WindowHwnd(), "")
			return
		LibIO.Import(info.fileName, info.lib)
		.print("IMPORTED: " $ info.fileName)
		}

	Destroy()
		{
		Unload()
		exit? = .exit?
		super.Destroy()
		.shutdown(exit?)
		}
	}
