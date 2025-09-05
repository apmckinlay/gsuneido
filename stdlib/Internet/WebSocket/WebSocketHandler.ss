// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
class
	{
	typeMap: #(
		0x1: #TEXT,
		0x2: #BINARY,
		0x8: #CLOSE,
		0x9: #PING,
		0xA: #PONG,
		TEXT: 0x1,
		BINARY: 0x2,
		CLOSE: 0x8,
		PING: 0x9,
		PONG: 0xA)
	closed?: false
	Id: false
	DBClosed?: false
	CallClass(env, socket, appHandler)
		{
		instance = new this(env, socket, appHandler)
		return instance.Run()
		}

	New(.env, .socket, .appHandler) {}

	Run()
		{
		if .handShake() is false
			return Object('BadRequest', Object(), '')

		.Id = Timestamp()

		open = false
		try
			open = (.appHandler.App)(
				.env.Copy().Append(Object(type: #OPEN, body: '', wsHandler: this)))
		catch (e)
			.errorHandler(e)

		if open isnt true
			{
			if .closed? isnt true and open isnt 'close'
				try .CloseSocket(reason: String?(open) ? open : 'Initilization failed')
			return .disconnect()
			}

		.Loop()
		return .disconnect()
		}

GetLevel()
	{
	return .level
	}

	QUITLOOP: 'WebSocket Handler: QUIT_LOOP'
	level: 0
	Loop()
		{
		.level++
		forever
			{
			req = .env.Copy()
			try
				{
				req = req.Merge(.receive().Add(this, at: #wsHandler))
				switch (req.type)
					{
				case #TEXT, #BINARY:
					(.appHandler.App)(req)
				case #CLOSE:
					.handleClose(req) // .handleClose calls .Terminate
				case #PING:
					.Pong(req.body)
				case #PONG: // ignore
					}
				.errorCount = 0
				}
			catch (e)
				{
				if .errorHandler(e, req)
					return
				}
			}
		}

	errorCount: 0
	// return true to exit
	errorHandler(e, req = "")
		{
		if .connectionError?(e)
			{
			if false is .appHandler.OnConnectionError(e, this)
				return false
			.Terminate()
			}
		if e is .QUITLOOP
			{
			Assert(.level greaterThan: 1) // to detect Loop problems
			.level--
			return true
			}
		stop? = ++.errorCount >= 3 /*=max consecutive errors*/
		.log("WebSocket Handler" $ (stop? ? ' (closing)' : '') $ ": " $ e, params: req)
		return stop?
		}

	Terminate(reason = false, e = false)
		{
		if e isnt false
			.log(e)
		if reason isnt false
			.CloseSocket(:reason)
		.disconnect()
		Thread.Exit()
		throw "After Thread.Exit(). Should not be here"
		}

	log(e, params = "")
		{
		try
			SuneidoLog('ERROR (CAUGHT): ' $ e, :params, caughtMsg: 'needs attention')
		catch (err)
			{
			s = e $ ' (err: ' $ err $ ')'
			if Type(e) is 'Except'
				s $= '\r\n' $ e.Callstack().Map({ Display(it.fn) }).Join('\r\n')
			ErrorLog(s)
			}
		}

	handShake()
		{
		if not .env.Member?('sec_websocket_key')
			return false

		key = .env['sec_websocket_key']
		accept = Base64.Encode(Sha1(key $ "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
		HttpSend(.socket, "HTTP/1.1 101 Switching Protocols",
			[Upgrade: 'websocket', Connection: 'Upgrade', Sec_WebSocket_Accept: accept],
			'')
		return true
		}

	receive()
		{
		byte1 = .readByte()
		fin? = (byte1 & 0x80) isnt 0
		opcode = byte1 & 0x0f

		byte2 = .readByte()
		mask? = (byte2 & 0x80) isnt 0
		Assert(mask?)

		payloadLen = byte2 & 0x7f
		payloadLen = .calcPayloadLen(payloadLen)

		mask = .read(4/*=32 bits mask*/)
		encoded = .read(payloadLen)
		decoded = StringXor(encoded, mask)

		if opcode is 0x0 // continuation frame
			return fin? ? decoded : decoded $ .receive()

		body = fin? ? decoded : decoded $ .receive()
		return Object(type: .typeMap[opcode], :body)
		}

	readByte()
		{
		return .read(1).Asc()
		}

	read(n)
		{
		if false is result = .socket.Read(n)
			throw 'socket.Read: unexpected EOF'
		return result
		}

	calcPayloadLen(payloadLen)
		{
		if payloadLen < 126 /*=magic number for 16 bits*/
			return payloadLen

		if payloadLen is 126 /*=magic number for 16 bits*/
			return (.readByte() << 8) + .readByte()

		if payloadLen is 127 /*=magic number for 64 bits*/
			{
			len = 0
			for .. 8
				len = (len << 8) + .readByte()
			return len
			}
		}

	closeReq: false
	handleClose(req)
		{
		e = false
		if not .closed?
			{
			if req.body.Has?('CloudFlare')
				throw 'lost connection' // to trigger reconnect

			try
				{
				(.appHandler.App)(req)
				.close(req.body)
				}
			catch (err)
				{
				// ignore connection error since we are already in closing
				if not .connectionError?(err)
					e = 'WebSocketHandler.handleClose - ' $ err
				}
			}
		.Terminate(:e)
		}

	CloseSocket(code = 1000, reason = '')
		{
		body = (code >> 8).Chr() $ (code & 0xff).Chr() $ reason
		.close(body)
		}

	close(body)
		{
		.closed? = true
		.Send(#CLOSE, body)
		}

	Pong(body)
		{
		.Send(#PONG, body)
		}

	Send(@args) // only suport single frame
		{
		type = args[0]
		if .closed? and type isnt #CLOSE
			return

		contents = args[1..]
		head = (0x80 /*=FIN*/ + .typeMap[type]).Chr() $
			.calcSendPayloadLen(contents.SumWith(#Size))
		.socket.Write(head)
		for content in contents
			.socket.Write(content)
		}

	calcSendPayloadLen(len)
		{
		if len < 126 /*=magic number for 7 bits*/
			return len.Chr()

		if len < 65536	/*=max 16 bits unsigned*/
			return 126.Chr() $ (len >> 8).Chr() $ (len & 0xff).Chr() /*=magic*/

		ob = Object()
		for i in .. 8 /*=8 bytes*/
			{
			ob[7 - i] = len & 0xff /*=8 bytes*/
			len >>= 8
			}
		return 127.Chr() $ ob.Map(#Chr).Join() /*=magic number for 64 bits*/
		}

	connectionError?(e)
		{
		return e.Has?('lost connection') or e.Prefix?('socket')
		}

	InternalError?(e)
		{
		return .connectionError?(e) or e is .QUITLOOP
		}

	SetSocket(.socket)
		{
		}

	GetSocket()
		{
		return .socket
		}

	disconnect()
		{
		(.appHandler.BeforeDisconnect)()
		return -1
		}
	}
