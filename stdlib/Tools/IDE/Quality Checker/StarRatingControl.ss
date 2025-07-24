// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 		StarRating
	maxStars: 	5
	New()
		{
		.starHorz = .FindControl('starHorz')
		.stars = Object()
		for i in .. .maxStars
			.stars.Add(.FindControl('star' $ i))
		.colors = [CLR.RED, CLR.RED, CLR.RED, CLR.BLACK, CLR.DARKGREEN, CLR.DARKGREEN]
		.starHorz.SetVisible(false)
		}

	starTypes: (
		empty: 'starRatingImageEmpty',
		half: 'starRatingImageHalf',
		full: 'starRatingImageFull')
	Controls()
		{
		horzCtrl = Object(#Horz, name: #starHorz)
		for i in .. .maxStars
			{
			horzCtrl.Add(Object(#Image, .starTypes.empty, xmin: 18, ymin: 18,
				name: #star $ i, alwaysReadOnly?:))
			if i isnt .maxStars - 1
				horzCtrl.Add(#(Skip, 5))
			}
		return horzCtrl
		}

	prevRating: false
	SetRating(rating)
		{
		if rating is .prevRating
			return
		.prevRating = rating
		.starHorz.SetVisible(false)
		if rating is false
			return
		.setStars(rating)
		.starHorz.SetVisible(true)
		}

	setStars(rating)
		{
		ratingScaleRatio = 2 //10:5 - Rating/10 converted to 5 stars
		fullStars = (rating / ratingScaleRatio).RoundDown(0)
		halfStars = rating % ratingScaleRatio

		color = .colors[fullStars]
		for (i = 0; i < .maxStars; i++)
			{
			if 0 < fullStars--
				.stars[i].Set(.starTypes.full, :color)
			else if 0 < halfStars--
				.stars[i].Set(.starTypes.half, :color)
			else
				.stars[i].Set(.starTypes.empty, :color)
			}
		}
	}
