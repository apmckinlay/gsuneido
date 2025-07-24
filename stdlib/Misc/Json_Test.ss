// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cases: (
		(true, 'true')
		(false, 'false')
		(123, '123')
		(-1e99, '-1e99')
		(.5, '0.5')
		(-.5, '-0.5')
		('foo bar', '"foo bar"')
		((), '[]')
		((1, 2, 3), '[1,2,3]')
		((1, 2, (3, 4)), '[1,2,[3,4]]')
		((a: 1, b: '\\'), '{"a":1,"b":"\\\\"}')
		((a: 1, b: '"'), '{"a":1,"b":"\\""}')
		((a: 1, b: '~!@#$%^&*()_+?:{}[]|,.'), '{"a":1,"b":"~!@#$%^&*()_+?:{}[]|,."}')
		((b: 'abc'), '{"b":"abc"}')
		((a: 1, b: 'abc'), '{"a":1,"b":"abc"}')
		((a: 0.850817, b: 'abc'), '{"a":0.850817,"b":"abc"}')
		((a: true, b: 'abc'), '{"a":true,"b":"abc"}')
		((a: false, b: 'abc'), '{"a":false,"b":"abc"}')
		((a: 1, b: #(c: 33, d: 'dddd')), '{"a":1,"b":{"c":33,"d":"dddd"}}')
		((a: 1, b: #('abc': 222, 'efg': '2009')), '{"a":1,"b":{"abc":222,"efg":"2009"}}')
		((a: 'Test with multiple lines and tabs:\r\n' $
				'line 1\r\n\r\n\tline 2 with tab\r\nline 3', b: 'hello world'),
			'{"a":"Test with multiple lines and tabs:' $
				'\\r\\nline 1\\r\\n\\r\\n\\t' $
				'line 2 with tab\\r\\nline 3","b":"hello world"}')
		(#('string{}[]', #( a: 1)), '["string{}[]",{"a":1}]')
		((a: .5, b: -.5, c: 1.5, d: 0, e: -1.5),
			'{"a":0.5,"b":-0.5,"c":1.5,"d":0,"e":-1.5}')
		("hello", '"hello"')
		(123, '123')
		(true, 'true')
		(false, 'false')
		)
	Test_Encode()
		{
		for c in .cases
			Assert(Json.Encode(c[0]) is: c[1])
		}

	Test_Decode()
		{
		for c in .cases
			Assert(Json.Decode(c[1]) is: c[0])

		Assert(Json.Decode('{}') is: [])

		// test decoding non-ascii characters in \uXXXX format
		s = '{"desc": "A\u0026B \u2013 hello world","name":"joe smith"}'
		Assert(Json.Decode(s)
			is: Object(
				desc: 'A&B \xE2\x80\x93 hello world'.FromUtf8(),
				name: 'joe smith'))

		Assert({Json.Decode('null')}
			throws: 'Invalid Json format: data should not contain null')
		Assert(Json.Decode('null', handleNull: 'empty') is: '')
		Assert({Json.Decode('null', handleNull: 'skip')}
			throws: 'Invalid Json format: data should not contain null')

		Assert({Json.Decode('{"a": null}')}
			throws: 'Invalid Json format: data should not contain null')
		Assert(Json.Decode('{"a": null, "b": [1, null]}', handleNull: 'empty')
			is: #(a: '', b: #(1, '')))
		Assert(Json.Decode('{"a": null, "b": 1, "c": [1, null, 2]}', handleNull: 'skip')
			is: #(b: 1, c: #(1, 2)))

		Assert(Json.Decode('[ null ]', 'skip') is: #())
		Assert(Json.Decode('[ null ]', 'empty') is: #(''))
		Assert({Json.Decode('[ null ]')}
			throws: 'Invalid Json format: data should not contain null')

		Assert(Json.Decode('[ null, 1, 2]', 'skip') is: #(1, 2))
		Assert(Json.Decode('[ null, 1, 2]', 'empty') is: #('', 1, 2))
		Assert({Json.Decode('[ null, 1, 2]')}
			throws: 'Invalid Json format: data should not contain null')
		}

	valueListCases: (
		((3, 8, 1), '[3,8,1]')
		((a: 1, b: 'b', c: ((x: 13, y: 23, z: 'z33'), (x: 23, y: 33, z: '43')), d: 4),
		'{"a":1,"b":"b","c":[{"x":13,"y":23,"z":"z33"},{"x":23,"y":33,"z":"43"}],"d":4}')
		)
	Test_valuesList()
		{
		for c in .valueListCases
			{
			Assert(Json.Encode(c[0]) is: c[1])
			Assert(Json.Decode(c[1]) is: c[0])
			}
		Assert({ Json.Decode(1234) } throws: 'Invalid Json format')
		Assert({ Json.Decode('hello world') } throws: 'Invalid Json format')
		}

	Test_decode_extraText()
		{
		Assert(Json.Decode(' {"abc": "eee"} ') is: #(abc: 'eee'))
		Assert(Json.Decode('\r\n{"abc": "eee"} ') is: #(abc: 'eee'))
		Assert({ Json.Decode('extra {"abc": "eee"} ') }
			throws: 'Invalid Json format: unexpected: extra')

		Assert(Json.Decode('{"abc": "eee"} ') is: #(abc: 'eee'))
		Assert(Json.Decode('{"abc": "eee", "ddd": 20 }\r\n') is: #(abc: 'eee', ddd: 20))
		Assert({ Json.Decode('{"abc": "eee"} extra') }
			throws: 'Invalid Json format: extra text at end')
		}
	}
