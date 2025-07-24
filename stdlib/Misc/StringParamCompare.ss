// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(param, curVal)
		{
		if false is Object?(param) or param.operation is ""
			return true

		if not String?(curVal)
			curVal = String(curVal)
		val = Field_string.Encode(param.value)
		val2 = Field_string.Encode(param.value2)
		return .compareVal(param, val, val2, curVal)
		}

	compareVal(param, val, val2, curVal)
		{
		switch (param.operation)
			{
		case "in list":
			return param.value.Has?(curVal)
		case "not in list":
			return not param.value.Has?(curVal)
		case "range":
			return curVal >= val and curVal <= val2
		case "not in range":
			return curVal < val or curVal > val2
		case "empty":
			return curVal is ""
		case "not empty":
			return curVal isnt ""
		case "equals":
			return val is curVal
		case "not equal to":
			return val isnt curVal
		case "greater than":
			return curVal > val
		case "greater than or equal to":
			return curVal >= val
		case "less than":
			return curVal < val
		case "less than or equal to":
			return curVal <= val
		case "contains":
			return curVal.Lower() =~ val.Lower()
		case "does not contain":
			return curVal.Lower() !~ val.Lower()
		case "starts with":
			return curVal.Lower().Prefix?(val.Lower())
		case "ends with":
			return curVal.Lower().Suffix?(val.Lower())
		case "matches":
			return curVal =~ val
		case "does not match":
			return curVal !~ val
			}
		}
	}