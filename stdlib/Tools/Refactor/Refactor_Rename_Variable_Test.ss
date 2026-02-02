// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cases:
		(
		(from, to, true, "from", "to")
		(from, to, true, ".from", ".from")
		(from, to, true, "from:", "from:")
		(from, to, false, "/* from */ from", "/* from */ to")
		(from, to, true, "/* from */ from",	"/* to */ to")
		)
	Test_main()
		{
		for c in .cases
			{
			text = "function () {\n" $ c[3] $ "\n}"
			result = "function () {\n" $ c[4] $ "\n}"
			Assert(Refactor_Rename_Variable.Rename(text, c[0], c[1], c[2]) is: result)
			}
		}
	Test_ToExists?()
		{
		te = Refactor_Rename_Variable.ToExists?
		src = 'function (a, .b, _c, ._d) { e .f f: #f }'
		Assert(te(src, 'a'))
//		Assert(te(src, 'b'))
//		Assert(te(src, 'c'))
//		Assert(te(src, 'd'))
		Assert(te(src, 'e'))
		Assert(te(src, 'f') is: false)
		}
	}