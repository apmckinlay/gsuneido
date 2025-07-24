// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Init()
		{
		r = new .test_refactor
		Assert(r.Init(Record(select: #(cpMin: 0, cpMax: 0))) is: false)
		Assert(r.Msg is: 'Please select the code you want to extract into a method')

		text = "class { F() { if (x is false) return false }"
		Assert(r.Init(Record(select: Object(cpMin: 0, cpMax: text.Size()), :text))
			is: true)
		Assert(r.Msg has: 'selection contains')

		text = "function () { stuff }"
		Assert(r.Init(Record(select: Object(cpMin: 0, cpMax: text.Size()), :text))
			is: false)
		Assert(r.Msg is: 'Extract method can only be used on a class')
		}
	Test_Errors()
		{
		r = new .test_refactor
		text = "class { Fred() { } }"
		Assert(r.Init(data = Record(select: Object(cpMin: 0, cpMax: 1), :text))
			is: true)

		data.method_name = ''
		Assert(r.Errors(data) is: "Invalid method name")

		data.method_name = 'Fred'
		Assert(r.Errors(data) is: "Method name already exists")
		}
	test_refactor: Refactor_Extract_Method
		{
		Msg: false
		Info(msg)
			{ .Msg = msg }
		Warn(msg)
			{ .Msg = msg }
		Refactor_Extract_Method_resetTimer() { }
		Refactor_Extract_Method_initPreview() { }
		Refactor_Extract_Method_updateList() { }
		}
	text:
"class
	{
	One(a, b, c)
		{
		d = a + b
		e = b + c
		A()
		f = 123 + 456
		return e
		}
	Two()
		{
		}
	}"
	Test_Extract()
		{
		// no inputs, no outputs
		selection = '\t\tA()\r\n'
		pos = .text.Find(selection)
		result = .text.
			Replace('\tTwo', '\tAdded()\r\n\t\t{\r\n' $ selection $ '\t\t}\r\n\tTwo').
			Replace('(?q)A()', '.Added()', 1)
		Assert(Refactor_Extract_Method.Extract(.text, '', pos, selection, 'Added')
			is: result)

		// inputs, one output
		selection = '\t\te = b + c\r\n'
		pos = .text.Find(selection)
		result = .text.
			Replace('\tTwo', '\tAdded(b, c)\r\n\t\t{\r\n' $ selection $
				'\t\treturn e\r\n\t\t}\r\n\tTwo').
			Replace('(?q)e = b + c', 'e = .Added(b, c)', 1)
		Assert(Refactor_Extract_Method.Extract(.text,'b, c', pos, selection, 'Added')
			is: result)

		// inputs, one output
		selection = '\t\td = a + b\r\n\t\te = b + c\r\n'
		pos = .text.Find(selection)
		result = .text.
			Replace('\tTwo', '\tAdded(a, b, c)\r\n\t\t{\r\n' $ selection $
				'\t\treturn e\r\n\t\t}\r\n\tTwo').
			Replace('(?q)d = a + b\r\n\t\te = b + c',
				'e = .Added(a, b, c)', 1)
		Assert(Refactor_Extract_Method.Extract(.text, 'a, b, c', pos, selection, 'Added')
			is: result)

		// no inputs, no outputs
		selection = '\t\tf = 123 + 456\r\n'
		pos = .text.Find(selection)
		result = .text.
			Replace('\tTwo', '\tAdded()\r\n\t\t{\r\n' $ selection $
				'\t\t}\r\n\tTwo').
			Replace('(?q)f = 123 + 456', '.Added()', 1)
		Assert(Refactor_Extract_Method.Extract(.text, '', pos, selection, 'Added')
			is: result)

// TODO: multiple outputs
		}
	text2:
"class
	{
	One(a, b)
		{
		c = 1
		d = 3
		++b
		a $= 'x'
		return d
		}
	}"
	Test_Inputs()
		{
		pos = .text2.Find('d')
		selection = .text2[pos ..].BeforeFirst('\t\ta')
		Assert(Refactor_Extract_Method.Inputs(.text2, pos, selection) is: #(b))
		}
	Test_Outputs()
		{
		pos = .text2.Find('d')
		selection = .text2[pos ..].BeforeFirst('\t\ta')
		Assert(Refactor_Extract_Method.Outputs(.text2, pos, selection) is: #(d))
		}
	}
