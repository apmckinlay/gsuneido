// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	fields = Object()
	for field in .custfield_fields_list
		{
		if KeyControl.IsKeyControl?(field)
			fields.Add(field)
		}
	return fields
	}