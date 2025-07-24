// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
AmazonAWS
	{
	region: 'us-east-1'
	Host(region = false)
		{
		return 'guardduty.' $ (region is false ? .region : region) $ '.amazonaws.com'
		}

	ContentType()
		{
		return 'application/x-amz-json-1.0'
		}

	Service()
		{
		return 'guardduty'
		}

	CanonicalQueryString(params /*unused*/)
		{
		return ''
		}

	PayloadHash(params)
		{
		return Sha256(params).ToHex()
		}

	makeRequest(call, params, path, content = '', region = false)
		{
		region = region is false ? .region : region
		params = AmazonAWS.UrlEncodeValues(params)
		url = 'https://' $ .Host(region) $ path $ Opt('?', params)
		extraHeaders = Object()
		extraHeaders.X_Amz_Security_Token = .SecurityToken()

		if false is header = .signRequest(call, region, content, path, extraHeaders)
			return false
		return .https(call, url, header, content)
		}

	https(call, url, header = #(), content = '')
		{
		return Https[call.Capitalize()](url, :header, :content)
		}

	signRequest(call, region, content, path, extraHeaders)
		{
		region = region is false ? .region : region
		return AmazonV4Signing(this, call, region, content, path,
			extraHeaders).AuthorizationHeader()
		}

	CreateMalwareProtectionPlan(bucket, role)
		{
		region = AmazonS3.GetBucketLocationCached(bucket)
		content = Object(
			'protectedResource': Object('s3Bucket': ['bucketName': bucket]),
			'role': role)

		return .makeRequest('POST', [], '/malware-protection-plan',
			content: Json.Encode(content), :region)
		}
	}