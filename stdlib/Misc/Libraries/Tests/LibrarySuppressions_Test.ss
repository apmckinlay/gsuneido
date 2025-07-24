// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		Suneido.LibrarySuppressionsOverride = #(StagedCached)
		}

	Test_main()
		{
		// tests must be ran in a certain order as they affect each other
		// failure tests create error suneidologs in the main and reset tests
		.test_main()
		.test_function_ResetRequired?()
		.test_class_ResetRequired?()
		.test_object_ResetRequired?()
		.test_errorHandling_suppressRef()
		.test_errorHandling_ResetRequired?()
		.test_errorHandling_suppressList()
		}

	test_main()
		{
		fake1lib = .MakeLibrary()
		fake2lib = .MakeLibrary()
		list = LibrarySuppressions.LibrarySuppressions_suppressList

		Assert(LibrarySuppressions.Func(fake1lib) is: #())
		Assert(list(fake1lib) is: false)
		Assert(LibrarySuppressions.Func(fake2lib) is: #())
		Assert(list(fake2lib) is: false)
		libraries = LibrarySuppressions.Libraries()
		Assert(libraries hasnt: fake1lib)
		Assert(libraries hasnt: fake2lib)

		.MakeLibraryRecord(
			[name: fake1lib.Capitalize() $ '_CheckLibrarySuppressions',
				text: `function ()
					{
					return #(TestRec1, TestRec2)
					}`],
			[name: 'TestRec1', text: `function () { }`],
			[name: 'TestRec2', text: `function () { }`],
			table: fake1lib)
		Assert(LibrarySuppressions.Func(fake1lib) is: #(TestRec1, TestRec2))
		Assert(LibrarySuppressions.Func(fake2lib) is: #())
		Assert(list(fake2lib) is: false)
		libraries = LibrarySuppressions.Libraries()
		Assert(libraries has: fake1lib)
		Assert(libraries hasnt: fake2lib)

		.MakeLibraryRecord(
			[name: fake2lib.Capitalize() $ '_CheckLibrarySuppressions',
				text: `#(TestRec3, TestRec4)`],
			[name: 'TestRec3', text: `function () { }`, group: -2], // Deleted
			[name: 'TestRec4', text: `function () { }`],
			[name: 'TestRec5', text: `function () { }`],
			table: fake2lib)
		Assert(LibrarySuppressions.Func(fake1lib) is: #(TestRec1, TestRec2))
		Assert(LibrarySuppressions.Func(fake2lib) is: #(TestRec3, TestRec4))
		libraries = LibrarySuppressions.Libraries()
		Assert(libraries has: fake1lib)
		Assert(libraries has: fake2lib)
		}

	test_class_ResetRequired?()
		{
		fakelib = .MakeLibrary()
		suppressName = fakelib.Capitalize() $ '_CheckLibrarySuppressions'
		.MakeLibraryRecord(
			[name: suppressName,
				text: `class
					{
					CallClass()
						{
						return QueryAll(.query()).Map({ it.name })
						}
					query()
						{
						return "` $ fakelib $
							` where group is -1 and name isnt '` $ suppressName $ `'"
						}
					Suppressed?(name)
						{
						return Query1(.query(), :name) isnt false
						}
					}`],
			[name: 'NotYetCached',
				text: `class { /* Should be suppressed, not yet cached */ }`],
			table: fakelib)
		Assert(LibrarySuppressions(fakelib) is: #(StagedCached))
		Assert(LibrarySuppressions.ResetRequired?(fakelib, 'NotYetCached'))
		}

	test_function_ResetRequired?()
		{
		fakelib = .MakeLibrary()
		suppressName = fakelib.Capitalize() $ '_CheckLibrarySuppressions'
		.MakeLibraryRecord(
			[name: suppressName,
				text: `function()
					{
					return #(StagedCached)
					}`],
			[name: 'NotCached', text: `class { }`],
			table: fakelib)
		Assert(LibrarySuppressions(fakelib) is: #(StagedCached))
		Assert(LibrarySuppressions.ResetRequired?(fakelib, 'NotYetCached') is: false)
		}

	test_object_ResetRequired?()
		{
		fakelib = .MakeLibrary()
		suppressName = fakelib.Capitalize() $ '_CheckLibrarySuppressions'
		.MakeLibraryRecord(
			[name: suppressName, text: `#(StagedCached)`],
			[name: 'NotCached', text: `class { }`],
			table: fakelib)
		Assert(LibrarySuppressions(fakelib) is: #(StagedCached))
		Assert(LibrarySuppressions.ResetRequired?(fakelib, 'NotYetCached') is: false)
		}

	test_errorHandling_suppressRef()
		{
		fakelib = .MakeLibrary()
		suppressName = fakelib.Capitalize() $ '_CheckLibrarySuppressions'
		.MakeLibraryRecord(
			[name: suppressName,
				text: `class
					{
					syntaxError
					CallClass()
						{
						return #(StagedCached)
						}
					}`],
			[name: 'NotCached', text: `class { }`],
			table: fakelib)
		spyon = .SpyOn(LibrarySuppressions.LibrarySuppressions_log)
		logs = spyon.Return(false).CallLogs()

		Assert(LibrarySuppressions.Func('doesNotExistLib') is: #())
		Assert(logs isSize: 0)

		Assert(LibrarySuppressions.Func(fakelib) is: #())
		log = logs.PopFirst()
		Assert(log.library is: fakelib)
		Assert(log.error startsWith: 'error loading')
		Assert(log.event is: 'reference')
		spyon.Close()
		}

	test_errorHandling_suppressList()
		{
		fakelib = .MakeLibrary()
		suppressName = fakelib.Capitalize() $ '_CheckLibrarySuppressions'
		.MakeLibraryRecord(
			[name: suppressName, text: `function () { }`],
			[name: 'NotCached', text: `class { }`],
			table: fakelib)
		spyon = .SpyOn(LibrarySuppressions.LibrarySuppressions_log)
		logs = spyon.Return(false).CallLogs()

		Assert(LibrarySuppressions.Func('doesNotExistLib') is: #())
		Assert(logs isSize: 0)

		Assert(LibrarySuppressions.Func(fakelib) is: #())
		log = logs.PopFirst()
		Assert(log.library is: fakelib)
		Assert(log.error is: 'no return value')
		Assert(log.event is: 'suppression list')
		spyon.Close()
		}

	test_errorHandling_ResetRequired?()
		{
		fakelib = .MakeLibrary()
		suppressName = fakelib.Capitalize() $ '_CheckLibrarySuppressions'
		.MakeLibraryRecord(
			[name: suppressName,
				text: `class
					{
					CallClass()
						{
						return #(StagedCached)
						}
					Suppressed?() // Missing argument: name
						{
						}
					}`],
			[name: 'NotCached', text: `class { }`],
			table: fakelib)
		spyon = .SpyOn(LibrarySuppressions.LibrarySuppressions_log)
		logs = spyon.Return(false).CallLogs()

		Assert(LibrarySuppressions.Func('doesNotExistLib') is: #())
		Assert(logs isSize: 0)

		Assert(LibrarySuppressions.Func(fakelib) is: #(StagedCached))
		Assert(logs isSize: 0)

		Assert(LibrarySuppressions.ResetRequired?(fakelib, 'NotCached') is: false)
		log = logs.PopFirst()
		Assert(log.library is: fakelib)
		Assert(log.error is: 'too many arguments')
		Assert(log.event is: 'reset required')
		spyon.Close()
		}

	Teardown()
		{
		Suneido.Delete('LibrarySuppressionsOverride') // Ensure override is removed
		super.Teardown()
		}
	}
