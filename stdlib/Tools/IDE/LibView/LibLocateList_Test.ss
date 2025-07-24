// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_GetMatches()
		{
		mock = Mock(LibLocateList)
		mock.When.getNames([anyArgs:]).Return([])
		mock.When.getNames(1, 'stdlib').Return([#Aaa_Bb_Cc, #Abc, #Abc_Test])
		mock.When.getNames(2, 'test_not_in_used_lib').Return([#Abc])
		mock.When.forceRun?().Return(false)
		mock.When.getList([anyArgs:]).CallThrough()

		m = mock.LibLocateList_getMatches
		prevPadding = false
		for extraLibs, expectedPadding in #(0: 1, 10: 2, 100: 3)
			{
			libs = Object('Builtin', 'stdlib', 'test_not_in_used_lib')
			for i in ..extraLibs
				libs.Add('test' $ i $ 'lib')
			info = Object(libNums: #(stdlib: 6174), :libs, list: mock.getList(libs, []))

			// Ensuring that IF the padding changes the names are as expected
			for name in info.list
				Assert(name.AfterLast('%') isSize: expectedPadding)
			// Ensuring that the padding actually differs between test loops
			Assert(prevPadding isnt: expectedPadding)
			prevPadding = expectedPadding

			msg = 'Libs: ' $ extraLibs $ ', Padding: ' $ expectedPadding
			Assert(m(info, 'ab', justName:) is: #(Aaa_Bb_Cc, Abc, Abc_Test), :msg)
			Assert(m(info, 'at', justName:) is: #(Abc_Test), :msg)
			Assert(m(info, 'am', justName:) is: #(), :msg)
			// exact match should be in the front
			Assert(m(info, 'abc', justName:) is: #(Abc, Aaa_Bb_Cc, Abc_Test), :msg)
			// exact match should be in the front
			Assert(m(info, 'abc')
				is: #('Abc - stdlib',
					'Abc - (test_not_in_used_lib)',
					'Aaa_Bb_Cc - stdlib',
					'Abc_Test - stdlib'), :msg)
			}
		}
	Test_validRecord()
		{
		fn = LibLocateList.LibLocateList_validRecord

		Assert(fn(false, 'fakelib') isnt: '', msg: 'invalid record')
		Assert(fn('somethingvalid', 'fakelib') is: '', msg: 'valid record')
		}

	Test_processName()
		{
		cl = LibLocateList
			{ LibLocateList_printError(unused) { } }
		m = cl.LibLocateList_processName

		list = Object()
		m(false, 'fakelib', 0, list)
		Assert(list.Empty?(), msg: 'invalid record, empty list')

		m('fakeRecord0', 'fakelib', 0, list)
		Assert(list isSize: 1)
		Assert(list[0] is: 'fakerecord0+fakeRecord0%0')

		m('FakeRecord1', 'fakelib', 1, list)
		Assert(list isSize: 3)
		Assert(list[1] is: 'fakerecord1+FakeRecord1%1')
		Assert(list[2] is: 'fr+FakeRecord1%1')

		m('fake_record2', 'fakelib', 2, list)
		Assert(list isSize: 4)
		Assert(list[3] is: 'fakerecord2+fake_record2%2')

		m('Fake_Record3', 'fakelib', 3, list)
		Assert(list isSize: 6)
		Assert(list[4] is: 'fakerecord3+Fake_Record3%3')
		Assert(list[5] is: 'fr+Fake_Record3%3')

		m('Fake_record4', 'fakelib', 4, list)
		Assert(list isSize: 7)
		Assert(list[6] is: 'fakerecord4+Fake_record4%4')
		}

	Test_getNames()
		{
		m = LibLocateList.LibLocateList_getNames

		Assert(m(0, '') is: BuiltinNames()) 		// Checks builtins
		Assert(m(1, 'notALib') is: #()) 			// throws / catches invalid query

		lib = .MakeLibrary(
			[name: 'Test1'], [name: 'Test2'], [name: 'Test3'],
			[name: 'Test4'], [name: 'Test5'], [name: 'Test6'])
		Assert(m(1, lib) equalsSet: #(Test1, Test2, Test3, Test4, Test5, Test6))
		}
	}