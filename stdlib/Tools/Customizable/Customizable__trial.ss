// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
class
	{
	lib: 'configlib'
	savein: 'customizable'
	user: ''
	PermissionOption()
		{
		return OptContribution('CustomizablePermission', '')
		}
	New(table, name = false, .availableFields = false, .defaultLayout = '', .user = '',
		.customKey = '')
		{
		.table = table
		if name is false
			name = table
		.name = name

		.EnsureSaveInTable()
		}

	GetName()
		{
		return .name
		}

	GetTable()
		{
		return .table
		}

	EnsureSaveInTable()
		{
		// need add key(name, tab) first to make sure table at least have a key
		// before delete the old key
		Database('ensure ' $ .savein $
			' (name, tab, layout, custom?, hidden?, table_name, user, customKey)
			key(name, tab, customKey, user)')
		if QueryKeys(.savein).Has?('name,tab')
			{
			Database('alter ' $ .savein $ ' drop key(name,tab)')
			QueryKeys.ResetCache()
			}
		if QueryKeys(.savein).Has?('name,tab,user')
			{
			Database('alter ' $ .savein $ ' drop key(name,tab,user)')
			QueryKeys.ResetCache()
			}
		}

	// NOTE: This will return ALL custom fields. IF returning to the user interface
	// 	it is preferable to use GetPermissableFields. This will ensure that custom tab
	// 	permissions are respected
	CustomFields(table = false)
		{
		table = table is false ? .table : table
		return QueryColumns(table).Filter(.CustomField?)
		}

	CustomFieldsByType(table, typeOb)
		{
		return .CustomFields(table).Filter({ typeOb.Has?(DatadictType(it)) })
		}

	UnProtectFields(table)
		{
		customFields = .CustomFields(table)
		QueryApply(.savein $ ' where table_name isnt "" and name is ' $ Display(table))
			{
			key = .CustomTableKey(table)
			customFields.Add(.customTableLinkField(key, it.table_name))
			}
		return customFields
		}

	CustomFieldsStructure(table)
		{
		resultOb = Object()
		fields = .CustomFields(table)
		for field in fields
			resultOb.Add(Object(:field, prompt: Prompt(field), type: DatadictType(field)))

		return resultOb
		}

	CustomField_AllowDelete(field)
		{
		msg = ''
		for func in Contributions('CustomField_AllowDelete')
			if '' isnt contMsg = (func)(.table, field)
				msg $= contMsg
		return msg
		}

	CustomField?(col, includeDeleted = false)
		{
		String(col) =~ '^custom_[0-9]+$' and
			(includeDeleted is true or not Internal?(col))
		}
	DeletedField?(col)
		{
		String(col) =~ '^custom_[0-9]+$' and Internal?(col)
		}

	Layout(sf, tab = '')
		{
		x = .CustomTab(tab)
		layout = x isnt false ? x.layout : .defaultLayout
		return .buildLayout(sf, layout)
		}

	buildLayout(sf, layout)
		{
		if Object?(layout)
			return LayoutToForm.Revert(layout, sf)

		return layout
		}

	DefaultLayout(sf, tab)
		{
		x = .getSavedDefault(tab)
		layout = x isnt false ? x.layout : .defaultLayout
		return .buildLayout(sf, layout)
		}

	LayoutExists?(tab = '')
		{
		x = .CustomTab(tab)
		return x isnt false and x.layout isnt '' or .defaultLayout isnt ''
		}

	CustomTab(tab)
		{
		if not TableExists?(.savein)
			return false

		// check if user has a layout saved, if not load the default one
		customKey = .customKeyForLookup(tab)
		if false is x = Query1(.savein, name: .name, :tab, user: .user, :customKey)
			x = .getSavedDefault(tab)
		return x
		}

	customKeyForLookup(tab)
		{
		return tab is CustomizeExpandControl.LayoutName ? .customKey : ''
		}

	getSavedDefault(tab)
		{
		if not TableExists?(.savein)
			return false

		customKey = .customKeyForLookup(tab)
		return Query1(.savein, name: .name, :tab, user: '', :customKey)
		}

	CustomTableTab?(tab)
		{
		tab = .CustomTab(tab)
		return tab isnt false and tab.table_name isnt ''
		}

	CustomTableName(table_name)
		{
		rec = Query1(.savein, :table_name, user: .user)
		return rec isnt false
			? GetTableName(rec.name) $ ' - ' $ rec.tab
			: table_name
		}

	Form(tab = '', limitHeight = false)
		{
		if false is rec = .CustomTab(tab)
			layoutOb = .defaultLayout
		else
			layoutOb = rec.layout

		if layoutOb is ''
			return ''
		if String?(layoutOb)
			{
			SuneidoLog('WARNING: Old layout format cannot be used', params: rec)
			flds = tab is CustomizeExpandControl.LayoutName and
				.availableFields isnt false
				? .availableFields
				: .GetPermissableFields(.table)
			sf = SelectFields(flds, joins: false)
			return LayoutToForm(layoutOb, sf, onlyCustomFields?:)
			}

		if Object?(layoutOb) and layoutOb.Readonly?()
			layoutOb = layoutOb.DeepCopy()
		.formatDeletedCustomFields(layoutOb)
		.handleRemovedFields(layoutOb)
		.removeEmptyStatic(layoutOb)
		return .ensureHeights(layoutOb, limitHeight)
		}

	ensureHeights(form, limitHeight)
		{
		if not limitHeight
			return form
		for item in form
			{
			if Object?(item)
				LayoutToForm.EnsureEditorHeightLimit(item[0], item)
			}
		return form
		}

	formatDeletedCustomFields(layoutOb)
		{
		for item in layoutOb
			if Object?(item) and .DeletedField?(item[0])
				{
				item[1] = SelectFields.GetFieldPrompt(item[0])
				item[0] = 'Static'
				}
		}

	removeEmptyStatic(layoutOb)
		{
		if not Object?(layoutOb)
			return
		layoutOb.RemoveIf({ Object?(it) and it[0] is 'Static' and it[1].Blank?() })
		}

	handleRemovedFields(layoutOb)
		{
		if .availableFields isnt false
			layoutOb.Each()
				{
				if Object?(it) and it[0] isnt 'Static' and
					not .availableFields.Has?(it[0])
					{
					if .CustomField?(it[0])
						.makeFieldStatic(it, SelectFields.GetFieldPrompt(it[0]))
					else
						{
						if it[0] is prompt = Prompt(it[0])
							{
							.programmerError('field not available', [field: it[0]],
								'datadict deleted or field was renamed')
							prompt = '???'
							}
						.makeFieldStatic(it, prompt)
						}
					}
				}
		}

	makeFieldStatic(fldOb, prompt)
		{
		fldOb[1] = prompt
		fldOb[0] = 'Static'
		}

	programmerError(msg, params, caughtMsg)
		{
		ProgrammerError(msg, params, caughtMsg)
		}

	Table(tab)
		{
		table = .CustomTab(tab).table_name
		key = .CustomTableKey()
		query = table $
			' rename custtable_num to custtable_num_new,
				bizuser_user to bizuser_user_cur
				sort custtable_num_new'
		columns = QuerySelectColumns(query).Remove(#custtable_TS, #custtable_FK)
		name = .customTableLinkField(key, table)
		return [#LineItemControl, :query, :columns, linkField: #custtable_FK,
			protectField: #custtable_protect, :name, headerFields: [key $ '_new'],
			keyField: 'custtable_num_new']
		}

	RebuildLayout(sf, layout)
		{
		newLayout = ''
		sf.ScanFormula(layout,
			{ |field|
			newLayout $= Prompt(field)
			},
			{|s|
			newLayout $= s
			})
		return newLayout
		}
	DeleteField(field)
		{
		QueryApply1(.lib, name: 'Field_' $ field)
			{ |x|
			if false is newtext = .setInternal(x.text)
				return
			x.text = newtext
			x.Update()
			}
		if not .NotCustomizableScreen?()
			QueryDo('delete customizable_fields
				where custfield_field is ' $ Display(field))
		.unloadField(field)
		CustomizableMap.ResetServerCache()
		ServerEval('CustomizableOnServer.NotifyCustomScreenChanges', field)
		SuneidoLog(.DeletedCustomFieldMessage(field))
		}
	DeletedCustomFieldMessage(field)
		{
		return 'INFO: Customizable.DeleteField ' $ field
		}
	setInternal(text)
		{
		if text.Has?('Internal: true')
			return false

		ob = text.SplitOnFirst('{')
		ob[1] = '\tInternal: true\n\t' $ ob[1].Trim()
		return ob[0] $ '{\n' $ ob[1]
		}

	unloadField(field)
		{
		LibUnload('Field_' $ field)
		ServerEval('LibUnload', 'Field_' $ field)
		}

	UpdateField(fieldobject)
		{
		f = fieldobject
		options = f.options

		types = CustomFieldTypes()
		t = types[types.FindIf({ it.name is f.ctllbl })]
		if t.Member?('customOptions')
			{
			customOptFn = Global(t.customOptions)
			options = customOptFn.UpdateProperties(f.colnme, f.colpro, .lib)
			}

		recName = 'Field_' $ f.colnme
		text = Display(.build_custom_dd_text(f.fldbse, f.colpro, recName, options))
		QueryDo('update ' $ .lib $ ' where name = ' $ Display(recName) $
			' set text = ' $ text)
		.unloadField(f.colnme)
		ServerEval('CustomizableOnServer.NotifyCustomScreenChanges', f.colnme)
		}

	build_custom_dd_text(base, prompt, recName, options)
		{
		if not Object?(options)
			gridHeading = ''
		else
			gridHeading = options.GetDefault('dd_gridheading', '')
		gridHeadingText = gridHeading isnt ''
			? '\tGridHeading: ' $ Display(gridHeading) $ '\n'
			: ''
		return base $ '\n' $
			'\t{\n' $
			'\tPrompt: ' $ Display(prompt) $ '\n' $
			gridHeadingText $
			.addSelectPromptIfDupPrompt(recName, prompt) $
			.fieldControlAndFormat(options) $
			'\t}'
		}

	ExtraCreateFieldChecking(name, type)
		{
		types = CustomFieldTypes()
		type = types[types.FindIf({ it.name is type })]

		if type.Member?('customOptions')
			{
			customOptFn = Global(type.customOptions)
			if customOptFn.Method?('ExtraCreateChecking')
				return customOptFn.ExtraCreateChecking(name, customizable: this)
			}
		return ''
		}

	CreateField(name, type, sf, options = #())
		{
		field = .Add_field_def(name, type, options)
		Database("alter " $ .table $ " create (" $ field $ ")")
		CustomizableMap.ResetServerCache()
		QueryColumns.ResetCache()
		.after_field_def(field, name, type, options)
		sf.AddField(field, name)
		.rebuildForms(field, name)
		return field
		}

	rebuildForms(field, name)
		{
		QueryApply('customizable where name is ' $ Display(.table), update:)
			{
			.addCustomFieldToForms(it, field, name)
			}
		}

	addCustomFieldToForms(rec, field, prompt)
		{
		form = .getModifiedForm(rec.layout, field, prompt)
		if rec.layout isnt form
			{
			rec.layout = form
			rec.Update()
			}
		}

	getModifiedForm(origLayout, field, prompt)
		{
		if origLayout is ''
			return origLayout
		form = origLayout.Copy()
		for idx in origLayout.Members()
			{
			ctrl = origLayout[idx]
			if ctrl[0] is 'Static'
				{
				str = ctrl[1]
				if str.Has?(prompt)
					{
					before = str.BeforeFirst(prompt)
					after = str.AfterFirst(prompt)
					beforeStatic = Object('Static', before)
					afterStatic = Object('Static', after)
					customCtrl = Object(field)
					if after isnt ''
						{
						form[idx] = afterStatic
						form.Add(customCtrl at: idx)
						}
					else
						form[idx] = customCtrl
					if before isnt ''
						form.Add(beforeStatic at: idx)
					break
					}
				}
			}
		return form
		}

	FilterTable(list)
		{
		return list.Filter( { it.GetDefault('table', .table) is .table  })
		}

	Add_field_def(name, type, options)
		{
		if not Object?(options)
			options = #()

		if not Libraries().Has?(.lib)
			{
			EnsureConfigLib(.lib)
			ServerEval("Use", .lib)
			Unload()
			}

		types = CustomFieldTypes()
		t = types[types.FindIf({ it.name is type })]

		field = options.Member?('field_name') ?
			options.field_name : .build_field_name(.lib)
		if t.Member?('customOptions')
			{
			customOptFn = Global(t.customOptions)
			options = customOptFn(
				field, name, .lib, sourceTable: .table, :options)
			}
		base = options.GetDefault('dd_base', 'Field_' $ t.base)
		OutputLibraryRecord(.lib,
			Record(name: recName = 'Field_' $ field,
				text: .build_custom_dd_text(base, name, recName, options))
			)
		return field
		}

	ensureLib()
		{
		EnsureConfigLib(.lib)
		ServerEval('Use', .lib)
		}

	after_field_def(field, name, type, options)
		{
		types = CustomFieldTypes()
		type = types[types.FindIf({ it.name is type })]

		if type.Member?('customOptions')
			{
			customOptFn = Global(type.customOptions)
			if customOptFn.Method?('AfterCreate')
				customOptFn.AfterCreate(field, name, options, customizable: this)
			}
		}

	addSelectPromptIfDupPrompt(recName, prompt)
		{
		result = ''
		QueryApply(.lib $ ' where name isnt ' $ Display(recName) $ ' and
			name =~ "^Field_custom_[0-9]+$"')
			{
			if Global(it.name).Prompt is prompt
				{
				result = '\tSelectPrompt: "' $ prompt $ ' ~ ' $
					GetTableName(.table) $ '"\n'
				break
				}
			}
		return result
		}

	fieldControlAndFormat(injectOpts = false)
		{
		if injectOpts is false
			return ''
		str = ''
		ctrl = injectOpts.GetDefault('control', #())
		for m, v in ctrl
			str $= '\tControl_' $ m $ ': ' $ Display(v) $ '\n'
		format = injectOpts.GetDefault('format', #())
		for m, v in format
			str $= '\tFormat_' $ m $ ': ' $ Display(v) $ '\n'
		return str
		}

	build_field_name(lib)
		{
		maxname = QueryMax(lib $
			' where name >= "Field_custom_000000"
			  and name < "Field_custom_999999"
			  and name =~ "^Field_custom_[0-9]+$"',
			'name', 'custom_0')
		num = 1 + Number(maxname.Extract('[0-9]+$'))
		do
			{
			field = 'custom_' $ num.Pad(6) /*= from 000001 to 999999 */
			LibUnload('Field_' $ field)
			dict = Datadict(field)
			num++
			}
		while(dict isnt Field_string)
		return field
		}

	TabCustom?(tab)
		{
		return Query1(.savein, name: .name, :tab, custom?: true, user: .user) isnt false
		}
	ListCustomTabs()
		{
		return QueryList(.savein $
			' where name = ' $ Display(.name) $ ' and custom? is true
				and hidden? isnt true', 'tab').Sort!()
		}

	SaveLayout(layout, sf, tab = '', custom = '', asDefault = false,
		onlyCustomFields? = false)
		{
		for fld in .CustomFields()
			sf.AddField(fld, Prompt(fld))
		form = LayoutToForm(layout, sf, :onlyCustomFields?)
		if form is #(Form)
			form = ''
		Transaction(update:)
			{ |t|
			if false is t.Query1(.saveInTabQuery(tab, asDefault))
				t.QueryOutput(.savein, [
					name: .name,
					:tab,
					customKey: .customKeyForLookup(tab),
					layout: form,
					hidden?: false,
					custom?: custom is true,
					user: asDefault ? '' : .user])
			else
				t.QueryApply1(.saveInTabQuery(tab, asDefault))
					{
					it.tab = tab
					it.layout = form
					it.user = asDefault ? '' : .user
					it.Update()
					}
			}
		}

	SaveTab(tab = '', customTableTab? = false)
		{
		table_name = customTableTab? ? .customTableName() : ''
		Transaction(update:)
			{ |t|
			if false is t.Query1(.saveInTabQuery(tab))
				{
				t.QueryOutput(.savein, [
					name: .name,
					:tab,
					customKey: .customKeyForLookup(tab),
					:table_name,
					hidden?: false,
					custom?:,
					user: .user])
				}
			else
				t.QueryApply1(.saveInTabQuery(tab))
					{
					it.user = .user
					it.hidden? = false
					// Not customTableTab so wipeout previous name if saved
					// else if we are now a customTableTab output the new name
					// or we have a saved value and a customTableTab so dont change
					if not customTableTab?
						it.table_name = ''
					else if it.table_name is ''
						it.table_name = table_name
					it.Update()
					}
				}
		if customTableTab?
			{
			.createCustomTable(tab, table_name)
			.resetCustomTableDataSources()
			}
		}

	CountVisibleCustomTabs()
		{
		return QueryCount(.savein $ ' where name is ' $ Display(.name) $ ' and ' $
			'tab > "" and custom? is true and hidden? isnt true')
		}

	customTableName()
		{
		return 'custom_table_' $ Display(Timestamp()).Tr('#.')
		}

	saveInTabQuery(tab, asDefault = false)
		{
		return .savein $ ' where name = ' $ Display(.name) $
			' and tab is ' $ Display(tab) $
			' and user is ' $ (asDefault is true ? '""' : Display(.user) $
			' and customKey is ' $ Display(.customKeyForLookup(tab)))
		}

	defaultFieldOptions: (
		format: (width: 25, height: 4),
		control: (width: 25, status: '', tabthrough: false, height: 4)
		)
	createCustomTable(tab, table_name)
		{
		.ensureLib()
		key = .CustomTableKey()
		ctrlName = .customTableLinkField(key, table_name)
		// Output the new custom table
		Database('ensure ' $ table_name $
			'(bizuser_user, custtable_TS, custtable_num, custtable_FK)
			key (custtable_num)
			index (custtable_FK) in ' $ .table $ '(' $ key $ ') cascade')

		// Output the Rule/Field required to properly link the
		// new custom table to the standard table.
		// Rule is required for the foreign key
		OutputLibraryRecord(.lib, [
			name: 'Rule_' $ ctrlName,
			text: 'function (){ return .' $ key $ '_new }'
			])
		// Field is required for the data validation error messages
		OutputLibraryRecord(.lib, [
			name: 'Field_' $ ctrlName,
			text: 'Field_' $ key $ '_new { Prompt: ' $ Display(tab) $ ' }'
			])

		// Output a basic custom field for user notes onto the new custom table
		Customizable(table_name).
			CreateField(tab $ ' Notes', 'Text, multi line', SelectFields(),
				.defaultFieldOptions)
		}

	CustomTableKey(table = false)
		{
		table = table is false ? .table : table
		return QueryKeys(table).FindOne({ it.Suffix?('_num') })
		}

	customTableLinkField(key, table)
		{
		return key $ '_' $ table.AfterLast('_')
		}

	resetCustomTableDataSources()
		{
		CustomTableDataSources.ResetCache()
		}

	RenameLayout(oldTab, newTab)
		{
		.updateTab(oldTab, [tab: newTab])
		}

	updateTab(tab, newValues)
		{
		resetCustomTableCache? = false
		QueryApply1(.saveInTabQuery(tab))
			{ |x|
			for m, v in newValues
				x[m] = v
			x.Update()
			if x.table_name isnt ''
				resetCustomTableCache? = true
			}
		if resetCustomTableCache?
			.resetCustomTableDataSources()
		}

	HideLayout(tab)
		{
		.showHideLayout(tab, hide:)
		}
	ShowLayout(tab)
		{
		.showHideLayout(tab, hide: false)
		}
	showHideLayout(tab, hide = false)
		{
		.updateTab(tab, [hidden?: hide])
		}
	LayoutHidden?(tab)
		{
		if false is rec = Query1(.saveInTabQuery(tab))
			return false
		return rec.hidden? is true
		}
	ListHiddenLayouts()
		{
		return QueryList(.savein $
			' where name = ' $ Display(.name) $ ' and hidden? is true', 'tab').Sort!()
		}

	ListVisibleLayouts()
		{
		return QueryList(.savein $
			' where name is ' $ Display(.name) $ ' and hidden? isnt true', 'tab').Sort!()
		}

	AvailableCustomFields()
		{
		ob = Object()
		form = .Form()
		for col in .CustomFields()
			if not form.Any?({ Object?(it) and it[0] is col })
				ob[Prompt(col)] = col
		return ob
		}
	PromptToField(prompt)
		{
		customFields = .CustomFields()
		columnIndex = customFields.FindIf({ Prompt(it) is prompt})
		return columnIndex isnt false ? customFields[columnIndex] : false
		}
	GetCustomFieldsPerTab(tab = 'Custom')
		{
		if TableExists?(.savein) and
			false isnt rec = Query1(.savein, :tab, name: .name, user: .user)
			{
			customFields = .CustomFields()
			fields = rec.layout.Filter({ Object?(it) and customFields.Has?(it[0]) })
			return fields.Map({ it[0] })
			}
		return #()
		}

	GetNonPermissableFields(query)
		{
		if query is false
			return #()

		table = QueryGetTable(query, nothrow:, orview:)
		if table is "" or not TableExists?('customizable') or not .hasCustomTabs?(table)
			return #()

		removeOb = Object()
		if false isnt cl = OptContribution('CustomTabPermissions', false)
			removeOb = cl.NonPermissableFieldsFromTable(table)
		return removeOb
		}

	GetPermissableFields(query)
		{
		customOb = .CustomFields(query)
		removeOb = .GetNonPermissableFields(query)
		return customOb.Difference(removeOb)
		}

	GetPermissableDataSources()
		{
		sources = ServerEval(#CustomTableDataSources)
		return sources.RemoveIf({ not AccessPermissions(it.auth) })
		}

	CustomTablePrefix: 'Custom Table > '
	AccessToDataSource?(sourceName, defaultRtn = false)
		{
		return sourceName.Prefix?(.CustomTablePrefix)
			? .GetPermissableDataSources().Member?(sourceName)
			: defaultRtn
		}

	CustomReporterReports(tab)
		{
		custtable_source = .CustomTablePrefix $ GetTableName(.name) $ ' > ' $ tab
		return QueryAll('params extend custtable_source', :custtable_source)
		}

	hasCustomTabs?(table)
		{
		return not QueryEmpty?('customizable', name: table, custom?: true)
		}

	GetEditableCustomFields(table, custfield_name)
		{
		custFields = Object()
		for field in Customizable.CustomFields(table)
			if false is Query1('customizable_fields
				where custfield_readonly is true or custfield_hidden is true',
				:custfield_name, custfield_field: field)
				custFields.Add(field)
		return custFields
		}


	// TODO: Move the following methods to CustomizeField or separate record

	DeleteCustomFields(t, oldrec, newrec, fields)
		{
		if .NotCustomizableScreen?()
			return

		for field in fields.Members()
			if oldrec[field] isnt newrec[field]
				t.QueryDo('delete customizable_fields
					where custfield_field is ' $ Display(fields[field].AfterFirst('_')))
		}

	SetRecordDefaultValues(key, rec, protectField = false, useDefaultsIfEmpty? = false)
		{
		if key is false
			return
		values = .CacheByKey(key, 'CustomizedDefaultValues')
			{|name|
			custom = false
			if not .NotCustomizableScreen?()
				{
				custom = Object()
				QueryApply('customizable_fields
					where custfield_name is ' $ Display(name) $ ' and
						custfield_default_value isnt ""
					sort custfield_num')
					{ |x|
					.SetDefaultValue(x, custom)
					}
				if custom.Empty?()
					custom = false
				}
			custom
			}

		if values isnt false
			{
			if protectField is false
				{
				if useDefaultsIfEmpty?
					rec.MergeNew(values)
				else
					rec.Merge(values)
				rec.CustomizableSetDefaultValues = true
				}
			else
				.handleProtectField(rec, protectField, values)

			.ensureDeps(rec, values)
			}
		}
	ensureDeps(rec, values)
		{
		if rec.CustomizableSetDefaultValues isnt true
			return

		newRec = rec.Copy()
		for field in values.Members()
			{
			if newRec.Member?(field)
				{
				newRec.Delete(field)
				newRec[field]
				if '' isnt deps = newRec.GetDeps(field)
					rec.SetDeps(field, deps)
				}
			}
		}
	handleProtectField(rec, protectField, values)
		{
		protect = rec[protectField]
		if .invalidProtect(protect)
			{
			.logError('invalid return type from protect rule')
			return
			}

		if protect is true or (String?(protect) and protect isnt '')
			return

		rec.Merge(values)
		rec.CustomizableSetDefaultValues = true

		protect = rec[protectField]
		.handleProtectedAfterMerged(protect, rec, values)
		}
	invalidProtect(protect)
		{
		return not String?(protect) and not Boolean?(protect) and not Object?(protect)
		}
	handleProtectedAfterMerged(protect, rec, values)
		{
		if Object?(protect)
			.handleProtectOb(rec, protect, values)
		else
			for field in values.Members()
				if rec.Member?(field) and rec[field $ '__protect'] is true
					rec.Delete(field)
		}
	handleProtectOb(rec, protect, values)
		{
		allbut? = protect.GetDefault(0, false) is 'allbut'
		for field in values.Members()
			if (rec.Member?(field) and (protect.Member?(field) isnt allbut? or
				rec[field $ '__protect'] is true))
				rec.Delete(field)
		}
	logError(str, prefix = "ERROR: ")
		{
		SuneidoLog(prefix $ str)
		}

	SetRecordDefaultValue(key, rec, field)
		{
		if .NotCustomizableScreen?()
			return

		QueryApply('customizable_fields
			where custfield_name is ' $ Display(key) $ ' and
				custfield_field is ' $ Display(field) $ ' and
				custfield_default_value isnt ""
			sort custfield_num')
			{ |x|
			.SetDefaultValue(x, rec)
			}
		}

	SetDefaultValue(x, custom)
		{
		dd = Datadict(x.custfield_field)
		control = dd.Control[0]

		if not control.Suffix?('Control')
			control $= 'Control'
		ctrl = Global(control)
		if ctrl.Method?('CustomizableSetDefaultValue') and
			true is ctrl.CustomizableSetDefaultValue(x, custom, dd)
				return
		custom[x.custfield_field] = x.custfield_default_value
		}

	GetCustomizedFields(key)
		{
		return .CacheByKey(key, 'CustomizedFields')
			{|name|
			custom = false
			if not .NotCustomizableScreen?()
				{
				custom = Object()
				QueryApply('customizable_fields', custfield_name: name)
					{|cust_rec|
					for opt in CustomFieldOptions()
						if cust_rec[opt.field] isnt false and cust_rec[opt.field] isnt ""
							{
							field = cust_rec.custfield_field
							if not custom.Member?(field)
								custom[field] = Object()
							custom[field][opt.field_option] = cust_rec[opt.field]
							}
					}
				if custom.Empty?()
					custom = false
				}
			custom
			}
		}

	GetTabPermissionName(tabName)
		{
		return .AuthPath(GetTableName(.name), tabName)
		}

	AuthPath(tableName, tabName)
		{
		// encode so that `/`'s in tab name will not conflict with permission option
		return "/Custom Tab/" $ tableName $ '/' $ Url.EncodeQueryValue(tabName)
		}

	CacheByKey(key, cache_name, block)
		{
		if key is false
			return false

		if not Suneido.Member?(cache_name)
			Suneido[cache_name] = Object()

		cache = Suneido[cache_name]
		if not cache.Member?(key)
			cache[key] = block(key)

		return cache[key]
		}

	cache_list: ('CustomizedFields', 'CustomizedDefaultValues', 'CustomizedFieldFormulas')
	ResetCustomizedCache(key)
		{
		if key is false
			return

		cache_names = .cache_list
		for cache_name in cache_names
			if Suneido.Member?(cache_name) and Suneido[cache_name].Member?(key)
				Suneido[cache_name].Delete(key)
		}

	ResetServerCustomizedCache(key)
		{
		ServerEval('Customizable.ResetCustomizedCache', key)
		}

	ResetAllCache()
		{
		Suneido.Delete(@.cache_list)
		}

	HandleChangesOnClient(changes)
		{
		customFieldsChanges = changes.GetDefault('custom_fields', #())
		for cache_name in .cache_list
			if Suneido.Member?(cache_name)
				for customKey in customFieldsChanges
					Suneido[cache_name].Delete(customKey)

		custom_screen = changes.GetDefault('custom_screen', #())
		for field in custom_screen
			LibUnload('Field_' $ field)
		}

	NotCustomizableScreen?()
		{
		if not TableExists?('customizable_fields')
			return true

		.EnsureTable()
		return false
		}

	EnsureTable()
		{
		Database("ensure customizable_fields
			(custfield_num, custfield_name, custfield_field, custfield_mandatory,
			custfield_readonly, custfield_hidden, custfield_tabover,
			custfield_first_focus, custfield_only_fillin_from, custfield_default_value,
			custfield_formula, custfield_formula_code, custfield_formula_fields,
			bizuser_user, custfield_date_modified, custfield_user_modified)
			key(custfield_name, custfield_field)
			key(custfield_num)")
		}

	IsCustomized?(query)
		{
		if query is ''
			return false

		table = QueryGetTable(query)
		if table is ''
			return false
		where = 'where name is ' $ Display(table) $
			' and tab isnt ' $ Display(CustomizeExpandControl.LayoutName)
		QueryApply('customizable ' $ where)
			{|x|
			if x.layout isnt '' or (x.table_name isnt '' and x.hidden? isnt true)
				return true
			}

		return false
		}

	TabCustomized?(tab)
		{
		if not TableExists?(.savein)
			return false

		tab = Query1(.savein, :tab, name: .name)
		return tab isnt false
			? tab.layout isnt '' and tab.hidden? isnt true
			: false
		}
	}