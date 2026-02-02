// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// An HTTP client implemented in Suneido
// INTERNAL - should only be used directly by Http
class
	{
	CallClass(method, url, content = "", header = #(),
		timeout = 60, timeoutConnect = 60)
		{
		Assert(url.Prefix?('http://'), "Http: url must have http:// prefix")
		Assert(#(GET, PUT, POST, DELETE, OPTIONS, TRACE, HEAD).Has?(method),
			'HttpClient: invalid method: ' $ method)
		return .socketClient(url, timeout, timeoutConnect, method, content, header)
		}
	socketClient(url, timeout, timeoutConnect, method, content, header)
		{
		a = Url.Split(url)
		port = a.GetDefault('port', Http.DefaultPort)
		Assert(Type(port) is: 'Number')
		SocketClient(a.host, port, timeout, :timeoutConnect)
			{|sc|
			.request(sc, method, a, content, header)
			return .response(sc, method)
			}
		}
	request(sc, method, a, content, header)
		{
		header = header.Copy()
		header['User-Agent'] = 'Suneido'
		header.Host = a.host $ Opt(':', a.GetDefault(#port, ''))
		header.Connection = 'close'
		HttpSend(sc,
			method $ " " $ a.GetDefault(#path, '/') $ " HTTP/1.1", header, content)
		}
	response(sc, method)
		{
		headerLines = InetMesg.ReadHeader(sc)
		code = .responseCode(headerLines.GetDefault(0, ""))
		props = .getResponseProperties(headerLines)
		content = ""
		if .hasContent(method, code)
			{
			if props.chunked
				content = .readChunked(sc)
			else if props.content_length isnt false
				{
				content = sc.Read(props.content_length)
				.checkContentLength(content, props.content_length)
				}
			else
				content = sc.Read()
			}
		return [header: headerLines.Join("\r\n"),
			content: content is false ? "" : content]
		}
	ContentLengthErrPrefix: 'HttpClient: content size ('
	checkContentLength(content, content_length)
		{
		size = content is false ? 0 : content.Size()
		if size isnt content_length
			throw .ContentLengthErrPrefix $ content.Size() $ ')' $
				' did not match content-length from header (' $ content_length $ ')'
		}
	responseCode(header)
		{
		return header.AfterFirst(' ').Trim().BeforeFirst(' ')
		}
	getResponseProperties(headerLines)
		{
		chunked = false
		content_length = false
		for line in headerLines
			{
			line = line.Lower()
			if line.Prefix?("content-length:")
				content_length = Number(line.AfterFirst(':'))
			if line.Prefix?("transfer-encoding:")
				chunked = line.Has?("chunked")
			}
		return Object(:content_length, :chunked)
		}
	hasContent(method, code)
		{
		return method isnt 'HEAD' and code !~ '^(1..|204|304)$'
		}
	readChunked(sc)
		{
		// see: http://en.wikipedia.org/wiki/Chunked_transfer_encoding
		// see: http://tools.ietf.org/html/rfc2616
		content = ""
		forever
			{
			line = sc.Readline()
			chunkSize = ("0x" $ line.BeforeFirst(';')).SafeEval()
			if chunkSize is 0
				break
			chunk = sc.Read(chunkSize)
			Assert(chunk.Size() is: chunkSize, msg: 'lineSize: ' $ line)
			content $= chunk
			line = sc.Readline()
			Assert(line is: "")
			}
		// for now, just skip trailing header
		while "" isnt sc.Readline()
			{
			}
		return content
		}
	}