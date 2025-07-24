// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// abstract base class for AmazonS3 and AmazonSQS
// Uses AmazonCredentialsHandler which must be contributed by the application
// and must be a function returning a record containing:
//	amazoncr_access_key: <access key id>
//	amazoncr_secret_key: <secret access key>
//	amazoncr_security_token: <session token>
class
	{
	//Error Msgs
	CredentialErrMsg: 'unable to get temporary credentials'

	ConnectionFailedMsg: 'There was a problem trying to complete your request.\n' $
		'Please check your internet connection ' $
		'(Business > Connection And Network Testers) and try again.'

	AccessKey()
		{
		return .getCredentialField('amazoncr_access_key')
		}

	SecretKey()
		{
		return .getCredentialField('amazoncr_secret_key')
		}

	SecurityToken()
		{
		return .getCredentialField('amazoncr_security_token')
		}

	getCredentialField(field)
		{
		if false is credentials = SoleContribution('AmazonCredentialsHandler')()
			return .CredentialErrMsg
		return credentials[field]
		}

	UnixTime() // overridden by tests
		{
		// use time from internet in case local clock is wrong
		return .retryWithFallback({ Htp.UnixTime() }, { UnixTime() })
		}

	GMTime()
		{
		// use time from internet in case local clock is wrong
		return .retryWithFallback(
			{ Htp.InternetFormatWithThrow() 	},
			{ Date().InternetFormat() 			})
		}

	retryWithFallback(defaultBlock, fallBackBlock)
		{
		try
			return Retry(defaultBlock, maxRetries: 3)
		catch (err)
			{
			SuneidoLog('INFO: AmazonAWS - ' $ err)
			return fallBackBlock()
			}
		}

	UrlEncodeValues(params)
		{
		return Url.EncodeValues(params.Copy().RemoveIf(#Blank?).Map!(#ToUtf8).Sort!())
		}
	}