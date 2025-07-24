// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_makecontrols()
		{
		m = VfieldsControl.VfieldsControl_makecontrols
		Assert(m(#()) is: #())
		Assert(m(#(#())) is: #())
		Assert(m(#(#(one two three)))
			is: #(#(one group: 1) 'nl' #(two group: 1) 'nl' #(three group: 1) 'nl'))
		Assert(m(#(one two three))
			is: #(#(one group: 1) 'nl' #(two group: 1) 'nl' #(three group: 1) 'nl'))
		Assert(m(#(#(one) #(two) #(three)))
			is: #(#(one group: 1) 'nl' #(two group: 1) 'nl' #(three group: 1) 'nl'))
		}
	}