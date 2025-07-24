// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if '' isnt msg = .eca_action_valid
		return msg

	for c in .eca_conditions
		{
		if false is op = Select2.Ops.FindOne({ it[0] is c.condition_op })
			continue
			
		if Select2.Invalid_operator?(op, Datadict(c.condition_field))
			return Display(c.condition_op) $ ' is not valid for ' $
				Prompt(c.condition_field)

		if Select2.InvalidOpValue?(op, c.condition_value)
			return Display(c.condition_value) $ ' is not valid for ' $
				Display(c.condition_op)
		}

	return ''
	}