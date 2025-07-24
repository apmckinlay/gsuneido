// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: 'ReportLayoutDesign'
	New(.rptName)
		{
		super(.layout())
		}

	layout()
		{
		return Object('BrowseFlipForm',
			query: 'report_layout_designs
				extend notUsed = ""
				where report is ' $ Display(.rptName) $
				' rename rptdesign_num to rptdesign_num_new
				sort report, rptdesign_name',
			form: .form(.formLayout()),
			columns: #(rptdesign_name),
			linkField: 'notUsed',
			protectField: 'rptdesign_protect',
			validField: ReportLayoutDesign.GetValidField(.rptName),
			extraFmts: Object(rptdesign_name:
				Object(width: 10, font: #(weight: 'bold'), color: CLR.Highlight)),
			keyField: 'rptdesign_num_new',
			stretchColumn:,
			name: .Name $ .rptName,
			preventCustomExpand?:
			)
		}

	form(format)
		{
		return Object('Vert',
				Object('Horz',
					Object('EnhancedButton',
						image: #('view_list.emf', 'view_form.emf', highlighted: 1)
						mouseEffect:, imagePadding: .1, command: 'Flip',
						tip: 'switch between list and form view (Alt+A)', alignTop:),
					#Skip,
					#(Heading3, '', name: 'titleText'))
				format
			)
		}

	formLayout()
		{
		return Object(.wrapper,
			ReportLayoutDesign.GetLayout(.rptName),
			ReportLayoutDesign.GetProtect(.rptName)
			name: 'rptdesign_layout')
		}

	wrapper: Controller
		{
		New(layout, .protectField)
			{
			super(Object('Record', layout))
			.Send(#Data)
			.Data.SetProtectField(protectField)
			}

		Get()
			{
			return .Data.GetControlData().DeepCopy()
			}

		Set(value)
			{
			if value is ''
				value = []
			.Data.Set(value.DeepCopy())
			}

		Record_NewValue(@unused)
			{
			.Send(#NewValue, .Get())
			}

		Valid()
			{
			return .Data.Valid()
			}

		Valid?()
			{
			return .Data.Valid() is true
			}
		SetValid(valid /*unused*/)
			{
			}

		SetEditMode()
			{
			return not .Data.GetReadOnly()
			}

		Destroy()
			{
			.Send(#NoData)
			super.Destroy()
			}
		}

	LineItem_NewRowAdded(record)
		{
		record.report = .rptName
		record.rptdesign_layout = []
		}

	LineItem_AllowDelete(record, tranFromSave, source)
		{
		if tranFromSave is false and source.GetLineItems().Size() < 2
			{
			.AlertInfo('Delete', 'At least one layout is required.')
			return false
			}
		if tranFromSave isnt false and .usedOnSchedReport?(tranFromSave, record.vl_origin)
			{
			AlertDelayed('Can not delete layout because it is used on Scheduled Reports',
				'Delete')
			return false
			}
		return true
		}

	LineItem_AfterSave(data, t)
		{
		if not data.Member?('vl_origin')
			return

		if data.rptdesign_name isnt data.vl_origin.rptdesign_name
			.updateSchedReports(t, data.vl_origin, data)
		}

	BrowseFlipForm_Flip_To_Form(rec)
		{
		if false isnt ctrl = .FindControl('titleText')
			ctrl.Set(rec.rptdesign_name)

		return true
		}

	usedOnSchedReport?(t, oldrec)
		{
		if false is rpt = .findSchedReportName(oldrec)
			return false
		used? = false
		t.QueryApply('biz_scheduled_reports where bizrpt_name is ' $ Display(rpt))
			{
			if it.bizrpt_params.params_report_layout is oldrec.rptdesign_name
				{
				used? = true
				break
				}
			}
		return used?
		}

	findSchedReportName(oldrec)
		{
		rptName = oldrec.report
		reportsOb = GetStandardScheduleReportsOb()
		if false is rpt = reportsOb.FindIf({ it.name is rptName })
			return false
		return rpt
		}

	updateSchedReports(t, oldrec, newrec)
		{
		if false is rpt = .findSchedReportName(oldrec)
			return
		t.QueryApply('biz_scheduled_reports where bizrpt_name is ' $ Display(rpt))
			{
			if it.bizrpt_params.params_report_layout is oldrec.rptdesign_name
				{
				it.bizrpt_params.params_report_layout = newrec.rptdesign_name
				it.Update()
				}
			}
		}
	}
