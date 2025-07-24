// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// shorthand
function (s)
	{
	return Display(s) is 'HtmlString()' ? s.Value() : XmlEntityEncode(s)
	}