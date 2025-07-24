// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_rename_report()
		{
		.TearDownIfTablesNotExist('params')
		Params.Ensure()

		name = 'Reporter_Test'
		QueryOutput('params', Record(
			params: [report_name: name, Source: "Recurring Payable Invoices"],
			report: "Reporter - " $ name))
		reporter = ReporterOpenDialog
			{
			ReporterOpenDialog_prefix: "Reporter - "
			}

		renamed_name = 'Reporter_Test_Renamed'
		Assert(reporter.ReporterOpenDialog_renameReport(name, renamed_name) is: '')

		newname = 'Reporter_Test2'
		QueryOutput('params', Record(
			params: [report_name: newname, Source: "Recurring Payable Invoices"],
			report: "Reporter - " $ newname))
		Assert(reporter.ReporterOpenDialog_renameReport(renamed_name, newname)
			is: 'already exists')
		}

	Test_save()
		{
		// testing if HandleReportMenu failed
		mock = .mockSave(data = [report: 'reporter_test_report_design'])
		rpt_rec = Record(report: 'reporter_test_report', params: data)
		mock.When.Reporter_handleReportMenu(rpt_rec).Return(false)
		mock.Eval(Reporter.Reporter_save)
		mock.Verify.Never().Reporter_output_params(rpt_rec)
		mock.Verify.Never().Dirty?([anyArgs:])

		// testing if output_params failed
		mock = .mockSave(data = [report: 'reporter_test_report_design'])
		rpt_rec = Record(report: 'reporter_test_report', params: data)
		mock.When.Reporter_handleReportMenu(rpt_rec).Return(true)
		mock.When.Reporter_output_params(rpt_rec).Throw("interrupt")
		Assert({ mock.Eval(Reporter.Reporter_save) } throws: "interrupt")
		mock.Verify.Reporter_output_params(rpt_rec)
		mock.Verify.Never().Dirty?([anyArgs:])

		// testing if everything is ok
		mock = .mockSave(data = [report: 'reporter_test_report_design'])
		rpt_rec = Record(report: 'reporter_test_report', params: data)
		mock.When.Reporter_handleReportMenu(rpt_rec).Return(true)
		mock.When.Reporter_output_params(rpt_rec).Return(true)
		mock.Eval(Reporter.Reporter_save)
		mock.Verify.Never().Reporter_saveWarningMessage()
		mock.Verify.Reporter_output_params(rpt_rec)
		mock.Verify.Dirty?(false)
		}

	mockSave(data)
		{
		mock = Mock(Reporter)
		mock.When.Reporter_prepare_to_save().Return(data)
		mock.When.Dirty?([anyArgs:]).Return(false)
		mock.When.Reporter_isScheduledReport?([anyArgs:]).Return(false)
		mock.When.CheckScheduled([anyArgs:]).CallThrough()
		mock.Reporter_rpt = Mock()
		mock.Reporter_rpt.When.GetSaveName().Return('reporter_test_report')
		mock.Reporter_source = 'reporter_test_source'
		return mock
		}

	Test_save_scheduled()
		{
		mock = Mock(Reporter)
		data = [report: 'reporter_test_report_design']
		mock.When.Reporter_prepare_to_save().Return(data)
		mock.When.Dirty?([anyArgs:]).Return(false)
		mock.When.Reporter_isScheduledReport?([anyArgs:]).Return(true)
		mock.When.CheckScheduled([anyArgs:]).CallThrough()
		mock.When.Reporter_saveWarningMessage().Return(false, true)
		mock.Reporter_rpt = Mock()
		mock.Reporter_rpt.When.GetSaveName().Return('Reporter - reporter_test_report')
		mock.Reporter_source = 'reporter_test_source'
		rpt_rec = Record(report: 'Reporter - reporter_test_report', params: data)
		mock.When.Reporter_handleReportMenu(rpt_rec).Return(true)
		mock.When.Reporter_output_params(rpt_rec).Return(true)

		// user picked cancel
		mock.Eval(Reporter.Reporter_save)
		mock.Verify.Reporter_saveWarningMessage()
		mock.Verify.Never().Reporter_output_params(rpt_rec)
		mock.Verify.Never().Dirty?([anyArgs:])

		// user picked ok
		mock.Eval(Reporter.Reporter_save)
		mock.Verify.Times(2).Reporter_saveWarningMessage()
		mock.Verify.Reporter_output_params(rpt_rec)
		mock.Verify.Dirty?(false)
		}

	Test_validFormula_validateFormulaName()
		{
		mock = Mock(Reporter)
		mock.When.validFormula([anyArgs:]).CallThrough()
		mock.When.validateFormulaName([anyArgs:]).CallThrough()
		mock.When.validateFormulaCode([anyArgs:]).Return('')

		prompt = #TestFormula
		formula = Object(calc: prompt, formula: #notTesting, form_val: true, type: '')

		results = mock.validFormula(formula, prompt, [], [], [])
		Assert(results is: `Format is required for formula: TestFormula`)

		formula.type = #Dollar
		results = mock.validFormula(formula, prompt, [], [], [])
		Assert(results is: `Should not have both Formula and Menu Option`)

		results = mock.validFormula(formula, prompt, [prompt], [], [])
		Assert(results
			is: `Formula field name TestFormula is a duplicate. ` $
				`Please rename one of the formulas.`)

		results = mock.validFormula(formula, prompt, [], [prompt], [])
		Assert(results
			is: `Formula field name TestFormula is in use. Please rename the formula.`)

		formula.form_val = false
		formula.calc = prompt = ''
		results = mock.validFormula(formula, prompt, [], [], [])
		Assert(results is: ``)

		formula.calc = prompt = 'TestFormula'
		results = mock.validFormula(formula, prompt, [], [], [])
		Assert(results is: ``)
		}

	Test_validateFormulaCode()
		{
		testCl = Reporter
			{
			Reporter_rpt: class
				{
				GetSelectFields()
					{
					return class
						{
						PromptToField(prompt)
							{
							return "not_a_summary_func_field"
							}
						}
					}
				}
			}
		res = testCl.Reporter_validateFormulaCode("formula_text", "A Fake Prompt")
		Assert(res
			is: 'Formula field name A Fake Prompt is in use. Please rename the formula.')
		}

	Test_buildHistoryString()
		{
		buildHisStr = Reporter.Reporter_buildHistoryString
		Assert(buildHisStr('created', []) is: '')
		// stored users prior to this, will never not have user so no history without it
		data = [created_on: #20220408]
		Assert(buildHisStr('created', data)	is: '')
		// lots of reports with only users creating them, so at least show that info
		data = [created_by: 'Joe ReportMaker']
		Assert(buildHisStr('created', data)
			is: 'Created by Joe ReportMaker\r\n')

		data.created_on = #20220408
		data.last_modified_on = #20220607
		data.last_modified_by = 'Blake'
		Assert(buildHisStr('created', data)
			is: 'Created on ' $ #20220408.ShortDate() $ ' by Joe ReportMaker\r\n')
		Assert(buildHisStr('last_modified', data)
			is: 'Last Modified on ' $ #20220607.ShortDate() $ ' by Blake\r\n')

		data = [created_by: 'Blake']
		allStringOldData = buildHisStr('created', data) $
			buildHisStr('last_modified', data) $ buildHisStr('last_ran', data)
		Assert(allStringOldData is: 'Created by Blake\r\n')

		data = [created_by: 'Blake', last_ran_on: #20220601, last_ran_by: 'Someone else']
		allStringPartialData = buildHisStr('created', data) $
			buildHisStr('last_modified', data) $ buildHisStr('last_ran', data)
		Assert(allStringPartialData is: 'Created by Blake\r\nLast Ran on ' $
			#20220601.ShortDate() $ ' by Someone else\r\n')
		}

	Test_save?()
		{
		mock = Mock(Reporter)
		// Methods being suppressed
		mock.When.AlertInfo([anyArgs:]).Do({ })
		mock.When.SetSaveName([anyArgs:]).Do({ })
		mock.When.On_Save_As([anyArgs:]).Return(false)
		mock.When.stdReportPrefix([anyArgs:]).Return(false)
		// Return values to be manipulated:
		mock.When.OverwriteReport?([anyArgs:]).Return(ID.NO)
		// Methods being tested:
		mock.When.save?([anyArgs:]).CallThrough()
		mock.When.CheckName([anyArgs:]).CallThrough()
		mock.When.saveStandardReport([anyArgs:]).CallThrough()

		// Name is '', On_Save_As is called
		Assert(.testSave(mock, '', '') is: false)
		mock.Verify.On_Save_As()
		mock.Verify.Never().CheckName([anyArgs:])
		mock.Verify.Never().saveStandardReport([anyArgs:])

		// Name is false, CheckName is called and immediately returns
		Assert(.testSave(mock, false, '') is: false)
		mock.Verify.CheckName(false, '')
		mock.Verify.Never().AlertInfo([anyArgs:])
		mock.Verify.Never().saveStandardReport([anyArgs:])

		// Name contains non alpha-numeric characters, user is informed
		name = 'My^Report1'
		stdReportPrefix = 'TestPrefix'
		Assert(.testSave(mock, name, stdReportPrefix) is: false)
		mock.Verify.AlertInfo('Save As',
			'You can only use alpha-numeric characters (letters and numbers) ' $
			'and spaces in the name.\n\nPlease save as a different name.')
		mock.Verify.Never().saveStandardReport([anyArgs:])

		// Name is prefixed by the system standard prefix, user is informed
		name = 'TestPrefix My Report1'
		Assert(.testSave(mock, name, stdReportPrefix) is: false)
		mock.Verify.AlertInfo('Save As',
			"Reporter name cannot start with 'TestPrefix'.\n\n" $
			"Please save as a different name.")
		mock.Verify.Never().saveStandardReport([anyArgs:])
		// Prefix is false, previous validation is bypassed
		Assert(.testSave(mock, name, false))
		mock.Verify.saveStandardReport([anyArgs:])

		// save_name is a "standard" report, and a non-standard variant exists
		// User chooses to rename the report via On_Save_As()
		name = 'My Report1'
		Assert(.testSave(mock, name, stdReportPrefix) is: false)
		mock.Verify.AlertInfo('Save As',
			'Standard TestPrefix reports cannot be modified.\n\nSaving as "' $ name $ '"')
		// Prefix is false, previous validation is bypassed
		Assert(.testSave(mock, name, false))
		mock.Verify.SetSaveName(name)
		mock.Verify.SetSaveName('~' $ stdReportPrefix $ ' ' $ name)
		mock.Verify.Times(2).On_Save_As()

		// save_name is a "standard" report, and a non-standard variant exists
		// User chooses to overwrite a pre-existing non-standard variant
		name = 'My Report2'
		mock.When.OverwriteReport?([anyArgs:]).Return(ID.YES)
		Assert(.testSave(mock, name, stdReportPrefix))
		mock.Verify.AlertInfo('Save As',
			'Standard TestPrefix reports cannot be modified.\n\nSaving as "' $ name $ '"')
		// Prefix is false, previous validation is bypassed
		Assert(.testSave(mock, name, false))
		mock.Verify.SetSaveName(name)
		mock.Verify.Never().SetSaveName('~' $ stdReportPrefix $ ' ' $ name)
		mock.Verify.Times(2).On_Save_As()

		// save_name is a "standard" report, and a non-standard variant exists
		// User chooses to cancel the process
		name = 'My Report3'
		mock.When.OverwriteReport?([anyArgs:]).Return(ID.CANCEL)
		Assert(.testSave(mock, name, stdReportPrefix) is: false)
		mock.Verify.AlertInfo('Save As',
			'Standard TestPrefix reports cannot be modified.\n\nSaving as "' $ name $ '"')
		// Prefix is false, as such, most of the checking is bypassed / unnecessary
		Assert(.testSave(mock, name, false))
		mock.Verify.SetSaveName(name)
		mock.Verify.SetSaveName('~' $ stdReportPrefix $ ' ' $ name)
		mock.Verify.Times(2).On_Save_As()
		}

	testSave(mock, name, prefix)
		{
		return mock.save?(name, 'Reporter - ~' $ prefix $ ' ' $ name, prefix)
		}

	Teardown()
		{
		QueryDo('delete params
			where report >= "Reporter - Reporter_Test"
			and report < "Reporter - Reporter_Test~"')
		super.Teardown()
		}
	}
