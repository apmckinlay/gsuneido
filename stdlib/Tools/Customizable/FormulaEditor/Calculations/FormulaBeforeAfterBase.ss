// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	DisplayName: ''
	CallClass(s, delimiter)
		{
		.Validate(Object(s.type), Object(delimiter.type))

		try
			value = .Calc(s, delimiter)
		catch
			throw "Formula: " $ .DisplayName $ " failed to extract substring " $
				Display(delimiter.value) $ ' from ' $ Display(s.value)

		return Object(type: FORMULATYPE.STRING, :value)
		}

	Calc(s /*unused*/, delimiter /*unused*/)
		{
		throw 'MUST IMPLEMENT'
		}

	Validate(@args)
		{
		if args.Size() < 2
			throw "Formula: " $ .DisplayName $ " missing arguments"

		if args.Size() > 2
			throw "Formula: " $ .DisplayName $ " too many arguments"

		if args[0] isnt Object(FORMULATYPE.STRING)
			throw "Formula: " $ .DisplayName $ " text must be a <String>"

		if args[1] isnt Object(FORMULATYPE.STRING)
			throw "Formula: " $ .DisplayName $ " delimiter must be a <String>"

		return Object(FORMULATYPE.STRING)
		}
	}
