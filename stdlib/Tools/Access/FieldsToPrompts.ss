// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (fields, map)
	// pre: map should be an empty object
	// post: returns list of fields
	// post: map has prompt members, field values
	{
	return fields.Map(
		{ |field|
		prompt = SelectPrompt(field).Tr(',').Trim()
		map[prompt] = field
		prompt
		}).Sort!()
	}