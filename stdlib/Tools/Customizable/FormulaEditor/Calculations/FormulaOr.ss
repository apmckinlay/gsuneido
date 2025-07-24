// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaAnd
	{
	UnsupportedText: '<op1> or <op2>'
	Calc(firstValue, secondValue)
		{
		return firstValue is true or secondValue is true
		}
	}