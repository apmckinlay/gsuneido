// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaConvert
	{
	DisplayName: 'TONUMBER'
	EmptyPlaceHolder:
		'This is a temporary place holder. Please replace it with a proper UOM'

	ReturnType(unused)
		{
		return FORMULATYPE.NUMBER
		}
	ReturnValue(converted, unused)
		{
		return converted
		}

	CheckEmptyPlaceHolder(code)
		{
		if code.Has?(.EmptyPlaceHolder)
			return 'There is a temporary place holder in Formulas. ' $
				'Please replace it with a proper UOM'
		return ""
		}
	}
