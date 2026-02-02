// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function (name, _outputType = #html)
	{
	if false is x = Query1(_table, path: '/res/.snippets', :name)
		return '!!!SNIPPET NOT FOUND!!!'

	if outputType is #md and BookContent.Type(_table) is #md
		return x.text

	return BookContent.ToHtml(_table, x.text)
	}