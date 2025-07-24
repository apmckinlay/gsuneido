// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
#(
Contributions:
	(
	(UI, menu, menu: IDE, order: 8)

	(UI, action, menu: IDE, name: ('Views'),
		target: IdeActions)
	(UI, viewsubmenu, menu: IDE, name: "&WorkSpace",
		target: function (@unused) { PersistentWindow(WorkSpaceControl) })
	(UI, viewsubmenu, menu: IDE, name: "&LibraryView",
		target: function (@unused) { PersistentWindow(LibViewControl) })
	(UI, viewsubmenu, menu: IDE, name: "&QueryView",
		target: function (@unused) { QueryViewControl() })
	(UI, viewsubmenu, menu: IDE, name: "&SchemaView",
		target: function (@unused) { PersistentWindow(SchemaView) })
	(UI, viewsubmenu, menu: IDE, name: "&ClassView",
		target: function (@unused) { PersistentWindow(ClassBrowser) })
	(UI, viewsubmenu, menu: IDE, name: "&PluginsView",
		target: function (@unused) { PersistentWindow(PluginsView) })

	(UI, action, menu: IDE, name: ("&Edit a Book"),
		target: IdeActions)
	(UI, action, menu: IDE, name: ('&Open a Book'),
		target: IdeActions)
	(UI, action, menu: IDE, name: "&MultiView a Query...",
		target: IdeActions)
	(UI, action, menu: IDE, name: "&TestRunner",
		target: function (@unused) { TestRunnerGui() })
	(UI, action, menu: IDE, name: "&Version Control",
		target: function (@unused) { SvcControl() })
	(UI, action, menu: IDE, name: "Show/Hide Override",
		target: function (@unused) { ShowHideDialog() })
	(UI, action, menu: IDE, name: 'Dump Libraries',
		target: function (@unused) { DumpLibrariesControl() })
	(UI, action, menu: IDE, name: 'Reset Caches',
		target: function (@unused) { ResetCaches() })
	(UI, action, menu: IDE, name: 'Preferences',
		target: function (@unused) { IDESettingsControl() })

	(UI, action, menu: IDE, name: ('Plugins'), target: IdeActions)
	(UI, pluginsubmenu, menu: IDE, name: "&Plugin Errors",
		target: function (@unused)
			{ Plugins().ShowErrors() })
	(UI, pluginsubmenu, menu: IDE, name: "&Reset Plugins",
		target: function (@unused)
			{ Plugins().Reset() })
	(UI, pluginsubmenu, menu: IDE, name: "PluginsView",
		target: function (@unused)
			{ PersistentWindow(PluginsView) })
	(UI, pluginsubmenu, menu: IDE, name: 'Create Plugin Wizard',
		target: function (@unused)
			{ PluginWizardControl() })

	(UI, attach, menu: IDE, to: WorkSpace)
	(UI, attach, menu: IDE, to: LibraryView)
	(UI, attach, menu: IDE, to: QueryView)
	(UI, attach, menu: IDE, to: SchemaView)
	(UI, attach, menu: IDE, to: ClassView)
	(UI, attach, menu: IDE, to: TestRunner)
	(UI, attach, menu: IDE, to: BookEdit)
	(UI, attach, menu: IDE, to: 'Version Control')

	(UI, action, menu: IDE, name: ('Switch Mode'), target: class
		{
		modes: (
			'Standalone': function() { IDESwitchMode(standalone:) }
			'Client Server with Random Port': function() { IDESwitchMode() }
			'Client Server with Default Port': function() { IDESwitchMode(defaultPort:) }
			)
		Menu_Switch_Mode()
			{
			return .modes.Members().Sort!()
			}
		On_Switch_Mode(option)
			{
			(.modes[option])()
			}
		})
	)
)
