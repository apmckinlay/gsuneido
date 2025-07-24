// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	testCases: (
		test1: (code: "
		class
			{
				x: 55
			PublicMethod()
				{
				//Should not show privateMethod2 as unused
				ob = Object()
				ob.TestMethod2() //should not read as private method
				//Undefined
				.undefinedPrivMethod1()
				.undefinedPrivMethod2( )
				.undefinedPrivMethod3 ()
				.undefinedPrivMethod4 ( )
				y = 5 + 55 + x - .
				undefinedPrivMethod5()
				//Class members
				.x = 'xxx'
				.usedMethod1()
				}
			unusedPrivMethod1()
				{
				}
			unusedPrivMethod2 ()
				{
				}
			unusedPrivMethod3   ()
				{
				randomOb.objectMember = function(){ }
				randomOb = 55 + .
					undefinedPrivMethod6  (arg1, arg2)
				}
			usedMethod1 ()
				{
				.UndefinedPubMethod1(x,y,z)
				Display(.usedMethod2)
				}
			usedMethod2()
			{
			Print('HI')
			if true
				return .undefinedPrivMethod7()
			else
				return .usedPrivMethod4()
			}
		usedPrivMethod4()
		{
		Print('HI')
		.useMyVal()
		}
			getter_font()
				{
				return #(font: arial, size: 15)
				}
			getter_myVal()
				{
				return 'test'
				}

			useMyVal()
				{
				test = .myVal + 34
				}

			myVal2: 15
			getter_myVal2()
				{
				.get_format()
				x = .get_plugins()
				return 26
				}
			get_plugins()
				{
				return #(5)
				}

			classMember: class
				{
				New()
					{
					.undefinedClassMethod1()
					}
				unusedClassMethod()
					{
					}
				}
				c: Controller
					{
					New()
						{
						.subClassMethod1()
						}
					subClassMethod()
						{
						}
					}
			contrib_foo()
				{
				}
			contrib5_bar()
				{
				}
			Callable: function() { return 'hello world' }
			OtherCallable : function () { return 'hello again' }
	ExtraWhitespace		: function	()	{ return 'tabs' }
			privateCallable: function() { return 'private' }
			}",
		warnings: (
			warnings:(
				[name: "lib:className:11 ERROR: .undefinedPrivMethod1 is undefined"],
				[name: "lib:className:12 ERROR: .undefinedPrivMethod2 is undefined"],
				[name: "lib:className:13 ERROR: .undefinedPrivMethod3 is undefined"],
				[name: "lib:className:14 ERROR: .undefinedPrivMethod4 is undefined"],
				[name: "lib:className:16 ERROR: .undefinedPrivMethod5 is undefined"],
				[name: "lib:className:31 ERROR: .undefinedPrivMethod6 is undefined"],
				[name: "lib:className:42 ERROR: .undefinedPrivMethod7 is undefined"],
				[name: "lib:className:68 ERROR: .get_format is undefined"],
				[name: "lib:className:21 WARNING: .unusedPrivMethod1 is unused"],
				[name: "lib:className:24 WARNING: .unusedPrivMethod2 is unused"],
				[name: "lib:className:27 WARNING: .unusedPrivMethod3 is unused"],
				[name: "lib:className:51 WARNING: .getter_font is unused"],
				[name: "lib:className:66 WARNING: .getter_myVal2 is unused"])
			desc: "There are unused or undefined private methods",
			rating: 0)
		fullWarnings: (
			warnings:(
				[name: ".undefinedPrivMethod1 is undefined"],
				[name: ".undefinedPrivMethod2 is undefined"],
				[name: ".undefinedPrivMethod3 is undefined"],
				[name: ".undefinedPrivMethod4 is undefined"],
				[name: ".undefinedPrivMethod5 is undefined"],
				[name: ".undefinedPrivMethod6 is undefined"],
				[name: ".undefinedPrivMethod7 is undefined"],
				[name: ".get_format is undefined"],
				[name: ".unusedPrivMethod1 is unused"],
				[name: ".unusedPrivMethod2 is unused"],
				[name: ".unusedPrivMethod3 is unused"],
				[name: ".getter_font is unused"],
				[name: ".getter_myVal2 is unused"])
			desc: "There are unused or undefined private methods",
			rating: 0),

		lines: ((11), (12), (13), (14), (16), (31), (42), (68),
			(21), (24), (27), (51), (66))),

		test2: (code: "class
			{
			CallClass(x,y,z)
				{
				a = x + y
				b = y + z
				.undefinedMeth(x,y,z)
				c = .undefinedMeth1 ()+ b
				}
			unusedmeth(x,y,z)
				{
				a = 5
				b = 55
				.undefinedMeth2(x,y,z,a,b)
				return 55
				}

			unusedMeth2(a,b)
				{
				return a + b
				.
				undefinedMeth3()
				}
			PublicMethod2()
				{
				.usedPrivMethod6 = function () { }
				Call(.usedPrivMethod4)
				try .usedPrivMethod5()
				while .usedPrivMethod6()
					throw .usedPrivMethod7()
				}
			usedPrivMethod4()
				{
				}
			usedPrivMethod5()
				{
				}
			usedPrivMethod6()
				{
				}
			usedPrivMethod7()
				{
				}
			}",

		warnings: (
			warnings: (
				[name: "lib:className:7 ERROR: .undefinedMeth is undefined"],
				[name: "lib:className:8 ERROR: .undefinedMeth1 is undefined"],
				[name: "lib:className:14 ERROR: .undefinedMeth2 is undefined"],
				[name: "lib:className:22 ERROR: .undefinedMeth3 is undefined"],
				[name: "lib:className:10 WARNING: .unusedmeth is unused"],
				[name: "lib:className:18 WARNING: .unusedMeth2 is unused"]),
			desc: "There are unused or undefined private methods",
			rating: 0),
		fullWarnings: (
			warnings: (
				[name: ".undefinedMeth is undefined"],
				[name: ".undefinedMeth1 is undefined"],
				[name: ".undefinedMeth2 is undefined"],
				[name: ".undefinedMeth3 is undefined"],
				[name: ".unusedmeth is unused"],
				[name: ".unusedMeth2 is unused"]),
			desc: "There are unused or undefined private methods",
			rating: 0),
		lines: ((7), (8), (14), (22), (10), (18))),

		test_empty_class: (code: "class{ }",
		warnings: (warnings: (), desc: "", rating: 5),
		fullWarnings: (warnings: (),
			desc: "There are no undefined or unused private methods", rating: 5),
		lines: ())

		test_function: (code: "function () { }",
		warnings: (warnings: (), desc: "", rating: 5),
		fullWarnings: (warnings: (),
			desc: "There are no undefined or unused private methods", rating: 5),
		lines: ()),

		test_static_ob: (code:
		"#('CK': (
			newIndex: 0
			desc: 'Supplier Check'
			headerFields: (browse_data, all_browse_data)
			columns: (apivc_invoice_payment, apchklin_currency_default,
				apchklin_amount_due, apchklin_amount_paid_default,
				apchklin_due_date, apchklin_exchgrate_default, apchklin_converted_amount)
			protectField: 'apchklin_protect_supplier'
			renameOb:  (
				apchklin_num: apchklin_num_new,
				apchk_num: apchk_num_hdr,
				apivc_invoice: apivc_invoice_payment,
				suffix: (default: #(apchklin_currency, apchklin_exchgrate)
				)
			buttons: ('Void Check', ('Void Check'), 'Copy and Void')
			))",
		warnings: (warnings: (), desc: "", rating: 5),
		fullWarnings: (warnings: (),
			desc: "There are no undefined or unused private methods", rating: 5),
		lines: ()),

		test_get_and_getter: (code:
		"class
			{
			getter_usedVar()
				{
				return 'this is a test'
				}
			get_anotherUsedVar()
				{
				return 'this is also a test'
				}
			test_fn()
				{
				Print(.usedVar)
				Print(.anotherUsedVar)
				.test_fn() // infinite recursion, don't do this for real
				}
			getter_unusedVar()
				{
				}
			get_alsoUnusedVar()
				{
				}
			}",
		warnings: (warnings: (
				[name: "lib:className:17 WARNING: .getter_unusedVar is unused"],
				[name: "lib:className:20 WARNING: .get_alsoUnusedVar is unused"]),
			desc: "There are unused or undefined private methods", rating: 0),
		fullWarnings: (warnings: (
				[name: ".getter_unusedVar is unused"],
				[name: ".get_alsoUnusedVar is unused"]),
			desc: "There are unused or undefined private methods", rating: 0),
		lines: ((17), (20)))
	)

	Test_Qccm_PrivMethods()
		{
		recordData = Record(recordName: "className", lib: "lib")
		for test in .testCases
			{
			recordData.code = test.code

			fullWarnings = Qc_PrivMethods(recordData, minimizeOutput?: false)
			lineWarnings = fullWarnings.Extract('lineWarnings')
			Assert(fullWarnings is: test.fullWarnings)
			Assert(lineWarnings isSize: 0)

			warnings = Qc_PrivMethods(recordData, minimizeOutput?:)
			lineWarnings = warnings.Extract('lineWarnings')
			Assert(warnings is: test.warnings)
			Assert(lineWarnings is: test.lines)
			}
		}

	Test_class_is_rule()
		{
		recordData = Record(recordName: "Rule_some_fake_class", lib: "lib",
			code: .testCases.test2.code)
		fullWarnings = Qc_PrivMethods(recordData, minimizeOutput?: false)
		lineWarnings = fullWarnings.Extract('lineWarnings')
		Assert(fullWarnings is: #(warnings: (),
			desc: "There are no undefined or unused private methods", rating: 5))

		warnings = Qc_PrivMethods(recordData, minimizeOutput?:)
		lineWarnings = warnings.Extract('lineWarnings')
		Assert(warnings is: #(warnings: (), desc: "", rating: 5))
		Assert(lineWarnings isSize: 0)
		}
	}