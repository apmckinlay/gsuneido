// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Menu_Edit_a_Book()
		{
		return BookTables().Add('New Book...')
		}
	On_Edit_a_Book(option)
		{
		if option is 'New Book...'
			{
			if false is newbook = Ask("New Book Name", "Create New Book")
				return
			if false is BookEditModel.Create(newbook)
				{
				Alert("Can't create book: " $ newbook, "Book Creation Error",
					flags: MB.ICONERROR)
				return
				}
			option = newbook
			}
		PersistentWindow(Object(BookEditControl, option))
		}

	Menu_Open_a_Book()
		{
		return BookTables()
		}
	On_Open_a_Book(option)
		{
		BookControl(option, help_book: option =~ 'Help|doc')
		}

	On_MultiView_a_Query()
		{
		if false isnt query = .choose_table("MultiView a Query")
			MultiViewControl(
				args: Object(
					query,
					title: query,
					option: query
					)
				accessArgs: Object()
				listArgs: Object(
					headerSelectPrompt: 'no_prompts'
					defaultColumns: #()
					)
				)
		}

	choose_table(title)
		{
		query = Ask("Query", title, ctrl: [#AutoChoose,
			QueryList('tables',	#table).SortWith!(#Lower), xmin: 225, allowOther:])
		return query is false or query is "" ? false : query
		}

	Menu_Plugins()
		{
		return .build_plugin_menu('pluginsubmenu')
		}
	On_Plugins(option)
		{
		.run_plugin_target('pluginsubmenu', option)
		}
	Menu_Views()
		{
		return .build_plugin_menu('viewsubmenu')
		}
	On_Views(option)
		{
		.run_plugin_target('viewsubmenu', option)
		}
	build_plugin_menu(plugin)
		{
		ob = Object()
		Plugins().ForeachContribution("UI", plugin)
			{ |c|
			if c.menu is 'IDE'
				ob.Add(c.name)
			}
		return ob
		}
	run_plugin_target(plugin, option)
		{
		target = ''
		Plugins().ForeachContribution("UI", plugin)
			{ |c|
			if c.name is option
				target = c.target
			}
		(target)()
		}
	}