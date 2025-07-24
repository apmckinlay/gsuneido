// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		vArray = Object()
		TdopTraverse(Tdop('"a" $ "b" $ 1', type: 'expression'),
			{ |node| vArray.Add(node.Token) })
		Assert(vArray is: #(BINARYOP, BINARYOP, STRING, CAT, STRING, CAT, NUMBER))
		}
	}