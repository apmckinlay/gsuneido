// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// this is tests for built-in methods
// Objects_Test is for methods defined in Objects
Test
	{
	Test_Add()
		{
		ob = #(12, 34, a: 56, b: 78)
		Assert(ob.Copy().Add() is: ob)
		Assert(ob.Copy().Add(45) is: #(12, 34, 45, a: 56, b: 78))
		Assert(ob.Copy().Add(44, 55) is: #(12, 34, 44, 55, a: 56, b: 78))
		Assert(ob.Copy().Add(99, at: 'c') is: #(12, 34, a: 56, b: 78, c: 99))
		Assert(ob.Copy().Add(88, 99, at: 8) is: #(12, 34, a: 56, b: 78, 8: 88, 9: 99))
		Assert(ob.Copy().Add(0, at: 0) is: #(0, 12, 34, a: 56, b: 78))
		Assert(ob.Copy().Add(23, at: 1) is: #(12, 23, 34, a: 56, b: 78))
		Assert(ob.Copy().Add(45, at: 2) is: #(12, 34, 45, a: 56, b: 78))
		}

	Test_Size()
		{
		ob = #(12, 34, a: 56, b: 78)
		Assert(#() isSize: 0)
		Assert(ob isSize: 4)
		Assert(ob.Size(list:) is: 2)
		}

	Test_Delete()
		{
		ob = #(12, 34, 56, a: 7, b: 8)
		Assert(ob.Copy().Delete(0) is: #(34, 56, a: 7, b: 8))
		Assert(ob.Copy().Delete(1) is: #(12, 56, a: 7, b: 8))
		Assert(ob.Copy().Delete(2) is: #(12, 34, a: 7, b: 8))
		Assert(ob.Copy().Delete(3) is: ob)
		Assert(ob.Copy().Delete('a') is: #(12, 34, 56, b: 8))
		Assert(ob.Copy().Delete('c') is: ob)
		Assert(ob.Copy().Delete(#a, 1) is: #(12, 56, b: 8))

		Assert(ob.Copy().Delete(all:) is: #())
		Assert(ob.Copy().Delete(all: false) is: ob)
		}

	Test_Erase()
		{
		ob = #(12, 34, 56, a: 7, b: 8)
		Assert(ob.Copy().Erase(2) is: #(12, 34, a: 7, b: 8))
		Assert(ob.Copy().Erase('a') is: #(12, 34, 56, b: 8))
		Assert(ob.Copy().Erase('c') is: ob)
		Assert(ob.Copy().Erase(#a, 2) is: #(12, 34, b: 8))
		x = ob.Copy()
		Assert(x.Erase(1) is: #(12, 2: 56, a: 7, b: 8))
		Assert(x.Add(34 at: 1) is: ob)
		}

	Test_Sort()
		{
		x = Object()
		for (i = 0; i < 20; ++i)
			x[i] = Random(10000)
		x.Sort!()
		prev = 0
		for (n in x)
			{
			Assert(n >= prev)
			prev = n
			}

		Assert({ #(1, 2, 3).Sort!() } throws: "readonly")
		}

	Test_Find()
		{
		x = Object()
		for (i = 0; i < 20; ++i)
			x[i] = i
		for (i = 0; i < 20; ++i)
			Assert(x.Find(i) is: i)
		x = Object(a: 12, b: 34, 99: 'x')
		Assert(x.Find(12) is: 'a')
		Assert(x.Find(34) is: 'b')
		Assert(x.Find('x') is: 99)
		}

	Test_Join()
		{
		Assert(Object(123, 456, 789).Join("/") is: "123/456/789")
		Assert(Object(Date.Begin(), "->", Date.End()).Join(" ")
			is: Display(Date.Begin()) $ " -> " $ Display(Date.End()))
		}

	Test_Join_Embedded_Nuls()
		{
		Assert(#("<", "hello\x00world", ">").Join('')
			is: "<hello\x00world>")
		}

	Test_Object?()
		{
		Assert(Object?(false) is: false)
		Assert(Object?("") is: false)
		Assert(Object?(45) is: false)
		Assert(Object?(Object()))
		Assert(Object?(#(test: "test")))
		Assert(Object?(Object(a:1, b:2)))
		Assert(Object?(class{}) is: false)
		Assert(Object?(class{}()) is: false)
		}

	Test_Reverse()
		{
		ob = Object(1, 2, 3, 4, 5)
		ob2 = Object(one: 1, two: 2, three: 3)
		Assert(ob.Reverse!() is: Object(5, 4, 3, 2, 1))
		Assert(Object().Reverse!() is: Object())
		Assert(ob2.Reverse!() is: Object(one: 1, two: 2, three: 3))

		Assert({ #(1, 2, 3).Reverse!() } throws: "readonly")
		}

	Test_Defaults()
		{
		ob = Object(one: 1, two: 2, three: 3)
		ob.Set_default("test")
		Assert(ob.name is: "test")
		Assert(ob.one is: 1)
		Assert(ob.three is: 3)

		ob = Object().Set_default(Object(1).Set_default(0))
		Assert(ob.a is: #(1))
		Assert(ob.a.b is: 0)
		}

	Test_Members()
		{
		Assert(Object(12, 34, a: 56, b: 78) members: #(0, 1, 'a', 'b'))
		Assert(Object() members: #())
		Assert(Object(1, 2, 3, 4, 90, 78).Members() is: Object(0, 1, 2, 3, 4, 5))
		Assert(Field_comment.Members(all:) equalsSet: #(Control, Encode, Format, Prompt))
		}

	Test_Member?()
		{
		ob = Object(one: 1, two: 2, three: 3, test: "test")
		Assert(ob hasMember: "test")
		Assert(ob hasntMember: "non_existent")
		}

	Test_for()
		{
		src = Object(11, 22, a: 33)

		dst = Object()
		for (x in src)
			dst.Add(x)
		Assert(dst is: #(11, 22, 33))

		dst = Object()
		for (x in src.Members())
			dst.Add(x)
		Assert(dst is: #(0, 1, 'a'))

		dst = Object()
		for (x in src.Values())
			dst.Add(x)
		Assert(dst is: #(11, 22, 33))

		dst = Object()
		for (x in src.Values(list:))
			dst.Add(x)
		Assert(dst is: #(11, 22))
		}

	Test_Iterator()
		{
		ob = Object(11, 22, a: 33)

		iter = ob.Iter()
		Assert(iter.Next() is: 11)
		Assert(iter.Next() is: 22)
		Assert(iter.Next() is: 33)
		Assert(iter.Next() is: iter)
		}

	i: 0 // to suppress warning
	Test_Eval()
		{
		func = function()
			{
			Assert(.i is: 5)
			}
		ob = Object(i: 5)
		ob.Eval(func)
		}

	Test_Copy()
		{
		ob = Object(one: 1, two: "two", two2: Object(1, 2, 3), three: 3)
		ob2 = Object(1, 2, 3, 4, 5, 6, 7, 8)
		Assert(Object().Copy() is: Object())
		Assert(ob.Copy()
			is: Object(one: 1, two: "two", three: 3, two2: Object(1, 2, 3)))
		Assert(ob2.Copy() is: Object(1, 2, 3, 4, 5, 6, 7, 8))
		}

	Test_Display()
		{
		ob = #(abc: 123, abc?: 456, "a b": 789)
		Assert(ob is: Display(ob).SafeEval())

		ob = Object()
		strings = #("plain", "single's", "double\"s", "back\\slash")
		for s in strings
			ob.Add(s)
		for s in strings
			ob[s] = true
		Assert(Display(ob).SafeEval() is: ob)
		}

	Test_AtArgs()
		{
		fn = function (@args) { return args }
		ob = #(12, 34, a: 56, b: 78)
		Assert(fn(@ob) is: ob)
		ms = ob.Members()
		Assert(fn(@ms) is: ms)
		}

	Test_SeqBug()
		{
		ob = Object(1, 2, a: 4, b: 5)
		m = ob.Members() // seq
		m[0] // force seq to build object
		n = 0
		for unused in m
			++n
		Assert(ob isSize: n)

		n = 0
		for unused in m
			for unused in m
				++n
		Assert(n is: ob.Size() * ob.Size())

		n = 0
		seq = Seq(4)
		for unused in seq
			for unused in seq
				++n
		Assert(n is: 4 * 4)
		}

	Test_FracBug()
		{
		x = Object()
		x[.5] = 123
		Assert(x.Members() is: #(.5))
		}

	Test_int_if_num_bug()
		{
		ob = Object(5)
		ob[.1] = 5.1
		Assert(ob[.1] is: 5.1)
		}

	Test_Object_as_Key()
		{
		for .. 10
			{
			ob = Object()
			k1 = Object(1, 2, a: 3, b: 4)
			ob[k1] = true
			k2 = Object(1, 2, a: 3, b: 4)
			Assert(ob hasMember: k2)
			}
		}

	Test_set_default_bug()
		{
		x = Object()
		x.Set_default(#())
		x.Set_readonly()
		Assert(x.fred is: #())
		Assert(not x.Member?(#fred))
		}

	Test_ranges()
		{
		Assert(#(a,b,c,d)[1 .. 3] is: #(b,c))
		Assert(#(a,b,c,d)[1 .. 9] is: #(b,c,d))
		Assert(#(a,b,c,d)[1 ..] is: #(b,c,d))
		Assert(#(a,b,c,d)[6 .. 9] is: #())
		Assert(#(a,b,c,d)[2 .. 1] is: #())
		Assert(#(a,b,c,d)[-3 .. -1] is: #(b,c))
		Assert(#(a,b,c,d)[1 .. -1] is: #(b,c))
		Assert(#(a,b,c,d)[.. -2] is: #(a,b))
		Assert(#(a,b,c,d)[-2 ..] is: #(c,d))

		Assert(#(a,b,c,d)[1 :: 2] is: #(b,c))
		Assert(#(a,b,c,d)[:: 2] is: #(a,b))
		Assert(#(a,b,c,d)[-2 :: 1] is: #(c))
		Assert(#(a,b,c,d)[1 :: -1] is: #())
		Assert(#(a,b,c,d)[1 :: 9] is: #(b,c,d))
		Assert(#(a,b,c,d)[1 ::] is: #(b,c,d))
		Assert(#(a,b,c,d)[9 :: 1] is: #())
		}

	a: 0 // to suppress warning
	Test_Eval_with_at()
		{
		f = function (b, aa) { .a + b + aa}
		Assert(Object(a: 1).Eval(@Object(f, 2, 3)) is: 6)
		}

	Test_Unique!()
		{
		Assert([].Unique!() is: [])
		Assert([1].Unique!() is: [1])
		Assert([1, 2].Unique!() is: [1, 2])
		Assert([1, 2, 2, 3].Unique!() is: [1, 2, 3])
		Assert([1, 1, 1, 1].Unique!() is: [1])
		}

	Test_GetDefault()
		{
		Assert([].GetDefault(#a, 123) is: 123)
		Assert([a: 456].GetDefault(#a, 123) is: 456)
		n = 0
		Assert([].GetDefault(#a, { ++n; 123 + 456 }) is: 579)
		Assert(n is: 1)
		Assert([a: 123].GetDefault(#a, { ++n; 123 + 456 }) is: 123)
		Assert(n is: 1)
		}

	Test_get_default_bug()
		{
		x = #()
		Assert(x.GetDefault(#fred, #()) is: #())
		Assert(not x.Member?(#fred))
		}

	Test_ModifiedDuringIteration()
		{
		ob = ''
		test =
			{|block|
			ob = Object(1, 2, 3, 4)
			Assert(block throws: 'object modified')
			}
		test({ for x in ob
				if x is 2
					ob.Delete(0) })
		test({ for x in ob
				if x is 2
					ob.Add(123) })

		test({ for x in ob
				if x is 3
					ob.Delete(0) })
		test({ for x in ob
				if x is 2
					ob.Sort!() })
		test({ for x in ob
				if x is 2
					ob.Reverse!() })
		}

	Test_methods()
		{
		// methods available to Objects, classes, and instances
		Assert(Object().GetDefault(#x, 123) is: 123)
		Assert(class{}.GetDefault(#x, 123) is: 123)
		Assert(class{}().GetDefault(#x, 123) is: 123)
		// methods available to Objects and instances
		Assert(Object(123).Copy() is: #(123))
		c = class
			{
			UseDeepEquals: true
			m: 123
			}
		i = new c
		Assert(i.Copy() is: i)
		Assert({ c.Copy() } throws: "method not found")
		i.foo = 456
		Assert(i.Delete(#foo) is: new c)
		Assert({ c.Delete() } throws: "method not found")
		// methods only available to Objects
		Assert(Object().Add(123) is: #(123))
		Assert({ class{}.Add(123) } throws: "method not found")
		Assert({ class{}().Add(123) } throws: "method not found")
		}

	Test_Join_displays()
		{
		Assert(#(true, 123, 'foo', #()).Join() is: 'true123foo#()')
		}

	Test_AssignOp()
		{
		ob = Object()
		Assert({ ob.n += 5 } throws: "member not found")
		Assert(ob isSize: 0)
		Assert({ ob.n $= "foo" } throws: "member not found")
		Assert(ob isSize: 0)
		}
	Test_for_m_v()
		{
		x = Seq(3).Map({ it })
		ob = Object()
		for m, v in x
			ob.Add([m, v])
		Assert(x is: #(0, 1, 2))
		Assert(ob is: #((0, 0), (1, 1), (2, 2)))
		}
	}