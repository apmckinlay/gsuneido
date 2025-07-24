// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
#(
ExtensionPoints:
	(
	(Toolbar)
	)
Contributions:
	(
//	(WorkSpace, Toolbar, 'Diff', D,
//		target: function (@unused) { LibDiffControl() })
	(WorkSpace, Toolbar, 'Preferences', custom_screen,
		target: function (@unused) { IDESettingsControl() })
	(WorkSpace, Toolbar, 'Schema_View', S,
		target: function (@unused) { PersistentWindow(SchemaView) })
	(WorkSpace, Toolbar, 'Test_Runner', T,
		target: function (@unused) { TestRunnerGui() })
	(WorkSpace, Toolbar, 'Version_Control', V,
		target: function (@unused) { SvcControl() })
	(WorkSpace, Toolbar, 'LibraryView', vsplit,
		target: function (@unused) { PersistentWindow(LibViewControl) })
	(WorkSpace, Toolbar, 'QueryView', hsplit,
		target: function (@unused) { QueryViewControl() })
	)
)
