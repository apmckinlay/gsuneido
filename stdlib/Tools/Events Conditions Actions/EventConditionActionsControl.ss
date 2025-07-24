// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New()
		{
		super(.layout())
		.data = .FindControl("ExplorerListView").GetView()
		}

	layout()
		{
		ctrls = Object('Vert')
		ctrls.Add(#(CenterTitle 'Business Triggers'))

		msg = EventConditionActions.PermissionToChange()
		readonly = false
		if Object?(msg) and msg.GetDefault('reason', '') isnt ''
			{
			readonly = true
			ctrls.Add(Object('Heading2', msg.reason), 'Skip')
			}
		ctrls.Add(Object('ExplorerListView'
			#(ExplorerListModel "event_condition_actions
				rename eca_num to eca_num_new
				sort eca_num_new"
				(eca_num_new))
			#(Vert
				(Heading2 'Event:')
					(Skip 4)
					(Horz Skip eca_event) Skip
				(Heading2 'Conditions:')
					(Horz Skip (ChooseConditions name: eca_conditions)) Skip
				(Heading2 'Actions:')
					(Skip 4)
					(Horz (ChooseActions name: eca_actions))
				)
			columns: Object('eca_event'),
			:readonly,
			columnsSaveName: 'event_condition_actions'
			protectField: 'eca_protect'
			validField: 'eca_valid'
			buttonBar: true))
		return ctrls
		}

	ChooseConditions_SourceFields()
		{
		return EventConditionActions.GetEventSourceFields(.data.Get().eca_event)
		}

	ChooseActions_ActionList()
		{
		return EventConditionActions.Actions(.data.Get().eca_event)
		}

	ExplorerListView_RecordChanged(member, data /*unused*/)
		{
		if member is 'eca_event'
			{
			.FindControl('eca_conditions').SetChooseFieldList()
			.FindControl('eca_actions').SetActionList()
			}
		}

	ExplorerListView_EntryLoaded()
		{
		.data.Dirty?(false)
		}

	Destroy()
		{
		ServerEval('Has_ECA_Event.ResetCache')
		super.Destroy()
		}
	}
