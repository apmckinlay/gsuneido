// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide
class
	{
	CallClass(port, endpoint, args = #(), ip = false, timeout = 60, timeoutConnect = 60,
		auth? = false)
		{
		if ip is false
			ip = .defaultIP()
		address = "http://" $ ip $ ':' $ port $ '/' $ endpoint
		if auth? is true
			args = args.Copy().Add(Auth.Token(), at: #runOnHttpToken)
		content = Pack(args)
		result = Http.Post(address, :content, :timeout, :timeoutConnect)
		return Unpack(result)
		}

	defaultIP()
		{
		if "" isnt ip = ServerIP()
			return ip
		return "127.0.0.1"
		}
	}
