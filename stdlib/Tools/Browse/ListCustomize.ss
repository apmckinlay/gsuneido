// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	InitCustomFields(table, columns)
		{
		custom_fields = Object()
		if not Object?(columns)
			custom_fields = false

		if false is .checkCustomFieldCreation(QueryGetTable(table, nothrow:))
			custom_fields = false

		return custom_fields
		}

	checkCustomFieldCreation(table)
		{
		if table is ''
			{
			Alert("Could not create a custom field for this table",
				"Create Custom Column", flags: MB.ICONERROR)
			SuneidoLog('WARNING: Could not create a custom field', calls:)
			return false
			}
		return true
		}

	SubTitle(ctrl, linkField)
		{
		if linkField is false
			return ''

		if 0 is sub_title = ctrl.Send('TabGetSelectedName')
			sub_title = ''
		return sub_title $ ' List'
		}

	AddCustomColumns(query, columns)
		{
		if columns is false or query is ''
			return columns
		new_cols = columns.Copy()
		for c in .getPermissableFields(query)
			new_cols.AddUnique(c)
		return new_cols
		}

	getPermissableFields(query)
		{
		return Customizable.GetPermissableFields(query)
		}

	GetControlOption(customFields, field, option)
		{
		field_def = Datadict(field)
		return (field_def.Control.GetDefault(option, false) is true) or
			(customFields isnt false and customFields.Member?(field) and
				customFields[field].GetDefault(option, false) is true)
		}

	HandleCustomizableFields(customKey, record, protectField, useDefaultsIfEmpty? = false)
		{
		Customizable.SetRecordDefaultValues(customKey, record, protectField,
			:useDefaultsIfEmpty?)
		CustomizeField.SetFormulas(customKey, record, protectField)
		}

	MandatoryAndEmpty?(rec, field, customFields, protectField)
		{
		return .mandatoryAndEmpty?(rec, field, customFields) and
			not FieldProtected?(field, rec, protectField)
		}

	mandatoryAndEmpty?(rec, field, customFields)
		{
		return .GetControlOption(customFields, field, 'mandatory') and
			CustomizeField.ConsideredEmpty?(field, rec) and
			not .GetControlOption(customFields, field, 'readonly')
		}

	InvalidField?(rec, field, customFields, protectField)
		{
		if FieldProtected?(field, rec, protectField)
			return false
		return .mandatoryAndEmpty?(rec, field, customFields) or
			not ControlValidData?(rec, field)
		}

	ReasonProtected(record, protectField, hwnd, query = false)
		{
		allowDeleteMsg = ''
		if (query isnt false and
			('' isnt allowDeleteMsg = RecordAllowDelete(query, record)))
			allowDeleteMsg = 'This record can not be deleted.\n\n' $ allowDeleteMsg

		protect = protectField is false ? '' : record[protectField]
		if .noInfo(protect) and allowDeleteMsg is ''
			{
			.alert('No Information', hwnd)
			return
			}
		if Object?(protect)
			protect = protect.reason
		.alert(allowDeleteMsg $ Opt('\n\n', protect), hwnd)
		}

	noInfo(protect)
		{
		return (protect is '' or protect is false or protect is true or
			(Object?(protect) and not protect.Member?('reason')) or
			(Object?(protect) and protect.Member?('reason') and protect.reason is ""))
		}

	alert(msg, hwnd)
		{
		Alert(msg, 'Reason Protected', hwnd, MB.ICONINFORMATION)
		}

	MandatoryColumn?(column, mandatoryFields, customFields)
		{
		if mandatoryFields.Has?(column)
			return true

		if Object?(customFields) and customFields.Member?(column) and
			customFields[column].GetDefault('mandatory', false) is true
			return true

		ctrl = Datadict(column).Control
		return ctrl.Member?('mandatory') and ctrl.mandatory is true
		}

	Customize(query, columns, customKey, linkField, ctrl)
		{
		sfOb = Object(cols: columns.Copy().Remove('listrow_deleted'), excludeFields: #())

		table = QueryGetTable(query, nothrow:)
		if false is custom_fields = .InitCustomFields(table, columns)
			return false
		dirty = CustomizeDialog(ctrl.Window.Hwnd, customKey, query, sfOb, true,
			Customizable(table), custom_fields,
			sub_title: .SubTitle(ctrl, linkField))
		return Object(:dirty, :custom_fields)
		}

	GetCustomizedFields(customKey)
		{
		return Customizable.GetCustomizedFields(customKey)
		}

	BuildCustomKeyFromQueryTitle(query, title)
		{
		// if title is false don't bother looking up the table
		if title is false
			return false
		table = QueryGetTable(query, nothrow:)
		return title $ ' ~ ' $ table
		}
	}
