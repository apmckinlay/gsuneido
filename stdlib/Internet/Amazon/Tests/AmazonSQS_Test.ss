// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_SendMessage()
		{
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return('postResult')
		mock.When.handleInvalidChars([anyArgs:]).CallThrough()
		mock.When.makeRequest([anyArgs:]).CallThrough()

		// send
		Assert(mock.Eval(AmazonSQS.SendMessage, 'send', '/testQueue', '')
			is: 'postResult')
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/testQueue',
			'Action=SendMessage', '/testQueue')
		Assert(mock.Eval(AmazonSQS.SendMessage, 'send', '/testQueue', 'test message')
			is: 'postResult')
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/testQueue',
			'Action=SendMessage&MessageBody=test%20message', '/testQueue')
		Assert(mock.Eval(AmazonSQS.SendMessage, 'send', '/testQueue', `!@#$^&*()+=<>/`)
			is: 'postResult')
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/testQueue',
			'Action=SendMessage&MessageBody=%21%40%23%24%5E%26%2A%28%29%2B%3D%3C%3E%2F',
			'/testQueue')
		Assert(mock.Eval(AmazonSQS.SendMessage, 'send', '/testQueue', 'Hello
			\t\r\nWorld')
			is: 'postResult')
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/testQueue',
			'Action=SendMessage&MessageBody=Hello%0D%0A%09%09%09%09%0D%0AWorld',
			'/testQueue')

		// receive
		Assert(mock.Eval(AmazonSQS.SendMessage, 'receive', '/testQueue', '')
			is: 'postResult')
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/testQueue',
			'Action=ReceiveMessage&MaxNumberOfMessages=' $
			AmazonSQS.AmazonSQS_maxToReceive, '/testQueue')

		// delete
		Assert(mock.Eval(AmazonSQS.SendMessage, 'delete', '/testQueue',
			'52UWh5+c2xTG2YW95TdV3YpznHU1oEhNNtt3wNcJD0J/fF')
			is: 'postResult')
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/testQueue',
			'Action=DeleteMessage&ReceiptHandle=' $
			'52UWh5%2Bc2xTG2YW95TdV3YpznHU1oEhNNtt3wNcJD0J%2FfF', '/testQueue')

		// errors
		mock = Mock(AmazonSQS)
		mock.When.handleInvalidChars([anyArgs:]).CallThrough()
		mock.When.makeRequest([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.SendMessage, 'invalidAction', '', '')
			is: 'ERROR: Invalid Action for SendMessage: invalidAction')
		mock.When.post([anyArgs:]).Throw("ERROR: POST FAILED")
		Assert(mock.Eval(AmazonSQS.SendMessage, 'send', '/testQueue', 'test message')
			is: 'ERROR: POST FAILED')

		msg = "ERROR: Could not build signed string - unable to get temporary credentials"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg))
		mock = Mock(AmazonSQS)
		mock.When.handleInvalidChars([anyArgs:]).CallThrough()
		mock.When.makeRequest([anyArgs:]).CallThrough()
		mock.When.post([anyArgs:]).Return(AmazonAWS.CredentialErrMsg)
		Assert(mock.Eval(AmazonSQS.SendMessage, 'send', '/testQueue', 'test message')
			is: AmazonAWS.CredentialErrMsg)
		}

	Test_CheckQueue()
		{
		mock = Mock(AmazonSQS)
		validXml = .checkQueueXmlStr('11')
		mock.When.post([anyArgs:]).Return(validXml)
		mock.When.xmlQueueResult([anyArgs:]).CallThrough()
		mock.When.queueAttributesResponse([anyArgs:]).CallThrough()

		// valid return values
		Assert(mock.Eval(AmazonSQS.CheckQueue, 'testQueue')
			is: XmlParser(validXml))
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/testQueue',
			'Action=GetQueueAttributes&AttributeName.1=All', '/testQueue')

		Assert(mock.Eval(AmazonSQS.CheckQueue, 'testQueue', 'random')
			is: XmlParser(validXml))
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/testQueue',
			'Action=GetQueueAttributes&AttributeName.1=random', '/testQueue')

		// errors
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return('<invalid>stuff<>')
		mock.When.xmlQueueResult([anyArgs:]).CallThrough()
		mock.When.queueAttributesResponse([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.CheckQueue, 'testQueue') is: false)

		mock = Mock(AmazonSQS)
		mock.When.xmlQueueResult([anyArgs:]).CallThrough()
		mock.When.queueAttributesResponse([anyArgs:]).CallThrough()
		mock.When.post([anyArgs:]).Throw('CONNECTION ERROR')
		Assert(mock.Eval(AmazonSQS.CheckQueue, 'testQueue') is: false)

		msg = "ERROR: Could not build signed string - unable to get temporary credentials"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg))
		mock = Mock(AmazonSQS)
		mock.When.xmlQueueResult([anyArgs:]).CallThrough()
		mock.When.queueAttributesResponse([anyArgs:]).CallThrough()
		mock.When.post([anyArgs:]).Return(AmazonAWS.CredentialErrMsg)
		Assert(mock.Eval(AmazonSQS.CheckQueue, 'testQueue') is: false)
		}

	checkQueueXmlStr(value)
		{
		return '<getqueueattributesresponse>' $
			'<getqueueattributesresult>' $
				'<attribute>' $
					'<value>' $
						value $
					'</value>' $
				'</attribute>' $
			'</getqueueattributesresult>' $
		'</getqueueattributesresponse>'
		}

	Test_ApproximateNumberOfMessages()
		{
		mock = Mock(AmazonSQS)
		validXml = .checkQueueXmlStr('11')
		mock.When.post([anyArgs:]).Return(validXml)
		mock.When.CheckQueue([anyArgs:]).CallThrough()

		// valid return values
		Assert(mock.Eval(AmazonSQS.ApproximateNumberOfMessages, 'testQueue')
			is: 11)
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/testQueue',
			'Action=GetQueueAttributes&AttributeName.1=ApproximateNumberOfMessages',
				'/testQueue')

		// errors
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return('<invalid>stuff<>')
		mock.When.CheckQueue([anyArgs:]).CallThrough()
		Assert({mock.Eval(AmazonSQS.ApproximateNumberOfMessages, 'testQueue')}
			throws: 'ApproximateNumberOfMessages: Problem accessing Message Queue')

		mock = Mock(AmazonSQS)
		invalidXmlForMessages = '<invalidxml>hello world</invalidxml>'
		mock.When.post([anyArgs:]).Return(invalidXmlForMessages)
		mock.When.CheckQueue([anyArgs:]).CallThrough()
		Assert({mock.Eval(AmazonSQS.ApproximateNumberOfMessages, 'testQueue')}
			throws: "ApproximateNumberOfMessages does not exist in response xml: " $
				"<invalidxml>\n\thello world\n</invalidxml>")

		msg = "ERROR: Could not build signed string - unable to get temporary credentials"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg))
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return(AmazonAWS.CredentialErrMsg)
		mock.When.CheckQueue([anyArgs:]).CallThrough()
		Assert({mock.Eval(AmazonSQS.ApproximateNumberOfMessages, 'testQueue')}
			throws: 'ApproximateNumberOfMessages: Problem accessing Message Queue')
		}

	createQueueXml(queue)
		{
return '<CreateQueueResponse>
	<CreateQueueResult>
		<QueueUrl>https://queue.amazonaws.com/123456789012/' $ queue $ '</QueueUrl>
	</CreateQueueResult>
	<ResponseMetadata>
		<RequestId>7a62c49f-347e-4fc4-9331-6e8e7a96aa73</RequestId>
	</ResponseMetadata>
</CreateQueueResponse>'
		}

	createQueueFailedXml: '<invalidResponse></invalidResponse>'

	Test_CreateQueue()
		{
		// valid return values
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return(.createQueueXml('testQueue'))
		mock.When.xmlQueueResult([anyArgs:]).CallThrough()

		Assert(mock.Eval(AmazonSQS.CreateQueue, 'testQueue'))
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/',
			'Action=CreateQueue&QueueName=testQueue', '/')

		// invalid return values
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return(.createQueueXml('createQueueFailed'))
		mock.When.xmlQueueResult([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.CreateQueue, 'testQueue') is: false)

		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return('<invalidxml></invalidxml>')
		mock.When.xmlQueueResult([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.CreateQueue, 'testQueue') is: false)

		// errors
		msg = "ERROR: Could not build signed string - unable to get temporary credentials"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg))
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return(AmazonAWS.CredentialErrMsg)
		mock.When.xmlQueueResult([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.CreateQueue, 'testQueue') is: false)
		}

	Test_processAttributes()
		{
		fn = AmazonSQS.AmazonSQS_processAttributes

		params = [Action: 'SetQueueAttributes'].Set_default(false)
		attributes = []
		fn(params, attributes)
		Assert(params is: [Action: 'SetQueueAttributes'])

		attributes = [NotAnAttribute: 'true']
		Assert({ fn(params, attributes) }
			throws: 'Invalid queue attribute: "NotAnAttribute"')
		Assert(params is: [Action: 'SetQueueAttributes'])

		attributes = [MessageRetentionPeriod: '345600' /*= 4 days in seconds (default)*/]
		fn(params, attributes)
		Assert(params.Action is: 'SetQueueAttributes')
		Assert(params['Attribute.1.Name'] is: 'MessageRetentionPeriod')
		Assert(params['Attribute.1.Value'] is: '345600')

		attributes = [
			MessageRetentionPeriod: '345600', 	/*= 14 days in seconds (our standard)*/
			MaximumMessageSize: 	'1024' 		/*= 1 kb, (256 or 262144 is the default)*/
			]
		fn(params, attributes)
		Assert(params.Action is: 'SetQueueAttributes')
		Assert(params['Attribute.1.Name']  is: 	'MaximumMessageSize')
		Assert(params['Attribute.1.Value'] is: 	'1024')
		Assert(params['Attribute.2.Name']  is: 	'MessageRetentionPeriod')
		Assert(params['Attribute.2.Value'] is: 	'345600')
		}

	deleteQueueXml:
'<DeleteQueueResponse>
	<ResponseMetadata>
		<RequestId>6fde8d1e-52cd-4581-8cd9-c512f4c64223</RequestId>
	</ResponseMetadata>
</DeleteQueueResponse>'
	Test_DeleteQueue()
		{
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return(.deleteQueueXml)
		mock.When.QueueExists([anyArgs:]).Return(true)
		mock.When.xmlQueueResult([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.DeleteQueue, 'testQueue'))
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/testQueue',
			'Action=DeleteQueue', '/testQueue')

		mock = Mock(AmazonSQS)
		mock.When.QueueExists([anyArgs:]).Return(false)
		Assert(mock.Eval(AmazonSQS.DeleteQueue, 'testQueue'))

		mock = Mock(AmazonSQS)
		mock.When.QueueExists([anyArgs:]).Return('UNKNOWN')
		Assert(mock.Eval(AmazonSQS.DeleteQueue, 'testQueue') is: false)

		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return('<invalidResponse>invalid</invalidResponse>')
		mock.When.QueueExists([anyArgs:]).Return(true)
		mock.When.xmlQueueResult([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.DeleteQueue, 'testQueue') is: false)
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/testQueue',
			'Action=DeleteQueue', '/testQueue')

		// errors
		msg = "ERROR: Could not build signed string - unable to get temporary credentials"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg))
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return(AmazonAWS.CredentialErrMsg)
		mock.When.QueueExists([anyArgs:]).Return(true)
		mock.When.xmlQueueResult([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.DeleteQueue, 'testQueue') is: false)
		}

listQueueXml:
'<ListQueuesResponse>
	<ListQueuesResult>
		<QueueUrl>https://sqs.us-east-2.amazonaws.com/123456789012/MyQueue</QueueUrl>
		<QueueUrl>https://sqs.us-east-2.amazonaws.com/123456789012/MyQueue2</QueueUrl>
	</ListQueuesResult>
	<ResponseMetadata>
		<RequestId>725275ae-0b9b-4762-b238-436d7c65a1ac</RequestId>
	</ResponseMetadata>
</ListQueuesResponse>'

listQueueXmlNoMatches:
'<ListQueuesResponse>
	<ListQueuesResult />
	<ResponseMetadata>
		<RequestId>725275ae-0b9b-4762-b238-436d7c65a1ac</RequestId>
	</ResponseMetadata>
</ListQueuesResponse>'

	Test_ListQueues()
		{
		// valid return values
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return(.listQueueXml)
		mock.When.listQueueNodes([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.ListQueues, 'MyQueue') is: #(MyQueue, MyQueue2))
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/',
			'Action=ListQueues&QueueNamePrefix=MyQueue','/')

		// invalid return values
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return('<invalidxml></invalidxml>')
		mock.When.listQueueNodes([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.ListQueues, 'MyQueue') is: #())
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/',
			'Action=ListQueues&QueueNamePrefix=MyQueue','/')

		// valid xml, no matches
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return(.listQueueXmlNoMatches)
		mock.When.listQueueNodes([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.ListQueues, 'MyQueue') is: #())
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/',
			'Action=ListQueues&QueueNamePrefix=MyQueue','/')

		// invalid queue
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return(.listQueueXmlNoMatches)
		mock.When.listQueueNodes([anyArgs:]).CallThrough()
		Assert({mock.Eval(AmazonSQS.ListQueues, 123456789)}
			throws: 'Assert FAILED: AmazonSQS: queueName must be a string')

		// errors
		msg = "ERROR: Could not build signed string - unable to get temporary credentials"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg))
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return(AmazonAWS.CredentialErrMsg)
		mock.When.listQueueNodes([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.ListQueues, 'MyQueue') is: false)
		}

	Test_QueueExists()
		{
		// queue exists
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return(.listQueueXml)
		mock.When.listQueueNodes([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.QueueExists, 'MyQueue'))
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/',
			'Action=ListQueues&QueueNamePrefix=MyQueue','/')

		// queue doesn't exist
		Assert(mock.Eval(AmazonSQS.QueueExists, 'NonExistentQueue') is: false)
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/',
			'Action=ListQueues&QueueNamePrefix=MyQueue','/')

		// invalid return values
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Throw("ERROR: POST FAILED")
		mock.When.listQueueNodes([anyArgs:]).CallThrough()

		Assert(mock.Eval(AmazonSQS.QueueExists, 'MyQueue') is: 'UNKNOWN')
		mock.Verify.post('https://sqs.us-east-1.amazonaws.com/',
			'Action=ListQueues&QueueNamePrefix=MyQueue','/')

		//errors
		msg = "ERROR: Could not build signed string - unable to get temporary credentials"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg))
		mock = Mock(AmazonSQS)
		mock.When.post([anyArgs:]).Return(AmazonAWS.CredentialErrMsg)
		mock.When.listQueueNodes([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonSQS.QueueExists, 'MyQueue') is: 'UNKNOWN')
		}

	Test_EncodeMessageRec()
		{
		test = function (msg, result)
			{
			Assert(result.Split('&')
				equalsSet: AmazonSQS.UrlEncodeValues(msg).Split('&'))
			}
		test(#(), '')
		test(#(Field1: 'value1'), 'Field1=value1')
		test(#(Field1: 'value1', Field2: ''), 'Field1=value1')
		test(#(Field1: 'value1', Field2: 'value2'),
			'Field1=value1&Field2=value2')
		test(#(field1: 'test1', Field2: 'TEst2'),
			'field1=test1&Field2=TEst2')
		test(#(field1: 'test1', Field2: 'TEst2', abc: 'abc', Abcd: 'Abcd', ABCE: 'ABCE'),
			'abc=abc&Abcd=Abcd&ABCE=ABCE&field1=test1&Field2=TEst2')
		test(#(Action: 'CreateQueue', QueueName: 'queue2',
			AWSAccessKeyId: '0A8BDF2G9KCB3ZNKFA82', MessageBody: 'test message'),
			'Action=CreateQueue&AWSAccessKeyId=0A8BDF2G9KCB3ZNKFA82' $
			'&MessageBody=test%20message&QueueName=queue2')
		test(#(Action: 'CreateQueue', QueueName: 'queue2',
			AWSAccessKeyId: '0A8BDF2G9KCB3ZNKFA82', MessageBody: '--- \ntest: message\n'),
			'Action=CreateQueue&AWSAccessKeyId=0A8BDF2G9KCB3ZNKFA82' $
			'&MessageBody=---%20%0Atest%3A%20message%0A&QueueName=queue2')
		}

	Test_invalidChars()
		{
		handleInvalidChars = AmazonSQS.AmazonSQS_handleInvalidChars
		Assert(handleInvalidChars(`123abcABC!@#(),{}[];"',`) is:
			`123abcABC!@#(),{}[];"',`)
		Assert(handleInvalidChars('efg\x02') is: 'efg')
		Assert(handleInvalidChars('hello\r\n\tworld') is: 'hello\r\n\tworld')
		}

	Test_LogFailedResponse_InvalidResponses()
		{
		result = ''
		Assert(AmazonSQS.LogFailedResponse(result) is: false)

		result = 'Test Response - wrong information <?xml version="1.0"?>\r\n' $
			'test blah blah'
		Assert(AmazonSQS.LogFailedResponse(result) is: false)

		// error for 'send'
		result = '<?xml version="1.0"?>\r\n' $
			'<ErrorResponse xmlns="http://queue.amazonaws.com/doc/2008-01-01/">' $
			'<Error><Type>Sender</Type><Code>SignatureDoesNotMatch</Code>' $
			'<Message>The request signature we calculated does not match the ' $
			'signature you provided. Check your AWS Secret Access Key and ' $
			'signing method. Consult the service documentation for details.' $
			'</Message><Detail/></Error>' $
			'<RequestID>131dd6e6-5ea4-453d-b3b6-eccf373253eb</RequestID>' $
			'</ErrorResponse>'
		Assert(AmazonSQS.LogFailedResponse(result) is: false)
		}

	Test_LogFailedResponse_SuccessSend()
		{
		result = '<?xml version="1.0"?>\r\n' $
			'<SendMessageResponse xmlns="http://queue.amazonaws.com/doc/2008-01-01/">' $
			'<SendMessageResult>' $
			'<MD5OfMessageBody>3dbb7943468f08d2f3026b2e1df9b402</MD5OfMessageBody>' $
			'<MessageId>248b8a77-05ad-406c-89fa-00c64bce4efe</MessageId>' $
			'</SendMessageResult><ResponseMetadata>' $
			'<RequestId>e4d56ac9-9142-40c5-aff5-baab44c7fe88</RequestId>' $
			'</ResponseMetadata></SendMessageResponse>'
		Assert(AmazonSQS.LogFailedResponse(result))

		// success for 'send' - from customer
		result = '<?xml version="1.0"?>' $
			'<SendMessageResponse xmlns="http://queue.amazonaws.com/doc/2008-01-01/">' $
			'<SendMessageResult>' $
			'<MD5OfMessageBody>16f4bf33ea33305e93a15047a740cf2c</MD5OfMessageBody>' $
			'<MessageId>832b453d-45ec-498b-803d-ce78bec75de9</MessageId>' $
			'</SendMessageResult><ResponseMetadata>' $
			'<RequestId>31145641-c546-4668-ab43-0d616edf36d1</RequestId>' $
			'</ResponseMetadata></SendMessageResponse>'

		Assert(AmazonSQS.LogFailedResponse(result))
		}

	Test_DeleteMessageBatch()
		{
		failures = .deleteMessageBatchRun([false])
		Assert(failures.failedReceipts isSize: 3)
		Assert(failures.failedMsgs.Empty?())

		// Fails to delete msg3, deletes the other messages
		xmlFail1 = `<?xml version="1.0"?><DeleteMessageBatchResponse>` $
			`<DeleteMessageBatchResult>` $
			`<DeleteMessageBatchResultEntry><Id>msg1</Id>` $
			`</DeleteMessageBatchResultEntry>` $
			`<DeleteMessageBatchResultEntry><Id>msg2</Id>` $ `
			</DeleteMessageBatchResultEntry>` $
			`<BatchResultErrorEntry>` $
				`<Id>msg3</Id>` $
				`<Code>ReceiptHandleIsInvalid</Code>` $
				`<Message>` $
				`The input receipt handle "testText3" is not a valid receipt handle.` $
				`</Message>` $
				`<SenderFault>true</SenderFault>` $
			`</BatchResultErrorEntry>` $
			`</DeleteMessageBatchResult><ResponseMetadata>` $
			`<RequestId>31145641-c546-4668-ab43-0d616edf36d1</RequestId>` $
			`</ResponseMetadata></DeleteMessageBatchResponse>`
		xmlFail2 = `<?xml version="1.0"?><DeleteMessageBatchResponse>` $
			`<DeleteMessageBatchResult>` $
			`<BatchResultErrorEntry>` $
				`<Id>msg1</Id>` $
				`<Code>ReceiptHandleIsInvalid</Code>` $
				`<Message>` $
				`The input receipt handle "testText3" is not a valid receipt handle.` $
				`</Message>` $
				`<SenderFault>true</SenderFault>` $
			`</BatchResultErrorEntry>` $
			`</DeleteMessageBatchResult><ResponseMetadata>` $
			`<RequestId>31145641-c546-4668-ab43-0d616edf36d1</RequestId>` $
			`</ResponseMetadata></DeleteMessageBatchResponse>`

		failures = .deleteMessageBatchRun([XmlParser(xmlFail1), XmlParser(xmlFail2)])
		Assert(failures.failedReceipts isSize: 1)
		Assert(failures.failedMsgs isSize: 1)
		Assert(failures.failedMsgs.First()
			is: `The input receipt handle "testText3" is not a valid receipt handle.`)

		// Fails to delete msg3 on first attempt, succeeds on retry
		xmlMsg3 = `<?xml version="1.0"?><DeleteMessageBatchResponse>` $
			`<DeleteMessageBatchResult>` $
			`<DeleteMessageBatchResultEntry><Id>msg3</Id>` $
			`</DeleteMessageBatchResultEntry>` $
			`</DeleteMessageBatchResult><ResponseMetadata>` $
			`<RequestId>31145641-c546-4668-ab43-0d616edf36d1</RequestId>` $
			`</ResponseMetadata></DeleteMessageBatchResponse>`
		.deleteMessageBatchRun([XmlParser(xmlFail1), XmlParser(xmlMsg3)], succeeds?:)


		// Fails to delete msg3 on first and second attempt, succeeds on final retry
		.deleteMessageBatchRun([XmlParser(xmlFail1), XmlParser(xmlFail2),
			XmlParser(xmlMsg3)], succeeds?:)

		// Successfully deletes all messages
		xmlPass = `<?xml version="1.0"?><DeleteMessageBatchResponse>` $
			`<DeleteMessageBatchResult>` $
			`<DeleteMessageBatchResultEntry><Id>msg1</Id>` $
			`</DeleteMessageBatchResultEntry>` $
			`<DeleteMessageBatchResultEntry><Id>msg2</Id>` $
			`</DeleteMessageBatchResultEntry>` $
			`<DeleteMessageBatchResultEntry><Id>msg3</Id>` $
			`</DeleteMessageBatchResultEntry>` $
			`</DeleteMessageBatchResult><ResponseMetadata>` $
			`<RequestId>31145641-c546-4668-ab43-0d616edf36d1</RequestId>` $
			`</ResponseMetadata></DeleteMessageBatchResponse>`
		.deleteMessageBatchRun([XmlParser(xmlPass)], succeeds?:)

		// No messages to delete
		.deleteMessageBatchRun([XmlParser(xmlPass)], receipts: #(), succeeds?:)
		}

	deleteMessageBatchRun(xmlReturn, receipts = #(testText1, testText2, testText3),
		succeeds? = false)
		{
		mock = Mock(AmazonSQS)
		mock.When.DeleteMessageBatch([anyArgs:]).CallThrough()
		mock.When.deleteBatches([anyArgs:]).CallThrough()
		mock.When.deleteBatch([anyArgs:]).CallThrough()
		mock.When.xmlQueueResult([anyArgs:]).Return(@xmlReturn)
		results = mock.Eval(AmazonSQS.DeleteMessageBatch, '/jlfandjnt', receipts)
		if succeeds?
			{
			Assert(results.Empty?())
			Assert(results hasntMember: #failedReceipts)
			Assert(results hasntMember: #failedReceipts)
			}
		return results
		}
	}