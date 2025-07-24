// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaUpper
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		.Check(fn, [s, ''], [s, ''])
		.Check(fn, [s, 'abcDEF'], [s, 'ABCDEF'])
		.Check(fn, [s, '123ABCdefghi'], [s, '123ABCDEFGHI'])

		.CheckError(fn, [b, true], "Formula: UPPER text must be a <String>")
		.CheckError(fn, [n, true], "Formula: UPPER text must be a <String>")
		.CheckError(fn, [u, true], "Formula: UPPER text must be a <String>")
		.CheckError(fn, [r, true], "Formula: UPPER text must be a <String>")
		}

	Test_Validate()
		{
		fn = FormulaUpper.Validate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		Assert(fn([s]) is: [s])
		Assert({ fn() } throws: "UPPER missing arguments")
		Assert({ fn([s], [s]) } throws: "UPPER too many arguments")
		Assert({ fn([n]) } throws: "UPPER text must be a <String>")
		Assert({ fn([b]) } throws: "UPPER text must be a <String>")
		Assert({ fn([u]) } throws: "UPPER text must be a <String>")
		Assert({ fn([r]) } throws: "UPPER text must be a <String>")
		}
	}