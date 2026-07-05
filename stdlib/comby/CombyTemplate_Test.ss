// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		// --- hole basics
		// single hole
		Assert(CombyTemplate(':[x]') is: Object(
			Object(type: #HOLE, value: 'x', text: ':[x]', start: 0, end: 4)))

		// two holes with literal text
		Assert(CombyTemplate('foo(:[first], :[second])') is: Object(
			Object(type: #IDENTIFIER, keyword?: false, value: 'foo', text: 'foo',
				start: 0, end: 3),
			Object(type: '', value: '(', text: '(', start: 3, end: 4),
			Object(type: #HOLE, value: 'first', text: ':[first]', start: 4, end: 12),
			Object(type: '', value: ',', text: ',', start: 12, end: 13),
			Object(type: #WHITESPACE, value: ' ', text: ' ', start: 13, end: 14),
			Object(type: #HOLE, value: 'second', text: ':[second]', start: 14, end: 23),
			Object(type: '', value: ')', text: ')', start: 23, end: 24)))

		// holes at start and end with operators
		Assert(CombyTemplate(':[x] + :[y]') is: Object(
			Object(type: #HOLE, value: 'x', text: ':[x]', start: 0, end: 4),
			Object(type: #WHITESPACE, value: ' ', text: ' ', start: 4, end: 5),
			Object(type: '', value: '+', text: '+', start: 5, end: 6),
			Object(type: #WHITESPACE, value: ' ', text: ' ', start: 6, end: 7),
			Object(type: #HOLE, value: 'y', text: ':[y]', start: 7, end: 11)))

		// --- non-hole edge cases
		for item in CombyTemplate('a:[x')
			Assert(item.type isnt #HOLE)

		for item in CombyTemplate(':[ ]')
			Assert(item.type isnt #HOLE)

		for item in CombyTemplate(':[]')
			Assert(item.type isnt #HOLE)

		for item in CombyTemplate(':[!]')
			Assert(item.type isnt #HOLE)

		// --- simple literal
		Assert(CombyTemplate('hello') is: Object(
			Object(type: #IDENTIFIER, keyword?: false, value: 'hello', text: 'hello',
				start: 0, end: 5)))

		// --- whitespace coalescing
		Assert(CombyTemplate('foo  x') is: Object(
			Object(type: #IDENTIFIER, keyword?: false, value: 'foo', text: 'foo',
				start: 0, end: 3),
			Object(type: #WHITESPACE, value: '  ', text: '  ', start: 3, end: 5),
			Object(type: #IDENTIFIER, keyword?: false, value: 'x', text: 'x',
				start: 5, end: 6)))

		Assert(CombyTemplate('a\r\nb') is: Object(
			Object(type: #IDENTIFIER, keyword?: false, value: 'a', text: 'a',
				 start: 0, end: 1),
			Object(type: #WHITESPACE, value: '\r\n', text: '\r\n', start: 1, end: 3),
			Object(type: #IDENTIFIER, keyword?: false, value: 'b', text: 'b',
				start: 3, end: 4)))

		Assert(CombyTemplate('a\r\n:[x]	b') is: Object(
			Object(type: #IDENTIFIER, keyword?: false, value: 'a', text: 'a',
				start: 0, end: 1),
			Object(type: #WHITESPACE, value: '\r\n', text: '\r\n', start: 1, end: 3),
			Object(type: #HOLE, value: 'x', text: ':[x]', start: 3, end: 7),
			Object(type: #WHITESPACE, value: '	', text: '	', start: 7, end: 8),
			Object(type: #IDENTIFIER, keyword?: false, value: 'b', text: 'b',
				start: 8, end: 9)))

		// leading whitespace preserved
		Assert(CombyTemplate('  foo') is: Object(
			Object(type: #WHITESPACE, value: '  ', text: '  ', start: 0, end: 2),
			Object(type: #IDENTIFIER, keyword?: false, value: 'foo', text: 'foo',
				start: 2, end: 5)))

		// trailing whitespace
		Assert(CombyTemplate('foo  ') is: Object(
			Object(type: #IDENTIFIER, keyword?: false, value: 'foo', text: 'foo',
				start: 0, end: 3),
			Object(type: #WHITESPACE, value: '  ', text: '  ', start: 3, end: 5)))

		// trailing whitespace after hole dropped
		Assert(CombyTemplate(':[x]  ') is: Object(
			Object(type: #HOLE, value: 'x', text: ':[x]', start: 0, end: 4),
			Object(type: #WHITESPACE, value: '  ', text: '  ', start: 4, end: 6)))

		// --- comment tokens
		Assert(CombyTemplate('/* comment */') is: Object(
			Object(type: #COMMENT, value: '/* comment */', text: '/* comment */',
				start: 0, end: 13)))

		Assert(CombyTemplate('/* comment */ foo') is: Object(
			Object(type: #COMMENT, value: '/* comment */', text: '/* comment */',
				start: 0, end: 13),
			Object(type: #WHITESPACE, value: ' ', text: ' ', start: 13, end: 14),
			Object(type: #IDENTIFIER, keyword?: false, value: 'foo', text: 'foo',
				start: 14, end: 17)))

		// --- empty and whitespace-only
		Assert(CombyTemplate('') is: Object())
		Assert(CombyTemplate(' \r\n  ')
			is: Object(Object(type: #WHITESPACE, value: ' \r\n  ', text: ' \r\n  ',
				start: 0, end: 5)))
		}
	}
