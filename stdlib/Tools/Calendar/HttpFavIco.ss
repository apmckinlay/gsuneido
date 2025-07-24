// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
function(@unused)
	{
	headers = Object()
	headers['Cache_Control'] = 'max-age=484200'
	headers['Last_Modified'] = #20000101
	headers['Expires'] = Date().Plus(years: 1)

	icons = OptContribution('HttpFavIco', #(book: 'imagebook', icon: 'suneido.png'))
	return ['OK', headers, GetBookText(icons.icon, icons.book)]
	}
