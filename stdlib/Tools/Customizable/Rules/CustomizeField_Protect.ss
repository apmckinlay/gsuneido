// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(rec)
		{
		protect = Object()
		dict = Datadict(rec.custfield_field)
		.protectFromControlDef(dict, protect, rec)
		.protectFormulaAndDefaultValue(rec, protect, dict)
		.handleUOM(dict, protect)

		if Customizable.CustomField?(rec.custfield_field) and
			DatadictType(rec.custfield_field) isnt 'image'
			protect.Delete('custfield_only_fillin_from')

		customizable = dict.Member?('AllowCustomizableOptions')
			? dict.AllowCustomizableOptions : #()
		for member in customizable.Members()
			{
			if customizable[member] is true
				protect.Delete('custfield_' $ member)
			else if customizable[member] is false
				protect['custfield_' $ member] = true
			}

		protect.allowDelete = true
		return protect
		}

	protectFormulaAndDefaultValue(rec, protect, dict)
		{
		try
			{
			Global('Rule_' $ rec.custfield_field)
			protect.custfield_formula = true
			}

		try
			{
			Global('Rule_' $ rec.custfield_field $ '__protect')
			protect.custfield_formula = true
			protect.custfield_default_value = true
			}

		if Suneido.User isnt 'default' and Suneido.User isnt 'axon'
			protect.custfield_formula = true

		if (dict.Member?('NoCustomDefaultValue') and dict.NoCustomDefaultValue is true)
			protect.custfield_default_value = true
		}

	handleUOM(dict, protect)
		{
		if dict.Control[0] is 'UOM' and
			((Object?(dict.Control[1]) and
				dict.Control[1].GetDefault('mandatory', false)) or
			(Object?(dict.Control[2]) and
				dict.Control[2].GetDefault('mandatory', false)))
			{
			protect.custfield_mandatory = true
			protect.custfield_readonly = true
			}
		}

	protectFromControlDef(dict, protect, rec)
		{
		ctrlClass = GetControlClass.FromControl(dict.Control)
		params = ctrlClass.Member?('CustomizableOptions')
			? ctrlClass.CustomizableOptions.Join(',')
			: ctrlClass.New.Params()

		for opt in CustomFieldOptions()
			if not (params =~  "\<" $ opt.field_option $ "\>")
				protect[opt.field] = true

		.handleTabover(dict, protect, rec)

		.handleFocus(dict, protect, rec)

		if dict.Control.GetDefault('mandatory', false) is true
			{
			protect.custfield_mandatory = true
			protect.custfield_readonly = true
			protect.custfield_hidden = true
			}

		if dict.Control.GetDefault('readonly', false) is true
			{
			protect.custfield_mandatory = true
			protect.custfield_readonly = true
			protect.custfield_formula = true
			}
		}

	handleTabover(dict, protect, rec)
		{
		if dict.Control.GetDefault('tabover', false) is true
			protect.custfield_tabover = true
		else if rec.custfield_browse? is true and protect.Member?('custfield_tabover')
			{
			// always allow tabover for fields from browse
			protect.Delete('custfield_tabover')
			}
		}

	handleFocus(dict, protect, rec)
		{
		// always allow first focus for fields from access
		if dict.Control.GetDefault('first_focus', false) is true
			protect.custfield_first_focus = true
		else if rec.custfield_browse? is false and
			protect.Member?('custfield_first_focus')
			protect.Delete('custfield_first_focus')
		}
	}