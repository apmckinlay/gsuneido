// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
ChooseDateControl
	{
	Name: ChooseDateTime
	Layout(field /*unused*/)
		{
		return Object('Date', showTime:)
		}

	Getter_DialogControl()
		{
		return Object('DateTime', .Field.Get())
		}

	DisplayValues(control /*unused*/, vals)
		{
		return vals.Map({ DateControl.FormatValue(it, showTime:) })
		}
	}
