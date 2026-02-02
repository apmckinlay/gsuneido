// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(book)
		{
		.book = book
		.index_name = 'index_' $ book
		.field = .FindControl('search')
		}
	Controls: (Horz
		Skip
		(ScintillaFieldReturn name: 'search' status: 'Type your search and press Enter')
		(Skip 2)
		(EnhancedButton command: #Search image: 'zoom.emf' tip: 'Search'
			imagePadding: 0.1 mouseEffect:)
		)

	FieldReturn()
		{ .search() }
	On_Search()
		{ .search() }

	search()
		{
		s = .field.Get()
		if s is ""
			{
			Beep()
			return
			}
		for c in GetContributions('BookSearchReplacements')
			s = s.Replace(c[0], c[1])
		.close_results()

		results = .search_ftsearch(s, .book, .index_name)
		// fallback only if ftsearch has an error, not just if results empty
		if results is false
			results = .search_query(s, .book)

		BookLog("Book Search (" $ results.Size() $ " results) " $ s)

		.display_results(results)
		}
	display_results(results)
		{
		if results.Size() > 0
			results.Add('Vert', at: 0)
		else
			results = #(Static, "No matches")
		.w = Window(['BookSearchResults', this, .field.Hwnd, results],
			style: WS.BORDER | WS.POPUP, exStyle: WS_EX.TOOLWINDOW)
		}
	w: false
	close_results()
		{
		if .w is false
			return
		.w.Destroy()
		.w = false
		}

	search_ftsearch(s, book, indexName)
		{
		if false is results = .Ftsearch(s, book, indexName)
			return false

		formattedResults = Object()
		for result in results
			.add_result(result.path $ '/' $ result.name, book, formattedResults)
		return formattedResults
		}

	Ftsearch(s, book, indexName)
		{
		if Sys.Client?()
			return ServerEval('BookSearchControl.Ftsearch', s, book, indexName)

		if not .index_available?(indexName)
			{
			SuneidoLog("ERROR: Ftsearch index not available for BookSearch in book: " $
				book $ " - index may need to be regenerated")
			return false
			}

		resIDs = #()
		search = String(s)

		// cache the file content in Suneido,
		// server should have it already
		index = Suneido.GetInit(indexName, { GetFile(indexName) })

		try
			{
			idx = Ftsearch.Load(index)
			resIDs = idx.Search(search)
			}
		catch (err)
			{
			SuneidoLog('ERROR: ' $ err, params: Object(:search, ftindex: indexName))
			return false
			}

		return resIDs.Map({ Query1(book, num: it) }).Remove(false)
		}

	index_available?(indexName)
		{
		return FileExists?(indexName)
		}

	searchSize: 20
	search_query(s, book)
		{
		results = Object()
		s = '(?i)(?q)' $ s
		try
			Transaction(read:)
				{|t| .search_query2(t, s, :book, :results) }
		catch (unused, "toomany")
			results.Add(#(Static, "Too many matches, first 20 shown."))
		return results
		}
	search_query2(t, s, book, results = #(), path = "") // recursive
		{
		if path =~ "^/res\>"
			return
		_table = book
		t.QueryApply(book $ " where path = " $ Display(path) $
			" sort order, name")
			{|x|
			name = x.path $ "/" $ x.name
			_path = x.path
			_name = x.name
			if (BookContent.Match(book, x.text) and Asup(x.text) =~ s)
				{
				.add_result(name, book, results)
				if (results.Size() > .searchSize)
					throw "toomany"
				}
			.search_query2(t, s, :book, :results, path: name) // do children (if any)
			}
		}
	add_result(name, book, results = #())
		{
		if false isnt BookModel(book).Get(name, no_query:)
			results.Add(Object('Html_ahref_',
				name[1..].Replace('/', ' > '), name))
		}

	Destroy()
		{
		.close_results()
		super.Destroy()
		}
	}
