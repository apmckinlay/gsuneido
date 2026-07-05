// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(query, field, maxRecords)
		{
		// pre: maxRecords must be greater than 0
		Assert(maxRecords > 0)
		value = ""
		count = 0
		QueryApply(query $ " sort reverse " $ field)
			{ |x|
			++count
			value = x[field]
			if count is maxRecords
				break
			}
		return Object(indexWhere: " where " $ field $ " >= " $ Display(value),
			startIndexValue: value
			maxRecords: count is maxRecords)
		}
	}
