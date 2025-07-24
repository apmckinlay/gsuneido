// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase
	{
	Default(@args)
		{
		first = args[1]
		second = args[2]
		return Object(type: FORMULATYPE.STRING,
			value: FormulaConvertToString(first.value) $
				FormulaConvertToString(second.value))
		}
	}