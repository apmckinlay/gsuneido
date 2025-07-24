// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'Select2'
	Numrows: 15
	New(sf, select_vals, printParams = false, option = '', title = '',
		.initialSelect = [], .menuOptions = false)
		{
		super(.layout(sf, printParams, option, title, menuOptions))
		.data = .Data

		// fill in values from previous select
		selected = select_vals.Copy()
		for m in selected.Members()
			if m.Prefix?('oplist')
				selected[m] = TranslateLanguage(selected[m])
		.data.Set(selected)
		.ignore_edit_change = false
		}
	layout(.sf, printParams, option, title, menuOptions)
		{
		.ignore_edit_change = true
		.ops = Select2.Ops
		.ops_desc = Select2.TranslateOps()
		return .layout2(printParams, option, title, menuOptions)
		}
	layout2(printParams, option, title, menuOptions)
		{
		ctrls = Object()
		headings = Object(#(Static '')
			#(Horz (Skip 4) (Static 'Field'))
			#(Horz (Skip 4) (Static 'Operator'))
			#(Horz (Skip 4) (Static 'Value')))
		if printParams is true
			headings.Add(#(Horz (Static 'Print?')))
		if menuOptions is true
			headings.Add(#(Horz (Static 'Menu Option')))

		ctrls.Add(headings)
		fieldlist = .sf.Prompts().Sort!()
		for (i = 0; i < .Numrows; ++i)
			{
			checkbox = Object('CheckBox', name: "checkbox" $ i)
			fields = Object('ChooseList', fieldlist, width: 20, name: "fieldlist" $ i)
			ops = Object('ChooseList', .ops_desc, name: "oplist" $ i)
			vals = Object('Field', name: "val" $ i)
			row = Object(checkbox, fields, ops, vals)
			if printParams is true
				row.Add(Object('CheckBox', name: "print" $ i))
			if menuOptions is true
				row.Add(Object('Horz' 'Fill'
					Object('CheckBox', name: "menu_option" $ i, align: 'auto') 'Fill'))
			ctrls.Add(row)
			}
		return Object('Record',
			Object('Vert',
				Object('Horz',
					#(Static "Checkmark the rows you want to apply"),
					'Fill',
					option is '' ? 'Skip'
						: Object('Presets', option, title, initial: #(Reset))
					)
				'Skip'
				Object("Grid", ctrls)))
		}

	LoadPresets(ignore?)
		{
		.ignore_edit_change = ignore?
		}

	On_Reset() // from Presets
		{
		sel = Record()
		AccessSelectMgr.Add_initial_selects(.initialSelect, sel)
		.data.Set(sel)
		}

	Record_NewValue(member, value)
		{
		if member.Has?('checkbox')
			return
		i = Number(member.Extract('[0-9]+$'))
		data = .data.Get()
		if data['fieldlist' $ i] isnt '' and data['oplist' $ i] isnt ''
			data['checkbox' $ i] = true
		.protect_menuOptions(member, value)
		}
	protect_menuOptions(member, value)
		{
		if .menuOptions is false
			return

		i = member.Tr('^0-9')
		// if the menu option checkbox is checked, the parameter should always be
		// included in the select
		if member.Prefix?('menu_option') and
			false isnt check_ctrl = .FindControl('checkbox' $ i)
			{
			if value is true
				.Data.SetField('checkbox' $ i, true)
			check_ctrl.SetReadOnly(value)
			}
		}
	Edit_Change(source)
		{
		if .ignore_edit_change or false is source.Name.Has?("val")
			return
		i = Number(source.Name.Extract('[0-9]+$'))

		// have to use RecordControl's GetField / SetField because
		// RecordControl.Get changes the focus (causes problems with the FieldControl)
		if .data.GetField('fieldlist' $ i) isnt '' and
			.data.GetField('oplist' $ i) isnt ''
			.data.SetField('checkbox' $ i, true)
		}
	Where(fields = false, except = false, extra_dd = #())
		{
		if .data.Valid() isnt true
			return false
		data = .data.Get()
		if true isnt checkMenu = .checkMenuOptions(data)
			{
			Alert(checkMenu, title: "Error", flags: MB.ICONERROR)
			return false
			}
		result = Select2(.sf).Where(data, fields, except, extra_dd)
		if result.errs isnt ""
			{
			Alert(result.errs, title: "Error", flags: MB.ICONERROR)
			return false
			}
		return result
		}

	checkMenuOptions(data)
		{
		menuOptions = Object()
		multiMenuNames = Object()
		for (i = 0; i < .Numrows; ++i)
			{
			selectField = data["fieldlist" $ i]
			if selectField isnt "" and data["menu_option" $ i] is true
				{
				if menuOptions.Has?(selectField)
					multiMenuNames.Add(selectField)
				menuOptions.Add(selectField)
				}
			}
		if not multiMenuNames.Empty?()
			return "Can't have the same field as a Menu Option more than once: " $
				multiMenuNames.Join(',')
		return true
		}

	Get()
		{
		data = .data.Get().Copy()
		for op in data.Members()
			if (op.Prefix?('oplist') and (false isnt pos = .ops_desc.Find(data[op])))
				data[op] = Select2.Ops[pos][0]
		return data
		}
	Set(data)
		{
		data_transl = data.Copy()
		for op in data_transl.Members()
			if op.Prefix?('oplist')
				data_transl[op] = TranslateLanguage(data_transl[op])
		.ignore_edit_change = true
		.data.Set(data_transl)
		.ignore_edit_change = false
		// need to do this after RecordControl.Set so that set readonly
		// does not get cleared out
		for op in data_transl.Members().Copy()
			.protect_menuOptions(op, data_transl[op])
		}
	ChangeSF(.sf)
		{
		fieldlist = .sf.Prompts().Sort!()
		for (i = 0; i < .Numrows; ++i)
			{
			ctrl = .data.GetControl('fieldlist' $ i)
			ctrl.SetList(fieldlist)
			ctrl.FieldReturn()
			}
		.data.Dirty?(true)
		}
	}
