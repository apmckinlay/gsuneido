// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
// TODO: Global > Reporter should fill the query (worked in 2005)
// TODO: move Summarize to a separate tab
// TODO: change Summarize, Sort, Select to use RepeatControl
// TODO: make Field wider on Summarize and Sort
// TODO: allow overriding prompt for Summarize Functions
// TODO: don't allow summarize by key
// TODO: don't allow summarize function on a summarize by field
// TODO: don't allow total on non-numeric fields
// TODO: wipe out Summarize Field if you choose Function of "count"
Controller
	{
	Title: Reporter
	Xmin: 790
	Ymin: 530
	SortRows: 8
	PortraitChars: 176
	LandscapeChars: 123
	DefaultColWidth: 15

	MaxSummarizeFields: 15	// if increased, create more Field_ definitions
	Ncalcs: 20

	CallClass(hwnd /*unused*/ = 0, rpt = '', printMode = false, addButtons = false,
		onDestroy = false, reporterMode = 'simple')
		{
		if rpt isnt '' and QueryEmpty?("params", report: rpt)
			{
			Alert("Reporter: Cannot find report. Has the report been renamed or deleted?",
				'Reporter', flags: MB.ICONWARNING)
			return
			}
		BookLog('Reporter')
		ModalWindow(Object(this, rpt, printMode, addButtons, dialog?:, :reporterMode),
			keep_size:, :onDestroy)
		return
		}

	New(rpt = '', .printMode = false, .addButtons = false, .dialog? = false,
		.reporterMode = 'simple')
		{
		super(ReporterLayout(printMode, reporterMode))
		.init(rpt)
		}

	init(rpt)
		{
		.Data.SetProtectField('reporter_protect')
		.Tabs = .Data.Vert.Tabs

		.source = .FindControl('DataSource')
		.select = .FindControl('Select2')
		.schedWarn = .FindControl('schedWarn')
		.reportHistory = .FindControl('reportHistory')

		.canvas = .FindControl('reporterCanvas')
		.rc = .FindControl('ReporterColumns')

		if .reporterMode is #form
			.rc.SetCanvas(.canvas)

		.open_report(rpt)

		if .printMode is true
			{
			if .addButtons is false
				.Defer(.setPrintMode)
			.Window.AddValidationItem(this)
			}
		Params.Ensure()
		}

	tabs: (Input: 0, Formulas: 1, Design: 2, Sort: 3, Select: 4)
	setPrintMode()
		{
		// if user is going through book options too fast,
		// sometimes this is destroyed already
		if .Member?("Tabs")
			.Tabs.Select(.tabs.Select)
		}

	Record_NewValue(field, value /*unused*/)
		{
		if field is 'Source'
			.setQuery()
		else if field.Prefix?('sort')
			{
			.Data.SetField(field.Replace('sort', 'total'),
				.subtotal?(field))
			.notifyDesignChange()
			}
		else if field is 'formulas' or field.Prefix?("summarize_")
			.setQuery()
		}

	setQuery()
		{
		.rpt.SetQuery(.source.Source(), .Data.Get())
		.select.ChangeSF(.rpt.GetSelectFields())
		.setFields()
		.notifyDesignChange()
		}

	notifyDesignChange()
		{
		if .canvas isnt false
			.canvas.DesignChanged()
		}

	setFields()
		{
		.Data.SetField('allcols', .rpt.GetAllCols())
		.Data.SetField('design_cols', .rpt.GetDesignCols())
		.Data.SetField('summarize_func_cols', .rpt.GetSummarizeFuncCols())
		.Data.SetField('nonsummarized_fields', .rpt.GetNonSummarizedFields())
		.Data.SetField('selectFields', .rpt.GetSelectFields())
		}

	On_Help()
		{
		OpenBook(LastContribution('HelpBook'),
			.reporterMode is #form
				? #(name: 'Reporter Form', path: 'Reporter Form')
				: #(name: 'Reporter', path: 'Reporter'))
		}
	On_New()
		{
		if false is .dirtysave()
			return

		// TODO: move the following into clear()?  See sugg. 14437
		.rpt.Clear()
		.select.ChangeSF(.rpt.GetSelectFields())
		.clear()

		if .canvas isnt false
			.canvas.Set(#(), defaultHeader:)
		}
	clear()
		{
		.Data.Set(Record(columns: #(), coloptions: Object(), select: #{},
			allcols: #(), design_cols: #(), heading1: PageHeadName()))
		.rc.Set(#())
		.select.Set(Record())
		.Window.SetTitle('Reporter')
		}
	dirtysave()
		{
		if not .Dirty?()
			return true
		result = Alert('Do you want to save your changes?', 'Reporter',
			hwnd: .Window.Hwnd, flags: MB.ICONEXCLAMATION | MB.YESNOCANCEL)
		switch (result)
			{
		case ID.CANCEL :
			return false
		case ID.YES :
			return .On_Save()
		case ID.NO :
			.Dirty?(false)
			return true
		default :
			return false
			}
		}
	On_Open()
		{
		if not .dirtysave()
			return
		rpt = ReporterOpenDialog(.Window.Hwnd, reporterMode: .reporterMode)
		if rpt is false
			return
		.open_report(rpt)
		}

	open_report(rpt, dirty = false)
		{
		.clear()
		.rpt = ReporterModel(rpt, defaultMode: .reporterMode)
		if rpt is ''
			{
			.Defer(.updateWindowTitle)
			return
			}
		if not .source.Authorized?() or
			not ReporterOpenDialog.HasPermission?(Record(params: .rpt.GetData()))
			{
			AlertDelayed("Sorry, you are not authorized to access this report.",
				"Reporter", hwnd: .Window.Hwnd, flags: MB.ICONEXCLAMATION)
			return
			}

		.Data.Set(.rpt.GetData())
		.clear_column_lists(.Data.Get())
		.setFields()

		// DeepCopy is required to separate .Data.Set from .rc.Set.
		// Otherwise Data and .rc are synced via the object reference,
		// making the comparison in .Dirty? irrelevant.
		.rc.Set(.rpt.GetColumns().DeepCopy())
		.select.ChangeSF(.rpt.GetSelectFields())
		.select.Set(.rpt.GetData().select isnt ''
			? .rpt.GetData().select.Copy() : Record())

		.setSchedWarning(.rpt)
		.setReportHistory(.rpt)
		// need to delay because the Dialog will replace title to Reporter
		.Defer(.updateWindowTitle)
		.Dirty?(dirty)
		}

	setSchedWarning(rpt)
		{
		if false isnt .schedWarn
			{
			if .isScheduledReport?(rpt.GetSaveName().AfterFirst('Reporter - '))
				.schedWarn.Set('Scheduled')
			else
				.schedWarn.Set('')
			}
		}

	isScheduledReport?(rpt_name)
		{
		return not QueryEmpty?('biz_scheduled_reports', bizrpt_name: rpt_name)
		}

	setReportHistory(rpt)
		{
		if false isnt .reportHistory
			{
			data = rpt.GetData()
			created = .buildHistoryString('created', data)
			lastMod = .buildHistoryString('last_modified', data)
			for field in data.report_last_ran.Members()
				data[field] = data.report_last_ran[field]
			lastRan = data.last_ran_on isnt false
				? .buildHistoryString('last_ran', data)
				: ''
			.reportHistory.Set(created $ lastMod $ lastRan)
			}
		}

	buildHistoryString(prefix, data)
		{
		str = prefix.CapitalizeWords().Replace('_', ' ')
		date = data.GetDefault(prefix $ '_on', '')
		if Date?(date)
			date = date.ShortDate()
		// cannot use Display due to '' => "''"
		str $= Opt(' on ', String(date))
		user = data.GetDefault(prefix $ '_by', '')
		return Opt(str, ' by ', user, '\r\n')
		}

	updateWindowTitle()
		{
		if not .Destroyed?()
			.Window.SetTitle(.rpt.GetTitle())
		}

	On_Save()
		{
		if not .Dirty?()
			return false

		.report_name = .rpt.GetSaveName().Trim('Reporter - ')
		if true isnt valid = .check_valid(from_save?:)
			{
			Alert(valid, "Reporter", .Window.Hwnd, MB.ICONERROR)
			return false
			}
		name = .Data.GetField('report_name')
		save_name = .rpt.GetSaveName()
		stdReportPrefix = .stdReportPrefix()
		return .save?(name, save_name, stdReportPrefix)
			? .save()
			: false
		}

	stdReportPrefix()
		{
		return LastContribution('Reporter').StandardReportPrefix
		}

	save?(name, save_name, stdReportPrefix)
		{
		if name is ''
			return .On_Save_As()
		else if .CheckName(name, stdReportPrefix) is false
			return false
		else if .saveStandardReport(name, save_name, stdReportPrefix) is false
			return false
		return true
		}

	saveStandardReport(name, save_name, stdReportPrefix)
		{
		if stdReportPrefix isnt false and
			save_name.Prefix?('Reporter - ~' $ stdReportPrefix)
			{
			.AlertInfo("Save As",
				'Standard ' $ stdReportPrefix $
				' reports cannot be modified.\n\nSaving as "' $ name $ '"')
			.SetSaveName(name)
			if false isnt overwrite? = .OverwriteReport?(name)
				switch overwrite?
					{
				// Need to return save_name back to original value or else save will
				// overwrite the report on next Save click without prompting
				case ID.CANCEL :
					.SetSaveName(save_name.RemovePrefix('Reporter - '))
					return false
				case ID.YES :
					return true
				case ID.NO :
					.SetSaveName(save_name.RemovePrefix('Reporter - '))
					return .On_Save_As()
					}
			}
		return true
		}

	InvalidMsg(type = 'load')
		{
		return "Unable to " $ type $ " report: Invalid Input Data Source" $
			LastContribution('Reporter').ReporterInvalidMsg
		}
	check_valid(from_save? = false)
		{
		if .rpt.GetQuery() is ''
			{
			.Tabs.Select(.tabs.Input)
			return .InvalidMsg(from_save? ? 'save' : 'load')
			}
		source = .source.Source()
		if source.Empty?()
			return "Please choose an Input Data Source"

		if source.Member?('validField')
			{
			rec = .Data.Get().Copy()
			rec.Invalidate(source.validField)
			errs = rec[source.validField]
			if errs isnt ''
				return errs
			}
		// use forceCheck to check valid even when not dirty in case invalid report
		// was loaded and may not have been modified so it is not dirty.
		// Didn't want to just set RecordControl to dirty, because then it will ask
		// the user to save even when nothing was changed.
		errs = .Data.Valid(forceCheck:)
		if errs isnt true
			return errs

		return .check_reporter_options() // select, formula, sort, summary
		}
	check_reporter_options()
		{
		if .select.Where() is false
			{
			.Tabs.Select(.tabs.Select)
			return "Invalid Select"
			}
		if '' isnt msg = .validate_formula_menu_options()
			return msg

		if '' isnt msg = .checkAllowedSortSummaryFields()
			return msg

		if '' isnt msg = .checkSummaryInfo()
			return msg

		if '' isnt msg = .checkCanvas()
			return msg

		return true
		}

	validate_formula_menu_options()
		{
		orig_prompts = Object()
		for f in .rpt.GetSelectFields().OrigFields()
			orig_prompts.Add(.rpt.GetSelectFields().FieldToPrompt(f))
		prompts = Object()
		form_ob = Object()
		for row in .Data.Get().formulas
			{
			field = row.calc.Trim()
			formula = row.formula.Trim()
			form_val = row.form_val
			if field isnt '' and (formula isnt '' or form_val is true)
				{
				if '' isnt result = .validFormula(row, field, prompts, orig_prompts,
					form_ob)
					return result
				}
			}
		return .validateFormulaValues(form_ob)
		}

	validFormula(formula, prompt, prompts, orig_prompts, form_ob)
		{
		if formula.type is ''
			return 'Format is required for formula: ' $ prompt

		if '' isnt valid = .validateFormulaName(prompt, prompts, orig_prompts)
			return valid

		if '' isnt valid = .validateFormulaCode(formula, prompt)
			return valid

		if formula.form_val is true
			{
			if not formula.formula.Blank?()
				return 'Should not have both Formula and Menu Option'
			form_ob.Add(formula.calc)
			}
		prompts.Add(formula.calc)
		return ''
		}

	validateFormulaName(prompt, prompts, orig_prompts)
		{
		if prompts.Has?(prompt)
			return 'Formula field name ' $ prompt $ ' is a duplicate. ' $
				'Please rename one of the formulas.'
		if orig_prompts.Has?(prompt)
			return .fieldNameInUseMsg(prompt)
		return ''
		}

	fieldNameInUseMsg(prompt)
		{
		return 'Formula field name ' $ prompt $ ' is in use. ' $
			'Please rename the formula.'
		}

	validateFormulaCode(formula, prompt)
		{
		sf = .rpt.GetSelectFields()
		field = sf.PromptToField(prompt)
		// check here again in case the prompt gives us a standard field name, which
		// means the prompt given for the formula conflicts with field not in the
		// original select fields (before joins, could be a _abbrev or _name field)
		if field !~ "^calc\d+$"
			return .fieldNameInUseMsg(prompt)
		newFn = FormulaEditor.ConstructFormula(
			sf, prompt, formula.type, formula.formula, field, quiet:)
		if newFn.formulaCode is false or
			(newFn.formulaCode is '' and formula.form_val isnt true)
			{
			.DeleteDD(newFn.ddName)
			return "Invalid or missing operators in formula " $ formula.calc $ '.'
			}
		if Object?(newFn.formulaCode) and newFn.formulaCode.Member?('err')
			{
			.DeleteDD(newFn.ddName)
			return newFn.formulaCode.err
			}
		if '' isnt msg = CustomizeField.ExtraCheck(newFn.formulaCode)
			{
			.DeleteDD(newFn.ddName)
			return msg
			}
		return ''
		}

	DeleteDD(name)
		{
		QueryApply1('configlib where name is ' $ Display(name))
			{
			it.Delete()
			}
		}

	validateFormulaValues(form_ob)
		{
		sel = .select.Get()
		for (i = 0; i < Select2Control.Numrows; ++i)
			if sel['menu_option' $ i] is true and form_ob.Has?(sel['fieldlist' $ i])
				return 'Formula Value can not used as Menu Option'

		return ''
		}

	checkAllowedSortSummaryFields()
		{
		data = .Data.Get()
		for i in .. .SortRows
			{
			prompt = data['sort' $ i]
			if prompt isnt '' and .editorCtrl?(prompt)
				return 'Sorting on large text field is not allowed.\n' $
					'Please remove the sort on ' $ prompt
			}

		prompts = data.summarize_by.Split(',')
		for prompt in prompts
			if .editorCtrl?(prompt)
				return 'Summarize by large text field is not allowed.\n' $
					'Please remove ' $ prompt $ ' from Summarize By'

		return ''
		}

	editorCtrl?(prompt)
		{
		if false is (field = .rpt.GetSelectFields().PromptToField(prompt))
			return false
		dd = Datadict(field)
		return dd.Control[0] is 'Editor' or dd.Control[0] is 'ScintillaAddonsEditor'
		}

	checkSummaryInfo()
		{
		data = .Data.Get()
		if .reportHasInputSummary?(data) and .anyFieldHasSummaryOption?(data)
			return "Summarize options (Min, Max, Average) on Design columns " $
				" can not be used when Summarize options are being used on the Input tab"
		return ""
		}

	reportHasInputSummary?(data)
		{
		if data['summarize_by'] isnt ''
			return true
		for (i = 0; i < .MaxSummarizeFields; ++i)
			if data['summarize_func' $ i] isnt ''
				return true
		return false
		}

	anyFieldHasSummaryOption?(data)
		{
		if not Object?(data.GetDefault('coloptions', false))
			return false

		for fieldOptions in data.coloptions
			{
			if not Object?(fieldOptions)
				continue
			for summaryField in #(min max average)
				if fieldOptions.GetDefault(summaryField, false) is true
					return true
			}
		return false
		}

	checkCanvas()
		{
		if .canvas is false
			return ''

		return .canvas.CheckContent()
		}

	resetKeys: false
	On_Save_As()
		{
		if true isnt (valid = .check_valid(from_save?:))
			{
			Alert(valid, "Reporter", .Window.Hwnd, MB.ICONERROR)
			return false
			}
		forever
			{
			name = Ask('Save As', title: 'Reporter', hwnd: .Window.Hwnd)
			if name is ""
				{
				Alert("Please enter a name to save as", "Reporter",
					.Window.Hwnd, MB.ICONWARNING)
				continue
				}
			if .CheckName(name) is false
				return false
			if false isnt (overwrite? = .OverwriteReport?(name))
				{
				switch overwrite?
					{
				case ID.CANCEL :
					return false
				case ID.YES :
					break
				case ID.NO :
					continue
					}
				}
			break
			}
		.report_name = .Data.Get().report_name
		.Data.SetField('report_name', name)
		.SetSaveName(name)
		.setSchedWarning(.rpt)
		.updateWindowTitle()
		return .save()
		}
	SetSaveName(name)
		{
		.rpt.SetSaveName('Reporter - ' $ name)
		.resetKeys = .report_name isnt "" and .report_name isnt name
		}
	CheckName(name, stdReportPrefix = false)
		{
		if name is false
			return false
		if name =~ "[^a-zA-Z0-9 -]"
			{
			.AlertInfo("Save As", "You can only use alpha-numeric characters " $
				"(letters and numbers) and spaces in the name.\n\n" $
				"Please save as a different name.")
			return false
			}

		if stdReportPrefix is false
			stdReportPrefix = .stdReportPrefix()

		if stdReportPrefix isnt false and name.Prefix?(stdReportPrefix)
			{
			.AlertInfo("Save As",
				"Reporter name cannot start with '" $ stdReportPrefix $ "'.\n\n" $
				"Please save as a different name.")
			return false
			}
		return true
		}

	OverwriteReport?(name)
		{
		if false is rec = QueryFirst('params
			where report is ' $ Display('Reporter - ' $ name) $ ' sort user')
			return false

		if .reporterMode isnt mode = rec.params.GetDefault(#reporterMode, 'simple')
			{
			.AlertWarn('Reporter', 'Reporter ' $ (mode is #form ? 'Form' : 'Report') $
				Display(name) $ ' already exists\nPlease save as a different name')
			return ID.NO
			}

		return Alert(
			Display(name) $ ' already exists\nDo you want to overwrite it?', 'Reporter',
			hwnd: .Window.Hwnd, flags: MB.ICONEXCLAMATION | MB.YESNOCANCEL)
		}

	save()
		{
		if .CheckScheduled(.rpt.GetSaveName())
			return false

		data = .prepare_to_save()
		rpt_rec = Record(report: .rpt.GetSaveName(), params: data)

		try
			{
			if false is .handleReportMenu(rpt_rec)
				return false
			}
		catch (errs, 'Reporter:')
			{
			.AlertWarn('Reporter', errs.Replace('Reporter: ', ''))
			return false
			}

		.output_params(rpt_rec)

		.Dirty?(false)
		return true
		}
	changeWarningMsg: 'This Reporter Report has been set up as a Business - Scheduled ' $
		'Reports. Making changes here could require changes to be made to the Scheduled' $
		' Report for it to run correctly.\n\n' $
		'Are you sure you want to continue?'
	saveWarningMessage()
		{
		return OkCancel(.changeWarningMsg, .Title)
		}
	CheckScheduled(name)
		{
		if .isScheduledReport?(name.Replace('^Reporter - ', ''))
			if not .saveWarningMessage()
				return true
		return false
		}

	handleReportMenu(rpt_rec)
		{
		c = LastContribution('Reporter')
		return c.HandleReportMenu(rpt_rec, .source, .Dirty?(), .Data,
			reporterMode: .reporterMode)
		}

	output_params(rpt_rec)
		{
		.clear_column_lists(rpt_rec.params)
		ReporterDataConverter.ToFields(rpt_rec.params, .rpt.GetSelectFields())
		RetryTransaction()
			{ |t|
			t.QueryDo("delete params where report is " $ Display(.rpt.GetSaveName()))
			t.QueryOutput("params", rpt_rec)
			}
		}

	prepare_to_save()
		{
		origData = .Data.Get()
		origData.columns = .rc.Get()
		origData.select = .select.Get()
		data = .prepare_data(origData)
		if .resetKeys is true
			{
			for m in data.formulas.Members()
				origData.formulas[m].key = data.formulas[m].key =
					Display(Timestamp()).Tr('#.')
			.setQuery()
			if Object?(items = data.GetDefault('reporterCanvas', []).items)
				{
				sf = .rpt.GetSelectFields()
				for item in items
					{
					if false is fld = item.FindOne({ String?(it) and it =~ "^calc\d+$" })
						continue
					item.Replace(fld, sf.PromptToField(Datadict.PromptOrHeading(fld)))
					}
				}
			}
		return data
		}

	prepare_data(data)
		{
		data = data.Copy()
		data.Delete('allcols')
		data.Delete('design_cols')
		data.Delete('reporter_cols')
		data.Delete('reporter_sortcolumns')
		data.Delete('nonsummarized_fields')
		data.Delete('summarize_func_cols')
		data.Delete('reporter_permission_list')
		data.Delete('reporter_summarizeby_cols')
		if data.formulas is ''
			data.formulas = #()
		else
			data.formulas = data.formulas.DeepCopy()
		if data.coloptions is ''
			data.coloptions = #()
		else
			data.coloptions = data.coloptions.DeepCopy()
		.clearHistoryFields(data)
		for i in .. .MaxSummarizeFields
			data.Delete('summarize_field' $ i $ '__protect')

		if .reporterMode is 'enhanced'
			data.reporterCanvas = .canvas.Get()
		data.reporterMode = .reporterMode
		if not data.Member?('created_by')
			{
			data.created_by = Suneido.User
			data.created_on = Timestamp()
			}
		else
			{
			data.last_modified_by = Suneido.User
			data.last_modified_on = Timestamp()
			}
		return data
		}

	clearHistoryFields(data)
		{
		data.Delete('report_last_ran')
		data.Delete('last_ran_on')
		data.Delete('last_ran_by')
		data.Delete('lastRan')
		data.Delete('last_modified_on')
		data.Delete('last_modified_by')
		if .resetKeys is true
			{
			data.Delete('created_on')
			data.Delete('created_by')
			}
		}

	clear_column_lists(data)
		{
		data.Delete('allcols')
		data.Delete('reporter_cols')
		data.Delete('selectFields')
		data.Delete('design_cols')
		data.Delete('reporter_summarizeby_cols')
		data.Delete('summarize_func_cols')
		data.Delete('nonsummarized_fields')
		}

	Dirty?(dirty = '')
		{
		if .Data.Dirty?(dirty)
			return true
		data = .Data.Get()
		if data.columns isnt .rc.Get() or
			data.select isnt .select.Get()
			return true
		return false
		}

	last_formula: 0
	FormulaKillFocus(idx)
		{
		.last_formula = idx
		}
	On_Add(option)
		{
		if .rpt.GetSelectFields().Fields.Empty?()
			return Alert("Please choose an Input Data Source", 'Reporter'
				hwnd: .Window.Hwnd, flags: MB.ICONWARNING)

		FormulaEditor['Add_a_' $ option.Tr(' ', '_')](.findControl(.last_formula),
			selectFields: .rpt.GetSelectFields(), hwnd: .Window.Hwnd)
		}

	On_Add_a_Function(fnName)
		{
		FormulaEditor.Add_a_Function(.findControl(.last_formula), fnName,
			.rpt.GetSelectFields())
		}

	On_Add_an_Operator(option)
		{
		FormulaEditor.Add_an_Operator(.findControl(.last_formula), option)
		}

	findControl(pos)
		{
		rows = .Data.GetControl('formulas').GetRows()
		while pos >= 0 and not rows.Member?(pos)
			pos--
		return pos < 0 ? false : rows[pos].GetControl('formula')
		}

	FormulaEditor_Click(source)
		{
		FormulaEditor.HighlightSelection(source, .rpt.GetSelectFields())
		}

	ChooseField(prompts, hwnd)
		{
		return ToolDialog(hwnd,
			Object(.choose_field, prompts.Sort!()),
			border: 0)
		}
	choose_field: Controller
		{
		Title: 'Add a Field'
		Xmin: 300
		Ymin: 400
		New(list)
			{
			super(Object('ListBox', list))
			.list = list
			}
		ListBoxSelect(i)
			{
			if .list.Member?(i)
				.Window.Result(.list[i])
			return 0 // no currently selected item
			}
		}

	On_AddRemove_Columns()
		{
		curlist = .rc.Get().Map({ it.text })
		colwidths = Object()
		for col in .rc.Get()
			colwidths[col.text] = Object(width: col.GetDefault('width', false))
		list = OkCancel(Object('TwoListDlg', .rpt.GetDesignCols(), curlist)
			title: 'Add/Remove Columns', hwnd: .Window.Hwnd)
		if list is false or list is curlist
			return
		list = list.Split(',')
		if list is curlist
			return
		coloptions = .Data.GetField('coloptions')
		for item in curlist.Difference(list)
			coloptions.Delete(item)
		.rc.Set(.buildColumns(list, colwidths, coloptions))
		}
	GetRPT()
		{
		return .rpt
		}
	GetRptSortCols()
		{
		try
			return .rpt.GetSortCols().Map({ it.prompt })
		catch (unused, 'Reporter:')
			return Object()
		}

	GetRptDesignCols()
		{
		return .rpt.GetDesignCols()
		}

	PromptToField(prompt)
		{
		return .rpt.GetSelectFields().PromptToField(prompt)
		}

	FieldToPrompt(field)
		{
		return .rpt.GetSelectFields().FieldToPrompt(field)
		}

	buildColumns(list, colwidths, coloptions)
		{
		cols = Object()
		for fld in list
			{
			if colwidths.Member?(fld)
				width = colwidths[fld].width
			else
				{
				// check for invalid field in list first before trying to get dd width
				if false is field_name = .rpt.GetSelectFields().PromptToField(fld)
					continue
				dd = Datadict(field_name)
				width = Object?(dd.Format) and dd.Format.Member?('width')
					? dd.Format.width : .DefaultColWidth
				if dd.Format[0] is 'Image'
					width = (width / 8.5.InchesInTwips()) * .LandscapeChars/*=page width*/

				coloptions[fld] = Record(heading: fld)
				}
			cols.Add(Object(text: fld, :width))
			}
		return cols
		}

	On_Clear_Select()
		{
		.select.Set(Record())
		}

	On_Print()
		{
		.print_report("On_Print")
		}
	On_PDF_Save_to_file()
		{
		.print_report("On_PDF_Save_to_file")
		}
	On_PDF_Download()
		{
		.On_PDF_Save_to_file()
		}
	On_PDF_Email_as_attachment()
		{
		.print_report("On_PDF_Email_as_attachment")
		}
	ClearReporterColumns()
		{
		.rc.Set(#())
		}

	On_Preview()
		{
		.print_report("On_Preview")
		}

	print_report(option)
		{
		dirtyState = .Dirty?()
		if false is report = .report()
			return
		report.previewWindow = .Window.Hwnd
		try
			Params[option](@report)
		catch (err)
			.checkAndAlert(err)
// Temporary logging for issue 33830
if .Destroyed?()
	{
	SuneidoLog('ERROR: (CAUGHT) Reporter destroyed', calls:)
	SuneidoLog('ERROR: (CAUGHT) Reporter destroyed - extra info', calls: .destroyCS,
		params: [at: .destroyAt, :option], caughtMsg: 'for 33830')
	}
else
	.Dirty?(dirtyState)
		}
	checkAndAlert(err)
		{
		if err.Has?('assertion failure: key.cursize()') or
			err.Has?('assertion failure: size')
			.AlertError("Reporter", "Invalid Summarize By/Sort")
		else if err.Prefix?('SHOW:')
			throw err
		}
	On_Page_Setup()
		{
		Params.On_Page_Setup(name: .rpt.GetTitle(), hwndOwner: .Window.Hwnd)
		}
	On_Export()
		{
		.print_report("On_Export")
		}
	tryCheckAndBuildQuery()
		{
		try
			{
			queryInfo = .checkAndBuildQuery()
			return queryInfo isnt false
				? queryInfo[0] // ReporterModel.BuildReport > rpt member
				: false
			}
		catch (errs, 'Reporter:')
			{
			.AlertWarn('Reporter', errs.Replace('Reporter: ', ''))
			return false
			}
		}
	report()
		{
		try
			return .checkAndBuildQuery()
		catch (errs, 'Reporter:')
			{
			if errs =~ '; select tab:'
				{
				tab = errs.AfterLast('; select tab: ')
				.Tabs.Select(.tabs.GetDefault(tab, 0))
				errs = errs.BeforeLast(';')
				}
			.AlertWarn('Reporter', errs.Replace('Reporter: ', ''))
			return false
			}
		}
	checkAndBuildQuery()
		{
		if true isnt (valid = .check_valid())
			throw "Reporter: " $ valid

		data = .Data.Get()
		data.columns = .rc.Get()
		data.select = .select.Get()

		return .rpt.BuildReport(data, source: .source.Source())
		}

	Click(item)
		{ .Properties(item) }
	Properties(item)
		{
		coloptions = .Data.GetField('coloptions')
		if not coloptions.Member?(item)
			coloptions[item] = Record(heading: item)
		options = ToolDialog(.Window.Hwnd,
			Object(.coloptionsControl, item, coloptions[item].Copy()))
		if options is false or options is coloptions[item]
			return
		coloptions[item] = options
		.Dirty?(true)
		}
	ClearProperties(item)
		{
		coloptions = .Data.GetField('coloptions')
		coloptions.Delete(item)
		.Dirty?(true)
		}
	coloptionsControl: Controller
		{
		Title: 'Column Properties'
		New(col, options)
			{
			super(.controls(col))
			// handle options as object or record, ensuring data is set to a record
			.Data.Set(Record().Merge(options))
			}
		controls(col)
			{
			return Object('Record'
				Object('Vert'
					Object('Heading' col)
					'Skip'
					#(Pair
						(Static Heading)
						(Editor xmin: 200 height: 5 name: heading))
					'Skip'
					#('Horz'
						#(CheckBox 'Total', name: 'total')
						'Skip'
						#(CheckBox 'Min', name: 'min')
						'Skip'
						#(CheckBox 'Max', name: 'max')
						'Skip'
						#(CheckBox 'Average', name: 'average')
						'Skip'
						#(CheckBox 'Count', name: 'count')
						)
					'Skip'
					#(CheckBox 'Print Summary Prompts In This Column',
						name: 'print_summary_prompts')
					'Skip'
					#(OkCancel))
				xstretch: 0, ystretch: 0)
			}
		On_OK()
			{
			data = .Data.Get()
			if data.print_summary_prompts is true and .anySummaryOptionsChecked?()
				{
				.AlertInfo(.Title, 'Summary prompts can not be printed in the ' $
					'same column as the summary options are used')
				return
				}
			.Window.Result(data)
			}

		anySummaryOptionsChecked?()
			{
			data = .Data.Get()
			return data.total is true or data.min is true or data.max is true or
				data.average is true or data.count is true
			}
		}
	WndPane_ContextMenu(x, y)
		{
		if .Tabs.TabGetSelectedName() isnt 'Design'
			return
		.rc.ContextMenu(x, y)
		}
	subtotal?(field)
		{
		data = .Data.Get()
		prompt = data[field]
		if prompt is '' or not .rpt.GetSelectFields().HasPrompt?(prompt)
			return false
		field = .rpt.GetSelectFields().PromptToField(prompt)
		// TODO: also don't total unique indexes
		return not .rpt.GetKeys().Has?(field)
		}

	On_View_Data()
		{
		if false is rpt = .tryCheckAndBuildQuery()
			return
		ReporterViewDataControl(rpt[1], rpt.exclude, .rpt.GetSelectFields(), .Window.Hwnd)
		}

	import_export_type: 'report design '
	filter: "Report Design (*.rpt)\x00*.rpt"
	On_Export_Report_Design()
		{
		data = .prepare_to_save()
		data.Delete('reporter_permissions')
		data.Delete('selectFields')
		ImportExportObject.Export('Export Report Design', data,
			.import_export_type $ .reporterMode, .filter, .Window.Hwnd,
			name: data.report_name)
		}
	On_Import_Report_Design()
		{
		if not .dirtysave()
			return
		try
			{
			data = ImportExportObject.Import('Import Report Design',
				.import_export_type $ .reporterMode, .filter, .Window.Hwnd)
			if data is false
				return
			// need to invalidate because it gets removed
			// (in prepare_data method) from data when exporting
			data.Invalidate('reporter_summarizeby_cols')
			.open_report(data, dirty:)
			}
		catch
			Alert("An error occurred during the import. " $
				"The import file may be invalid.\n" $
				"Only reports exported from Reporter can be imported",
				"Reporter", .Window.Hwnd, MB.ICONERROR)
		}

	Ok_to_CloseWindow?()
		{
		return .dirtysave()
		}

	GetDataSrcCtrl()
		{
		return .source
		}

	ConfirmDestroy()
		{
		if .dialog?
			return true
		else
			return .dirtysave()
		}

	On_Cancel()
		{ // so you can't close with ESC
		}

destroyAt: false
destroyCS: ''
	Destroy()
		{
		if .printMode
			.Window.RemoveValidationItem(this)
		super.Destroy()
.destroyAt = Date()
.destroyCS = FormatCallStack(GetCallStack(limit: 20), levels: 20)
		}
	}
