// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
CustomizableFieldDialogPropertiesEditor_FieldControl
	{
	RangeFrom: 5
	RangeTo: 200
	DefaultWidth: 25
	heightRangeFrom: 1
	heightRangeTo: 20
	defaultHeight: 4

	New()
		{
		super(excludeTip:)
		.field_height = .FindControl('field_height')
		.tabthrough = .FindControl('tabthrough')
		}

	Controls()
		{
		ctrls = super.Controls()
		ctrls.Add(Object('Static' 'Select Text Height (from ' $
			Display(.heightRangeFrom) $ ' to ' $ Display(.heightRangeTo) $ ')'))
		ctrls.Add(#Skip)
		ctrls.Add(Object('Number' rangefrom: .heightRangeFrom, rangeto: .heightRangeTo,
			set: .defaultHeight, name: 'field_height'))
		ctrls.Add(#Skip)
		ctrls.Add(Object('CheckBox', text: 'Tab through', name:'tabthrough'))
		return ctrls
		}

	Valid?()
		{
		value = .field_height.Get()
		return super.Valid?() and Number?(value) and
			.heightRangeFrom <= value and value <= .heightRangeTo
		}

	Get()
		{
		height = .field_height.Get()
		options = super.Get()
		options.control.height = height
		options.format.height = height
		options.control.tabthrough = .tabthrough.Get()
		return options
		}

	Set(object)
		{
		super.Set(object)
		.field_height.Set(object.GetDefault('Control_height', .defaultHeight))
		.tabthrough.Set(object.GetDefault('Control_tabthrough', false))
		}
	}
