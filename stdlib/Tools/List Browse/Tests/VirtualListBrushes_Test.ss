// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		brushes = VirtualListBrushes
			{
			VirtualListBrushes_createSystemBrush(color)
				{
				return color
				}
			}()

		Assert(brushes.GetBrush([]) is: false)
		Assert(brushes.GetBrush([abc: 'efg']) is: false)

		brushes.HighlightValues('abc', #('efg', 'hij'), 'red')
		brushes.HighlightValues('abc', #('opq', 'rst'), 'yellow')
		Assert(brushes.GetBrush([]) is: false)
		Assert(brushes.GetBrush([abc: 'efg']) is: 'red')
		Assert(brushes.GetBrush([abc: 'hij']) is: 'red')
		Assert(brushes.GetBrush([abc: 'opq']) is: 'yellow')
		Assert(brushes.GetBrush([abc: 'rst']) is: 'yellow')
		Assert(brushes.GetBrush([abc: 'klm']) is: false)
		}
	}