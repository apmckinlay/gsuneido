// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
AmazonAWS
	{
	DynamoDBExpirySeconds: 60

	Host(region = false)
		{
		region = region is false ? .region : region
		return 'dynamodb.' $ region $ '.amazonaws.com'
		}

	ContentType()
		{
		return 'application/x-amz-json-1.0'
		}

	Service()
		{
		return 'dynamodb'
		}

	CanonicalQueryString(unused)
		{
		return ''
		}

	PayloadHash(params)
		{
		return Sha256(params).ToHex()
		}

	region: 'us-east-1'
	makeRequest(target, content = '', expectedResponse = '200',
		fullResponse? = false, extraHeaders = #())
		{
		url = 'https://' $ .Host()
		extraHeaders = extraHeaders.Copy()
		extraHeaders.X_Amz_Content_Sha256 = .PayloadHash(content)
		extraHeaders.X_Amz_Security_Token = .SecurityToken()
		extraHeaders.X_Amz_Target = 'DynamoDB_20120810.' $ target
		if false is header = .signRequest('POST', .region, content, extraHeaders)
			return false

		return .throttle('POST', :expectedResponse, :fullResponse?)
			{
			.https('POST', url, header, content)
			}
		}

	signRequest(call, region, params, extraHeaders)
		{
		return AmazonV4Signing(this, call, region, params,
			extraHeaderInfo: extraHeaders).AuthorizationHeader()
		}

	https(call, url, header = #(), content = '')
		{
		result = Https(call, url, :content, :header)
		return result
		}

	DescribeTable(tablename)
		{
		params = [TableName: tablename]
		return .makeRequest('DescribeTable', Json.Encode(params))
		}

//	Query(table)
//		{
//
//		}

	throttle(action, block, expectedResponse = '200', fullResponse? = false)
		{
		resultOb = false
		try
			.throttleRetry()
				{
				resultOb = (block)()
				params = Locals(3) /*= call levels up to throttle */
				status = .checkResponseAndLog(action, resultOb, expectedResponse, params)
				return status is true
					? fullResponse?
						? resultOb
						: resultOb.GetDefault(#content, '').Trim()
					: false
				}
		catch (err, "Retry failed")
			{
			detail = resultOb is false ? ''
				: resultOb.header $ '\r\n\r\n' $ resultOb.GetDefault(#content, '')
			.addToLog(action, err, detail, params)
			return false
			}
		return false
		}

	// extrated for tests
	throttleRetry(block)
		{
		Retry(block, maxRetries: 3, minDelayMs: 100,
			retryException: 'Bad HTTP Status Code (503)')
		}

	checkResponseAndLog(action, resultOb, expectedResponse = '200', params = '')
		{
		code = Http.ResponseCode(resultOb.header)
		if code is expectedResponse or
			Object?(expectedResponse) and expectedResponse.Has?(code)
			return true
		msg = 'Bad HTTP Status Code (' $ code $ ')'
		// 503 is either 'SlowDown' or 'ServiceUnavailable' - in both cases
		// amazon requsts reducing the request rate (switch to Exponential Fallback)
		if code is '503'
			throw msg
		detail = resultOb.header $ '\r\n\r\n' $ resultOb.GetDefault(#content, '')
		.addToLog(action, msg, detail, params)
		return false
		}

	addToLog(action, msg, detail = '', params = '')
		{
		if Object?(params)
			params = params.Copy().Delete(#resultOb, #result)
		if .exceptionInsteadOfLog?()
			throw 'AmazonDynamoDB - ' $ action $ ': ' $ msg
		SuneidoLog('ERRATIC: AmazonDynamoDB - ' $ action $ ': ' $ msg, calls:,
			:params, switch_prefix_limit: 5)
		CreateDir('logs')
		AddFile('logs/amazonDynamoDB.log', Display(Date())[1..] $ ', ' $
			action $ ' > ' $ msg $ Opt('\r\n', detail.Trim()) $ '\r\n\r\n')
		}

	exceptionInsteadOfLog?()
		{
		try
			return _exceptionOnFailure
		return false
		}
	}
