// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(report, destinationFile, params = false,
		orientation = 'Portrait', type = 'pdf', maxPages = false)
		{
		result = ReportStatus.NODATA
		try
			{
			result = this['GenerateReport' $ type.Upper()](report,
				destinationFile, params, orientation, :maxPages)
			QueryApply1('params', :report)
				{
				it.params.lastRan = Timestamp()
				it.Update()
				}
			}
		catch (err)
			{
			if err isnt .InvalidOptionsMsg and
				not err.Prefix?(Report.MaxPagesMessagePrefix)
				SuneidoLog('ERROR: (CAUGHT) CreateReportFile - ' $ err, calls:,
					caughtMsg: 'unattended; error returned; report status/notify updated')
			return err
			}
		return result
		}

	// public only so method can be called on 'this'
	GenerateReportPDF(report, destinationFile, params, orientation, maxPages = false)
		{
		if not #(Portrait Landscape).Has?(orientation)
			throw 'Invalid report orientation: ' $ orientation

		reportOb = .buildReportObject(report, params, orientation)
		return Report(@reportOb).PrintPDF(destinationFile, quiet?:, :maxPages)
		}

	// public only so method can be called on 'this'
	GenerateReportCSV(report, destinationFile, params, orientation)
		{
		reportOb = .buildReportObject(report, params, orientation)
		return Report(@reportOb).ExportCSV(destinationFile, quiet?:)
		}

	buildReportObject(report, params, orientation)
		{
		reportOb = report.Prefix?('Reporter')
			? .buildReporterObject(report, params)
			: .buildSystemReportObject(report, params)
		reportOb.paramsdata.ReportDestination = 'pdf'
		reportOb.default_orientation = orientation
		return reportOb
		}

	buildReporterObject(report, params)
		{
		rpt = ReporterModel(report)
		.checkReportOptions(rpt, params)
		params = .getConvertedParamsDateCodes(params, rpt)

		paramsOb = Object()
		if params isnt false
			for mem, val in params
				paramsOb[mem] = Object?(val) ? val.Copy() : val

		reportOb = rpt.BuildReport(rpt.GetData(), paramsOb)
		reportOb.printParams.MergeUnion(rpt.Menu_print_params())
		reportOb.paramsdata.MergeNew(paramsOb)
		return reportOb
		}

	buildSystemReportObject(report, params)
		{
		standardReports = GetStandardScheduleReportsOb()
		if not standardReports.Member?(report)
			throw "standard report not found: " $ report
		params = .getConvertedParamsDateCodes(params)
		rptDef = standardReports[report]
		reportOb = Object()
		reportOb.paramsdata = params isnt false ? params : []
		reportOb[0] = rptDef[1]
		if rptDef.Member?('printParams')
			reportOb.printParams = rptDef.printParams
		ctrl = GetControlClass.FromControl(rptDef)
		if ctrl isnt false and ctrl.Base?(Params)
			ctrl.SetExtraParamsData(reportOb.paramsdata)

		reportOb.title = rptDef.Member?('title') ? rptDef.title : report
		if rptDef.Member?('header')
			reportOb.header = rptDef.header
		return reportOb
		}

	InvalidOptionsMsg: 'Invalid Report Options'
	checkReportOptions(rpt, params)
		{
		if not Object?(params)
			return
		sf = rpt.GetSelectFields()

		fields = sf.Fields.Values().Map({ it.RemoveSuffix('_param').RemoveSuffix('?') })
		nums = fields.Filter({ it.Has?('_abbrev') }).Map({ it.Replace('_abbrev', '_num')})
		fields.Append(nums)
		paramFields = params.Members().Map({ it.RemoveSuffix('_param') })

		menuParams = rpt.Menu_params_fields()
		missing = menuParams.Difference(paramFields)
		for m in missing
			params.Add(#(operation: "", value: '', value2: '') at: m $ "_param")

		remainingParams = paramFields.Difference(fields)
		if not remainingParams.Empty?()
			{
			converted = sf.GetConverted().Values()
			if not remainingParams.Difference(converted).Empty?()
				throw .InvalidOptionsMsg
			}
		}

	getConvertedParamsDateCodes(params, reporterModel = false)
		{
		return GetConvertedParamsDateCodes(params, reporterModel)
		}
	}
