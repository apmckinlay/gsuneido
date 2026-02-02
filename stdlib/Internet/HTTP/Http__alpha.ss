// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// HTTP client that uses either HttpClient or Curl
class
	{
	DefaultPort: 80

	// Returns the response body. Throws if response code isnt 2xx
	Get(@args)
		{
		args.Add('GET' at: 0)
		return .checkCode(@args)
		}

	// Returns the response body. Throws if response code isnt 2xx
	Post(@args)
		{
		args.Add('POST' at: 0)
		return .checkCode(@args)
		}

	// Returns the response body. Throws if response code isnt 2xx
	Put(@args)
		{
		args.Add('PUT' at: 0)
		return .checkCode(@args)
		}

	checkCode(@args)
		{
		result = .CallClass(@args)
		if .ResponseCode(result.header) !~ `^2\d\d$`
			throw "Http." $ args[0] $ " failed: " $ result.header.BeforeFirst('\n')
		return result.GetDefault(#content, '')
		}

	// overridden by Https
	// NOTE: Using the Get/Post/Put methods is preferred, it handles checking the response
	// code. When using the CallClass directly, the response code must be checked by
	// the calling code.
	CallClass(method, url, content = '', fromFile = '', toFile = '', header = #(),
		timeout = 60, timeoutConnect = 60, curl = false)
		{
		method = method.Upper()
		if not curl and (toFile isnt "" or fromFile isnt "")
			throw "HttpClient does not support toFile or fromFile"
		impl = curl is false ? HttpClient : Curl.Http
		return impl(method, url, :content, :header,	:fromFile, :toFile,
			:timeout, :timeoutConnect)
		}

	Piped(method, url, block, header = #(), timeoutConnect = 60,
		user = '', pass = '', cookies = '', checkCode? = false)
		{
		method = method.Upper()
		Curl.HttpPiped(method, url, { |pipe|
			pipeWrapper = new .pipeWrapper(pipe, method, checkCode?)
			block(pipeWrapper) },
			:header, :timeoutConnect, :user, :pass, :cookies)
		}

	pipeWrapper: class
		{
		New(.pipe, .method, .checkCode?) { }
		Write(s)
			{
			.pipe.Write(s)
			}

		CloseWrite()
			{
			.pipe.CloseWrite()
			return .checkCode()
			}

		checked?: false
		checkCode()
			{
			.checked? = true
			// probably a curl error, return false to tell the block to not continue
			header = InetMesg.ReadHeader(.pipe).Join('\n')
			if .checkCode? and Http.ResponseCode(header) !~ `^2\d\d$`
				throw "Http." $ .method $ " failed: " $ header.BeforeFirst('\n')

			return header
			}

		readSize: 1024
		Read()
			{
			if not .checked?
				throw '.CloseWrite() was not called'
			return .pipe.Read(.readSize)
			}

		Readline()
			{
			if not .checked?
				throw '.CloseWrite() was not called'
			return .pipe.Readline()
			}
		}

	// e.g. Http.ResponseCode('HTTP/1.1 200 OK') => '200'
	// NOTE: result is a string, NOT a number, BEWARE of e.g. is 200
	ResponseCode(header)
		{
		code = header.Extract("^HTTP/[\d.]+ (\d\d\d)")
		if code is false
			throw 'Invalid HTTP response code in: ' $
				(header.Blank?() ? '(empty header)' : header)
		return code
		}
	}
