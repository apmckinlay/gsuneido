// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Host()
		{
		return 'route53.amazonaws.com'
		}

	ContentType()
		{
		return 'application/x-www-form-urlencoded; charset=utf-8'
		}

	Service()
		{
		return 'route53'
		}

	CanonicalQueryString(params)
		{
		return params
		}

	PayloadHash(params)
		{
		params = '' // params is sent as query string, payload/body is empty
		return Sha256(params).ToHex()
		}

	AccessKey()
		{
		return AmazonKeys.Access()
		}

	SecretKey()
		{
		return AmazonKeys.Secret()
		}

	ListResourceRecordSets(domainId)
		{
		mapping = Object()
		next = Object()
		while false isnt result = .list(domainId, next)
			if false is next = .extractMapping(result, mapping)
				break
		return mapping
		}

	list(domainId, next)
		{
		path = '/2013-04-01/hostedzone/' $ domainId $ '/rrset'
		params = Object().Merge(next)
		return .makeRequest('GET', params, path)
		}

	region: 'us-east-1'
	makeRequest(call, params, path)
		{
		params = AmazonAWS.UrlEncodeValues(params)
		url = 'https://' $ .Host() $ path $ Opt('?', params)
		extraHeaders = Object()
		if false is header = .signRequest(call, .region, params, path, extraHeaders)
			return false

		return .https(call, url, header)
		}

	signRequest(call, region, params, path, extraHeaders)
		{
		return AmazonV4Signing(
			this, call, region, params, path, extraHeaders).AuthorizationHeader()
		}

	https(call, url, header = #())
		{
		return Https[call.Capitalize()](url, :header)
		}

	extractMapping(responseContent, mapping)
		{
		next = Record(more: false)
		response = XmlParser(responseContent)
		for node in response.Children()
			{
			if node.Name() is 'resourcerecordsets'
				for recSets in node.Children()
					if recSets.Name() is "resourcerecordset"
						.extractNameIp(recSets, mapping)
			.extractNextStart(node, next)
			}
		return next.more is true and next.name $ next.type isnt ''
			? next.Project('name', 'type')
			: false
		}

	extractNextStart(node, next)
		{
		if node.Name() is 'istruncated' and
			node.Text().Trim().Lower() is 'true'
			next.more = true
		if node.Name() is 'nextrecordname'
			next.name = node.Text().Trim()
		if node.Name() is 'nextrecordtype'
			next.type = node.Text().Trim()
		}

	extractNameIp(recSets, mapping)
		{
		name = ip = false
		for recordSet in recSets.Children()
			{
			recordSetMem = recordSet.Name().Lower()
			if recordSetMem is 'name'
				name = recordSet.Text().Trim().Lower().RemoveSuffix('.')
			if recordSetMem is 'resourcerecords'
				for record in recordSet.Children()
					{
					ip = record.Children()[0].Text().Trim()
					break
					}
			}
		if ip isnt false and name isnt false
			mapping[ip] = name
		}
	}
