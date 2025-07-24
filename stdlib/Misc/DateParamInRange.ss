// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(date_param, start, end)
		{
		if not Date?(start) or not Date?(end)
			return false

		if not Object?(date_param) or date_param.operation is ""
			return false

		val = Field_date.Encode(date_param.value)
		val2 = Field_date.Encode(date_param.value2)

		fn = date_param.operation.Replace(' ', '_').Capitalize()
		if .Method?(fn)
			return this[fn](:date_param, :start, :end, :val, :val2)
		throw "unhandled operation: " $ date_param.operation
		}

	In_list(date_param, start, end)
		{
		return date_param.value.Any?({ it >= start and it <= end })
		}

	Not_in_list(date_param, start, end)
		{
		for (date = start; date <= end; date = date.Plus(days:1))
			if not date_param.value.Has?(date)
				return true
		return false
		}

	Range(start, end, val, val2)
		{
		return ((val <= start and val2 >= start) or (val >= start and val2 <= end) or
			(val <= end and val2 >= end) or (val <= start and val2 >= end))
		}

	Not_in_range(start, end, val, val2)
		{
		return not (start >= val and end <= val2)
		}

	Empty(date_param)
		{
		return date_param.operation.Prefix?('not')
		}

	Not_empty(date_param)
		{
		return .Empty(date_param)
		}

	Equals(start, end, val)
		{
		return val >= start and val <= end
		}

	Not_equal_to(start, end, val)
		{
		return not (start is end and start is val)
		}

	Greater_than(end, val)
		{
		return end > val
		}

	Less_than(start, val)
		{
		return start < val
		}

	Greater_than_or_equal_to(val, end)
		{
		return end >= val
		}

	Less_than_or_equal_to(start, val)
		{
		return start <= val
		}
	}