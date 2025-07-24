// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_updateAttachmentsList()
		{
		method = EmailAttachment.EmailAttachment_updateAttachmentsList
		data = []
		//Testing rearrangement
		data.attachments = Object("(APPEND) file1.pdf", "(APPEND) file2.pdf",
			"(APPEND) file4.pdf", "file3.txt")
		data.attachmentsList = Object(Object(File: "(APPEND) file1.pdf"),
			Object(File: "(APPEND) file2.pdf"), Object(File: "(APPEND) file4.pdf"),
			Object(File: "file3.txt"))
		mergeableFiles = Object("file1.pdf", "file2.pdf", "file4.pdf")
		newAppend = Object("file2.pdf", "file1.pdf", "file4.pdf")
		method(mergeableFiles, newAppend, data)
		expectedAttachments = Object("(APPEND) file2.pdf", "(APPEND) file1.pdf",
			"(APPEND) file4.pdf", "file3.txt")
		Assert(data.attachments is: expectedAttachments)
		//Testing removal of attachments
		data.attachments = Object("(APPEND) file1.pdf", "(APPEND) file2.pdf",
			"(APPEND) file4.pdf", "file3.txt")
		data.attachmentsList = Object(Object(File: "(APPEND) file1.pdf"),
			Object(File: "(APPEND) file2.pdf"), Object(File: "(APPEND) file4.pdf"),
			Object(File: "file3.txt"))
		mergeableFiles = Object("file1.pdf", "file2.pdf", "file4.pdf")
		newAppend = Object("file2.pdf", "file4.pdf")
		method(mergeableFiles, newAppend, data)
		expectedAttachments = Object("(APPEND) file2.pdf", "(APPEND) file4.pdf",
			"file1.pdf", "file3.txt")
		Assert(data.attachments is: expectedAttachments)
		}

	Test_updateAttachmentsList2()
		{
		method = EmailAttachment.EmailAttachment_updateAttachmentsList
		data = []
		//Testing rearrangment when some arent in attachmentsList
		data.attachments = Object("(APPEND) file1.pdf", "(APPEND) file2.pdf",
			"(APPEND) file4.pdf", "file3.txt")
		data.attachmentsList = Object(Object(File: "(APPEND) file1.pdf"),
			Object(File: "(APPEND) file2.pdf"), Object(File: "(APPEND) file5.jpg"),
			Object(File: "(APPEND) file4.pdf"), Object(File: "file3.txt"))
		mergeableFiles = Object("file1.pdf", "file2.pdf", "file4.pdf")
		newAppend = Object("file2.pdf", "file1.pdf", "file4.pdf")
		method(mergeableFiles, newAppend, data)
		expectedAttachments = Object("(APPEND) file2.pdf", "(APPEND) file1.pdf",
			"(APPEND) file4.pdf", "file3.txt")
		Assert(data.attachments is: expectedAttachments)
		//Testing removal when some arent in attachmentsList and comma in filename
		data.attachments = Object("(APPEND) file,1.pdf", "(APPEND) file2.pdf",
			"(APPEND) file4.pdf", "file3.txt")
		data.attachmentsList = Object(Object(File: "(APPEND) file,1.pdf"),
			Object(File: "(APPEND) file2.pdf"), Object(File: "(APPEND) file5.jpg"),
			Object(File: "(APPEND) file4.pdf"), Object(File: "file3.txt"))
		mergeableFiles = Object("file,1.pdf", "file2.pdf", "file4.pdf")
		newAppend = Object("file2.pdf", "file4.pdf")
		method(mergeableFiles, newAppend, data)
		expectedAttachments = Object("(APPEND) file2.pdf", "(APPEND) file4.pdf",
			"file,1.pdf", "file3.txt")
		Assert(data.attachments is: expectedAttachments)
		}

	Test_attachmentFilesExist()
		{
		Assert(EmailAttachment.AttachmentFilesExist(#()) is: '')
		Assert(EmailAttachment.AttachmentFilesExist('') is: '')

		.name = .TempTableName()
		name2 = .TempTableName()
		Assert(EmailAttachment.AttachmentFilesExist(Object(.name, .name, name2))
			is: "Could not find the following file(s):\n\n" $ .name $ '\n' $ name2 $
			'\n\nPlease ensure that they exist and ' $
			'that you have permission to access them.')

		// Teardown is deleting .name for the PutFile
		.PutFile(.name, 'test')
		Assert(EmailAttachment.AttachmentFilesExist(Object(.name, .name, name2))
			is: "Could not find the following file(s):\n\n" $ name2 $
			'\n\nPlease ensure that they exist and ' $
			'that you have permission to access them.')

		.PutFile(.name, 'test')
		Assert(EmailAttachment.AttachmentFilesExist(Object(.name, .name)) is: "")

		Assert(EmailAttachment.AttachmentFilesExist(#()) is: '')
		}

	Test_attachmentList()
		{
		f = EmailAttachment.EmailAttachment_attachmentList
		Assert(f(#()) is: #(attachments: #(), filename: ''))
		Assert(f(#(attachments: '')) is: #(attachments: #(), filename: ''))
		Assert(f(#(attachments: #())) is: #(attachments: #(), filename: ''))

		Assert(f(#(attachments: #(#(attachFileName: 'f.txt', fileName: 'g.txt'))))
			is: #(attachments: #('g.txt'), filename: ''))

		Assert(f(#(attachments: '', filename: "C:/arinvoice_115850283.pdf")) is:
				#(filename: "C:/arinvoice_115850283.pdf", attachments: #()))

		// mimicking std email attachments. attachments is object of strings
		stdAttachOb = #(attachments:
			#(`C:\work\eta\rpt2\Order Confirmation.pdf`,
			`C:\work\eta\rpt2\Print Financial Statement Specs.pdf`,
			`C:\work\eta\output.txt`),
			filename: "C:/Users/Gerry/AppData/Local/Temp/su944107614.pdf")

		Assert(f(stdAttachOb) is:
			#(attachments:
				#(`C:\work\eta\rpt2\Order Confirmation.pdf`,
				`C:\work\eta\rpt2\Print Financial Statement Specs.pdf`,
				`C:\work\eta\output.txt`),
			filename: "C:/Users/Gerry/AppData/Local/Temp/su944107614.pdf"))

		// mimicking batch email attachments; attachments is object of ojects
		batchEmailAttachOb = #(attachments: #(
			#(attachFileName: `C:\Order Confirmation.pdf`,
				fileName: `C:\HELLO WORLD.pdf`),
			#(attachFileName: `Z:\BLAH.pdf`,
				fileName: `C:\FRED.pdf`),
			#(attachFileName: `C:\\output.txt`,
				fileName: `C:\output.txt`)),
			filename: "C:/AXONETA/workqueue/arinvoice/arinvoice_20210817_115850283.pdf")

		Assert(f(batchEmailAttachOb) is:
			#(attachments: #(`C:\HELLO WORLD.pdf`, `C:\FRED.pdf`, `C:\output.txt`),
			filename: `C:/AXONETA/workqueue/arinvoice/arinvoice_20210817_115850283.pdf`))
		}

	Teardown()
		{
		DeleteFile(.name)
		super.Teardown()
		}
	}
