// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	// Must be named like this in order to be found by CustomizableFieldDialog
	Name:'peditor'

	New()
		{
		.format = .FindControl('format')
		.tooltip = .FindControl('tooltip')
		}
	Controls:	(Vert
					Skip
					(Static 'Please select date format')
					Skip
					(RadioButtons  'Short' 'Long' name: 'format')
					Skip
					(Static 'Tooltip')
					Skip
					(Field name: 'tooltip')
					Skip
				)
	Valid?()
		{
		return true
		}

	// should return an object with options i.e. (list:('a' 'b') width:50 )
	Get()
		{
		format = .format.Get()
		status = .tooltip.Get()
		return Object(control: Object(:status) format: Object(long?: format is 'Long'))
		}
	Set(object)
		{
		.format.Set(object.GetDefault('Format_long?', false) is true ? 'Long' : 'Short')
		.tooltip.Set(object.GetDefault('Control_status', ''))
		}
	}
