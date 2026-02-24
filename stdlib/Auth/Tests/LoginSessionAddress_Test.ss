// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		address = LoginSessionAddress
			{
			LoginSessionAddress_databaseSessionId() { return '127.0.0.15' }
			LoginSessionAddress_getOsName() { return 'Windows 10 Pro' }
			LoginSessionAddress_js?() { return false }
			}
		Assert(address() is: '127.0.0.15')

		address = LoginSessionAddress
			{
			LoginSessionAddress_databaseSessionId() { return '127.0.0.15' }
			LoginSessionAddress_getOsName() { return 'Windows 11 Pro' }
			LoginSessionAddress_js?() { return false }
			}
		Assert(address() is: '127.0.0.15')

		address = LoginSessionAddress
			{
			LoginSessionAddress_wts_GetSessionId() { return 10 }
			LoginSessionAddress_getOsName() { return 'Windows Server 2012 R2' }
			LoginSessionAddress_js?() { return false }
			}
		Assert(address() is: 'wts10')

		address = LoginSessionAddress
			{
			LoginSessionAddress_wts_GetSessionId() { return 12 }
			LoginSessionAddress_getOsName() { return 'Windows Server 2016' }
			LoginSessionAddress_js?() { return false }
			}
		Assert(address() is: 'wts12')

		address = LoginSessionAddress
			{
			LoginSessionAddress_databaseSessionId() { return 'user@ip<token>(jsS)' }
			LoginSessionAddress_getOsName() { return 'Windows Server 2016' }
			LoginSessionAddress_js?() { return true }
			}
		Assert(address() is: 'user@ip<token>(jsS)')
		}
	}