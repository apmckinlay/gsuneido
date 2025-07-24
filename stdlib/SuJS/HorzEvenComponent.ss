// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
HorzComponent
	{
	Name: "HorzEven"
	Xstretch: 1
	Recalc()
		{
		children = .GetChildren()
		for child in children
			{
			if child.Base?(SkipComponent)
				continue
			child.El.SetStyle('flex-basis', 0)
			child.Xstretch = 1
			}
		super.Recalc()
		}
	}
