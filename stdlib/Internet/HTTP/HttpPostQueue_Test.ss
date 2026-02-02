// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		queueTable = .MakeTable(HttpPostQueue.QueueTableSchema)
		queueCl = HttpPostQueue
			{
			New()
				{ .SendLog = Object() }
			HttpPostQueue_httpPost(record)
				{
				if record.address.Has?('httppost_queue_testing')
					{
					.SendLog.Add(Object(address: record.address,
						contents: record.contents))
					if record.address is 'httppost_queue_testing3'
						return "post failed"
					}
				else if record.address.Has?('httppost_queue_error')
					throw "program error"
				return ""
				}
			}
		.postQueue = new queueCl
		.postQueue.QueueTableName = queueTable
		}

	Test_main()
		{
		.testAddToQueue()
		.testErrorOnFirstPost()
		.testSendFromQueue()
		}

	testAddToQueue()
		{
		addr = contents = 'httppost_queue_error1'
		.postQueue.AddToQueue(addr, contents)
		rec = .queuedRec(addr, contents)
		Assert(rec isnt: false)

		addr = contents = 'httppost_queue_testing1'
		.postQueue.AddToQueue(addr, contents)
		rec = .queuedRec(addr, contents)
		Assert(rec isnt: false)

		addr = contents = 'httppost_queue_testing2'
		.postQueue.AddToQueue(addr, contents)
		rec = .queuedRec(addr, contents)
		Assert(rec isnt: false)

		addr = contents = 'httppost_queue_error2'
		.postQueue.AddToQueue(addr, contents)
		rec = .queuedRec(addr, contents)
		Assert(rec isnt: false)

		msg = "ERROR: HttpPostSendQueue.HttpPost: post failed"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg, msg, msg))
		addr = contents = 'httppost_queue_testing3'
		.postQueue.AddToQueue(addr, contents)
		rec = .queuedRec(addr, contents)
		Assert(rec isnt: false)
		}

	testErrorOnFirstPost()
		{
		.postQueue.Send()
		size = .postQueue.SendLog.Size()
		Assert(size is: 0)
		QueryDo('delete ' $ .postQueue.QueueTableName $
			' where address is "httppost_queue_error1"')
		}

	testSendFromQueue()
		{
		.postQueue.Send()
		size = .postQueue.SendLog.Size()
		Assert(size is: 3)
		// First two posts succeed
		// Third post is program error
		// Fourth post fails with content
		for (i = 0; i < size; ++i)
			{
			queueStr = "httppost_queue_testing" $ (i + 1)
			Assert(.postQueue.SendLog[i].address is: queueStr)
			Assert(.postQueue.SendLog[i].contents is: queueStr)
			}
		Assert(.queuedRec("httppost_queue_testing1", "httppost_queue_testing1") is: false)
		Assert(.queuedRec("httppost_queue_testing2", "httppost_queue_testing2") is: false)
		Assert(.queuedRec("httppost_queue_testing3", "httppost_queue_testing3")
			isnt: false)
		Assert(.queuedRec("httppost_queue_error2", "httppost_queue_error2") isnt: false)
		}

	queuedRec(addr, contents)
		{
		return Query1(.postQueue.QueueTableName, address: addr, :contents)
		}
	}
