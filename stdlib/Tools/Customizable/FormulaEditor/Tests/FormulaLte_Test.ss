// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_main()
		{
		fn = FormulaLte
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE
		b = FORMULATYPE.BOOLEAN
		d = FORMULATYPE.DATE

//		.Check(fn, [n, ''], [n, 1], [b, false])
//		.Check(fn, [n, ''], [n, -1], [b, false])
//		.Check(fn, [b, ''], [b, true], [b, false])
//		.Check(fn, [b, ''], [b, false], [b, false])

		.Check(fn, [n, 1], [n, 0], [b, false])
		.Check(fn, [n, 1], [n, 1], [b, true])
		.Check(fn, [n, 0], [n, 1], [b, true])

		.Check(fn, [s, "b"], [s, "a"], [b, false])
		.Check(fn, [s, "a"], [s, "a"], [b, true])
		.Check(fn, [s, "a"], [s, "b"], [b, true])

		.Check(fn, [u, "2000 g"], [u, "1 kg"], [b, false])
		.Check(fn, [u, "1000 g"], [u, "1 kg"], [b, true])
		.Check(fn, [u, "1 kg"], [u, "2000 g"], [b, true])

		.Check(fn, [r, "2 g"], [r, "1000 kg"], [b, false])
		.Check(fn, [r, "1 g"], [r, "1000 kg"], [b, true])
		.Check(fn, [r, "1000 kg"], [r, "2 g"], [b, true])

		.Check(fn, [b, true], [b, false], [b, false])
		.Check(fn, [b, false], [b, false], [b, true])
		.Check(fn, [b, false], [b, true], [b, true])

//		.Check(fn, [n, 1], [s, ""], [b, true])
//		.Check(fn, [n, 1], [n, ""], [b, true])
		.Check(fn, [n, ""], [n, ""], [b, true])
//		.Check(fn, [d, #20180101], [s, ""], [b, false])
//		.Check(fn, [d, #20180101], [d, ""], [b, false])
		.Check(fn, [d, ""], [d, ""], [b, true])
		.Check(fn, [u, ""], [u, ""], [b, true])
		.Check(fn, [u, ""], [u, "10 kg"], [b, true])
		.Check(fn, [u, "10 kg"], [u, ""], [b, false])

		.CheckError(fn, [n, 0], [s, "1"], "Operation not supported")
		.CheckError(fn, [u, "1 kg"], [r, "2 kg"], "Operation not supported")
		.CheckError(fn, [u, "1 kg"], [u, "2 km"], "Incompatible unit of measure")
		}
	}