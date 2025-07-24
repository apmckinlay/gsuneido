// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	C: class {  }
	Test_one()
		{
		Assert(Name(Print) is: 'Print') // builtin
		Assert(Name(Stack) is: 'Stack')
		Assert(Name(Stack().Base()) is: 'Stack')
		Assert(Name(NameTest.C) is: 'NameTest.C')
		Assert(Name(this.C) is: 'NameTest.C')
		Assert(Name(Stack.Push) is: 'Stack.Push')
		Assert(Name(Stack().Push) is: 'Stack.Push')
		}
	Test_gsuneido()
		{
		f = function (){}
		Assert(Name(f) is: "NameTest.Test_gsuneido f")

		c = class { F(){} }
		Assert(Name(c) is: "NameTest.Test_gsuneido c")
		Assert(Name(c.F) is: "NameTest.Test_gsuneido c.F")
		}
	}