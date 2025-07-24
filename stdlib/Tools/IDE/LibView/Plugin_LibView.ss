// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
#(
ExtensionPoints:
	(
	(Tools)
	(New)
	)
Contributions:
	(
	(LibView, Tools, Find_References_to_Current, R, "Ctrl+R",
		target: function (libview) { FindReferencesControl(libview.CurrentName()) })
	(LibView, Tools, Go_To_Documentation, D,
		target: function (libview) { GotoDocumentation(libview.CurrentName()) })
	(LibView, Tools, Run_Associated_Tests, T, "Ctrl+T",
		target: function (libview) { RunAssociatedTests.FromLibView(libview) } )
	(LibView, Tools, Debug_Test, K, "Ctrl+K",
		target: function (libview)
			{
			LibView_DebugTest(libview, libview.CurrentTable(), libview.CurrentName())
			})
	(LibView, Tools, Version_History, H,
		target: function (libview)
			{
			libview.Save() // for diff to current
			VersionHistoryControl(libview.CurrentTable(), libview.CurrentName())
			} )
	(LibView, Tools, Quality_Checker_Window, check,
		target: function (libview)
			{
			libview.Save() // for diff to current
			Window(Object('QualityChecker', libview),
				w: 800, h: 1000, keep_placement:)
			} )
	(LibView, Tools, Version_Control_Settings, "",
		target: function () { SvcSettings(openDialog:) } )
	(LibView, Tools, Svc_Statistics, "",
		target: function(libview) { SvcStatsDisplay(libview.CurrentName(),
			libview.CurrentTable()) })
	(LibView, Tools, Delete_Empty_Folders, "",
		target: function(libview /*unused*/)
			{
			if LibEmptyFolders.RemoveAll()  > 0
				SvcTable.Publish('TreeChange', type: 'lib', force:)
			})

	(LibView, New, 'function',
'function ()
	{
	}')
	(LibView, New, 'class',
'class
	{
	}')
	(LibView, New, 'Test',
'Test
	{
	Test_one()
		{
		}
	}')
	(LibView, New, 'Plugin',
'#(
ExtensionPoints:
	(
	(myExtensionPoint1)
	(myExtensionPoint2)
	)
Contributions:
	(
	(MyPlugin, myExtensionPoint1, "...")
	(MyPlugin, myExtensionPoint1, "...")
	(MyPlugin, myExtensionPoint2, "...")
	(OtherPlugin, itsExtensionPoint, "...")
	)
)')
	)
)
