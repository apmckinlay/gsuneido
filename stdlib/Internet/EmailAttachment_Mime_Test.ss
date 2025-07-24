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
	}