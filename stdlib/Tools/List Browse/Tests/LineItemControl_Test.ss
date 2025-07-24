// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_GetLineItems()
		{
		// use Mock to avoid re-defining different ctrl classes
		mock = Mock(LineItemControl)
		mock.When.GetAllLineItems().CallThrough()
		mock.When.filterItem([anyArgs:]).CallThrough()
		mock.LineItemControl_exclude = false
		mock.LineItemControl_list = class
			{
			GetLoadedData()
				{
				return Object(
					[a: 1], [vl_expand?:], [vl_expand?:], [a:2], [a:3, vl_deleted:])
				}
			}
		Assert(mock.Eval(LineItemControl.GetAllLineItems)
			is: #([a: 1], [a: 2], [a: 3, vl_deleted:]))
		Assert(mock.Eval(LineItemControl.GetAllLineItems, includeAll?:)
			is: #([a: 1], [a: 2], [a: 3, vl_deleted:]))
		Assert(mock.Eval(LineItemControl.GetLineItems)
			is: #([a: 1], [a: 2]))
		Assert(mock.Eval(LineItemControl.GetDeleted)
			is: #([a: 3, vl_deleted:]))

		mock.LineItemControl_exclude = { it.a < 2 }
		Assert(mock.Eval(LineItemControl.GetAllLineItems)
			is: #([a: 2], [a: 3, vl_deleted:]))
		Assert(mock.Eval(LineItemControl.GetAllLineItems, includeAll?:)
			is: #([a: 1], [a: 2], [a: 3, vl_deleted:]))
		Assert(mock.Eval(LineItemControl.GetLineItems)
			is: #([a: 2]))
		Assert(mock.Eval(LineItemControl.GetDeleted)
			is: #([a: 3, vl_deleted:]))
		}
	}
