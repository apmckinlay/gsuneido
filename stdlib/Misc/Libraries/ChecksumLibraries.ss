// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(libs = false, record_results = false, books = false)
		{
		cksums = Object()
		if record_results is true
			DeleteFile('checksum_results')

		.calc_libs_cksums(libs, cksums, record_results)
		.calc_books_cksums(books, cksums, record_results)
		return cksums
		}
	calc_libs_cksums(libs, cksums, record_results)
		{
		if libs is false
			libs = Libraries()
		for lib in libs
			.Calc_cksums(lib, lib $ ' where group is -1 sort name',
				record_results, cksums)
		}
	calc_books_cksums(books, cksums, record_results)
		{
		if books is false
			books = BookTables()
		for book in books
			{
			BookModel.Create(book) // need this to make sure the plugin column exists
			.Calc_cksums(book, book $ ' where plugin isnt true and
				not path.Has?("Reporter Reports") and
				not path.Has?("Reporter Forms") and
				not (name is "libCommitted")
				sort path, name', record_results, cksums)
			}
		}
	Calc_cksums(table, query, record_results, cksums)
		{
		max_mod = Date.Begin()
		cksum = Adler32()
		if TableExists?(table)
			{
			QueryApply(query)
				{ |x|
				max_mod = Max(max_mod, x.lib_committed, x.lib_modified)
				if x.path !~ "^/res\>"
					x.text = x.text.Tr('\r').Trim()
				cksum.Update(x.name.Trim())
				cksum.Update(x.text)
				if record_results isnt false
					AddFile('checksum_results', x.name.RightFill(30) $ '\t\t' $
						.fmt(Adler32(x.text)) $ '\r\n')
				}
			cksums.Add(Object(lib: table, max_mod: max_mod.StdShortDateTime(),
				cksum: .fmt(cksum.Value())))
			}
		}
	fmt(n)
		{ FormatChecksum(n) }
	}