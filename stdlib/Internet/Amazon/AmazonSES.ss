// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
/*
e.g.
	msg = MimeText('hello').
		Subject('ses test').
		From('email@mydomain.com'). // must be verified with SES
		To('sue@gmail.com').
		Reply_To('joe@abc.com').
		ToString()					// SES supplies Date and Message-Id
	AmazonSES.SendRawEmail(msg)

amazon_access_keys must be configured correctly
*/
class
	{
	Host()
		{
		return .getEndPoint()
		}

	ContentType()
		{
		return  'application/x-www-form-urlencoded'
		}

	Service()
		{
		return 'ses'
		}

	AccessKey()
		{
		return AmazonKeys.Access()
		}

	SecretKey()
		{
		return AmazonKeys.Secret()
		}

	CanonicalQueryString(unused)
		{
		return ''
		}

	PayloadHash(params)
		{
		return Sha256(params).ToHex()
		}

	SendRawEmail(mime_msg)
		{
		params = Object(Action: 'SendRawEmail',
			'RawMessage.Data': Base64.Encode(mime_msg))
		result = .doRequest(params)
		return result.code is '200' ? true : result
		}

	GetSendQuotaInfo()
		{
		params = Object(Action: 'GetSendQuota')
		result = .doRequest(params)
		if result.code isnt '200'
			return false
		quotaInfoOb = Object().Set_default(0)
		xmlResponse = XmlParser(result.content)
		for node in xmlResponse.Children()
			{
			if node.Name() is "getsendquotaresult"
				for attrib in node.Children()
					if attrib.Text().Number?()
						quotaInfoOb[attrib.Name()] = Number(attrib.Text())
			}
		return quotaInfoOb
		}

	retryErrorMsg: 'Bad HTTP Status Code'
	doRequest(params)
		{
		message = AmazonAWS.UrlEncodeValues(params)
		result = Object().Set_default('')
		try
			{
			Retry(maxRetries: 3, minDelayMs: 100, retryException: .retryErrorMsg)
				{
				.retryRequest(message, result)
				}
			}
		catch (err, "Retry failed")
			{
			if result.code is ''
				result.code = '500'
			if result.header is ''
				result.header = err
			result.content = Opt(result.content, '\r\n\r\n') $ err
			}
		return result
		}
	retryRequest(message, result)
		{
		try
			{
			response = .postRequest(message)
			result.Merge(response)
			.parseResponse(result)
			}
		catch (err, 'curl: (35) gnutls_handshake() failed:')
			{
			// treat the Curl tls error like a 500 error so it can Retry
			SuneidoLog('AmazonSES: ' $ err)
			throw .retryErrorMsg
			}
		}

	postRequest(message)
		{
		endPoint = .getEndPoint()
		region = endPoint.Split('.')[1]
		header = AmazonV4Signing(this, 'POST', region, message).AuthorizationHeader()
		return Https('POST', 'https://' $ endPoint, content: message, :header)
		}

	us_east_endpoint: 'email.us-east-1.amazonaws.com'
	us_west_endpoint: 'email.us-west-2.amazonaws.com'
	endPointTable: 'amazon_ses_endpoint'
	getEndPoint()
		{
		if TableExists?(.endPointTable) and false isnt ep = Query1(.endPointTable)
			return ep.aws_ses_endpoint
		return .us_east_endpoint
		}

	OutputEndPointRecord(endPoint)
		{
		Database('ensure ' $ .endPointTable $ '(aws_ses_endpoint) key ()')
		QueryDo('delete ' $ .endPointTable)
		QueryOutput(.endPointTable, Record(aws_ses_endpoint: endPoint))
		}

	parseResponse(result)
		{
		result.code = Http.ResponseCode(result.header)
		if result.code.Prefix?('5')
			throw .retryErrorMsg
		}

	InvalidAttachmentExtensions: (ade, adp, app, asp, bas, bat, cer, chm, cmd,
		com, cpl, crt, csh, der, exe, fxp, gadget, hlp, hta, inf, ins, isp, its,
		js, jse, ksh, lib, lnk, mad, maf, mag, mam, maq, mar, mas, mat, mau, mav,
		maw, mda, mdb, mde, mdt, mdw, mdz, msc, msh, msh1, msh2, mshxml, msh1xml,
		msh2xml, msi, msp, mst, ops, pcd, pif, plg, prf, prg, reg, scf, scr, sct,
		shb, shs, sys, ps1, ps1xml, ps2, ps2xml, psc1, psc2, tmp, url, vb, vbe,
		vbs, vps, vsmacros, vss, vst, vsw, vxd, ws, wsc, wsf, wsh, xnk)

	// other actions -----------------------------------------------------------

	VerifyEmailAddress(address)
		{
		params = Object(Action: 'VerifyEmailAddress', EmailAddress: address)
		result = .doRequest(params)
		return result.code is '200'
		}

	DeleteVerifiedEmailAddress(address)
		{
		params = Object(Action: 'DeleteVerifiedEmailAddress', EmailAddress: address)
		result = .doRequest(params)
		return result.code is '200'
		}

	ListVerifiedEmailAddresses()
		{
		params = Object(Action: 'ListVerifiedEmailAddresses')
		result = .doRequest(params)
		if result.code isnt '200'
			return false
		emails = Object()
		parsedResponse = XmlParser(result.content)
		for node in parsedResponse.Children()
			{
			if node.Name() is 'listverifiedemailaddressesresult'
				for email in node.Children()
					if email.Name() is "verifiedemailaddresses"
						for child in email.Children()
							emails.AddUnique(child.Text())
			}
		return emails
		}

	SourceEmail(from)
		{
		return Display(from) $ ' <' $ SoleContribution('AmazonSESSourceEmail') $ '>'
		}
	}
