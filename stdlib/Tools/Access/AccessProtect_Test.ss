// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_buildProtectedReason()
		{
		fn = AccessProtect.AccessProtect_buildProtectedReason
		Assert(fn(AccessControl, false)
			is: 'This record can not be deleted.')

		acc = AccessProtect
			{
			AccessProtect_protectResult(@unused) { return true }
			}
		Assert(acc.AccessProtect_buildProtectedReason(AccessControl, true)
			is: 'This record can not be deleted.')

		acc = AccessProtect
			{
			AccessProtect_protectResult(@unused)
				{ return 'Transaction dated prior to protection date' }
			}
		Assert(acc.AccessProtect_buildProtectedReason(AccessControl, true)
			is: 'This record can not be deleted.\n\n' $
			'Transaction dated prior to protection date')

		acc = AccessProtect
			{
			AccessProtect_protectResult(@unused)
				{ return Object('allbut', testfield:) }
			}
		Assert(acc.AccessProtect_buildProtectedReason(AccessControl, true)
			is: 'This record can not be deleted.')

		acc = AccessProtect
			{
			AccessProtect_protectResult(@unused)
				{ return Object('allbut', testfield:, reason: 'Pay finalized') }
			}
		Assert(acc.AccessProtect_buildProtectedReason(AccessControl, true)
			is: 'This record can not be deleted.\n\nPay finalized')
		}
	}