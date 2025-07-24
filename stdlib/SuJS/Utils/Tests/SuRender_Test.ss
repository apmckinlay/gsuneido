// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_SetMouseMoveCB()
		{
		mock = Mock(SuRender)
		mock.When.SetMouseMoveCB([anyArgs:]).CallThrough()
		mock.When.ClearMouseMoveCB([anyArgs:]).CallThrough()

		// .mousemoveCB is false
		mock.SetMouseMoveCB(.call)
		mock.Verify.Never().ClearMouseMoveCB()

		// .mousemoveCB is still true from the previous test
		mock.SetMouseMoveCB(.call)
		mock.Verify.ClearMouseMoveCB()

		// .mousemoveCB is now false from being cleared via the previous test
		// Cleared properly via: ClearMouseMoveCB
		mock.ClearMouseMoveCB()
		mock.Verify.Times(2).ClearMouseMoveCB()
		mock.SetMouseMoveCB(.call)
		mock.Verify.Times(2).ClearMouseMoveCB()
		}

	call(unused)
		{ }

	Test_SetMouseUpCB()
		{
		mock = Mock(SuRender)
		mock.When.SetMouseUpCB([anyArgs:]).CallThrough()
		mock.When.ClearMouseUpCB([anyArgs:]).CallThrough()
		mock.When.restoreIframes([anyArgs:]).Do({ })
		mock.When.freezeIframes([anyArgs:]).Do({ })

		// .mouseupCB is false
		mock.SetMouseUpCB(.call)
		mock.Verify.Never().restoreIframes()
		mock.Verify.freezeIframes()

		// .mouseupCB is still true from the previous test
		mock.SetMouseUpCB(.call)
		mock.Verify.restoreIframes()
		mock.Verify.Times(2).freezeIframes()

		// .mouseupCB is still true from the previous tests
		// Cleared properly via: ClearMouseUpCB
		mock.ClearMouseUpCB()
		mock.Verify.Times(2).restoreIframes()
		mock.SetMouseUpCB(.call)
		mock.Verify.Times(2).restoreIframes()
		mock.Verify.Times(3).freezeIframes()
		}
	}