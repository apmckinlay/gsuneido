// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cases:
		(
		(from, to, true, ".from", ".to")
		(from, to, true, "from", "from")
		(from, to, true, "this.from", "this.to")
		(from, to, true, "x.from", "x.from")
		(from, to, true, "x().from", "x().from")
		(from, to, true, "x[0].from", "x[0].from")
		(from, to, true, "from:", "to:")
		(from, to, true, "{ from: }", "{ from: }")
		(from, to, false, "/* from */ .from", "/* from */ .to")
		(from, to, true, "/* from */ .from",	"/* to */ .to")
		)
	Test_main()
		{
		for c in .cases
			{
			text = "class {\n" $ c[3] $ "\n}"
			result = "class {\n" $ c[4] $ "\n}"
			Assert((new Refactor_Rename_Member).Rename(text, c[0], c[1], c[2]) is: result)
			}
		}
	}