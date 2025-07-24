// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		plugins = new .mockPlugins
		Assert(plugins.Contributions('One')
			is: #((One, menu, from: 'lib0:Plugin_Two')))
		Assert({ plugins.Contributions('xyxyxyx') }
			throws: 'nonexistent plugin')
		Assert({ plugins.Contributions('Two') }
			throws: 'nonexistent plugin')
		Assert({ plugins.Contributions('One', 'menuu') }
			throws: 'nonexistent extension point')
		}
	mockPlugins: Plugins
		{
		Plugins_foreachPluginLibraryRecord(block)
			{
			block(#(name: Plugin_One, lib: lib0,
				text: "#(
					ExtensionPoints:
						(
						(menu)
						)
					)"))
			block(#(name: Plugin_Two, lib: lib0,
				text: "#(
					Contributions:
						(
						(One, menu)
						)
					)"))
			}
		Plugins_log_error(unused)
			{
			}
		}
	Test_multiple_extenpts()
		{
		plugins = new .mockPlugins2
		Assert(plugins.Errors has: 'multiple definitions')
		}
	mockPlugins2: Plugins
		{
		Plugins_foreachPluginLibraryRecord(block)
			{
			block(#(name: Plugin_One, lib: lib1,
				text: "#(
					ExtensionPoints:
						(
						(menu)
						)
					)"))
			block(#(name: Plugin_One, lib: lib2,
				text: "#(
					ExtensionPoints:
						(
						(menu)
						)
					)"))
			}
		Errors: ''
		Plugins_log_error(err)
			{
			.Errors $= err
			}
		}

	Test_foreachPluginLibraryRecord()
		{
		lib1 = .MakeLibrary(
			[name: 'Plugin_RealValidPluginLib1', group: -1],
			[name: 'Plugin_InvalidDeleted', group: -2],
			[name: 'Plugin_InvalidFolder', group: 1],
			[name: 'Plugins_InvalidName', group: 1],
			[name: SvcTable.MaxCommitName, group: -3]
			)
		lib2 = .MakeLibrary(
			[name: 'Plugin_RealValid1PluginLib2', group: -1],
			[name: 'Plugin_RealValid2PluginLib2', group: -1],
			[name: 'Plugin_InvalidDeleted', group: -2],
			[name: 'Plugin_InvalidFolder', group: 1],
			[name: SvcTable.MaxCommitName, group: -3]
			)
		lib3 = .MakeLibrary() // empty library should not cause issues

		mock = Mock(Plugins)
		mock.When.libraries().Return([lib1, lib2, lib3])
		mock.When.foreachPluginLibraryRecord([anyArgs:]).CallThrough()

		plugins = Object()
		mock.foreachPluginLibraryRecord({ plugins.Add(it.name) })
		Assert(plugins isSize: 3)
		Assert(plugins has: 'Plugin_RealValidPluginLib1')
		Assert(plugins has: 'Plugin_RealValid1PluginLib2')
		Assert(plugins has: 'Plugin_RealValid2PluginLib2')
		}
	}