// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
LibViewAddon
	{
	ContextMenu(item)
		{
		recordSelected? = .Explorer.Tree.Selection.Size() is 1 and
			not .Explorer.Container?(item)
		return .buildContextMenu(name: .Explorer.Tree.GetName(item), :recordSelected?
			rootSelected?: .Explorer.RootSelected?() is true)
		}

	buildContextMenu(name, recordSelected?, rootSelected?)
		{
		// Base menu must always be new in order to avoid issues with
		// "past by reference" object modification.
		menu = []
		if recordSelected?
			menu.Add(@.recordOptions())
		if rootSelected?
			menu.Add(@.rootOptions(name))
		menu.Add(@.newOptions())
		return menu
		}

	recordOptions()
		{
		return [
			[name: '&Run', def:, order: 1],
			[name: '', order: 2],
			[name: '', order: 80],
			[name: 'Export Record', order: 81]]
		}

	newOptions()
		{
		return [
			[name: '&New', order: 40],
			['&Folder', '&Item', order: 41]]
		}

	rootOptions(name)
		{
		return [
			[name: '&Delete', order: 30],
			['&Delete Library', order: 31],
			[name: '', order: 90],
			[name: 'Dump', order: 91],
			[name: 'Import Records...', order: 92],
			[name: name.Has?('(') ? 'Use' : 'Unuse', order: 93],
			[name: 'Undelete...', order: 94],
			[name: '', order: 110],
			[name: 'Check Records', order: 111]]
		}

	On_Context_Delete_Library()
		{ .Explorer.On_Delete_Item(allowLibraryDelete?:) }

	On_Context_Export_Record()
		{ .LibExportFile() }

	getLibDisplayName()
		{ return .Explorer.Tree.GetName(.Explorer.Tree.Selection[0]) }

	getLibName()
		{ return .getLibDisplayName().Tr('()') }

	On_Context_Run()
		{ .On_Run() }

	On_Context_Unuse()
		{ .On_Unuse_Library(.getLibDisplayName()) }

	On_Context_Use()
		{ .On_Use_Library(.getLibName()) }

	On_Context_Undelete()
		{
		lib = .getLibName()
		if false isnt x = .getUndeleteRec(lib)
			.undelete(lib, x)
		}

	getUndeleteRec(lib)
		{
		deleted = Object().Set_default([tran_asof: ''])
		if false is history = QueryHistory.GUI(lib $ ' where group is -1')
			return false
		for rec in history.Reverse!()
			{
			// QueryHistory searches from the most recent record backwards
			// (But the result is sorted by tran_asof, so we need .Reverse!())
			// This means once we've identified a record as deleted it will have
			// the most recent tran_asof
			if deleted.Member?(rec.name)
				continue
			if QueryEmpty?(lib $ ' where group in (-1, -2)', name: rec.name)
				deleted[rec.name] = rec
			}
		if false is result = .libViewUndeleteControl(lib, deleted.Values())
			return false
		return QueryHistory(lib $
			' where group is -1 and name is ' $ Display(result.name),
			asof: result.tran_asof)[0]
		}

	libViewUndeleteControl(lib, deleted)
		{ return ToolDialog(.Window.Hwnd, Object(LibViewUndeleteControl, lib deleted)) }

	undelete(lib, x)
		{
		if QueryEmpty?(lib, num: x.parent)
			x.parent = 0
		svcTable = SvcTable(lib)
		// If no longer staged for deletion then the delete has been committed.
		// Clear lib_committed and treat as a new record
		if x.lib_committed isnt '' and false is svcTable.Get(x.name, deleted:)
			x.lib_committed = ''
		svcTable.Output(x)
		.ResetCtrls()
		}

	On_Context_Check_Records()
		{
		CheckRecordsControl(.getLibName())
		}
	}
