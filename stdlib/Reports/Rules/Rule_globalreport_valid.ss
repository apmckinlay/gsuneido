// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	check_fields = function (fields, sf, type = "sort")
		{
		if fields.Blank?()
			return ""
		selected_fields = fields.Split(',')
		if type is "sort" and selected_fields.Size() > 4
			return "Can not sort by more than 4 fields"
		for field in selected_fields
			{
			field = field.Trim()
			if not sf.HasPrompt?(field)
				return type is "sort"
					? "Invalid Sort By field: " $ field
					: "Invalid Column: " $ field
			field = sf.PromptToField(field)
			if type is "sort" and Datadict(field).Base?(Field_image)
				return "Can not sort by image fields"
			}
		return ""
		}
	if ("" isnt (result = check_fields(.fields, .sf, "fields")))
		return result
	return check_fields(.sort, .sf)
	}