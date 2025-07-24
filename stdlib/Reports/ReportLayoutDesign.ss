// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
class
	{
	query: report_layout_designs

	EnsureTable()
		{
		Database('ensure report_layout_designs
			(rptdesign_num, report, rptdesign_name, rptdesign_layout)
			key (rptdesign_num)
			key (report, rptdesign_name)')
		}

	Customizable?(rptName)
		{
		return GetContributions('ReportLayoutDesign').Member?(rptName)
		}

	CustomizationTab(rptName, embedded = false)
		{
		if not .Customizable?(rptName)
			throw rptName $ ' is not set up to be customizable'
		layout = Object(#ReportLayoutDesign, rptName)
		if not embedded
			layout.Tab = .tab(rptName)
		return layout
		}

	GetLayout(rptName)
		{
		return Global(GetContributions('ReportLayoutDesign')[rptName].design)
		}

	GetProtect(rptName)
		{
		cont = GetContributions('ReportLayoutDesign')[rptName]
		return cont.Member?(#protect)
			? cont.protect
			: false
		}

	GetValidField(rptName)
		{
		cont = GetContributions('ReportLayoutDesign')[rptName]
		return cont.GetDefault(#validField, false)
		}

	Access(rptName)
		{
		return GetContributions('ReportLayoutDesign')[rptName].access
		}

	DefaultValue(reportName, previousValue)
		{
		return QueryCount(query = .GetQuery(reportName, withSort?:)) is 1
			? QueryFirst(query).rptdesign_name
			: previousValue
		}

	GetQuery(rptName, withSort? = false)
		{
		return .query $ ' where ' $ 'report is "' $ .baseReport(rptName) $ '"' $
			(withSort? ? ' sort rptdesign_name' : '')
		}

	tab(rptName)
		{
		return GetContributions('ReportLayoutDesign')[rptName].tab
		}

	// If a Layout shares a report with another layout, the member "report" can be
	// specified in the contribution to
	baseReport(rptName)
		{
		return GetContributions('ReportLayoutDesign')[rptName].GetDefault('report',
			rptName)
		}

	GetDefaultLayout(rptName)
		{
		return .GetLayout(rptName).SetupDefault()
		}

	OutputDefault(rptName)
		{
		.EnsureTable()
		layoutCtrl = .GetLayout(rptName)
		QueryOutputIfNew('report_layout_designs', [
			report: rptName,
			rptdesign_name: 'Default',
			rptdesign_num: Timestamp(),
			rptdesign_layout: layoutCtrl.SetupDefault()])
		}
	}
