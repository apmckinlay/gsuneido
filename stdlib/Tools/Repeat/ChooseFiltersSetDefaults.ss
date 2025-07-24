// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(defaultFilterFields, savedFilters)
		{
		if defaultFilterFields is false or defaultFilterFields.Empty?()
			return savedFilters

		if savedFilters is ""
			savedFilters = Object()

		for field in defaultFilterFields
			{
			if false is .repeatConditionExists(savedFilters, field)
				{
				pos = defaultFilterFields.Find(field)
				ob = Record(condition_field: field)
				ob[field] = #(operation: '', value: '', value2: '')
				savedFilters.Add(ob, at: pos)
				}
			}

		return savedFilters
		}

	repeatConditionExists(savedFilters, condition_field)
		{
		return savedFilters.FindIf({ it.condition_field is condition_field })
		}
	}