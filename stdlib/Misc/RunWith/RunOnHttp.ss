// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide
class
	{
	CallClass(port, endpoint, args = #(), ip = false,
		timeout = 60, timeoutConnect = 60, asyncCompletion = false)
		{
		if ip is false
			ip = .defaultIP()
		address = 'http://' $ ip $ ':' $ port $ '/' $ endpoint
		content = Pack(args)
		if asyncCompletion isnt false
			{
			// Calls Http instead of Http.Post, in order to circumvent Http.ResponseCode.
			// As when using asyncCompletion, a thread is started and nothing is returned
			Http('POST', address, :content, :timeout, :timeoutConnect, :asyncCompletion)
			return true
			}
		result = Http.Post(address, :content)
		return Unpack(result)
		}

	defaultIP()
		{
		if '' isnt ip = ServerIP()
			return ip
		return '127.0.0.1'
		}
	}
