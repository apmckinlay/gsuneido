// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_All()
		{
		Assert(StringParamCompare(Record(), '') is: true)
		Assert(StringParamCompare(#(operation: '', value: '', value2: ''), ''))
		}
	Test_Equals()
		{
		param = Object(operation: 'equals', value: '110', value2: '')
		Assert(StringParamCompare(param, '100') is: false)
		Assert(StringParamCompare(param, '111') is: false)
		Assert(StringParamCompare(param, '110') is: true)
		}
	Test_InList()
		{
		param = Object(operation: "in list",
			value2: "",
			value: #("100", "101", "105"))
		Assert(StringParamCompare(param, '100')  is: true)
		Assert(StringParamCompare(param, '110') is: false)

		param = Object(operation: "not in list",
			value2: "",
			value: #("100", "101", "105"))
		Assert(StringParamCompare(param, '100') is: false)
		Assert(StringParamCompare(param, '110') is: true)
		}
	Test_Range()
		{
		param = Object(operation: 'range', value: '100', value2: '105')
		Assert(StringParamCompare(param, '100') is: true)
		Assert(StringParamCompare(param, '101') is: true)
		Assert(StringParamCompare(param, '105') is: true)
		Assert(StringParamCompare(param, '110') is: false)

		param = Object(operation: 'not in range', value: '100', value2: '105')
		Assert(StringParamCompare(param, '100') is: false)
		Assert(StringParamCompare(param, '101') is: false)
		Assert(StringParamCompare(param, '105') is: false)
		Assert(StringParamCompare(param, '110') is: true)
		}
	Test_Empty()
		{
		param = Object(operation: 'empty', value: '', value2: '')
		Assert(StringParamCompare(param, '') is: true)
		Assert(StringParamCompare(param, '100') is: false)

		param = Object(operation: 'not empty', value: '', value2: '')
		Assert(StringParamCompare(param, '') is: false)
		Assert(StringParamCompare(param, '100') is: true)
		}
	Test_LessThan()
		{
		param = Object(operation: 'less than', value: '110', value2: '')
		Assert(StringParamCompare(param, '100') is: true)
		Assert(StringParamCompare(param, '110') is: false)
		Assert(StringParamCompare(param, '99') is: false)

		param = Object(operation: 'less than or equal to', value: '110', value2: '')
		Assert(StringParamCompare(param, '100') is: true)
		Assert(StringParamCompare(param, '110') is: true)
		Assert(StringParamCompare(param, '99') is: false)
		}
	Test_GreaterThan()
		{
		param = Object(operation: 'greater than', value: '110', value2: '')
		Assert(StringParamCompare(param, '100') is: false)
		Assert(StringParamCompare(param, '110') is: false)
		Assert(StringParamCompare(param, '120') is: true)
		Assert(StringParamCompare(param, '99') is: true)

		param = Object(operation: 'greater than or equal to', value: '110', value2: '')
		Assert(StringParamCompare(param, '100') is: false)
		Assert(StringParamCompare(param, '110') is: true)
		Assert(StringParamCompare(param, '120') is: true)
		Assert(StringParamCompare(param, '99') is: true)
		}
	Test_Contains()
		{
		param = Object(operation: 'contains', value: '110ABC', value2: '')
		Assert(StringParamCompare(param, '2110abc') is: true)
		Assert(StringParamCompare(param, '21105abc') is: false)

		param = Object(operation: 'does not contain', value: '110ABC', value2: '')
		Assert(StringParamCompare(param, '2110abc') is: false)
		Assert(StringParamCompare(param, '21105abc') is: true)
		}
	Test_Matches()
		{
		param = Object(operation: 'matches', value: '110ABC', value2: '')
		Assert(StringParamCompare(param, '2110abc') is: false)
		Assert(StringParamCompare(param, '2110ABC') is: true)
		Assert(StringParamCompare(param, '21105abc') is: false)

		param = Object(operation: 'does not match', value: '110ABC', value2: '')
		Assert(StringParamCompare(param, '2110abc') is: true)
		Assert(StringParamCompare(param, '2110ABC') is: false)
		Assert(StringParamCompare(param, '21105abc') is: true)
		}
	Test_StartEndWith()
		{
		param = Object(operation: 'starts with', value: '125', value2: '')
		Assert(StringParamCompare(param, '2125') is: false)
		Assert(StringParamCompare(param, '1251') is: true)
		Assert(StringParamCompare(param, '1255') is: true)

		param = Object(operation: 'ends with', value: '25', value2: '')
		Assert(StringParamCompare(param, '1234') is: false)
		Assert(StringParamCompare(param, '4535') is: false)
		Assert(StringParamCompare(param, '4225') is: true)
		Assert(StringParamCompare(param, '4325') is: true)
		}
	Test_Unhandled()
		{
		param = Object(operation: 'undhandled', value: '25', value2: '')
		Assert({ StringParamCompare(param, '1234') } throws: 'unhandled')
		}
	}