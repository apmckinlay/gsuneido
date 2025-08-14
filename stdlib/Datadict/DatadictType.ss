// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
function (fieldname)
	{
	dd = Datadict(fieldname)
	baseTypes = Object(number: Field_number, string: Field_string,
		date: Field_date, boolean: Field_boolean, image: Field_image, info: Field_info)

	// sort because info needs to be checked before string (info field is base of both)
	for type in baseTypes.Members().Sort!()
		if dd.Base?(baseTypes[type])
			return type
	return 'string'
	}