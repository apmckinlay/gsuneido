// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	getter_controls()
		{
		return .controls = Object() // once only
		}
	Data(source)
		{
		.controls.Add(source)
		.Send(#Data, :source)
		}
	NewValue(value, source)
		{
		if value is false and
			.controls.Every?({ it.Get() is false })
			{
			if .controls.Size() is 2
				{
				i = .controls[0] is source ? 1 : 0
				.controls[i].Set(true)
				.Send(#NewValue, true, source: .controls[i])
				}
			else
				{
				source.Set(true)
				return
				}
			}
		else
			.Send(#NewValue, value, :source)
		}
	}