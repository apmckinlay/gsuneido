// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.MakeTable('(test) key()') // to ensure teardown for temp params table
		.paramsTable = .TempTableName()
		Params.Ensure(.paramsTable)
		userParamsTestCl = UserParams
			{ SetParamsTable(table) { .UserParams_table = table } }
		.userParamsTester = new userParamsTestCl
		.userParamsTester.SetParamsTable(.paramsTable)
		}

	Test_Sync()
		{
		name1 = .TempName()
		name2 = .TempName()

		name2Preset1 = name2 $ '~presets~' $ 'fred'
		name2Preset2 = name2 $ '~presets~' $ 'barney'
		name1Preset1 = name1 $ '~presets~' $ 'barney1'

		// no saved params
		.userParamsTester.Sync(name1, name2, #(abc), '')

		.userParamsTester.Save(name2Preset1, [abc: 111, efg: 666], user: '')
		.userParamsTester.Save(name2Preset2, [abc: 555, efg: 222], user: '')

		.userParamsTester.Save(name1, [abc: 1, efg: 2])
		.userParamsTester.Sync(name2, name1, #(abc), '') // empty to exist, no change

		Assert(.userParamsTester.Load(name1) is: [abc: 1, efg: 2])
		Assert(.userParamsTester.Load(name1 $ '~presets~' $ 'fred', user: '')
			is: [abc: 111])
		Assert(.userParamsTester.Load(name1 $ '~presets~' $ 'barney', user: '')
			is: [abc: 555])

		.userParamsTester.Save(name2, [abc: 1, efg: 3])
		.userParamsTester.Sync(name2, name1, #(abc), '') // no change for abc
		Assert(.userParamsTester.Load(name1) is: [abc: 1, efg: 2])

		.userParamsTester.Save(name2, [abc: 2, efg: 3])
		.userParamsTester.Sync(name2, name1, #(abc), '') // changes should be synced
		Assert(.userParamsTester.Load(name1) is: [abc: 2, efg: 2])

		// new params
		.userParamsTester.Save(name1Preset1, [abc: 999, efg: 555], user: '')
		.userParamsTester.Sync(name1, name2, #(abc), '')
		Assert(.userParamsTester.Load(name2) is: [abc: 2, efg: 3])
		Assert(.userParamsTester.Load(name2 $ '~presets~' $ 'fred', user: '')
			is: [abc: 111, efg: 666])
		Assert(.userParamsTester.Load(name2 $ '~presets~' $ 'barney', user: '')
			is: [abc: 555, efg: 222])
		Assert(.userParamsTester.Load(name2 $ '~presets~' $ 'barney1', user: '')
			is: [abc: 999])

		// deleting params
		QueryDo('delete ' $ .paramsTable $ ' where report is ' $ Display(name1Preset1))
		.userParamsTester.Sync(name1, name2, #(abc), '')
		Assert(.userParamsTester.Load(name2) is: [abc: 2, efg: 3])
		Assert(.userParamsTester.Load(name2 $ '~presets~' $ 'fred', user: '')
			is: [abc: 111, efg: 666])
		Assert(.userParamsTester.Load(name2 $ '~presets~' $ 'barney', user: '')
			is: [abc: 555, efg: 222])
		Assert(.userParamsTester.Load(name2 $ '~presets~' $ 'barney1', user: '')
			is: [])

		// rename params
		QueryDo('update ' $ .paramsTable $ ' where report is ' $ Display(name2Preset2) $
			' set report = ' $ Display(name2 $ '~presets~' $ 'test'))
		.userParamsTester.Sync(name2, name1, #(abc), '')
		Assert(.userParamsTester.Load(name1) is: [abc: 2, efg: 2])
		Assert(.userParamsTester.Load(name1 $ '~presets~' $ 'fred', user: '')
			is: [abc: 111])
		Assert(.userParamsTester.Load(name1 $ '~presets~' $ 'test', user: '')
			is: [abc: 555])
		Assert(.userParamsTester.Load(name1 $ '~presets~' $ 'barney', user: '')
			is: [])
		Assert(.userParamsTester.Load(name2 $ '~presets~' $ 'barney1', user: '')
			is: [])
		}

	Test_PresetsNotDuplicatedAfterSync()
		{
		.testForSavedUser('', 'fred')
		.testForSavedUser('fred', 'fred')
		.testForSavedUser('barney', 'fred')
		.testForSavedUser('fred', 'barney')
		.testForSavedUser('barney', '')
		.testForSavedUser('', 'fred')
		}

	testForSavedUser(savedUser, forUser)
		{
		report1 = .TempName()
		report2 = .TempName()
		preset1 = report1 $ '~presets~' $ 'test preset name'
		preset2 = report2 $ '~presets~' $ 'test preset name'
		.userParamsTester.Save(preset2, [yyy: 1, zzz: 2], user: savedUser)
		.userParamsTester.Save(preset1, [yyy: 11, zzz: 22], user: savedUser)
		.userParamsTester.Sync(report1, report2, #(yyy), forUser)
		q = .paramsTable $ ' where report is ' $ Display(preset2)
		Assert(QueryCount(q) is: 1)
		rec = Query1(q)
		Assert(rec.user is: savedUser)
		Assert(rec.params.yyy is: 11)
		Assert(rec.params.zzz is: 2)
		}
	}
