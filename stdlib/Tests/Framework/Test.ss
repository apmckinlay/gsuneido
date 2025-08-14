// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
print(@args)
{
//Print(@args)
}
	CallClass(quiet = false, timeEachMethod? = false)
		{
		name = Display(this).BeforeFirst(' ')
		observer = .RunTest(name, :quiet, :timeEachMethod?)
		return observer.Result
		}

	RunTest(name, observer = false, quiet = false, timeEachMethod? = false)
		{
		prevConditions = .beforeRun()
		if observer is false
			observer = new TestObserverString(:quiet)
		TestRunner.Run1(name, observer, check_tables:, :timeEachMethod?)
		// Call Test.AfterRun instead of .AfterRun because the current this class' base
		// class may not exist after running the test
		Test.AfterRun(prevConditions)
		return observer
		}

	beforeRun()
		{
		prevConditions = Object()
		Plugins().ForeachContribution('TestRunner', 'beforeRun')
			{
			prevConditions.MergeNew((it.func)())
			}
		return prevConditions
		}

	AfterRun(prevConditions)
		{
		Plugins().ForeachContribution('TestRunner', 'afterRun')
			{
			(it.func)(prevConditions)
			}
		}

	Debug()
		{
		_debugTest = true
		if Type(this) is "Class"
			return (new this).Debug()

		.debug_block()
			{
			.Foreach_test_method()
				{ |member|
				this[member]()
				}
			}
		}

	Foreach_test_method(block)
		{
		base = .Base()
		privatePrefix = Display(base).BeforeFirst(' ') $ '_'
		for member in .getbaseMembers(base).Sort!()
			.debug_one(member, privatePrefix, base, block)
		}

	getbaseMembers(base)
		{
		members = Object()
		do
			{
			members.MergeUnion(base.Members())
			base = base.Base()
			}
			while base isnt Test

		return members
		}

	DebugOne(test_method)
		{
		_debugTest = true
		if Type(this) is "Class"
			return (new this).DebugOne(test_method)

		.debug_block()
			{
			base = .Base()
			privatePrefix = Display(base).BeforeFirst(' ') $ '_'
			.debug_one(test_method, privatePrefix, base, {|member| this[member]() })
			}
		}

	debug_block(block)
		{
		Suneido.SkipRetrySleep = teardown? = passed? = true
		try
			{
			.Setup()
			block()
			}
		catch(e)
			{
			teardown? = .failureHandler(e)
			passed? = false
			}
		Suneido.Delete(#SkipRetrySleep)
		if teardown?
			.Teardown()
		return passed?
		}

	failureHandler(error)
		{
		teardown? = false
		if Sys.Win32?()
			Debugger.Window(0, error, calls: .callsBeforeDebug(error.Callstack()),
				onDestroy: .teardownOnDestroy)
		else
			{
			Print(ERROR: error)
			Print(FormatCallStack(.callsBeforeDebug(error.Callstack())))
			teardown? = true
			}
		return teardown?
		}
	callsBeforeDebug(stack)
		{
		calls = Object()
		for err in stack
			{
			if String(err.fn) is "Test.Debug /* stdlib block */"
				break
			calls.Add(err)
			}
		return calls
		}

	teardownOnDestroy()
		{
		// need to delete Debugger, so Debugger.Window does not destroy existing window
		// which caused onDestroy being called twice
		Suneido.Delete(#Debugger)
		.Teardown()
		}

	runningTestMethod: false
	debug_one(member, privatePrefix, base, block)
		{
		if member.Prefix?("Test") and
			not member.Prefix?(privatePrefix) and
			Function?(base[member])
			Finally(
				{
				.runningTestMethod = true
				block(member)
				},
				{
				.runningTestMethod = false
				.TeardownAfterEachMethod()
				})
		}

	New()
		{
		.teardowns = Object()
		.teardownsAfterEachMethod = Object()
		.suneidolog_exists? = TableExists?('suneidolog')
		ServerSuneido.Set('TestRunningLogs', Object())
		}

	Setup()
		{ }

	TempName()
		{ return .TempTableName().Capitalize() }

	TempTableName()
		{ return 'tests' $ Display(Timestamp()).Tr('#.', '_') }

	MakeTable(@args)
		{
		.AddTeardown(.teardown_tables)
		table = .TempTableName()
		Database("create " $ table $ " " $ args[0])
		for r in args.Delete(0)
			QueryOutput(table, r)
		return table
		}
	teardown_tables(_testname = "")
		{
		libs = Libraries()
		resetCache? = false
		QueryApply('tables
			where table > "tests_20" and table < "tests_2099"
			sort reverse table') // reverse to handle foreign keys
			{ |x|
			if libs.Has?(x.table)
				{
				resetCache? = true
				ServerEval('Unuse', x.table)
.print(testname, "teardown_tables")
				}
			try
				Database("drop " $ x.table)
			catch (e)
				Print(e)
			}
		if resetCache?
			SvcDisabledLibraries.ResetCache()
		}

	MakeView(definition)
		{
		.AddTeardown(.teardown_views)
		table = .TempTableName()
		Database("view " $ table $ " = " $ definition)
		return table
		}
	teardown_views()
		{
		QueryApply('views
			where view_name > "tests_20" and view_name < "tests_2099"')
			{ |x|
			Database("destroy " $ x.view_name)
			}
		}

	testlib: 'Test_lib'
	TestLibName()
		{ return .testlib }

	EnsureLibrary(lib, _testname = "")
		{
		if not TableExists?(lib)
			.makeLibrary(lib)
		if ServerEval('Use', lib)
			{
.print(testname, "EnsureLibrary 1", lib)
//StackTrace()
			Unload()
			}
		else if lib is .testlib and Libraries().Last() isnt lib
			{
.print(testname, "EnsureLibrary 2", lib)
			ServerEval('Unuse', lib)
			ServerEval('Use', lib)
			Unload()
			}
		}

	// example: .MakeLibraryRecord([name: "Class_Name", text: `class { }`])
	MakeLibraryRecord(@records)
		{
		table = .testlib
		if records.Member?('table')
			{
			// TODO: Teardown does not clean up records if the table is a standard lib
			table = records.table
			records.Delete('table')
			.AddTeardown(.teardown_tables)
			}

		.EnsureLibrary(table)
		for x in records
			{
			QueryDo('delete ' $ table $ ' where name is ' $ Display(x.name))
			OutputLibraryRecord(table, x)
			.checkLibraryName(x.name)
			}
		.AddTeardown(.teardown_test_lib)
		return table
		}
	MakeLibrary(@records)
		{
		.AddTeardown(.teardown_tables)
		return .makeLibrary(.TempTableName(), records)
		}
	makeLibrary(library, records = #())
		{
		LibTreeModel.Create(library)
		for x in records
			{
			OutputLibraryRecord(library, x)
			.checkLibraryName(x.name)
			}
		return library
		}
	checkLibraryName(name)
		{
//		if name.GlobalName?() and
//			not name.Has?(Date().Format("yyyyMMdd")) and
//			not Uninit?(name) and
//			not Suneido.GetInit(#MakeLibrary, Object).Member?(name)
//			{
//			Print("should not use MakeLibrary to override existing name: " $
//				name, "in e.g.", Display(this))
//			Suneido.MakeLibrary[name] = true
//			}
		}

	teardown_test_lib()
		{
		if not TableExists?(.testlib)
			return
		QueryApply(.testlib)
			{|x|
			Unload(x.name)
			ServerEval('Unload', x.name)
			}
		QueryDo('delete ' $ .testlib)
		}

	MakeBook()
		{
		.AddTeardown(.teardown_tables)
		BookModel.Create(book = .TempTableName())
		return book
		}

	MakeBookRecord(book, text, path = '', extrafields = #())
		{
		.AddTeardown({ .teardownBook(book) })
		num = QueryMax(book, 'num') + 1
		name = .TempName()
		rec = [:num, :path, :name, :text]
		rec.MergeNew(extrafields)
		QueryOutput(book, rec)
		return rec
		}

	teardownBook(book)
		{
		// If MakeBook is called before MakeBookRecord, the table will be dropped
		// prior to this teardown (rendering the below unnecessary)
		if TableExists?(book)
			QueryDo('delete ' $ book $ ' where name > "Tests_" and name < "Tests_z"')
		}

	MakeRecordFromDB(rec) // so .New?() returns false
		{
		return QueryFirst('columns sort column').Delete(all:).Merge(rec)
		}

	MakeFile(content = false)
		{
		.AddTeardown(.teardown_files)
		filename = .TempTableName()
		if content isnt false
			.PutFile(filename, content)
		return filename
		}

	// Do not use, unless the test calling this is cleaning the file up
	PutFile(name, content)
		{
		Retry(maxRetries: 3, minDelayMs: 100)
			{
			PutFile(name, content)
			}
		}

	teardown_files()
		{
		for file in Dir('./tests_20*')
			DeleteFile(file)
		}

	MakeDir(parent = '')
		{
		.AddTeardown(.teardown_dir)
		dirName = Opt(parent, '/') $ .TempTableName()
		EnsureDir(dirName)
		return dirName
		}

	teardown_dir()
		{
		for folder in Dir('./tests_20*').Filter({ it.Suffix?('/') })
			{
			result = RetryBool(maxretries: 5, min: 50) { DeleteDir(folder) }
			if result isnt true
				throw result
			}
		}

	// type: Plugin_FieldTypes
	MakeCustomField(tableName, type, prompt = false, options = #())
		{
		destroyCustomizable? = not TableExists?('customizable')
		c = Customizable(tableName)
		if prompt is false
			prompt = .TempName()

		unuseConfiglib? = not Libraries().Has?('configlib')
		destroyConfiglib? = not TableExists?('configlib')

		field = c.CreateField(prompt, type, SelectFields(), :options)

		.AddTeardown({
			Database('alter ' $ tableName $ ' drop (' $ field $ ')')
			if unuseConfiglib?
				{
.print(_testname, "MakeCustomField")
				ServerEval('Unuse', 'configlib')
				Unload()
				}
			if destroyConfiglib?
				Database('destroy configlib')
			else
				{
				QueryDo('delete configlib where name in (' $
					Display('Access_' $ field) $ ', ' $ Display('Field_' $ field) $ ')')
				Unload('Field_' $ field)
				}
			if destroyCustomizable?
				Database('destroy customizable')
			CustomizableMap.ResetServerCache()
			QueryColumns.ResetCache()
			})
		return [:field, :prompt]
		}

	MakeCustomizeField(tableName, field, formula = '', key = false, selectFields = false,
		extrafields = #())
		{
		if not tableName.Lower().Prefix?('tests_') and
			not field.Lower().Prefix?('tests_') and
			not Customizable.CustomField?(field)
			throw "MakeCustomizeField does not currently handle non-custom fields"

		if not TableExists?('customizable_fields')
			{
			Customizable.EnsureTable()
			.AddTeardown(.tearDownCustomizableFieldsTable)
			}
		.AddTeardown(.tearDownCustomizeFields)
		if selectFields is false
			selectFields = QueryColumns(tableName)
		fn = CustomizeField.TranslateFormula(
			SelectFields(selectFields), formula, field)
		if key is false
			key = tableName
		else
			Customizable.ResetCustomizedCache(key)
		QueryOutput('customizable_fields', [custfield_num: .TempTableName(),
			custfield_name: key, custfield_field: field
			custfield_formula: formula,
			custfield_formula_code: fn.formulaCode
			custfield_formula_fields: fn.fields].Merge(extrafields))
		}

	MakeDatadict(@args)
		{
		fieldName = args.Extract('fieldName', .TempName().Lower())
		baseClass = args.Extract('baseClass', 'Field_string')
		control = args.Extract('control', 'Field')

		fieldArgs = Object()
		for m, v in args
			fieldArgs.Add(m $ `: ` $ (String?(v) ? Display(v) : v))
		joinStr = '\r\n\t\t\t\t'
		.MakeLibraryRecord([name: `Field_` $ fieldName,
			text: baseClass $ `
				{` $
				Opt(joinStr, fieldArgs.Join(joinStr)) $ `
				Control: (` $ control $ `)
				Format: (Text)
				}`])
		return fieldName
		}

	MakeIdField()
		{
		field = .TempTableName()
		fieldPrefix = field.Tr('_')
		num = fieldPrefix $ '_num'
		name = fieldPrefix $ '_name'
		abbrev = fieldPrefix $ '_abbrev'
		table = .MakeTable('(' $ Object(num, name, abbrev).Join(',') $ ')
			key (' $ num $ ')
			key (' $ name $ ')
			index unique (' $ abbrev $ ')')
		.MakeLibraryRecord([name: "Field_" $ num,
			text: `Field_num
				{
				Prompt: "` $ field $ `"
				Control: (Id "` $ table $ `"
					columns: (` $ name $ ', ' $ abbrev $ `),
					field: ` $ num $ `)
				Format: (Id query: '` $ table $ `' numField: ` $ num $ `)
				}`])
		leftjoin = ' leftjoin by(' $ num $ ') (' $ table $
			' project ' $ num $ ', ' $ name $ ', ' $ abbrev $ ')'
		return Object(:table, :num, :name, :abbrev, :leftjoin)
		}

	MakeNextNum(field = 'num', num = 1)
		{
		.AddTeardown(.teardown_tables)
		GetNextNum.Create(table = .TempTableName(), field, num)
		return table
		}

	MakeEventConditionActions(name, conditions, actions)
		{
		.AddTeardown(.teardown_eca)
		num = .TempTableName() // should be suitable for test key as well
		QueryOutput('event_condition_actions', Record(
			eca_num: num,
			eca_event: name,
			eca_conditions: conditions,
			eca_actions: actions))
		}

	teardown_eca()
		{
		QueryDo('delete event_condition_actions
			where eca_num > "tests_20" and eca_num < "tests_2099"')
		}

	WatchTable(table)
		{
		watchMember = 'TestWatchTable_' $ table $ .TempName()
		ServerSuneido.Set(watchMember, Object())
		.AddTeardown({.teardownWatchTable(table, watchMember) })
		trigger = "Trigger_" $ table
		callSuper = .getCallSuperAndTrackWatch(table, trigger, watchMember)
		code = .watchTableOverrideCode(watchMember, callSuper)
		ServerEval("LibraryOverride", .testlib, trigger, code)
		return watchMember
		}
	teardownWatchTable(table, watchMember)
		{
		recs = ServerSuneido.Get(watchMember, Object())
		ServerSuneido.DeleteMember(watchMember)
		if TableExists?(table)
			for rec in recs
				QueryDelete(table, rec)
		ob = Suneido.GetDefault('WatchTable_' $ table, Record())
		newmems = ob.members.Copy()
		newmems.Remove(watchMember)
		if newmems.Empty?()
			Suneido.Delete('WatchTable_' $ table)
		else
			Suneido['WatchTable_' $ table] = Object(members: newmems,
				callSuper: ob.callSuper)
		}
	watchTableOverrideCode(watchMember, callSuper)
		{
		return `function(t, oldrec, newrec)
			{
			member = "` $ watchMember  $ `"
			` $ callSuper $ `
			if not Suneido.Member?(member) // for teardown
				return
			if oldrec isnt false
				Suneido[member].Remove(oldrec)
			if newrec isnt false
				Suneido[member].Add(newrec)
			}`
		}
	getCallSuperAndTrackWatch(table, trigger, watchMember)
		{
		callSuper = ''
		if not Suneido.Member?('WatchTable_' $ table)
			{
			if not Uninit?(trigger)
				callSuper = "_" $ trigger $ '(t, oldrec, newrec)'
			Suneido['WatchTable_' $ table] = Object(:callSuper,
				members: Object(watchMember))
			}
		else
			{
			ob = Suneido.GetDefault('WatchTable_' $ table, Record())
			callSuper = ob.callSuper
			newmembers = ob.members.Copy().Add(watchMember)
			Suneido['WatchTable_' $ table] = Object(:callSuper, members: newmembers)
			}
		return callSuper
		}

	GetWatchTable(watchMember)
		{
		return ServerSuneido.Get(watchMember, Object())
		}

	SpyOn(target)
		{
		if .runningTestMethod isnt true
			throw "Usage: .SpyOn must be called within test methods"
		.AddTeardownAfterEachMethod(.teardownSpys)
		.spys.Add(spy = Spy(target))
		return spy
		}
	teardownSpys()
		{
		.spys.Each()
			{ |spy|
			if spy.NotEmpty?()
				spy.Close()
			}
		.spys = Object()
		}
	getter_spys()
		{
		return .spys = Object() // once only
		}

	TearDownIfTablesNotExist(@tables)
		{
		tables = tables.Filter({ not TableExists?(it) }).Instantiate()
		.AddTeardown({ .dropTables(tables) })
		}

	dropTables(tables)
		{
		tables.Filter(TableExists?).Each({ Database('drop ' $ it) })
		}

	tearDownCustomizableFieldsTable()
		{
		Database('destroy customizable_fields')
		Customizable.ResetAllCache()
		}

	tearDownCustomizeFields()
		{
		QueryDo('delete customizable_fields
			where custfield_num > "tests_20" and custfield_num < "tests_2099"')
		Customizable.ResetAllCache()
		}

	AddTeardownAfterEachMethod(func)
		{
		.teardownsAfterEachMethod.AddUnique(func)
		}
	TeardownAfterEachMethod()
		{
		try
			{
			for (i = .teardownsAfterEachMethod.Size() - 1; i >= 0; --i)
				(.teardownsAfterEachMethod[i])()
			}
		.teardownsAfterEachMethod = Object()
		}
	AddTeardown(@args)
		{
		fn = args.Size() is 1
			? args[0]
			: Curry(@args) // need Curry for AddUnique since blocks are always unique
		.teardowns.AddUnique(fn)
		}
	Teardown()
		{
		errors = Object()
		tries = 0
		maxtries = 5
		allteardowns = .teardowns.Copy()
		while (tries < maxtries and not allteardowns.Empty?())
			{
			errors = Object()
			for (i = allteardowns.Size() - 1; i >= 0; --i)
				{
				try
					{
					(allteardowns[i])()
					allteardowns.Delete(i)
					}
				catch (err, '*foreign key')
					{
					if tries is maxtries - 1 // last try, record errors for user
						errors.Add(err)
					}
				}
			++tries
			}
		.cleanUp()
		ServerEval("LibraryOverrideClear")
		if not allteardowns.Empty?()
			throw 'did not run all teardowns (see errors object) : ' $
				(errors.Join('; ')[..100])
		}
	cleanUp()
		{
		if not .suneidolog_exists?
			{
			if TableExists?('suneidolog')
				Database('destroy suneidolog')
			}
		else
			{
			for l in ServerSuneido.Get('TestRunningLogs', #())
				{
				QueryApply1('suneidolog', sulog_timestamp: l)
					{ |rec|
					if rec.sulog_message.Prefix?('ERROR') and
						not rec.sulog_message.Has?('CAUGHT') and
						not Database.SessionId().Has?("continuous_tests")
						Print(rec)
					rec.Delete()
					}
				}
			ServerSuneido.DeleteMember('TestRunningLogs')
			}
		}
	}