// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(field, prompt, customLib)
		{
		table = field $ "_table"
		Database("ensure " $ table $ " (name, desc) key(name)")

		Transaction(update:)
			{ |t|
			maxNum = t.QueryMax(customLib, 'num', 0)
			.createAccess(t, customLib, field, table, prompt, maxNum)
			}

		return Object(control: Object(customField: field))
		}

	createAccess(t, customLib, field, table, prompt, maxNum)
		{
		t.QueryOutput(customLib, Record(
			name: 'Access_' $ field,
			text: '#(Browse, ' $ table $ ', title: ' $ Display(prompt) $ ')'
			num: maxNum + 1
			group: -1, parent: 0))
		}

	UpdateProperties(field, prompt, customLib)
		{
		access_field = 'Access_' $ field
		text = '#(Browse, '$ field $ '_table, title: ' $ Display(prompt) $')'
		QueryDo('update ' $ customLib $ ' where name = ' $ Display(access_field) $
			' set text = ' $ Display(text))
		LibUnload(access_field)
		ServerEval('LibUnload', access_field)
		return Object(control: Object(customField: field))
		}

	}
