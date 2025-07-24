// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// Interface to curl.exe, see: http://curl.haxx.se/docs/manpage.html
// INTERNAL - should only be used directly by Http, Https, and Ftp
class
	{
	// static method, does not require an instance
	// does not use MappedOptions
	Http(method, url, content = '', fromFile = '', toFile = '', header = #(),
		timeoutConnect = false, user = '', pass = '', cookies = '', limitRate = '')
		{
		Assert(url matches: '^https?://')
		Assert(content is "" or fromFile is "",
			'Curl Http should not have content AND fromFile')
		Assert((user is "") is (pass is ""),
			"Curl user and pass should be specified together")

		args = .buildArgs(method, url, :content, :fromFile, :toFile, :header,
			:timeoutConnect, :user, :pass, :cookies, :limitRate)

		result = Object(header: '', content: '')
		error = .runCommand(args)
			{ |p|
			if content isnt ''
				p.Write(content)
			p.CloseWrite()

			result.header = InetMesg.ReadHeader(p).Join('\n')
			if toFile is ''
				if false is result.content = p.Read()
					result.content = ""
			}
		if error isnt ''
			throw error
		return result
		}

	buildArgs(method, url, content = '', fromFile = '', toFile = '', header = #(),
		timeoutConnect = false, user = '', pass = '', cookies = '', pipe? = false,
		limitRate = '')
		{
		args = '"' $ url $ '"'
		if limitRate isnt ''
			args $= ' --limit-rate ' $ limitRate
		if toFile isnt ''
			args $= ' -o "' $ toFile $ '"'
		args $= Opt(' -u ', user $ (pass isnt '' ? ':' $ pass : ''))
		args $= ' -D -'
		args $= .addAdditionalArgs(url, timeoutConnect, header, cookies)
		args $= .handleMethodAndParams(method, :content, :fromFile, :pipe?)
		return args
		}

	addAdditionalArgs(url, timeoutConnect, header, cookies)
		{
		args = ''
		if url.Prefix?('https')
			args $= ' -k' // insecure
		if timeoutConnect isnt false
			args $= ' --connect-timeout ' $ timeoutConnect
		for h in header.Members().Sort!()
			args $= ' -H "' $ h.Tr('_', '-') $ ': ' $ header[h] $ '"'
		if cookies isnt ''
			args $= ' -b ' $ cookies $ ' -c ' $ cookies
		return args
		}
	handleMethodAndParams(method, content = '', fromFile = '', pipe? = false)
		{
		args = ''
		switch method
			{
		case 'GET':
			Assert(content is '' and fromFile is '', 'Http GET should not have body')
		case 'PUT':
			if pipe?
				args $= ' -g -T - -H "Transfer-Encoding: "'
			else if content isnt ""
				args $= ' -g -T - -H "Transfer-Encoding: " -H "Content-Length: ' $
					content.Size() $ '"'
			else
				{
				Assert(fromFile isnt '', 'Curl Http PUT requires content or fromFile')
				args $= ' -g -T "' $ fromFile $ '"'
				}
		case 'POST', 'PATCH':
			// use --data-binary so curl doesn't alter the data
			args $= ' --data-binary "@' $ (fromFile isnt '' ? fromFile : '-') $ '"'
			// suppress default Content-Type
			args $= ' -H Content-Type:'
		case 'POSTFILES':
			// nothing to do for this; files required in options for -F
			Assert(.options.Member?('files'), 'files required in options')
		case 'DELETE': // Needed for Amazon AWS
			args $= ' -X DELETE'
		// fake method, still act same as PUT
		// differentiate with PUT, which still forces non-empty upload file or content
		case 'EMPTYPUT':
			args $= ' -X PUT' // Needed for Amazon AWS
		case 'HEAD':
			args $= ' -I'
			}
		if method is 'PATCH'
			args $= ' -X PATCH'
		return args
		}

	HttpPiped(method, url, block, header = #(), timeoutConnect = false,
		user = '', pass = '', cookies = '')
		{
		Assert(url matches: '^https?://')
		Assert((user is "") is (pass is ""),
			"Curl user and pass should be specified together")

		args = .buildArgs(method, url, :header, :timeoutConnect,
			:user, :pass, :cookies, pipe?:)

		error = .runCommand(args, :block)
		if error isnt ''
			throw error
		}

	// ftp ---------------------------------------------------------------------

	// timeout is only used by Dir ???
	New(.protocol, server, user = '', pass = '', .timeout = 60,
		.timeoutConnect = 60, .options = #(), .exeSubFolder = '')
		{
		Assert(protocol is 'ftp' or protocol is 'sftp' or protocol is 'ftps' or
			protocol is 'https')
		.server = server.Prefix?(protocol $ '://')
			? server
			: protocol $ '://' $ server
		if .server isnt '' and not .server.Suffix?('/')
			.server $= '/'
		.userPass = .initUserPass(user, pass)
		}

	initUserPass(user, pass)
		{
		pwd = Opt(':', pass)
		pwd.Has?(`"`)
			pwd = pwd.Replace(`"`,`\\"`)
		userPass = not user.Blank?() and not pwd.Blank?() ? user $ pwd :
			pwd.Blank?() and not user.Blank?() ? user : ''
		return Opt(' -u ', '"', userPass, '"')
		}

	Get(remName, locName = '', retries = 0)
		{
		if locName is '' and String?(remName)
			locName = Paths.Basename(remName).Trim()
		if Object?(remName)
			{
			// "#1" - to use the remote name as local name
			locName $= locName is '' ? '#1' : '\#1'
			remName = '{' $ Paths.ToUnix(remName.Join(",")) $ '}'
			}
		args = '-o "' $ locName $ '" "' $ .server $ remName $ '"'
		return .runCommand(args, :retries)
		}

	GetMultiple(fileList, localPath, ftpPath)
		{
		Assert(.protocol.Has?('ftp'))
		if fileList.Empty?()
			return ''

		argsfile = GetAppTempFullFileName("curl")
		files = .buildGetScripts(fileList, localPath, ftpPath)

		PutFile(argsfile, files)

		.copyExistingFiles(fileList, localPath)

		result = .runCommand('-K ' $ argsfile) // read curl args from file
		DeleteFile(argsfile)
		return result
		}
	buildGetScripts(fileList, folderPath, receivingPath)
		{
		get = ''
		for filename in fileList
			{
			outputFileName = Paths.Basename(filename)
			get $= '-o "' $ Paths.ToUnix(folderPath) $ outputFileName $ '"\n' $
				'url = "' $ .server $ Opt(receivingPath, '/') $
					Url.EncodePreservePath(filename) $ '"\n'
			}
		return get
		}
	// TODO move this out of here (not part of curl interface)
	copyExistingFiles(fileList, folderPath)
		{
		// some edi companies (like BNSF) use the same file name over and over
		// so if something happens and the file does not get imported it just
		// gets overridden the next time the download runs and it gets lost
		// check for files with the same name in the directory before downloading
		// and if they exist, make a copy and add them back in to the import list
		for filename in fileList.Copy()
			{
			if FileExists?(folderPath $ filename)
				{
				newname = filename $
					Display(Timestamp()).Tr('#.')
				CopyFile(folderPath $ filename, folderPath $ newname, false)
				fileList.Add(newname)
				}
			}
		}

	Put(locName, remName = '', retries = 0)
		{
		Assert(.protocol.Has?('ftp'))
		if remName is '' and String?(locName)
			remName = Paths.Basename(locName).Trim()
		if Object?(locName)
			locName = '{' $ Paths.ToUnix(locName.Join(",")) $ '}'
		cmd = '-T ' $ Display(locName) $ ' ' $ Display(.server $ remName)
		return .runCommand(cmd, :retries)
		}

	Del(fileName, path = '', notFromRoot = false)
		{
		Assert(.protocol.Has?('ftp'))

		if .protocol is 'sftp'
			{
			optionalPath = Opt('/', path)
			result = .runCommand(' ' $ .server $ ' -Q "rm ' $ optionalPath $
				'/' $ fileName $ '"')
			return result.Prefix?('curl: ') ? result : ''
			}

		extraDir = notFromRoot ? "" : "/"
		if path isnt ''
			path = ' -Q "CWD ' $ extraDir $ path $ '"'
		result = .runCommand(' ' $ .server $ path $ ' -Q "DELE ' $ fileName $ '"')
		return result.Prefix?('curl: ') ? result : ''
		}

	// If a problem occurs, the process will stop immediately and not delete
	// the remaining files (this is due to the -Q option)
	// could change to use asterisk (*) in front of ftp command to make curl continue
	// but it won't report failure on any of the individual deletes (see curl help for -Q)
	DeleteMultiple(fileList, ftpPath, notFromRoot = false)
		{
		Assert(.protocol.Has?('ftp'))
		if fileList.Empty?()
			return ''

		deletefile = .makeTempFile("curl", .buildDeleteScripts(fileList, ftpPath,
			notFromRoot))

		ipaddress = .server.Suffix?('/') ? .server.BeforeLast('/') : .server
		result = .runCommand(' ' $ ipaddress $ ' -K ' $ deletefile)
		DeleteFile(deletefile)
		// curl returns the list of folder as default result
		return result.Prefix?('curl: ') ? result : ''
		}
	makeTempFile(prefix, text)
		{
		Retry()
			{
			deletefile = GetAppTempFullFileName(prefix)
			PutFile(deletefile, text)
			}
		return deletefile
		}
	buildDeleteScripts(fileList, receivingPath, notFromRoot)
		{
		delete = ''
		if .protocol isnt 'sftp'
			{
			extraDir = notFromRoot ? "" : "/"
			delete $= ' -Q "CWD ' $ extraDir $ receivingPath $ '"\n'
			}
		for filename in fileList
			{
			if .protocol is 'sftp'
				delete $= ' -Q "rm /' $ receivingPath $ '/' $ filename $ '"\n'
			else
				delete $= '-Q "DELE ' $ filename $ '"\n'
			}
		return delete
		}

	Ren(oldName, newName)
		{
		Assert(.protocol.Has?('ftp'))
		result = .runCommand(' ' $ .server $
			' -Q "RNFR ' $ oldName $ '" -Q "RNTO ' $ newName $ '"')
		return result.Prefix?('curl: ') ? result : ''
		}

	Dir(path = '*.*', details = false, caseSense = false, regExp = '',
		listOnly = true)
		{
		Assert(.protocol.Has?('ftp'))
		dirText = ""
		folder = path.BeforeLast('/')
		if folder isnt '' and not folder.Suffix?('/')
			folder $=  '/'
		cmd = .server $ folder $ ' -m ' $ .timeout $
			(details is false ? (listOnly is true ? ' -l' : '') : '')
		dirText = .runCommand(cmd)

		if dirText.Prefix?('curl: ')
			{
			// treat 2xx codes (other than 200) as empty success, not errors
			// (return empty file list)
			return dirText =~ `response: 2\d\d` ? Object() : false
			}
		return .dirList(dirText, FtpClient.BuildFMask(path, caseSense), regExp, details)
		}
	sizeCol: 4
	dirList(dirText, fmask, regExp, details)
		{
		list = Object()
		for line in dirText.Lines()
			{
			// BUG: this does not handle file names that contain spaces.
			// We probably should not be doing AfterLast if listOnly option was used.
			// Not sure if we can safely handle this when details are included unless
			// we can be guaranteed of the output format. There is likely some
			// differences between ftp servers
			name = line.Has?(' ') ? line.Trim().AfterLast(' ') : line
			if name !~ fmask or name !~ regExp
				continue
			if not details
				list.Add(name)
			else
				{
				ob = line.Split(' ').Remove('')
				if ob.Size() < .sizeCol + 1 or not ob[.sizeCol].Numeric?()
					continue
				list.Add(Object(:name, size: Number(ob[.sizeCol])))
				}
			}
		return list
		}

	DirMultiple(directories)
		{
		Assert(.protocol.Has?('ftp'))
		if directories.Empty?()
			return false

		argsfile = GetAppTempFullFileName("curl")
		PutFile(argsfile, .buildDirectoryScripts(directories, argsfile))

		dirText = .runCommand('-K ' $ argsfile) // read curl args from file
		if dirText.Prefix?('curl: ')
			{
			SuneidoLog('Curl.DirMultiple - ' $ dirText $ directories.Join(','))
			.cleanUpDirFiles(directories, argsfile)
			DeleteFile(argsfile)
			return false
			}

		dirLists = Object()
		for i, dir in directories
			{
			if false is listFile = GetFile(argsfile $ i)
				listFile = ''
			dirLists.Add(Object(:dir, list: listFile.Lines()))
			DeleteFile(argsfile $ i)
			}

		DeleteFile(argsfile)
		return dirLists
		}
	buildDirectoryScripts(directories, outputFilePrefix)
		{
		dirCmd = '-l\n'
		for i, dir in directories
			{
			dirCmd $= 'url = "' $ .server $ dir $ '"\n'
			dirCmd $= 'output = "' $ outputFilePrefix $ i $ '"\n'
			}
		return dirCmd
		}
	cleanUpDirFiles(directories, outputFile)
		{
		for i in .. directories.Size()
			DeleteFile(outputFile $ i)
		}

	// internal ----------------------------------------------------------------

	runCommand(args, block = false retries = 0)
		{
		cmd = .buildCommand(args)
		if cmd.Prefix?(`"false"`) // can not find curl.exe
			return 'missing curl.exe'
		cmd = .addExtraDebugging(cmd)
		for (i = 1; ; ++i)
			{
			result = block is false
				? .runPiped(cmd)
				: .runPipedWithBlock(cmd, block)

			result = result.Replace('(?q)curl: (56) OpenSSL SSL_read: ' $
				'error:0A000126:SSL routines::unexpected eof while reading, ' $
				'errno 0(?-q)\r?\n?', '') // ignore this error
			if not .retry?(result, i, retries)
				return .extractDebugging(result)
			RetrySleep(i, 200) /*= 200, 400, 800 - not too short, not too long */
			}
		}
	retry?(result, i, retries)
		{
		if result.Has?('SSL_ERROR_SYSCALL') or
			result.Has?('OpenSSL SSL_read: Connection was reset') or
			result.Has?('OpenSSL SSL_connect: Connection was reset')
			return i < Max(4, retries) /*= max retries for this */
		return result.Prefix?('curl: ') and i < retries
		}

	runPipedWithBlock(cmd, block, _curlDebugFile = false)
		{
		curlOutputFile = GetAppTempFullFileName('curl')
		cmd $= ' --stderr "' $ curlOutputFile $ '"'

		result = ''
		exitValue = 0
		Finally(
			{
			RunPiped(cmd, { |p|
				block(p)
				exitValue = p.ExitValue() })
			}, {
			if exitValue isnt 0 or curlDebugFile isnt false
				result = GetFile(curlOutputFile)

			LogErrors('Curl DeleteFile curlOutputFile',
				asErratic: #('Access is denied'))
				{
				DeleteFile(curlOutputFile)
				}
			})
		return result
		}

	timeoutConnect: 60
	userPass: ''
	options: ()
	buildCommand(args)
		{
		// WARNING: runCommand depends on this returning '"false"' if not found
		ct = ''
		if not args.Has?('--connect-timeout')
			ct = ' --connect-timeout ' $ .timeoutConnect
		return .app() $
			' -Y 1 -y 480' $ // abort if slower than 1 byte/sec for 8 min
			ct $
			.userPass $
			.curlOptions() $
			' ' $ args $
			' -s' $ // silent
			' -S' // show_errors
		}
	exeSubFolder: ''
	app()
		{
		// WARNING: runCommand depends on this returning '"false"' if not found
		app = Sys.Linux?()
			? 'curl'
			: .exeSubFolder isnt ''
				? ExeDir() $ '/' $ .exeSubFolder $ '/curl'
				: '"' $ ExternalApp("curl") $ '"'
		.checkVersion(app)
		return app
		}

	checkVersion(app)
		{
		if app.Prefix?(`"false"`) // can not find curl.exe
			return

		if Suneido.GetDefault('Curl_VersionChecked', false) is true
			return

		Suneido.Curl_VersionChecked = true
		if false isnt min = .minVersion()
			{
			cur = .version(app)
			if not Object?(cur)
				SuneidoLog('ERROR: (CAUGHT) Unexpected curl version: ' $ cur,
					params: [:app], caughtMsg: "One time check without throw", calls:)
			else if cur < min
				SuneidoLog('ERROR: (CAUGHT) The curl (' $ cur.Join('.') $
					') is lower than the minimum requirement (' $ min.Join('.') $ ')',
					params: [:app], caughtMsg: "One time check without throw", calls:)
			}
		}

	minVersion()
		{
		return OptContribution('MinCurlVersion', false)
		}

	version(app)
		{
		s = .runPiped(app $ ' -V')
		versionString = s.RemovePrefix('curl ').BeforeFirst(' ')
		if versionString.Suffix?('-DEV')
			return versionString
		return versionString.Split('.').Map(Number)
		}

	addExtraDebugging(cmd, _curlDebugFile = false)
		{
		if curlDebugFile is false
			return cmd

		Assert(curlDebugFile hasnt: ' ')
		debug_cmd = cmd $ ' -v '
		AddFile(curlDebugFile, '\r\n' $ Display(Timestamp()) $ ': ' $ debug_cmd $ '\r\n')
		return debug_cmd
		}
	runPiped(cmd) // overridden by tests
		{
		return RunPipedOutput(cmd)
		}
	extractDebugging(result, _curlDebugFile = false)
		{
		if curlDebugFile is false
			return result.Trim()
		AddFile(curlDebugFile, result $ '\r\n')
		returnVal = result.Extract('^curl: \([0-9]+\).*$')
		return String?(returnVal) ? returnVal : ''
		}

	// these are for external use
	// only used for Ftp, not Http
	MappedOptions()
		{
		options = Object(
			'append': '-a' // used by carslib Honda
			'insecure': '-k'
			'non_passive': '-P-'
			'ssl': '--ssl-reqd'
			'quote_user': 'USER'
			'quote_password': 'PASS'
			'disable_epsv': '--disable-epsv'
			'use_ascii': '--use-ascii'
			'key': '--key'
			'pubkey': '--pubkey'
			'cacert': '--cacert'
			'cert': '--cert'
			'files': '-F'
			'ignore_content_length': '--ignore-content-length'
			'key_pass_phrase': '--pass'
			)
		return options
		}

	curlOptions(options = false)
		{
		if options is false
			options = .options
		mappedOptions = .MappedOptions()
		str = quoteStr = ''
		for option in options.Members()
			{
			value = options[option]
			curlOption = mappedOptions[option]
			if not mappedOptions.Member?(option)
				ProgrammerError('Curl invalid option: ' $ option)
			else if value is ''
				ProgrammerError('Curl empty option value for ' $ option)
			else if option.Prefix?('quote_')
				{
				// need to add -Q "PASS" after -Q "USER"; otherwise, will get curl error
				if option is 'quote_password'
					continue
				quoteStr $= ' -Q "' $ curlOption $ ' ' $ value $ '"'
				if option is 'quote_user'
					quoteStr $= ' -Q "' $ mappedOptions.quote_password $ ' ' $
						options.quote_password $ '"'
				}
			else if option is 'files'
				{
				for val in value
					str $= ' ' $ curlOption $ ' "' $ val $ '"'
				}
			else
				str $= ' ' $ curlOption $ ' ' $ (value is true ? '' : value)
			}
		return quoteStr $ str
		}
	}