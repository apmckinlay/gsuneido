// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase
	{
	UnsupportedText: '<op1> * <op2>'
	Calc(firstValue, secondValue)
		{
		return firstValue * secondValue
		}

	Number_UOM(first, second)
		{
		return .Calc_number_uom(first, second, FORMULATYPE.UOM)
		}

	Number_UOMRate(first, second)
		{
		return .Calc_number_uom(first, second, FORMULATYPE.UOM_RATE)
		}

	Number_Number(first, second)
		{
		return .Calc_number_number(first, second, FORMULATYPE.NUMBER)
		}

	UOM_Number(first, second)
		{
		return .Calc_number_uom(second, first, FORMULATYPE.UOM, rev:)
		}

	UOM_UOMRate(first, second)
		{
		return .Calc_uom_uom(second, first, FORMULATYPE.NUMBER)
		}

	UOMRate_Number(first, second)
		{
		return .Calc_number_uom(second, first, FORMULATYPE.UOM_RATE, rev:)
		}

	UOMRate_UOM(first, second)
		{
		return .Calc_uom_uom(first, second, FORMULATYPE.NUMBER)
		}
	}
