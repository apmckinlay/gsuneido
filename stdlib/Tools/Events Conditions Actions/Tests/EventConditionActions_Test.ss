// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	actionsTableExisted?: true
	conditionsTableExisted?: true
	Setup()
		{
		super.Setup()
		.actionsTableExisted? = TableExists?('event_actions')
		.conditionsTableExisted? = TableExists?('event_condition_actions')
		EventConditionActions.EnsureTables()
		}
	Test_validate_plugins()
		{
		getPluginInfo = EventConditionActions.EventConditionActions_getPluginInfo
		invalid = false
		Plugins().ForeachContribution('ECA', 'event')
			{ |c|
			for infoType in c.GetDefault('produce', #()).Members()
				if getPluginInfo(infoType, 'info_type') is false
					invalid = infoType
			}
		Assert(invalid is false, 'Info Type: ' $ Display(invalid) $ ' not defined')

		invalid = false
		Plugins().ForeachContribution('ECA', 'action')
			{ |c|
			for infoType in c.GetDefault('requires', #())
				if getPluginInfo(infoType, 'info_type') is false
					invalid = infoType
			}
		Assert(invalid is false, 'Info Type: ' $ Display(invalid) $ ' not defined')

		.checkPluginArgs()
		}

	checkPluginArgs()
		{
		invalid = ''
		Plugins().ForeachContribution('ECA', 'action')
			{ |c|
			try
				{
				requiredArgs = c.GetDefault('requires', Object()).Copy()
				for arg in c.GetDefault('setting', #())
					requiredArgs.Add(arg)

				f = Global(c.func)
				func = Class?(f) ? Global(c.func $ '.CallClass') : f

				// Removing logError from expected params as it is defaulted to false
				// As well is populated by OptContribution
				if func.Params().Split(',').Remove('logError=false)').Size() isnt
					requiredArgs.Size()
					invalid $= 'Action: ' $ Display(func) $
						' does not match expecting arguments: ' $
						requiredArgs.Join(',') $ '; '
				}
			catch(e)
				invalid $= 'Action ' $ c.func $ ' does not exist: ' $ e $ '; '
			}
		Assert(invalid is '', invalid)
		}

	Test_GetEventSourceFields()
		{
		.mockCase('EventConditionActions', 'GetEventSourceFields',
			args: #('test'), expect: #(),
			when: [ ['getPluginInfo', [[anyArgs:]], false]] )

		table = .MakeTable('(name, abbrev) key ()')
		.mockCase('EventConditionActions', 'GetEventSourceFields',
			args: #('Test Event'), expect: #(
				"Test Result": #()
				"Test Result2": #(Name: 'name', Abbreviation: 'abbrev')),
			when: [
				[#getPluginInfo, ['Test Event', 'event'],
					#(ECA, event, name: 'Test Event',
					produce: (
						'Test Result': function (testarg) { return testarg },
						'Test Result2': function (testarg) { return testarg } ))],
				[#getPluginInfo, ['Test Result', 'info_type'],
					#(ECA, info_type, name: 'Test Result', type: 'string')],
				[#getPluginInfo, ['Test Result2', 'info_type'],
					Object(#ECA, #info_type, name: 'Test Result2', source: table)],
				[#extraFields, [[anyArgs:]], '']
				] )
		}

	Test_execute()
		{
		Suneido.Delete('EventConditionActions_Test_actionFunc')
		actionFunc = "function (@args) {
			Assert(args[0] is: 'testing2')
			Assert(args[1] is: 'test setting')
			Suneido.EventConditionActions_Test_actionFunc = true
			}"
		.MakeLibraryRecord([name: 'TestActionFunc', text: actionFunc])
		actionInfo = Object(#ECA, #action, name: 'Test Action',
			func: 'TestActionFunc', requires: #('Test Result'))
		.mockCase('EventConditionActions', 'execute',
			args: #('Test Action', 'Test Event', #(testarg: 'testing1'),
				#([field: 'test_setting', value: 'test setting']))
			when: [
				[#getPluginInfo, ['Test Action', 'action'], actionInfo],
				[#getActionRequired, ['Test Event', actionInfo, #(testarg: 'testing1')],
					Object('testing2')]
				])
		Assert(Suneido.EventConditionActions_Test_actionFunc)
		Suneido.Delete('EventConditionActions_Test_actionFunc')
		}

	Test_getActionRequired()
		{
		.mockCase('EventConditionActions', 'getActionRequired',
			args: #('Test Event', [requires: #('Test Result')],
				#(testarg: 'testing'))
			expect: #('testing111'),
			when: [
				[#getPluginInfo, ['Test Event', 'event'],
					Object(#ECA, #event, name: 'Test Event', produce: Object(
						'Test Result': function (testarg) { return testarg $ '111' }))]
				])

		.mockCase('EventConditionActions', 'getActionRequired',
			args: #('Test Event', [requires: #('Test Result')],
				#(field: 'testing_field'))
			expect: #('testing_field'),
			when: [
				[#getPluginInfo, ['Test Event', 'event'],
					Object(#ECA, #event, name: 'Test Event', produce: Object(
						'Test Result': 'Prompt'))]
				])
		}

	Test_satisfied?()
		{
		.mockCase('EventConditionActions', 'satisfied?',
			args: #('Test Event', [], #()), expect: true)

		.mockCase('EventConditionActions', 'satisfied?',
			args: #('Test Event',
				[[condition_source: 'Test Result',
					condition_op: 'equals', condition_value: '30']],
				#(testarg: 40)),
			expect: false
			when: [
				[#getPluginInfo, ['Test Result', 'info_type'],
					#(ECA, info_type, name: 'Test Result', type: 'number')],
				[#getDetail, ['Test Event', 'Test Result', #(testarg: 40)], 40],
				])

		.mockCase('EventConditionActions', 'satisfied?',
			args: #('Test Event',
				[[condition_source: 'Test Result',
					condition_op: 'empty', condition_value: '']],
				#(testarg: 40)),
			expect: false
			when: [
				[#getPluginInfo, ['Test Result', 'info_type'],
					#(ECA, info_type, name: 'Test Result', type: 'number')],
				[#getDetail, ['Test Event', 'Test Result', #(testarg: 40)], 40],
				])

		.mockCase('EventConditionActions', 'satisfied?',
			args: #('Test Event',
				[[condition_source: 'Test Result',
					condition_op: 'greater than', condition_value: '30']],
				#(testarg: 40)),
			expect: true
			when: [
				[#getPluginInfo, ['Test Result', 'info_type'],
					#(ECA, info_type, name: 'Test Result', type: 'number')],
				[#getDetail, ['Test Event', 'Test Result', #(testarg: 40)], 40],
				])

		.mockCase('EventConditionActions', 'satisfied?',
			args: #('Test Event',
				[[condition_source: 'Test Result',
					condition_op: 'contains', condition_value: '9']],
				#(testarg: '40')),
			expect: false
			when: [
				[#getPluginInfo, ['Test Result', 'info_type'],
					#(ECA, info_type, name: 'Test Result', type: 'string')],
				[#getDetail, ['Test Event', 'Test Result', #(testarg: '40')], '40'],
				])
		}

	Test_source_satisfied?()
		{
		table = .MakeTable('(abc, efg, rec_field) key ()')
		.mockCase('EventConditionActions', 'satisfied?',
			args: #('Test Event',
				[[condition_source: 'Test Result',
					condition_field: 'rec_field'
					condition_op: 'contains', condition_value: 'ZZZZZZ']],
				#(testarg: [rec_field: 'testing record'])),
			expect: false
			when: [
				[#getPluginInfo, ['Test Result', 'info_type'],
					Object(#ECA, #info_type, name: 'Test Result', source: table)],
				[#getDetail, ['Test Event', 'Test Result',
					#(testarg: [rec_field: 'testing record'])],
					[rec_field: 'testing record']],
				])

		.mockCase('EventConditionActions', 'satisfied?',
			args: #('Test Event',
				[[condition_source: 'Test Result',
					condition_field: 'rec_field'
					condition_op: 'contains', condition_value: 'ZZZZZZ']],
				#(testarg: [rec_field: 'testing record'])),
			expect: false
			when: [
				[#getPluginInfo, ['Test Result', 'info_type'],
					Object(#ECA, #info_type, name: 'Test Result', source: table)],
				[#getDetail, ['Test Event', 'Test Result',
					#(testarg: [rec_field: 'testing record'])],
					false],
				])

		.mockCase('EventConditionActions', 'satisfied?',
			args: #('Test Event',
				[[condition_source: 'Test Result',
					condition_field: 'rec_field'
					condition_op: 'contains', condition_value: 'rec']],
				#(testarg: [rec_field: 'testing record'])),
			expect: true
			when: [
				[#getPluginInfo, ['Test Result', 'info_type'],
					Object(#ECA, #info_type, name: 'Test Result', source: table)],
				[#getDetail, ['Test Event', 'Test Result',
					#(testarg: [rec_field: 'testing record'])],
					[rec_field: 'testing record']],
				])
		}


	Test_getDetail()
		{
		.mockCase('EventConditionActions', 'getDetail',
			args: #('Test Event', 'Test Result', #(testarg: 'testing args')),
			expect: 'testing args 111'
			when: [
				[#getPluginInfo, ['Test Event', 'event'],
				Object(#ECA, #event, name: 'Test Event', produce: Object(
					'Test Result': function (testarg) { return testarg $ ' 111' }))]])
		}

	Test_extraFields()
		{
		func = EventConditionActions.EventConditionActions_extraFields
		rec = []
		Assert(func(rec, '') is: '')

		rec = [extends: #(field1, field2)]
		Assert(func(rec, '') is: '')
		Assert(func(rec, 'extends') is: ' extend field1, field2')

		rec = [extends: #(field1, field2), renames: #('field3 to field4')]
		Assert(func(rec, 'extends') is: ' extend field1, field2')
		Assert(func(rec, 'renames') is: ' rename field3 to field4')
		Assert(func(rec, 'not proper type') is: '')
		}

	mockCase(className, testMethod, args, expect = 'no return', when = #())
		{
		mock = Mock(EventConditionActions)
		for w in when
			{
			method = w[0]
			whenArgs = w[1]
			result = w[2]
			if not method.Capitalized?()
				method = className $ '_' $ method
			mock.When[method](@whenArgs).Return(result)
			}
		mock.When.EventConditionActions_buildCondition([anyArgs:]).CallThrough()
		if not testMethod.Capitalized?()
			testMethod = className $ '_' $ testMethod
		testMethod = Global(className)[testMethod]
		args = args.Copy().Add(testMethod, at: 0)

		if expect isnt 'no return'
			Assert(mock.Eval(@args) is: expect)
		else
			mock.Eval(@args)
		}

	Test_applyForeignData()
		{
		rec = []
		expectedField = 'random_field'
		func = EventConditionActions.EventConditionActions_applyForeignData
		func(rec, expectedField)
		Assert(rec is: [])
		rec = [no_foreignTable: 'just this']
		func(rec, expectedField)
		Assert(rec is: [no_foreignTable: 'just this'])

		spy = .SpyOn(FindForeignRecWithAbbrevNameOrNum)
		spy.Return([othertable_name: 'Bob'])

		rec = [mastertable_job: 'a thing']
		func(rec, expectedField)
		Assert(rec is: [mastertable_job: 'a thing', random_field: ''])

		rec = [mastertable_job: 'a thing']
		expectedField = 'othertable_name'
		func(rec, expectedField)
		Assert(rec is: [mastertable_job: 'a thing', othertable_name: 'Bob'])

		rec = [mastertable_job: 'a thing']
		expectedField = 'othertable_name_renamed'
		func(rec, expectedField)
		Assert(rec is: [mastertable_job: 'a thing', othertable_name_renamed: 'Bob'])

		spy.Close()
		spy = .SpyOn(FindForeignRecWithAbbrevNameOrNum)
		spy.Return(false)
		rec = [mastertable_job: 'a thing', othertable_num: 'NotLinkedButInDB']
		expectedField = 'othertable_name'
		func(rec, expectedField)
		Assert(rec is: [mastertable_job: 'a thing', othertable_num: 'NotLinkedButInDB',
			othertable_name: 'NotLinkedButInDB'])
		}

	Test_populateEventActions()
		{
		// simple mechanism on event name match lets us control whether conditions
		// are satisfied without having to set up actual conditions that pass/fail.
		// Could replace this with actual conditions down the road if we want to,
		// but the condition checking is already tested by Test_satisfied?
		testCl = EventConditionActions
			{
			EventConditionActions_satisfied?(
				eventName, conditions /*unused*/, details /*unused*/)
				{
				return eventName.Has?('SATISFIED')
				}
			}
		func = testCl.EventConditionActions_populateEventActions

		.MakeEventConditionActions('Test', #(),
			#((action_name: 'Send Email')
				(action_name: 'Notify User')))
		.MakeEventConditionActions('Test SATISFIED', #(),
			#((action_name: 'Send Email')
				(action_name: 'Notify User')
				(action_name: 'Random Action')))

		// non-existent event name, shouldn't match to anything, nothing output
		eventActions = .WatchTable('event_actions')
		func('TEST_NON_EXISTENT', #(), false)
		actions = .GetWatchTable(eventActions)
		Assert(actions is: #())

		// event name exists but the conditions are not satisfied, nothing output
		eventActions = .WatchTable('event_actions')
		func('Test', #(), false)
		actions = .GetWatchTable(eventActions)
		Assert(actions is: #())

		// event exists and conditions satisfield, so actions are created
		eventActions = .WatchTable('event_actions')
		func('Test SATISFIED', #(), false)
		actions = .GetWatchTable(eventActions)
		Assert(actions isSize: 3)
		Assert(actions[0].event_name is: 'Test SATISFIED')
		Assert(actions[0].action is: 'Send Email')
		Assert(actions[1].event_name is: 'Test SATISFIED')
		Assert(actions[1].action is: 'Notify User')
		Assert(actions[2].event_name is: 'Test SATISFIED')
		Assert(actions[2].action is: 'Random Action')

		// same as previous but supplying our own transaction
		Transaction(update:)
			{ |t|
			eventActions = .WatchTable('event_actions')
			func('Test SATISFIED', #(), t)
			actions = .GetWatchTable(eventActions)
			Assert(actions isSize: 3)
			Assert(actions[0].event_name is: 'Test SATISFIED')
			Assert(actions[0].action is: 'Send Email')
			Assert(actions[1].event_name is: 'Test SATISFIED')
			Assert(actions[1].action is: 'Notify User')
			Assert(actions[2].event_name is: 'Test SATISFIED')
			Assert(actions[2].action is: 'Random Action')
			}
		}
	Teardown()
		{
		super.Teardown()
		if not .actionsTableExisted?
			try Database('drop event_actions')
		if not .conditionsTableExisted?
			try Database('drop event_condition_actions')
		}
	}