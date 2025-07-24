// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
HorzControl
	{
	Name: 'Pair'
	hidden: false
	New(left, right)
		{
		super(@.layout(left, right))
		.setLeft()
		}
	layout(left, right)
		{
		if Object?(left) and left.GetDefault('hidden', false) is true and
			Object?(right) and right.GetDefault('hidden', false) is true
			.hidden = true
		return Object(left, #(Skip 6), right)
		}
	setLeft()
		{
		children = .GetChildren()
		.Left = children[0].Xmin + children[1].Xmin
		}
	SetFocus()
		{
		.GetChildren()[2].SetFocus()
		}
	Recalc()
		{
		.setLeft()
		super.Recalc()
		if .hidden
			.Xmin = .Xstretch = .Ystretch = .Ymin = 0
		}
	MakeSummary()
		{
		control = .GetChildren()[2]
		return control.Method?(#MakeSummary) ? control.MakeSummary() : ""
		}
	}