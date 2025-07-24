// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase
	{
	Validate?: false
	UnsupportedText: '<op1> > <op2>'
	Calc(firstValue, secondValue)
		{
		return firstValue > secondValue
		}

	UOM_UOM(first, second)
		{
		if false isnt res = .handleEmptyUOM(first, second)
			return res
		return .Calc_uom_uom(first, second, FORMULATYPE.BOOLEAN)
		}

	UOMRate_UOMRate(first, second)
		{
		if false isnt res = .handleEmptyUOM(first, second)
			return res
		return .Calc_uom_uom(first, second, FORMULATYPE.BOOLEAN)
		}

	handleEmptyUOM(first, second)
		{
		return first.value is '' or second.value is ''
			?.GenerateElement(type: FORMULATYPE.BOOLEAN,
				value: .Calc(first.value, second.value))
			: false
		}

	Default(func, first, second)
		{
		return (first.type is second.type or first.value is "" or second.value is "")
			? .GenerateElement(type: FORMULATYPE.BOOLEAN,
				value: .Calc(first.value, second.value))
			: super.Default(func, first, second)
		}
	}