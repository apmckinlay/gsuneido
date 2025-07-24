// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(path)
		{
		book = path[1 ..].BeforeFirst('/')
		s = '<p><a href="/' $ book $ '/General/Help for New Users">' $
			'<img src="suneido:/' $ book $ '/res/newusers.png" align="middle"> Help ' $
			'for New Users</a></p>'
		fullpath = path
		path = path[book.Size() + 1 ..]
		records = BookModel(book).Children(path).Copy().
			RemoveIf({|x| x.name is 'Cover' or x.name is 'Contents' })
		for rec in records
			{
			authorize = BookEnabled(book, path $ '/' $ rec.name)
			if authorize is "hidden"
				continue
			s $= .GetLink(rec, fullpath, book) $ '\n'
			}
		return s
		}
	GetLink(rec, path, book)
		{
		image = .get_image(rec.name, book)
		if image is ''
			return ''
		return '<p><a href="' $ path $ '/' $ rec.name $ '">' $ image $ rec.name $
			(not (BookModel(book).Children(path[book.Size() + 1 ..] $
				'/' $ rec.name)).Empty?() ? ' <font face="Wingdings">\xd8</font>' : '') $
			'</a></p>'
		}
	images: (
		'Concepts': 'key.png',
		'How Do I ...?': 'question.png',
		'Reference': 'book2.png',
		'Tips': 'blacklightbulb.png',
		'Troubleshooting': 'trouble.png'
		)
	get_image(str, book)
		{
		for s in .images.Members()
			if str.Has?(s) and .query?(book, .images[s])
				return '<img src="suneido:/' $ book $ '/res/' $ .images[s] $
					'" align="middle"> '
		return ''
		}
	query?(book, name)
		{
		return false isnt Query1Cached(book $
			' where path =~ "^/res\>" and name is ' $ Display(name))
		}
	}
