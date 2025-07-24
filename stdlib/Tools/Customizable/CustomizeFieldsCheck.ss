// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	msg = "Checking Customize Fields settings "
	if Customizable.NotCustomizableScreen?()
		return msg $ 'SUCCEEDED \n'

	detail_msg = ""
	QueryApply('customizable_fields')
		{|x|
		if '' isnt valid = x.custfield_valid
			detail_msg $= '\t' $ x.custfield_name $ ' - ' $
				x.custfield_field $ ': ' $ valid $ '\n'

		if x.custfield_formula isnt '' and x.custfield_formula_code is ''
			detail_msg $= '\tFormula code does not match formula text on field ' $
				x.custfield_field $ '\n'

		if x.custfield_formula isnt '' and x.custfield_formula.Blank?()
			detail_msg $= '\tFormula only has white spaces on field ' $
				x.custfield_field $ '\n'
		}

	return detail_msg is ""
		? msg $ 'SUCCEEDED \n'
		: msg $ 'FAILED \n' $ detail_msg
	}