// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaAnd
	{
	UnsupportedText: "not <op1>"
	CallClass(first)
		{
		return super.CallClass(first,
			.GenerateElement(type: FORMULATYPE.BOOLEAN, value: true))
		}

	Calc(firstValue, unused)
		{
		return firstValue isnt true
		}
	}