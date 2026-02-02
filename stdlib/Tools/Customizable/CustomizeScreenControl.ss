// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// TODO: list context menu: modify, delete
Controller
	{
	Name: "CustomizeScreen"

	tabs_ctrl: false
	New(.c, browse_custom_fields = false, .tabs = false, readonly = false,
		.allowCustomTabs = false, .sfOb = false, .custFieldName = false)
		{
		super(.layout(browse_custom_fields, readonly))
		.selectFields = SelectFields(sfOb.cols, sfOb.excludeFields, false)
		.list = .FindControl('ListBox')
		.tabs_ctrl = .FindControl('Tabs')
		.init_fields()
		if .tabs isnt false // only needed if there are tabs
			.init_rec()
		}

	layout(browse_custom_fields, readonly)
		{
		.readonly = readonly
		.browse_custom_fields = browse_custom_fields

		listBoxLayout = .buildLayout(.browse_custom_fields)

		buttons_ob = Object('Horz'
			#(Button 'Create a Field...'),
			'Skip',
			#(Button 'Field Properties'))
		if .browse_custom_fields is false
			buttons_ob.Add('Skip', #(Button 'Add Field to Layout'))
		buttons_ob.Add('Fill', 'Skip', 'OkCancel')

		if .browse_custom_fields isnt false
			return listBoxLayout.Add(Object('Border', buttons_ob))

		tabSpec = Object('Tabs')
		for tab in .tabs.custom_tabs
			{
			if .c.LayoutHidden?(tab)
				continue
			tabSpec.Add(.editor(tab))
			}

		return Object('Vert',
			Object('Record'
				Object('HorzSplit',
					listBoxLayout,
					Object('Vert',
						.layoutHeader(),
						tabSpec,
						xstretch: 3), name: 'Vert2')
				)
			'Skip'
			buttons_ob
			)
		}

	layoutHeader()
		{
		layoutHeader = Object('Horz' #(Border (Static Layouts) 3) 'Fill')
		if .readonly isnt true and .allowCustomTabs is true
			{
			options = Object('Add Custom Tab', 'Add Custom Table Tab', '',
				'Rename Current Tab', '',
				'Remove Current Tab', 'Restore Tab')
			if .c.CustomTableKey() is false
				options.Remove('Add Custom Table Tab')
			layoutHeader.Add(Object('MenuButton', 'Custom Tabs', options))
			}
		return layoutHeader
		}

	buildLayout(browse_custom_fields = false)
		{
		listBoxLayout = Object('Vert')
		if browse_custom_fields is false
			listBoxLayout.Add(#(Border (Static 'Custom Fields') 3))
		listBoxLayout.Add(#(ListBox sort:))
		return listBoxLayout
		}
	editor(tab)
		{
		readonly = .readonly or .c.CustomTableTab?(tab)
		return Object('ScintillaAddonsEditor', Tab: tab, :readonly, Addon_detab:,
			Addon_highlight_customFields: Object(customizable: .c), font: StdFonts.Mono(),
			name: .editorName(tab))
		}
	editorName(tab)
		{
		return tab.Tr(' ')
		}
	resetList()
		{
		.list.DeleteAll()
		.init_fields()
		}

	init_fields()
		{
		.fieldList = .c.CustomFields()
		// need the Copy because AddItem modifies fieldList ???
		for m, v in .fieldList.Copy()
			.list.AddItem(Prompt(v), m)
		}
	origRec: false
	init_rec()
		{
		rec = Record()
		for(i = 0; i < .tabs_ctrl.GetAllTabCount(); i++)
			{
			tabName = .tabs_ctrl.TabName(i)
			rec[.editorName(tabName)] = .c.Layout(SelectFields(.c.CustomFields()),
				tabName)
			}
		.origRec = rec.Copy()
		.Vert.Data.Set(rec)
		}
	ListBox_ContextMenu(x, y)
		{
		if .readonly or .list.GetCurSel() is -1
			return

		ContextMenu(#('Delete')).ShowCall(this, x, y)
		}
	fieldsChanged?: false
	FieldsChanged?()
		{
		return .fieldsChanged?
		}
	On_Context_Delete()
		{
		if false is field = .getFieldForDelete()
			return

		prompt = .list.Get()
		warning = 'Are you sure? Deleting a field cannot be undone. \r\n' $
			'If you proceed you will no longer be able to access this field.\r\n'
		if false isnt tab = .check_field_used(field)
			{
			if tab.Prefix?('a tab that you do not have permission to')
				{
				.AlertInfo('Can Not Delete', 'This field exists on ' $ tab)
				return
				}
			warning = prompt  $ ' field is already selected on ' $
				tab $ '\r\n' $ warning
			}

		if ToolDialog(.Window.Hwnd, [.warningDialog, warning])
			{
			.list.DeleteItem(.list.GetCurSel())
			.c.DeleteField(field)
			.markFieldsChanged()
			if Object?(.browse_custom_fields)
				.browse_custom_fields.Add(field)
			}
		}

	getFieldForDelete()
		{
		if .readonly
			return false

		msg = ''
		if false is field = .fieldList.GetDefault(.list.GetSelected(), false)
			msg = 'Could not find field to delete'

		if msg isnt '' or '' isnt msg = .c.CustomField_AllowDelete(field)
			{
			.AlertInfo('Can Not Delete', 'In order to delete (' $ Prompt(field) $
				') please remove from the following places:\r\n' $ .formatMsg(msg))
			return false
			}
		return field
		}

	formatMsg(msg)
		{
		msgOb = msg.Lines()
		helpReq = Object()
		for m in msgOb
			if m =~ '^Axon:.*'
				{
				helpReq.Add(m.AfterFirst('Axon:'))
				msgOb.Remove(m)
				}
		return msgOb.Join('\r\n') $
			Opt('\r\nAxon Support Required for the following places:\r\n',
				helpReq.Join('\r\n'))
		}

	warningDialog: Controller
		{
		Title: "Delete custom field"
		New(prompt)
			{
			super(['Vert',
				Object('Static' prompt size: '+2', weight: 'bold', color: CLR.RED),
				'Skip',
				#(Horz #(Button "Yes, delete this field" command: 'yes')
					#Fill #(Button "No, do not delete" command: 'no')) ])
			}
		On_yes()
			{ .Window.Result(true) }
		On_no()
			{ .Window.Result(false) }
		}

	handleRename(field, newName, block)
		{
		sf = SelectFields(.c.CustomFields())
		tab = .check_field_used(field)
		.selectFields.AddField(field, newName)

		block()

		.renameInFormula(field, sf)

		if .c.LayoutExists?(tab)
			{
			rec = .getRecord()
			rec[.editorName(tab)] = .c.Layout(SelectFields(.c.CustomFields()), tab)
			}

		.resetList()
		.markFieldsChanged()
		.Send('FieldRenamed', field)
		}
	renameInFormula(field ,sf)
		{
		.c.EnsureTable()
		// if list is false than the other tab isn't constructed.
		if false isnt list = .Send('GetList')
			{
			for item in list
				{
				if item.custfield_formula_fields.Split(',').Has?(field)
					item.custfield_formula = .c.RebuildLayout(sf, item.custfield_formula)
				}
			}
		else // if tab not construted, need to update data diretly
			{
			QueryApplyMulti('customizable_fields
				where custfield_name is ' $ Display(.custFieldName), update:)
				{
				if it.custfield_formula_fields.Split(',').Has?(field)
					{
					it.custfield_formula = .c.RebuildLayout(sf, it.custfield_formula)
					it.Update()
					}
				}
			}
		}

	ListBoxDoubleClick(i)
		{
		.add_field(i)
		}

	On_Add_Field_to_Layout()
		{
		if false isnt csf = .getCurrentlySelectedField()
			.add_field(csf)
		}

	getCurrentlySelectedField()
		{
		if -1 is sel = .list.GetCurSel()
			{
			.AlertInfo(.AlertTitle, 'Please select a field')
			return false
			}
		return sel
		}

	validateTabName(name)
		{
		if name is ''
			return "Please Choose a name"
		return ""
		}

	On_Custom_Tabs_Add_Custom_Tab()
		{
		name = .validateNewCustomTab('Add Custom Tab', 'Custom Tab Name')
		if name isnt false
			.insertCustomTab(name)
		}

	validateNewCustomTab(title, prompt)
		{
		if false is .checkCurrentVisibleTabCount(title)
			return false

		if false is name = Ask(prompt, title, valid: .validateTabName)
			return false

		if false is .validateCustomTabName(name, .tabs, .tabs_ctrl, .c, allowHidden:)
			return false

		if false isnt cl = OptContribution('CustomTabPermissions', false)
			{
			option = .c.GetTabPermissionName(name)
			cl.AddPermission(option)
			if QueryEmpty?('biz_permissions', bizperm_option: option)
				{
				.AlertInfo(title, 'Unable to create permissions. Custom tab not added')
				return false
				}
			}
		return name
		}

	maxVisibleCustomTabs: 25
	checkCurrentVisibleTabCount(title)
		{
		if .c.CountVisibleCustomTabs() >= .maxVisibleCustomTabs
			{
			.AlertInfo(title, 'Tabs cannot be added or restored if there are already ' $
				'at least ' $ .maxVisibleCustomTabs $ ' visible custom tabs.')
			return false
			}
		return true
		}

	// allowHidden used so Add Custom Tab can restore hidden tabs
	// but still prevent renaming a tab to a hidden one.
	validateCustomTabName(name, tabs, tabs_ctrl, c, allowHidden = false)
		{
		if ((tabs_ctrl.GetAllTabNames().Has?(name) or tabs.all_tabs.Has?(name)) or
			(c.TabCustom?(name) and
				(allowHidden is false or not c.LayoutHidden?(name))))
			{
			.AlertInfo('Add Custom Tab',
				'Reserved name or Tab name already exists. Choose another name')
			return false
			}
		return true
		}

	insertCustomTab(name, customTableTab? = false)
		{
		.c.SaveTab(name, :customTableTab?)
		.tabs_ctrl.Append(name, .editor(name))
		rec = .getRecord()
		rec[.editorName(name)] = .c.Layout(SelectFields(.c.CustomFields()), name)
		.markFieldsChanged()
		}

	On_Custom_Tabs_Add_Custom_Table_Tab()
		{
		name = .validateNewCustomTab('Add Custom Table Tab', 'Custom Table Tab Name')
		if name isnt false
			.insertCustomTab(name, customTableTab?:)
		}

	On_Custom_Tabs_Rename_Current_Tab()
		{
		if -1 is selected = .tabs_ctrl.GetSelected()
			{
			.AlertInfo('Rename Tab', 'Please select a tab to rename.')
			return
			}
		oldName = .tabs_ctrl.TabName(selected)
		if not .c.TabCustom?(oldName)
			{
			.AlertInfo('Rename Tab', 'Only Custom tabs can be renamed.')
			return
			}
		if .checkCustomTableTabForReports(oldName, 'rename')
			return
		if false is newName = Ask("Custom Tab Name", "Add Custom Tab",
			ctrl: Object('Field' set: oldName), valid: .validateTabName)
			return

		if false is .validateCustomTabName(newName, .tabs, .tabs_ctrl, .c)
			return

		if false isnt cl = OptContribution('CustomTabPermissions', false)
			cl.RenamePermission(
				.c.GetTabPermissionName(oldName), .c.GetTabPermissionName(newName))
		.renameTab(oldName, newName, selected)
		.markFieldsChanged()
		}

	markFieldsChanged()
		{
		.fieldsChanged? = true
		}

	checkCustomTableTabForReports(tab, action)
		{
		reports = .c.CustomReporterReports(tab)
		if reports.Empty?()
			return false
		msgOb = ['Cannot ' $ action $ ' as it is used on the following Reporter Reports:']
		msgOb.Add(@reports.Map!({ it.report.AfterFirst('Reporter ') }))
		.AlertInfo(action.Capitalize() $ ' Tab', msgOb.Join('\r\n\t'))
		return true
		}

	renameTab(oldName, newName, selected)
		{
		.c.RenameLayout(oldName, newName)
		.tabs_ctrl.Remove(selected)
		.tabs_ctrl.Insert(newName, .editor(newName) at: selected)
		rec = .getRecord()
		neweditor = .editorName(newName)
		oldeditor = .editorName(oldName)
		rec[neweditor] = rec[oldeditor]
		rec.Delete(oldeditor)
		}
	On_Custom_Tabs_Remove_Current_Tab()
		{
		if -1 is selected = .tabs_ctrl.GetSelected()
			{
			.AlertInfo('Remove Tab', 'Please select a tab to remove.')
			return
			}
		name = .tabs_ctrl.TabName(selected)
		if not .c.TabCustom?(name)
			{
			.AlertInfo('Remove Tab', 'Only custom tabs can be removed.')
			return
			}
		if .checkCustomTableTabForReports(name, 'remove')
			return
		if false isnt cl = OptContribution('CustomTabPermissions', false)
			cl.RemovePermission(.c.GetTabPermissionName(name))

		.c.HideLayout(name)
		.tabs_ctrl.Remove(selected)
		if .tabs_ctrl.GetTabCount() isnt 0
			.tabs_ctrl.Select(selected is 0 ? 0 : selected - 1)
		.markFieldsChanged()
		}
	On_Custom_Tabs_Restore_Tab()
		{
		if false is .checkCurrentVisibleTabCount('Restore Tab')
			return

		list = .c.ListHiddenLayouts()
		if false is tabName =
			Ask("Tab to Restore", "Restore Tab", ctrl: Object('ChooseList', :list),
				valid: .validateTabName)
			return
		if false isnt cl = OptContribution('CustomTabPermissions', false)
			cl.AddPermission(.c.GetTabPermissionName(tabName))

		.c.ShowLayout(tabName)
		.tabs_ctrl.Insert(tabName, .editor(tabName))
		rec = .getRecord()
		rec[.editorName(tabName)] = .c.Layout(SelectFields(.c.CustomFields()), tabName)
		.markFieldsChanged()
		}

	AlertTitle: "Customize"
	On_Field_Properties()
		{
		if .readonly
			{
			.permissionMsg()
			return
			}
		i = .getCurrentlySelectedField()
		if false is i
			return
		field = .fieldList[.list.GetData(i)]
		textSelectedRow = .list.GetText(i)

		fo = .evaluateClass(field)
		fs = .evaluateSuperClass(field)
		fd = Object(
			custpe: textSelectedRow
			colnme: field
			flddef: fo
			fldbse: fs
			)
		_customizable = .c
		result = CustomizableFieldDialog(fd)

		if false isnt result
			{
			// Prompt Changed
			if result.colpro isnt fd.custpe
				{
				if "" isnt msg = .checkForDuplicate(.list, result.colpro)
					{
					.AlertError('Create a Field', msg)
					return
					}
				.handleRename(fd.colnme, result.colpro, { .c.UpdateField(result) })
				}
			else
				.c.UpdateField(result)
			.markFieldsChanged()
			}
		}
	permissionMsg()
		{
		.AlertInfo(.AlertTitle, "You must have permission to " $
			Customizable.PermissionOption()[1..].Replace('/', ' > ') $
			" to access this option")
		}
	evaluateClass(className)
		{
		return Global('Field_' $ className)
		}
	evaluateSuperClass(className)
		{
		originalSource = Query1('configlib', name: 'Field_' $ className).text
		return ClassHelp.SuperClass(originalSource)
		}

	add_field(i)
		{
		if .browse_custom_fields isnt false
			return

		if .readonly
			{
			.permissionMsg()
			return
			}

		if false is editor = .tabs_ctrl.GetControl()
			{
			.AlertInfo(.AlertTitle, 'Please select a tab to add the Custom Field to.')
			return
			}

		field = .fieldList[.list.GetData(i)]
		prompt = .list.GetText(i)
		if false isnt tab = .check_field_used(field)
			.AlertError(.AlertTitle, 'This field is already selected on ' $ tab)
		else
			editor.Paste(prompt $ ' ')
		}
	getRecord() // for tests
		{
		return .Vert.Data.Get()
		}
	check_field_used(field)
		{
		if .tabs_ctrl is false
			return false

		if .nonPermissableField?(field)
			return "a tab that you do not have permission to, " $
				"please contact an administrator."

		sf = SelectFields(.c.CustomFields())

		rec = .getRecord()
		for(i = 0; i < .tabs_ctrl.GetAllTabCount(); i++)
			{
			editor = .editorName(.tabs_ctrl.TabName(i))
			fields = sf.FormulaFields(rec[editor].Trim())
			if fields.Has?(field)
				return .tabs_ctrl.TabName(i)
			}
		return false
		}

	nonPermissableField?(field)
		{
		tableName = .c.GetName()
		if false isnt cl = OptContribution('CustomTabPermissions', false)
			return cl.NonPermissableFieldsFromTable(tableName).Has?(field)

		return false
		}

	On_Create_a_Field()
		{
		if .readonly
			{
			.permissionMsg()
			return
			}
		_customizable = .c
		result = CustomizableFieldDialog(text: 'Create Field')
		if result is false
			return

		if "" isnt msg = .checkForDuplicate(.list, result.colpro)
			{
			.AlertError('Create a Field', msg)
			return
			}
		if "" isnt msg = .c.ExtraCreateFieldChecking(result.colpro, result.ctllbl)
			{
			.AlertError('Create a Field', msg)
			return
			}
		field = .c.CreateField(result.colpro, result.ctllbl, .selectFields,
			result.options)
		.markFieldsChanged()
		.resetList()
		if Object?(.browse_custom_fields)
			.browse_custom_fields.Add(field)
		}

	checkForDuplicate(list, fieldName)
		{
		duplicateFound? = false
		if -1 isnt index = list.FindString(fieldName)
			while(duplicateFound? is false and index < list.GetCount())
				{
				duplicateFound? = fieldName =~  '(?q)' $ list.GetText(index)
				index++
				}
		return duplicateFound?
			? 'Field name ' $ Display(fieldName) $ ' already exists.\nDuplicate ' $
				'field names are not allowed - fields must have different names.'
			: ""
		}

	check_restrictions()
		{
		sf = SelectFields(.c.CustomFields())
		if .tabs_ctrl is false
			return false

		return .checkRestrictionsPerTab(.tabs_ctrl, sf, .logError)
		}
	maxLinesOnHeader: 6
	maxFieldsOnRow: 8
	maxTotalFields: 500
	checkRestrictionsPerTab(tabsCtrl, sf, logError)
		{
		rec = .getRecord()
		fields = Object()
		for(i = 0; i < tabsCtrl.GetAllTabCount(); i++)
			{
			tab_name = tabsCtrl.TabName(i)
			layout = rec[.editorName(tab_name)].Trim()
			if tab_name.Has?('Header') and layout.Lines().Size() > .maxLinesOnHeader
				{
				logError('You cannot add more than ' $ .maxLinesOnHeader $
					' lines of customized fields to the Header screen.')
				return false
				}

			for line in layout.Lines()
				if sf.FormulaFields(line).Size() > .maxFieldsOnRow
					{
					logError('You cannot add more than ' $ .maxFieldsOnRow $
						' customized fields to a row (' $ tab_name $ ' tab.)')
					return false
					}

			for field in sf.FormulaFields(layout)
				if not fields.Has?(field)
					fields.Add(field)
				else
					{
					logError('You cannot add the same field (' $
						Prompt(field) $	') to a screen more than once.')
					return false
					}
			}
		return .checkOverallFieldCount(fields, logError)
		}

	checkOverallFieldCount(fields, logError)
		{
		if fields.Size() > .maxTotalFields
			{
			logError('The total number of fields on all tab layouts cannot exceed ' $
				.maxTotalFields $ ' (Currently there are ' $ fields.Size() $ ')')
			return false
			}
		return true
		}
	// broken out for tests
	logError(msg)
		{
		.AlertError(.AlertTitle, msg)
		}
	GetFieldsOnLayout()
		{
		sf = SelectFields(.c.CustomFields())
		fields = Object()
		rec = .getRecord()
		for(i = 0; i < .tabs_ctrl.GetAllTabCount(); i++)
			{
			layout = rec[.editorName(.tabs_ctrl.TabName(i))]
			fields.MergeUnion(sf.FormulaFields(layout))
			}
		return fields
		}
	On_Cancel()
		{
		.Send('OnCancel')
		}
	On_OK()
		{
		.Send('OnOK')
		}
	Save()
		{
		if .readonly
			return false

		if Object?(.browse_custom_fields)
			return true
		else if .tabs_ctrl isnt false
			{
			if false is .check_restrictions()
				return 'invalid'

			if .origRec is false
				return false

			sf = SelectFields(.sfOb.cols, .sfOb.excludeFields, false)
			rec = .getRecord()
			dirty? = false
			for(i = 0; i < .tabs_ctrl.GetAllTabCount(); i++)
				{
				name = .tabs_ctrl.TabName(i)
				editor = .editorName(name)
				// new tab, retored hidden tab, or tab change
				if not .origRec.Member?(editor) or .origRec[editor] isnt rec[editor]
					{
					.c.SaveLayout(rec[editor].Trim(), sf, name, onlyCustomFields?:)
					dirty? = true
					}
				}
			return dirty? or .fieldsChanged?
			}
		else
			return false
		}

	TabContextMenu(x, y, hover = false, source = false)
		{
		idx = hover
		if .readonly is true
			return 0

		if idx is false
			{
			ScreenToClient(source.Hwnd, pt = Object(:x, :y))
			idx = SendMessageTabHitTest(source.Hwnd, TCM.HITTEST, 0, tch = Object(:pt))

			if tch.flags is TCHT.NOWHERE
				return 0
			}

		tab = .tabs_ctrl.TabName(idx)
		if not .c.TabCustom?(tab)
			return 0
		tabmenu = Object('Rename', 'Remove')

		if 0 is i = ContextMenu(tabmenu).Show(source.Hwnd, x, y)
			return 0
		.tabs_ctrl.Select(idx)
		switch (tabmenu[i-1])
			{
		case 'Rename':
			.On_Custom_Tabs_Rename_Current_Tab()
		case 'Remove':
			.On_Custom_Tabs_Remove_Current_Tab()
			}
		return 0
		}
	}
