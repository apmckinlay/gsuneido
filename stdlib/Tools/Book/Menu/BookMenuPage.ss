// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// TODO: make headings work for multiple columns (they currently only work for 1 column)
class
	{
	orderMultiplier: 10
	CallClass(path, cols, before = false, after = false, headingLevel = 1, headings = #())
		{
		Assert(cols > 0)
		fullpath = path
		book = path[1..].BeforeFirst('/')
		path = path[book.Size() + 1..]
		title = path.AfterLast('/')
		if title is ''
			title = 'Contents'
		page = ''
		if headingLevel isnt 0
			page = .getHeading(title, headingLevel, book)
		if (before isnt false)
			page $= before
		records = BookModel(book).Children(path).Copy().
			RemoveIf({ it.name is 'Cover' or it.name is 'Contents' })
		done = Object()
		rows = (records.Size() / cols).Ceiling()
		page $= '\n<table width="100%">\n'
		nextorder = .orderMultiplier
		for (row = 0; row < rows; ++row)
			{
			s = '\n<tr>\n'
			for (col = 0; col < cols; ++col)
				{
				rec = col * rows + row
				item = records.Member?(rec) and not done.Has?(rec)
					? rec : (records.Member?(rec - 1) and not done.Has?(rec - 1)
						? rec - 1 : false)
				if item isnt false
					{
					if records[item].order is "" or
						records[item].order >= nextorder + .orderMultiplier
						nextorder += .orderMultiplier
					if headings.Member?(nextorder) and
						records[item].order.Int() % .orderMultiplier is 0 and
						records[item].order is nextorder
						{
						s $= '\t' $ Xml('td',
							Xml('h' $ (headingLevel + 1), headings[nextorder],
							style: 'padding-top=8px; margin-bottom=0')) $ '\n'
						nextorder += .orderMultiplier
						--row
						continue
						}

					authorize = BookEnabled(book, path $ '/' $ records[item].name)
					if authorize is "hidden"
						continue
					s $= '\t' $ .getLink(records[item], fullpath, book) $ '\n'
					done.Add(item)
					}
				}
			s $= '\n</tr>\n'
			if s.Has?('<td>')
				page $= s
			}
		page $= '\n</table>\n'
		if (after isnt false)
			page $= after
		return page
		}
	getHeading(title, level, book)
		{
		title = .getImage(title, book) $ title
		return Xml('h' $ level, title) $ '\n'
		}
	getLink(rec, path, book)
		{
		children = BookModel(book).Children(path[book.Size() + 1..] $ '/' $ rec.name)
		sideArrow = children.Empty?() ? '' : ' &raquo;'

		image = .getImage(rec.name, book)
		return Xml('td',
			Xml('a', image $ rec.name $ sideArrow, href: path $ '/' $ rec.name))
		}
	getImage(str, book)
		{
		images = GetContributions('BookMenuImages')
		for s in images.Members()
			if str.Has?(s) and .imageFound?(book, images[s])
				return '<img src="suneido:/' $ book $ '/res/' $ images[s] $
					 '" align="middle"> '
		return ''
		}
	imageFound?(book, name)
		{
		return false isnt Query1Cached(book $
			' where path is "/res" and name is ' $ Display(name))
		}

	}
