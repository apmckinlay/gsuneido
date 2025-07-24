// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(TdopIsConstant(Tdop('function(){}', type: 'expression')) is: true)
		Assert(TdopIsConstant(Tdop('class{}', type: 'expression')) is: true)

		Assert(TdopIsConstant(Tdop('1', type: 'expression')) is: true)
		Assert(TdopIsConstant(Tdop('+1', type: 'expression')) is: true)
		Assert(TdopIsConstant(Tdop('-1', type: 'expression')) is: true)
		Assert(TdopIsConstant(Tdop('--1', type: 'expression')) is: false)

		Assert(TdopIsConstant(Tdop('#20170816', type: 'expression')) is: true)
		Assert(TdopIsConstant(Tdop('true', type: 'expression')) is: true)
		Assert(TdopIsConstant(Tdop('false', type: 'expression')) is: true)

		Assert(TdopIsConstant(Tdop('#()', type: 'expression')) is: true)
		Assert(TdopIsConstant(Tdop('#{}', type: 'expression')) is: true)

		Assert(TdopIsConstant(Tdop('"abc"', type: 'expression')) is: true)
		Assert(TdopIsConstant(Tdop('"abc" $ "efg" $ "h"', type: 'expression')) is: true)
		Assert(TdopIsConstant(Tdop('"abc" $ "efg" $ h', type: 'expression')) is: false)
		}
	}