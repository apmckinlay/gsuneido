// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Circlog(10).ToString() is: "")

		cl = Circlog(5)
		for i in ..5
			cl.Append(i)
		Assert(cl.ToString() is: Seq(5).Join('\n') $ '\n')

		cl = Circlog(5)
		for i in ..8
			cl.Append(i)
		Assert(cl.ToString() is: Seq(3, 8).Join('\n') $ '\n')
		}
	}