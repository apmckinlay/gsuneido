// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(@args)
		// pre:	arg[0] is a tag name string,
		//		arg[1] is optional content string AND remaining args are attributes
		//		OR
		//		arg[1] is object with optional content and optional attributes
		// post:	returns a string in the form: <tagname /> OR
		//			<tagname name="value" ...> content </tagname>
		{
		attrs = content = ""
		tagname = args.PopFirst()
		if args.Member?(0) and Object?(args[0])
			args = args[0]
		for name in args.Members().Sort!()
			// sort members so tests get consistent order
			{
			if name is 0
				content = args[0]
			else if name is 'block'
				content = args.block
			else
				attrs $= ' ' $ name $ '="' $ XmlEntityEncode(String(args[name])) $ '"'
			}
		return .buildXml(tagname, attrs, content)
		}

	buildXml(tagname, attrs, content)
		{
		result = '<' $ tagname $ attrs
		if Type(content) is 'Block'
			content = content()

		if tagname.Prefix?('?') // processing instruction
			result $= "?>"
		else if tagname is '!--' // comment
			result $= content $ "-->"
		else if content is ""
			result $= " />"
		else
			result $= '>' $ content $ '</' $ tagname $ '>'
		return result
		}
	}