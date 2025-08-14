// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_model_name?()
		{
		Assert(AccessModel.AccessModel_model_name?('stdlib') is: false)
		Assert(AccessModel.AccessModel_model_name?('Record'))
		}

	Test_init()
		{
		model = AccessModel(#('stdlib'), false)
		Assert(model.GetQuery() is: 'stdlib')
		Assert(model.GetFields() has: 'num')
		rec = [num: 'test_num']
		Assert(model.GetKeyField() is: 'num')
		Assert(model.GetLockKey(rec) is: 'test_num')
		Assert(model.SetKeyQuery(rec) sameText: 'stdlib where num is "test_num"')

		model.SetQuery('stdlib where num is "1"')
		Assert(model.GetQuery() is: 'stdlib where num is "1"')
		model.GetCursor().Close()
		}

	Test_query()
		{
		model = AccessModel(#('stdlib remove num sort name'), false)
		Assert(model.GetQuery() is: 'stdlib remove num sort name')
		Assert(model.GetFields() hasnt: 'num')
		rec = [num: 'test_num', name: 'Test_num', group: -1]
		Assert(model.GetKeyField() is: 'name,group')
		Assert(model.GetLockKey(rec) is: 'Test_num\x01-1')
		Assert(model.SetKeyQuery(rec)
			sameText: 'stdlib remove num where name is "Test_num" where group is -1')

		Assert(model.AddMoreToQuery('where name.Has?("Init")'))
		Assert(model.GetQuery()	sameText:
			'stdlib remove num where name.Has?("Init") sort name')

		rec = model.LookupRecord('name', 'SuneidoLog')
		Assert(rec.name is: 'SuneidoLog')
		Assert(rec.group is: -1)

		model.GetCursor().Close()
		}

	Test_locate()
		{
		model = AccessModel(#('stdlib', locate: #(keys: #('num', 'name,group'))), false)
		Assert(model.GetLocateKeys() is: #('Date/Time Created': 'num'))
		model.GetCursor().Close()
		}

	Test_FindGotoField()
		{
		model = AccessModel(#('stdlib rename name to name_new'), false)
		Assert(model.FindGotoField('name_new') is: 'name_new')
		model.GetCursor().Close()

		model = AccessModel(#('stdlib'), false)
		Assert(model.FindGotoField('name') is: 'name')
		Assert(model.FindGotoField('test_not_exist') is: false)
		model.GetCursor().Close()
		}
	}
