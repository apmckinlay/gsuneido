// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	test(type, resp, expected)
		{
		result = AmazonSQSParseMessages(type, resp)
		Assert(result is: expected)
		}
	Test_delete()
		{
		.test('delete', '', true)
		}
	Test_send()
		{
		.test('send', '', #(md5: '', id: ''))
		.test('send',
			'<other>
			<MD5OfMessageBody>md5</MD5OfMessageBody>
			<MessageId>123</MessageId>
			</other>',
			#(md5: 'md5', id: '123'))
		}
	Test_receive_missing_start()
		{
		.test('receive', '
			<MessageId>123</MessageId>
			<MD5OfBody>md5</MD5OfBody>
			<ReceiptHandle>foo</ReceiptHandle>
			<Body>now is the time</Body>
			</Message>',
			#())
		}
	Test_receive_missing_end()
		{
		.test('receive', '
			<Message>
			<MessageId>123</MessageId>
			<MD5OfBody>md5</MD5OfBody>
			<ReceiptHandle>foo</ReceiptHandle>
			<Body>now is the time</Body>',
			#((id: '123', md5: 'md5', receipt: 'foo', body: 'now is the time')))
		}
	Test_receive_single()
		{
		.test('receive', '
			<Message>
			<MessageId>123</MessageId>
			<MD5OfBody>md5</MD5OfBody>
			<ReceiptHandle>foo</ReceiptHandle>
			<Body>now is the time</Body>
			</Message>',
			#((id: '123', md5: 'md5', receipt: 'foo', body: 'now is the time')))
		}
	Test_receive_multiple()
		{
		.test('receive', '
			<Message>
			<MessageId>123</MessageId>
			<MD5OfBody>md5</MD5OfBody>
			<ReceiptHandle>foo</ReceiptHandle>
			<Body>now is the time</Body>
			</Message>
			<Message>
			<MessageId>456</MessageId>
			<MD5OfBody>md5too</MD5OfBody>
			<ReceiptHandle>bar</ReceiptHandle>
			<Body>one &amp; two</Body>
			</Message>',
			#((id: '123', md5: 'md5', receipt: 'foo', body: 'now is the time'),
				(id: '456', md5: 'md5too', receipt: 'bar', body: 'one & two')))
		}
	}