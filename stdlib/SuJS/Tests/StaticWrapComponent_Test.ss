// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_ConvertText()
		{
		mock = Mock(StaticWrapComponent)
		mock.When.BestFit([anyArgs:]).CallThrough()
		mock.StaticComponent_measure = { |@unused|  63 }
		mock.Orig_xmin = 200
		Assert(mock.Eval(StaticWrapComponent.ConvertText, 'Hello World')
			is: 'Hello&nbspWorld')

		measure = Mock()
		// need to fake the Text Extent as BestFit tries each sub-string
		measure.When.Call([anyArgs:]).Return(
			170, 258, 211, 189, 196, 204,
			 76, 116, 132, 143, 143)
		mock.StaticComponent_measure = measure
		mock.Orig_xmin = 230
		Assert(mock.Eval(StaticWrapComponent.ConvertText, .msg2).Split('<br>') isSize: 2)

		measure = Mock()
		measure.When.Call([anyArgs:]).Return(
			390, 184, 273, 230, 202, 195, 199,
			311, 147, 226, 187, 206, 187, 199,
			199, 306, 252, 224, 209, 202,
			105, 166, 192, 215, 201,
			 13,  23,  23)
		mock.StaticComponent_measure = measure
		mock.Orig_xmin = 200
		Assert(mock.Eval(StaticWrapComponent.ConvertText, .msg3).Split('<br>') isSize: 5)
		}
	msg2: 'This text will end up being wider, and should end up on two lines'
	msg3: 'If you have time to enter a short description of what you were ' $
		'doing when this problem occurred, it will help us improve the ' $
		'software. Thank you.'
	}