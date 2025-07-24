// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		// test that record change isnt done if destroyed
		acc = AccessControl
			{
			Send(@unused) { throw "Send Called"}
			}

		// Send should not be called, since record_change_member is false
		acc.AccessControl_record_change()
		}

	Test_get_control_ob()
		{
		get_control_ob = AccessControl.AccessControl_control_ob
		mock = Mock()
		mock.AccessControl_types = Object(DynamicTypes: false)
		mock.AccessControl_fields = #(a, b, c)
		Assert(mock.Eval(get_control_ob, #('stdlib'))
			is: #(Vert, a, b, c))

		Assert(mock.Eval(get_control_ob, #('stdlib', c, d))
			is: #(Form, c, d))

		Assert(mock.Eval(get_control_ob, #('stdlib', #(Vert e, f)))
			is: #(Vert e, f))
		}
	}
