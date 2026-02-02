// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	// Customizable Controls can use:
	// Customizable_fieldSetVal to format/set data before displaying
	// Customizable_fieldGetVal to format/set data before saving
	Name: "CustomizeFields"

	New(.key, .query, sfOb, .browse? = false, .readonly = false,
		.disable_default = false, .virtual_list? = false)
		{
		super(.layout(sfOb))
		.default_value_ctrl = .FindControl('custfield_default_value')
		.explorer_list_view = .FindControl('ExplorerListView')
		}
	browse?: false
	dirty?: false
	layout(sfOb)
		{
		Customizable.EnsureTable()

		.table = QueryGetTable(.query)
		.fields_list = SelectFields(sfOb.cols, sfOb.excludeFields, false).Fields.Values()
		formulaExcludeFields = .getFormulaExcludeFields(sfOb)
		.formulaSelectFields = SelectFields(sfOb.cols, formulaExcludeFields, false)

		return Object('Vert'
			Object('ExplorerListView',
				.explorerModel(),
				.build_layout_controls(.disable_default)
				columns: #(custfield_field, custfield_options,
					bizuser_user, custfield_date_modified, custfield_user_modified),
				readonly: .readonly,
				validField: 'custfield_valid',
				protectField: 'custfield_protect'
				buttonBar: true
				extraButtons: #(Fill OkButton)
				columnsSaveName: .Name
				historyFieldsPrefix: 'custfield'
				))
		}
	getFormulaExcludeFields(sfOb)
		{
		numAndExcludeFields = sfOb.cols.Copy().RemoveIf(
			{ not (it.Suffix?('_num') or it.Has?('_num_')) })
		excludeControls = Object('Key', 'Id').Append(
			GetContributions('FormulaExcludeControls'))
		numAndExcludeFields.RemoveIf({ not excludeControls.Has?(Datadict(it).Control[0])})
		for field in sfOb.cols.Copy()
			if .addToExclude?(field, excludeControls)
				numAndExcludeFields.Add(field)

		return sfOb.excludeFields.Copy().Append(numAndExcludeFields)
		}

	addToExclude?(field, excludeControls)
		{
		dd = Datadict(field)
		if field.Prefix?('custom_')
			{
			if excludeControls.Has?(dd.Control[0])
				{
				parent = Name(dd.Base()) // check immediate inheritance
				grandparent = Name(dd.Base().Base()) // and one level up
				if parent.Suffix?('_num') or parent.Has?('_num_') or
					grandparent.Suffix?('_num') or grandparent.Has?('_num_')
					return true
				}
			}
		return dd.Member?('NoFormulas') and dd.NoFormulas
		}

	build_layout_controls(disable_default)
		{
		form = Object('Form', #(custfield_field group: 0) #nl #nl #nl
			#(custfield_mandatory group: 0), #(custfield_readonly group: 1) #nl
			#(custfield_tabover group: 0))

		if .browse? isnt true or .virtual_list?
			form.Add(#(custfield_hidden group: 1))
		if .browse? isnt true
			form.Add(#(custfield_first_focus group: 2))

		form.Add(#nl, #(custfield_only_fillin_from group: 0))

		ob = Object(#Vert #Skip form)

		defVal = Object('CustomizeFieldsDefaultValue' name: 'custfield_default_value')
		if disable_default is false
			{
			defVal.Add('Field')
			ob.Add(#(Skip 3))
			ob.Add(Object('Pair',
				#(Static 'Default Value', name: 'default_value_prompt'),
				defVal))
			ob.Add(#(Skip 3))
			ob.Add(#('custfield_formula'))
			if LastContribution('AllowFormulaEdit')()
				ob.Add(FormulaEditor.AddFormulaButtons())
			}
		else
			ob.Add(defVal)

		ob.Add(#Skip)
		return ob
		}

	explorerModel()
		{
		return Object('ExplorerListModel'
			'customizable_fields
				rename custfield_num to custfield_num_new
				extend custfield_fields_list = ' $ Display(.fields_list) $
				' , custfield_browse? = ' $ Display(.browse?) $
				' where custfield_name is ' $ Display(.key) $
				' sort custfield_num_new',
			#(custfield_name, custfield_field))
		}

	Recalc()
		{
		if false isnt promptCtrl = .FindControl('default_value_prompt')
			{
			formula = .FindControl('custfield_formula')
			promptCtrl.Xmin = formula.Parent.Left
			if false isnt addFieldSkip = .FindControl('add_field_skip')
				addFieldSkip.Xmin = formula.Parent.Left
			}
		super.Recalc()
		}

	ExplorerListView_CurrentPrintSavedName()
		{
		return .Name
		}

	ExplorerListView_AddRecord(rec)
		{
		rec.custfield_name = .key
		rec.custfield_browse? = .browse?
		rec.custfield_fields_list = .fields_list
		rec.bizuser_user = Suneido.User

		.default_value_ctrl.RemoveAll()
		.append_control('string')
		}

	ExplorerListView_AfterSave(rec)
		{
		.processRecordBeforeLoad(rec)
		}
	ExplorerListView_Selection(rec)
		{
		.ExplorerListView_BeforeEntryLoaded(rec)
		}
	ExplorerListView_BeforeEntryLoaded(rec)
		{
		if .Destroyed?()
			return

		.processRecordBeforeLoad(rec)
		.default_value_ctrl.RemoveAll()
		if rec isnt false
			.append_control(rec.custfield_field)

		.custom_options(rec) {|fld, deflt| if deflt rec.PreSet(fld, true)}
		rec.PreSet('prev_custfield_field', rec.custfield_field)
		}
	processRecordBeforeLoad(rec)
		{
		.handleControlCustomizedMessage('Customizable_fieldGetVal', rec, Object(preset?:))
		}
	ExplorerListView_RecordChanged(member, data)
		{
		if member is 'custfield_field'
			{
			.default_value_ctrl.RemoveAll()
			field = data.custfield_field
			if field isnt ''
				{
				data.custfield_default_value = ''
				data.custfield_formula = ''
				.append_control(field)
				data.Invalidate('custfield_protect')
				}
			.custom_options(data) {|fld, deflt| data[fld] = deflt ? true : ""}
			}
		if member is 'custfield_default_value'
			{
			v = data.custfield_default_value
			if String?(v) and v.Trim() isnt v
				data.custfield_default_value = v.Trim()
			}
		if member is 'custfield_formula'
			{
			fn = CustomizeField.TranslateFormula(.formulaSelectFields,
				data.custfield_formula, data.custfield_field, quiet:)
			data.custfield_formula_code = fn.formulaCode
			data.custfield_formula_fields = fn.fields
			}
		.dirty? = true
		}
	ExplorerListView_BeforeSave(rec)
		{
		.custom_options(rec) {|field, deflt| if deflt rec.PreSet(field, "") }
		.handleControlCustomizedMessage('Customizable_fieldSetVal', rec)
		}
	ExplorerListView_AfterRestore(rec)
		{
		.ExplorerListView_BeforeEntryLoaded(rec)
		.handleControlCustomizedMessage('Customizable_fieldGetVal', rec)
		}
	custom_options(rec, block)
		{
		dict = Datadict(rec.custfield_field)
		for opt in CustomFieldOptions()
			block(opt.field, dict.Control.GetDefault(opt.field_option, false) is true)
		}
	ExplorerListView_BeforeRestore(rec)
		{
		.ExplorerListView_BeforeSave(rec)
		}
	handleControlCustomizedMessage(method, rec, args = false)
		{
		if args is false
			args = Object()

		.invokeControlMethod(method, GetControlClass.FromField(rec.custfield_field),
			args.MergeNew([:rec]))
		}
	invokeControlMethod(method, ctrl, args)
		{
		if ctrl.Method?(method)
			(ctrl[method])(@args)
		}
	no_default_value: #(OpenImage, OpenImageWithLabels)
	append_control(field)
		{
		control = Datadict(field).Control.Copy()
		if .no_default_value.Has?(control[0])
			return

		control.saveInfoName = field
		.Set_options(control)
		if .disable_default is false
			.default_value_ctrl.Append(control)
		}
	Set_options(control)
		{
		control.mandatory = false
		ctrlOpts = .controlOptions().GetDefault(control[0].RemoveSuffix(#Control), false)
		if false isnt ctrlOpts
			ctrlOpts(control)
		}

	controlOptions()
		{
		return GetContributions(#CustomizableControlOptions).MergeNew(Object(
			Info: 			{ |control| control.allowOnlyType = true },
			ChooseMany: 	{ |control| control.saveNone = false },
			RadioButtons:
				{ |control|
				if control.GetDefault(#noInitalValue, false)
					control.mandatory = true
				},
			UOM:
				{ |control|
				#(1, 2).Each()
					{
					if not control.Member?(it)
						continue
					control[it] = control[it].Copy()
					control[it].mandatory = false
					}
				}))
		}

	FormulaEditor_Click(source)
		{
		FormulaEditor.HighlightSelection(source, .formulaSelectFields)
		}
	On_Add(option)
		{
		FormulaEditor['Add_a_' $ option.Tr(' ', '_')](.findControl(),
			selectFields: .formulaSelectFields, hwnd: .Window.Hwnd)
		.explorer_list_view.Valid?(evalRule?:)
		}

	findControl()
		{
		if .formulaDisabled()
			return false
		return .FindControl('custfield_formula')
		}

	formulaDisabled(rec = false)
		{
		if not LastContribution('AllowFormulaEdit')()
			return true

		data = rec is false ? .getCurrent() : rec
		protect = data.custfield_protect
		if protect.GetDefault('custfield_formula', false) is true
			{
			.AlertInfo('Customize Fields',
				'You cannot add Formula since if there is a system formula on it')
			return true
			}
		return false
		}
	getCurrent()
		{
		return .FindControl('ExplorerListView').GetView().Get()
		}
	On_Add_an_Operator(option)
		{
		FormulaEditor.Add_an_Operator(.findControl(), option)
		.explorer_list_view.Valid?(evalRule?:)
		}
	On_Add_a_Function(fnName)
		{
		FormulaEditor.Add_a_Function(.findControl(), fnName, .formulaSelectFields)
		.explorer_list_view.Valid?(evalRule?:)
		}
	GetMandatoryCustomFields()
		{
		data = .FindControl('ExplorerListView').GetList().Get()
		return data.
			Filter({
				Customizable.CustomField?(it.custfield_field) and
				it.custfield_mandatory is true }).
			Map({ it.custfield_field })
		}
	GetList()
		{
		return .FindControl('ExplorerListView').GetList().Get()
		}
	On_OK()
		{
		.Send('OnOK')
		}
	Save()
		{
		if not .readonly and .dirty?
			{
			Customizable.ResetCustomizedCache(.key)
			Customizable.ResetServerCustomizedCache(.key)
			}
		ctrl = .FindControl('ExplorerListView')
		data = ctrl.GetList().Get()
		if 1 < data.MembersIf({|m| data[m].custfield_first_focus is true }).Size()
			{
			.AlertInfo('Customize', 'Only one field can be set to First Focus.')
			return 'invalid'
			}

		if .trySave(ctrl)
			return .dirty?
		return 'invalid'
		}
	trySave(ctrl)
		{
		if ctrl is false
			return false
		if ctrl.SaveRecord()
			return true
		if CloseWindowConfirmation() is true
			{
			.FindControl('ExplorerListView').On_Current_Restore()
			return true
			}
		return false
		}

	FieldRenamed(field)
		{
		if false is view = .FindControl('ExplorerListView').GetView()
			return

		if false is fieldCtrl = view.FindControl('custfield_field')
			return

		if fieldCtrl.Get() is field
			{
			// call FieldPromptControl.Set to trigger refreshing the prompt list and
			// updating the display prompt value
			fieldCtrl.Set(field)
			}
		}
	}
