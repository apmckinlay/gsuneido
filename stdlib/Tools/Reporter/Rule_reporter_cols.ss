// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if not Object?(.allcols)
		return #()
	cols = .allcols.Copy()
	for col in .allcols
		{
		field = .selectFields.PromptToField(col)
		if UnsortableField?(field) or
			.nonsummarized_fields.Has?(col) or
			field.Suffix?("_lower!")
			cols.Remove(col)
		}
	return cols
	}
