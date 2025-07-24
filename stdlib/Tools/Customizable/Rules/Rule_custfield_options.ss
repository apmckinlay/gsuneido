// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	ob = Object()
	for opt in CustomFieldOptions()
		if this[opt.field] is true and not .custfield_protect.Member?(opt.field)
			ob.Add(Prompt(opt.field))
	str = ob.Join(', ')

	if .custfield_default_value isnt '' and .custfield_default_value isnt false
		{
		default_value = CustomizeField_DisplayDefaultValue(this)
		str = Opt(str, ', ') $ 'Default=' $ Display(default_value)
		}

	if .custfield_only_fillin_from isnt ''
		str = Opt(str, ', ') $ 'Only fills in from ' $
			SelectPrompt(.custfield_only_fillin_from)

	if .custfield_formula isnt ''
		str = Opt(str, ', ') $ 'Formula=' $ Display(.custfield_formula)
	return str
	}
