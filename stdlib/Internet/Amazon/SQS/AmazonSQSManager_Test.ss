// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_processEachMsg()
		{
		mock = Mock()
		logErr = "AmazonSQSManager_logErr"
		mock.When[logErr]([anyArgs:]).Return(false)
		f = AmazonSQSManager.AmazonSQSManager_processEachMsg

		rec = Record(body: '{"abc": 123}')
		mock.Eval(f, rec, {|unused| true })
		mock.Verify.Never()[logErr]([anyArgs:])

		rec = Record()
		mock.Eval(f, rec, { throw 'should not get here' })
		mock.Verify[logErr]('Invalid Json format: unexpected end of string', rec)

		rec = Record(body: 'wrong message')
		mock.Eval(f, rec, { throw 'should not get here' })
		mock.Verify[logErr]('Invalid Json format: unexpected: wrong', rec)

		rec = Record(body: '{"abc": 123}')
		mock.Eval(f, rec, {|unused| false })
		mock.Verify[logErr]('processor failed', rec)

		rec = Record(body: '{"abc": 123}')
		mock.Eval(f, rec, {|unused| throw 'process failed mid run' })
		mock.Verify[logErr]('process failed mid run', rec)
		}

	Test_handleBatchDelete()
		{
		mock = Mock(AmazonSQSManager)
		mock.When.amazonSqsDeleteMessageBatch([anyArgs:]).Return(#(),
			Object(failedMsgs: [
				'msg1 failed to remove, bad request'
				'msg3 failed to remove, bad request'
				]))
		mock.When.handleBatchDelete([anyArgs:]).CallThrough()
		mock.When.log([anyArgs:]).Return('')
		mock.When.logErr([anyArgs:]).Return('')

		mock.handleBatchDelete(true, ['msg1', 'msg2', 'msg3'])
		mock.Verify.log('ReceivedMessage - successfully deleted batch',
			params: Object(Receipts: 'msg1\r\nmsg2\r\nmsg3'))

		mock.handleBatchDelete(false, ['msg1', 'msg2', 'msg3'])
		mock.Verify.logErr('', 'msg1 failed to remove, bad request')
		mock.Verify.logErr('', 'msg3 failed to remove, bad request')
		}

	Test_Receive()
		{
		queue = .TempName()
		count = 0
		processor = { |rec|
			Assert(rec.id is: count++)
			true
			}

		suneidologWatch = .WatchTable('suneidolog')
		table = .MakeTable('(amazonsqstestrlogcount) key()')
		QueryOutput(table, Record(amazonsqstestrlogcount: 0))
		_amazonsqstestrlogcount = table
		cl = AmazonSQSManager
			{
			AmazonSQSManager_handleBatchDelete(@unused) { return '' }
			AmazonSQSManager_send(@unused)
				{
				msg = _amazonsqstestmessages[0]
				_amazonsqstestmessages.Delete(0)
				return msg
				}
			AmazonSQSManager_log(@unused)
				{
				QueryApply1(_amazonsqstestrlogcount)
					{
					++it.amazonsqstestrlogcount
					it.Update()
					}
				}
			}
		fn = cl.Receive

		count = 0
		_amazonsqstestmessages = .makeMsgsOb(0)
		Assert(fn(queue, processor, batchDelete?:))
		Assert(Query1(table).amazonsqstestrlogcount is: 0)
		Assert(.GetWatchTable(suneidologWatch) is: #())

		count = 0
		_amazonsqstestmessages = .makeMsgsOb(5)
		Assert(fn(queue, processor, batchDelete?:))
		Assert(Query1(table).amazonsqstestrlogcount is: 5)
		Assert(.GetWatchTable(suneidologWatch) is: #())

		count = 0
		QueryDo('delete ' $ table)
		QueryOutput(table, Record(amazonsqstestrlogcount: 0))
		_amazonsqstestmessages = .makeMsgsOb(9)
		Assert(fn(queue, processor, batchDelete?:) )
		Assert(Query1(table).amazonsqstestrlogcount is: 9)
		Assert(.GetWatchTable(suneidologWatch) is: #())

		count = 0
		QueryDo('delete ' $ table)
		QueryOutput(table, Record(amazonsqstestrlogcount: 0))
		_amazonsqstestmessages = .makeMsgsOb(10)
		Assert(fn(queue, processor, batchDelete?:) is: false)
		Assert(Query1(table).amazonsqstestrlogcount is: 10)
		Assert(.GetWatchTable(suneidologWatch) is: #())
		}

	makeMsgsOb(loop)
		{
		msgsOb = Object()
		start = 0
		for ..loop
			{
			n = Random(9) + 1
			msgs = Object()
			for i in ..n
				msgs.Add([receipt: start + i, body: Json.Encode([id: start + i])])
			msgsOb.Add(msgs)
			start += n
			}
		while msgsOb.Size() < 10
			msgsOb.Add(#())
		return msgsOb
		}
	}