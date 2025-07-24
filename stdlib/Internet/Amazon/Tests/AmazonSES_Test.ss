// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getEndPoint()
		{
		sesCl = new AmazonSES
		endPoint = sesCl.AmazonSES_getEndPoint

		// table does not exist
		sesCl.AmazonSES_endPointTable = "amazonses_test_endpoint"
		Assert(endPoint() is: 'email.us-east-1.amazonaws.com')

		// table exists but no rec
		sesCl.AmazonSES_endPointTable = .MakeTable('(aws_ses_endpoint) key ()')
		Assert(endPoint() is: 'email.us-east-1.amazonaws.com')

		// table exists and has endpoint record
		aws_ses_endpoint = 'sample_endpoint_in_table'
		QueryOutput(sesCl.AmazonSES_endPointTable, Record(:aws_ses_endpoint))
		Assert(endPoint() is: aws_ses_endpoint)

		// test OutputEndPointRecord method
		aws_ses_endpoint = 'another_sample_endpoint_in_table'
		sesCl.OutputEndPointRecord(aws_ses_endpoint)
		Assert(endPoint() is: aws_ses_endpoint)

		// deleting endpoint rec
		QueryDo("delete " $ sesCl.AmazonSES_endPointTable)
		Assert(endPoint() is: 'email.us-east-1.amazonaws.com')

		// drop endpoint table
		Database("drop " $ sesCl.AmazonSES_endPointTable)
		Assert(endPoint() is: 'email.us-east-1.amazonaws.com')
		}

	Test_doRequest()
		{
		mock = Mock(AmazonSES)
		mock.When.doRequest([anyArgs:]).CallThrough()

		count = 0
		mock.When.postRequest([anyArgs:]).Do(
			{
			if count++ < 2
				throw 'curl: (35) gnutls_handshake() failed: testing'
			Object(header: 'HTTP/1.1 200 OK')
			})
		Assert(mock.doRequest(#()) is: #(header: "HTTP/1.1 200 OK", code: "200"))

		mock.When.postRequest([anyArgs:]).
			Throw('curl: (35) gnutls_handshake() failed: testing')
		result = mock.doRequest(#())
		Assert(result.header is: result.content)
		Assert(result.content
			has: 'Retry failed - too many retries, last error: Bad HTTP Status Code')

		mock.When.postRequest([anyArgs:]).Throw('Bad HTTP Status Code')
		result = mock.doRequest(#())
		Assert(result.header is: result.content)
		Assert(result.content
			has: 'Retry failed - too many retries, last error: Bad HTTP Status Code')

		mock.When.postRequest([anyArgs:]).Throw('test error')
		Assert({ mock.doRequest(#()) } throws: 'test error')

		mock.When.postRequest([anyArgs:]).Return(Object(header: 'HTTP/1.1 200 OK'))
		Assert(mock.doRequest(#()) is: #(header: "HTTP/1.1 200 OK", code: "200"))

		mock = Mock(AmazonSES)
		mock.When.doRequest([anyArgs:]).CallThrough()
		mock.When.postRequest([anyArgs:]).Do(
			{
			Object(header: 'HTTP/1.1 502 server unavailable', content: 'server failed')
			})
		result = mock.doRequest(#())
		Assert(result.code is: '502')
		Assert(result.header is: 'HTTP/1.1 502 server unavailable')
		Assert(result.content
			has: 'server failed\r\n\r\n' $
				'Retry failed - too many retries, last error: Bad HTTP Status Code')
		}
	}