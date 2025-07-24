// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.

// TODO: add option to create instance
RefactorConvert
	{
	Name: 'Convert Function To Class'

	CanConvert(data)
		{
		if ScannerWithContext(data.text).Next() isnt 'function'
			return 'Can only convert functions'
		return true
		}

	Convert(text)
		{
		head = text.BeforeFirst('{').
			Replace('function[ \t]*', 'class\r\n\t{\r\n\tCallClass')
		body = text.AfterFirst('{').Replace('^\t', '\t\t').RightTrim()
		return head $ '\t{' $ body $ '\r\n\t}'
		}
	}