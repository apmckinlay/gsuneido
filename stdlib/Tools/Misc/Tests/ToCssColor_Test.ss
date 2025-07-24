// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_one()
		{
		testCases = Object(
			Object(CLR.RED, #(r: 255, g: 0, b: 0), '#ff0000')
			Object(CLR.BLUE, #(r: 0, g: 0, b: 255), '#0000ff')
			Object(CLR.GREEN, #(r: 0, g: 255, b: 0), '#00ff00')
			Object(CLR.YELLOW, #(r: 255, g: 255, b: 0), '#ffff00')
			Object(CLR.ErrorColor, #(r: 255, g: 127, b: 127), '#ff7f7f')
			Object(CLR.BLACK, #(r: 0, g: 0, b: 0), '#000000')
			Object(CLR.WHITE, #(r: 255, g: 255, b: 255), '#ffffff')
			)
		testCases.Each()
			{
			Assert(ToCssColor(it[0]) is: it[2])
			Assert(ToCssColor(it[1]) is: it[2])
			}

		Assert(ToCssColor("red") is: "red")
		Assert({ ToCssColor(true) } throws: "unhandled")
		}
	}