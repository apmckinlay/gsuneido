// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
class
	{
	QueueTableName: 'post_queue'
	QueueTableSchema: ' (num, address, header, contents) key (num)'

	EnsureQueueTable()
		{
		Database('ensure ' $ .QueueTableName $ .QueueTableSchema)
		}

	AddToQueue(addr, contents = '', header = #(), t = false)
		{
		//FIXME don't do an ensure every time
		.EnsureQueueTable()
		DoWithTran(:t, update:)
			{ |tran|
			tran.QueryOutput(.QueueTableName,
				[num: Timestamp(), address: addr, :header, :contents])
			}
		}

	Send()
		{
		//FIXME don't do an ensure every time
		.EnsureQueueTable()
		counter = 0
		for record in QueryAll(.QueueTableName $ ' sort num')
			{
			try
				{
				++counter
				result = .httpPost(record)
				// NOTE: assumes non-blank result is failure
				if result.Blank?()
					.removeFromQueue(record)
				else
					SuneidoLog('ERROR: HttpPostSendQueue.HttpPost: ' $ result)
				}
			catch (err)
				{
				SuneidoLog('ERRATIC: HttpPostSendQueue: ' $ err)
				// first attempt failed; assume the service is down and stop sending
				if counter is 1
					return
				}
			}
		}

	httpPost(rec) // overridden by test
		{
		return Http.Post('http://' $ rec.address, rec.contents, header: rec.header)
		}

	removeFromQueue(record)
		{
		RetryTransaction()
			{ |t|
			t.QueryDo('delete ' $ .QueueTableName $
				' where num is ' $ Display(record.num))
			}
		}
	}
