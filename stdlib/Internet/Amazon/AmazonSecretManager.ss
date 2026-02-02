// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
AmazonAWS
	{
	Service()
		{
		return 'secretsmanager'
		}

	Host(region = false)
		{
		return .Service() $ '.' $ (region is false ? .region : region) $ '.amazonaws.com'
		}

	ContentType()
		{
		return 'application/x-amz-json-1.1'
		}

	PayloadHash(params)
		{
		return Sha256(params).ToHex()
		}

	CanonicalQueryString(unused)
		{
		return ''
		}

	region: 'us-east-1'
	makeRequest(target, params, extraHeaderInfo = #(), region = false)
		{
		region = region is false ? .region : region
		params = Json.Encode(params)
		extraHeaderInfo = extraHeaderInfo.Copy()
		extraHeaderInfo.X_Amz_Security_Token = .SecurityToken()
		extraHeaderInfo.X_Amz_Target = 'secretsmanager.' $ target
		return false isnt (header = .header(region, params, extraHeaderInfo))
			? .https(region, params, header)
			: false
		}

	header(region, params, extraHeaderInfo)
		{
		return AmazonV4Signing(this, 'POST', region, params, :extraHeaderInfo).
			AuthorizationHeader()
		}

	https(region, content, header)
		{
		result = Https('POST', .url(region), :content, :header, timeoutConnect: 60)
		Http.ResponseCode(result.header)
		result = Json.Decode(result.content)
		if result.GetDefault('__type', '').Has?('Exception')
			throw result.MapMembers({ it.Lower() }).message
		return result
		}

	url(region)
		{
		return 'https://' $ .Host(region)
		}

	ListSecrets(filters = #(/*(key: "string", values: "string")*/),
		includePlannedDeletion = false, maxResults = 100, nextToken = '',
			sortOrder = 'desc' /*or asc*/)
		{
		if maxResults > 100 /*= max allowed value*/
			throw 'maxResults cannot exceed 100'
		params = Object(
			Filters: filters,
			IncludePlannedDeletion: includePlannedDeletion,
			MaxResults: maxResults,
			SortOrder: sortOrder)
		if nextToken isnt ''
			params.NextToken = nextToken
		// Expected result:
		//		#(
		// 			SecretList:	<List of secrets>,
		//			NextToken:	<Only returned if there are more secrets then maxResult>)
		return .makeRequest('ListSecrets', params)
		}

	GetSecretValue(secretId)
		{
		result = .makeRequest('GetSecretValue', Object(SecretId: secretId))
		return Json.Decode(result.SecretString).apiKey
		}

	GetRandomPassword(excludeCharacters = '', excludeLowercase = false,
		excludeNumbers = false, excludePunctuation = false,
		excludeUppercase = false, includeSpace = false, passwordLength = 5,
		requireEachIncludedType = false)
		{
		result = .makeRequest('GetRandomPassword',
			Object(
				ExcludeCharacters: excludeCharacters,
				ExcludeLowercase: excludeLowercase,
				ExcludeNumbers: excludeNumbers,
				ExcludePunctuation: excludePunctuation,
				ExcludeUppercase: excludeUppercase,
				IncludeSpace: includeSpace,
				PasswordLength: passwordLength,
				RequireEachIncludedType: requireEachIncludedType))
		return result.RandomPassword
		}
	}