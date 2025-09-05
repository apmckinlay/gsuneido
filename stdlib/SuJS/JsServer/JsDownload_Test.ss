// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_CallClass()
		{
		Assert(JsDownload([queryvalues: #()])
			is: Object('Unauthorized', [], 'Your session is invalid or expired'))

		cl = JsDownload
			{
			JsDownload_invalidToken?(env /*unused*/) { return false }
			}
		Assert(cl([queryvalues: #()])
			is: ['BadRequest', [], 'Invalid request, missing file name'])

		Assert(cl([queryvalues: Object(Base64.Encode(
			.TempTableName().Xor(EncryptControlKey())))])
			is: #("404 Not Found", #(), "not found"))
		}
	}