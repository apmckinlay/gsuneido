// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
HorzComponent
	{
	Name: 'Pair'
	hidden: false
	New(left, skip, right)
		{
		super(left, skip right)
		children = .GetChildren()
		.hidden = children[0].GetHidden() and children[2].GetHidden()
		if .hidden
			.El.SetStyle(#display, 'none')
		}

	Recalc()
		{
		super.Recalc()
		if .hidden
			{
			.Xmin = .Xstretch = .Ystretch = .Ymin = 0
			return
			}
		children = .GetChildren()
		.Left = children[0].Xmin + children[1].Xmin
		// this is to align the label and multi-line editors properly
		// html textarea is a block and doesn't work with baseline alignment properly
		// so we need to set the label's line-height
		if children[0].Name is 'Static' and children[0].Get().LineCount() is 1
			children[0].El.SetStyle('line-height',
				(children[0].Ymin + 8/*=padding*/) $ 'px')
		}

	SetFocus()
		{
		.GetChildren()[2].SetFocus()
		}

	MakeSummary()
		{
		control = .GetChildren()[1]
		return control.Method?(#MakeSummary) ? control.MakeSummary() : ""
		}
	}
