// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
// NOTE: this is deliberately named so it will NOT be found by the test runner
// because it is slow and accesses the internet
// See: http://httpbin.org/
class
	{
	tmp: 'test.tmp'
	robot: 'User-agent: *\nDisallow: /deny\n'
	CallClass()
		{
		result = .checkGet()
		result $= .checkPut()
		result $= .checkConnections()
		DeleteFile(.tmp)
		failed = result.Count('FAILED') + result.Count('ERROR')
		passed = result.Count('PASSED')
		pre = 'TESTING WITH HTTP BIN. '
		pre $= failed is 0
			? 'All Tests PASSED'
			: failed $ ' Tests FAILED out of ' $ (passed + failed)
		return pre $ '\r\n\r\n' $ result
		}

	checkGet()
		{
		url = 'httpbin.org/robots.txt'
		failedPrefix = '\tFAILED.\r\n\r\nEXPECTED:\r\n' $ .robot $ '\r\nGOT:\r\n'
		result = .checkHttpClientGet(url, failedPrefix)
		result $= .checkCurlHttpHttpsGet(url, failedPrefix)
		result $= .checkHttpGet(url, failedPrefix)
		result $= .checkHttpGetCurl(url, failedPrefix)
		result $= .checkHttpsGet(url, failedPrefix)
		result $= .checkGetToFile(url, failedPrefix)
		return result $= .checkGetChunked()
		}

	checkHttpClientGet(url, failedPrefix)
		{
		s = 'Testing HttpClient GET\r\n'
		return s $= TimedTryCatch()
			{
			(getResult = HttpClient('GET', 'http://' $ url).content) is .robot
				? '\tPASSED'
				: failedPrefix $ getResult $ '\r\n'
			}
		}

	checkCurlHttpHttpsGet(url, failedPrefix)
		{
		s = ''
		for p in #('http', 'https')
			{
			s $= 'Testing Curl.Http GET (' $ p $ ')\r\n' $ TimedTryCatch()
				{
				(getResult = Curl.Http('GET', p $ '://' $ url, '').content) is .robot
					? '\tPASSED'
					: failedPrefix $ getResult $ '\r\n'
				}
			}
		return s
		}

	checkHttpGet(url, failedPrefix)
		{
		s = 'Testing Http.Get\r\n'
		return s $= TimedTryCatch()
			{
			(getResult = Http.Get('http://' $ url)) is .robot
				? '\tPASSED'
				: failedPrefix $ getResult $ '\r\n'
			}
		}

	checkHttpGetCurl(url, failedPrefix)
		{
		s = 'Testing Http.Get with curl option\r\n'
		return s $= TimedTryCatch()
			{
			(getResult = Http.Get('http://' $ url, curl:)) is .robot
				? '\tPASSED'
				: failedPrefix $ getResult $ '\r\n'
			}
		}

	checkHttpsGet(url, failedPrefix)
		{
		s = 'Testing Https.Get\r\n'
		return s $= TimedTryCatch()
			{
			(getResult = Https.Get('https://' $ url)) is .robot
				? '\tPASSED'
				: failedPrefix $ getResult $ '\r\n'
			}
		}

	checkGetToFile(url, failedPrefix)
		{
		s = 'Testing Http.Get with toFile option (curl)\r\n'
		return s $= TimedTryCatch()
			{
			DeleteFile(.tmp)
			Http.Get('http://' $ url, toFile: .tmp, curl:)
			(contents = GetFile(.tmp)) is .robot
				? '\tPASSED'
				: failedPrefix $ contents $ '\r\n'
			}
		}

	checkGetChunked()
		{
		s = ''
		for curl in #(false, true)
			{
			s $= 'Testing Http.Get Chunked (curl: ' $ curl $ ')\r\n'
			// response is not specifying chunked (not tested?)
			s $= TimedTryCatch()
				{
				x = Http.Get('http://httpbin.org/stream-bytes/1024', :curl)
				(size = x.Size()) is 1024 /*=expected chunk size*/
					? '\tPASSED'
					: '\tFAILED. Unexpected chunk size per packet. Got chunk size: ' $
						size $ '. Expected chunk size: 1024'
				}
			}
		return s
		}

	// PUT
	checkPut()
		{
		result = .checkPutPost()
		return result $= .checkPutPostFromFile()
		}

	checkPutPost()
		{
		s = ''
		data = 'the\r\ndata'
		for method in #(Put, Post)
			for curl in #(false, true)
				{
				s $= 'Testing Http.' $ method $ ' (curl: ' $ curl $ ')\r\n'
				s $= TimedTryCatch()
					{
					x = Http[method]('http://httpbin.org/' $ method.Lower(),
						data, :curl)
					x = ('#' $ x).SafeEval()
					x.data is data
						? '\tPASSED'
						: '\tFAILED.\r\nExpected:\r\n' $ data $ '\r\n\r\nGot: ' $ x.data
					}
				}
		return s
		}
	checkPutPostFromFile()
		{
		s = ''
		data = 'hello\nworld'
		for method in #(Put, Post)
			{
			s $= 'Testing Http.' $ method $ ' from file\r\n'
			s $= TimedTryCatch()
				{
				PutFile(.tmp, data)
				x = Http[method]('http://httpbin.org/' $ method.Lower(),
					fromFile: .tmp, curl:)
				x = ('#' $ x).SafeEval()
				x.data is data
					? '\tPASSED'
					: '\tFAILED.\r\nExpected:\r\n' $ data $ '\r\n\r\nGot: ' $ x.data
				}
			}
		return s
		}

	checkConnections()
		{
		result = .checkCantConnect()
		return result $= .checkResponseCode()
		}

	checkCantConnect()
		{
		s = 'Testing bad http connection\r\n'
		s $= TimedTryCatch()
			{
			thrown? = false
			msg = 'Expected error to be thrown from failed connection but it was not'
			try
				Http.Get('http://dkjfkdjdjfd.com')
			catch (err)
				{
				expectedErrOb = LocalCmds.HttpTestThrowMatch.Split('|')
				if expectedErrOb.Any?({ err.Has?(it) })
					thrown? = true
				else
					msg = 'Expected an error contining: "' $
						expectedErrOb.Join('" or "') $ '". Got: ' $ Display(err.Trim())
				}
			thrown?	? '\tPASSED' : '\tFAILED. ' $ msg
			}

		s $= 'Testing bad curl connection\r\n'
		s $= TimedTryCatch()
			{
			thrown? = false
			msg = 'Expected error to be thrown from failed connection but it was not'
			try
				Http.Get('http://dkjfkdjdjfd.com', curl:)
			catch (err)
				{
				if err.Lower().Has?('could not resolve')
					thrown? = true
				else
					msg = 'Expected an error containing "could not resolve". ' $
						'Got: ' $ Display(err.Trim())
				}
			thrown?	? '\tPASSED' : '\tFAILED. ' $ msg
			}
		}

	checkResponseCode()
		{
		s = ''
		for curl in #(false, true)
			{
			s $= 'Testing Http.Get response codes (curl option is: ' $ curl $ ')\r\n'
			s $= TimedTryCatch()
				{
				thrown? = false
				msg = 'Expected a 404 response to be thrown but it was not\r\n'
				try
					Http.Get('http://httpbin.org/status/404', :curl)
				catch (err)
					{
					if err.Has?('404')
						thrown? = true
					else
						msg = 'Expected a 404 response. Got: ' $ Display(err.Trim())
					}
				thrown?	? '\tPASSED' : '\tFAILED. ' $ msg
				}

			s $= 'Testing Http GET response codes (curl option is: ' $ curl $ ')\r\n'
			s $= TimedTryCatch()
				{
				msg = 'Expected a 404 response\r\n'
				x = Http('GET', 'http://httpbin.org/status/404', :curl)
				response = Http.ResponseCode(x.header)
				response is '404'
					? '\tPASSED'
					: '\tFAILED. Expected a 404 response, Got a ' $ response $ ' response'
				}
			}
		return s
		}
	}