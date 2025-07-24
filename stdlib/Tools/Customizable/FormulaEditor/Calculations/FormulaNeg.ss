// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase
	{
	UnsupportedText: "- <op2>"
	CallClass(first)
		{
		return super.CallClass(.GenerateElement(type: FORMULATYPE.NUMBER, value: 0),
			first)
		}

	Calc(unused, secondValue)
		{
		return -secondValue
		}

	Number_Number(first, second)
		{
		return .Calc_number_number(first, second, type: FORMULATYPE.NUMBER)
		}

	Number_UOM(first, second)
		{
		return .Calc_number_uom(first, second, type: FORMULATYPE.UOM)
		}

	Number_UOMRate(first, second)
		{
		return .Calc_number_uom(first, second, type: FORMULATYPE.UOM_RATE)
		}
	}