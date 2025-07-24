// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(custom)
		{
		default_value = custom.custfield_default_value
		dd = Datadict(custom.custfield_field)
		ctrl = dd.Control
		if ctrl[0] is 'Id' or ctrl[0] is 'IdControl'
			{
			field = ctrl.Member?('field') ? ctrl.field : ctrl[2]
			nameField = ctrl.GetDefault(#nameField,
				field.Replace("(_num|_name|_abbrev)$", "_name"))
			query = KeyControl.Key_BuildQuery(ctrl[1],
				ctrl.GetDefault('restrictions', false),
				ctrl.GetDefault('invalidRestrictions', false),
				noSend:)
			where = ' where ' $ field $ ' is ' $ Display(custom.custfield_default_value)
			if false isnt rec = Query1(QueryAddWhere(query, where))
				default_value = rec[nameField]
			}
		else if dd.Format[0].Has?('Date') and Date?(default_value)
			default_value = default_value.NoTime() is default_value
				? default_value.StdShortDate()
				: default_value.StdShortDateTime()


		return .displayExtra(custom, default_value)
		}

	displayExtra(custom, default_value)
		{
		ctrl = GetControlClass.FromField(custom.custfield_field)
		if ctrl.Method?('Customizable_custfield_options_extra')
			{
			if '' isnt extraVal = ctrl.Customizable_custfield_options_extra(custom)
				return extraVal
			}
		return default_value
		}
	}