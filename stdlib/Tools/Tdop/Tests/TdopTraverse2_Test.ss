// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		vArray = Object()
		endVArray = Object()
		TdopTraverse2(Tdop('"a"$ 1', type: 'expression'),
			{ |node| vArray.Add(node.Token); true},
			{ |node| endVArray.Add(node.Token); true })
		Assert(vArray is: #(BINARYOP, STRING, CAT, NUMBER))
		Assert(endVArray is: #(STRING, CAT, NUMBER, BINARYOP))

		vArray = Object()
		endVArray = Object()
		TdopTraverse2(Tdop('"a"$ 1', type: 'expression'),
			{ |node| vArray.Add(node.Token) },
			{ |node| endVArray.Add(node.Token) }, reverse:)
		Assert(vArray is: #(BINARYOP, NUMBER, CAT, STRING))
		Assert(endVArray is: #(NUMBER, CAT, STRING, BINARYOP))

		vArray = Object()
		endVArray = Object()
		TdopTraverse2(Tdop('"a"$ 1', type: 'expression'),
			{ |node|
				vArray.Add(node.Token)
				node.Token isnt #BINARYOP },
			{ |node| endVArray.Add(node.Token) }, reverse:)
		Assert(vArray is: #(BINARYOP))
		Assert(endVArray is: #(BINARYOP))
		}
	}