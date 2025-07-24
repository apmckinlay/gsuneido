// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	a: #(1, 2, 3,    1, 2, 2, 1)
	b: #(      3, 2, 1, 2,    1, 3)
	Test_Diff()
		{
		edits = #((7 I 3) (5 D 2) (3 I 2) (1 D 2) (0 D 1))
		Assert(Diff(.a, .b) is: edits)
		}
	Test_Apply()
		{
		Assert(Diff.Apply(.a, Diff(.a, .b)) is: .b)
		}
	Test_SideBySide()
		{
		Assert(Diff.SideBySide(.a, .b)
			is: #(
				(1,  "<", ""),
				(2,  "<", ""),
				(3,  "",  3),
				("", ">", 2),
				(1,  "",  1),
				(2,  "",  2),
				(2,  "<", ""),
				(1,  "",  1),
				("", ">", 3)))

		a = #(1, 2, 3, 4)
		b = #(1, 22, 33, 4)
		Assert(Diff.SideBySide(a, b)
			is: #(
				(1, "",  1),
				(2, "#", 22),
				(3, "#", 33),
				(4, "", 4)))
		}
	Test_Three()
		{
		base = #(2, 4, 6, 8)
		a = #(2, 3, 4, 6, 8)
		b = #(2, 4, 8)
		Assert(Diff.Three(base, a, b)
			is: #(
				(2,  "",   2),
				("", "+1", 3),
				(4,  "",   4),
				(6,  "-2", ""),
				(8,  "",   8)))

		base = #(2, 4, 6, 8)
		a = #(2, 3, 4, 8)
		b = #(2, 3, 4, 8)
		Assert(Diff.Three(base, a, b)
			is: #(
				(2, "", 2),
				("", "+", 3),
				(4, "", 4),
				(6, "-", ""),
				(8, "", 8)))

		base = #(1, 2, 3, 4)
		a = #(1, 2, 3, 4)
		b = #(1, b2, b3, 4)
		Assert(Diff.Three(base, a, b)
			is: #(
				(1, "", 1),
				(2, "#2", "b2"),
				(3, "#2", "b3"),
				(4, "", 4)))

		base = #(1, 2, 3, 4)
		a = #(1, 22, 33, 4)
		b = #(1, 22, 33, 4)
		Assert(Diff.Three(base, a, b)
			is: #(
				(1, "", 1),
				(2, "#", 22),
				(3, "#", 33),
				(4, "", 4)))

		base = #(1, 2, 3, 4)
		a = #(1, a2, a3, 4)
		b = #(1, b2, b3, 4)
		Assert(Diff.Three(base, a, b)
			is: #(
				(1, "", 1),
				(2, "-", ""),
				(3, "-", ""),
				("", "+1", "a2"),
				("", "+1", "a3"),
				("", "+2", "b2"),
				("", "+2", "b3"),
				(4, "", 4)))
		}
	Test_Merge()
		{
		base = #(2, 4, 6, 8)
		a = #(2, 3, 4, 6, 8)
		b = #(2, 4, 8)
		Assert(Diff.Merge(base, a, b) is: #(2, 3, 4, 8))
		}
	}