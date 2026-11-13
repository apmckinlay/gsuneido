// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CopyFieldName(fieldName)
		{
		if fieldName isnt false
			WndProc.Copy_Field_Name(fieldName)
		}

	GoToFieldDefinition(fieldName)
		{
		if fieldName isnt false
			GotoLibView("Field_" $ fieldName)
		}

	GoToQueryView(query, rec = false, keys = #())
		{
		if rec isnt false and not keys.Empty?()
			query = .keyQuery(query, keys, rec)
		GotoQueryView(query)
		}

	keyQuery(query, keys, rec)
		{
		where = keys.Map({ it $ " = " $ Display(rec[it]) }).Join(" and ")
		return QueryHelper.StripSort(query) $ Opt(" where ", where)
		}
	}
