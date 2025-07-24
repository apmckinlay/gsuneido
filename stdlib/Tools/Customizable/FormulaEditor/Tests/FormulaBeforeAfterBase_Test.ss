// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_Validate()
		{
		fn = FormulaBeforeAfterBase.Validate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING

		Assert(fn([s], [s]) is: [s])

		Assert({ fn([s]) } throws: "missing arguments")
		Assert({ fn([s], [s], [s]) } throws: "too many arguments")
		Assert({ fn([s, n], [s]) }
			throws: "text must be a <String>")
		Assert({ fn([s], [s, n]) }
			throws: "delimiter must be a <String>")
		}
	}