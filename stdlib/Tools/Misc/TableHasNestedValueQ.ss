// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
FindValueInX
	{
	CallClass(query, fieldsToCheck, value)
		{
		if String?(fieldsToCheck)
			return .queryCheckOneField(query, fieldsToCheck, value)

		for field in fieldsToCheck
			if .queryCheckOneField(query, field, value)
				return true
		return false
		}

	queryCheckOneField(query, fieldToCheck, value)
		{
		QueryApply(query)
			{
			targetOb = it[fieldToCheck]
			if .Reference?(targetOb, value)
				return true
			}
		return false
		}
	}