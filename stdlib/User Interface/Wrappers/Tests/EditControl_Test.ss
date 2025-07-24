// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_callContext()
		{
		_result = Object()
		cl = EditControl
			{
			Default(@args)
				{
				_result.Add(args[0])
				}
			}
		fCE = cl.EditControl_callContext
		contextExtra = Object(
			"General Item 1"
			"General Item 2"
			""
			"Regular Cascade"
			#("Regular Option 1", "Regular Option 2")
			""
			Object(name: 'Item 1', status: 0, runFunc: { fnResult = "Item 1"}),
			Object(name: 'Item 2', status: 0, runFunc: { fnResult = "Item 2"}),
			Object(name: '', status: 0, runFunc: false),
			Object(name: 'Item 3', status: 0, runFunc: { fnResult = "Item 3" }),
			Object(name: 'Cascade', status: 0, runFunc: false),
			Object(name: Object(
				Object(name: 'Option 1', status: 0, runFunc: { fnResult = "Option 1" })
				Object(name: 'Option 2', status: 0, runFunc: { fnResult = "Option 2" }))
				status: 0, runFunc: false),
			Object(name: 'Item 4', status: 0, runFunc: { fnResult = "Item 4" })
			"General Item 3"
			)

		// test ContextMenuCall
		fCE(contextExtra, 0)
		Assert(_result.Last() is: "On_General_Item_1")
		fCE(contextExtra, 1)
		Assert(_result.Last() is: "On_General_Item_2")
		fCE(contextExtra, 4)
		Assert(_result.Last() is: "On_Regular_Option_1")
		fCE(contextExtra, 5)
		Assert(_result.Last() is: "On_Regular_Option_2")
		fCE(contextExtra, 15)
		Assert(_result.Last() is: 'On_General_Item_3')

		// test runFunc()
		fCE(contextExtra, 7)
		Assert(fnResult is: 'Item 1')
		fCE(contextExtra, 8)
		Assert(fnResult is: 'Item 2')
		fCE(contextExtra, 10)
		Assert(fnResult is: 'Item 3')
		fCE(contextExtra, 12)
		Assert(fnResult is: 'Option 1')
		fCE(contextExtra, 13)
		Assert(fnResult is: 'Option 2')
		fCE(contextExtra, 14)
		Assert(fnResult is: 'Item 4')
		}
	}