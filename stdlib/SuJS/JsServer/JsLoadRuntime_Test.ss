// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		.book = .MakeBook()
		SvcTable(.book, svcEnsure:)
		}
	cl: JsLoadRuntime
		{
		New(book, .info)
			{
			.JsLoadRuntime_book = book
			}
		JsLoadRuntime_maxChunkSize: 10
		JsLoadRuntime_tags()
			{
			return .info.tags
			}
		JsLoadRuntime_tempName()
			{
			return 'import' $ .info.n++
			}
		JsLoadRuntime_fileSize(name)
			{
			return .info.files[name].Size()
			}
		JsLoadRuntime_readFile(name, block)
			{
			FakeFile(.info.files[name], block)
			}
		JsLoadRuntime_putFile(name, s)
			{
			Assert(not .info.putFiles.Member?(name))
			.info.putFiles[name] = s
			}
		JsLoadRuntime_deleteFile(name)
			{
			.info.deleteFiles.Add(name)
			}
		importCl: ImportSvcTableText
			{
			ImportSvcTableText_getFile(filename, hwnd/*unused*/, quiet/*unused*/,
				skipSize/*unused*/)
				{
				if filename.Prefix?('import')
					return _info.putFiles[filename]
				else
					return _info.files[filename]
				}
			}
		JsLoadRuntime_import(@args)
			{
			_info = .info
			(.importCl)(@args)
			}
		}
	Test_prepareFiles()
		{
		info = Object(n: 0,
			files: Object(
				'short.js': '1234567890',
				'long.js': '1234567890'.Repeat(3) $ '1234'),
			putFiles: Object())
		cl = new .cl(.book, info)
		fn = cl.JsLoadRuntime_prepareFiles

		Assert(fn('short.js') is: #([path: "short.js", temp?: false, partNum: ""]))
		Assert(fn('long.js') is: [[path: "import0", temp?:, partNum: 0],
			[path: "import1", temp?:, partNum: 1],
			[path: "import2", temp?:, partNum: 2],
			[path: "import3", temp?:, partNum: 3]])
		Assert(info.putFiles.import0 is: '1234567890')
		Assert(info.putFiles.import1 is: '1234567890')
		Assert(info.putFiles.import2 is: '1234567890')
		Assert(info.putFiles.import3 is: '1234')
		}

	Test_QueryRec()
		{
		// queryRec_not_found
		info = Object(tags: #())
		cl = new .cl(.book, info)

		result = cl.QueryRec('missing.js')
		Assert(result is: false)

		// queryRec_not_found_with_tags
		info = Object(tags: #('__prod', '__test'))
		cl = new .cl(.book, info)

		result = cl.QueryRec('missing.js')
		Assert(result is: false)

		// queryRec_single_record
		.outputRec('file content', 'file__prod.js', #20250101, #20250102)

		info = Object(tags: #('__prod'))
		cl = new .cl(.book, info)
		result = cl.QueryRec('file.js')
		Assert(result hasSubset: [text: 'file content',
			lib_committed: #20250102, lib_modified: #20250101])

		// queryRec_two_chunks
		.outputRec('part one ', 'chunk__prod.part0.js', #20250101, #20250102)
		.outputRec('part two', 'chunk__prod.part1.js', #20250103, #20250104)

		info = Object(tags: #('__prod'))
		cl = new .cl(.book, info)
		result = cl.QueryRec('chunk.js')
		Assert(result hasSubset: [text: 'part one part two',
			lib_modified: #20250101, lib_committed: #20250102])

		// queryRec_three_chunks
		.outputRec('AAA', 'three__prod.part0.js', #20250101)
		.outputRec('BBB', 'three__prod.part1.js')
		.outputRec('CCC', 'three__prod.part2.js')

		info = Object(tags: #('__prod'))
		cl = new .cl(.book, info)
		result = cl.QueryRec('three.js')
		Assert(result hasSubset: [text: 'AAABBBCCC', lib_modified: #20250101])

		// queryRec_tag_priority - tags searched in reverse order
		.outputRec('prod version', 'priority__prod.js')
		.outputRec('test version', 'priority__test.js')

		info = Object(tags: #('__prod', '__test'))
		cl = new .cl(.book, info)
		result = cl.QueryRec('priority.js')
		Assert(result hasSubset: [text: 'test version']) // __test has priority

		// queryRec_tag_fallback - first tag not found, falls back to second
		.outputRec('fallback prod', 'fallback__prod.js')

		info = Object(tags: #('__prod', '__test'))
		cl = new .cl(.book, info)
		result = cl.QueryRec('fallback.js')
		Assert(result hasSubset: [text: 'fallback prod'])

		// queryRec_chunked_over_single
		.outputRec('chunked ', 'both__prod.part0.js', #20250101)
		.outputRec('content', 'both__prod.part1.js')
		.outputRec('single version', 'both__prod.js', #20250102)

		info = Object(tags: #('__prod'))
		cl = new .cl(.book, info)
		result = cl.QueryRec('both.js')
		Assert(result hasSubset: [text: 'chunked content', lib_modified: #20250101])
		}

	Test_importFiles_queryRec()
		{
		// import_query_small_file - file smaller than maxChunkSize
		info = Object(
			n: 0,
			tags: #('__prod'),
			files: Object('small.js': 'hello'),
			putFiles: Object(),
			deleteFiles: Object())
		cl = new .cl(.book, info)

		cl.JsLoadRuntime_importFiles('small.js', 'small__prod.js')

		// Should create single record (no chunks)
		Assert(info.putFiles.Member?('import0') is: false) // no temp file for small
		result = cl.QueryRec('small.js')
		Assert(result isnt: false)
		Assert(result.text is: 'hello')
		Assert(info.deleteFiles is: #())

		// import_query_large_file - file larger than maxChunkSize (10 chars)
		info = Object(
			n: 0,
			tags: #('__prod'),
			files: Object('large.js': '1234567890123456789012345'), // 25 chars
			putFiles: Object(),
			deleteFiles: Object())
		cl = new .cl(.book, info)

		cl.JsLoadRuntime_importFiles('large.js', 'large__prod.js')

		// Should create 3 chunks: 10 + 10 + 5
		Assert(info.putFiles.import0 is: '1234567890')
		Assert(info.putFiles.import1 is: '1234567890')
		Assert(info.putFiles.import2 is: '12345')
		Assert(info.deleteFiles is: #("import0", "import1", "import2"))
		result = cl.QueryRec('large.js')
		Assert(result isnt: false)
		Assert(result.text is: '1234567890123456789012345')

		// import_query_exact_boundary - file exactly at chunk boundary
		info = Object(
			n: 0,
			tags: #('__prod'),
			files: Object('exact.js': '1234567890'), // exactly 10 chars
			putFiles: Object(),
			deleteFiles: Object())
		cl = new .cl(.book, info)

		cl.JsLoadRuntime_importFiles('exact.js', 'exact__prod.js')

		result = cl.QueryRec('exact.js')
		Assert(result isnt: false)
		Assert(result.text is: '1234567890')
		Assert(info.putFiles is: #())

		// import_replaces_existing - re-import replaces old chunks
		info = Object(
			n: 0,
			tags: #('__prod'),
			files: Object('replace.js': 'updated content'),
			putFiles: Object(),
			deleteFiles: Object())
		cl = new .cl(.book, info)

		// First import
		cl.JsLoadRuntime_importFiles('replace.js', 'replace__prod.js')
		result = cl.QueryRec('replace.js')
		Assert(result.text is: 'updated content')

		// Re-import with different content
		info.files['replace.js'] = 'new data here'
		info.putFiles = Object()
		info.n = 0
		cl.JsLoadRuntime_importFiles('replace.js', 'replace__prod.js')
		result = cl.QueryRec('replace.js')
		Assert(result.text is: 'new data here')

		// import_with_tag - verify tag is used in record name
		info = Object(
			n: 0,
			tags: #('__prod'),
			files: Object('tagged.js': 'tagged content'),
			putFiles: Object(),
			deleteFiles: Object())
		cl = new .cl(.book, info)

		cl.JsLoadRuntime_importFiles('tagged.js', 'tagged__prod.js')

		// Query with matching tag should find it
		result = cl.QueryRec('tagged.js')
		Assert(result isnt: false)
		Assert(result.text is: 'tagged content')

		// Query with different tag should not find it
		info = Object(tags: #('__test'))
		cl = new .cl(.book, info)
		result = cl.QueryRec('tagged.js')
		Assert(result is: false)
		}

	outputRec(text, name, lib_modified = false, lib_committed = false)
		{
		rec = [:text, path: '/res', :name, num: NextTableNum(.book)]
		if lib_modified isnt false
			rec.lib_modified = lib_modified
		if lib_committed isnt false
			rec.lib_committed = lib_committed
		QueryOutput(.book, rec)
		}
	}