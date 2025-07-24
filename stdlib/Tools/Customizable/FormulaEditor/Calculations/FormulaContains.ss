// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Calc(text, substring)
		{
		return  text.value =~ "(?i)(?q)" $ substring.value
		}

	CallClass(text, substring)
		{
		.Validate(Object(text.type), Object(substring.type))

		try
			value = .Calc(text, substring)
		catch
			throw "Formula: " $ .DisplayName $ " failed to search " $
				Display(substring.value) $ ' from ' $ Display(text.value)

		return FormulaBase.GenerateElement(type: FORMULATYPE.BOOLEAN, :value)
		}

	DisplayName: 'CONTAINS'
	Validate(@args)
		{
		if args.Size() < 2
			throw "Formula: " $ .DisplayName $ " missing arguments"

		if args.Size() > 2
			throw "Formula: " $ .DisplayName $ " too many arguments"

		typesAllowed = Object(FORMULATYPE.STRING, FORMULATYPE.UOM, FORMULATYPE.UOM_RATE)

		if not args[0].Difference(typesAllowed).Empty?()
			throw "Formula: " $ .DisplayName $
				" Text must be a <String>, <Quantity> or <Rate>"

		if not args[1].Difference(typesAllowed).Empty?()
			throw "Formula: " $ .DisplayName $
				" Substring must be a <String>, <Quantity> or <Rate>"

		return Object(FORMULATYPE.BOOLEAN)
		}
	}