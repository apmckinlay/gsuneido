// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(rec)
		{
		if not Object?(rec.design_cols)
			return #()
		cols = rec.design_cols.Copy()
		cols.RemoveIf({|col| .unsortable?(rec, col) })
		.move_designcols_to_start(rec, cols)
		return cols
		}

	unsortable?(rec, col)
		{
		return UnsortableField?(rec.selectFields.PromptToField(col))
		}

	move_designcols_to_start(rec, cols)
		{
		cols.Add('')
		rec_columns = rec.columns.Map({ it.text })
		cols.SortWith!({|col| rec_columns.Has?(col) ? 0 : col is '' ? 1 : 2 })
		cols.Trim!('')
		}
	}