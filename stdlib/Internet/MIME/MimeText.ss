// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
MimeBase
	{
	New(text = "", subtype = "plain", charset = "us-ascii")
		{
		super('text', subtype)
		.AddExtra('Content-Type', :charset)
		.Content_Transfer_Encoding('7bit')
		.SetPayload(text.ChangeEol('\r\n'))
		}

	ToString()
		{
		return .MimeVersion $
			.ExtraHeader(.Fields.Copy().Remove('Content-Transfer-Encoding')) $
			.ContentType() $ .ExtraHeader(#('Content-Transfer-Encoding')) $
			.MessageContent()
		}
	}
