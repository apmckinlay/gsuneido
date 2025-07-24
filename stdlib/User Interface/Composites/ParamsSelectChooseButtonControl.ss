// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
ChooseButtonControl
	{
	New(text, list, width = false, emptyValue = 'all')
		{
		super(text, list, width)
		.emptyValue = emptyValue
		}
	Set(value)
		{
		super.Set(value is '' ? .emptyValue : value)
		.ChooseButton.Grayed?(value is "" or value is .emptyValue)
		}
	Get()
		{
		value = super.Get()
		return value is .emptyValue ? '' : value
		}
	}