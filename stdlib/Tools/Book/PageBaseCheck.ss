// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
class
	{
	ForeachBookOption(book, block = function (@unused) { })
		{
		errs = ""
		bookname = book.BeforeFirst(' ')
		QueryApply(book $ .where())
			{ |x|
			if .ignorePage(bookname, x)
				continue
			try
				{
				ctrl = x.text.Eval() // Eval handles object, function, function call
				if not String?(ctrl) or not ctrl.Prefix?('<')
					block(ctrl, x.name)
				}
			catch (err)
				{
				errs $= x.path $ '/' $ x.name $ '\t' $ err $ '\n'
				SujsAdapter.CallOnRenderBackend(#CancelAllReserved)
				}
			}
		return errs
		}

	where()
		{
		contribs = GetContributions(#BookCheck).Join('|')
		pathsWhere = ' extend pathName = path $ "/" $ name
			where pathName =~ "^/(' $ (contribs.Blank?() ? 'nonbook' : contribs) $ ')"'
		return pathsWhere $ ' sort path, order'
		}

	ignorePage(book, x)
		{
		for fn in GetContributions('BookCheck_IgnorePage')
			if fn(:book, :x) is true
				return true
		return false
		}
	}