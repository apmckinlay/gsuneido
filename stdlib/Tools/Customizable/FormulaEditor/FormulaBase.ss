// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Validate?: true
	CallClass(first, second) // formulafn
		{
		if .Validate?
			{
			.validate(first)
			.validate(second)
			}

		firstType = .typeMapping(first)
		secondType = .typeMapping(second)
		return this[firstType $ '_' $ secondType](first, second)
		}

	validate(elm)
		{
		valid? = false
		switch (elm.type)
			{
		case FORMULATYPE.NUMBER:
			valid? = .validateNumber(elm.value)
		case FORMULATYPE.STRING:
			valid? = String?(elm.value)
		case FORMULATYPE.UOM, FORMULATYPE.UOM_RATE:
			valid? = .validateUOM(elm.value)
		case FORMULATYPE.DATE:
			valid? = Date?(elm.value)
		case FORMULATYPE.BOOLEAN:
			valid? = .validateBoolean(elm.value)
		default:
			}
		if valid? is false
			.InvalidValue()
		}

	validateNumber(value)
		{
		return Number?(value) or value is ''
		}

	validateUOM(value)
		{
		uom = Split_UOM(value)
		return Number?(uom.value) and uom.uom isnt ''
		}

	validateBoolean(value)
		{
		return Boolean?(value) or value is ''
		}

	typeMapping(rec)
		{
		switch (rec.type)
			{
		case FORMULATYPE.NUMBER:
			return #Number
		case FORMULATYPE.STRING:
			return #String
		case FORMULATYPE.UOM:
			return #UOM
		case FORMULATYPE.UOM_RATE:
			return #UOMRate
		case FORMULATYPE.DATE:
			return #Date
		case FORMULATYPE.BOOLEAN:
			return #Boolean
			}
		}

	Default(@args)
		{
		.Unsupport(args[0])
		}

	Unsupport(operation)
		{
		types = operation.Split('_').Map!(.typeToDisplay)
		msg = "Formula: Operation not supported: " $
			Display(.UnsupportedText.Replace('op1', types[0]).Replace('op2', types[1]))
		throw msg
		}

	typeToDisplay(type)
		{
		switch (type)
			{
		case #UOM:
			return 'Quantity'
		case #UOMRate:
			return 'Rate'
		default:
			return type
			}
		}

	Calc_number_uom(first, second, type, rev = false)
		{
		uom = Split_UOM(second.value)
		number = first.value

		res = rev is false ? .Calc(number, uom.value) : .Calc(uom.value, number)
		return .GenerateElement(:type, value: .generateValue(res, type, uom.uom))
		}

	Calc_uom_uom(first, second, type)
		{
		firstUOM = Split_UOM(first.value)
		secondUOM = Split_UOM(second.value)

		convertedSecondValue = .ConvertValue(secondUOM.value, secondUOM.uom,
			firstUOM.uom, rate?: second.type is FORMULATYPE.UOM_RATE)

		if convertedSecondValue is false
			.IncompatibleUOM(firstUOM.uom, secondUOM.uom)

		res = .Calc(firstUOM.value, convertedSecondValue)
		return .GenerateElement(:type,
			value: .generateValue(Number?(res) ? res.Round(2): res, type, firstUOM.uom))
		}

	ConvertValue(fromValue, fromUom, toUom, rate?)
		{
		convertedSecondValue = CompatibleUOMFactor.Calc(fromValue, fromUom, toUom, rate?)

		if convertedSecondValue isnt false
			return convertedSecondValue

		if fromUom.Suffix?('s') or fromUom.Suffix?('S')
			return .ConvertValue(fromValue, fromUom[..-1], toUom, rate?)

		if toUom.Suffix?('s') or toUom.Suffix?('S')
			return .ConvertValue(fromValue, fromUom, toUom[..-1], rate?)

		return false
		}

	generateValue(n, type, unit = false)
		{
		return type in (FORMULATYPE.UOM, FORMULATYPE.UOM_RATE)
			? n $ ' ' $ unit
			: n
		}

	Calc_number_number(first, second, type)
		{
		return .GenerateElement(:type, value: .Calc(first.value, second.value))
		}

	GenerateElement(type, value)
		{
		return Object(:type, :value)
		}

	IncompatibleUOM(unit1 = '', unit2 = '')
		{
		throw 'Formula: Incompatible unit of measure' $ Opt('(', unit1, ', ', unit2, ')')
		}

	InvalidValue()
		{
		throw 'Formula: Invalid Value'
		}
	}
