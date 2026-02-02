// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(ValidEmailAddresses?('') is: false)

		Assert(ValidEmailAddresses?('test@abc.com'))

		Assert(ValidEmailAddresses?('test@abc.com,test2@abc.com'))
		Assert(ValidEmailAddresses?('test@abc.com;test2@abc.com'))
		Assert(ValidEmailAddresses?('test@abc.com;test2@abc.com,abc@abc.com'))

		Assert(ValidEmailAddresses?('test@abc.com, test2@abc.com'))
		Assert(ValidEmailAddresses?('test@abc.com; test2@abc.com'))
		Assert(ValidEmailAddresses?('test@abc.com; test2@abc.com, abc@abc.com'))
		Assert(ValidEmailAddresses?('test@abc.com ; test2@abc.com;'))

		Assert(ValidEmailAddresses?('test@abc.com; test2@abc. com') is: false)
		Assert(ValidEmailAddresses?('test.abc.com') is: false)
		Assert(ValidEmailAddresses?('test@abc.com,test2.abc.com') is: false)
		Assert(ValidEmailAddresses?('test@abc.com;test2.abc.com') is: false)

		Assert(ValidEmailAddresses?('test@abc.com;;; test2@abc.com') is: false)
		Assert(ValidEmailAddresses?('test@abc.com,,, test2@abc.com') is: false)
		Assert(ValidEmailAddresses?('test@abc.com;,test2@abc.com') is: false)
		}
	}