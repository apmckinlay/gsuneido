// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Test UI"
	Xmin: 1000
	Ymin: 700
	CallClass()
		{
		Dialog(0, TestUI)
		}
	Controls: (Scroll (Border (Vert,
		(Heading1 'Baseline Align')
		Skip,
		(Static 'The bottom of the text in these should all be aligned')
		Skip,
		(Horz,
			(Static Hello), Skip,
			(Static Hello size: '+4'), Skip,
			(Field set: 'Hello' width: 6), Skip,
			(Field set: 'Hello' size: '+4', width: 6), Skip,
			(ChooseList (Hello World) set: 'Hello', width: 6), Skip,
			(ChooseList (Hello World) set: 'Hello', size: '+4',  width: 6), Skip,
			(ScintillaAddonsField set: 'Hello', width: 6), Skip,
			(ScintillaAddonsField set: 'Hello', fontSize: '+4', width: 6), Skip,
			),
		Skip, Skip,
		(Heading1 'Sizing & Margins')
		Skip,
		(Static 'The text in these should be aligned vertically and ' $
			'should fill the fields')
		(Static 'And the fields should be the same size, ' $
			'excluding the ChooseList buttons')
		Skip,
		(Horz
			(Vert,
				(Field set: 'MMMMMMMMMM' width: 10),
				(ScintillaAddonsField set: 'MMMMMMMMMM', width: 10),
				(ChooseList (MMMMMMMMMM) set: 'MMMMMMMMMM', width: 10),
				(Editor set: 'MMMMMMMMMM', width: 10, height: 2,
					xstretch: false, ystretch: false),
				),
			Skip, Skip,
			(Vert,
				(Field set: 'HelloWorld', font: '@mono', width: 10),
				(ScintillaAddonsField set: 'HelloWorld', font: '@mono', width: 10),
				(ChooseList (HelloWorld) set: 'HelloWorld', font: '@mono', width: 10),
				),
			Skip, Skip,
			(Vert,
				(Field set: 'MMMMMMMMMM' size: '+4', width: 10),
				(ScintillaAddonsField set: 'MMMMMMMMMM', fontSize: '+4', width: 10),
				(ChooseList (Hello World) set: 'MMMMMMMMMM', size: '+4', width: 10),
				)
			Skip, Skip,
			(Vert
				(Number mask: false, set: '1234567890', width: 10)
				(Number mask: '##########', set: '1234567890')
				)
			)
		Skip, Skip,
		(Heading1 'Multi-line Static')
		Skip,
		(Static 'The text should just fit the background')
		Skip,
		(Static
			'now is the time\nfor all good men\nto come to the aid\nof their party',
			bgndcolor: 0xaaaaaa)
		)))
	}