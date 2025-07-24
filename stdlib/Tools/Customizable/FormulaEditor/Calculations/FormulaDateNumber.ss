// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	DisplayName: false
	CallClass(date)
		{
		.Validate(Object(date.type))

		if not Date?(date.value)
			FormulaBase.InvalidValue()

		return Object(type: FORMULATYPE.NUMBER, value: .Calc(date.value))
		}

	Calc(dateValue /*unused*/)
		{
		throw 'MUST IMPLEMENT'
		}

	Validate(@args)
		{
		if args.Size() < 1
			throw "Formula: " $ .DisplayName $ " missing arguments"

		if args.Size() > 1
			throw "Formula: " $ .DisplayName $ " too many arguments"

		if not args[0].Difference(Object(FORMULATYPE.DATE)).Empty?()
			throw "Formula: " $ .DisplayName $ " Field must be a <Date>"

		return Object(FORMULATYPE.NUMBER)
		}
	}