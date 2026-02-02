// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Base?()
		{
		class11 = class {}
		class2 = ObjectTestClass { }
		Assert(class2.Base?(ObjectTestClass))
		Assert(class2.Base?(class11) is: false)
		}
	Test_Base()
		{
		class11 = class {}
		class2 = ObjectTestClass { }
		ob = ObjectTestClass()
		ob2 = class2()
		ob3 = class11()
		Assert(ob.Base() is: ObjectTestClass)
		Assert(ob2.Base() is: class2)
		Assert(ob3.Base() is: class11)
		}
	Test_Method()
		{
		otc = ObjectTestClass { }
		for x in [ObjectTestClass, otc, ObjectTestClass(), otc()]
			{
			Assert(x.Method?(#Public_method))
			Assert(x.MethodClass(#Public_method) is: ObjectTestClass)
			Assert(x.Method?(#Public_member) is: false)
			Assert(x.MethodClass(#Public_member) is: false)
			}
		}
	Test_MethodClass()
		{
		class2 = ObjectTestClass { method2() { } Public_method2() { } }
		ob = class2()
		Assert(ob.MethodClass("method1") is: false)
		Assert(ob.MethodClass("Public_method") is: ObjectTestClass)
		Assert(ob.MethodClass("method2") is: false)
		Assert(ob.MethodClass("Public_method2") is: class2)
		Assert(ob.MethodClass("non_existent") is: false)
		}
	Test_AutoMethod()
		{
		c = class
			{
			New() { .X = 456 }
			F() { .X }
			X: 123
			}
		m = 'F'

		for f in [c.F, c[m]]
			Assert(f() is: 123)

		x = c()
		for f in [x.F, x[m]]
			Assert(f() is: 456)
		}

	Test_Default()
		{
		c = class
			{
			}
		Assert({ c.F() } throws: 'method not found')

		c = class
			{
			F() { return 123 }
			Default(@args) { return args }
			}
		Assert(c.F() is: 123)
		Assert(c.X() is: #(X))
		Assert(c.X(1, 2, a: 3, b: 4) is: #(X, 1, 2, a: 3, b: 4))
		}
	Test_GetDefault()
		{
		c = class
			{
			X: 123
			Getter_Y()
				{ return 456 }
			}
		x = c()
		Assert(x.GetDefault(#nonexistent, 123) is: 123)

		// in the parent class
		Assert(x.GetDefault(#X, false) is: 123)

		// getter
		Assert(x.Member?(#Y) is: false)
		Assert(x.GetDefault(#Y, false) is: 456)
		Assert(x.Y is: 456)

		Assert(c.GetDefault(#nonexistent, 123) is: 123)
		Assert(c.GetDefault(#X, false) is: 123)
		Assert(c.Member?(#Y) is: false)
		Assert(c.GetDefault(#Y, false) is: 456)
		Assert(c.Y is: 456)

		n = 0
		Assert(c.GetDefault(#a, { ++n; 123 + 456 }) is: 579)
		Assert(n is: 1)
		Assert(c.GetDefault(#X, { ++n; 123 + 456 }) is: 123)
		Assert(n is: 1)
		}
	Test_Getter2()
		{
		c = class
			{
			data: (a: 123)
			Getter_(member)
				{
				if .data.Member?(member)
					return .data[member]
				else
					return Nothing()
				}
			}
		ob = c()
		Assert(ob.a is: 123)
		Assert(ob.GetDefault(#a, false) is: 123 msg: 'GetDefault')
		Assert(ob.GetDefault(#b, false) is: false msg: 'GetDefault no value')
		}
	Test_Getter_()
		{
		c = class
			{
			X: 123
			Getter_Y() { return 456 } // will never be used
			Getter_(m) { return m } // has priority
			}
		ob = c()
		Assert(ob.X is: 123)
		Assert(ob.Y is: 'Y')
		Assert(ob.Z is: 'Z')

		Assert(c.X is: 123)
		Assert(c.Y is: 'Y')
		Assert(c.Z is: 'Z')

		c1 = class
			{
			A: 123
			Getter_X()
				{ return .A }
			getter_x()
				{ return 789 }
			F()
				{ return .x }
			}
		ob = c1()
		ob.A = 456
		Assert(ob.X is: 456)
		Assert(ob.F() is: 789)

		c2 = class
			{
			A: 123
			Getter_(unused)
				{ return .A }
			}
		ob = c2()
		ob.A = 456
		Assert(ob.X is: 456)

		Assert(c1.X is: 123)
		Assert(c1.F() is: 789)
		Assert(c2.X is: 123)
		}

	Test_class_members_must_be_named()
		{
		Assert({ "class { x }".Compile() } throws: "class members must be named")
		}
	Test_arguments_to_New()
		{
		c = class { New(a/*unused*/,b/*unused*/,c/*unused*/) { } }
		c(1,2,3)
		c(@[1,2,3])
		new c(1,2,3)
		new c(@[1,2,3])
		Assert({ c() } throws: "missing argument")
		Assert({ c(@[]) } throws: "missing argument")
		Assert({ c(1,2,3,4) } throws: "too many arguments")
		Assert({ c(@[1,2,3,4]) } throws: "too many arguments")
		Assert({ new c() } throws: "missing argument")
		Assert({ new c(@[]) } throws: "missing argument")
		Assert({ new c(1,2,3,4) } throws: "too many arguments")
		Assert({ new c(@[1,2,3,4]) } throws: "too many arguments")

		c = class { }
		c()
		c(a: 1)
		c(@[])
		c(@[a: 1])
		new c()
		new c(a: 1)
		new c(@[])
		new c(@[a: 1])
		Assert({ c(1,2,3) } throws: "too many arguments")
		Assert({ c(@[1,2,3]) } throws: "too many arguments")
		Assert({ new c(1,2,3) } throws: "too many arguments")
		Assert({ new c(@[1,2,3]) } throws: "too many arguments")
		}
	Test_equals()
		{
		// a class is only equal to itself
		c = class { }
		Assert(c is: c)
		Assert(c isnt: class { })
		// instances are equal if same class and members
		c = class
			{
			UseDeepEquals: true
			New(.x) {}
			}
		i1 = c(123)
		i2 = c(123)
		Assert(i1 is: i2)
		}
	Test_ToString()
		{
		c = class {
			New(x, y) { .x = x; .y = y }
			ToString() { 'c(x:' $ Display(.x) $ ', y: ' $ Display(.y) $ ')' } }
		x = c(12, c(34, 'abc'))
		Assert(Display(x) is: 'c(x:12, y: c(x:34, y: "abc"))')
		}
	Test_ToString2()
		{
		c = class
			{
			New(.x)
				{
				}
			ToString()
				{
				return "MyClass(" $ .x $ ')'
				}
			}
		x = c(123)
		Assert(Display(x) is: "MyClass(123)")
		Assert('=' $ x $ '=' is: "=MyClass(123)=")

		c = class
			{
			ToString()
				{
				// no return value
				}
			}
		x = c()
		Assert({ Display(x) } throws: "ToString should return a string")

		c = class
			{
			ToString()
				{
				return #()
				}
			}
		x = c()
		Assert({ Display(x) } throws: "ToString should return a string")
		}
	Test_method_not_found()
		{
		Assert({ Stack.FooBar() } throws: "method not found: class.FooBar")
		Assert({ Stack().FooBar() } throws: "method not found: instance.FooBar")
		Assert({ Object().FooBar() } throws: "method not found: object.FooBar")
		Assert({ Record().FooBar() } throws: "method not found: record.FooBar")
		}
	Test_instance_of_instance()
		{
		c = class { }
		x = new c
		Assert({ new x } throws: "can't create instance of instance")
		}
	Test_inherit_Objects_methods()
		{
		Assert(class{X: 123}.Val_or_func("X") is: 123)
		}
	Test_Synchronized_Unload()
		{
		// can't really tell if this is "working" without inserting debugging in gSuneido
		// but at least we can make sure it doesn't throw an error

		// instance
		.Synchronized({})
		Unload(#ClassTest)
		.Synchronized({})

		// class = this
		ClassTest.SynchTest()

		c = ClassTest
			{
			F() { .SynchTest() }
			}
		c.F()
		}
	SynchTest()
		{
		.Synchronized({})
		Unload(#ClassTest)
		.Synchronized({})
		}
	}