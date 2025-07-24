// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	TopLeft_Extra: (Skip 0)
	StandardReportPrefix: false
	ReporterInvalidMsg: ''
	ReporterBook: ''

	HandleReportMenu(@unused)
		{ return true }

	HasPermission?(x)
		{ return Customizable.AccessToDataSource?(x.params.Source, defaultRtn:) }

	AfterDeleteReport(@unused)
		{ }

	Protect(unused)
		{ return false }

	GetPath(@unused)
		{
		return ''
		}
	}
