// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (book, path, name)
	{
	// exact match?
	if false isnt x = Query1(book $
		' where BookPathStripReference(path) is ' $ Display(path) $
		' and name is ' $ Display(name) $
		' and text !~ "<!-- Hide -->"')
		return x

	pos = 999
	do
		{
		query = book $ ' where name is ' $ Display(name) $
			' and BookPathStripReference(path) =~ ' $ Display(path[pos..]) $
			' and text !~ "^GetBookPage" and text !~ "<!-- Hide -->"'
		WithQuery(query)
			{ |q|
			x = q.Next()
			if x isnt false and q.Next() is false
				return x
			}
		pos = path[.. pos].FindLast('/')
		} while pos isnt false

	// if you can't find page, try for parent
	if path isnt ''
		return BookPageFind(book, path.BeforeLast('/'), path.AfterLast('/'))

	return false
	}