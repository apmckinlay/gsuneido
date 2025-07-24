// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Mock
	{
	New(columns = false)
		{
		super(ListControl)
		.data = Object()
		.columns = columns is false ? Object() : columns
		}

	AddRow(row)
		{
		.data.Add(row)
		}

	AddRecord(record)
		{
		.data.Add(record)
		}

	GetRow(row)
		{
		.data[row]
		}

	Get()
		{
		return .data
		}

	GetAllBrowseData()
		{
		return .data
		}

	Set(ob)
		{
		.data = ob
		}

	GetCol(col)
		{
		return .columns.GetDefault(col, false)
		}

	FindRowIdx(field, value)
		{
		return .data.FindIf({ it[field] is value })
		}

	GetNumRows()
		{
		return .data.Size()
		}

	RepaintRow(@unused) {}
	}
