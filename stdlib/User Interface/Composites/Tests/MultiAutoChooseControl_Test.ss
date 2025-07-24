// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getValues()
		{
		getValues = MultiAutoChooseControl.MultiAutoChooseControl_getValues
		Assert(getValues('a,b,c') is: #(a, b, c))
		Assert(getValues('a,,,,b\r\n,,,,,\nc\r\nd') is: #(a, b, c, d))
		Assert(getValues('a,  b,  ,  c') is: #(a, b, c))

		// test alternate Join char (; for email addresses)
		testCl = MultiAutoChooseControl { AlternateJoinChar: ';' }
		getValues = testCl.MultiAutoChooseControl_getValues
		Assert(getValues('test@test.com, test2@test.com')
			is: #('test@test.com', 'test2@test.com'))
		Assert(getValues('test@test.com, test2@test.com;test3@test.com')
			is: #('test@test.com', 'test2@test.com', 'test3@test.com'))
		Assert(getValues('test@test.com; test2@test.com;test3@test.com')
			is: #('test@test.com', 'test2@test.com', 'test3@test.com'))
		}
	Test_ValidData()
		{
		args = Object('str')
		Assert(MultiAutoChooseControl.ValidData?(@args) is: false, msg: 'no list')

		args.list = #('str')
		Assert(MultiAutoChooseControl.ValidData?(@args), msg: 'with ob list')

		args.list = #('str2')
		Assert(MultiAutoChooseControl.ValidData?(@args) is: false, msg: 'not in ob list')

		args.list = 'test'
		Assert(MultiAutoChooseControl.ValidData?(@args) is: false, msg: 'with str list')

		args.record = Record(test: #('str'))
		Assert(MultiAutoChooseControl.ValidData?(@args), msg: 'with record')

		args[0] = 'str,str1'
		Assert(MultiAutoChooseControl.ValidData?(@args) is: false, msg: 'missing str1')

		args.record = Record(test: #('str', 'str1'))
		Assert(MultiAutoChooseControl.ValidData?(@args), msg: 'multi valid')
		}
	}