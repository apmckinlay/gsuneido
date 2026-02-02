// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(TdopIsConstant(Tdop('function(){}', type: 'expression')))
		Assert(TdopIsConstant(Tdop('class{}', type: 'expression')))

		Assert(TdopIsConstant(Tdop('1', type: 'expression')))
		Assert(TdopIsConstant(Tdop('+1', type: 'expression')))
		Assert(TdopIsConstant(Tdop('-1', type: 'expression')))
		Assert(TdopIsConstant(Tdop('--1', type: 'expression')) is: false)

		Assert(TdopIsConstant(Tdop('#20170816', type: 'expression')))
		Assert(TdopIsConstant(Tdop('true', type: 'expression')))
		Assert(TdopIsConstant(Tdop('false', type: 'expression')))

		Assert(TdopIsConstant(Tdop('#()', type: 'expression')))
		Assert(TdopIsConstant(Tdop('#{}', type: 'expression')))

		Assert(TdopIsConstant(Tdop('"abc"', type: 'expression')))
		Assert(TdopIsConstant(Tdop('"abc" $ "efg" $ "h"', type: 'expression')))
		Assert(TdopIsConstant(Tdop('"abc" $ "efg" $ h', type: 'expression')) is: false)
		}
	}