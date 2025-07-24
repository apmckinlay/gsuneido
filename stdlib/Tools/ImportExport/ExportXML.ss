// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Export
	{
	// Redefine for XML output format (want tags and header "columns" to match)
	GetHead(prompts)
		{
		.Head = Object()
		if .HeaderType is "Prompts"
			for field in .Fields
				.Head[field] = prompts is false ? Prompt(field) : prompts.Find(field)
		else
			for field in .Fields
				.Head[field] = field
		}

	Before()
		{
		.Putline('<?xml version="1.0"?>')
		.Putline('<!--  suneido xml export  -->')
		.Putline('<table>')
		}
	Export1(x)
		{
		.Putline('<record>')
		for (field in .Fields)
			{
			tag = .HeaderType is "Prompts" ? .Head[field] : field
			tag = .mapToXmlName(tag)
			s = XmlEntityEncode(String(x[field]))
			if s isnt ""
				.Putline(
					"<" $ tag $ ' type="' $ Type(x[field]) $ '">' $ s $ "</" $ tag $ ">")
			}
		.Putline('</record>')
		}
	After()
		{
		.Putline('</table>')
		}
	Ext: 'xml'

	mapToXmlName(str)
		{
		if str is ''
			return ':A:' // XML name cannot be empty.
		newstr = ''
		start_pos = 0
		if str[0] =~ '[:_a-zA-Z]'
			{
			newstr $= str[0]
			start_pos = 1
			}
		else  // put ':A:' at start of string to indicate Axon auto-formatted 1st char
			newstr $= ':A:'
		for (i = start_pos; i < str.Size() ; i++)
			if str[i] =~ '[:_a-zA-Z0-9.-]'
				newstr $= str[i]
			else
				newstr $= .charToHex(str[i])
		return newstr
		}

	charToHex(chr) // UTF-8 is the default encoding
		{
		return '0x' $ chr.Asc().Hex()
		}
	}
