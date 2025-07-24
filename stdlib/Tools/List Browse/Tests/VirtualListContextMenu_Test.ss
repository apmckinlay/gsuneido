// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.old_user = Suneido.User
		Suneido.User = 'test'

		.fieldName = .TempName()
		.MakeLibraryRecord([name: .fieldName $ 'Format', text: `class {
			List_ExtraContext()
				{
				return 'Special Format ContextMenu'
				}
			}`],
			[name: "Field_" $ .fieldName, text: `Field_num {
				Prompt: "Special Field"
				Format: (` $ .fieldName $ `) }`])
		}

	Test_buildMenu()
		{
		menu = VirtualListContextMenu(#())
		control = Mock()
		control.Addons = class
			{
			Collect(unused) { return  #() }
			}
		.SpyOn(RecordMenuManager.RecordMenuManager_setMenu).Return(false)
		control.When.GetHeaderSelectPrompt().Return('prompts')
		recMenu = RecordMenuManager(false, false, false, control)

		control.When.Editable?().Return(false)
		Assert(menu.VirtualListContextMenu_buildMenu(control, Object("Test"))
			is: Object("Reset Columns", "Test", "", "Print...", "Reporter..."))

		menu.ContextRec = []
		control.When.Editable?().Return(true)
		control.When.GetSelectedRecords().Return(#())
		Assert(menu.VirtualListContextMenu_buildMenu(control)
			is: #(New))

		control.When.GetSelectedRecords().Return(#(1, 2))
		Assert(menu.VirtualListContextMenu_buildMenu(control)
			is: #(New))

		menu2 = VirtualListContextMenu(#(), :recMenu, addCurrentMenu?:,
			addGlobalMenu?:)
		menu2.ContextRec = []
		menu2.ContextCol = .fieldName

		control.When.GetSelectedRecords().Return(#(1))
		Assert(menu.VirtualListContextMenu_buildMenu(control)
			is:  #("Edit Field", "New"))
		// default record menu + current menu + global menu + format menu
		Assert(menu2.VirtualListContextMenu_buildMenu(control)
			is: Object("Edit Field", "New",
				"Save", "Print...", "", "Restore", "", "Delete", #("Delete"), "", "Global"
				Object('Reporter...', 'Summarize...', 'CrossTable...', 'Export...'),
				"", "Special Format ContextMenu"))

		menu.ContextRec = false
		menu2.ContextRec = false
		menu2.ContextCol = ''
		control.When.GetSelectedRecords().Return(#())
		Assert(menu.VirtualListContextMenu_buildMenu(control)
			is: Object("Reset Columns", "", "Print...", "Reporter...", "New"))
		// default header menu + only global menu
		Assert(menu2.VirtualListContextMenu_buildMenu(control)
			is: Object("Reset Columns", "", "Print...", "Reporter...", "New", "", "Global"
				Object('Reporter...', 'Summarize...', 'CrossTable...', 'Export...')))
		}

	Teardown()
		{
		super.Teardown()
		Suneido.User = .old_user
		}
	}