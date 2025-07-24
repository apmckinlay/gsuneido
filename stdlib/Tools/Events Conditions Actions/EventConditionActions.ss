// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Events()
		{
		return Plugins().Contributions('ECA', 'event').Map({ it.name })
		}

	Actions(eventName)
		{
		if false is eventInfo = .getPluginInfo(eventName, 'event')
			return #()
		infoTypes = eventInfo.GetDefault('produce', #()).Members()
		actions = Object()
		Plugins().ForeachContribution('ECA', 'action')
			{ |c|
			if c.GetDefault('requires', #()).Subset?(infoTypes) and
				not c.GetDefault('excludeEvents', #()).Has?(eventName)
				actions.Add(c.name)
			}
		return actions
		}

	PermissionToChange()
		{
		Plugins().ForeachContribution('ECA', 'permission')
			{ |c|
			if '' isnt msg = (c.allowToChange)()
				return msg
			}
		return ''
		}

	EnsureTables()
		{
		Database('ensure event_condition_actions (eca_num, eca_event,
			eca_conditions, eca_actions) key (eca_num)')

		Database('ensure event_actions (event_num, event_name,
			event_details, action, action_setting) key (event_num)')
		}

	DisabledOption: 'PopulateEventDisabled' // used by test
	PopulateEvent(name, details = false, t = false)
		{
		if ServerSuneido.Get(.DisabledOption) is true
			return true
		try
			.populateEventActions(name, details, t)
		catch(err)
			SuneidoLog('ERROR: ' $ err)
		return true
		}

	populateEventActions(name, details, t)
		{
		DoWithTran(t, update:)
			{|t|
			t.QueryApply('event_condition_actions where eca_event is ' $ Display(name))
				{ |eca|
				if .satisfied?(name, eca.eca_conditions, details)
					for (action in eca.eca_actions)
						t.QueryOutput('event_actions', [event_num: Timestamp(),
							event_name: name, event_details: details,
							action: action.action_name,
							action_setting: action.GetDefault('action_setting', #())])
				}
			}
		}

	satisfied?(eventName, conditions, details)
		{
		for condition in conditions
			{
			conditionInfo = .getPluginInfo(condition.condition_source, 'info_type')
			if conditionInfo.Member?('type') // type
				{
				value = .getDetail(eventName, conditionInfo.name, details)
				encode = Datadict(conditionInfo.type, #(Encode)).Encode
				}
			else // source
				{
				conditionField = condition.condition_field
				if false is rec = .getDetail(eventName, conditionInfo.name, details)
					return false
				if false is rec.Member?(conditionField) or rec[conditionField] is ''
					.applyForeignData(rec, conditionField)
				value = rec[conditionField]
				encode = Datadict(conditionField, #(Encode)).Encode
				}

			conditionCode = .buildCondition(condition, value, encode)
			try
				{
				if false is conditionCode.Eval() // Eval should be okay here
					return false
				}
			catch(err /*unused*/)
				{
				SuneidoLog('ERROR: EventConditionAction condition check failed: ' $
					conditionCode, params: [:eventName, :conditions, :details])
				return false
				}
			}
		return true
		}

	buildCondition(condition, value, encode)
		{
		pos = Select2.Ops.FindIf({|x| x[0] is condition.condition_op })
		op = Select2.Ops[pos]
		operator = op[1]

		// TODO fix Select2 instead of tweaking here
		if operator in ('=', '==')
			operator = 'is'
		if operator is '!='
			operator = 'isnt'

		return Display(value) $ " " $ operator $ " " $
			Display( op[0].Suffix?('empty') // does not require value
				? ""
				: op[1].Suffix?('~') // if string operation
					? op.pre $ condition.condition_value $ op.suf
					: encode(condition.condition_value))
		}

	applyForeignData(rec, expectedField)
		{
		if false is foreignRec = FindForeignRecWithAbbrevNameOrNum(rec,
			expectedField, useNum?:)
			{
			if expectedField.Has?('name') or expectedField.Has?('abbrev')
				{
				numVal = rec[expectedField.Replace('name|abbrev', 'num')]
				if numVal isnt '' and false is Date?(numVal)
					rec[expectedField] = numVal
				}
			return
			}

		value = foreignRec[expectedField]
		if value is ''
			value = foreignRec[expectedField.BeforeLast('_')]
		rec[expectedField] = value
		}

	getDetail(eventName, infoTypeName, details)
		{
		event = .getPluginInfo(eventName, 'event')
		return (event.produce[infoTypeName])(@details)
		}

	PerformActions()
		{
		taskName = 'EventConditionActions.PerformActions'
		ScheduleNextEvent.LogTask(taskName, 'SchedulerLastProcessStarted',
			'10MinuteTasks')
		forever
			{
			if false is event = QueryFirst('event_actions sort event_num')
				return true
			try
				.execute(event.action, event.event_name, event.event_details,
					event.action_setting)
			catch(err)
				{
				SuneidoLog('ERROR: EventConditionAction execute failed on ' $
						event.action $ ' - ' $ err,
					params: [event.action, event.event_name, event.event_details,
						event.action_setting])
				}
			QueryDelete('event_actions', event)
			}
		ScheduleNextEvent.LogTask(taskName $ ' Completed',
			'SchedulerLastProcessCompleted', '10MinuteTasks')
		return true
		}

	execute(action, eventName, details, setting)
		{
		actionInfo = .getPluginInfo(action, 'action')
		args = .getActionRequired(eventName, actionInfo, details)
		for arg in setting
			args.Add(arg.value)
		Global(actionInfo.func)(@args)
		}

	getActionRequired(eventName, actionInfo, details)
		{
		args = Object()
		eventInfo = .getPluginInfo(eventName, 'event')
		for require in actionInfo.GetDefault('requires', #())
			{
			producer = eventInfo.produce[require]
			if String?(producer)
				producer = Global(producer)
			args.Add(producer(@details))
			}
		return args
		}

	getPluginInfo(name, extension)
		{
		Plugins().ForeachContribution('ECA', extension)
			{ |c|
			if c.name is name
				return c
			}
		return false
		}

	GetActionSettings(action)
		{
		if action is '' or false is actionInfo = .getPluginInfo(action, 'action')
			return #()
		return actionInfo.GetDefault('setting', #())
		}

	GetEventSourceFields(event)
		{
		if false is eventInfo = .getPluginInfo(event, 'event')
			return #()
		product = eventInfo.GetDefault('produce', #())
		availableFields = Object().Set_default(Object())
		for infoTypeName in product.Members()
			{
			condition = .getPluginInfo(infoTypeName, 'info_type')
			if false is condition.Member?('type')
				{
				columns = condition.GetDefault('columnsFunc', [])
				if Object?(columns) is false
					columns = columns()

				sf = SelectFields(QueryColumns(condition.source $
					.extraFields(columns, 'renames') $
					.extraFields(columns, 'extends')),
					columns.excludes)
				availableFields[infoTypeName] = sf.Fields
				}
			else
				availableFields[infoTypeName] = #()
			}
		return availableFields
		}

	extraFields(condition, type)
		{
		queryArg = type is 'extends'
			? ' extend '
			: ' rename '
		return Opt(queryArg, condition.GetDefault(type, #()).Join(', '))
		}
	}
