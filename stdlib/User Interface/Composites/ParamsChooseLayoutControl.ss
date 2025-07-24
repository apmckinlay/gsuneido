// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: 'ParamsChooseLayout'
	New(reportName, xstretch = 1)
		{
		super(.layout(reportName, xstretch))
		.ReportName = reportName
		}
	layout(reportName, xstretch)
		{
		return Object('Vert',
			Object('Horz'
				#(Heading, 'Layout') 'Skip'
					Object('Key',
						query: ReportLayoutDesign.GetQuery(reportName),
						access: ReportLayoutDesign.Access(reportName),
						mandatory:,
						field: 'rptdesign_name',
						columns: #(rptdesign_name)
						keys: #(rptdesign_name),
						excludeSelect: #(report, rptdesign_num, rptdesign_layout),
						:xstretch,
						prefixColumn: 'rptdesign_name',
						name: 'params_report_layout'),
				)
			Object('StaticText', '', name: 'layoutError' textStyle: 'error')
			#Skip
			Object('StaticText', '', name: 'layoutWarning' textStyle: 'warn')
			)
		}
	}
