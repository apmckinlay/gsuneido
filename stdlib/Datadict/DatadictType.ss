// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
function (fieldname)
	{
	dd = Datadict(fieldname)
	baseTypes = Object(number: Field_number, string: Field_string,
		date: Field_date, boolean: Field_boolean, image: Field_image)

	for type in baseTypes.Members()
		if dd.Base?(baseTypes[type])
			return type
	return 'string'
	}