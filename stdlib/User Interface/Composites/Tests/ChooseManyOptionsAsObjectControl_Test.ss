// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_formatValues()
		{
		mock = Mock(ChooseManyOptionsAsObjectControl)
		mock.When.number?([anyArgs:]).CallThrough()
		mock.When.object?([anyArgs:]).CallThrough()
		mock.When.formatValues([anyArgs:]).CallThrough()
		mock.origValues = #()
		mock.ChooseManyOptionsAsObjectControl_stringMembers = #()
		mock.formatValues(ob = Object(), [value: '1', option: "companynum"])
		Assert(ob.companynum is: 1)

		mock.origValues = #(companynum: 1)
		mock.formatValues(ob = Object(), [value: '2', option: "companynum"])
		Assert(ob.companynum is: 2)

		mock.origValues = #()
		mock.formatValues(ob = Object(), [value: '#(3147,3149)', option: "ports"])
		Assert(ob.ports is: #(3147, 3149))

		mock.origValues = #(ports: #(3147))
		mock.formatValues(ob = Object(), [value: '#(3147,3149)', option: "ports"])
		Assert(ob.ports is: #(3147, 3149))

		mock.ChooseManyOptionsAsObjectControl_stringMembers = #(testname)
		mock.origValues = #(companynum: 1)
		ob = Object()
		mock.formatValues(ob, [value: '2', option: "companynum"])
		mock.formatValues(ob, [value: '234', option: "testname"])
		Assert(ob.companynum is: 2)
		Assert(ob.testname is: '234')
		}
	}