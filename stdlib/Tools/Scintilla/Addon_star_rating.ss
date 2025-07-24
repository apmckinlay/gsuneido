// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
CodeViewAddon
	{
	Name:		StarRating
	Inject:		topRight
	Controls(addonControls)
		{
		if QcIsEnabled()
			addonControls.Add(#(Horz, Fill, (StarRating)), at: 5)
		}

	Set()
		{ .CheckCode_QualityChanged([rating: false]) }

	Addon_RedirMethods()
		{ return #(CheckCode_QualityChanged) }

	CheckCode_QualityChanged(checks)
		{
		if .AddonControl isnt false
			.AddonControl.SetRating(checks.rating)
		}
	}
