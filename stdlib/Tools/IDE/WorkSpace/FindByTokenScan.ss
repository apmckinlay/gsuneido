// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(s)
		{
		return Sequence(new .iterator(s))
		}
	Iterator(s)
		{
		return new .iterator(s)
		}
	iterator: class
		{
		New(.s)
			{
			.scan = Scanner(s)
			}
		Next()
			{
			while .scan isnt type = .scan.Next2()
				{
				if type in (#COMMENT, #WHITESPACE, #NEWLINE)
					continue
				tok = .scan.Text()
				if type is #STRING
					tok = .stdQuote(tok)
				return tok
				}
			return this
			}
		Position()
			{
			return .scan.Position()
			}
		stdQuote(s)
			{
			//NOTE: result may not be correctly escaped
			if s[0] in ('"', '`', "'")
				s = '"' $ s[1 .. -1] $ '"'
			return s
			}
		Dup()
			{
			return new (.Base())(.s)
			}
		Infinite?()
			{
			return false
			}
		}
	}