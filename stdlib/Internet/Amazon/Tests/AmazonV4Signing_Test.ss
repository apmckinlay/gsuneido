// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	sampleAWS: class
		{
		Host() { return "sampleAWSHost" }
		ContentType() { return  'plain-text' }
		Service() { return 'sample_service' }
		AccessKey() { return "ACCESS" }
		SecretKey() { return "secret" }
		SecurityToken() { return 'token' }
		CanonicalQueryString(params) { return params }
		PayloadHash(params) { return Sha256(params).ToHex() }
		}

	v4SignCl: AmazonV4Signing
			{
			AmazonV4Signing_amazonDate()
				{
				return #20180702.125901555
				}
			}

	Test_getCredentialScope()
		{
		ascl = AmazonV4Signing
			{ AmazonV4Signing_amazonDate() { return #20180702.125901555 } }
		cl = ascl(.sampleAWS, "POST", "us-east1", "TEST_PARAMS")
		cl.AmazonV4Signing_date = #20180701.125601555
		result = cl.AmazonV4Signing_getCredentialScope()
		Assert(result is: "20180701/us-east1/sample_service/aws4_request")
		}

	Test_buildAuthorizationHeader()
		{
		ascl = AmazonV4Signing
			{ AmazonV4Signing_amazonDate() { return #20180702.125901555 } }
		cl = ascl(.sampleAWS, "POST", "us-east1", "MORE_PARAMS")
		cl.AmazonV4Signing_date = #20180701.125601556
		result = cl.AmazonV4Signing_buildAuthorizationHeader("sample_signature")
		Assert(result
			is: 'AWS4-HMAC-SHA256 Credential=ACCESS/20180701/us-east1/sample_service/' $
				'aws4_request, SignedHeaders=content-type;host;x-amz-date, ' $
				'Signature=sample_signature')
		}

	Test_createCanonicalRequest()
		{
		params = AmazonAWS.UrlEncodeValues([param1: "SAMPLE_PARAMS"])
		cl = (.v4SignCl)(.sampleAWS, "POST", "us-west1", params)
		result = cl.AmazonV4Signing_createCanonicalRequest()
		Assert(result is: 'POST\n/\n' $
			'param1=SAMPLE_PARAMS\ncontent-type:plain-text\nhost:sampleAWSHost\n' $
			'x-amz-date:20180702T125901Z\n\n' $
			'content-type;host;x-amz-date\n' $
			'8f8daea5b1f0ce47ee02ab32ce6cec0948cb4f01a258ccce3bdebceca4e09817')
		}

	Test_getSignatureKey()
		{
		// WARNING - not testing actual output of the method, just the flow since sign
		// method is overridden to do simple concatenation
		testCl = AmazonV4Signing {
			AmazonV4Signing_sign(key, msg) { return key $ msg }
			AmazonV4Signing_amazonDate() { return #20180702.125901555 }
			}
		cl = testCl(.sampleAWS, "GET", "us-west1", "TEST_CONTENT")
		cl.AmazonV4Signing_date = #20180704.125901555
		result = cl.AmazonV4Signing_getSignatureKey()
		Assert(result is: 'AWS4secret20180704us-west1sample_serviceaws4_request')
		}

	sqsHeader: #(Host: "sqs.us-east-1.amazonaws.com", X_Amz_Date: "20180702T125901Z",
		Content_Type: "application/x-www-form-urlencoded; charset=utf-8",
		X_Amz_Security_Token: "token",
		Authorization: "AWS4-HMAC-SHA256" $
			" Credential=ACCESS/20180702/us-east-1/sqs/aws4_request," $
			" SignedHeaders=content-type;host;x-amz-date;x-amz-security-token, " $
			"Signature=f568bd1ace29fe1430ace9d75a248944023a4b83a82dff56a9b2b4552b6a52ce")
	Test_AuthorizationHeader_SQS()
		{
		sqsCl = AmazonSQS
			{
			AccessKey() { return "ACCESS" }
			SecretKey() { return "secret" }
			}

		body = 'Action=ReceiveMessage&MaxNumberOfMessages=10'
		path = '/fredsQueue'
		hdr = (.v4SignCl)(sqsCl, 'POST', 'us-east-1', body, path,
			#(X_Amz_Security_Token: 'token')).AuthorizationHeader()
		Assert(hdr is: .sqsHeader)
		}

	s3Header: #(Host: "s3.us-east-1.amazonaws.com", X_Amz_Date: "20180702T125901Z",
		Content_Type: "multipart/form-data",
		X_Amz_Security_Token: "token",
		X_Amz_Content_Sha256: "UNSIGNED-PAYLOAD",
		Authorization: "AWS4-HMAC-SHA256 " $
			"Credential=ACCESS/20180702/us-east-1/s3/aws4_request, " $
			"SignedHeaders=content-type;host;x-amz-content-sha256;x-amz-date;" $
			"x-amz-security-token, " $
			"Signature=1a60531168437259ced864d2f47e888f338df5a4460f38ebb5091da75e226fe4")
	Test_AuthorizationHeader_S3()
		{
		sqsCl = AmazonS3
			{
			AccessKey() { return "ACCESS" }
			SecretKey() { return "secret" }
			}

		hdr = (.v4SignCl)(sqsCl, 'PUT', 'us-east-1', '', '/fredsBucket/testFile.txt',
			#(X_Amz_Security_Token: "token", X_Amz_Content_Sha256: "UNSIGNED-PAYLOAD")).
				AuthorizationHeader()
		Assert(hdr is: .s3Header)
		}

	Test_AmazonV4Signing()
		{
		s3 = AmazonS3
				{
				SecurityToken() { return 'secretToken' }
				AccessKey() { return "ACCESS" }
				SecretKey() { return "secret" }
				}
		c = AmazonV4Signing
			{
			AmazonV4Signing_amazonDate() { return #20000101 }
			}
		Assert((new c(s3, 'GET', 'us-east-1', '', '/test.png')).PresignUrl(
			'testbucket.s3.amazonaws.com')
			is: 'https://testbucket.s3.amazonaws.com/test.png?X-Amz-Algorithm=' $
				'AWS4-HMAC-SHA256&X-Amz-Credential=ACCESS%2F20000101%2F' $
				'us-east-1%2Fs3%2Faws4_request&X-Amz-Date=20000101T000000Z&' $
				'X-Amz-Expires=3600&X-Amz-Security-Token=secretToken&' $
				'X-Amz-SignedHeaders=host&' $
				'response-content-disposition=inline%3B%20filename%3D%22test.png%22&' $
				'response-content-type=image%2Fpng&' $
				'X-Amz-Signature=' $
				'a8629fb8dc4885bdd70b3c46297def3a3ab11b1534a246f4b1bd694db293e264')
		Assert((new c(s3, 'PUT', 'us-east-1', '', '/test.png')).PresignUrl(
			'testbucket.s3.amazonaws.com')
			is: 'https://testbucket.s3.amazonaws.com/test.png?X-Amz-Algorithm=' $
				'AWS4-HMAC-SHA256&X-Amz-Credential=ACCESS%2F20000101%2F' $
				'us-east-1%2Fs3%2Faws4_request&X-Amz-Date=20000101T000000Z&' $
				'X-Amz-Expires=3600&X-Amz-Security-Token=secretToken&' $
				'X-Amz-SignedHeaders=host&' $
				'response-content-disposition=inline%3B%20filename%3D%22test.png%22&' $
				'response-content-type=image%2Fpng&' $
				'X-Amz-Signature=' $
				'3c40d35bfd435d41de2c88cedbbf969665f8f2dc3298367e2c7c77f5de8421c0')
		}
	}
