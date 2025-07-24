// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
class
	{
	DisplayName: 'UPPER'
	CallClass(s)
		{
		.Validate(Object(s.type))
		try
			value = s.value.Upper()
		catch
			throw "Formula: " $ .DisplayName $ " failed to convert to upper case"
		return Object(type: FORMULATYPE.STRING, :value)
		}

	Validate(@args)
		{
		if args.Size() < 1
			throw "Formula: " $ .DisplayName $ " missing arguments"

		if args.Size() > 1
			throw "Formula: " $ .DisplayName $ " too many arguments"

		if args[0] isnt Object(FORMULATYPE.STRING)
			throw "Formula: " $ .DisplayName $ " text must be a <String>"

		return Object(FORMULATYPE.STRING)
		}
	}
