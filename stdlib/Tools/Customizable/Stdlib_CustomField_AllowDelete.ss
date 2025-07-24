// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
// Tested in a different library's contribution
function (unused, field)
	{
	msg = ''
	if not QueryEmpty?('customizable_fields
		where (custfield_field is "' $ field $ '" or
			custfield_formula_code.Has?("' $ field $ '") or
			custfield_formula_fields.Has?("' $ field $ '"))')
		msg  $= 'Axon:\tFormulas\r\n'

	if TableHasNestedValue?('event_condition_actions', 'eca_conditions', field)
		msg $= 'Axon:\tBusiness Triggers\r\n'

	return msg
	}