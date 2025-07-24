// Copyright (C) 2013 Axon Development Corporation. All rights reserved worldwide.
RepeatControl
	{
	Name: 'Filters'
	useCheckBoxes: false
	// NOTE: could also use SelectFields by handling FieldPrompt_GetSelectFields msg
	New(.query = false, exclude = #(), additionalFields = #(), .defaultFilterFields = #(),
		.fieldOptional = false, .useCheckBoxes = false, maxRecords = false)
		{
		super(.layout(exclude, additionalFields),
			minusAtFront:, plusAtBottom:, disableFieldProtectRules:, :maxRecords)

		// Suggestion 27034 added aliasing paramSelect fields when no checkboxes but
		// didnt document reason. Need a way to stop aliasing without checkbox for 33474
		// so FieldPrompt DD's without checkboxes match FieldPrompt DD's with checkboxes
		if 0 is .aliasParamSelectFields = .Send('ChooseFilters_AliasParamSelectFields?')
			.aliasParamSelectFields = not .useCheckBoxes
		}

	prevConditions: false
	Set(data)
		{
		.prevConditions = Object?(data) ? data.Copy() : Object()
		data = ChooseFiltersSetDefaults(.defaultFilterFields, data)
		super.Set(data)
		}

	layout(exclude, additionalFields)
		{
		horz = Object(#Horz)
		if .useCheckBoxes
			horz.Add(#(Skip, large:) #(CheckBox, name: 'check'), #(Skip, small:))
		horz.Add(
			Object('FieldPrompt',
				fields: additionalFields,
				table: .query,
				exclude_fields: exclude,
				width: 15,
				mandatory: not .fieldOptional,
				name: 'condition_field'),
			Object('Vert', name: 'paramSelectWrapper'))
		// use WndPane for correct tab order
		return Object('WndPane', horz, windowClass: "SuBtnfaceArrow")
		}

	Record_NewValue(field, value, source = false)
		{
		if field is 'condition_field'
			.appendParamSelect(source, value)
		if .useCheckBoxes and Object?(value) and value.Member?('operation') and
			source isnt false and not .focusOnCheckBox?(source)
			source.SetField('check', true)
		.notifyIfChanged()
		}

	focusOnCheckBox?(source)
		{
		checkBox = source.FindControl('check')
		return checkBox isnt false and checkBox.HasFocus?()
		}

	notifyIfChanged()
		{
		if .conditionChanged(.prevConditions)
			.Send('ChooseFiltersChanged')
		.prevConditions = .Get().Copy()
		}

	appendParamSelect(source, selectedField)
		{
		wrapper = source.FindControl('paramSelectWrapper')
		Assert(wrapper.Tally() lessThanOrEqualTo: 1)

		.keepFocus(source)
			{
			wrapper.Remove(0)
			.resetFilter(source)
			if not .aliasParamSelectFields
				ctrlField = source.GetField('condition_field')
			else
				ctrlField = source.FindControl('condition_field').AliasParamSelectField()
			wrapper.Insert(0, Object('ParamsSelect', ctrlField, paramPrompt: '',
				name: selectedField))
			}
		}

	keepFocus(source, block)
		{
		prevFocus = GetFocus()
		chooseBtnName = Sys.SuneidoJs?() ? 'operation' : 'ChooseButton'
		chooseBtn = source.FindControl(chooseBtnName)
		chooseBtnFocused = chooseBtn isnt false and
			(prevFocus is chooseBtn.Hwnd)

		block()

		if chooseBtnFocused and
			false isnt newChooseBtn = source.FindControl(chooseBtnName)
			newChooseBtn.SetFocus()
		else
			SetFocus(prevFocus)
		}

	resetFilter(source)
		{
		data = source.Get()
		if not data.Member?('condition_field')
			return
		field = data.condition_field
		if field isnt ""
			data[field] = #(operation: '', value: '', value2: '')
		data.DeleteIf({ |mem| not Object(field, 'condition_field').Has?(mem) })
		}

	On_Minus(source)
		{
		super.On_Minus(source)
		.notifyIfChanged()
		}

	On_Plus(source)
		{
		super.On_Plus(source)
		source = .GetRows().Last()
		.appendParamSelect(source, "")
		}

	// do not want use Get in super because remove in RepeatControl kicks in rules
	Get()
		{
		rows = .GetRows().Map(#Get)
		newConditions = Object()
		for row in rows
			{
			if row.Empty?()
				continue
			field = row.condition_field
			if field is ''
				continue
			newCondition = [condition_field: field]
			if row.Member?(field)
				newCondition[field] = row[field]
			else // dont want to save "" because of rules kicking in on fields
				newCondition[field] = #(operation: "", value: "", value2: "")
			if .useCheckBoxes
				newCondition['check'] = row.check is true
			newConditions.Add(newCondition)
			}
		return newConditions
		}

	HasFieldCondition?(filters, field)
		{
		return .GetFilterValue(filters, field) isnt false
		}

	GetFilterValue(filters, field)
		{
		if not Object?(filters)
			return false

		return filters.FindOne(
			{ it.condition_field is field and it[field].operation isnt '' })
		}

	conditionChanged(prevConditions)
		{
		return .ActiveConditions(prevConditions) isnt .ActiveConditions()
		}

	ActiveConditions(conditions = false)
		{
		if conditions is false
			conditions = .Get()
		if not Object?(conditions)
			return #()
		return conditions.Filter(.activeCondition?)
		}

	activeCondition?(row)
		{
		if '' is field = row.condition_field
			return false

		if not row.Member?(field)
			return false

		condition = row[field]
		if not Object?(condition)
			return false

		return condition.GetDefault("operation", "") isnt ""
		}

	BeforeRowSet(rowCtrl, rowData)
		{
		wrapper = rowCtrl.FindControl('paramSelectWrapper')
		field = rowData.condition_field
		if not .aliasParamSelectFields
			ctrlField = field
		else
			ctrlField =
				rowCtrl.FindControl('condition_field').AliasParamSelectField(field)
		wrapper.Remove(0)
		wrapper.Insert(0,
			Object('ParamsSelect', ctrlField, paramPrompt: '', name: field))
		}

	ClearFilter()
		{
		.setFilterValue([operation: '', value: '', value2: ''])
		}

	SetFilterValue(valueOb, field)
		{
		.setFilterValue(valueOb, field)
		}

	setFilterValue(valueOb, field = false)
		{
		updated? = false
		for row in .GetRows()
			{
			rowData = row.Get()
			rowField = rowData.condition_field
			if field isnt false and rowField isnt field or rowField is ''
				continue

			if ((false is paramSelect = row.FindControl(rowField)) or
				(not paramSelect.Base?(ParamsSelectControl)))
				continue

			paramSelect.Set(valueOb)
			paramSelect.Send("NewValue", paramSelect.Get())
			updated? = true
			valueOb = [operation: '', value: '', value2: '']
			}
		return updated?
		}

	AppendFilter(valueOb, field)
		{
		super.On_Plus(false)
		source = .GetRows().Last()
		source.Get().condition_field = field
		.appendParamSelect(source, field)
		.setFilterValue(valueOb, field)
		}

	BuildWhereFromFilter(filter = false, skipFieldConditionFn = false, fields = false,
		useCheckBoxes = false, conditionFields = false)
		{
		if filter is false
			filter = .Get()
		if useCheckBoxes is false
			useCheckBoxes = .useCheckBoxes
		where = ''
		invalidConditions = Object()
		for condition in filter
			{
			field = .fieldFromCondition(condition, invalidConditions)
			if .skip?(field, useCheckBoxes, condition)
				continue
			if Object?(fields) and not fields.Has?(field)
				{
				SuneidoLog("ERROR: Invalid condition field: " $ field, calls:)
				continue
				}
			if .hasCondition(field, skipFieldConditionFn, condition)
				{
				where $= GetParamsWhere(condition.condition_field, data: condition)
				.collectConditionFields(conditionFields, condition)
				}
			}
		.logInvalidConditions(invalidConditions)
		return where
		}

	collectConditionFields(conditionFields, condition)
		{
		if conditionFields isnt false
			conditionFields.AddUnique(condition.condition_field)
		}

	logInvalidConditions(invalidConditions)
		{
		if not invalidConditions.Empty?()
			SuneidoLog("ERROR: Invalid conditions", calls:, params: invalidConditions)
		}

	fieldFromCondition(condition, invalidOb)
		{
		if condition.Member?('condition_field')
			return condition.condition_field
		invalidOb.Add(condition)
		return ''
		}

	skip?(field, useCheckBoxes, condition)
		{
		return field is '' or (useCheckBoxes and condition.check isnt true)
		}

	hasCondition(field, skipFieldConditionFn, condition)
		{
		return not .skipField(field, skipFieldConditionFn) and
			condition.Member?(field) and
			Object?(condition[field]) and
			condition[field].GetDefault('operation', "") isnt ""
		}

	GetQuery()
		{
		return .query $ ' ' $ .BuildWhereFromFilter()
		}

	skipField(field, conditionFn)
		{
		return conditionFn is false ? false : (conditionFn)(field)
		}

	GetFieldName()
		{
		return .Name
		}

	Valid?()
		{
		return super.Valid?(checkAll: .fieldOptional)
		}

	ForceValid()
		{
		invalid = Object()
		.GetRows().Each()
			{|c|
			if false isnt ctrl = c.FindControl('condition_field')
				.checkRow(c, ctrl.Field.Get(), invalid)
			}
		missingField = invalid.Extract('missingField', false)
			? 'Please select a field\n'
			: ''
		return missingField $ Opt('Invalid: ', invalid.Join(', '))
		}

	checkRow(c, fieldPrompt, invalid)
		{
		rowRec = c.Get()
		field = rowRec.condition_field
		if not .skip?(field, .useCheckBoxes, rowRec) and c.Valid(forceCheck:) isnt true
			invalid.AddUnique(fieldPrompt)
		else if field is '' and fieldPrompt isnt ''
			invalid.AddUnique(fieldPrompt)
		else if .useCheckBoxes and rowRec.check is true and field is ''
			invalid.missingField = true
		}

	FocusRow(rowIndex, value? = false)
		{
		rows = .GetRows()
		if rows.Member?(rowIndex)
			{
			rowCtrl = rows[rowIndex]
			fieldList = rowCtrl.FindControl('condition_field')
			if not value?
				fieldList.SetFocus()
			else
				if false isnt selCtrl = rowCtrl.FindControl(fieldList.Get())
					selCtrl.FocusValue()
			}
		}
	}
