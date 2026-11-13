// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cases: (
		((), ""),
		((a: 123), "a: 123\n"),
		((a: 123), "--- \na: 123\n"),
		((a: 123, b: fred), 'a: 123\nb: "fred"\n'),
		((a: "'lisp"), 'a: "\'lisp"\n'),
		((a: (b: 123)), "a:\n  b: 123\n")
		)

	Test_FromObject()
		{
		for c in .cases
			Assert(Yaml.FromObject(c[0]) is: c[1].Replace("^--- \n", ""))
		}

	Test_FromObject_numeric_string_keys()
		{
		Assert(Yaml.FromObject(Object("123": #value))
			is: '"123": "value"\n')
		Assert(Yaml.FromObject(Object("123": #value, a: 1, "456": 2))
			is: '"123": "value"\n"456": 2\na: 1\n')
		Assert(Yaml.FromObject(Object(nested: Object("99": #val)))
			is: "nested:\n  \"99\": \"val\"\n")
		Assert(Yaml.FromObject(Object("1.2": 1, "-1": 2, "+3": 3, "1e5": 4, ".5": 5))
			is: '"+3": 3\n"-1": 2\n".5": 5\n"1.2": 1\n"1e5": 4\n')
		Assert(Yaml.FromObject(Object("1a": 1, a1: 2))
			is: "1a: 1\na1: 2\n")
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
			is: "a: 123\n" $ 'b: "testyaml"\n' $ "c:\n" $ "  c1:\n" $
				'    field: "testfield"\n' $ "  c2:\n" $ '    type: "date"\n' $
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

	Test_FromObject_block()
		{
		Assert(Yaml.FromObject([1, 2, 3], blockStyle?:)
			is: "- 1\n- 2\n- 3\n")
		Assert(Yaml.FromObject(['x', 'y'], blockStyle?:)
			is: "- x\n- y\n")
		Assert(Yaml.FromObject([1, [2, 3], Object(b: 4)], blockStyle?:)
			is: "- 1\n- - 2\n  - 3\n- b: 4\n")
		Assert(Yaml.FromObject(Object(), blockStyle?:) is: "")
		Assert(Yaml.FromObject([1, 2, 3], "    ", blockStyle?:)
			is: "    - 1\n    - 2\n    - 3\n")
		Assert(Yaml.FromObject(Object(x: 1, y: [2, 3]), blockStyle?:)
			is: "x: 1\ny:\n  - 2\n  - 3\n")
		Assert(Yaml.FromObject(Object(a: Object()), blockStyle?:)
			is: "a: []\n")
		Assert(Yaml.FromObject(Object(a: ['x', 'y'], b: [1, Object(z: 2)]),
				blockStyle?:)
			is: "a:\n  - x\n  - y\nb:\n  - 1\n  - z: 2\n")
		}

	Test_FromObject_block_empty()
		{
		Assert(Yaml.FromObject([Object(), [Object()]], blockStyle?:)
			is: "- []\n- - []\n")
		Assert(Yaml.FromObject(Object(), "    ", blockStyle?:)
			is: "    []\n")
		}

	Test_FromObject_block_scalars()
		{
		Assert(Yaml.FromObject(
				[#fred, "", 'a "b" c', 123, 0, 1.5, -1.5,
					true, false, #20200101], blockStyle?:)
			is: '- fred\n- \n- a "b" c' $
				"\n- 123\n- 0\n- 1.5\n- -1.5\n- true\n- false\n- #20200101\n")
		for value in #("true", FALSE, yes, null, "123", "1.5",
			"a: b", "a #comment", "", " leading", "trailing ", "- item", "[1]")
			{
			Assert(Yaml.FromObject([value], blockStyle?:)
				is: "- " $ value $ '\n')
			Assert(Yaml.FromObject(Object(text: value), blockStyle?:)
				is: "text: " $ value $ '\n')
			}
		}

	Test_FromObject_block_hash_strings()
		{
		for value in #('#', "#/components/schemas/Item", '#text "quoted"')
			{
			Assert(Yaml.FromObject([value], blockStyle?:)
				is: "- " $ Display(value) $ '\n')
			Assert(Yaml.FromObject(Object(ref: value), blockStyle?:)
				is: "ref: " $ Display(value) $ '\n')
			}
		Assert(Yaml.FromObject(["text#hash"], blockStyle?:)
			is: "- text#hash\n")
		}

	Test_FromObject_block_mappings()
		{
		Assert(Yaml.FromObject([Object(b: 2, a: 1), Object("123": #value)],
				blockStyle?:)
			is: '- a: 1\n  b: 2\n- "123": value\n')
		Assert(Yaml.FromObject([[0, 1, a: 2]], blockStyle?:)
			is: "- 0: 0\n  1: 1\n  a: 2\n")
		Assert(Yaml.FromObject([Object(a: Object(c: [1, [2]]))],
				blockStyle?:)
			is: "- a:\n    c:\n      - 1\n      - - 2\n")
		}

	Test_FromObject_block_servers()
		{
		ob = Object(servers: [Object(description: "ETA REST API",
				url: "http://127.0.0.1:3247/EtaRestAPI/")])
		Assert(Yaml.FromObject(ob, blockStyle?:)
			is: "servers:\n  - description: ETA REST API\n" $
				"    url: http://127.0.0.1:3247/EtaRestAPI/\n")
		}
	}
