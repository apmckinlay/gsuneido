// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// only handles create_credentials (GetFederationToken)
class
	{
	RequestCredentials(user, accessKey/*unused*/ = '', secretKey/*unused*/ = '',
		tags = #())
		{
		response = .RequestCredentialsRaw(user, tags)
		if response is ''
			return false
		return .extractCredentials(response)
		}

	RequestCredentialsRaw(user, tags = #())
		{
		if user.Blank?()
			throw 'AmazonIAM: user must not be empty'

		requestOb = .BuildRequest(user, :tags)
		if not Object?(requestOb)
			return ''
		try
			return Https.Post(requestOb.url, requestOb.params, header: requestOb.header)
		catch (e)
			{
			.log("Error requesting credentials: " $ e)
			return ''
			}
		}

	BuildRequest(user, actions = false, tags = #())
		{
		messageRec = .buildMessageRec(user, actions, tags)
		params = AmazonAWS.UrlEncodeValues(messageRec)

		if AmazonAWS.CredentialErrMsg is header = AmazonV4Signing(this, 'POST',
			'us-east-1', params).AuthorizationHeader()
			return AmazonAWS.CredentialErrMsg

		url = 'https://' $ .Host() $ '/?' $ params
		return Object(:header, :url, :params)
		}

	Host(region = false)
		{
		region = region is false ? 'us-east-1' : region
		return 'sts.' $ region $ '.amazonaws.com'
		}

	Service()
		{
		return 'sts'
		}

	ContentType()
		{
		return 'text/xml; charset=UTF-8'
		}

	CanonicalQueryString(messageRec)
		{
		return messageRec
		}

	PayloadHash(params)
		{
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

	extractCredentials(response)
		{
		credentials = Object()
		parsedResponse = XmlParser(response)
		for node in parsedResponse.Children()
			{
			if node.Name() is 'getfederationtokenresult'
				for token in node.Children()
					if token.Name() is "credentials"
						for child in token.Children()
							credentials[child.Name()] = child.Text().Trim()
			}
		return credentials
		}

	buildMessageRec(user, actions = false, tags = #())
		{
		if actions is false
			actions = .defaultActions(user)
		rec = [Version: '2011-06-15'
			Action: 'GetFederationToken'
			Name: String(user)
			DurationSeconds: String(129600) /*= 36 hours*/
			Policy: Json.Encode(
				[Version: '2012-10-17'
					Statement: [Object(Effect: 'Allow', Action: actions, Resource: '*')]
				]),
			]
		i = 1
		for key in tags.Members()
			{
			rec['Tags.member.' $ i $ '.Key'] = key
			rec['Tags.member.' $ i $ '.Value'] = tags[key]
			i++
			}
		return rec
		}

	defaultActions(user)
		{
		actions = ['s3:GetBucketLocation', 'sqs:*', 'ses:*']
		func = OptContribution('AdditionalIAMPolicy',
			function (@unused)  { })
		func(user, actions)
		return actions
		}

	log(@args)
		{
		Rlog('amazonIAM', args.Join('\t'))
		}
	}
