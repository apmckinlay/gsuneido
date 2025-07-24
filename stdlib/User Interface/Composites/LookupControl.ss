// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
ChooseListControl
	{
	New(query, field)
		{
		super(.makelist(query, field))
		}

	makelist(query, field)
		{
		list = Object()
		QueryApply(query)
			{|x|
			list.Add(x[field])
			}
		return list
		}
	}
