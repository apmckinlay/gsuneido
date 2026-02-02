// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// used by TestRunnerGui
	Run1(name, observer = false, check_tables = false, timeEachMethod? = false)
		{
		Test.EnsureLibrary(Test.TestLibName())
		instance = Construct(name)
		if observer is false
			observer = TestObserverPrint()
		.RunTests()
			{
			.wrap(observer)
				{
				.run(name, instance, observer, check_tables, :timeEachMethod?)
				}
			}
		}

	// used by TestRunnerGui
	Run1OnServer(name, observer, check_tables = false, lib = "")
		{
		// cant pack the observer, so need to pass the existing data through as object
		// code assumes you are using TestObserverOnServer
		Assert(observer.Base?(TestObserverOnServer))
		result = ServerEval('TestRunner.Run1FromServer', lib, name,
			observer.Values(), check_tables)
		observer.SetValues(result)
		}

	// Designed to be called from Run1OnServer
	Run1FromServer(lib, name, observerData, check_tables)
		{
		observer = TestObserverOnServer()
		observer.SetValues(observerData)
		try
			{
			x = Query1(lib, :name, group: -1)
			if not CodeTags.Matches(x.text)
				{
				observer.BeforeTest(name)
				observer.Warning(name, 'Skipped because of tags: ' $
					CodeTags.ExtractTags(x.text).Join(' '))
				observer.AfterTest(name, 0, 0, 0)
				}
			else
				.Run1(name, observer, check_tables)
			}
		catch (err)
			{
			observer.BeforeTest(name)
			observer.Error(name, err)
			observer.AfterTest(name, 0, 0, 0)
			}
		return observer.Values()
		}

	ensureSetup(libs)
		{
		Test.EnsureLibrary(Test.TestLibName())
		Suneido.Delete(#Config)
		SystemSummary() // can be slow, so ensure it's cached
		return libs is false
			? Libraries()
			: libs
		}

	// Designed to be called from ContinuousTests
	RunOnServer(libs)
		{
		libs = .ensureSetup(libs)
		serverObserver = TestObserverForServerTests()
		serverObserver.SetValues(ServerEval('TestRunner.RunAllFromServer', libs))

		return serverObserver.DisplayValue('errors') $
			serverObserver.DisplayValue('warnings') $
			serverObserver.FinalResults()
		}

	RunAllFromServer(libs)
		{
		observer = TestObserverForServerTests()
		tests = Object()
		for lib in libs
			.ForeachTest(lib)
				{ tests.Add(it) }
		tests.Shuffle!()
		.RunTests({ .runList(tests, observer) })
		return observer.Values()
		}

	Run(observer = false, quit_on_failure = false, libs = false, inorder = false)
		{
		libs = .ensureSetup(libs)

		if observer is false
			observer = TestObserverPrint(quiet:)

		tests = Object()
		for lib in libs
			.ForeachTest(lib)
				{ tests.Add(it) }
		if not inorder
			tests.Shuffle!()
		.RunTests({ return .runList(tests, observer, :quit_on_failure) })
		}

	RunTests(block)
		{
		if not Suneido.GetDefault(#TestRunner, false)
			{
			.check_demodata()
			Finally(
				{
				Suneido.TestRunner = true
				ServerSuneido.Set(#TestRunner, true)
				return block()
				},
				{
				Suneido.Delete(#TestRunner)
				ServerSuneido.DeleteMember(#TestRunner)
				})
			}
		else
			return block()
		}
	check_demodata()
		{
		if Libraries().Remove("Test_lib") is #(stdlib)
			return
		// demodata is output by Create_DemoData
		if TableExists?('demodata') and
			false isnt (x = Query1('demodata')) and
			x.libraries isnt Libraries().Remove("Test_lib")
			throw "Libraries changed since demo data was created (" $
				Display(x.when) $ ")\n" $ Display(x.libraries) $ "\n" $
				Display(Libraries())
		}

	ForeachTest(lib, block) // also used by TestRunnerGui
		{
		QueryApply(lib $ " where name =~ '^[A-Z].*Test$' and group is -1 sort name")
			{|x|
			if not CodeTags.Matches(x.text) or CheckLibrary.BuiltDate_skip?(x.text)
				continue
			try
				{
				val = x.name.Eval() // needs Eval to handle Name?_Test
				if Class?(val) and val.Base?(Test)
					block(x.name, :x)
				}
			catch (e)
				Print(lib $ ':' $ x.name $ " -", e)
			}
		}

	RunningTests?()
		{
		Suneido.GetDefault("TestRunner", false) is true
		}

	wrap(observer, block)
		{
		m = .measure(block)
		return observer.After(m.time, dbgrowth: m.dbgrowth, memory: m.memory)
		}

	runList(tests, observer, quit_on_failure = false)
		{
		return .wrap(observer)
			{
			before = SystemChanges.GetState()
			for name in tests
				{
				.run(name, Construct(name), observer)
				if quit_on_failure and observer.HasError?()
					break
				}
			if "" isnt s = SystemChanges.CompareState(before)
				observer.Error('', s)
			if not observer.HasError?()
				.OutputSuccess()
			}
		}

	testRunTable: testrunner_success
	OutputSuccess(date_time = false)
		{
		if date_time is false
			date_time = Date()
		Database('ensure ' $ .testRunTable $ '(date_time) key()')
		QueryEnsure(.testRunTable, [:date_time])
		}

	RequireRun()
		{
		if not .RunningTests?()
			try QueryDo('delete ' $ .testRunTable)
		}

	LastSuccess()
		{
		x = TableExists?(.testRunTable)
			? Query1(.testRunTable)
			: false
		return x is false
			? Date.Begin()
			: x.date_time
		}

	run(name, instance, observer, check_tables = false, timeEachMethod? = false)
		{
		observer.BeforeTest(name)

		transBefore = Database.Transactions().Size()
		td = Database.TempDest()
		nc = Database.Cursors()
		before = check_tables ? SystemChanges.GetState() : false

		m = .measure()
			{
			.run1(instance, observer, timeEachMethod?)
			}

		maxTime = 10
		if m.time > maxTime
			{
			secondsPrecision = 3
			observer.Warning("",
				"SLOWTEST: Test took longer than " $ maxTime $ " seconds" $
					" (" $ m.time.RoundToPrecision(secondsPrecision) $ ")")
			}

		.checkTransactions(transBefore, observer)
		.checkTables(check_tables, before, observer)
		.checkTempDest(td, observer)
		.checkCursors(nc, observer)

		observer.AfterTest(name, m.time, dbgrowth: m.dbgrowth, memory: m.memory)
		}

	checkTransactions(transBefore, observer)
		{
		retries = 3
		minSleep = 30
		if false is RetryBool(retries, minSleep,
			{ Database.Transactions().Size() <= transBefore })
			.error(observer, "", "didn't complete all its transactions")
		}

	checkTables(check_tables, before, observer)
		{
		if not check_tables
			return

		if "" isnt s = SystemChanges.CompareState(before)
			.error(observer, "", s)
		}

	checkTempDest(td, observer)
		{
		curTempDest = Database.TempDest()
		if curTempDest > td
			.error(observer, "", "didn't close all queries (tempdest grew by " $
				(curTempDest - td) $ ")")
		}

	checkCursors(nc, observer)
		{
		if ((nc = Database.Cursors() - nc) > 0)
			.error(observer, "", "didn't close " $ nc $ " cursor" $ (nc > 1 ? "s" : ""))
		}

	error(observer, where, msg)
		{
		observer.Error(where, msg)
		}
	run1(instance, observer, timeEachMethod?)
		{
		if .runMethod(observer, instance, 'Setup', true, timeEachMethod?)
			instance.Foreach_test_method
				{ |method|
				.runMethod(observer, instance, method, true, timeEachMethod?)
				}
		.runMethod(observer, instance, 'Teardown', false, timeEachMethod?)
		}
	runMethod(observer, instance, method, skipRetrySleep, timeEachMethod?)
		{
		Suneido.SkipRetrySleep = skipRetrySleep
		observer.BeforeMethod(method)
		time = false
		if timeEachMethod?
			time = Timer()
				{ passed? = .tryMethod(observer, instance, method) }
		else
			{
// CountSlow(Curry(observer.Warning, method))
//	{
			passed? = .tryMethod(observer, instance, method)
//	}
			}
		observer.AfterMethod(method, :time)
		Suneido.Delete(#SkipRetrySleep)
		return passed?
		}
	tryMethod(observer, instance, method)
		{
		try
			instance[method]()
		catch (e)
			{
			.error(observer, method, e)
			return false
			}
		return true
		}
	measure(block)
		{
		dbsize = Database.CurrentSize()
		mem = MemoryAlloc()
		t = Timer(block)
		return Object(
			time: t,
			dbgrowth: Database.CurrentSize() - dbsize,
			memory: MemoryAlloc() - mem)
		}
	}
