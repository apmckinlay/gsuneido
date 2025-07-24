// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		.TearDownIfTablesNotExist('userselects')
		AccessSelectMgr.Ensure()
		.user = Suneido.User
		Suneido.User = 'AccessSelectMgrTestUser'
		.MakeLibraryRecord([name: "Field_number_withprompt",
			text: `Field_number { Prompt: "Number With Prompt" }`])
		}
	Test_new()
		{
		mgr = AccessSelectMgr()
		Assert(mgr.Select_vals() is: #())

		mgr = AccessSelectMgr(#(
			(date, '==', 123),
			(number_withprompt, '==', '20190101'),
			(city, '>', 456)))
		Assert(mgr.Select_vals()
			is: #(
			[condition_field: 'date', check:,
				date: #(operation: 'equals', value: 123, value2: '')],
			[condition_field: 'number_withprompt', check:,
				number_withprompt: #(operation: 'equals', value: 20190101, value2: '')],
			[condition_field: 'city', check:,
				city: #(operation: 'greater than', value: '456', value2: '')]))
		}

	Test_reset()
		{
		mgr = AccessSelectMgr(#(
			(date, '==', 123),
			(city, '>', 456)))
		mgr.Reset(#(
			(city, '==', 'Saskatoon'),
			(date, '<=', '#20190101')))
		Assert(mgr.Select_vals()
			is: #(
			[condition_field: 'city', check:,
				city: #(operation: 'equals', value: 'Saskatoon', value2: '')],
			[condition_field: 'date', check:,
				date: #(operation: 'less than or equal to',
					value: #20190101, value2: '')]))
		mgr.Reset(#())
		Assert(mgr.Select_vals() is: #())
		}
	Test_SaveSelects()
		{
		// selects empty, delete should be done but not output
		acc = AccessSelectMgr
			{
			Ensure() { }
			AccessSelectMgr_delete_selects(t /*unused*/) { }
			AccessSelectMgr_output_selects(t /*unused*/, selects)
				{
				if selects.Empty?()
					return
				throw "output called"
				}
			}
		mgr = acc()
		mgr.SaveSelects()

		mgr = .createAccessSelectMgr(addSelects: Object(.selectVal(check:)),
			saveSelects: Object(.selectVal()))
		mgr.SaveSelects()
		mgr.Finalize()

		mgr = .createAccessSelectMgr(addSelects: false, saveSelects: false)
		mgr.SaveSelects()
		mgr.Finalize()

		mgr = .createAccessSelectMgr(
			addSelects: Object(.selectVal(check:), .initVal(check:)),
			saveSelects: Object(.selectVal(), .initVal()),
			defaultSelects: Object([condition_field: 'state_prov', check:,
				state_prov: #(value: 'SK', operation: 'equals', value2: '')]))
		mgr.SaveSelects()
		mgr.Finalize()
		}

	accessMgr: AccessSelectMgr
		{
		transaction: false
		SetTransaction(transaction)
			{
			/* expects a MockObject */
			.transaction = transaction
			}
		defaultSel: #()
		SetDefaultSel(defaultSel)
			{
			.defaultSel = defaultSel
			}
		Finalize()
			{
			.transaction.Finalize()
			}
		Ensure() { }
		AccessSelectMgr_retrySave(block)
			{
			block(.transaction)
			}
		AccessSelectMgr_querySelect(savedDefault? = false)
			{
			if savedDefault? is true
				return Object(userselect_selects: .defaultSel)
			return false
			}
		}

	// mgr.Select_vals() simulates user adding to the select screen
	// mgr.LoadSelects simulates loading the prev saved selects
	createAccessSelectMgr(addSelects, saveSelects, defaultSelects = false)
		{
		mgr = (.accessMgr)(#((city, '==', 'Saskatoon')), name: 'Test')
		deleteQuery = mgr.AccessSelectMgr_userselects_query()
		calls = Object(Object('QueryDo', 'delete ' $ deleteQuery))

		// needs to happen before LoadSelects
		if defaultSelects isnt  false
			mgr.SetDefaultSel(defaultSelects)

		sf = SelectFields(#(city, state_prov, number_withprompt))
		ctrl = FakeObject(GetSelectFields: sf)
		option = .nextName()
		ctrl.Option = option
		mgr.LoadSelects(ctrl)

		if false isnt addSelects
			{
			calls.Add(Object('QueryOutput', 'userselects',
				Record(userselect_user: Suneido.User, userselect_title: 'Test',
					userselect_selects: saveSelects)))
			mgr.Select_vals().Add(@addSelects)
			}
		calls.Add(#(Finalize))
		transaction = MockObject(calls)
		mgr.SetTransaction(transaction)
		return mgr
		}
	// can't inherit from BizTests
	nextName()
		{
		return 'biztests_' $ Display(Timestamp())[1..].Tr('.', '_')
		}
	selectVal(check = false)
		{
		return Record(condition_field: 'number_withprompt', :check,
			number_withprompt: #(operation: 'equals', value: 20190101, value2: ''))
		}
	initVal(check = false)
		{
		return [condition_field: 'city', :check,
			city: #(value: 'Saskatoon', operation: 'equals', value2: '')]
		}

	Test_FilterInitial()
		{
		select_vals = Object(
			[check: false, condition_field: "etaequip_terminate_date",
				etaequip_terminate_date: #(value: "", value2: "", operation: "empty")],
			[check: false, condition_field: "etaequip_type",
				etaequip_type: #(value: "trailer", value2: "", operation: "equals")],
			[check: true, condition_field: "etaequip_terminate_date",
				etaequip_terminate_date: #(value: "", operation: "", value2: "")])
		initial_vals = [
			[check:, condition_field: "etaequip_terminate_date",
				etaequip_terminate_date: #(value: "", operation: "empty", value2: "")]]
		acc = AccessSelectMgr.AccessSelectMgr_filterInitial
		selects = acc(select_vals, initial_vals)
		Assert(selects is: #(
			[check: false, condition_field: "etaequip_type",
				etaequip_type: #(value: "trailer", value2: "", operation: "equals")],
			[check: false, condition_field: "etaequip_terminate_date",
				etaequip_terminate_date: #(value: "", operation: "", value2: "")]))
		}

	Test_buildCondition()
		{
		m = AccessSelectMgr.AccessSelectMgr_buildCondition
		Assert(m('string', 'equals', 'test')
			is: [condition_field: 'string',
				'string': #(operation: 'equals', value: 'test', value2: '')])
		Assert(m('string', 'not equal to', 'test')
			is: [condition_field: 'string',
				'string': #(operation: 'not equal to', value: 'test', value2: '')])

		Assert(m('string', 'not empty', 'test')
			is: [condition_field: 'string',
				'string': #(operation: 'not empty', value: '', value2: '')])
		Assert(m('string', 'empty', 'test')
			is: [condition_field: 'string',
				'string': #(operation: 'empty', value: '', value2: '')])

		Assert(m('string', 'equals', '')
			is: [condition_field: 'string',
				'string': #(operation: 'empty', value: '', value2: '')])
		Assert(m('string', 'not equal to', '')
			is: [condition_field: 'string',
				'string': #(operation: 'not empty', value: '', value2: '')])

		Assert(m('string', 'greater than', '')
			is: [condition_field: 'string',
				'string': #(operation: 'not empty', value: '', value2: '')])

		Assert(m('string', 'less than or equal to', '')
			is: [condition_field: 'string',
				'string': #(operation: 'empty', value: '', value2: '')])

		// matches "" should be converted to all
		Assert(m('string', 'matches', '')
			is: [condition_field: 'string',
				'string': #(operation: '', value: '', value2: '')])

		// starts with "" should be converted to all
		Assert(m('string', 'starts with', '')
			is: [condition_field: 'string',
				'string': #(operation: '', value: '', value2: '')])

		Assert(m('string', 'matches', 'test')
			is: [condition_field: 'string',
				'string': #(operation: 'matches', value: 'test', value2: '')])
		Assert(m('string', 'does not match', 'test')
			is: [condition_field: 'string',
				'string': #(operation: 'does not match', value: 'test', value2: '')])

		Assert(m('string', '', 'test')
			is: [condition_field: 'string',
				'string': #(operation: '', value: '', value2: '')])
		Assert(m('string', '', '')
			is: [condition_field: 'string',
				'string': #(operation: '', value: '', value2: '')])

		Assert(m('date', 'equals', '2019-01-01')
			is: [condition_field: 'date',
				'date': #(operation: 'equals', value: #20190101, value2: '')])
		}

	Test_LoadSelects()
		{
		.WatchTable('userselects')
		selName = .TempName()
		ctrl = class
			{
			Option: 'test'
			GetSelectFields()
				{
				return SelectFields(#(city, country, date))
				}
			}
		mgr = AccessSelectMgr(#((city, '=', 'Saskatoon')), name: selName)
		mgr.LoadSelects(ctrl)
		Assert(mgr.Select_vals()
			is: #([condition_field: 'city', check:,
				city: #(operation: 'equals', value: 'Saskatoon', value2: '')]))

		QueryOutput('userselects', [
			userselect_user: 'AccessSelectMgrTestUser'
			userselect_title: selName
			userselect_selects: #([condition_field: 'country', check: false,
				country: #(operation: 'not empty', value: '', value2: '')])
			])

		mgr = AccessSelectMgr(#((city, '=', 'Saskatoon')), name: selName)
		mgr.LoadSelects(ctrl)
		Assert(mgr.Select_vals()
			is: #([condition_field: 'city', check:,
				city: #(operation: 'equals', value: 'Saskatoon', value2: '')],
				[condition_field: 'country', check: false,
				country: #(operation: 'not empty', value: '', value2: '')]))

		QueryOutput('userselects', [
			userselect_user: 'AccessSelectMgrTestUser'
			userselect_title: selName $ '~default'
			userselect_selects: #([condition_field: 'date', check:,
				date: #(operation: 'equals', value: #20190101, value2: '')],
				[condition_field: 'invalid', check:,
					invalid: #(operation: 'not empty', value: '', value2: '')])
			])
		mgr = AccessSelectMgr(#((city, '=', 'Saskatoon')), name: selName)
		mgr.LoadSelects(ctrl)
		Assert(mgr.Select_vals()
			is: #([condition_field: 'date', check:,
				date: #(operation: 'equals', value: #20190101, value2: '')],
				[condition_field: 'country', check: false,
				country: #(operation: 'not empty', value: '', value2: '')]))

		QueryDo('delete userselects where userselect_user is "AccessSelectMgrTestUser"')
		QueryOutput('userselects', [
			userselect_user: 'AccessSelectMgrTestUser'
			userselect_title: selName
			userselect_selects: #(
				[condition_field: 'city', check: true,
					city: #(operation: 'equals', value: 'Saskatoon', value2: '')],
				[condition_field: 'country', check: true,
					country: #(operation: 'not empty', value: '', value2: '')])
			])
		mgr = AccessSelectMgr(#((city), (date)), name: selName)
		mgr.LoadSelects(ctrl)
		Assert(mgr.Select_vals()
			is: #(
				[condition_field: 'date', check: false,
					date: #(operation: '', value: '', value2: '')],
				[condition_field: 'city', check: false,
					city: #(operation: 'equals', value: 'Saskatoon', value2: '')],
				[condition_field: 'country', check: false,
					country: #(operation: 'not empty', value: '', value2: '')]))

		QueryDo('delete userselects where userselect_user is "AccessSelectMgrTestUser"')
		QueryOutput('userselects', [
			userselect_user: 'AccessSelectMgrTestUser'
			userselect_title: selName
			userselect_selects: #(
				[condition_field: 'city', check: false,
					city: #(operation: '', value: '', value2: '')],
				[condition_field: 'country', check: false,
					country: #(operation: 'not empty', value: '', value2: '')])
			])
		mgr = AccessSelectMgr(#((city), (date)), name: selName)
		mgr.LoadSelects(ctrl)
		Assert(mgr.Select_vals()
			is: #(
				[condition_field: 'city', check: false,
					city: #(operation: '', value: '', value2: '')],
				[condition_field: 'date', check: false,
					date: #(operation: '', value: '', value2: '')],
				[condition_field: 'country', check: false,
					country: #(operation: 'not empty', value: '', value2: '')]))

		// screen initial select and user history select more than max
		cl = AccessSelectMgr
			{ AccessSelectMgr_maxSelectRecords() { return 5 } }
		QueryDo('delete userselects where userselect_user is "AccessSelectMgrTestUser"')
		QueryOutput('userselects', [
			userselect_user: 'AccessSelectMgrTestUser'
			userselect_title: selName
			userselect_selects: #(
				[condition_field: 'city', check: false,
					city: #(operation: 'contains', value: 'a', value2: '')],
				[condition_field: 'country', check: false,
					country: #(operation: 'contains', value: 'b', value2: '')]
				[condition_field: 'city', check: false,
					city: #(operation: 'contains', value: 'c', value2: '')],
				[condition_field: 'country', check: false,
					country: #(operation: 'contains', value: 'd', value2: '')]
				[condition_field: 'city', check: false,
					city: #(operation: 'contains', value: 'e', value2: '')])
			])
		mgr = cl(#((city), (date)), name: selName)
		mgr.LoadSelects(ctrl)
		Assert(mgr.Select_vals()
			is: #(
				[condition_field: 'date', check: false,
					date: #(operation: '', value: '', value2: '')],
				[condition_field: 'city', check: false,
					city: #(operation: 'contains', value: 'a', value2: '')],
				[condition_field: 'country', check: false,
					country: #(operation: 'contains', value: 'b', value2: '')],
				[condition_field: 'city', check: false,
					city: #(operation: 'contains', value: 'c', value2: '')],
				[condition_field: 'country', check: false,
					country: #(operation: 'contains', value: 'd', value2: '')]
				))

		// user default select and user history select more than max
		QueryOutput('userselects', [
			userselect_user: 'AccessSelectMgrTestUser'
			userselect_title: selName $ '~default'
			userselect_selects: #([condition_field: 'date', check:,
				date: #(operation: 'equals', value: #20190101, value2: '')],
				[condition_field: 'invalid', check:,
					invalid: #(operation: 'not empty', value: '', value2: '')])
			])
		mgr = cl(#(), name: selName)
		mgr.LoadSelects(ctrl)
		Assert(mgr.Select_vals()
			is: #(
				[condition_field: 'date', check:,
					date: #(operation: 'equals', value: #20190101, value2: '')]
				[condition_field: 'city', check: false,
					city: #(operation: 'contains', value: 'a', value2: '')],
				[condition_field: 'country', check: false,
					country: #(operation: 'contains', value: 'b', value2: '')],
				[condition_field: 'city', check: false,
					city: #(operation: 'contains', value: 'c', value2: '')],
				[condition_field: 'country', check: false,
					country: #(operation: 'contains', value: 'd', value2: '')]
				))
		}

	Test_convertSimpleSelects()
		{
		convert = AccessSelectMgr.AccessSelectMgr_convertSimpleSelects
		convert(#(), sel = Object())
		Assert(sel is: #())

		convert(#(('date')), sel = Object())
		Assert(sel is: #([check: false, condition_field: "date",
			date: #(value: "", value2: "", operation: "")]))

		convert(#(('date', '=', #20200303)), sel = Object())
		Assert(sel is: #([check:, condition_field: "date",
			date: #(value: #20200303, value2: "", operation: "equals")]))

		convert(#(('date', '>', #20200303), (city, '=', 'Saskatoon')), sel = Object())
		Assert(sel is: #([check:, condition_field: "date",
			date: #(value: #20200303, value2: "", operation: "greater than")],
			[check:, condition_field: "city",
			city: #(value: 'Saskatoon', value2: "", operation: "equals")]))
		}

	Teardown()
		{
		Suneido.User = .user
		super.Teardown()
		}
	}
