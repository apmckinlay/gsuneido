// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
// This is a base class, derived classes must define .Tables()
MemoizeSingle
	{
	Func()
		{
		return .Tables().Filter(.disabled)
		}

	disabled(table)
		{
		cols = QueryList("columns where table is " $ Display(table), 'column')
		return not cols.Intersects?(SvcTable.SvcColumns)
		}
	}