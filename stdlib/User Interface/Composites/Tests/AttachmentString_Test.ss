// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		attach = ''
		Assert(AttachmentString(attach) is: '')
		attach = Object()
		Assert(AttachmentString(attach) is: '')
		attach = Object([attachment0: "This is just some file"])
		Assert(AttachmentString(attach) is: 'This is just some file')
		attach = Object([attachment1: "This is a second File",
			attachment0: "This is just some file"])
		Assert(AttachmentString(attach)
			is: "This is just some file, This is a second File")
		attach = Object([attachment1: "This is on the first row"],
			[attachment0: 'This is on the second row',
				attachment1: 'This is on the first row']) // test duplicate attachment
		Assert(AttachmentString(attach)
			is: "This is on the first row, This is on the second row")
		}
	}