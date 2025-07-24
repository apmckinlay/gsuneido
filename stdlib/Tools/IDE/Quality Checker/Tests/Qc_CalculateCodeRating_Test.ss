// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getRating()
		{
		threshold = #(70, 50, 30)

		starRatings = #(1: (.13, .35, .60),
						2: (.10, .30, .55),
						3: (.08, .26, .50),
						4: (.07, .22, .44),
						5: (.05, .15, .30))

		getRating = Qc_CalculateCodeRating
		testMethodSizes = .buildThreshold(#(5: 71, 10: 51, 15: 31, 70: 15))
		Assert(getRating(testMethodSizes, threshold, starRatings) is: 5)
		testMethodSizes[5] = 71
		Assert(getRating(testMethodSizes, threshold, starRatings) is: 4)

		testMethodSizes = .buildThreshold(#(7: 71, 15: 51, 22: 31, 56: 15))
		Assert(getRating(testMethodSizes, threshold, starRatings) is: 4)
		testMethodSizes[7] = 71
		Assert(getRating(testMethodSizes, threshold, starRatings) is: 3)

		testMethodSizes = .buildThreshold(#(8: 71, 18: 51, 24: 31, 50: 15))
		Assert(getRating(testMethodSizes, threshold, starRatings) is: 3)
		testMethodSizes[8] = 71
		Assert(getRating(testMethodSizes, threshold, starRatings) is: 2)

		testMethodSizes = .buildThreshold(#(10: 71, 20: 51, 25: 31, 45: 15))
		Assert(getRating(testMethodSizes, threshold, starRatings) is: 2)
		testMethodSizes[10] = 71
		Assert(getRating(testMethodSizes, threshold, starRatings) is: 1)

		testMethodSizes = .buildThreshold(#(13: 71, 22: 51, 25: 31, 40: 15))
		Assert(getRating(testMethodSizes, threshold, starRatings) is: 1)
		testMethodSizes[13] = 71
		Assert(getRating(testMethodSizes, threshold, starRatings) is: 0)
		}
	buildThreshold(thresholds)
		{
		testMethodSizes = Object()
		for threshold in thresholds.Members()
			for .. threshold
				testMethodSizes.Add(thresholds[threshold])
		return testMethodSizes
		}

	Test_passStarRating?()
		{
		threshold = #(70, 50, 30)
		passStarRating = Qc_CalculateCodeRating.Qc_CalculateCodeRating_passStarRating?
		allowedPercentages = Object(.25, .5, .75)
		conditions = #(
			#(methodSizes: #(15, 15, 15, 15), result: true),
			#(methodSizes: #(71, 51, 31, 15), result: true),
			#(methodSizes: #(71, 71, 31, 15), result: false),
			#(methodSizes: #(71, 51, 51, 15), result: false),
			#(methodSizes: #(71, 51, 31, 31), result: false))
		for condition in conditions
			Assert(passStarRating(allowedPercentages, condition.methodSizes, threshold)
				is: condition.result)
		}
	}
