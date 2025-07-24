// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: 'skip'
	SkipSetFocus: true
	New(amount = 10)
		{
		.CreateElement('div')
		if _parent is false or not _parent.Member?("Dir")
			.Xmin = .Ymin = amount
		else if _parent.Dir is "vert"
			.Ymin = amount
		else if _parent.Dir is "horz"
			.Xmin = amount
		.SetMinSize()
		}

	GetReadOnly() // read-only not applicable to skip
		{
		return true
		}

	CalcXminByControls(plusCtrls, minusCtrls)
		{
		.Xmin =	.DoCalcXminByControls(plusCtrls, minusCtrls)
		.SetMinSize()
		}
	}
