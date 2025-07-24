// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: RadioGroups
	Left: 18
	args: () // default in class for tests that destroy the ctrl right away
	New(@args)
		// pre: args is list of controls with label: member
		{
		super(.controls(args))
		.args = args
		.Picked(.args[0].label)
		.Send("Data")
		// ensure non-active groups are disabled initially
		.Defer(.disableNonActiveGroups)
		}
	controls(args)
		{
		dir = args.Member?('horz') and args.horz is true ? 'Horz' : 'Vert'
		controls = Object(dir)
		for i in args.Members(list:)
			controls.Add(Object('Vert',
				Object('RadioButton', args[i].label),
				Object('Horz', #(Skip 18), args[i]),
				name: "Rgc" $ i))
		controls.name = 'Group'
		controls[1][1].first = true
		.value = args[0].label
		return controls
		}
	Picked(name)
		{
		.value = name
		.Send("NewValue", name)
		for i in .args.Members(list:)
			{
			ctrl = .Group['Rgc' $ i]
			selected? = .args[i].label is .value
			ctrl.RadioButton.Set(selected?)
			ctrl.Horz.SetEnabled(selected?)
			}
		}
	Get()
		{
		return .value
		}
	Set(.value)
		{
		for i in .args.Members(list:)
			{
			ctrl = .Group['Rgc' $ i]
			selected? = .args[i].label is .value
			ctrl.RadioButton.Set(selected?)
			if .readOnly isnt true
				ctrl.Horz.SetEnabled(selected?)
			}
		}
	SetEnabled(enabled)
		{
		// handle in the default manner, but then have to ensure
		// that non-active groups are disabled
		super.SetEnabled(enabled)
		if enabled is true
			.disableNonActiveGroups()
		}
	disableNonActiveGroups()
		{
		for i in .args.Members(list:)
			if .args[i].label isnt .value
				.Group['Rgc' $ i].Horz.SetEnabled(false)
		}
	readOnly: false
	SetReadOnly(.readonly)
		{
		.SetEnabled(not readonly)
		}
	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}
