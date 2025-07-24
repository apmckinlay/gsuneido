// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	RandomCLR()
		{
		colors = Object()
		stdColors = CLR.Copy()
		for clr in #(NONE, INVALID, DEFAULT)
			stdColors.Delete(clr)
		for item in stdColors.Members()
			{
			rgb = stdColors[item].ToRGB()
			if not .BadContrast?(rgb)
				colors.Add(stdColors[item])
			}
		return colors[Random(colors.Size())]
		}

	// "nearest" arg reduces chance that next RGB color generated will be within
	// a few digits of the last one (making them hard to tell apart)
	RandomRGB(nearest = 25)
		{
		c = Object()
		for (i = 0; i < 3; i++)
			c.Add(Random(255).RoundToNearest(nearest))
		return RGB(c[0], c[1], c[2])
		}

	GetContrast(color, dark = 0x000000, light = 0xffffff)
		{
		colorNumber = String?(color) and color.Prefix?('0x')
			? color.SafeEval() 		// convert hex to dec
			: color
		if .BadContrast?(colorNumber.ToRGB())
			return light
		return dark
		}

	// determines how well color passed in shows up against dark colors
	BadContrast?(rgb)
		{
		contrast = 128 * 3 		// represents accum of R,G,B
		return rgb.Sum() < contrast and not rgb.Any?({ it > 230 })
		}
	}