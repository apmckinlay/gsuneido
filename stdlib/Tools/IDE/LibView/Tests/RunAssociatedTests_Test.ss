// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_tests()
		{
		m = RunAssociatedTests.RunAssociatedTests_tests
		Assert(m(tabsOb = []) is: [])

		tabsOb.Add([path: "/Folder", group:])
		Assert(m(tabsOb) is: #())

		tabsOb.Add([path: "/RecA", group: false])
		Assert(m(tabsOb) equalsSet: #(RecA_Test, RecATest))

		tabsOb.Add([path: "/RecB?", group: false])
		tabsOb.Add([path: "/RecTest", group: false])
		tabsOb.Add([path: "/Rec_Test", group: false])
		Assert(m(tabsOb)
			equalsSet: #(RecA_Test, RecATest, RecBTest, RecB_Test, RecTest, Rec_Test))
		}
	}
