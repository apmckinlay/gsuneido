// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Split()
		{
		m = FixedLength.Split
		Assert(m('', #()) is: [])
		Assert(m('', #(field1: (pos: 0, len: 10))) is: [field1: ''])
		Assert(m("abcdefghijklmnopqrstuvwxyz", #(field2: (pos: 0, len: 10)))
			is: [field2: 'abcdefghij'])
		Assert(m("abcdefghijklmnopqrstuvwxyz", #(field3: (pos: 10, len: 2)))
			is: [field3: 'kl'])
		Assert(m("abcdefghijkl   mnopqrstuvwxyz", #(field3: (pos: 10, len: 5)))
			is: [field3: 'kl   '])
		Assert(m("abcdefghijkl   mnopqrstuvwxyz", #(field3: (pos: 10, len: 5,
			type: 'string' trim:)))	is: [field3: 'kl'])
		Assert(m("abcdefghijklmnopqrstuvwxyz",
			#(fld3: (pos: 2, len: 5) fld4: (pos: 7, len: 2) fld5: (pos: 9, len: 12)))
			is: [fld3: 'cdefg' fld4: 'hi' fld5: 'jklmnopqrstu'])
		Assert(m("20170831", #(date: (pos: 0, len: 8, type: 'date')))
			is: [date: #20170831])
		Assert(m("17/31/08", #(date: (pos: 0, len: 10, type: 'date',
			datefmt: 'yy/dd/MM'))) is: [date: #20170831])
		Assert({m("2a17/31/08", #(date: (pos: 0, len: 10, type: 'date',
			datefmt: 'yy/dd/MM')))} throws: 'invalid date')
		Assert(m("123456789", #(number: (pos: 0, len: 9, precision: 4, type: 'number')))
			is: [number: 12345.6789])
		Assert(m("52.61", #(number: (pos: 0, len: 9, type: 'number')))
			is: [number: 52.61])
		}

	Test_Build()
		{
		m = FixedLength.Build
		Assert(m([], #()) is: '')
		Assert(m([], #(fld1: (pos: 0, len: 10))) is: '          ')
		Assert(m([], #(fld1: (pos: 0, len: 10, padChar: 'x'))) is: 'xxxxxxxxxx')
		Assert(m([fld1: 'abc'], #(fld1: (pos: 0, len: 10, padChar: 'x')))
			is: 'abcxxxxxxx')
		Assert(m([fld1: 'abc'],
			#(fld1: (pos: 0, len: 10, padChar: 'x', justify: 'right')))
			is: 'xxxxxxxabc')
		Assert({ m([fld1: 'abc'],
			#(fld1: (pos: 0, len: 10, padChar: 'x', justify: 'invalid'))) }
			throws: 'FixedLength: invalid justify in map')
		Assert(m([fld1: 'abc', fld2: 'defg', fld3: 11111],
			#(fld1: (pos: 0, len: 10, padChar: 'x', justify: 'right')
				fld2: (pos: 10, len: 5, padChar: 'x', justify: 'left')
				fld3: (pos: 15, len: 7, padChar: '0', justify: 'right')
				fld4: (pos: 22, len: 2, padChar: 'Y', justify: 'right')))
			is: 'xxxxxxxabcdefgx0011111YY')
		}
	}