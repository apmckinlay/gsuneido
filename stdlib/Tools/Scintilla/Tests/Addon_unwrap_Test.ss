// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_unwrap()
		{
		unwrap = Addon_unwrap.Addon_unwrap_getUnwrapped
		Assert(unwrap('') is: '')
		Assert(unwrap('"Some text"') is: 'Some text')
		Assert(unwrap('[rec: item]') is: 'rec: item')
		Assert(unwrap('(layoutOb = MsgLayouts().GetDefault(a, false))')
			is: 'layoutOb = MsgLayouts().GetDefault(a, false)')
		Assert(unwrap('{ return inlineReturn }') is: 'return inlineReturn')
		Assert(unwrap('{
			return 0
			}') is: 'return 0')
		Assert(unwrap('{
			var1 = true
			var2 = 2
			}') is: 'var1 = true
			var2 = 2')
		Assert(unwrap('{
			if condition
				{
				var1 = true
				var2 = 2
				}
			return 0
			}') is: 'if condition
				{
				var1 = true
				var2 = 2
				}
			return 0')
		Assert(unwrap("{
			if customerFileValues().ports.Size() > 1
				return true

			try
				return LastContribution('TestContrib')
			return false
			}") is: "if customerFileValues().ports.Size() > 1
				return true

			try
				return LastContribution('TestContrib')
			return false")
		}
	}