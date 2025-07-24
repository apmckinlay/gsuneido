// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		// test invalid
		Assert(DateParamInRange("", #19000101, #19000102) is: false)
		param = Object(operation: "equals", value: #19000101, value2: '')
		Assert(DateParamInRange(param, '', #19000102) is: false)
		Assert(DateParamInRange(param, #19000102, '') is: false)
		param.operation = ''
		Assert(DateParamInRange(param, #19000101, #19000102) is: false)
		param.operation = 'test'
		Assert({ DateParamInRange(param, #19000101, #19000102) }
			throws: "unhandled operation")

		// test valid
		param = Object(operation: "equals", value: #19000101, value2: '')
		Assert(DateParamInRange(param, #19000101, #19000102))
		Assert(DateParamInRange(param, #19000102, #19000102) is: false)

		param.operation = 'empty'
		Assert(DateParamInRange(param, #19000101, #19000102) is: false)

		param.operation = 'not empty'
		Assert(DateParamInRange(param, #19000101, #19000102))

		param.operation = 'greater than'
		Assert(DateParamInRange(param, #19000101, #19000102))
		Assert(DateParamInRange(param, #19000101, #19000101) is: false)

		param.operation = 'less than'
		param.value = #19000115
		Assert(DateParamInRange(param, #19000117, #19000119) is: false)
		Assert(DateParamInRange(param, #19000101, #19000105))

		param.operation = 'greater than or equal to'
		Assert(DateParamInRange(param, #19000110, #19000114) is: false)
		Assert(DateParamInRange(param, #19000115, #19000115))

		param.operation = 'less than or equal to'
		Assert(DateParamInRange(param, #19000117, #19000119) is: false)
		Assert(DateParamInRange(param, #19000101, #19000115))

		param.operation = 'in list'
		param.value = #(#19000110 #19000113)
		Assert(DateParamInRange(param, #19000110, #19000111))
		Assert(DateParamInRange(param, #19000112, #19000114))
		Assert(DateParamInRange(param, #19000101, #19000105) is: false)
		Assert(DateParamInRange(param, #19000115, #19000115) is: false)

		param.operation = 'range'
		param.value = #19000110
		param.value2 = #19000120
		Assert(DateParamInRange(param, #19000101, #19000105) is: false)
		Assert(DateParamInRange(param, #19000115, #19000115))
		Assert(DateParamInRange(param, #19000119, #19000121))
		Assert(DateParamInRange(param, #19000120, #19000121))
		Assert(DateParamInRange(param, #19000109, #19000110))
		Assert(DateParamInRange(param, #19000121, #19000123) is: false)
		Assert(DateParamInRange(param, #19000105, #19000123))
		}

	Test_not_equal_to()
		{
		param = Object(operation: "not equal to", value: #20110114, value2: '')
		Assert(DateParamInRange(param, #20110101, #20110131))
		Assert(DateParamInRange(param, #20110114, #20110114) is: false)
		}

	Test_not_in_list()
		{
		param = Object(operation: "not in list",
			value: #(#20110112, #20110114), value2: "")
		Assert(DateParamInRange(param, #20110101, #20110131))
		Assert(DateParamInRange(param, #20110112, #20110112) is: false)
		Assert(DateParamInRange(param, #20110114, #20110114) is: false)

		param = Object(operation: "not in list",
			value: #(#20110112, #20110113, #20110114), value2: "")
		Assert(DateParamInRange(param, #20110112, #20110114) is: false)
		Assert(DateParamInRange(param, #20110112, #20110112) is: false)
		Assert(DateParamInRange(param, #20110114, #20110114) is: false)
		Assert(DateParamInRange(param, #20110120, #20110120))
		}

	Test_not_in_range()
		{
		param = Object(operation: "not in range",
			value: #20110110, value2: #20110115)

		Assert(DateParamInRange(param, #20110112, #20110114) is: false)
		Assert(DateParamInRange(param, #20110110, #20110115) is: false)
		Assert(DateParamInRange(param, #20110101, #20110110))
		Assert(DateParamInRange(param, #20110110, #20110113) is: false)
		Assert(DateParamInRange(param, #20110110, #20110110) is: false)
		Assert(DateParamInRange(param, #20110113, #20110115) is: false)
		Assert(DateParamInRange(param, #20110115, #20110115) is: false)
		Assert(DateParamInRange(param, #20110115, #20110125))
		Assert(DateParamInRange(param, #20110101, #20110105))
		Assert(DateParamInRange(param, #20110120, #20110125))
		Assert(DateParamInRange(param, #20110101, #20110113))
		Assert(DateParamInRange(param, #20110113, #20110125))
		Assert(DateParamInRange(param, #20110101, #20110125))
		}
	}
