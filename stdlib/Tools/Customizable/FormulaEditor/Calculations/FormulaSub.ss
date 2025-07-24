// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaAdd
	{
	UnsupportedText: '<op1> - <op2>'
	Calc(firstValue, secondValue)
		{
		return firstValue - secondValue
		}

	Sign: -1
	Date_Date(firstDate, secondDate)
		{
		try
			secs = firstDate.value.MinusSeconds(secondDate.value)
		catch (unused, '*interval too large')
			secs = CompatibleUOMFactor.Calc(
				firstDate.value.MinusDays(secondDate.value), 'days', 'secs')

		return Object(type: FORMULATYPE.UOM, value: secs $ ' secs')
		}

	UOM_Date(@unused)
		{
		.Unsupport('UOM_Date')
		}
	}