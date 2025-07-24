// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.SpyOn(OpenImageWithLabelsControl.OpenImageWithLabelsControl_getCopyTo).
			Return('', '', `Z:\copyTo\`)
		params = [EmailAttachments: Object()]
		data = []
		EmailLabelledAttachments(params, data, '', '')
		Assert(params is: [EmailAttachments: Object()])

		data = [attachment: [['bobsFile']]]
		EmailLabelledAttachments(params, data, '', 'attachment')
		Assert(params is: [EmailAttachments: Object()])

		params = [ReportDestination: 'pdf', EmailAttachments: Object()]
		EmailLabelledAttachments(params, data, '', 'attachment')
		Assert(params.EmailAttachments is: Object())

		data = [attachment: [['bobsFile  Axon Label: email with invoice']]]
		EmailLabelledAttachments(params, data, '', 'attachment')
		Assert(params.EmailAttachments is: Object())

		EmailLabelledAttachments(params, data, 'email with invoice', 'attachment')
		Assert(params.EmailAttachments is: Object('bobsFile'))

		// bobsFile is already in params, so it is not added a second time
		data = [attachment: [['bobsFile  Axon Label: email with invoice',
			'larrysFile', 'jakesFile Axon Label: email with customer report']]]
		EmailLabelledAttachments(params, data, 'email with invoice', 'attachment')
		Assert(params.EmailAttachments is: Object(`bobsFile`))

		params = [ReportDestination: 'pdf', EmailAttachments: Object()]
		data = [attachment: [['bobsFile  Axon Label: email with invoice',
			'larrysFile', 'jakesFile Axon Label: email with customer report']]]
		// Not an actual label, but EmailLabelledAttachments only matches what it is given
		EmailLabelledAttachments(params, data, 'email with customer report', 'attachment')
		Assert(params.EmailAttachments is: Object(`Z:\copyTo/jakesFile`))
		}
	}
