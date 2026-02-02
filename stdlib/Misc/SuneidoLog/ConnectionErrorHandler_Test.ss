// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	errorOb: #("couldn't connect to host",
		"Could not resolve host",
		"We encountered an internal error. Please try again.")
	Test_TreatConnectionErrorAsInfo()
		{
		cl = ConnectionErrorHandler
			{
			MessageManager_ClassId: 'MessageManager_ClassId_Test'
			New() { .LogOb = Object() }
			AddToLog(msg, caughtMsg) { .LogOb.Add(Object(:msg, :caughtMsg)) }
			CalcNumberOfConnectionErrs(classId /*unused*/, errorOb /*unused*/,
				errMsg /*unused*/)
				{ return 3 }
			}

		.connectionErrorAsInfo(cl, .errorOb, "curl: (7) couldn't connect to host")
		.connectionErrorAsInfo(cl, .errorOb, "Could not resolve host")
		.connectionErrorAsInfo(cl, .errorOb,
			'"1.0"?>' $
			'<ErrorResponse xmlns="http://queue.amazonaws.com/doc/2008-01-01/">' $
			'<Error><Type>Receiver</Type><Code>InternalError</Code>' $
			'<Message>We encountered an internal error. Please try again.' $
			'</Message><Detail/></Error>'$
			'<RequestId>27453add-4e5a-42ed-bfb7-ee58c5a23bfa</RequestId>' $
			'</ErrorResponse>')
		}

	connectionErrorAsInfo(cl, errorOb, errMsg)
		{
		msg = " (send)"
		c = new cl
		c.Process(errMsg, errorOb, msg, 'MessageManager_ClassId_Test')
		err = errMsg.Has?('<Message>')
			? errMsg.AfterFirst('<Message>').BeforeFirst('</Message>')
			: errMsg
		Assert(c.LogOb[0].msg
			has: "INFO: MessageManager_ClassId_Test: " $ msg $ ' - ' $ err)
		Assert(c.LogOb[0].caughtMsg is: '')
		}

	Test_TreatConnectionErrorAsError()
		{
		cl = ConnectionErrorHandler
			{
			MessageManager_ClassId: 'MessageManager_ClassId_Test'
			New() { .LogOb = Object() }
			AddToLog(msg, caughtMsg) { .LogOb.Add(Object(:msg, :caughtMsg)) }
			CalcNumberOfConnectionErrs(classId /*unused*/, errorOb /*unused*/,
				errMsg /*unused*/)
				{ return 11 }
			}
		msg = " (send)"
		c = new cl
		c.Process('Could not resolve host', .errorOb, msg, 'MessageManager_ClassId_Test')
		Assert(c.LogOb[0].msg
			has: "ERROR: (CAUGHT) MessageManager_ClassId_Test: " $ msg $
				' - Could not resolve host')
		Assert(c.LogOb[0].caughtMsg
			is: "generic connection error handling; may need attention")
		}

	Test_NonConnectionError()
		{
		cl = ConnectionErrorHandler
			{
			MessageManager_ClassId: 'MessageManager_ClassId_Test'
			New() { .LogOb = Object() }
			AddToLog(msg, caughtMsg) { .LogOb.Add(Object(:msg, :caughtMsg)) }
			}
		errMsg = '<?xml version="1.0"?>\r\n' $
			'<ErrorResponse xmlns="http://queue.amazonaws.com/doc/2008-01-01/">' $
			'<Error><Type>Sender</Type><Code>SignatureDoesNotMatch</Code>' $
			'<Message>The request signature we calculated does not match the ' $
			'signature you provided. Check your AWS Secret Access Key and ' $
			'signing method. Consult the service documentation for details.' $
			'</Message><Detail/></Error>' $
			'<RequestID>131dd6e6-5ea4-453d-b3b6-eccf373253eb</RequestID>' $
			'</ErrorResponse>'
		msg = "Response"
		c = new cl
		c.Process(errMsg, .errorOb, msg, 'MessageManager_ClassId_Test')
		Assert(c.LogOb[0].msg
			has: "ERROR: (CAUGHT) MessageManager_ClassId_Test: " $ msg $
				' - The request signature we calculated does not match the ' $
				'signature you provided. Check your AWS Secret Access Key and ' $
				'signing method. Consult the service documentation for details.')
		Assert(c.LogOb[0].caughtMsg
			is: "generic connection error handling; may need attention")
		}

	Teardown()
		{
		id = ConnectionErrorHandler.ConnectionErrorHandler_idPrefix $
			'MessageManager_ClassId_Test'
		ServerSuneido.DeleteMember(id)
		super.Teardown()
		}
	}