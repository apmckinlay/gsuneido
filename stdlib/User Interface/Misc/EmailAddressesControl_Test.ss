// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_KillFocus()
		{
		.SpyOn(MultiAutoChooseControl.KillFocus).Return(false)
		.SpyOn(EmailAddressesControl.Get).Return('', ' , ', 'a@b.cc', 'ab.c , ,d@e.ff, ')
		spy = .SpyOn(EmailAddresses.OutputAddr).Return(false)
		callLog = spy.CallLogs()

		EmailAddressesControl.KillFocus()
		Assert(callLog isSize: 0) // testing ''

		EmailAddressesControl.KillFocus()
		Assert(callLog isSize: 0) // testing ' , '

		EmailAddressesControl.KillFocus()
		Assert(callLog isSize: 1)
		Assert(callLog[0] is: #(addr: 'a@b.cc', t: false)) // testing 'a@b.cc'

		EmailAddressesControl.KillFocus()
		Assert(callLog isSize: 2)
		Assert(callLog[1] is: #(addr: 'd@e.ff', t: false)) // 'ab.c , ,d@e.ff, '
		}
	}