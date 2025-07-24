// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_CalcLines()
		{
		mock = Mock(StaticWrapControl)
		mock.Hwnd = 'NULL'
		mock.Hwnd_hfont = 0
		mock.When.StaticWrapControl_bestFit([anyArgs:]).CallThrough()
		mock.StaticWrapControl_measure = { |@unused|  63 }
		Assert(mock.Eval(StaticWrapControl.CalcLines, 'Hello World', 200) is: 1)

		measure = Mock()
		// need to fake the Text Extent as BestFit tries each sub-string
		measure.When.Call([anyArgs:]).Return(
			170, 258, 211, 189, 196, 204,
			 76, 116, 132, 143, 143)
		mock.StaticWrapControl_measure = measure
		Assert(mock.Eval(StaticWrapControl.CalcLines, .msg2, 230) is: 2)

		measure = Mock()
		measure.When.Call([anyArgs:]).Return(
			390, 184, 273, 230, 202, 195, 199,
			311, 147, 226, 187, 206, 187, 199,
			199, 306, 252, 224, 209, 202,
			105, 166, 192, 215, 201,
			 13,  23,  23)
		mock.StaticWrapControl_measure = measure
		Assert(mock.Eval(StaticWrapControl.CalcLines, .msg3, 200) is: 5)
		}
	msg2: 'This text will end up being wider, and should end up on two lines'
	msg3: 'If you have time to enter a short description of what you were ' $
		'doing when this problem occurred, it will help us improve the ' $
		'software. Thank you.'
	}