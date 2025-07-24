// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// used for Mock_Test
class
	{
	Func(a)
		{
		val = .Func1(a)
		.func3()
		val2 = .func2(a)
		return val + val2 + 1
		}

	Func1(a)  // testing public method
		{
		.f1 = a * a
		return .f1
		}

	func2(a) // testing private method
		{
		return a * .f1
		}

	func3()
		{
		.f1 = .f1 + 1 // testing members
		return // also testing no return value
		}

	F() { .f1()	}

	f1()
		{
		sum = .f2()
		sum += .f3()
		.F1CalledThrough = true
		return sum
		}
	f2()
		{
		.F2CalledThrough = true
		return 2
		}

	f3()
		{
		throw 'this method should not be called'
		}

	m: 10
	M: 5
	Foo()
		{
		.foo // to fix unused method

		return .m + .M
		}

	foo() // used by Mock_Test
		{ return .M * .m }

	M1(ob)
		{
		return .m1(ob)
		}
	m1(ob)
		{
		ob.M2()
		return .m2()
		}

	m2()
		{
		return 'from m2'
		}

	Foo1()
		{
		ob = Object()
		ob.Add(.Foo2(1))
		ob.Add(.foo3(2))
		ob.Add(#(3, 4).Map(.Foo2)).Add(#(5, 6).Map(.foo3))

		ob.Add(.Foo4(1))
		ob.Add(.Foo4(2))
		ob.Add(#(3, 4).Map(.Foo4)).Add(#(5, 6).Map(.foo5))
		return ob
		}

	Foo2(unused)
		{
		throw 'should not reach here'
		}

	foo3(unused)
		{
		throw 'should not reach here'
		}

	Foo4(unused)
		{
		throw 'should not reach here'
		}

	foo5(unused)
		{
		throw 'should not reach here'
		}
	}
