// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// NOTE: src could be a path or github repo, like "https://github.com/apmckinlay/gsuneido"
	CallClass(src, book)
		{
		new this(src, book)
		return
		}
	New(src, book)
		{
		origSrc = src
		if src.Prefix?('https://github.com/')
			src = .downloadFromGithub(src)
		src = Paths.ToStd(src)
		src = src.RightTrim('/')
		.n = 0
		.totalSize = 0
		if not TableExists?(book)
			BookModel.Create(book)
		svcTable = SvcTable(book)
		QueryDo('delete ' $ book)
		.importDir(svcTable, src, src)
		if origSrc.Prefix?('https://github.com/')
			DeleteDir(src)
		Print('Imported from "' $ origSrc $ '" to ' $ book $
			', files: ' $ .n $ ', size: ' $ ReadableSize(.totalSize))
		}

	downloadFromGithub(url)
		{
		rest = url.AfterFirst('github.com/')
		parts = rest.Split('/')
		owner = parts[0]
		repo = parts[1]
		r = Json.Decode(Https.Get('https://api.github.com/repos/' $ owner $ '/' $ repo),
			handleNull: 'skip')
		branch = r.GetDefault("default_branch", "main")
		zipUrl = 'https://github.com/' $ owner $ '/' $ repo $
			'/archive/refs/heads/' $ branch $ '.zip'
		tempDir = GetTempPath().RightTrim(`\/`) $ '/ImportSrcCode_' $
			Timestamp().Format('yyyyMMddHHmmss')
		tempZip = tempDir $ '.zip'
		EnsureDirectories(tempDir)
		try
			{
			return .unzip(zipUrl, tempZip, tempDir, repo, branch)
			}
		catch (e)
			{
			DeleteDir(tempDir)
			DeleteFile(tempZip)
			throw e
			}
		}

	unzip(zipUrl, tempZip, tempDir, repo, branch)
		{
		cmd = PowerShell() $
			` -NoProfile -Command "Invoke-WebRequest -UseBasicParsing -Uri '` $
			zipUrl $ `' -OutFile '` $ tempZip $ `'; Expand-Archive -LiteralPath '` $
			tempZip $ `' -DestinationPath '` $ tempDir $ `' -Force"`
		result = RunPipedOutput.WithExitValue(cmd)
		if result.exitValue isnt 0
			throw 'ImportSourceCodeToBook: unzip failed: ' $ tempZip
		DeleteFile(tempZip)
		return tempDir $ '/' $ repo $ '-' $ branch
		}

	importDir(svcTable, rootSrc, dir)
		{
		for item in Dir(dir $ '/*.*', details:)
			{
			if item.name.Suffix?('/')
				{
				name = item.name.Trim('/')
				path = .makePath(rootSrc, dir $ '/' $ name)
				importRec = Record(:name, :path, text: '')
				.output(svcTable, importRec)
				.importDir(svcTable, rootSrc, dir $ '/' $ name)
				}
			else if item.size < 1_000_000 /*= too big*/
				{
				text = GetFile(dir $ '/' $ item.name)
				if text isnt false
					{
					name = item.name
					path = .makePath(rootSrc, dir $ '/' $ item.name)
					importRec = Record(:name, :path, :text)
					.output(svcTable, importRec)
					}
				}
			else
				{
				Print('Skipped large file', dir, item.name)
				}
			}
		}
	makePath(src, filepath)
		{
		dir = Paths.ToStd(Paths.ParentOf(filepath))
		dir = dir.RightTrim('/')
		if dir is src
			return ''
		return '/' $ Paths.AbsToRel(src, dir)
		}
	output(svcTable, importRec)
		{
		svcTable.Output(importRec)
		++.n
		.totalSize += importRec.text.Size()
		}
	}
