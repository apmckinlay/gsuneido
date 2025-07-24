// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	New(valueField, overrideField,
		valid = function (value/*unused*/) { return '' })
		{
		super(.controls(valueField, overrideField))
		.valuefield = .Horz.GetChildren()[0]
		.checkbox = .Horz.GetChildren()[2]
		.valid = valid
		.Top = .Horz.Top
		}
	controls(valueField, overrideField)
		{
		return Object('Horz'
			Object('Number' mask: "-###,###,###.##", readonly:, width: 9,
				name: valueField)
			'Skip'
			Object('CheckBox', name: overrideField, top: 11)
			#(Skip 4)
			#(Static Override))
		}
	NewValue(value, source)
		{
		if (source is .checkbox)
			{
			if (value is true)
				{
				n = Ask('Enter override value', 'Override', .Window.Hwnd,
					#(Number mask: "-###,###,###.##", width: 9), .valid)
				if (n is false)
					{
					.checkbox.Set(false)
					return
					}
				.valuefield.Set(n)
				.Send('Override', n)
				}
			.Controller.Msg(Object('NewValue', value, :source))
			}
		}
	}