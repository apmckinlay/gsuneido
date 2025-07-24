// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		warningsAllMethods = #()
		Assert(QcContinuousWarningsOutput(warningsAllMethods) is: '')

		warningsAllMethods = #(#(warnings: #(), desc: "", rating: 5),
			#(warnings: #(), desc: "", rating: 5),
			#(warnings: #(), desc: "", rating: 5),
			#(warnings: #(), desc: "", rating: 5),
			#(warnings: #(), desc: "", rating: 5),
			#(warnings: #(), desc: "", rating: 5), lineWarnings: #())
		Assert(QcContinuousWarningsOutput(warningsAllMethods) is: '')

		warningsAllMethods = #(#(warnings: #(), desc: "", rating: 5),
			#(warnings: #([name: "stdlib:Init:34 - McCabe complexity is 9 for startup"]),
				desc:"McCabe function complexity",
				rating: 5),
			#(size: 1, warnings: #(), maxRating: 4,
				desc: "A test class was not found for this record. Please create one"),
				lineWarnings: #())
		Assert(QcContinuousWarningsOutput(warningsAllMethods)
			is: '\nMcCabe function complexity:\n' $
			"stdlib:Init:34 - McCabe complexity is 9 for startup\n\n" $
			"A test class was not found for this record. Please create one:\n")

		Assert(QcContinuousWarningsOutput(#((warnings: ([name: 'test warnings'],
				[name: "second test warning"], [name: "third test warning"]),
			desc: 'test description')))
			is: "\ntest description:\ntest warnings\nsecond test warning\n" $
				"third test warning\n")

		Assert(QcContinuousWarningsOutput("") is: "")
		}
	}