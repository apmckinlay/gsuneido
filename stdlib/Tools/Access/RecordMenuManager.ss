// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(protectField, option, .historyFields, .ctrl = false, warningField = false)
		{
		.current_plugins = Object()
		.buildHistoryMenu()
		.buildCurrentMenu(option, protectField, ctrl, warningField)
		.global_plugins = Object()
		.Global = Object('Reporter...', 'Summarize...', 'CrossTable...', 'Export...')
		.setMenu(.Global, 'Global', .global_plugins, option)
		}

	buildHistoryMenu()
		{
		.historyFields = .historyFields is false ? Object() : .historyFields.Copy()
		if .ctrl isnt false
			.ctrl.Addons.Collect('HistoryMenu').Each()
				{
				for m in it.Members()
					.historyFields[m] = it[m]
				}
		}

	buildCurrentMenu(option, protectField, ctrl, warningField)
		{
		.Current = Object('Save', 'Print...')
		.setMenu(.Current, 'Current', .current_plugins, option)

		for menus in ctrl.Addons.Collect('CurrentMenu')
			for menu in menus
				.Current.AddUnique(menu)
		if protectField isnt false
			.Current.Add('Reason Protected')
		if warningField isnt false
			.Current.Add('', 'View Warnings')
		//added at the bottom to mitigate accidental deletions
		.Current.Add('', 'Restore', '', 'Delete', #('Delete'))
		}

	setMenu(menu, menuOption, plugins, option)
		{
		Plugins().ForeachContribution('AccessMenus', menuOption)
			{|c|
			if c.Member?('option') and c.option isnt option
				continue
			if .noHistoryMenu?(c)
				continue
			if .hideDevMenu?(c) or .hideMenu?(c)
				continue
			menu.AddUnique(c[2])
			plugins[c[2]] = c[2] is 'History'
				? Object(view: c.view, update: c.update)
				: c[3]  /*= plugin function */
			}
		}

	noHistoryMenu?(c)
		{
		return c[2] is 'History' and (.historyFields is false or
			Object?(.historyFields) and .historyFields.Empty?())
		}

	hideDevMenu?(c)
		{
		return c.GetDefault('devel', false) is true and	Suneido.User isnt 'default'
		}

	hideMenu?(c)
		{
		return c.GetDefault('showhide?', function() {return true})(ctrl: .ctrl) isnt true
		}

	UpdateHistory(data, newrecord?)
		{
		if .current_plugins.Member?('History') and not newrecord?
			(.current_plugins.History.update)(data, .historyFields)
		}

	On_Current(option, data, ctrl)
		{
		hwnd = ctrl.Window.Hwnd
		if .current_plugins.Member?(option)
			if option is 'History'
				(.current_plugins[option].view)(data, hwnd, .historyFields)
			else
				(.current_plugins[option])(:data, :hwnd, access: ctrl)
		else
			{
			if not option.Prefix?('On_Current')
				option = ToIdentifier('On_Current_' $ option)
			ctrl.Addons.Send(option)
			}
		}

	On_Global(option, hwnd)
		{
		if .global_plugins.Member?(option)
			(.global_plugins[option])(:hwnd)
		option = option.Tr('.')
		if .Member?('On_Global_' $ option) and .ctrl.Save()
			this['On_Global_' $ option]()
		}

	On_Global_Summarize()
		{
		.Summarize(.ctrl)
		}

	Summarize(ctrl)
		{
		SummarizeControl(ctrl)
		}

	On_Global_Reporter()
		{
		Reporter()
		}

	On_Global_CrossTable()
		{
		.CrossTable(.ctrl)
		}

	CrossTable(ctrl)
		{
		ToolDialog(ctrl.Window.Hwnd, Crosstab(
			ctrl.GetQuery(), ctrl.GetTitle(), ctrl.GetExcludeSelectFields(),
			columns: ctrl.GetFields()))
		}

	On_Global_Export()
		{
		.Export(.ctrl)
		}

	Export(ctrl)
		{
		GlobalExportControl(
			ctrl.GetQuery(), excludeSelectFields: ctrl.GetExcludeSelectFields()
			columns: ctrl.GetFields())
		}
	}
