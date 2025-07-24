// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(book, option, module)
		{
		origQuery = bookOptionQuery = .buildQuery(book, option)
		if 0 is count = QueryCount(bookOptionQuery)
			return false
		if count > 1
			{
			prefix = "/" $ module
			extraWhere = ' where path.Prefix?(' $ Display(prefix) $ ')'
			bookOptionQuery = .buildQuery(book, option, extraWhere)
			}
		if false is rec = QueryFirst(bookOptionQuery $ ' sort path')
			{
			menus = .bookMenus(book)
			bookOptions = .sortBookOptions(QueryAll(origQuery), menus)
			for bookRec in bookOptions
				{
				if true is perm = .permission(bookRec)
					return bookRec
				else if perm is 'readOnly' and rec is false
					rec = bookRec
				}
			}
		return rec
		}

	buildQuery(book, option, extraWhere = '')
		{
		return QueryAddWhere(BookOptionQuery(book, option), extraWhere)
		}

	bookMenus(book)
		{
		return QueryAll(book $ ` where path is "" and name isnt "res" ` $
				`project name, order sort order`)
		}

	sortBookOptions(bookOptions, menus)
		{
		bookOptions.Each()
			{|x|
			// should always be able to find a match, but just in case
			if false isnt match = menus.FindOne({ it.name is x.path.AfterFirst('/') })
				x.order = match.order
			}
		return bookOptions.Sort!({|x,y| x.order < y.order})
		}

	permission(bookRec)
		{
		option = bookRec.path $ '/' $ bookRec.name
		return AccessPermissions(option)
		}
	}