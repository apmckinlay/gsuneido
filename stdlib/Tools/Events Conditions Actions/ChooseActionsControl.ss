// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
RepeatControl
	{
	New()
		{
		super(.layout(), minusAtFront:, plusAtBottom:)
		}
	layout()
		{
		return #(Horz Skip
			(Vert
				(ChooseList, list: #(), width: 20,
					mandatory:, weight: 'bold', name: action_name)
				(Skip 5)
				(Vert name: 'setting_group')
				Skip))
		}

	Get()
		{
		actions = Object()
		for row in super.Get()
			{
			actionRec = Record()
			actionRec.action_name = row.action_name
			settings = EventConditionActions.GetActionSettings(actionRec.action_name)
			actionRec.action_setting = Object()
			for field in settings
				actionRec.action_setting.Add(Object(:field, value: row[field]))
			actions.Add(actionRec)
			}
		return actions
		}

	On_Plus(source)
		{
		super.On_Plus(source)
		row = .GetRows().Last()
		row.Set([action_setting: #(), action_name: ""])
		.SetActionList()
		}

	Set(data)
		{
		super.Set(data)
		.SetActionList()
		rows = .GetRows()
		for i in .. rows.Size()
			{
			row = rows[i]
			actionRecord = row.Get()
			settings = EventConditionActions.GetActionSettings(actionRecord.action_name)
			settingGroup = row.FindControl('setting_group')
			settingGroup.RemoveAll()
			for(fieldIndex = 0; fieldIndex < settings.Size(); ++fieldIndex)
				{
				settingGroup.Append(settings[fieldIndex])
				ctrl = settingGroup.FindControl(settings[fieldIndex])
				// Setting ctrl since RecordControl does not handle object with member values
				value = (data[i].action_setting)[fieldIndex].value
				ctrl.Set(value)
				actionRecord[settings[fieldIndex]] = value
				}
			row.SetReadOnly(.GetReadOnly())
			}
		}

	SetActionList()
		{
		actions = .Send('ChooseActions_ActionList')
		for row in .GetRows()
			{
			actionName = row.FindControl('action_name')
			actionName.SetList(actions)
			}
		}

	RepeatRecord_Changed(member)
		{
		super.RepeatRecord_Changed(member)
		if member is 'action_name'
			.Set(.Get())
		}
	}