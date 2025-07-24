// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_valid?()
		{
		m = SvcServerManager.SvcServerManager_valid?
		Assert(m(data = []) is: 'Password is required')

		data.svcmng_password = 'password'
		Assert(m(data) is: 'Password must match Verify')

		data.svcmng_password_verify = 'password2'
		Assert(m(data) is: 'Password must match Verify')

		data.svcmng_password_verify = data.svcmng_password
		Assert(m(data) is: '')

		data.svcmng_create_custom_library = true
		Assert(m(data) is: 'Custom SVC library is required')

		data.svcmng_custom_library = 'custlib'
		Assert(m(data) is: '')
		}
	}
