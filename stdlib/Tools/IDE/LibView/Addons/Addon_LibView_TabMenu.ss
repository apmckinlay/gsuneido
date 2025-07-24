// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
LibViewAddon
	{
	CloseTab(tab = 0)
		{ .Explorer.Tab_Close(tab) }

	TabMenuOptions(options)
		{
		options.Add('Move To Other Library View', 'Copy Library And Name', 'Copy Name',
			'', 'Find References', 'Diff to Overridden', 'Diff to Original',
			'Restore to Original', #('Restore %1'), 'Version History', '',
			'Go To Documentation', 'Edit Documentation', 'Go To Associated Test', '',
			'Run Associated Test', 'Debug %1', '&Run %1', '&Profile %1', '', 'Export %1')
		}

	TabMenu_CopyLibraryAndName(tabData)
		{
		name = tabData.name
		table = tabData.table
		ClipboardWriteString(table is name ? table : table $ Opt(':', name))
		}

	TabMenu_CopyName(tabData)
		{ ClipboardWriteString(tabData.name is '' ? tabData.table : tabData.name) }

	TabMenu_Debug(tabData)
		{ LibView_DebugTest(.Parent, tabData.table, tabData.name) }

	TabMenu_EditDocumentation(tabData)
		{
		if false isnt page = QueryFirst('suneidoc
			where name is ' $ Display(tabData.name) $ ' sort path')
			OpenBook('suneidoc', 'suneidoc' $ page.path  $ '/' $ page.name, bookedit?:)
		}

	TabMenu_Export(tabData)
		{
		if not QueryEmpty?(tabData.table, name: tabData.name, group: -1)
			.LibExportFile(tabData.table, tabData.name, tabData.item)
		}

	TabMenu_FindReferences(tabData)
		{ FindReferencesControl(tabData.name) }

	TabMenu_DifftoOverridden(tabData)
		{
		if false isnt rec = Query1(tabData.table, name: tabData.name, group: -1)
			LibDiffOverriddenControl(tabData.name, tabData.table, rec.lib_current_text)
		}

	TabMenu_DifftoOriginal(tabData)
		{
		if DiffToOriginalControl(tabData.table, tabData.name) is 'Restore'
			SvcTable(tabData.table).Restore(tabData.name)
		}

	TabMenu_Restore(tabData)
		{ SvcTable(tabData.table).Restore(tabData.name) }

	TabMenu_GoToAssociatedTest(tabData)
		{
		name = LibraryTags.RemoveTagsFromName(tabData.name.Tr('?'))
		items = Object(name $ 'Test', name $ '_Test')
		libs = Libraries()
		list = #()
		for item in items
			{
			list = Gotofind(item, libs)
			if not list.Empty?()
				break
			}
		name = list.Empty?() ? name $ '_Test' : list[0].AfterLast('/')
		GotoLibView(name, .Editor, :libs, libview: .Parent, path: tabData.path)
		}

	TabMenu_GoToDocumentation(tabData)
		{ GotoDocumentation(tabData.name) }

	Tab_SupportSeparate?()
		{
		return true
		}

	Tab_Separate(tab)
		{
		tabData = .Explorer.Tabs.GetData(tab)
		tabData.idx = tab
		.TabMenu_MoveToOtherLibraryView(tabData)
		}

	TabMenu_MoveToOtherLibraryView(tabData)
		{
		// Gets a position of best fit
		line = .Editor isnt false
			? .Editor.GetFirstVisibleLine() + .Editor.LinesOnScreen() / 2
			: 0
		.CloseTab(tabData.idx)
		libview = GotoPersistentWindow('LibViewControl',
			LibViewControl, except: .Window.Hwnd)
		libview.GotoPathLine(tabData.path.Tr('()'), line)
		if .Explorer.Tabs.Count() is 0
			.On_Close()
		}

	TabMenu_Run(unused)
		{ .Try_run("all") }

	TabMenu_Profile(tabData)
		{
		if tabData.name.Suffix?("Test")
			RunWithProfile()
				{ TestRunner.Run1(tabData.name, observer: TestObserver) }
		else
			.Try_run('all', quiet?:,
				wrapper: function (block){ return { RunWithProfile(block); #() }})
		}

	TabMenu_RunAssociatedTest(tabData)
		{
		observer = RunAssociatedTests([tabData.table $ '/' $ tabData.name], .Parent)
		.AlertTestResult(observer)
		}

	TabMenu_VersionHistory(tabData)
		{ VersionHistoryControl(tabData.table, tabData.name) }
	}
