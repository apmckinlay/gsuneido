// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		w = Object()
		name = 'TestClass'
		newCode = RemoveUnderscoreRecordName(name, '', w)
		Assert(newCode is: '')
		Assert(w is: [])

		code = "//_TestClass
_TestClass
	{
	New()
		{
		Print('_TestClass' $ '_TestClasses' $ 'test_TestClass')
		}
	}"
		expectedCode = "// TestClass
 TestClass
	{
	New()
		{
		Print(' TestClass' $ '_TestClasses' $ 'test_TestClass')
		}
	}"
		newCode = RemoveUnderscoreRecordName(name, code, w)
		Assert(newCode is: expectedCode)
		Assert(w is: [2])

		newCode = RemoveUnderscoreRecordName('TestClass__alpha', code, w = Object())
		Assert(newCode is: expectedCode)
		Assert(w is: [2])
		}
	}