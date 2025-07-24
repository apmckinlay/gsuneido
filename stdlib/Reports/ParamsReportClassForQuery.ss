// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	GetParamsWhere(@fields)
		{
		return .Eval(@Object(Report.GetParamsWhere).Append(fields))
		}
	GetUserDefinedLayout(report)
		{
		return .Eval(Report.GetUserDefinedLayout, report)
		}
	}