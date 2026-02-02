// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Customize'

	CallClass(hwnd, key, query, sfOb, browse? = false,
		customizable = #(), browse_custom_fields = false, tabs = false, sub_title = '',
		disable_default = false, allowCustomTabs = false, virtual_list? = false)
		{
		return ToolDialog(hwnd,	Object(this, key, query, sfOb, browse?,
			customizable, browse_custom_fields, tabs, sub_title, disable_default,
			allowCustomTabs, virtual_list?), closeButton?: false)
		}
	New(key, query, sfOb, .browse? = false,
		customizable = #(), browse_custom_fields = false, tabs = false, sub_title = '',
		disable_default = false, allowCustomTabs = false, .virtual_list? = false)
		{
		super(.layout(key, query, sfOb, browse?,
			customizable, browse_custom_fields, tabs, sub_title, disable_default,
			allowCustomTabs))
		BookLog('Customize Dialog - Start')
		.tabs_ctrl = .FindControl('Tabs')
		}
	layout(key, query, sfOb, browse?, customizable, browse_custom_fields,
		tabs, sub_title, disable_default, allowCustomTabs)
		{
		layout = Object('Tabs')
		readonly = AccessPermissions(Customizable.PermissionOption()) isnt true

		if .customizeScreen?(browse_custom_fields, allowCustomTabs, tabs)
			layout.Add(Object('CustomizeScreen', customizable,
				browse_custom_fields, tabs, readonly,
				Tab: browse? ? 'Create/Edit Custom Column' : 'Customize Screen',
				:allowCustomTabs, :sfOb, custFieldName: key))

		if key isnt false
			layout.Add(Object('CustomizeFields',
				key, query,
				:sfOb, :browse?, :readonly, :disable_default,
				virtual_list?: .virtual_list?
				Tab: 'Customize Fields'))

		return Object('Vert',
			Object('TitleNotes', .build_title(key, sub_title), name: 'title')
			layout)
		}

	customizeScreen?(browse_custom_fields, allowCustomTabs, tabs)
		{
		if browse_custom_fields isnt false or allowCustomTabs
			return true
		return Object?(tabs)
			? tabs.GetDefault('custom_tabs', #()).NotEmpty?()
			: false
		}

	build_title(key, sub_title)
		{
		// coming from access and linked browse
		if key isnt false and key.Has?(' ~ ')
			{
			title = key.BeforeFirst(' ~ ')
			if key.AfterFirst(' ~ ').Has?(' | ')
				title = title.BeforeFirst(' | ')
			return .Title $ ' - ' $ title.Trim() $ Opt(TabsControl.TabSplit, sub_title)
			}
		return .Title
		}

	HelpButton_HelpPage()
		{
		return "/General/Reference/Customization Options"
		}

	switching_from: false
	TabControl_SelChanging()
		{
		.switching_from = .tabs_ctrl.GetSelected()
		return false
		}
	AllowSelectTab(i)
		{
		if .switching_from isnt false and i isnt .switching_from
			{
			result = .switching_from_tab().Save()
			if result is 'invalid'
				{
				.tabs_ctrl.Select(.switching_from)
				return false
				}
			}
		return true
		}
	switching_from_tab()
		{
		return .tabs_ctrl.GetControl(.switching_from)
		}

	GetList()
		{
		if .tabs_ctrl.Constructed?(1)
			return .tabs_ctrl.GetControl(1).GetList()
		return false
		}

	OnOK()
		{
		if false is .checkMandatoryFields()
			return

		retOb = Object()
		ctrlNames = #('CustomizeScreen', 'CustomizeFields')
		for (i = 0; i < 2; i++)
			{
			ret = false
			if false isnt ctrl = .FindControl(ctrlNames[i])
				ret = ctrl.Save()
			if ret is 'invalid'
				return
			retOb[i] = ret
			}
		BookLog('Customize Dialog - End (Ok)')
		.Window.Result(Object(screen: retOb[0], fields: retOb[1]))
		}
	checkMandatoryFields()
		{
		if false is customizeFields = .FindControl('CustomizeFields')
			return true

		if false is customizeScreen = .FindControl('CustomizeScreen')
			return true

		// using Select because constructing CustomizeFields tab does not list data
		cur = .tabs_ctrl.GetSelected()
		.tabs_ctrl.Select(cur is 0 ? 1 : 0)
		.tabs_ctrl.Select(cur)

		if .browse? or .virtual_list?
			return true

		mandatoryFields = customizeFields.GetMandatoryCustomFields()
		fields = customizeScreen.GetFieldsOnLayout()
		missingMandatory = mandatoryFields.Difference(fields)
		if not missingMandatory.Empty?()
			{
			.AlertInfo('Customize',
				'The following mandatory fields must be added to Customize Screen:\n' $
				missingMandatory.Map(Prompt).Join(', '))
			return false
			}
		return true
		}
	OnCancel()
		{
		if false is .checkMandatoryFields()
			return

		ret = false
		if false isnt ctrl = .FindControl('CustomizeFields')
			ret = ctrl.Save()
		if ret is 'invalid'
			return

		screen = false isnt (customizeScreen = .FindControl('CustomizeScreen')) and
			customizeScreen.FieldsChanged?()
		BookLog('Customizing Fields - End (Cancel)')
		.Window.Result(Object(:screen, fields: ret))
		}
	On_Cancel()
		{
		// disable Esc key
		}
	FieldRenamed(field)
		{
		if false isnt customizeFields = .FindControl('CustomizeFields')
			customizeFields.FieldRenamed(field)
		}
	// prevent close without using either ok or cancel button
	ConfirmDestroy()
		{
		return false
		}
	}
