// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		cl = ConfiglibHttp
			{
			ActionsHttp_runOnHttp?() { return false }

			ActionsHttp_getActions()
				{
				return #(
					test_add: function(a, b) { return a $ " + " $ b $ " = " $ (a + b) },
					test_args: function(@args) { return args },
					test_false: function() { return false },
					test_max: Max,
					test_fail: function(@unused) { throw "test failure" }
					)
				}
			}
		Assert(cl.Run(#test_add, 1, 2) is: "1 + 2 = 3")
		Assert(cl.Run(#test_add, b: 4, a: 3) is: "3 + 4 = 7")
		Assert(cl.Run(#test_args, 1, 2, name: #value)
			is: #(1, 2, name: value))
		Assert(cl.Run(#test_false) is: false)
		Assert(cl.Run(#test_max, 2, 5) is: 5)

		ServerSuneido.Set(#TestRunningExpectedErrors,
			[
				"ERROR: ConfiglibHttp: action not allowed: unknown",
				"ERROR: ConfiglibHttp: action not allowed: ",
				"ERROR: ConfiglibHttp: test_fail: test failure"])

		Assert({ cl.Run() } throws: "missing action")
		Assert({ cl.Run(#unknown) } throws: "ConfiglibHttp: action not allowed: unknown")
		result = cl.Run("") // return throw
		Assert(result is: "ConfiglibHttp: action not allowed: ")
		Assert({ cl.Run(#test_fail, 1, name: #value) }
			throws: "ConfiglibHttp: test_fail: test failure")

		Assert(ServerSuneido.Get(#TestRunningExpectedErrors) is: #())
		}
	}
