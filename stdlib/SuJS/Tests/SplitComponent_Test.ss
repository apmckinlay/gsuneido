// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_initSplit()
		{
		m = SplitComponent.SplitComponent_initSplit

		// Horizontal / Xstretch test:
		Assert(m(#(Xstretch: false), #(Xstretch: false), #horz) is: #(.5, .5))
		Assert(m(#(Xstretch: .25), #(Xstretch: .25), #horz) is: #(.5, .5))
		Assert(m(#(Xstretch: 4), #(Xstretch: false), #horz) is: #(1, 0))
		Assert(m(#(Xstretch: false), #(Xstretch: 4), #horz) is: #(0, 1))
		Assert(m(#(Xstretch: 4), #(Xstretch: 0), #horz) is: #(1, 0))
		Assert(m(#(Xstretch: 0), #(Xstretch: 4), #horz) is: #(0, 1))
		Assert(m(#(Xstretch: 1), #(Xstretch: 4), #horz) is: #(.2, .8))
		Assert(m(#(Xstretch: 3), #(Xstretch: 1), #horz) is: #(.75, .25))
		Assert(m(#(Xstretch: 2), #(Xstretch: 4), #horz)
			is: #(.3333333333333333, .6666666666666666))

		// Vertical / Ystretch test:
		Assert(m(#(Ystretch: false), #(Ystretch: false), #vert) is: #(.5, .5))
		Assert(m(#(Ystretch: .75), #(Ystretch: .75), #vert) is: #(.5, .5))
		Assert(m(#(Ystretch: 4), #(Ystretch: false), #vert) is: #(1, 0))
		Assert(m(#(Ystretch: false), #(Ystretch: 4), #vert) is: #(0, 1))
		Assert(m(#(Ystretch: 4), #(Ystretch: 0), #vert) is: #(1, 0))
		Assert(m(#(Ystretch: 0), #(Ystretch: 4), #vert) is: #(0, 1))
		Assert(m(#(Ystretch: 1), #(Ystretch: 4), #vert) is: #(.2, .8))
		Assert(m(#(Ystretch: 3), #(Ystretch: 1), #vert) is: #(.75, .25))
		Assert(m(#(Ystretch: 2), #(Ystretch: 4), #vert)
			is: #(.3333333333333333, .6666666666666666))
		}
	}