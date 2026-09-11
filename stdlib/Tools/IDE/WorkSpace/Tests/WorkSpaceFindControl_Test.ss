// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_isExpr?()
		{
		fn = WorkSpaceFindControl.WorkSpaceFindControl_isExpr?
		Assert(fn("a+b"))
		Assert(fn('a[..d]+b() * ("a" $ c).Size()'))

		Assert(not fn("try a + b"))
		Assert(not fn("if a + b"))
		Assert(not fn("return a + 1"))

		Assert(not fn("a|"))
		Assert(not fn("a |b| test()"))
		Assert(not fn('a| a+b() * ("a" $ c).Size()'))
		}
	}
