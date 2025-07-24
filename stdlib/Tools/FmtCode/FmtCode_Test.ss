// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
// BuiltDate > 30000101
Test
	{
	Test_format()
		{
		.t('function(a, b = "foo", c = #())
				{
				}')
		.t("function(unused)
				{
				}")
		.t("function(@unused)
				{
				}")
		.t("function(unused /*unused*/) { }",
			"function(unused)
				{
				}")
		.t("function(x /*unused*/)
				{
				}")
		.t("function(x /*unused*/ = 123)
				{
				}")
		.t("function(f = function() { })
				{
				}")
		.t('function
			(
			)
			{
			}',
			`function()
				{
				}`)
		.t("function()
				{
				function() { function() { 123 } }
				}")
		.t("123")
		.t("foo")
		.t(`s = "\t\n"`) // escape
		.t("return false")
		.t("-x")
		.t("not x")
		.t("++x")
		.t("x--", "--x")
		.t("return x--")
		.t("return .x")
		.t("x = y = 123")
		.t("x = #foo $ f(#bar)")
		.t('s =~ "^x"')
		.t("a * b * c")
		.t("a + b - c + d")
		.t("x - 1")
		.t("x + false")
		.t("a * b / c")
		.t("x in (1, 2, #foo, #(9))")
		.t("x not in (1, 2, #foo, #(9))")
		.t("x ? y : z")
		.t("x.y")
		.t("x.y + z")
		.t("x.y.z")
		.t("x[2]")
		.t("x[2..-1]")
		.t("x[1..]")
		.t("x[..n]")
		.t("x[f() .. to()]")
		.t("x[i::9]")
		.t("x[1::]")
		.t("x[::n]")
		.t("x = #(1, (2), foo)")
		.t("f()")
		.t("f = function() { f();; }")
		.t("f(1, 2, 3)")
		.t("f(1, 2, a: 3)")
		.t("f(:a, :b)")
		.t("f(a: a, b: b)", "f(:a, :b)")
		.t('f("a-b": 123)')
		.t("f(a:, b:)")
		.t("f(@args)")
		.t("f(@+1args)")
		.t("(.f)()")
		.t(".f()")
		.t("x.y()")
		.t("[]")
		.t("Record(x: 1, y: 2)", "[x: 1, y: 2]")
		.t("Object(1, 2)", "[1, 2]")
		.t("Object(x: 1, y: 2)")
		.t('Object(Foo: Foo)')
		.t('f = function(x = "foo", y = false) { 123 }')
		.t("f = function() { one; two }")
		.t("f = function() { one; two; }", "f = function() { one; two }")
		.t("c = class
				{
				}")
		.t("b = { it + x }")
		.t("b = {|it| it + x }", "b = { it + x }")
		.t("b = {|y| y + x }")
		.t("b = { one; two }")
		.t("x = fn({ block })")
		.t("if fn({ block })
				stmt")
		.t("if (fn() {|x| block })
				stmt",
			"if fn(block: {|x| block })
				stmt")
		.t("QueryApply(query)
				{|x|
				one
				}")
		.t("f();;")
		.t('throw "foo"')
		.t('new this')
		.t('new this()', 'new this')
		.t('new this(123)')
		.t('(new this)()')
		.t('(new this(123))(456)')
		.t("a and (b or c)")

		// try-catch
		.t("try
				f()")
		.t("try
				{
				one
				two
				}")
		.t("try
				one
			catch
				two")
		.t("try
				{
				one
				two
				}
			catch
				{
				three
				four
				}")
		.t("try
				one
			catch(e)
				two")
		.t('try
				one
			catch(e /*unused*/, "pat")
				two')

		.t("if cond
				stmt")
		.t("if ((x = y) is true)
				stmt")
		.t("if cond
				stmt
			else
				stmt")
		.t("if cond
				stmt
			else if cond
				stmt")
		.t("if cond
				{
				stmt1
				stmt2
				}
			else
				stmt")
		.t("if cond
				stmt
			else
				{
				stmt1
				stmt2
				}")
		.t("if cond
				{
				stmt1
				stmt2
				}
			else
				{
				stmt1
				stmt2
				}")
		.t("if cond
				{
				stmt1
				stmt2
				}
			else if cond
				{
				stmt1
				stmt2
				}
			else
				{
				stmt1
				stmt2
				}")
		.t("if cond
				{
				if cond2
					stmt
				}
			else
				stmt")
		.t("switch
				{
				}")
		.t("switch n
				{
			case 0:
				zero
			case 1, 2, 3:
				one
				two
				three
			default:
				panic
				}")
		.t("forever
				stmt")
		.t("while x < y
				++x")
		.t("while i < n
				{
				i += x
				++z
				}")
		.t("while cond
				if x
					break
				else
					continue")
		.t("do
				stmt
				while cond",
			"do
				{
				stmt
				} while cond")
		.t("do
				{
				stmt1
				stmt2
				} while cond")
		.t("for x in y
				f(x)")
		.t("for (;;)
				stmt",
			"forever
				stmt")
		.t("while true
				stmt",
			"forever
				stmt")
		.t("for (i = 0; i < 9; ++i)
				stmt")
		.t("for (i = 0, j = n; i < j; ++i, --j)
				{
				stmt1
				stmt2
				}")
		.t("for (i++; i < 9; i++)
				stmt",
			"for (++i; i < 9; ++i)
				stmt")
		.t("x = function() { }")
		.t("x = function() { 123 }")
		.t("x = function() { }.foo")
		.t("x = function() { 123 }.foo + x")

		.t("#()")
		.t('#(1, 2, a: 3, "@#": 4)')
		.t("#(a:, b:)")
		.t("#(a: true, b: true)", "#(a:, b:)")
		.t("#{x: 1, y: 2}")
		.t("class
				{
				}")
		.t("class
				{
				foo: 123
				bar: false
				}")
		.t("class
				{
				ob: (1, 2, 3)
				}")
		.t("class
				{
				fn(x)
					{
					return #()
					}
				}")
		.t('class
				{
				foo: "Foo Bar"
				"foo bar": 123
				bar()
					{
					}
				Getter_()
					{
					}
				Getter_Public()
					{
					}
				getter_public()
					{
					}
				}')
		}
	Test_comments()
		{
		// constants
		.t('/*foo*/ function()
				{
				} //bar')
		.t('/*foo*/ class
				{
				} //bar')
		.t('/*foo*/ #() //bar')
		.t('/*foo*/ #{} //bar')
		// statements
		.t('/*foo*/ 123 //bar')
		.t('/*foo*/ return //bar')
		.t('1 //
			2')
		.t('1
			/**/ 2')
		.t('/*a*/ 1 //b
			/*c*/ 2 //d
			/*d*/ 3 //e')
		.t('/*a*/ if foo
				/*b*/ { //1
				/*c*/ } //2')
		.t('if foo
				{ //1
				} //2
			else if bar
				{ //3
				} //4')
		.t('if a //
				b
			else //
				c')
		.t('if a
				//
				b
			else
				//
				c')
		.t('forever
				{
				}
			//')
		.t('try //
				foo
			catch //
				bar')
		.t('try
				//
				foo
			catch
				//
				bar')
		// expressions
		.t('/*a1*/ a /*a2*/ + /*b1*/ b /*b2*/')
		.t('x = //
				.y + //
				.z //')
		// arguments
		.t('f(/*1*/ 12 /*2*/, /*3*/ a: 34 /*4*/)')
		// params
		.t('function(/*1*/ a /*2*/, /*3*/ b = 0 /*4*/)
				{
				}')
		.t('function(/*1*/ @abc /*2*/)
				{
				}')
		.t('class
				{
				a: 1 //one
				b: 2 //two
				c: 3 //three
				}')
		.t('class
				{
				/*A*/ a: 1 //one
				/*B*/ b: 2 //two
				/*C*/ c: 3 //three
				}')
		.t('x //1
				? y //2
				: z //3')
		// mid pos
		.t('function /*1*/ (/*2*/) /*3*/
				{ //4
				}'
			'function() /*1*/ /*2*/ /*3*/
				{ //4
				}')
		.t('function() //foo
				{ //bar
				}')
		.t('class //foo
				{ //bar
				}')
		.t('switch /*1*/
				{ //2
			/*3*/ case x: //4
				/*5*/ dostuff //6
			/*7*/ default: //8
				/*9*/ otherwise //10
				} //11',
			'switch /*1*/
				{ //2
			case /*3*/ x: //4
				/*5*/ dostuff //6
			default: /*7*/ //8
				/*9*/ otherwise //10
				} //11')
		.t('try //1
				/*2*/ dostuff //3
			/*4*/ catch //5
				/*6*/ otherwise //7',
			'try //1
				/*2*/ dostuff //3
			catch /*4*/ //5
				/*6*/ otherwise //7')
		.t('try //1
				{ //2
				}
			catch(e) //3
				{ //4
				}')

		.t('foo
			//
			bar')
		.t('foo()
			//
			bar')
		.t('if foo
				//
				bar')
		.t('if foo()
				//
				bar')

		.t('foo //
			bar')
		.t('foo() //
			bar')
		.t('if foo //
				bar')
		.t('if foo() //
				bar')
		}
	Test_line_breaks()
		{
		.t('function
				(
				)
				{

				}',
			'function()
				{
				}')
		// parameters
		.t('function(a,
				b,
				c)
				{
				}')
		// arguments
		.t('f(x,
				y,
				z)')
		.t('f(1
				g(a))',
			'f(1,
				g(a))')
		.t('f(1,
				g(a))')
		.t('f(
				/*A*/ a: 1, //one
				/*B*/ b:, //two
				/*C*/ :c //three
				)')
		.t('f(
				/*A*/ a: 1, //one
				/*B*/ b: /*two*/
				/*C*/ :c //three
				)',
			'f(
				/*A*/ a: 1, //one
				/*B*/ b:, /*two*/
				/*C*/ :c //three
				)')
		.t('[
				/*A*/ a: 1, //one
				/*B*/ b:, //two
				/*C*/ :c //three
				]')
		// members
		.t('#(1, 2,
				a: 3,
				b: 4)')
		.t('#(
				/*A*/ a: 1, //one
				/*B*/ b: 2, //two
				/*C*/ c: //three
				)')
		.t('#{
				/*A*/ a: 1, //one
				/*B*/ b: 2, //two
				/*C*/ c: //three
				}')
		// expressions
		.t('a +
				b')
		.t('a *
				-b')
		.t('a and
				(b or c)')
		.t('a.
				b().
				c()')
		.t('x
				? y : z')
		.t('x
				? y
				: z')
		.t('x ?
				y :
				z',
			'x ? y : z')
		.t('x =
				.y +
				.z')
		// nesting
		.t('#((a,
				b),
				(c,
					d),
				e)')
		.t('f(g(a,
				b),
				h(c,
					d),
				e)')
		}
	Test_blank_lines()
		{
		// statements
		.t('function()
				{

				123

				456

				}',
			'function()
				{
				123

				456
				}')
		.t('function()
				{

				123;

				456;

				}',
			'function()
				{
				123

				456
				}')
		// class members
		.t('class
				{

				a: 1

				b: 2

				}',
			'class
				{
				a: 1

				b: 2
				}')
		.t('class
				{

				a: 1;

				b: 2;

				}',
			'class
				{
				a: 1

				b: 2
				}')
		// object members
		.t('#(

				1

				2

				)',
			'#(
				1,

				2
				)')
		.t('#(

				1,

				2,

				)',
			'#(
				1,

				2
				)')
		.t('#(

				a: 1

				b: 2

				)',
			'#(
				a: 1,

				b: 2
				)')
		.t('#(

				a: 1,

				b: 2,

				)',
			'#(
				a: 1,

				b: 2
				)')
		}
	Test_comment_lines()
		{
		.t('//before
			function()
				{
				123
				//a
				/*a*/ //c
				456
				}
			//after')
		}
	Test_mix()
		{
		// newlines
		.t('// one f
			/* two f */ function(x /* x three */, y /* y four */)
				{
				// five
				x = 1
				F() // six

				Fn(123 /* 123 seven */, /* eight */
					456 /* nine 456 */) // ten

				z = x /* x eleven */ + // twelve
					y // y thirteen
				return // fourteen
				}')
		}
	t(src, expected = false)
		{
		src = src.Tr('\r')
		funcWrap = .first(src) not in ("function", "class", "#")
		if funcWrap
			src = "function()\n\t{\n\t" $ src.Replace("^\t\t") $ "\n\t}\n"
		else
			src = src.Replace("^\t\t\t") $ '\n'
		if expected is false
			expected = src
		else
			{
			expected = expected.Tr('\r')
			if funcWrap
				expected = "function()\n\t{\n\t" $ expected.Replace("^\t\t") $ "\n\t}\n"
			else
				expected = expected.Replace("^\t\t\t") $ '\n'
			}
		actual = FmtCode(src)
		Assert(actual is: expected)
		}
	first(src)
		{
		scan = Scanner(src)
		while scan isnt (tok = scan.Next2()) and
			tok in (#COMMENT, #WHITESPACE, #NEWLINE)
			{}
		return tok is scan ? false : scan.Text()
		}
	Test_args()
		{
		before = "function()\n\t{\n\tf(\n\t\t"
		after = "\n\t\t)\n\t}\n"
		for arg1 in #("1 + 2", "n: 1 + 2", "n:", ":n")
			for arg2 in #("1 + 2", "m: 1 + 2", "m:", ":m")
				{
				if arg1 isnt "1 + 2" and arg2 is "1 + 2"
					continue
				for sep in #(' ', ', ', "\n\t\t", ",\n\t\t")
					{
					src = before $ arg1 $ sep $ arg2 $ after
					// Print(src)
					Assert(FmtCode(src).Tr(',') is: src.Tr(','))
					}
				}
		}
	}