// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase
	{
	UnsupportedText: '<op1> + <op2>'
	Calc(firstValue, secondValue)
		{
		return firstValue + secondValue
		}

	Number_Number(first, second)
		{
		return .Calc_number_number(first, second, FORMULATYPE.NUMBER)
		}

	UOM_UOM(first, second)
		{
		return .Calc_uom_uom(first, second, FORMULATYPE.UOM)
		}

	UOMRate_UOMRate(first, second)
		{
		return .Calc_uom_uom(first, second, FORMULATYPE.UOM_RATE)
		}

	Sign: 1
	Date_UOM(first, second)
		{
		uom = Split_UOM(second.value)
		value = uom.uom is CustomizeField.FormulaTestUnit
			? first.value
			: .datePlus(first.value, .Sign * uom.value, uom.uom)

		return .GenerateElement(type: FORMULATYPE.DATE, :value)
		}

	datePlus(date, offset, unit)
		{
		switch (unit.Lower())
			{
		case 'second', 'seconds', 'secs':
			return date.Plus(seconds: offset)
		case 'minute', 'minutes', 'mins':
			return date.Plus(minutes: offset)
		case 'hour', 'hours', 'hrs':
			return date.Plus(hours: offset)
		case 'day', 'days':
			return date.Plus(days: offset)
		case 'month', 'months':
			return date.Plus(months: offset)
		case 'year', 'years', 'yrs':
			return date.Plus(years: offset)
		default:
			.IncompatibleUOM('Date', unit)
			}
		}

	UOM_Date(first, second)
		{
		return .Date_UOM(second, first)
		}
	}


