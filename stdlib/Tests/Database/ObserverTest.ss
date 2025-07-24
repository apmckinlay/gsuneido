// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		s = ""
		f = {|member| s = member $ " = " $ this[member] }
		r = Record()
		r.Observer(f)
		r.a = 34
		Assert(s is: "a = 34")

		r.RemoveObserver(f)
		s = ""
		r.b = 56
		Assert(s is: "")
		}
	Test_invalidate()
		{
		calls = Object().Set_default(0)
		r = Record()
		r.Observer({|member| ++calls[member] })
		r.Invalidate('a')
		Assert(calls is: #(a: 1))
		}
	}