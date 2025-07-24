// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		key = 'decrypttest'
		Assert(Decrypt('', key) is: "")
		Assert(Decrypt(123, key) is: 123)
		Assert(Decrypt(#20080101, key) is: #20080101)
		Assert(Decrypt(false, key) is: false)
		Assert(Decrypt('123-456-789', key) is: "UWP_MEBYRKM")
		}
	}