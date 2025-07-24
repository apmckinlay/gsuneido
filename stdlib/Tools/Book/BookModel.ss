// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(table)
		{
		bms = Suneido.GetInit(#BookModels, { Object() })
		return bms.GetInit(table, { new this(table) })
		}
	ClearCache(book = false)
		{
		ClearBookImageCache(book)
		}
	New(table)
		{
		.table = table
		}
	Load()
		{
		.children
		}
	getter_toc() // once only
		{
		.toc = Object()
		.depthfirst(.children)
		return .toc
		}
	loadedAt: false
	getter_children() // once only
		{
		.loadedAt = Date()
		if not TableExists?(.table)
			return Object().Set_default(#())
		if Sys.Client?()
			{
			.children = ServerEval("BookModel.GetAllChildren", .table)
			.children.Set_default(#())
			}
		else
			.load_children(.table)
		return .children
		}
	GetAllChildren(table)
		{
		return BookModel(table).BookModel_children
		}
	Create(table)
		{
		Database("ensure " $ table $
			" (path, name, num, order, text, plugin, hasSubmenu?, lib_modified)
			key(num) key (path, name) index(name)
			index(order, name) index(path, order, name) index (plugin)")
		SvcDisabledBooks.ResetCache()
		}
	load_children(table)
		{
		Assert(not Sys.Client?())
		if not TableExists?(table)
			return .toc = Object()
		.Create(table) // needed because old tables don't have plugin field
		.HandlePlugins(table)
		children = Object().Set_default(#())
		QueryApply(table $
			" where	name isnt 'res' and path !~ '^/res\>'
			project path, name, num, order, plugin, hasSubmenu?
			sort order, name")
			{|x|
			if BookEnabled(table, x.path $ "/" $ x.name) isnt "hidden"
				children[x.path].Add(x)
			}
		.children = children
		}
	depthfirst(toc, path = '')
		{
		for x in toc[path]
			{
			.toc.Add(x)
			p = x.path $ '/' $ x.name
			if toc.Member?(p)
				.depthfirst(toc, p)
			}
		}
	HandlePlugins(table)
		{
		RetryTransaction()
			{|t|
			pages = Object()
			Plugins().ForeachContribution('BookPages', 'pages')
				{|x|
				if x.book isnt table
					continue
				pages[x.path $ '/' $ x.name] = x
				}
			t.QueryApply(table $ ' where plugin is true')
				{|x|
				key = x.path $ '/' $ x.name
				if not pages.Member?(key)
					{
					x.Delete()
					continue
					}
				update = false
				for field in x.Members()
					if .setBookRecField?(x, pages[key], field)
						{
						x[field] = pages[key][field]
						update = true
						}
				if pages[key].GetDefault(#hasSubmenu?, false) is true
					{
					x[#hasSubmenu?] = true
					update = true
					}
				if update is true
					x.Update()
				pages.Delete(key)
				}
			for x in pages
				{
				if x.Member?('func') and (x.func)() is false
					continue
				rec = Record(name: x.name, order: x.order, path: x.path,
					text: x.text, hasSubmenu: x.GetDefault(#hasSubmenu?, false),
					plugin: true)
				rec.num = NextTableNum(table, t)

				if not t.QueryEmpty?(table, path: x.path, name: x.name)
					t.QueryDo('delete ' $ table $ ' where path is ' $ Display(x.path) $
						' and name is ' $ Display(x.name))// switching from book to plugin

				t.QueryOutput(table, rec)
				}
			}
		}

	setBookRecField?(bookRec, page, field)
		{
		return page.Member?(field) and page[field] isnt bookRec[field]
		}

	Get(pathname, no_query = false)
		{
		name = pathname.Extract('[^/]*$')
		path = pathname[.. -name.Size() - 1]
		i = .toc.FindIf({ it.path is path and it.name is name })
		if i is false and no_query
			return false
		// "i" will be false for new reporter report thus need to look it up in the table
		return i is false
			? Query1(.table $ " where path is " $ Display(path) $
				" and name is " $ Display(name))
			: .toc[i]
		}
	GetAllItems()
		{
		return .toc.Map({ it.path $ '/' $ it.name })
		}
	GetLoadedAt()
		{
		return .loadedAt
		}
	First()
		{
		return .toc.GetDefault(0, false)
		}
	Children(path = '') // '' for root
		{
		return .children[path]
		}
	Next(cur)
		{
		if false is i = .toc.Find(cur)
			return false
		return .toc.GetDefault(i + 1, false)
		}
	Prev(cur)
		{
		if false is i = .toc.Find(cur)
			return false
		return .toc.GetDefault(i - 1, false)
		}
	Parent(cur)
		{
		if cur is false
			return false
		path = cur.path.BeforeLast('/')
		name = cur.path.AfterLast('/')
		if not .children.Member?(path)
			return false
		i = .children[path].FindIf({ it.name is name })
		return i is false ? false : .children[path][i]
		}
	Children?(pathname)
		{
		return .children.Member?(pathname)
		}
	}
