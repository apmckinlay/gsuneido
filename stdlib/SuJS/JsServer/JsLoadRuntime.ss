// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	maxChunkSize: 102400 // 100.Kb()
	book: 'imagebook'
	CallClass(env)
		{
		name = env.path.AfterFirst('runtime/')
		headers = Object('Cache-Control': 'max-age=1209600'/*=two weeks*/)

		if false is rec = .getRec(name)
			{
			SuneidoLog("ERRATIC: JsLoadRuntime - Can't find " $ Display(name))
			return ['404' /*= not found */, #(), '']
			}

		if NotModified?(env, rec.GetDefault('lib_modified', rec.lib_committed))
			return ['304', headers, '']

		return ['200', headers, rec.text]
		}

	getRec(name)
		{
		switch (name)
			{
		case 'su_code_bundle.js':
			return [text: SuCode().CodeBundle.code,
				lib_modified: Max(
					SuCode().CodeBundle.code_built.lib_modified,
					SuCode().CodeBundle.code_built.lib_committed),
				hash: SuCode().CodeBundle.code_built.hash]
		case 'su_bundle.js', 'su_bundle.min.js', 'su_bundle.min.js.map',
			'codemirror.css', 'foldgutter.css', 'codemirror_bundle.js',
			'suneido.ttf', 'suneido2.ttf':
			return .queryRecCached(name)
		default:
			return false
			}
		}

	queryRecCached(name)
		{
		if not Suneido.Member?('JsLoadRuntime.queryRecCached')
			Suneido['JsLoadRuntime.queryRecCached'] = LruCache(.QueryRec)
		return Suneido['JsLoadRuntime.queryRecCached'].Get(name)
		}

	QueryRec(name)
		{
		tags = .tags()
		for (i = tags.Size() - 1; i >= 0; i--)
			{
			taggedName = .buildName(name, tags[i])
			result = false
			partNum = 0
			while false isnt rec = Query1(.book, path: '/res',
				name: .buildChunkName(taggedName, partNum++))
				{
				// Preserve metadata from first chunk
				if result is false
					result = rec
				else
					result.text $= rec.text
				}
			if result isnt false
				return result

			// fall back to the name with part
			if false isnt rec = Query1(.book, path: '/res', name: taggedName)
				return rec
			}
		return false
		}

	buildName(name, tag)
		{
		return name.BeforeFirst('.') $ tag $ '.' $ name.AfterFirst('.')
		}

	buildChunkName(recName, partNum)
		{
		if partNum is ''
			return recName
		return recName.BeforeFirst('.') $ '.part' $ partNum $ '.' $
			recName.AfterFirst('.')
		}

	GetUrl(name)
		{
		if false is rec = .getRec(name)
			return "/runtime/" $ name

		if rec.Member?(#hash)
			return "/runtime/" $ name $ '?id=' $ rec.hash

		return "/runtime/" $ name $ '?date=' $
			rec.GetDefault('lib_modified', rec.lib_committed).Format("yyyyMMddHHmmss")
		}


	Import(projectDir, tag = '') // Called manually
		{
		path = Paths.Combine(projectDir, 'runtime')
		for file in #('su_bundle.js', 'su_bundle.min.js', 'su_bundle.min.js.map',
				'su_global_builtins.json')
			{
			filename = Paths.Combine(path, file)
			recName = .buildName(file, Opt('__', tag))
			.importFiles(filename, recName)
			}
		ResetCaches()
		}

	importFiles(filename, recName)
		{
		.deleteChunks(recName)
		files = .prepareFiles(filename)
		Finally(
			{ files.Each({ .import(it.path, .book, '/res',
					quiet:, recName: .buildChunkName(recName, it.partNum)) }) },
			{ files.Filter({ it.temp? }).Each({ .deleteFile(it.path) }) })
		}

	prepareFiles(filename)
		{
		if .fileSize(filename) <= .maxChunkSize
			return [[path: filename, temp?: false, partNum: '']]
		files = []
		partNum = 0
		.readFile(filename)
			{ |f|
			while false isnt chunk = f.Read(.maxChunkSize)
				{
				tempFile = .tempName()
				.putFile(tempFile, chunk)
				files.Add([path: tempFile, temp?:, partNum: partNum++])
				}
			}
		return files
		}

	deleteChunks(recName)
		{
		recName = '/res/' $ recName
		svcTable = SvcTable(.book)

//		if false isnt svcTable.Get(recName)
//			svcTable.StageDelete(recName)

		partNum = 0
		while true
			{
			chunkName = .buildChunkName(recName, partNum++)
			if false is svcTable.Get(chunkName)
				return

			svcTable.StageDelete(chunkName)
			}
		}

	ImportCodeMirror(projectDir, tag = '') // Called manually
		{
		file = 'codemirror_bundle.js'
		filename = Paths.Combine(projectDir, 'CodeMirror', file)
		recName = .buildName(file, Opt('__', tag))
		.importFiles(filename, recName)
		ResetCaches()
		}

	// for test
	import(@args)
		{
		ImportSvcTableText(@args)
		}

	tags()
		{
		return LibraryTags.GetTagsInUse().Filter({ not it.Has?(#webgui) })
		}

	tempName()
		{
		return GetTempFileName(GetTempPath(), "import")
		}

	fileSize(name)
		{
		return FileSize(name)
		}

	readFile(name, block)
		{
		File(name, :block)
		}

	putFile(name, s)
		{
		PutFile(name, s)
		}

	deleteFile(name)
		{
		if true isnt result = DeleteFile(name)
			Print('delete ' $ name $ ' error: ' $ result)
		}
	}
