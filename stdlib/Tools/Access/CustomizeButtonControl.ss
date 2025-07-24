// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
/* README
This control is constructed outside of the normal scope of controls.
As a result, the majority of values need to be set/reset on the fly.
Otherwise, the control risks referencing outdated data/controls.

This is why the majority of values are set during: On_Customize
*/
Controller
	{
	Controls: #('EnhancedButton', command: 'Customize', image: 'custom_screen.emf',
		mouseEffect:, imagePadding: .1, tip: 'Customize')
	New(.disable_default = false)
		{ }

	Startup()
		{
		.image = .EnhancedButton
		.updateImage()
		}

	updateImage()
		{
		if .customized?()
			.image.SetImageColor(CLR.Active, CLR.Active)
		else
			.image.SetImageColor(CLR.Inactive, CLR.Inactive)
		}

	customized?()
		{
		customKey = .Send('GetAccessCustomKey')
		return Customizable.IsCustomized?(.Send('GetQuery')) or
			CustomizeField.IsCustomized?(customKey)
		}

	On_Customize()
		{
		if false is .Send('Access_Save')
			return

		query = .Send('GetQuery')
		Assert(query isnt 0)
		table = QueryGetTable(query)
		.collectTabs(.Window, tabs = Object())
		customizable = .getCustomizable(table)

		key = .Send('GetAccessCustomKey')
		sfOb = Object(
			cols: .getCustomizeFieldsCols(tabs, query, key, customizable),
			excludeFields: .excludeFields())
		dirty = CustomizeDialog(.Window.Hwnd, key, query, sfOb,
			:customizable,
			tabs: .customDlgTabs(tabs),
			disable_default: .disable_default,
			allowCustomTabs: .Send('AllowCustomTabs?') is true)

		.refresh(dirty)
		}

	collectTabs(ctrl, tabs)
		{
		ctrl.GetChildren().Each()
			{
			if it.Base?(TabsControl)
				tabs.Add(it)
			.collectTabs(it, tabs)
			}
		}

	getCustomizeFieldsCols(tabs, query, key, customizable)
		{
		tabFields = Object()
		tabs.Each({ tabFields.MergeUnion(it.CollectFields(customizable)) })
		cols = QueryColumns(query).
			Filter({ .findHeaderControl(it) isnt false or tabFields.Has?(it) })

		if key isnt false
			cols.MergeUnion(
				QueryList('customizable_fields where custfield_name is ' $ Display(key),
					'custfield_field'))

		if Object?(extra = .Send('CustomizeFields_ExtraFields'))
			cols.MergeUnion(extra)
		return cols
		}

	findHeaderControl(ctrl)
		{
		return .Window.FindControl(ctrl, exclude: TabsControl)
		}

	excludeFields()
		{
		excludeFields = .Send('GetExcludeSelectFields')
		return excludeFields isnt 0
			? excludeFields.Filter({ Datadict(it).Control[0] not in ('Key', 'Id') })
			: Object()
		}

	getCustomizable(table)
		{
		if 0 is name = .Send('GetCustomizableName')
			name = false
		if 0 is customKey = .Send('GetAccessCustomKey')
			customKey = ''
		return Customizable(table, name, :customKey)
		}

	customDlgTabs(tabs)
		{
		all_tabs = Object('Header')
		tabs.Each({ all_tabs.MergeUnion(it.GetAllTabNames()) })
		custom_tabs = .getCustomTabs(tabs)
		return [:all_tabs, :custom_tabs]
		}

	getCustomTabs(tabs)
		{
		custom_tabs = Object()
		tabs.Each({ custom_tabs.MergeUnion(it.CustomizableTabs()) })
		if false isnt ctrl = .findHeaderControl('Customizable')
			{
			customHeader = ctrl.TabName is false ? 'Header' : ctrl.TabName
			if customHeader isnt CustomizeExpandControl.LayoutName
				custom_tabs.Remove(customHeader).Add(customHeader, at: 0)
			}
		return custom_tabs.Empty?() ? false : custom_tabs
		}

	refresh(dirty)
		{
		if not Object?(dirty) or dirty.screen or dirty.fields
			.Send('BookRefresh')
		else
			.updateImage()
		}
	}
