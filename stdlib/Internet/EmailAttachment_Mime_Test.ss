// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_buildCompressedAttachments()
		{
		buildCompressedAttachments = EmailAttachment_Mime.BuildCompressedAttachments
		attachOb = Object()
		attachOb.attachments = #(
			"C:/Local/Temp/1.pdf_temp_compressed20220310152543697",
			`C:\Desktop\test\2.pdf`,
			`C:\Users\Desktop\test\3.pdf`,
			"C:/Local/Temp/4.pdf_temp_compressed20220310152543698",
			`C:\Users\Desktop\test\5.pdf`)
		buildCompressedAttachments(x = Object(), attachOb)
		Assert(x is: Object(attachments: Object(
			Object(fileName: "C:/Local/Temp/1.pdf_temp_compressed20220310152543697",
				attachFileName: 'C:/Local/Temp/1.pdf'),
			Object(fileName: `C:\Desktop\test\2.pdf`,
				attachFileName: `C:\Desktop\test\2.pdf`),
			Object(fileName: `C:\Users\Desktop\test\3.pdf`,
				attachFileName: `C:\Users\Desktop\test\3.pdf`),
			Object(fileName: "C:/Local/Temp/4.pdf_temp_compressed20220310152543698",
				attachFileName: 'C:/Local/Temp/4.pdf'),
			Object(fileName: `C:\Users\Desktop\test\5.pdf`,
				attachFileName: `C:\Users\Desktop\test\5.pdf`)
			)))

		attachOb = Object()
		attachOb.attachments = #()
		buildCompressedAttachments(x = Object(), attachOb)
		Assert(x is: Object(attachments: #()))

		attachOb = Object()
		attachOb.attachments = #("C:/Local/Temp/1.pdf_temp_compressed20220310152543697")
		buildCompressedAttachments(x = Object(), attachOb)
		Assert(x is: Object(attachments:
			Object(#(fileName: "C:/Local/Temp/1.pdf_temp_compressed20220310152543697",
				attachFileName: 'C:/Local/Temp/1.pdf'))))

		attachOb = Object()
		attachOb.attachments = #("C:/Local/Temp/1.pdf")
		buildCompressedAttachments(x = Object(), attachOb)
		Assert(x is: Object(attachments:
			Object(#(fileName: "C:/Local/Temp/1.pdf",
				attachFileName: 'C:/Local/Temp/1.pdf'))))
		}

	Test_collectPresignedUrls()
		{
		if AttachmentS3Bucket() is ''
			return
		cl = EmailAttachment_Mime
			{
			EmailAttachment_Mime_preSignedUrl(@args)
				{
				return 's3://' $ args[1]
				}
			EmailAttachment_Mime_token()
				{
				return 'token'
				}
			}
		fn = cl.EmailAttachment_Mime_collectPresignedUrls
		fn([], history = Object(), 'bucket', preUrls = Object())
		Assert(history is: #())
		Assert(preUrls is: #())

		data = [attachments: Object('linkedattachments\202104/test.pdf',
			'202104/test2.pdf', `(APPEND) 202104\test3.pdf`, `(APPEND) temp/test4.pdf`)]
		fn(data, history = Object(), 'bucket', preUrls = Object())
		Assert(data.attachments is: Object('linkedattachments/202104/test.pdf',
			'202104/test2.pdf', '(APPEND) 202104/test3.pdf', `(APPEND) temp/test4.pdf`))

		curYearMonth = Date().Format('yyyyMM')
		Assert(history["test.pdf"].tmpName startsWith: curYearMonth $ "/test")
		Assert(history["test.pdf"].url startsWith: "s3://" $ curYearMonth $ "/test")
		Assert(history["test2.pdf"].tmpName startsWith: curYearMonth $ "/test2")
		Assert(history["test2.pdf"].url startsWith: "s3://" $ curYearMonth $ "/test2")

		Assert(preUrls[0] is: "s3://linkedattachments/202104/test.pdf")
		Assert(preUrls[1] is: "s3://202104/test2.pdf")
		Assert(preUrls[2] is: "s3://202104/test3.pdf")
		Assert(preUrls[3] matches: "download?.*&token=token&saveName=test4.pdf")
		}
	}