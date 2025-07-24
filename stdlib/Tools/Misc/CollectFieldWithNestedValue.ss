// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
FindValueInX
	{
	CallClass(query, fieldsToCheck, value, collectFieldName)
		{
		if String?(fieldsToCheck)
			return .queryCheckOneField(query, fieldsToCheck, value, collectFieldName)

		ob = Object()
		for field in fieldsToCheck
			{
			result = .queryCheckOneField(query, field, value, collectFieldName)
			if not result.Empty?()
				ob.Add(result)
			}

		return ob.Flatten()
		}

	queryCheckOneField(query, fieldToCheck, value, collectFieldName)
		{
		ob = Object()
		QueryApply(query)
			{
			targetOb = it[fieldToCheck]
			if .Reference?(targetOb, value)
				ob.Add(it[collectFieldName])
			}
		return ob
		}
	}