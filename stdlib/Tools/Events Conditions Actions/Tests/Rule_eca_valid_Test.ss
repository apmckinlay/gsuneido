// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		rec = [eca_actions: '',
			eca_conditions: #([condition_field: "date", condition_op: "equals"])]
		Assert(rec.eca_valid is: 'Please select at least one action.')

		rec.eca_actions = #()
		Assert(rec.eca_valid is: 'Please select at least one action.')

		rec.eca_actions = #(
			[action_setting: #(
				#(value: "a@a.c", field: "action_email_from"),
				#(value: "", field: "action_email_to"))])
		Assert(rec.eca_valid is: 'Please select at least one action.')
		
		rec.eca_actions = #(
			[action_setting: #(
				#(value: "a@a.c", field: "action_email_from"),
				#(value: "", field: "action_email_to")),
			action_name: "Test"])
		rec.eca_conditions = Object(
			[condition_field: "string", condition_op: "contains", condition_value: "(ab"])
		Assert(rec.eca_valid is: '')
		
		rec.eca_conditions.Add(
			[condition_field: "string", condition_op: "matches", condition_value: `\(ab`])
		rec.Invalidate(#eca_valid)
		Assert(rec.eca_valid is: '')
		
		rec.eca_conditions.Add(
			[condition_field: "string", condition_op: "matches", condition_value: "(ab"])
		rec.Invalidate(#eca_valid)
		Assert(rec.eca_valid
			is: Display("(ab") $ " is not valid for " $ Display("matches"))
		}
	}