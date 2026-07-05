// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
CustomizableFieldDialogPropertiesEditor
	{
	RangeFrom: 5
	RangeTo: 40
	DefaultWidth: 20

	New(.excludeTip = false)
		{
		.field_width = .FindControl('field_width')
		.tooltip = .FindControl('tooltip')
		}
	GetControls()
		{
		layoutOb = Object('Vert',
			#Skip,
			Object('Static' 'Select Text Width (from ' $ Display(.RangeFrom) $
				' to ' $ Display(.RangeTo) $ ')'),
			#Skip,
			Object('Number' rangefrom: .RangeFrom, rangeto: .RangeTo,
				set: .DefaultWidth, name: 'field_width'))

		if not .excludeTip
			layoutOb.Add(#Skip,
				Object('Static' 'Tooltip'),
				#Skip,
				Object('Field' name: 'tooltip'),
				#Skip)
		return layoutOb
		}
	Valid?()
		{
		if false is super.Valid?()
			return false
		value = .field_width.Get()
		return Number?(value) and .RangeFrom <= value and value <= .RangeTo
		}

	// should return an object with options i.e. (list:('a' 'b') width:50 )
	Get()
		{
		x = super.Get()
		width = .field_width.Get()
		status = false is .tooltip ? '' : .tooltip.Get()
		x.control = Object(:width, :status)
		x.format = Object(:width)
		return x
		}
	Set(object)
		{
		super.Set(object)
		.field_width.Set(object.GetDefault('Control_width', .DefaultWidth))
		if false isnt .tooltip
			.tooltip.Set(object.GetDefault('Control_status', ''))
		}
}
