// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Display(token, indent = 0)
		{
		s = '  '.Repeat(indent) $
			token.Token $ Opt('(', token.GetDefault(#Value, ''), ')') $ '\r\n'
		for mem in token.Children.Members()
			s $= .Display(token.Children[mem], indent+1)
		return s
		}
	}
