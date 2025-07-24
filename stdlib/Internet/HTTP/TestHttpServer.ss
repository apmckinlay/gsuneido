// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(ip = false, port = false)
		{
		if true isnt result = .testHttpServer(ip, port)
			SuneidoLog("ERRATIC: TestHttpServer - " $ result)
		return result
		}

	testHttpServer(ip, port)
		{
		if ip is false
			ip = ServerIP()
		if port is false
			port = HttpPort()
		addr = 'http://' $ ip $ ':' $ port $ '/TestResponse'
		try
			{
			resp = Http.Get(addr, timeout: 5, timeoutConnect: 5)
			return resp is 'TestResponse' ? true : "bad response: " $ resp
			}
		catch (err)
			return err
		}
	}