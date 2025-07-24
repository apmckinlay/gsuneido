// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Type()
		{
		return FORMULATYPE.UOM
		}

	CallClass(value, unit)
		{
		.Validate(Object(value.type), Object(unit.type))
		return Object(type: .Type(), value: Opt(value.value, ' ', unit.value))
		}

	DisplayName: 'QUANTITY'
	Validate(@args)
		{
		if args.Size() < 2
			.throwError("missing arguments")

		if args.Size() > 2
			.throwError("too many arguments")

		if args[0] isnt Object(FORMULATYPE.NUMBER)
			.throwError("Value must be a <Number>")

		if args[1] isnt Object(FORMULATYPE.STRING)
			.throwError("Unit must be a <string>")

		return Object(.Type())
		}

	throwError(err)
		{
		throw "Formula: " $ .DisplayName $ " " $ err
		}
	}