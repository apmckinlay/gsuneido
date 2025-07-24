// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	DisplayName: 'CONVERT'
	EmptyPlaceHolder: false

	CallClass(rec, unit)
		{
		.Validate(Object(rec.type), Object(unit.type))

		uom = Split_UOM(rec.value)
		if uom.uom is '' or uom.value is ''
			FormulaBase.InvalidValue()

		if unit.value is ''
			throw "Formula: " $ .DisplayName $ " unit must not be empty"

		if .EmptyPlaceHolder isnt false and unit.value is .EmptyPlaceHolder
			return Object(type: .ReturnType(rec.type),
				value: .ReturnValue(uom.value, unit.value))

		if false is converted = FormulaBase.ConvertValue(
			uom.value, uom.uom, unit.value, rate?: rec.type is FORMULATYPE.UOM_RATE)
			FormulaBase.IncompatibleUOM(uom.uom, unit.value)

		return Object(type: .ReturnType(rec.type),
			value: .ReturnValue(converted, unit.value))
		}

	Validate(@args)
		{
		if args.Size() < 2
			throw "Formula: " $ .DisplayName $ " missing arguments"

		if args.Size() > 2
			throw "Formula: " $ .DisplayName $ " too many arguments"

		if not args[0].Difference(Object(FORMULATYPE.UOM, FORMULATYPE.UOM_RATE)).Empty?()
			throw "Formula: " $ .DisplayName $ " Field must be a <Quantity> or <Rate>"

		if args.Member?(1) and args[1] isnt Object(FORMULATYPE.STRING)
			throw "Formula: " $ .DisplayName $ " Unit must be a <String>"

		return args[0].Map(.ReturnType).Unique!()
		}

	ReturnType(recType)
		{
		return recType
		}
	ReturnValue(converted, unitStr)
		{
		return converted $ ' ' $ unitStr
		}
	}
