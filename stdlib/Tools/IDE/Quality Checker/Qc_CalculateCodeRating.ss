// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	MaxRating: 5

	CallClass(checkCounts, thresholds, starRatings)
		{
		if not checkCounts.Sorted?(Gt)
			checkCounts.Sort!(Gt)
		rating = 5
		while rating > 0
			{
			if .passStarRating?(starRatings[rating], checkCounts, thresholds)
				break
			rating--
			}
		return rating
		}

	passStarRating?(ratingAllowedPercentages, checkCounts, thresholds)
		{
		for (i = 0; i < thresholds.Size(); i++)
			{
			violations = checkCounts.Filter({ it > thresholds[i] }).Size()
			percentage = violations / checkCounts.Size()

			if percentage > ratingAllowedPercentages[i] //If any threshold fails
				return false
			}
		return true
		}
	}


