// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		rec = []
		ob = AttachmentList(rec, 'attachments')
		Assert(ob isSize: 0)
		ob = AttachmentList.ListWithLabels(rec, 'attachments')
		Assert(ob isSize: 0)

		rec = [attachments: #()]
		ob = AttachmentList(rec, 'attachments')
		Assert(ob isSize: 0)
		ob = AttachmentList.ListWithLabels(rec, 'attachments')
		Assert(ob isSize: 0)

		rec = [
			attachments:
				#(
				[attachment3: `C:\attachments\test.pdf`, // full path
					attachment4: `C:\attachments\201911\cc.pdf`,
					attachment0: "C:\attachments\app1.jpg",
					attachment1: "C:\attachments\app2.jpg",
					attachment2: "C:\attachments\app3.jpg"],
				[attachment3: `C:\attachments\201911\a.pdf`],
				[],
				[attachment3: `C:\attachments\201911\spreadsheet.xlsm`])
			]
		ob = AttachmentList(rec, 'attachments')
		Assert(ob[0] is: "C:\attachments\app1.jpg")
		Assert(ob[1] is: "C:\attachments\app2.jpg")
		Assert(ob[2] is: "C:\attachments\app3.jpg")
		Assert(ob[3] is: "C:\attachments\\test.pdf")
		Assert(ob[4] is: "C:\attachments\201911\cc.pdf")
		Assert(ob[5] is: "C:\attachments\201911\a.pdf")
		Assert(ob[6] is: "C:\attachments\201911\spreadsheet.xlsm")
		Assert(ob isSize: 7)

		ob = AttachmentList.ListWithLabels(rec, 'attachments')
		Assert(ob[0] is: [path: "C:\attachments\app1.jpg", labels: ''])
		Assert(ob[1] is: [path: "C:\attachments\app2.jpg", labels: ''])
		Assert(ob[2] is: [path: "C:\attachments\app3.jpg", labels: ''])
		Assert(ob[3] is: [path: "C:\attachments\\test.pdf", labels: ''])
		Assert(ob[4] is: [path: "C:\attachments\201911\cc.pdf", labels: ''])
		Assert(ob[5] is: [path: "C:\attachments\201911\a.pdf", labels: ''])
		Assert(ob[6] is: [path: "C:\attachments\201911\spreadsheet.xlsm", labels: ''])
		Assert(ob isSize: 7)
		}
	}