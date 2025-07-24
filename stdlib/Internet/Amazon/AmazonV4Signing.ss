// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
// NOTE: when params is sent through query string with GET request (like: aws.com?abc=1),
// NOTE: PayloadHash params should be empty
class
	{
	CallClass(awsClass, method, region, params, path = '/', extraHeaderInfo = #())
		{
		return new this(awsClass, method, region, params, path, extraHeaderInfo)
		}

	New(.awsClass, .method, .region, .params, .path = '/', extraHeaderInfo = #())
		{
		.service = .awsClass.Service()
		.accessKey = .awsClass.AccessKey()
		.secretKey = .awsClass.SecretKey()
		.date = .amazonDate()
		.headerInfo = Object(Content_Type: .awsClass.ContentType(),
			X_Amz_Date: .fmtAmazonDateTime(.date),
			Host: .awsClass.Host(:region)).MergeNew(extraHeaderInfo)
		}

	AuthorizationHeader()
		{
		if .accessKey is false or .secretKey is false
			return false
		canonicalRequest = .createCanonicalRequest()
		stringToSign = .createStringToSign(canonicalRequest)
		signature = .buildSignature(stringToSign)
		return .headerInfo.Add(.buildAuthorizationHeader(signature), at: 'Authorization')
		}

	// https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-query-string-auth.html
	PresignUrl(host, download? = false)
		{
		if .accessKey is false or .secretKey is false
			return false

		params = Object()
		params.X_Amz_Algorithm = "AWS4-HMAC-SHA256"
		params.X_Amz_Credential = .accessKey $ '/' $ .getCredentialScope()
		params.X_Amz_Date = .fmtAmazonDateTime(.date)
		params.X_Amz_Expires = '3600'
		params.X_Amz_SignedHeaders = 'host'
		params.X_Amz_Security_Token = .awsClass.SecurityToken()
		filename = Paths.Basename(.path)
		params.response_content_disposition = (download? ? 'attachment' : 'inline') $
			`; filename="` $ filename $ `"`
		params.response_content_type =
			MimeTypes.GetDefault(filename.AfterLast('.').Lower(), 'multipart/form-data')
		headers = Object()
		for m, v in params
			headers[m.Replace('_', '-')] = v

		canonicalRequest = .method $ '\n' $
			.path $ '\n' $
			(q = AmazonAWS.UrlEncodeValues(headers)) $ '\n' $
			'host:' $ host $ '\n\n' $
			'host\n' $
			.awsClass.PayloadHash(params)

		stringToSign = .createStringToSign(canonicalRequest)

		signature = .buildSignature(stringToSign)
		return 'https://' $ host $ .path $ '?' $ q $ '&X-Amz-Signature=' $ signature
		}

	canonicalHeaders()
		{
		headerStr = ''
		for header in .headerInfo.Members().Sort!()
			headerStr $= .fmtHeader(header) $ ':' $ .headerInfo[header] $ '\n'
		return headerStr
		}

	fmtHeader(header)
		{
		return header.Lower().Replace('_', '-')
		}

	createCanonicalRequest()
		{
		return .method $ '\n' $
			.path $ '\n' $
			.awsClass.CanonicalQueryString(.params) $ '\n' $
			.canonicalHeaders() $ '\n' $
			.signedHeadersOb() $ '\n' $
			.awsClass.PayloadHash(.params)
		}

	signedHeadersOb()
		{
		return .headerInfo.Members().Sort!().Map({ .fmtHeader(it) }).Join(';')
		}

	amazonDate()
		{
		return Date(AmazonAWS.GMTime().AfterFirst(', ').BeforeFirst(' GMT'))
		}

	fmtAmazonDate(date)
		{
		return date.Format('yyyyMMdd')
		}

	fmtAmazonDateTime(date)
		{
		return date.Format('yyyyMMddTHHmmssZ')
		}

	hashAlgorithm: 'AWS4-HMAC-SHA256'

	createStringToSign(canonicalRequest)
		{
		credentialScope = .getCredentialScope()
		return .hashAlgorithm $ '\n' $
			.fmtAmazonDateTime(.date) $ '\n' $
			credentialScope $ '\n' $
			Sha256(canonicalRequest).ToHex()
		}

	sigTerminator: 'aws4_request'

	getCredentialScope()
		{
		return .fmtAmazonDate(.date) $ '/' $
			.region.ToUtf8() $ '/' $
			.service.ToUtf8() $ '/' $
			.sigTerminator
		}

	buildSignature(stringToSign)
		{
		signingKey = .getSignatureKey()
		return Hmacsha256(stringToSign.ToUtf8(), signingKey).ToHex()
		}

	getSignatureKey()
		{
		dateStr = .fmtAmazonDate(.date)
		kDate = .sign(('AWS4' $ .secretKey).ToUtf8(), dateStr)
		kRegion = .sign(kDate, .region)
		kService = .sign(kRegion, .service)
		return .sign(kService, .sigTerminator)
		}

	sign(key, msg)
		{
		return Hmacsha256(msg.ToUtf8(), key)
		}

	buildAuthorizationHeader(signature)
		{
		credentialScope = .getCredentialScope()
		return .hashAlgorithm $ ' ' $ 'Credential=' $ .accessKey $ '/' $
			credentialScope $ ', ' $ 'SignedHeaders=' $ .signedHeadersOb() $ ', ' $
			'Signature=' $ signature
		}
	}