// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cases: (
		((), ""),
		((a: 123), "a: 123\n"),
		((a: 123), "--- \na: 123\n"),
		((a: 123, b: fred), 'a: 123\nb: "fred"\n'),
		((a: "'lisp"), 'a: "\'lisp"\n'),
		((a: (b: 123)), "a: \n    b: 123\n")
		)

	Test_FromObject()
		{
		for c in .cases
			Assert(Yaml.FromObject(c[0]) is: c[1].Replace("^--- \n", ""))
		}

	Test_FromObject_nested()
		{
		ob = #(a: 123, b: testyaml,
			c: (
				c1: (field: testfield),
				c2: (type: date)
				),
			d: yes)
		Assert(Yaml.FromObject(ob),
			is: "a: 123\n" $ 'b: "testyaml"\n' $ "c: \n" $ "    c1: \n" $
				'        field: "testfield"\n' $ "    c2: \n" $ '        type: "date"\n' $
				'd: "yes"\n')
		}

	Test_FromObject_flow()
		{
		Assert(Yaml.FromObject([1, 2, 3]) is: "[1, 2, 3]\n")
		Assert(Yaml.FromObject(['x', 'y']) is: "['x', 'y']\n")
		Assert(Yaml.FromObject([1, [2, 3], Object(b: 4)])
			is: "[1, [2, 3], {'b': 4}]\n")
		Assert(Yaml.FromObject(Object()) is: "")
		Assert(Yaml.FromObject([1, 2, 3], "    ")
			is: "    [1, 2, 3]\n")
		Assert(Yaml.FromObject(Object(x: 1, y: [2, 3]))
			is: "x: 1\ny: [2, 3]\n")
		Assert(Yaml.FromObject(Object(a: Object())) is: "a: []\n")
		Assert(Yaml.FromObject(Object(a: ['x', 'y'], b: [1, Object(z: 2)]))
			is: "a: ['x', 'y']\nb: [1, {'z': 2}]\n")
		}

	Test_flow_scalar()
		{
		f = Yaml.Yaml_flow
		Assert(f(#fred) is: '"fred"')
		Assert(f("") is: '""')
		Assert(f('a "b" c') is: `'a "b" c'`)
		Assert(f(123) is: "123")
		Assert(f(0) is: '0')
		Assert(f(1.5) is: "1.5")
		Assert(f(-1.5) is: "-1.5")
		Assert(f(true) is: "true")
		Assert(f(false) is: "false")
		Assert(f(Date(#20200101)) is: "#20200101")
		}

	Test_flow_empty()
		{
		f = Yaml.Yaml_flow
		Assert(f(Object()) is: "[]")
		Assert(f(Object(a: Object())) is: "{'a': []}")
		ob = Object(a: 1)
		ob.Delete('a')
		Assert(f(ob) is: "[]")
		}

	Test_flow_sequence()
		{
		f = Yaml.Yaml_flow
		Assert(f(['x', 'y']) is: "['x', 'y']")
		Assert(f([1, [2, [3, 4]]]) is: "[1, [2, [3, 4]]]")
		Assert(f([Date(#20200101)]) is: "[#20200101]")
		}

	Test_flow_mapping()
		{
		f = Yaml.Yaml_flow
		Assert(f(Object(a: 1, b: 2)) is: "{'a': 1, 'b': 2}")
		Assert(f(Object(b: 2, a: 1)) is: "{'a': 1, 'b': 2}")
		Assert(f(Object(2: 'b', 1: 'a')) is: "{1: 'a', 2: 'b'}")
		Assert(f([0, 1, a: 2]) is: "{0: 0, 1: 1, 'a': 2}")
		Assert(f(Object(a: -1.5, b: 10_000_000_000))
			is: "{'a': -1.5, 'b': 10000000000}")
		}

	Test_flow_nested()
		{
		f = Yaml.Yaml_flow
		Assert(f(Object(a: Object(c: 1), b: Object(d: Object(e: 2))))
			is: "{'a': {'c': 1}, 'b': {'d': {'e': 2}}}")
		}
	}
