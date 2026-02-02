// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Controls: #(Horz, #('ConfigLocate', name: 'search'))

	fields: #()
	Startup()
		{
		if 0 is c = .Send('GetAccess1RecordControl')
			{
			.Horz.Remove(0)
			return
			}
		.wg = WaitGroup()
		.wg.Thread({
			try
				layout = c.GetControlLayout()
			catch (unused, '*member not found') // screen destroyed
				layout = #()
			.fields = CollectFields(layout, path?:)
			})
		}

	GetConfigFields()
		{
		.wg.Wait()
		return .fields
		}

	SearchControl(control, type, name)
		{
		for c in control.GetChildren()
			{
			match? = type in ('Field', 'Button', 'ResetButton')
				? c.Name is name
				: type is 'Static' and c.Base?(StaticControl) // also handles Heading
					? c.Get() is name
					: false
			if match? or false isnt (c = .SearchControl(c, type, name))
				return c
			}
		return false
		}

	LocateConfigField(path)
		{
		ctrl = .Send("GetRecordControl")
		.fields.Each()
			{ |f|
			if f.path is path
				{
				for i in f.Members(list:)
					{
					field = f[i]
					if field.type is 'Tab'
						{
						ctrl = ctrl.FindControl('Tabs')
						ctrl.Select(ctrl.FindTab(field.section))
						ctrl = ctrl.TabsControl_GetCurrentControl().Parent
						}
					if field.type is 'Accordion'
						{
						ctrl = ctrl.FindControl('Accordion')
						ctrl.ExpandByName(field.section)
						}
					if field.type in (#Button, #Static, #Heading, #Field, #ResetButton)
						{
						if false is target = .SearchControl(ctrl, field.type, field.name)
							continue
						if field.type is 'ResetButton'
							target = target.FindControl('Reset')
						.Send('ScrollToView', target)
						return
						}
					}
				}
			}
		}
	}
