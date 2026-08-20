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

	Test_isThirdPartyError?()
		{
		fn = SuRender.SuRender_isThirdPartyError?

		// external stack (e.g. browser extension) - classified as third party
		Assert(fn(#(error:
			(stack: 'Error: oops\r\n    at chrome-extension://abc/foo.js:1:1'))))

		// translated JS uses $f frames; Blink/Chrome
		Assert(not fn(#(error: (stack: 'Error: bad\r\n    at $f (script:10:5)'))))
		// ...and WebKit/Safari
		Assert(not fn(#(error: (stack: 'Error: bad\r\n$f@script:10:5'))))

		// bundled/minified JS filenames - not third party
		Assert(not fn(#(error:
			(stack: 'Error: bad\r\n    at stuff (https://x/su_bundle.min.js:1:42)'))))
		Assert(not fn(#(error:
			(stack: 'Error: bad\r\n    at stuff (https://x/su_code_bundle.js:1:42)'))))

		// missing/short stacks - treated as ours (do not suppress)
		Assert(not fn(#(error: (stack: ''))))
		Assert(not fn(#(error: (stack: 'Error: just a message'))))

		// missing event.error - protected by outer try
		Assert(not fn(Object()))
		}
	}
