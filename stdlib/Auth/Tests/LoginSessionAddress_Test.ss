// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		address = LoginSessionAddress
			{
			LoginSessionAddress_databaseSessionId() { return '127.0.0.15' }
			LoginSessionAddress_getOsName() { return 'Windows 10 Pro' }
			}
		Assert(address() is: '127.0.0.15')

		address = LoginSessionAddress
			{
			LoginSessionAddress_databaseSessionId() { return '127.0.0.15' }
			LoginSessionAddress_getOsName() { return 'Windows 11 Pro' }
			}
		Assert(address() is: '127.0.0.15')

		address = LoginSessionAddress
			{
			LoginSessionAddress_wts_GetSessionId() { return 10 }
			LoginSessionAddress_getOsName() { return 'Windows Server 2012 R2' }
			}
		Assert(address() is: 'wts10')

		address = LoginSessionAddress
			{
			LoginSessionAddress_wts_GetSessionId() { return 12 }
			LoginSessionAddress_getOsName() { return 'Windows Server 2016' }
			}
		Assert(address() is: 'wts12')
		}
	}