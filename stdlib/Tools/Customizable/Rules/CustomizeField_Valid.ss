// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(rec)
		{
		if rec.custfield_mandatory is true and
			(rec.custfield_readonly is true or rec.custfield_hidden is true)
			return "can not make mandatory field read-only/hidden."

		if rec.custfield_hidden is true and rec.custfield_default_value isnt ""
			return "hidden field can not have default value."

		if '' isnt invalid = .checkCustField?(rec)
			return invalid

		if .customField?(rec)
			return 'Only Fill-in From is not allowed for this field'

		if '' isnt msg = .check_additional(rec)
			return msg

		return .checkFormula(rec)
		}

	checkCustField?(rec)
		{
		if rec.custfield_field is '' or rec.custfield_field is false
			return ''

		field_def = Datadict(rec.custfield_field)
		if '' isnt invalid = .checkMandatory(rec, field_def)
			return invalid

		if .noCustomDefaultValue?(rec, field_def)
			return "default value is not allowed for this field"

		return .checkConfigLib(rec, field_def)
		}

	checkMandatory(rec, field_def)
		{
		if field_def.Control[0] is 'CheckBox' and rec.custfield_mandatory is true
			return 'mandatory is not allowed for check box.'

		if field_def.Control.GetDefault('mandatory', false) is true and
			(rec.custfield_readonly is true or rec.custfield_hidden is true)
			return "can not make mandatory field read-only/hidden."

		return ''
		}

	noCustomDefaultValue?(rec, field_def)
		{
		return field_def.Member?('NoCustomDefaultValue') and
			field_def.NoCustomDefaultValue is true and rec.custfield_default_value isnt ""
		}

	checkConfigLib(rec, field_def)
		{
		if Libraries().Has?('configlib') and
			false isnt configRec =
				Query1('configlib', name: 'Field_' $ rec.custfield_field)
			{
			code = configRec.text.RemovePrefix('_')
			field_def_lib = code.Compile()
			if .different?(field_def, field_def_lib, 'Prompt') or
				.different?(field_def, field_def_lib, 'Custom')
				return "Another user has renamed " $ field_def.Prompt $ " to " $
					field_def_lib.Prompt
			}
		return ''
		}

	different?(field_def1, field_def2, option)
		{
		opt1 = field_def1.Member?(option) ? field_def1[option] : false
		opt2 = field_def2.Member?(option) ? field_def2[option] : false
		return opt1 isnt opt2
		}

	customField?(rec)
		{
		return rec.custfield_field isnt '' and
			not Customizable.CustomField?(rec.custfield_field) and
			rec.custfield_only_fillin_from isnt ''
		}

	check_additional(rec)
		{
		valid = OptContribution('CustomizeField_Valid', function (unused) { return '' })
		return valid(rec)
		}

	checkFormula(rec)
		{
		if .hasAssignment?(rec)
			return "Can not assign value to field in formula"

		valid = CustomizeField.ValidateCode(rec.custfield_formula_code)
		if not valid.Blank?()
			return valid

		return ''
		}

	hasAssignment?(rec)
		{
		return Object?(rec.custfield_fields_list) and rec.custfield_fields_list.HasIf?(
			{ |field|
			rec.custfield_formula.Split('\n').Any?(
				{ |code_line| code_line.Trim() =~ '^' $ SelectPrompt(field) $ '\s=' })
			})
		}
	}