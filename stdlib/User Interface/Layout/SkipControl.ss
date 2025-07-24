// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Control
	{
	New(amount = 10, small = false, medium = false)
		{
		amount = .getScaledSkipAmount(amount, small, medium)
		if not .Parent.Member?("Dir")
			.Xmin = .Ymin = amount
		else if .Parent.Dir is "vert"
			.Ymin = amount
		else if .Parent.Dir is "horz"
			.Xmin = amount
		}

	small: 2
	medium: 5
	getScaledSkipAmount(amount, small, medium)
		{
		if small
			amount = .small
		else if medium
			amount = .medium
		return ScaleWithDpiFactor(amount)
		}

	GetReadOnly() // read-only not applicable to skip
		{
		return true
		}

	CalcXminByControls(@args)
		{
		.Xmin =	.DoCalcXminByControls(@args)
		}
	}
