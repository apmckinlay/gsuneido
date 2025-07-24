// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
HorzComponent
	{
	Name: "HorzEqual"
	New(@args)
		{
		super(@.init(args))
		}

	init(args)
		{
		.pad = args.GetDefault('pad', 20/*=default pad*/)
		return args
		}

	Recalc()
		{
		super.Recalc()
		maxXmin = 0
		for c in .GetChildren()
			if .isButton?(c)
				maxXmin = Max(maxXmin, c.Xmin)
		maxXmin += .pad
		.Xmin = 0
		for c in .GetChildren()
			{
			if .isButton?(c)
				{
				c.Xmin = maxXmin
				c.El.SetStyle('flex-basis', maxXmin $ 'px')
				}
			.Xmin += c.Xmin
			}
		.SetMinSize()
		}

	isButton?(c)
		{
		return c.Base?(ButtonComponent) or c.Base?(EnhancedButtonComponent)
		}
	}
