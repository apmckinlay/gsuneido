// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
XmlContentHandler
	{
	CallClass(env)
		{
		if env.content_type isnt 'text/xml'
			throw "XmlRpc: content type must be text/xml, not " $ env.content_type
		headers = Object()
		headers['Content_Type'] = 'text/xml'
		call = (new this).Decode2(env.body)
		ok = 200
		try
			return [ok, headers, .EncodeResponse((XmlRpcMap()[call[0]])(@+1call))]
		catch (err)
			return [ok, headers, .EncodeFault(err)]
		}

	preamble: '<?xml version="1.0"?>\r\n'
	EncodeResponse(result)
		{
		return .preamble $
			Xml('methodResponse',
				Xml('params',
					Xml('param', .EncodeValue(result)))) $ '\r\n'
		}
	EncodeFault(err)
		{
		return .preamble $
			Xml('methodResponse',
				Xml('fault',
					Xml('value',
						Xml('string', err)))) $ '\r\n'
		}
	Call(@args)
		// XML-RPC call
		// args[0] is addr, args[1] is method
		{
		addr = args[0]
		if not addr.Prefix?('http://')
			addr = 'http://' $ addr
		response = Http('Post', addr, .EncodeCall(args),
			header: #(Content_Type: 'text/xml'))
		if response.content.Blank?()
			throw "XmlRpc: no response from server"
		return (new this).Decode(response.header, response.content)
		}
	EncodeCall(args)
		// convert call with Suneido values to xml
		{
		return .preamble $
			Xml('methodCall',
				Xml('methodName', args[1]) $ '\r\n' $
				'<params>' $
				args[2..].Map{ Xml('param', .EncodeValue(it)) $ '\r\n' }.Join() $
				'</params>\r\n') $ '\r\n'
		}
	EncodeValue(x)
		{
		if Boolean?(x)
			x = Xml('boolean', x is true ? 1 : 0)
		else if Number?(x)
			x = Xml(x.Int() is x ? 'int' : 'double', x)
		else if String?(x)
			x = Xml('string', XmlEntityEncode(x))
		else if Date?(x)
			x = Xml('dateTime.iso8601', x.Format("yyyyMMddTHH:mm:ss"))
		else if Object?(x)
			{
			if x.HasNamed?()
				x = Xml('struct',
					x.Map2{|m,v| Xml('member', Xml('name', m) $ .EncodeValue(v)) }.Join())
			else
				x = Xml('array', Xml('data', x.Map(.EncodeValue).Join()))
			}
		else
			throw "XmlRpc: unhandled value: " $ Display(x)
		return Xml('value', x)
		}
	Decode(header, content)
		// convert xml to Suneido values
		{
		// strip header
		line = header.FirstLine()
		if line !~ "^HTTP.* 200"
			throw "XmlRpc: " $ line.AfterFirst(' ')
		type = "-none-"
		for line in header.Lines()
			{
			if line is ""
				break
			if line.Prefix?("Content-Type:")
				type = line[13 ..].Trim()
			}
		if type isnt "text/xml"
			throw "XmlRpc: unexpected Content-Type: " $ type
		return .Decode2(content)[0]
		}
	Decode2(msg)
		{
		.stack = new Stack
		.stack.Push(Object())
		.name = false

		xr = new XmlReader
		xr.SetContentHandler(this)
		xr.Parse(msg)

		return .stack.Top()
		}
	StartElement(qname, atts /*unused*/)
		{
		switch qname
			{
		case 'array', 'struct' :
			.stack.Push(Object())
		default:
			}
		}
	value: ''
	Characters(string)
		{
		.value = XmlEntityDecode(string)
		}
	EndElement(qname)
		{
		switch qname
			{
		case 'methodname' :
			.stack.Top().Add(.value)
		case 'array', 'struct' :
			.value = .stack.Pop()
		case 'boolean' :
			.value = .value is '1'
		case 'int', 'i4', 'double' :
			.value = Number(.value)
		case 'datetime.iso8601' :
			.value = Date('#' $ .value.Tr('T', '.').Tr(':'))
		case 'name' :
			.name = .value
		case 'value' :
			if .name is false
				.stack.Top().Add(.value)
			else
				{
				if .name.Number?()
					.name = Number(.name)
				.stack.Top()[.name] = .value
				.name = false
				}
		case 'fault' :
			throw 'XmlRpc: ' $ .value
		case 'string', 'member', 'data', 'param', 'params',
			'methodcall', 'methodresponse' :
			// ignore
		default :
			throw qname $ ' not handled'
			}
		}
	}