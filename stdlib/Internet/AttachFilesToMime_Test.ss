// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_sendAsLink?()
		{
		fn = AttachFilesToMime.AttachFilesToMime_sendAsLink?

		Assert(fn(true, false, #()) is: false)
		Assert(fn(false, false, #()) is: false)
		Assert(fn(false, false, #('contribution')) is: false)
		Assert(fn(false, true, #('contribution')))
		Assert(fn(false, true, #()) is: false)
		}

	Test_allowSend?()
		{
		fn = AttachFilesToMime.AttachFilesToMime_allowSend?

		Assert(fn(true, false))
		Assert(fn('Attachments with extension', false) is: false)
		Assert(fn('Attachments with extension', true) is: false)
		Assert(fn('File size ', false) is: false)
		Assert(fn('File size ', true))
		}
	}