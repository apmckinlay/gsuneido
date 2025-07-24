// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
/*
README:
- This control is designed to be run ONCE
-- Running it multiple times can have unexpected side effects
--- For that reason, most of the vital methods are public and callable from the WorkSpace

- After the process completes, you can connect to this SVC Server via SVC > Settings
-- The "Server" value will be the IP address of the new SvcServer's device

- You will need to add your user from the client via: SvcUsers
-- You will need to use the password you set during the setup process for "serverPassword"
*/
Controller
	{
	Title: 'SVC Server Manager'
	Controls: #(Record,
		(Vert,
			(Heading3, 'Usage:'),
			(StaticWrap,
				'\tUse this control to convert this database\r\n' $
				'\tinto a fully functional Suneido Version\r\n\tControl (SVC)',
				xmin: 350)
			Skip,
			Skip, EtchedLine, Skip,
			(Heading3, 'Server Password:'),
			svcmng_password, svcmng_password_verify,

			Skip, EtchedLine, Skip,
			(Heading3, 'Create master tables for:'),
			svcmng_libraries,
			svcmng_books,

			Skip, EtchedLine, Skip,
			(Horz,
				(Heading3, 'Create custom SVC Library: '),
				(Skip, small:),
				svcmng_create_custom_library
				)
			svcmng_custom_library,

			Skip, EtchedLine, Skip,
			(Heading3, 'Additional Options:'),
			(Horz, Skip, svcmng_local_only),
			(Horz, Skip, svcmng_unuse_excess_tables),
			(Horz, Skip, svcmng_drop_excess_tables),
			(Horz, Skip, svcmng_output_service)
		))
	CallClass()
		{
		if not OkCancel(Object(this), .Title)
			return
		if .Service() // If the service exists, a restart is required
			.restartIDE()
		else
			{
			Unload()
			ResetCaches()
			}
		}

	// vvvvvvvvvvvvvvvvvvvvv GUI / Control related code vvvvvvvvvvvvvvvvvvvv
	New()
		{
		.Data.Set([svcmng_custom_library: #svclib])
		.svcmng_local_only = .FindControl(#svcmng_local_only)
		.svcmng_unuse_excess_tables = .FindControl(#svcmng_unuse_excess_tables)
		.svcmng_drop_excess_tables = .FindControl(#svcmng_drop_excess_tables)
		.svcmng_output_service = .FindControl(#svcmng_output_service)
		.svcmng_custom_library = .FindControl(#svcmng_custom_library)
		.svcmng_custom_library.SetReadOnly(true)
		.data = .Data.Get()
		}

	Record_NewValue(value, source /*unused*/)
		{
		if value is #svcmng_local_only
			{
			standalone = .data.svcmng_local_only
			if standalone is true
				{
				.svcmng_unuse_excess_tables.Set(false)
				.svcmng_drop_excess_tables.Set(false)
				.svcmng_output_service.Set(false)
				}
			.svcmng_unuse_excess_tables.SetReadOnly(standalone)
			.svcmng_drop_excess_tables.SetReadOnly(standalone)
			.svcmng_output_service.SetReadOnly(standalone)
			}
		if value is #svcmng_create_custom_library
			.svcmng_custom_library.
				SetReadOnly(.data.svcmng_create_custom_library isnt true)
		if value is #svcmng_drop_excess_tables
			{
			.svcmng_unuse_excess_tables.Set(.data.svcmng_drop_excess_tables)
			.svcmng_unuse_excess_tables.SetReadOnly(.data.svcmng_drop_excess_tables)
			}
		}

	OK()
		{
		convert = false
		_dir = .currentDir()
		if '' isnt msg = .valid?(.data)
			.AlertError(.Title, msg)
		else if convert = .verifyOptions(.data)
			try
				{
				.AlertInfo(.Title, 'Output from this process can be found at: ' $
					.filePath(_dir, .logFile))
				.convertIDE()
				}
			catch (e)
				.print(.Title $ ': FAILED to complete the conversion: ' $ e)
		return convert
		}

	currentDir()
		{ return GetCurrentDirectory() }

	valid?(data)
		{
		if data.svcmng_password is ''
			return 'Password is required'
		if data.svcmng_password isnt data.svcmng_password_verify
			return 'Password must match Verify'
		if data.svcmng_create_custom_library is true and data.svcmng_custom_library is ''
			return 'Custom SVC library is required'
		return ''
		}

	verifyOptions(data)
		{
		if data.svcmng_drop_excess_tables is true
			return YesNo(
				'WARNING: "Drop excess tables" will result in all unnecessary ' $
				'tables being dropped. This cannot be undone.\r\n\r\nContinue?',
				title: .Title, flags: MB.ICONWARNING)
		return true
		}

	logFile: 'svcServerConvert.log'
	print(msg, mode = 'a', _dir = false)
		{
		if dir is false
			dir = .currentDir()
		Print(msg)
		File(.filePath(dir, .logFile), mode)
			{ it.Writeline(msg) }
		}

	// vvvvvvvvvvvvvvvvvvvvvv Convert process begins vvvvvvvvvvvvvvvvvvvvvvv
	convertIDE()
		{
		.print(.Title $ ': Setting IDE to run as a SVC Server', mode: 'w')
		.createBackUp()
		.setSvcSettings()
		.InitializeMasterTables(
			.data.svcmng_libraries.Split(','),
			.data.svcmng_books.Split(','),
			overwrite:)
		if .data.svcmng_drop_excess_tables is true
			.dropExcessTables()
		else if .data.svcmng_unuse_excess_tables is true
			.unuseExcessTables()
		.SetPassword(.data.svcmng_password)
		if .data.svcmng_create_custom_library is true
			.CustomSVCLib(.data.svcmng_custom_library)
		.InitiateSvcServer(.data.svcmng_output_service is true)
		.print(.Title $ ': SUCCESSFULLY converted IDE to run as a SVC Server')
		}

	createBackUp()
		{
		backup = 'svcConvertBackup_' $ Display(Timestamp()).Tr('#.') $ '.su'
		.print(.Title $ ': Creating back up: ' $ backup)
		Database.Dump()
		CopyFile('database.su', backup, true)
		}

	setSvcSettings()
		{
		QueryDo('update svc_settings set svc_local? = true')
		PubSub.Publish('SvcSettings_ConnectionModified')
		}

	SetPassword(password)
		{
		Database('ensure svc_passwords (svc_password) key (svc_password)')
		svc_password = PassHash('', password)
		QueryEmpty?('svc_passwords')
			? QueryOutput('svc_passwords', [:svc_password])
			: QueryDo('update svc_passwords set svc_password = ' $ Display(svc_password))
		}

	// vvvvvvvvvv Initialize master tables based on local tables vvvvvvvvvvv
	InitializeMasterTables(libraries, books, overwrite = false)
		{
		for type, tables in Object(library: libraries, book: books)
			{
			if tables.Empty?()
				continue
			.print(.Title $ ': Initializing master tables for ' $ type $ ' tables')
			for table in tables
				try
					.initializeMasterTable(table, :overwrite)
				catch (e)
					.print('\t\tError: ' $ e)
			}
		}

	initializeMasterTable(table, overwrite)
		{
		if TableExists?(masterTable = table $ '_master')
			{
			if not overwrite and not YesNo('Master table already exists for: ' $ table $
				'\r\nDo you wish to overwrite this master table?', title: .Title)
				return
			else
				Database('drop ' $ masterTable)
			}
		.print('\tOutputting table: ' $ masterTable)
		svcTable = SvcTable(table, svcEnsure:)
		SvcCore.EnsureMaster(masterTable, svcTable.Type)
		QueryApply(svcTable.Query())
			{
			.initRecord(svcTable, name: svcTable.MakeName(it))
			}
		svcTable.SetMaxCommitted(Timestamp(), force:)
		}

	initRecord(svcTable, name)
		{
		if false is rec = svcTable.Get(name)
			rec = []
		rec.name = name
		rec.id = 'svc'
		rec.type = '+'
		rec.comment = 'initializing master table'
		rec.lib_committed = rec.GetDefault('lib_committed', Timestamp())
		QueryOutput(svcTable.Table() $ '_master', rec)
		}

	// vvvvvvvvvvvvvvvvvvvv Drop / Unuse "extra" tables vvvvvvvvvvvvvvvvvvvv
	dropExcessTables()
		{
		excludes = QueryList("tables where table.Suffix?('_master')", "table")
		if .data.svcmng_create_custom_library is true
			excludes.Add(.data.svcmng_custom_library)
		.print(.Title $ ': Dropping excess tables')
		DropTables(excludes, dropLibraries:, dropBooks:, quiet:)
		}

	unuseExcessTables()
		{
		libraries = Libraries().Filter({ it isnt #stdlib and it isnt #configlib })
		if .data.svcmng_create_custom_library is true
			libraries.Remove(.data.svcmng_custom_library)
		if libraries.Empty?()
			return
		.print(.Title $ ': Unusing excess libraries')
		for lib in libraries
			ServerEval(#Unuse, lib)
		}

	// vvvvvvvvvv Output the custom SVC library with basic records vvvvvvvvv
	customLibRecs: #(
		SvcCommitHooks: #(
			`// USAGE: Use this record to automate behavior on record commit`,
			`// 	table: 	master table changes were committed to`,
			`// 	type: 	will be Put (update or add), or Remove (deleted)`,
			`function (table /*unused*/, type /*unused*/)`,
			`	{`,
			`	}`),
		RackRoutes: #(
			`// USAGE: Use this record to define rack routes and their behaviors`
			`#(`,
			`	// #('[Get | Post]', '/<access route>', ` $
				`['callable name'] | [function (env) { }])`,
			`)`))
	CustomSVCLib(lib)
		{
		if lib is ''
			return
		if not TableExists?(lib)
			{
			.print(.Title $ ': Creating custom SVC library: ' $ lib)
			LibTreeModel.Create(lib)
			}
		else
			.print(.Title $ ': Updating custom SVC library: ' $ lib)
		prefix = lib[0].Upper() $ lib[1 ..] $ '_'
		for m, v in .customLibRecs
			.outputRecord(prefix $ m, v.Join('\r\n'), lib)
		ServerEval(#Use, lib)
		}

	outputRecord(name, text, lib)
		{
		try
			{
			rec = [parent: 0, :name, :text, group: false]
			model = TreeModel(lib)
			model.EnsureUnique(rec)
			model.NewItem(rec)
			.print('\tOutput record: ' $ name)
			}
		catch (e)
			.print('\tERROR: encountered when outputting: ' $ name $ ', ' $ e)
		}

	// vvvvvvvvvvvvvvvv Adjust files to run as a SVC server vvvvvvvvvvvvvvvv
	// If running as a service, the IDE will require a restart via
	// svcRestart.bat (.restartFile). This will launch the service and
	// load the database
	InitiateSvcServer(runAsService = false, _dir = false)
		{
		if dir is false
			dir = .currentDir()
		.OutputFiles(runAsService, dir)
		.RemoveService()
		.ExtraProcesses(start: not runAsService)
		if runAsService
			.OutputService(dir)
		}

	goFile: 'svcServer.go'
	startFile: 'startSvcServer.bat'
	OutputFiles(runAsService = false, _dir = false)
		{
		if dir is false
			dir = .currentDir()
		.putFile(dir, .goFile, .serverGoText(runAsService))
		.putFile(dir, .startFile, .serverBatText(dir))
		}

	putFile(dir, file, text)
		{
		PutFile(.filePath(dir, file), text)
		.print(.Title $ ': ouput file: ' $ file $ ', to: ' $ dir)
		}

	filePath(dir, file)
		{ return Paths.ToLocal(Paths.Combine(dir, file)) }

	serverGoText(runAsService)
		{
		text = Libraries().Map({ 'Use(#' $ it $ ')' }).
			Add('LibraryTags.Reset()',
				'Suneido.User = Suneido.User_Loaded = #none',
				'Suneido.user_roles = #(none)',
				'ExtraProcesses.Start()')
		if not runAsService
			text.Add('CSDevServerWindow()')
		return text.Join('\r\n')
		}

	port: 3147
	serverBatText(dir)
		{ return .filePath(dir, 'gsport.exe') $ ' -s -p ' $ .port $ ' ' $ .goFile }

	ExtraProcesses(start = false)
		{
		ExtraProcesses.Ensure()
		ExtraProcesses.EnsureProcess(#scheduler,
			'Scheduler({ GetContributions(#SystemTasks) })',
			:start)
		ExtraProcesses.EnsureProcess(#svc, 'SvcServer()', :start)
		ExtraProcesses.EnsureProcess(#http, 'RunHttpServer()', :start)
		}

	service: 		svc_server
	ideFile: 		'Launch SVC IDE.bat'
	restartFile: 	'svcRestart.bat'
	OutputService(_dir = false)
		{
		if dir is false
			dir = .currentDir()
		.putFile(dir, .ideFile, .filePath(dir, 'gsuneido.exe') $ ' -c -p ' $ .port)
		cmd = 'sc create "' $ .service $
			'" displayname= "Suneido Version Control Server"' $
			' type= own binpath= "' $ .filePath(dir, 'gsport.exe') $
			' -p ' $ .port $ ' -s ' $ .goFile $ '" start= delayed-auto'
		if restart = .runCommand('Outputting service', cmd) is 0
			.putFile(dir, .restartFile, .restartBatText(dir))
		return restart
		}

	runCommand(msg, command, skipPrint = false)
		{
		if not skipPrint
			.print(.Title $ ': ' $ msg $ ':')
		result = RunPipedOutput.WithExitValue(command)
		if not skipPrint
			result.output.Lines().RemoveIf({ it.Blank?() }).Each({ .print('\t' $ it) })
		return result.exitValue
		}

	restartBatText(dir)
		{
		serverExe = .filePath(dir, 'gsport.exe')
		return ['sleep 10',
			serverExe $ ' -load',
			'sc start ' $ .service,
			'sleep 10',
			.filePath(dir, Display(.ideFile))].Join('\r\n')
		}

	RemoveService()
		{
		if not .Service()
			return // Query failed to find service, nothing to remove
		.runCommand('Stopping old service', 'sc stop ' $ .service)
		.runCommand('Deleting old service', 'sc delete ' $ .service)
		}

	Service(print = false)
		{
		return 0 is .runCommand('Service Status',
			'sc query ' $ .service,
			skipPrint: not print)
		}

	restartIDE()
		{
		.AlertInfo(.Title, 'A restart is required to finalize the conversion.\r\n' $
			'The IDE will re-launch automatically')
		.print(.Title $ ': Dumping database for server restart')
		Database.Dump()
		Spawn(P.NOWAIT, .restartFile)
		Shutdown(true)
		}
	}
