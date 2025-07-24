// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
class
	{
	MaxSendCount: 3

	// NOTE: must hava a slash on the front of the queue name (???)
	Send(queue, msg = '')
		{
		return .send('send', queue, msg)
		}

	// NOTE: must have a slash on the front of the queue name
	// NOTE: messages must be in JSON format
	// processor is passed a message object, must return true or false
	requestsPerReceiveCall: 10
	Receive(queue = '', processor = false, batchDelete? = false)
		{
		if queue is ''
			queue = .Queue

		receipts = Object()
		lastRequestIsEmpty? = true
		for (i = 1; i <= .requestsPerReceiveCall; i++)
			{
			msgs = .send('receive', queue)
			if msgs is AmazonAWS.CredentialErrMsg
				break
			newReceipts = .processIncomingMessages(queue, msgs, processor, batchDelete?)
			lastRequestIsEmpty? = newReceipts.Empty?()
			receipts.Add(@newReceipts)
			}
		if batchDelete?
			.handleBatchDelete(queue, receipts)
		return lastRequestIsEmpty?
		}

	handleBatchDelete(queue, receipts)
		{
		results = .amazonSqsDeleteMessageBatch(queue, receipts.Copy())
		if results.Empty?()
			.log('ReceivedMessage - successfully deleted batch',
				params: Object(Receipts: receipts.Join('\r\n')))
		else if results.Member?(#failedMsgs) and not results.failedMsgs.Empty?()
			results.failedMsgs.Each({ .logErr('', it) })
		}
	amazonSqsDeleteMessageBatch(queue, receipts)
		{
		AmazonSQS.DeleteMessageBatch(queue, receipts)
		}

	MinSleep: 500
	send(type, queue, msg = '')
		{
		result = response = false
		RetryBool(.MaxSendCount, .MinSleep)
			{ |count|
			result = AmazonSQS.SendMessage(type, queue, msg)
			response = .handleSqsResponse(result, type, count)
			// last condition is the block return
			response isnt AmazonAWS.CredentialErrMsg and response isnt false
			}
		if response is AmazonAWS.CredentialErrMsg
			return false
		return response
		}

	defaultErrorThreshold: 10
	receiveErrorThreshold: 20 // "receive" requests need higher threshold (higher volume)
	handleSqsResponse(result, type, count = 0)
		{
		if result is AmazonAWS.CredentialErrMsg
			return result

		threshold = type is 'receive' ? .receiveErrorThreshold : .defaultErrorThreshold
		if '' is (response = result.AfterFirst('<?xml version='))
			{
			if count is .MaxSendCount
				.processErr(type, threshold, result, result)
			return false
			}
		if not response.Has?(type.Capitalize() $ "MessageResponse")
			{
			if count is .MaxSendCount
				.processErr(type, threshold, response, result)
			return false
			}
		return AmazonSQSParseMessages(type, response)
		}

	processErr(type, threshold, response, result)
		{
		ConnectionErrorHandler.Process(response, .ConnectionErrStr,
			" (" $ type $ ")", 'AmazonSQSManager', :threshold)
		if Type(result) is 'Except'
			.log('send failed - ' $ type,
				params: LogFormatEntry(result.Callstack()[0].locals, maxStrSize: 1000))
		}

	processIncomingMessages(queue, msgs, processor, batchDelete?)
		{
		if msgs is false or msgs.Empty?()
			return #()
		.log('ReceivedMessage - received ' $ Display(msgs.Size()) $
			' messages', params: msgs)

		receipts = Object()
		for msg in msgs
			{
			.processEachMsg(msg, processor)
			receipts.Add(msg.receipt)
			if not batchDelete? and .send('delete', queue, msg.receipt) is true
				.log('ReceivedMessage - successfully deleted',
					params: Object('Receipt': msg.receipt))
			}
		return receipts
		}

	processEachMsg(msg, processor)
		{
		try
			{
			rec = Json.Decode(msg.body)
			if not Object?(rec)
				.logErr('invalid format for Message Body, this message will be skipped',
					msg)
			else if false is processor(rec)
				.logErr('processor failed', msg)
			}
		catch(e)
			.logErr(e, msg)
		}

	logErr(e, msg)
		{
		if e isnt ''
			e = ' - ' $ e // only $ concatenation keeps the exception information
		SuneidoLog('ERROR: AmazonSQSManager - cannot process message' $ e, params: msg)
		}

	log(s, params)
		{
		logName = GetContributions('LogPaths').GetDefault('amazonsqs', 'AmazonSQSManager')
		formatedParams = LogFormatEntry(params, 1000 /*= max size for log*/)
		Rlog(logName,
			Opt(s is false ? '' : s , ', ') $ 'Message: ' $ Display(formatedParams))
		}

	ConnectionErrStr: #("couldn't connect to host",
		"Could not resolve host",
		"We encountered an internal error. Please try again.",
		"unknown address: queue.amazonaws.com",
		"The action ReceiveMessage is not valid for this endpoint.")
	}
