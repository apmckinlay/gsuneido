// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_saveWorst()
		{
		worst = Object()
		for .. 10
			worst.Add(#(name: '', rating: 5))
		Qc_Stats.Qc_Stats_saveWorst(4, worst, "fourStar")
		Qc_Stats.Qc_Stats_saveWorst(3, worst, "threeStar")
		Qc_Stats.Qc_Stats_saveWorst(4.5, worst, "4.5Star")
		Qc_Stats.Qc_Stats_saveWorst(1, worst, "oneStar")
		Qc_Stats.Qc_Stats_saveWorst(5, worst, "shouldNotSeeMe")
		Qc_Stats.Qc_Stats_saveWorst(2, worst, "twoStar")
		Qc_Stats.Qc_Stats_saveWorst(2.5, worst, "2.5Star")
		Qc_Stats.Qc_Stats_saveWorst(3.5, worst, "3.5Star")
		Qc_Stats.Qc_Stats_saveWorst(1, worst, "oneStar2")
		Qc_Stats.Qc_Stats_saveWorst(4.5, worst, "4.5Star")
		Qc_Stats.Qc_Stats_saveWorst(1.5, worst, "1.5Star")
		Qc_Stats.Qc_Stats_saveWorst(4, worst, "fourStar")
		Qc_Stats.Qc_Stats_saveWorst(5, worst, "fiveStar")

		expectedResult = #([name: "4.5Star", rating: 4.5], [name: "fourStar", rating: 4],
			[name: "fourStar", rating: 4], [name: "3.5Star", rating: 3.5],
			[name: "threeStar", rating: 3], [name: "2.5Star", rating: 2.5],
			[name: "twoStar", rating: 2], [name: "1.5Star", rating: 1.5],
			[name: "oneStar2", rating: 1], [name: "oneStar", rating: 1])
		calculatedResult = worst
		Assert(calculatedResult is: expectedResult)
		}
	}