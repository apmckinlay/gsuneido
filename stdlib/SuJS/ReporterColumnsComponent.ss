// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: 'ReporterColumns'
	Xmin: 700

	New()
		{
		.CreateElement('table')
		.SetStyles(#(
			'table-layout': 'fixed',
			'user-select': 'none',
			'width': '100%'))
		.header = .Construct('ListHeader', .El, buttonStyle:)
		.Ymin = SuRender().GetTextMetrics(.El, 'M').height
		.SetMinSize()
		}

	UpdateHead(headCols)
		{
		.header.Update(headCols)
		}

	SetColWidth(col, width)
		{
		.header.SetColWidth(col, width)
		}
	}
