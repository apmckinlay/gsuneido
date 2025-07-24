// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		oldComments = Suneido.GetDefault('svc_comments', false)
		.AddTeardown({
			if oldComments is false
				Suneido.Delete('svc_comments')
			else
				Suneido.svc_comments = oldComments
			})
		}

	Test_getItem()
		{
		get = SvcGetDescription.SvcGetDescription_getComment
		date1 = Date().Format('HH:mm:ss')
		date2 = Date().Plus(minutes: 1).Format('HH:mm:ss')
		Suneido.svc_comments = Object(
			date1 $ ' - minor refactor - fixed long method',
			date2 $ ' - minor refactor - solved magic number')
		Assert(get(date1 $ ' - minor refactor - fixed ') is:
			'minor refactor - fixed long method')
		Assert(get(date2 $ ' - minor refactor') is:
			'minor refactor - solved magic number')
		}

	Test_buildChangeList()
		{
		list = SvcGetDescription.SvcGetDescription_buildChangeList
		Assert(list(#()) is: "")

		Assert(list(#(
			#(lib: "test1lib", type: " ", name: "Record1"),
			#(lib: "test1lib", type: " ", name: "Record2"),
			#(lib: "test2lib", type: " ", name: "Record1")))
			is: 'test1lib:Record1\r\ntest1lib:Record2\r\ntest2lib:Record1')

		Assert(list(#(
			#(lib: "test1lib", type: "+", name: "Record1"),
			#(lib: "test1lib", type: "-", name: "Record2"),
			#(lib: "test2lib", type: "+", name: "Record1")))
			is: '+test1lib:Record1\r\n-test1lib:Record2\r\n+test2lib:Record1')

		}
	}