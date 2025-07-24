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
		return QueryEmpty?("columns", :table, column: "lib_committed")
		}
	}