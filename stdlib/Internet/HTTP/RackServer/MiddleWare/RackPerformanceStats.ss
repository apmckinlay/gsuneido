// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
// Use this as the first middle-ware if we need to track the performance
RackComposeBase
	{
	StatsMem: 'RackPerfStats'
	Call(env)
		{
		start = Date()
		if -1 is result = .App(:env)
			return result

		bucket = (.getResTime(start).Log10() * 3.3/*=base 2 and 10 conversion*/).Ceiling()
		.Synchronized()
			{
			if not ServerSuneido.HasMember?(.StatsMem)
				ServerSuneido.Set(.StatsMem, Object().Set_default(0))
			cur = ServerSuneido.GetAt(.StatsMem, bucket, 0)
			ServerSuneido.Add(.StatsMem, ++cur, bucket)
			}
		return result
		}

	getResTime(start)
		{
		return Date().MinusSeconds(start) * 1000 /*= to millisecond conversion*/
		}
	}