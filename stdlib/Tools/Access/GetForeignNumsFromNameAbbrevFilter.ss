// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// this is to avoid the query string exceeding the length limit (~1M).
	// users can have multiple conditions all converted
	// 3000 * 20(length of TS) * 16 = ~0.96M
	limit: 3000
	CallClass(field, sf, operation, value = '', value2 = '')
		{
		if operation is ''
			return false
		if false is numField = sf.GetJoinNumField(field)
			return false

		condition = [condition_field: field]
		condition[field] = [:operation, :value, :value2]
		if '' is where = ChooseFiltersControl.BuildWhereFromFilter([condition])
			return false

		query = sf.Joins([field]).Replace(' leftjoin by(.*?) ', '') $ where
		records = QueryAll(query, limit: .limit)
		if records.Size() is .limit
			return false

		nums = records.Map({ it[numField] })
		// if the original condition accept "", the converted condition needs to accept ""
		whereFn = GetParamsWhere(condition.condition_field, data: condition,
			build_callable:)
		if whereFn([])
			nums.Add('')
		return Object(:numField, :nums)
		}
	}
