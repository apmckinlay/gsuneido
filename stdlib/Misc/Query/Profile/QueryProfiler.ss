// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Query Profiler'
	New(.query)
		{
		.mshtml = .FindControl("Mshtml")
		.metrics = .FindControl("qp_metrics")
		.NewValue(0)
		.FindControl('Editor').Set(FormatQuery(query))
		}
	Controls: (Tabs
		(Vert
			(Border (Horz qp_metrics Skip RefreshButton))
			(Mshtml)
			Tab: 'Profile')
		(QueryCodeControl Tab: 'Query')
		constructAll:)
	NewValue(unused)
		{
		metrics = .metrics.Get()
		if metrics is .prev
			return
		.prev = metrics
		.refresh()
		}
	On_Refresh()
		{
		.refresh()
		}
	prev: ()
	refresh()
		{
		metrics = .metrics.Get().Split(' | ')
		if metrics.Size() is 0
			return
		html = QueryTreeHtml(.query, metrics)
		.mshtml.Set(html)
		}
	}