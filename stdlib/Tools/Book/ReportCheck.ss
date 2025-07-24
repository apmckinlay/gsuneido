// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
PageBaseCheck
	{
	CallClass(book)
		{
		QueryApply(book $ " where path.Has?('Reporter Reports') and
			name isnt 'Reporter'")
			{ |x|
			.CheckReport(x.text.SafeEval(), x.name)
			}

		QueryApply(book $ " where path.Has?('Reporter Forms') and
			name isnt 'Reporter Form'")
			{ |x|
			.CheckReport(x.text.SafeEval(), x.name)
			}

		.ForeachBookOption(book)
			{ |rptClass, name|
			.CheckReport(rptClass, name)
			}
		}

	CheckReport(rptClass, name)
		{
		.runReport(rptClass, name)
			{ |paramsCtrl|
			.testPreview(paramsCtrl)

			Params.SetNullPdc()
			paramsCtrl.PrintReport()

			paramsCtrl.SavePDF('rpt/' $ name $ '.pdf')

			if false isnt paramsCtrl.FindControl('Export')
				paramsCtrl.SaveCSV('rpt/' $ name $ '.csv')
			}
		}

	testPreview(paramsCtrl)
		{
		paramsCtrl.DisablePreviewDialog()
		previewWnd = paramsCtrl.On_Preview()
		if Instance?(previewWnd)
			{
			if previewWnd.Member?('Ctrl')
				previewWnd.Ctrl.On_Last()
			previewWnd.Destroy()
			}
		}

	defaultSetup: #(
		setup: function () { }
		teardown: function () { }
		)
	runReport(rptClass, name, block)
		{
		if not .isReport?(rptClass)
			return

		setup = GetContributions('ReportCheckSetups').GetDefault(name, .defaultSetup)
		(setup.setup)()

		Finally({
			rptWnd = Window(rptClass, show: false)
			paramsCtrl = rptWnd.Ctrl
			if not paramsCtrl.Base?(ParamsControl)
				return

			EnsureDir('rpt')

			block(paramsCtrl)

			rptWnd = rptWnd.Destroy()
			}, {
			(setup.teardown)()
			})
		}

	isReport?(rptClass)
		{
		if Object?(rptClass)
			return rptClass.GetDefault(0, false) is 'Params'
		if Class?(rptClass)
			return rptClass.Base?(QueryFormat)
		return false
		}
	}