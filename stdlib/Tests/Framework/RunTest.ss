// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// simplest possible test runner
// primiarily for testing new implementations of Suneido
// usage e.g. RunTest(ObjectTest)
// Requirements:
// - classes with methods (for test class)
// - class.Members() to get test methods
// - object.Sort! (could be stub)
// - for in
// - string.Prefix?
// - try catch
// - string concatenation
// Plus whatever the tests require e.g. Assert and Catch
function (test)
	{
	s = ""
	mems = test.Members().Sort!()
	for (i = 0; i < mems.Size(); ++i) // don't require for-in loop
		{
		m = mems[i]
		if m.Prefix?('Test_')
			{
			try
				{
				test[m]()
				s $= 'PASS - ' $ m $ '\n'
				}
			catch (e)
				s $= 'FAIL ' $ m $ '\n\t' $ e $ '\n'
			}
		}
	return s
	}