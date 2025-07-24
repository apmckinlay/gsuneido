// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase
	{
	UnsupportedText: '<op1> and <op2>'
	Calc(firstValue, secondValue)
		{
		return firstValue is true and secondValue is true
		}
	Boolean_Boolean(first, second)
		{
		return .GenerateElement(type: FORMULATYPE.BOOLEAN,
			value: .Calc(first.value, second.value))
		}
	}