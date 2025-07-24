// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
AmazonAWS
	{
	Host(region = false)
		{
		region = region is false ? .region : region
		return 'sqs.' $ region $ '.amazonaws.com'
		}

	ContentType()
		{
		return 'application/x-www-form-urlencoded; charset=utf-8'
		}

	Service()
		{
		return 'sqs'
		}

	CanonicalQueryString(unused)
		{
		return ''
		}

	PayloadHash(params)
		{
		return Sha256(params).ToHex()
		}

	maxToReceive: "10"
	SendMessage(action, queue, body, extraParams = #())
		{
		params = []
		if action is 'send'
			params = [Action: 'SendMessage', MessageBody: .handleInvalidChars(body)]
		else if action is 'receive'
			{
			maxToReceive = .maxToReceive
			if Object?(body)
				maxToReceive = body.GetDefault('maxToReceive', .maxToReceive)
			params = [Action: 'ReceiveMessage', MaxNumberOfMessages: maxToReceive]
			}
		else if action is 'delete'
			params = [Action: 'DeleteMessage', ReceiptHandle: .handleInvalidChars(body)]
		else
			return 'ERROR: Invalid Action for SendMessage: ' $ action
		params = params.Merge(extraParams)
		try
			return .makeRequest(params, queue)
		catch (err)
			{ return err }
		}

	handleInvalidChars(body)
		{
		// SEE: http://docs.aws.amazon.com/AWSSimpleQueueService/
		// 			latest/APIReference/Query_QuerySendMessage.html
		// NOTE: Suneido only supports 8 bit characters (no unicode)
		// WARNING: this just deletes invalid characters
		return body.Tr('^\t\r\n\x20-\xff')
		}

	region: 'us-east-1'
	makeRequest(params, path)
		{
		body = AmazonAWS.UrlEncodeValues(params)
		url = 'https://' $ .Host() $ path
		return .post(url, body, path)
		}

	post(url, body, path)
		{
		extraHeaderInfo = Object(X_Amz_Security_Token: .SecurityToken())
		if .CredentialErrMsg is hdr = AmazonV4Signing(
			this, 'POST', .region, body, path, extraHeaderInfo).AuthorizationHeader()
			return .CredentialErrMsg
		return Https.Post(url, body, header: hdr)
		}

	// returns true even if the queue already exists
	// returns false if creation failed
	CreateQueue(queueName, attributes = #())
		{
		if not attributes.Member?('SqsManagedSseEnabled')
			{
			attributes = attributes.Copy()
			attributes.SqsManagedSseEnabled = 'false'
			}
		params = [Action: 'CreateQueue', QueueName: queueName]
		.processAttributes(params, attributes)
		if false is xml = .xmlQueueResult(params, '/')
			return false
		return xml.createqueueresult.queueurl.Text().Suffix?('/' $ queueName)
		}

	CreateFifoQueue(queueName, attributes = #())
		{
		if not queueName.Suffix?('.fifo')
			queueName $= '.fifo'
		attributes = Object('FifoQueue': 'true', 'ContentBasedDeduplication': 'true').
			Merge(attributes)
		return .CreateQueue(queueName, attributes)
		}

	/* See https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/
		API_SetQueueAttributes.html

		When you change a queue's attributes, the change can take up to 60 seconds for
		most of the attributes to propagate throughout the Amazon SQS system.
		Changes made to the MessageRetentionPeriod attribute can take up to 15 minutes. */
	SetQueueAttributes(queueName, attributes)
		{
		.processAttributes(params = [Action: 'SetQueueAttributes'], attributes)
		xml = .xmlQueueResult(params, '/' $ queueName)
		return .queueAttributesResponse(xml)
		}

	queueAttributes: #(Policy, VisibilityTimeout, MaximumMessageSize,
		MessageRetentionPeriod, ApproximateNumberOfMessages,
		ApproximateNumberOfMessagesNotVisible, CreatedTimestamp,
		LastModifiedTimestamp, QueueArn, ApproximateNumberOfMessagesDelayed,
		DelaySeconds, ReceiveMessageWaitTimeSeconds, RedrivePolicy, FifoQueue,
		ContentBasedDeduplication, KmsMasterKeyId, KmsDataKeyReusePeriodSeconds,
		SqsManagedSseEnabled)
	processAttributes(params, attributes)
		{
		count = 1 // Starting at 1 is based on AmazonSQS example
		for name in attributes.Members().Sort!()
			{
			value = attributes[name]
			if not .queueAttributes.Has?(name)
				throw 'Invalid queue attribute: ' $ Display(name)
			params['Attribute.' $ count $ '.Name'] = name
			params['Attribute.' $ count $ '.Value'] = value
			count++
			}
		}

	QueueExists(queueName)
		{
		if false is queueNodes = .listQueueNodes(queueName)
			return "UNKNOWN"
		return queueNodes.Any?({ it.Text().Suffix?('/' $ queueName) })
		}
	listQueueNodes(queueName)
		{
		Assert(String?(queueName), "AmazonSQS: queueName must be a string")
		params = [Action: 'ListQueues', QueueNamePrefix: queueName]
		if false is xml = .xmlQueueResult(params, '/')
			return false
		return xml.listqueuesresult.queueurl.List()
		}
	// WARNING: this only return the first 1000 records
	ListQueues(prefix)
		{
		if false is queueNodes = .listQueueNodes(prefix)
			return false
		return queueNodes.Map({ it.Text().AfterLast('/') })
		}
	// returns true if queue didn't exist
	// if delete failed for ANY REASON, return false
	// (including if .QueueExists() failed)
	DeleteQueue(queueName)
		{
		if false is exists = .QueueExists(queueName)
			return true
		if exists is "UNKNOWN"
			return false

		params = [Action: 'DeleteQueue']
		if false is xml = .xmlQueueResult(params, '/' $ queueName)
			return false
		return xml.responsemetadata.requestid.Text() isnt ''
		}

	DeleteMessageBatch(queueName, receipts)
		{
		RetryBool(AmazonSQSManager.MaxSendCount, AmazonSQSManager.MinSleep)
			{ | count |
			// After the first attempt, we are dealing with failed receipts
			if count isnt 1
				receipts = receipts.failedReceipts
			receipts = .deleteBatches(queueName, receipts)
			receipts.Empty?()
			}
		// If this fails every attempt, receipts will contain a list of failure messages
		return receipts
		}

	batchAttrib: 'DeleteMessageBatchRequestEntry.'
	deleteBatches(queueName, receipts)
		{
		batchSize = 10 // Max number of batch size, as per AmazonSQS specifications
		iterations = receipts.Size() / batchSize
		failures = Object(failedReceipts: Object(), failedMsgs: Object())
		for (i = 0; i < iterations;)
			{
			batch = [Action: 'DeleteMessageBatch']
			id = 1
			receipts[i * batchSize .. ++i * batchSize].Each(
				{
				batch[.batchAttrib $ id $ '.Id'] = 'msg' $ id
				batch[.batchAttrib $ id++ $ '.ReceiptHandle'] = .handleInvalidChars(it)
				})
			.deleteBatch(queueName, batch, failures)
			}
		return failures.failedReceipts.Empty?()
			? #()
			: failures
		}

	deleteBatch(queueName, batch, failures)
		{
		basePath = Object(#deletemessagebatchresponse, #deletemessagebatchresult,
			#batchresulterrorentry)
		if false is xml = .xmlQueueResult(batch, queueName)
			failures.failedReceipts =
				batch.DeleteIf({ it.Suffix?('.Id') or it is 'Action' }).Flatten()
		else if not XmlFind.All(xml, basePath).Empty?()
			{
			XmlFind.All(xml, basePath.Copy().Add(#id)).FlatMap(#Children).Each(
				{
				id = batch.Find(it.Text()).AfterFirst('.').BeforeFirst('.')
				failures.failedReceipts.Add(batch[.batchAttrib $ id $ '.ReceiptHandle'])
				})
			XmlFind.All(xml, basePath.Copy().Add(#message)).FlatMap(#Children).Each(
				{
				failures.failedMsgs.AddUnique(it.Text())
				})
			}
		}

	CheckQueue(queueName, attributes = 'All')
		{
		params = [Action: 'GetQueueAttributes', 'AttributeName.1': attributes]
		xml = .xmlQueueResult(params, '/' $ queueName)
		return .queueAttributesResponse(xml)
		}

	queueAttributesResponse(xml)
		{
		return Instance?(xml) and xml.Base?(XmlNode) and xml.Name() isnt 'errorresponse'
			? xml
			: false
		}
	ApproximateNumberOfMessages(queue)
		{
		if false is xmlNode = .CheckQueue(queue, 'ApproximateNumberOfMessages')
			throw 'ApproximateNumberOfMessages: Problem accessing Message Queue'

		if false is num = XmlFind.First(xmlNode,
			#(getqueueattributesresponse, getqueueattributesresult, attribute, value))
			throw 'ApproximateNumberOfMessages does not exist in response xml: ' $ xmlNode
		return Number(num.Text())
		}

	xmlQueueResult(params, path)
		{
		if false is result = .sqsQueueManager(params, path)
			return false

		try
			result = XmlParser(result)
		catch (err)
			{
			SuneidoLog('ERROR: ' $ err)
			return false
			}
		return result
		}

	sqsQueueManager(params, path)
		{
		try
			{
			if .CredentialErrMsg is result = .makeRequest(params, path)
				{
				SuneidoLog('ERROR: Could not build signed string - ' $ .CredentialErrMsg)
				return false
				}
			return result
			}
		catch (err)
			{
			SuneidoLog('Amazon SQS: ' $ err)
			return false
			}
		}

	LogFailedResponse(result)
		{
		if result is .CredentialErrMsg
			return false
		if '' is response = result.AfterFirst('<?xml version=')
			{
			SuneidoLog('Amazon SQS: ' $ result, calls:)
			return false
			}
		if not response.Has?("MessageResponse")
			{
			SuneidoLog('Amazon SQS: ' $ response, calls:)
			return false
			}
		return true
		}
	}