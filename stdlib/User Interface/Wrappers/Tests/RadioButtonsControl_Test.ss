// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidData?()
		{
		vd = RadioButtonsControl.ValidData?

		Assert({ vd() } throws: "member not found: 0")
		args = Object()
		Assert({ vd(@args) } throws: "member not found: 0")

		args = [#test1, #test1, #test2]
		Assert(vd(@args))

		args = [#test1, #test2]
		Assert(vd(@args) is: false)

		args = ["", #test1, #test2]
		Assert(vd(@args))

		args = ["", #test1, #test2, mandatory:]
		Assert(vd(@args) is: false)
		}
	}
