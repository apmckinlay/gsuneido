// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
KeyControl
	{
	FieldControl: 'IdField'
	Setup_KeyControl()
		{
		}

	GetIdFieldText()
		{
		return GetWindowText(.Field.Hwnd)
		}

	// used by ParamsSelect "in list" option
	DisplayValues(control, vals)
		{
		query = control.GetDefault('query', control.GetDefault(1, false))
		if Function?(query)
			query = query()
		field = control.GetDefault('field', control.GetDefault(2, false))
		nameField = control.GetDefault('nameField',
			field.Replace("(_num|_name|_abbrev)$", "_name"))
		idAllowOther = control.GetDefault('allowOther', false)

		return .FormatValues(vals, query, field, nameField, idAllowOther)
		}

	FormatValues(vals, query, field, nameField, idAllowOther = false)
		{
		displayNames = Object()
		for val in vals
			{
			if false is rec = Query1(QueryAddWhere(query,
				" where " $ field $ " = " $ Display(val)))
				val = idAllowOther ? val : '???'
			else
				val = rec[nameField]
			displayNames.Add(val)
			}
		return displayNames
		}

	KeyIdField_Access()
		{
		if .Dirty?()
			{
			.Field.Process_newvalue()
			.NewValue(.Get())
			}
		super.KeyIdField_Access()
		}
	}
