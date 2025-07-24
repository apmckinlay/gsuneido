// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Default_pass_through()
		{
		Assert({ ListViewControl.Foo() } throws: "method not found")
		}
	Test_Default_sends()
		{
		mock = Mock()
		mock.Eval(ListViewControl.Default, 'On_Context_Foo', item: 'Foo')
		mock.Verify.Send('On_Foo')
		mock.Verify.Send('On_Context', 'Foo')
		}
	Test_Default_on_method()
		{
		mock = Mock { On_Foo() { } }()
		mock.Eval(ListViewControl.Default, 'On_Context_Foo', item: 'Foo')
		mock.Verify.Never().Send('On_Foo')
		mock.Verify.Send('On_Context', 'Foo')
		}
	Test_Default_on_send_original_item()
		{
		mock = Mock { On_Test_Foo() { } }()
		mock.Eval(ListViewControl.Default, 'On_Context_Test_Foo', item: 'Test Foo')
		mock.Verify.Never().Send('On_Test_Foo')
		mock.Verify.Send('On_Context', 'Test Foo')
		}
	}