// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// Unlike Http, this always uses Curl
// since we don't have built-in support for SSL
Http
	{
	// NOTE: Using the Get/Post/Put methods is preferred, it handles checking the response
	// code. When using the CallClass directly, the response code must be checked by
	// the calling code.
	CallClass(method, addr, content = '', user = '', pass = '',
		cookies = '', fromFile = '', toFile = '', header = #(), timeoutConnect = false,
		limitRate = '')
		{
		Assert(addr.Prefix?('https://'), 'Https url must have https:// prefix')
		Assert(Object?(header), 'Https header argument must be an object')

		return Curl.Http(method, addr, :user, :pass, :content, :fromFile, :toFile,
			:header, :cookies, :timeoutConnect, :limitRate)
		}
	PostWithOptions(method, addr, content = '', user = '', pass = '',
		cookies = '', fromFile = '', toFile = '', header = #(), options = #())
		{
		return Curl('https', addr, user, pass, :options).Http(method, addr, :content,
			:fromFile, :toFile,	:header, :cookies)
		}

	Piped(method, url, block, user = '', pass = '', cookies = '',
		header = #(), timeoutConnect = false, checkCode? = false)
		{
		Assert(url.Prefix?('https://'), 'Https url must have https:// prefix')
		Assert(Object?(header), 'Https header argument must be an object')

		super.Piped(method, url, block, :user, :pass, :cookies, :header, :timeoutConnect,
			:checkCode?)
		}
	}
