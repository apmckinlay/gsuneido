// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (path)
	{
	book = path.BeforeFirst('/')
	name = path.AfterLast('/')
	path = '/' $ path.AfterFirst('/').BeforeLast('/')
	page = Query1(book $
		' where path = ' $ Display(path) $
		' and name = ' $ Display(name))
	if page is false
		return ""
	return page.text.Prefix?('<') ? page.text : page.text.Eval() // needs Eval
	}