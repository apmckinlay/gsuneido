// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(TdopIsExpression(Tdop('function(){}', type: 'expression')))
		Assert(TdopIsExpression(Tdop('class{}', type: 'expression')))

		Assert(TdopIsExpression(Tdop('1', type: 'expression')))
		Assert(TdopIsExpression(Tdop('+1', type: 'expression')))
		Assert(TdopIsExpression(Tdop('-1', type: 'expression')))
		Assert(TdopIsExpression(Tdop('--1', type: 'expression')))

		Assert(TdopIsExpression(Tdop('#20170816', type: 'expression')))
		Assert(TdopIsExpression(Tdop('true', type: 'expression')))
		Assert(TdopIsExpression(Tdop('false', type: 'expression')))

		Assert(TdopIsExpression(Tdop('#()', type: 'expression')))
		Assert(TdopIsExpression(Tdop('#{}', type: 'expression')))

		Assert(TdopIsExpression(Tdop('"abc"', type: 'expression')))
		Assert(TdopIsExpression(Tdop('"abc" $ "efg" $ "h"', type: 'expression')))
		Assert(TdopIsExpression(Tdop('"abc" $ "efg" $ h', type: 'expression')))

		Assert(TdopIsExpression(Tdop('return 1', type: 'statement')) is: false)
		}
	}