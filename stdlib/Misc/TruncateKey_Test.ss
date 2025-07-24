// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		key = "Test123456789"
		Assert(TruncateKey(key) is: "Test")
		key = "This (set of #ch@racters) is a te$t!"
		Assert(TruncateKey(key, replace: "[!|@|#|$|%|^|&|*|(|)]")
			is: "This set of chracters is a tet")
		key = "This should stay the same"
		Assert(TruncateKey(key, replace: "") is: "This should stay the same")
		key = "this is more than 50".RightFill(150, "B")
		Assert(TruncateKey(key) is: "this is more than 50".RightFill(50, "B"))
		key = "abcdefghijklmnopqrstuvwxyz"
		Assert(TruncateKey(key, length: 10) is: "abcdefghij")
		}
	}
