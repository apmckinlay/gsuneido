// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	totalSize: 0
	maxSize: 2000
	outerObjectSizeLimit: 25
	nestedObjectSizeLimit: 100

	CallClass(ob, maxStrSize = 100)
		{
		return (new this).FormatLog(ob, maxStrSize)
		}

	FormatLog(ob, maxStrSize)
		{
		if not Instance?(ob) and not Object?(ob)
			return .format_value("", ob, maxStrSize)
		result = Object()
		// Copy is required because referencing an invalidated rule
		// can cause the object to be modified during iteration
		// even though we're only reading it
		for var_name in ob.Members()[.. .outerObjectSizeLimit].Copy()
			{
			try
				{
				val = .processValue(var_name, ob[var_name], maxStrSize)
				if .hasBufferMsg?(val)
					{
					result[var_name] = val
					break
					}
				}
			catch
				val = '???'
			result[var_name] = val
			}
		if ob.Size() > .outerObjectSizeLimit
			result.Append(Object('...': '...'))
		return result
		}

	// Buffer overflow occurs based on all contents of the Object, not just the values
	// Need to account for members also
	// i.e. when you Display() an object it converts the whole thing to a string which
		// must fit inside the buffer 65536
	bufferMsg: 'Stopped logging to prevent buffer overflow'
	processValue(member, val, maxStrSize)
		{
		if Instance?(val) or Object?(val)
			{
			.totalSize += String(member).Size()
			log_ob = Object()
			for var2 in val.Members()[.. .nestedObjectSizeLimit].Copy()
				{
				tmp = .format_value(var2, val[var2], maxStrSize)
				.totalSize += String(var2).Size() + String(tmp).Size()
				if .totalSize > .maxSize
					{
					log_ob[var2] = .bufferMsg
					break
					}
				log_ob[var2] = tmp
				}
			if val.Size() > .nestedObjectSizeLimit
				log_ob.Append(Object('...': '...'))
			val = log_ob
			}
		else
			{
			val = .format_value(member, val, maxStrSize)
			.totalSize += String(member).Size() + String(val).Size()
			if .totalSize > .maxSize
				return .bufferMsg
			}
		Pack(val) // ensure it's save-able
		return val
		}

	format_value(member, val, maxStrSize)
		{
		if .privateData?(member)
			val = '***'
		else if Instance?(val)
			val = '<instance>'
		else if Object?(val)
			val = .formatObject(val)
		else if Function?(val) or Class?(val)
			val = Display(val)
		else if String?(val)
			val = .formatString(val, maxStrSize)
		else if not .packable?(val)
			val = "<" $ Type(val) $ ">"
		return val
		}

	formatObject(val)
		{
		if val.Values().Size() > 0
			return '<object>'
		else
			return '<emptyobject>'
		}

	privateData?(member)
		{
		if not String?(member)
			return false
		member = member.Lower()
		memberParts = member.Split('_')
		if #(pass passwd pw pwd).Intersects?(memberParts)
			return true
		if #(sin ssn sinssn ssnsin).Intersects?(memberParts)
			return true

		otherPrivateMembers = #(card, cvv, account, token, password, passphrase,
			passhash, hiddenlines, authorization)

		return member =~ otherPrivateMembers.Join('|')
		}

	formatString(val, maxStrSize)
		{
		if val.Size() < maxStrSize
			return val
		if not val.Has?('<')
			return val.Ellipsis(maxStrSize)
		return .ensureBraceClosed(val[.. maxStrSize/2]) $ '...' $ val[-maxStrSize/2 ..]
		}
	ensureBraceClosed(str)
		{
		if not str.Has?('<')
			return str
		start = str.BeforeLast('<')
		end = str.AfterLast('<')
		return start $ (end.Has?('>') ? end : '')
		}

	packable?(val)
		{
		return Type(val) in ('Number', 'Date', 'Boolean', 'String')
		}

	hasBufferMsg?(val)
		{
		return val is .bufferMsg or (Object?(val) and val.Has?(.bufferMsg))
		}
	}
