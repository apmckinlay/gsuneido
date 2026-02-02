// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// convert Suneido values to JSON
	Encode(value)
		{
		return Object?(value)
			? (value.HasNamed?() ? .encodeObject(value) : .encodeArray(value))
			: .encodeValue(value)
		}

	encodeArray(ob)
		{
		return '[' $ ob.Map(.Encode).Join(',') $ ']'
		}

	encodeObject(ob)
		{
		s = ''
		for m in ob.Members().Sort!()
			s $= ',' $ .encodeString(m) $ ':' $ .Encode(ob[m])
		return '{' $ s[1..] $ '}'
		}

	encodeValue(value)
		{
		if Number?(value)
			return value <= -1 or value is 0 or 1 <= value
				? String(value)
				: (value < 0 ? "-0" : "0") $ value.Abs()
		else
			return Boolean?(value)
				? String(value)
				: .encodeString(value)
		}

	encodeString(value)
		{
		return '"' $ .escape(String(value)) $ '"'
		}
	escape(s)
		{
		// escape double quote, backslash, and control chars < space
		// NOTE: this only handles ASCII, not UTF-8
		return s.Replace('[\\\\"\x00-\x1f]', function (c)
			{
			#('"': `\"`, `\`: `\\`,
				'\x00': '\u0000', '\x01': '\u0001', '\x02': '\u0002', '\x03': '\u0003',
				'\x04': '\u0004', '\x05': '\u0005', '\x06': '\u0006', '\x07': '\u0007',
				'\x08': '\u0008', '\x09': '\\t', 	'\x0a': '\\n',    '\x0b': '\u000b',
				'\x0c': '\u000c', '\x0d': '\\r',    '\x0e': '\u000e', '\x0f': '\u000f',
				'\x10': '\u0010', '\x11': '\u0011', '\x12': '\u0012', '\x13': '\u0013',
				'\x14': '\u0014', '\x15': '\u0015', '\x16': '\u0016', '\x17': '\u0017',
				'\x18': '\u0018', '\x19': '\u0019', '\x1a': '\u001a', '\x1b': '\u001b',
				'\x1c': '\u001c', '\x1d': '\u001d', '\x1e': '\u001e', '\x1f': '\u001f')[c]
			})
		}

	EncodeWithFormat(value, level = 0, indent = '\t')
		{
		return Object?(value)
			? (value.HasNamed?()
				? .encodeObjectWithFormat(value, level, indent)
				: .encodeArrayWithFormat(value, level, indent))
			: .encodeValue(value)
		}

	encodeArrayWithFormat(ob, level = 0, indent = '\t')
		{
		return '[' $
			ob.Map({ '\n' $ indent.Repeat(level+1) $
				.EncodeWithFormat(it, level+1, indent) }).Join(',') $ '\n' $
			indent.Repeat(level) $ ']'
		}

	encodeObjectWithFormat(ob, level = 0, indent = '\t')
		{
		s = '{'
		for m in ob.Members().Sort!()
			s $= '\n' $ indent.Repeat(level+1) $ .encodeString(m) $ ': ' $
				.EncodeWithFormat(ob[m], level+1, indent) $ ','
		return s[..-1] $ '\n' $ indent.Repeat(level) $ '}'
		}

	//====================================================================================

	errMsg: 'Invalid Json format'

	// convert JSON strings to Suneido values
	// handleNull should be:
	// 		'throw' 	- throws program error
	//		'empty'		- treats null as ''
	//		'skip' 		- won't create the member
	Decode(s, handleNull = 'throw')
		{
		try
			{
			if not String?(s)
				throw "string required"
			toks = Scanner(s)
			_handleNull = handleNull
			if .skipNull is value = .decode(toks)
				throw .unexpectedNull
			if .next(toks) isnt toks
				throw "extra text at end"
			return value
			}
		catch (e)
			throw .errMsg $ ": " $ e
		}
	decode(toks, delim = false)
		{
		switch .next(toks)
			{
		case '-':
			return .decodeMinus(toks)
		case #NUMBER:
			return Number(toks.Value())
		case #STRING:
			return .decodeString(toks.Value())
		case #IDENTIFIER:
			return .decodeIdentifier(toks.Text())
		case '[':
			return .decodeContainer(toks, .decodeArrayEntry)
		case '{':
			return .decodeContainer(toks, .decodeObjectEntry)
		case delim: // used for empty array/object
			return toks
		case toks:
			throw "unexpected end of string"
		default:
			throw "unexpected: " $ toks.Text()
			}
		}
	decodeMinus(toks)
		{
		if #NUMBER isnt .next(toks)
			throw "unexpected: " $ toks.Text()
		return -Number(toks.Value())
		}
	decodeString(s)
		{
		return s.Replace('\\\u\d\d\d\d')
			{
			lastTwo = it[4::2] /*= last two */
			WideCharToMultiByte((lastTwo $ it[2::2]).FromHex() $ '\x00\x00')
			}
		}
	decodeIdentifier(tok)
		{
		switch tok
			{
		case 'true':	return true
		case 'false':	return false
		case 'null': 	return .null()
		default:		throw "unexpected: " $ tok
			}
		}
	skipNull: class {  } // to not conflict with real json values
	unexpectedNull: 'data should not contain null'
	null(_handleNull)
		{
		switch(handleNull)
			{
			case 'throw' :
				throw .unexpectedNull
			case 'empty':
				return ''
			case 'skip':
				return .skipNull
			}
		}
	decodeContainer(toks, entry)
		{
		end = toks.Text() is '[' ? ']' : '}'
		ob = Object()
		if toks is entry(toks, ob) // end delimiter
			return ob
		while end isnt tok = .next(toks)
			{
			if tok is toks
				throw "unexpected end of string"
			if tok isnt ','
				throw "missing comma"
			entry(toks, ob)
			}
		return ob
		}
	decodeArrayEntry(toks, ob)
		{
		if toks is value = .decode(toks, delim: ']')
			return toks // end delimiter
		if value isnt .skipNull
			ob.Add(value)
		return true
		}
	decodeObjectEntry(toks, ob)
		{
		if toks is name = .decode(toks, delim: '}')
			return toks // end delimiter
		Assert(String?(name))
		if ':' isnt .next(toks)
			throw "missing ':'"
		if .skipNull isnt value = .decode(toks)
			ob[name] = value
		return true
		}
	next(toks)
		{
		while #WHITESPACE is (tok = toks.Next2()) or #NEWLINE is tok
			{ }
		return tok is "" ? toks.Text() : tok
		}
	}
