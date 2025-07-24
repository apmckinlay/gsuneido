// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_compare_types()
		{
		// test comparing objects with sequences (e.g. Members)
		values = [
			// [value, order]
			[#(), 0],
			[#().Members(), 0, 'm'],
			[#(0), 1],
			[#(5).Members(), 1, 'm'],
			[#(0,1), 2],
			[#(5,6).Members(), 2, 'm'],
			[#(9), 3],
			]
		for x in values
			for y in values
				{
				disp = {|c| Display(x) $ ' ' $ c $ ' ' $ Display(y) }
				Assert((x[0] is y[0]) is: (x[1] is y[1]), msg: disp('='))
				Assert((x[0] < y[0]) is: (x[1] < y[1]), msg: disp('<'))
				Assert((x[0] > y[0]) is: (x[1] > y[1]), msg: disp('>'))
				}
		}
	}