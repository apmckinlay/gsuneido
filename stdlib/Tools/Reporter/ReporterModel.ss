// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
class
	{
	errMsg: ''
	New(rpt = '', .defaultMode = 'simple')
		{
		if rpt is ''
			{
			.setEmpty()
			return
			}

		if String?(rpt)
			{
			if false is x = Query1("params", report: rpt)
				{
				.setEmpty()
				.errMsg = "Reporter: Cannot find report. " $
					"Has the report been renamed or deleted?"
				return
				}
			x = x.params
			}
		else
			{
			x = rpt
			rpt = "Reporter - " $ x.report_name
			}
		.reporterMode = x.GetDefault(#reporterMode, .defaultMode)
		if x.heading1 is "" and .reporterMode is 'simple'
			x.heading1 = PageHeadName()
		.data = x
		.source = GetCustomReportsSource(x.Source)
		Assert(.source isnt false)
		.source.Set_default('')
		.SetQuery(.source)
		if not Object?(x.columns.GetDefault(0, #())) // converted old reports
			x.columns.Map!({Object(text: it, width: Reporter_datadict.GetWidth(.sf, it))})
		.title = 'Reporter - ' $ x.report_name
		.save_name = rpt
		}

	BuildSelectFields(sourceName)
		{
		source = GetCustomReportsSource(sourceName)
		query = Reporter_query.GetQuery(source)
		sf = SelectFields(Reporter_query.GetColumns(query), .BuildExclude(source),
			convertCustomNumFields:)
		return sf
		}

	setEmpty()
		{
		.Clear()
		.data = []
		.source = []
		.reporterMode = .defaultMode
		}

	SetQuery(source, data = false)
		{
		if data isnt false
			.data = data

		.query = Reporter_query.GetQuery(source)
		if .query is ""
			{
			.Clear()
			return
			}
		.keys = Reporter_query.GetKeys(.query)
		.sf = SelectFields(Reporter_query.GetColumns(.query), .BuildExclude(source),
			convertCustomNumFields:)

		.infoFields = Reporter_extend_info.AddFields(.sf)
		.add_calc_fields()
		ReporterDataConverter.ToPrompts1(.data, .sf)
		.add_summarize_fields()
		.failedConverts = _failedConverts = Object()
		ReporterDataConverter.ToPrompts2(.data, .sf)
		.allcols = .sf.Prompts().Sort!()
		.design_cols = .summarize_cols.Empty?() ? .allcols : .design_cols_with_summarize()
		.calc_fields_not_allowed_to_summarize()

		// make available to rules
		.data.allcols = .allcols
		.data.design_cols = .design_cols
		.data.summarize_func_cols = .summarize_func_cols
		.data.nonsummarized_fields = .nonsummarized_fields
		.data.selectFields = .sf
		}

	BuildExclude(source = false)
		{
		if source is false
			source = .source
		return Reporter_query.BuildExclude(source)
		}

	add_summarize_fields()
		{
		.summarize_cols = Object()
		.summarize_func_cols = Object()
		.summarize_fields = Object()
		.summarize_joinfields = Object()
		.addSummaryByFields()
		for i in .. Reporter.MaxSummarizeFields
			{
			if "" is (func = .data['summarize_func' $ i])
				continue
			if func is "count"
				.add_summarize_field("count", "Count", summaryFunc:)
			else
				{
				prompt = .data['summarize_field' $ i]
				if prompt is "" or not .sf.HasPrompt?(prompt) or
					not .func_map.Member?(func)
					continue
				field = .sf.PromptToField(prompt)
				.add_summarize_field(.func_map[func] $ '_' $ field,
					func.Capitalize() $ " " $ prompt, summaryFunc:)
				.summarize_joinfields.Add(field)
				}
			}
		}

	addSummaryByFields()
		{
		for prompt in .data.summarize_by.Split(',')
			if .sf.HasPrompt?(prompt)
				{
				.summarize_cols.Add(prompt)
				.summarize_joinfields.Add(.sf.PromptToField(prompt))
				}
		}
	add_summarize_field(field, prompt, summaryFunc = false)
		{
		if .summarize_fields.Has?(field)
			return
		.summarize_cols.Add(prompt)
		if summaryFunc
			.summarize_func_cols.Add(prompt)
		.summarize_fields.Add(field)
		.sf.AddField(field, prompt)
		}
	design_cols_with_summarize()
		{
		cols = .summarize_cols.Copy()
		.ForEachCalc()
			{ |key /*unused*/, field|
			cols.AddUnique(field)
			}
		return cols
		}

	calc_fields_not_allowed_to_summarize()
		{
		// don't allow user to pick formula field that contains summarized field
		.nonsummarized_fields = Object()
		.ForEachCalc()
			{ |key /*unused*/, field, formula|
			.sf.FormulaPromptsToFields(formula, fields = Object())
			// does not contain summarized field
			if not fields.Intersect(.summarize_fields).Empty?()
				.nonsummarized_fields.Add(field)
			}
		}

	add_calc_fields()
		{
		.ForEachCalc()
			{ |key, field|
			.sf.AddField('calc' $ key, field)
			}
		}
	ForEachCalc(block)
		{
		for row in .data.formulas
			{
			field = row.calc.Trim()
			formula = row.formula.Trim()
			form_val = row.form_val
			if field isnt '' and (formula isnt '' or form_val is true)
				block(row.key, field, :formula, :row)
			}
		}

	GetColumns()
		{
		return .data.columns
		}
	GetData()
		{
		return .data
		}
	GetSelectFields()
		{
		return .sf
		}
	GetTitle()
		{
		return .reporterMode is #form
			? .title.Replace('^Reporter', 'Reporter Form')
			: .title
		}
	GetReportTitle()
		{
		return .title
		}
	GetHeader()
		{
		return .reporterMode is 'simple'
			? Object('PageHead', .data.heading2, title2: .data.heading1)
			: false
		}
	GetFooter()
		{
		return .reporterMode is 'enhanced'
			? ReporterFormat.LoadCanvasFromSavedData(.data.reporterCanvas, 'page_footer')
			: false
		}
	GetSaveName()
		{
		return .save_name
		}
	SetSaveName(name)
		{
		.save_name = .title = name
		// TODO: get rid of one of this two
		}
	GetAllCols()
		{
		return .allcols
		}
	GetDesignCols()
		{
		return .design_cols
		}
	GetSummarizeFields()
		{
		return .summarize_fields
		}
	GetSummarizeFuncCols()
		{
		return .summarize_func_cols
		}
	GetNonSummarizedFields()
		{
		return .nonsummarized_fields
		}
	GetQuery()
		{
		return .query
		}
	GetKeys()
		{
		return .keys
		}
	GetSummarizeJoinFields()
		{
		return .summarize_joinfields
		}
	Clear()
		{
		.query = ''
		.keys = #()
		.sf = SelectFields()
		.allcols = #()
		.design_cols = #()

		.summarize_func_cols = #()
		.nonsummarized_fields = #()
		.infoFields = Object()

		.title = 'Reporter'
		.save_name = ''
		}

//====================================================================
	Valid?: true
	Report(paramselects = #(), quiet = false, checkOnly = false)
		{
		.Valid? = true
		try
			{
			return .checkAndBuildQuery(paramselects, :checkOnly)
			}
		catch (errs, 'Reporter:')
			{
			.Valid? = false

			if errs =~ '; select tab:'
				errs = errs.BeforeLast(';')
			if quiet isnt true
				Alert(errs.Replace('Reporter: ', ''), 'Reporter', 0, MB.ICONWARNING)
			else
				SuneidoLog('Error: ' $ errs, calls:)

			// use a query that produces no results so we get an empty report
			report = Object(ReporterFormat, 'tables where table is ""',
				Columns: #(), Data: Record(), Sf: .sf, Totalfields: #(),
				Countfields: #(), Minfields: #(), Maxfields: #(),
				Averagefields: #(), Printsummaryfields: #(), exclude: #())
			return Object(report, name: .title,
				header: Object('PageHead', .data.heading1, title2: ''),
				printParams: [], paramsdata: [])
			}
		}

	checkAndBuildQuery(paramselects, checkOnly = false)
		{
		.check_report()
		return .BuildReport(.data, paramselects, :checkOnly)
		}

	check_report()
		{
		if .query is ''
			throw .errMsg isnt '' ? .errMsg : "Reporter: " $ Reporter.InvalidMsg()
		if .source.Member?('validField')
			{
			rec = .data.Copy()
			rec.Invalidate(.source.validField)
			errs = rec[.source.validField]
			if errs isnt ''
				throw 'Reporter: ' $ errs
			}
		}

	BuildReport(data, paramselects = #(), source = false, checkOnly = false)
		{
		.data = data
		sort = .make_sort()
		cols = .make_cols()
		try
			query = .make_query(cols, sort, paramselects, source, :checkOnly)
		catch(err/*unused*/, 'Retry failed -|RetryTransaction: too many retries')
			{
			throw 'Reporter: there was a problem generating the report, please try again'
			}
		summaryObject = .make_sort_summary(data)

		rpt = Object(.reporterMode is 'form' ? ReporterCanvasFormat : ReporterFormat,
			query, Columns: cols, Data: .data, Sf: .sf,
			Totalfields: .sf.PromptsToFields(summaryObject.total),
			Minfields: .sf.PromptsToFields(summaryObject.min),
			Maxfields: .sf.PromptsToFields(summaryObject.max),
			Averagefields: .sf.PromptsToFields(summaryObject.average),
			Countfields: .sf.PromptsToFields(summaryObject.count),
			Printsummaryfields: .sf.PromptsToFields(summaryObject.summary_prompt_cols),
			Headertext: .data.header, Footertext: .data.footer,
			exclude: .BuildExclude(source))

		params = .params(.sf, .data.GetDefault(#select, []))
		return Object(rpt, name: .title, header: .GetHeader(), footer: .GetFooter(),
			printParams: .printParams(params),
			paramsdata: params.Copy().Add(.data.print_lines, at: 'PrintLines'))
		}
	GetSortCols()
		{
		sort = ''
		for i in .. Reporter.SortRows
			{
			if .data['sort' $ i] is ''
				continue
			else if .sf.HasPrompt?(.data['sort' $ i])
				sort $= ', ' $ .sf.PromptToField(.data['sort' $ i])
			else
				throw "Reporter: Invalid Sort; select tab: Sort"
			}
		return sort.Split(',').Filter({ |x| x isnt '' }).Map!(
			{ Object(field: it.Trim(), prompt: .sf.FieldToPrompt(it.Trim())) })
		}
	make_sort()
		{
		sort = ''
		for i in .. Reporter.SortRows
			{
			if .data['sort' $ i] is ''
				continue
			else if .sf.HasPrompt?(.data['sort' $ i])
				sort $= ', ' $ .sf.PromptToField(.data['sort' $ i])
			else
				throw "Reporter: Invalid Sort; select tab: Sort"
			}
		sort = sort.Replace(',', ' sort' $
			(.data.reverse is true ? ' reverse' : ''), 1)
		return sort
		}
	make_cols()
		{
		cols = .valid_cols()
		if cols.Empty?()
			throw "Reporter: No valid columns in Design; select tab: Design"
		cols.Map!({|c| Object(text: .sf.PromptToField(c.text),
			width: c.GetDefault('width', false)) })
		return cols
		}
	valid_cols()
		{
		cols = .data.columns.Copy()
		cols.RemoveIf({|c| not .design_cols.Has?(c.text) })
		return cols
		}
	make_extend(formula_fields, where, checkOnly = false)
		{
		if not .failedConverts.Empty?()
			throw 'Reporter: No permission to the following fields:\r\n ' $
				.failedConverts.Map!(Datadict.PromptOrHeading).Join(', ')

		return Reporter_make_extend(.data, .sf, .ForEachCalc, formula_fields, where,
			.summarize_by_fields(.data), .summarize_fields, :checkOnly)
		}
	summarize_by_fields(data)
		{
		return data.summarize_by.Split(',').Map!(.summary_field_from_prompt)
		}
	make_query(cols, sort, paramselects, source = false, checkOnly = false)
		{
		menuParams = .preProcessMenuParams(paramselects)
		exclude_menu_options = menuParams.paramselects.NotEmpty?()
		extends = .make_extend(formula_fields = Object(), formula_where = Object(),
			:checkOnly)
		where = Select2(.sf).Where(.data.GetDefault('select', []),
			except: .summarize_fields.Copy().MergeUnion(formula_where),
			extra_dd: .calc_dd(), :exclude_menu_options)
		.menu_params_extend(extends, menuParams.paramselects)
		joinflds = where.joinflds.MergeUnion(menuParams.options)
		joinfields = Flatten(
			cols.Map({ it.text }),
			sort.Replace('sort ','').Replace('reverse ', '').Split(','),
			joinflds,
			formula_fields,
			.summarize_joinfields).Join(',')

		joinob = .sf.JoinsOb(joinfields, withDetails?:)

		query = .tableHint(joinob, source) $
			'(' $ BuildQueryJoin(.query, joinob) $  ')\n' $
			Reporter_extend_info.Extend(.infoFields, joinfields) $ '\n' $
			extends[0] $ '\n' $
			where.where $ '\n' $
			.menu_paramswhere_before_summarize(menuParams) $ '\n' $
			.menu_params_calc_where_before_summarize(menuParams.paramselects, extends[1],
				.summarize_fields) $ '\n' $
			.make_summarize() $ '\n' $
			extends[1] $ '\n' $
			Select2(.sf).Where(.data.GetDefault(#select, []), formula_where,
				extra_dd: .calc_dd()).where $ '\n' $
			Select2(.sf).Where(.data.GetDefault(#select, []), .summarize_fields,
				extra_dd: .calc_dd()).where $ '\n' $
			.menu_paramswhere_after_summarize(menuParams.paramselects) $ '\n' $
			.menu_params_calc_where(menuParams.paramselects, extends[1],
				.summarize_fields) $
			sort
		// TODO: don't allow user to sort by columns that aren't in summary or not
		// in extends after the summary (only if they are using summary option)
		Reporter_query.CheckQuery(query)
		return query
		}

	preProcessMenuParams(paramselects)
		{
		paramselects = Object?(paramselects) ? paramselects : Object()
		fields = .Menu_params_fields()
		options = Object()

		for field in fields.Copy()
			{
			.convertParamSelectsField(field, paramselects)

			paramField = paramselects.Member?(field $ .param_field_suffix)
				? field $ .param_field_suffix
				: paramselects.Member?(field)
					? field
					: ''

			if paramField is ''
				continue

			param = paramselects[paramField]
			if not Object?(param) or '' is param.GetDefault(#operation, '')
				continue

			field = .convertNameAbbrevToNum(field, param, fields, paramselects)
			options.Add(field)
			}
		return Object(:fields, :options, :paramselects)
		}

	convertParamSelectsField(field, paramselects)
		{
		if not field.Suffix?("?")
			return
		fd = field.RemoveSuffix("?") $ .param_field_suffix
		if not paramselects.Member?(fd)
			return
		paramselects[field $ .param_field_suffix] = paramselects[fd]
		}

	convertNameAbbrevToNum(field, param, fields, paramselects)
		{
		if ((false is numField = .sf.GetJoinNumField(field)) or fields.Has?(numField))
			return field

		if false is joinNums = GetForeignNumsFromNameAbbrevFilter(field, .sf,
			param.operation, param.value, param.value2)
			return field

		paramselects[joinNums.numField] = joinNums.nums.Empty?()
			? [operation: 'less than', value: '', value2: '']
			: [operation: 'in list', value: joinNums.nums, value2: '']
		fields.Replace(field, joinNums.numField)
		return joinNums.numField
		}

	tableHint(joinob, source)
		{
		if .query.Has?('/* tableHint: ')
			return ''

		rptSource = .source
		if .source.Empty?() and source isnt false
			rptSource = source
		tables = rptSource.GetDefault('tables', #())
		if not tables.Empty?()
			return '/* tableHint: ' $ tables[0] $ ' */ '

		tableHint = ''
		// call QueryGetTable to make sure that query is not a complex query
		if not joinob.Empty?() or not .summarize_joinfields.Empty?()
			tableHint = '/* tableHint: ' $ QueryGetTable(.query) $ ' */ '
		return tableHint
		}

	menu_params_extend(extends, paramselects)
		{
		for field in paramselects.Members()
			if .calc_field?(field) and not Object?(paramselects[field])
				for i in extends.Members()
					if extends[i].Has?(field)
						extends[i] = extends[i].Replace('extend ' $ field $ ' = .*\n',
							'extend ' $ field $ ' = ' $
								Display(paramselects[field]) $ '\n')

		}
	menu_paramswhere_before_summarize(menuParams)
		{
		fields = menuParams.fields.Difference(.summarize_fields)
		return .menu_paramswhere(fields, menuParams.paramselects)
		}
	menu_paramswhere_after_summarize(paramselects)
		{
		return .menu_paramswhere(.summarize_fields.Copy(), paramselects)
		}
	menu_paramswhere(args, paramselects)
		{
		args.data = paramselects
		args.no_encode = true
		return GetParamsWhere(@args)
		}
	menu_params_calc_where_before_summarize(paramselects, extends, summarize_fields)
		{
		fields = Object()
		for field in paramselects.Members()
			if .calc_field?(field) and Object?(paramselects[field]) and
				(not summarize_fields.Has?(field) and not extends.Has?(field))
				fields.AddUnique(field)

		return fields.Empty?() ? '' : .menu_paramswhere(fields, paramselects)
		}
	menu_params_calc_where(paramselects, extends, summarize_fields)
		{
		fields = Object()
		for field in paramselects.Members()
			if .calc_field?(field) and Object?(paramselects[field]) and
				(summarize_fields.Has?(field) or extends.Has?(field))
				fields.AddUnique(field)

		return fields.Empty?() ? '' : .menu_paramswhere(fields, paramselects)
		}
	calc_dd()
		{
		dd = Object()
		.ForEachCalc()
			{ |key, field /*unused*/|
			type = .formulaDataType(key)
			if false isnt fieldDD = Reporter_datadict.GetCalcDD(type)
				dd['calc' $ key] = fieldDD
			}
		return dd
		}
	calc_field?(fd)
		{
		return fd =~ .Calc_prefix
		}
	GetFormulaType(fd)
		{
		if not .calc_field?(fd)
			return ''
		key = fd.AfterLast('calc')
		return .formulaDataType(key)
		}

	formulaDataType(key)
		{
		idx = .data.formulas.FindIf({ it.key is key })
		return .data.formulas[idx].GetDefault('type', 'Unknown')
		}

	func_map: (total: 'total', average: 'average', maximum: max, minimum: min)
	make_summarize()
		{
		summarize = .summarize_by_fields(.data)
		nby = summarize.Size()
		for i in .. Reporter.MaxSummarizeFields
			{
			if "" is func = .data['summarize_func' $ i]
				continue
			prompt = .data['summarize_field' $ i]
			if func is "count"
				summarize.Add("count")
			else if prompt isnt ""
				summarize.Add(.func_map[func] $ ' ' $ .summary_field_from_prompt(prompt))
			}
		op = summarize.Size() is nby
			? 'project /*CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/'
			: 'summarize'
		return summarize.Empty?() ? '' : op $ ' ' $ summarize.Join(',')
		}

	summary_field_from_prompt(prompt)
		{
		summary_field = .sf.PromptToField(prompt)
		if summary_field is false
			throw "Reporter: Unable to generate report. Invalid summary field(s)."
		return summary_field
		}

	make_sort_summary(data)
		{
		sortSummaryObject = Object(total: Object(), min: Object(), max: Object(),
			average: Object(), count: Object(), summary_prompt_cols: Object())
		for item in data.coloptions.Members()
			if .design_cols.Has?(item)
				{
				if data.coloptions[item].Member?('print_summary_prompts') and
					data.coloptions[item].print_summary_prompts is true
					sortSummaryObject.summary_prompt_cols.Add(item)
				for sortSummaryItem in #('total', 'min', 'max', 'average', 'count')
					if data.coloptions[item].Member?(sortSummaryItem)
						if data.coloptions[item][sortSummaryItem] is true
							sortSummaryObject[sortSummaryItem].Add(item)
				}
		return sortSummaryObject
		}

	Calc_prefix: '^calc\d|^total_calc\d|^max_calc\d|^min_calc\d|^average_calc\d'
	printParams(params)
		{
		return params.Members().
			Map!({ |fd| .calc_field?(fd) ? .sf.FieldToPrompt(fd) : fd })
		}
	params(sf, sel)
		{
		params = Object()
		for (i = 0; i < Select2Control.Numrows; ++i)
			{
			if sel['checkbox' $ i] is true and
				sel['oplist' $ i] isnt "" and sel['fieldlist' $ i] isnt "" and
				sel['print' $ i] is true
				{
				if false is field = sf.PromptToField(sel['fieldlist' $ i])
					throw 'Reporter: Invalid field in Select'
				if .calc_field?(field)
					field = sel['fieldlist' $ i]

				.format_params(params, field, sel['oplist' $ i], sel['val' $ i], i)
				}
			}
		return params.Set_default('')
		}

	format_params(params, field, op, value, i)
		{
		member = params.Member?(field) ? field $ '_' $ i : field
		params[member] = Object(operation: op, :value)
		}

	// need a quiet parameter so that the AutoUpdate_ReportMenuOptions can run
	// without getting alerts from the .Report method
	// checkOnly will not recreate formula records
	BuildReportText(quiet = false, noHeader = false, checkOnly = false, report = false)
		{
		if report is false
			report = .Report(:quiet, :checkOnly)
		reporterOb = Object('Params',
			Object(.reporterMode is 'form' ? 'ReporterCanvasFormat' : 'ReporterFormat',
				.save_name),
			Params: .menu_params(:checkOnly),
			printParams: .Menu_print_params(),
			header: .reporterMode is 'form'
				? false
				: report.GetDefault('header', false)
			name: report.name.Replace('Reporter', 'Reporter Report'),
			title: .data.report_name,
			disableFieldProtectRules:)
		if false isnt footer = report.GetDefault('footer', false)
			reporterOb.footer = footer
		// form needs to set header to false to avoid printing any headers
		if noHeader and .reporterMode isnt 'form'
			reporterOb.Delete('header')
		return reporterOb
		}

	menu_params(checkOnly = false)
		{
		sel = .data.GetDefault(#select, [])
		params = Object('Vert')
		for (i = 0; i < Select2Control.Numrows; ++i)
			{
			prompt = sel['fieldlist' $ i]
			if prompt isnt "" and sel['menu_option' $ i] is true
				.addFieldParam(prompt, params, :checkOnly)
			}
		.addFormulaParams(params)
		if .data.report_description isnt ''
			params.Add('Skip', Object('StaticWrap', .data.report_description))
		params.Add('Skip',
			#('Static',
				'Note:  There may be additional Select options in Reporter.',
				weight: 'bold'))
		return params
		}

	addFieldParam(prompt, params, checkOnly = false)
		{
		if false is (field = .sf.PromptToField(prompt))
			return
		originalField = field
		field = .abbrev_to_num_field(field, .sf.GetConverted())
		// handle calc fields (no field definition)
		if .calc_field?(field)
			field = .calc_param(field, prompt)
		if field isnt false
			{
			paramField = Object?(field) or .calc_field?(field)
				? field
				: .getParamField(field, originalField, prompt, :checkOnly)
			params.Add(Object('ParamsSelect', paramField))
			}
		}

	getParamField(field, originalField, prompt, checkOnly = false)
		{
		baseField = originalField
		prefix = ''
		suffix = ''
		if baseField.BeforeFirst('_') in ('total', 'min', 'max', 'average')
			{
			baseField = baseField.AfterFirst('_')
			prefix = prompt.BeforeFirst(' ')
			}

		if baseField.Prefix?(Reporter_extend_info.Prefix)
			{
			suffix = baseField.AfterLast('_')
			baseField = baseField.RemovePrefix(Reporter_extend_info.Prefix).
				BeforeLast('_')
			promptMethod = 'Reporter_extend_info.PrefixPrompt'
			}
		else
			{
			if false isnt numField = .sf.GetJoinNumField(baseField)
				{
				suffix = prompt.Suffix?('Name')
					? 'Name'
					: prompt.Suffix?('Abbrev')
						? 'Abbrev'
						: ''
				baseField = numField
				}

			promptMethod = not baseField.Prefix?('custom_') or
				prompt.Has?(SelectPrompt(baseField))
				? 'SelectPrompt'
				: 'Prompt'
			}

		return .buildParam(field, [:baseField, :promptMethod, :prefix, :suffix],
			:checkOnly)
		}

	buildParam(field, promptInfo, checkOnly = false)
		{
		return Reporter_datadict.BuildParam(field, .param_field_suffix, :promptInfo,
			:checkOnly)
		}

	addFormulaParams(params)
		{
		.ForEachCalc()
			{ |key, field /*unused*/, row|
			formula = row
			if formula.form_val is true
				params.Add(Object('Pair', Object('Static', formula.calc),
					.calc_param('calc' $ key, formula.calc)))
			}
		}

	param_field_suffix: '_param'
	MakeDD(field, prompt, baseType)
		{
		return Reporter_datadict.MakeDD(field, prompt, baseType)
		}

	Menu_print_params()
		{
		sel = .data.GetDefault(#select, [])
		fields = Object()
		for (i = 0; i < Select2Control.Numrows; ++i)
			{
			prompt = sel['fieldlist' $ i]
			if prompt isnt "" and sel['menu_option' $ i] is true and
				sel['print' $ i] is true
				{
				if false is field = .sf.PromptToField(prompt)
					continue
				if .calc_field?(field) // for printParams from PageHeadFormat
					param = Object(paramField: field, paramPrompt: prompt,
						paramFormat: .getCalcFormat(prompt))
				else
					param = .abbrev_to_num_field(field,
						.sf.GetConverted()).RemoveSuffix('?') $ .param_field_suffix
				fields.Add(param)
				}
			}
		return fields
		}
	getCalcFormat(prompt)
		{
		if false is calc = .getCalcField(prompt, .data)
			return false

		paramFormat = false
		for c in CustomFieldTypes(reporter:)
			{
			if calc.type is c.name
				{
				type = .formulaDataType(calc.key)
				paramFormat = Reporter_datadict.GetCalcDef(type, c)
				break
				}
			}
		return paramFormat
		}
	getCalcField(prompt, data)
		{
		if "" is formulas = data.formulas
			return false

		if false isnt (idx = formulas.FindIf({ it.calc is prompt }))
			return formulas[idx]

		if false is pos = data.summarize_func_cols.Find(prompt)
			return false

		val = data['summarize_field' $ pos]
		if false isnt idx = formulas.FindIf({ it.calc is val })
			return formulas[idx]

		return false
		}
	Menu_params_fields()
		{
		sel = .data.GetDefault(#select, [])
		fields = Object()
		for (i = 0; i < Select2Control.Numrows; ++i)
			if sel['fieldlist' $ i] isnt "" and sel['menu_option' $ i] is true
				if false isnt field = .sf.PromptToField(sel['fieldlist' $ i])
					fields[.abbrev_to_num_field(field, .sf.GetConverted())] = ''
		return .printParams(fields)
		}

	abbrev_to_num_field(field, convertedob = #())
		{
		if field.Suffix?('_abbrev') or field.Has?('_abbrev_')
			{
			numField = field.Replace('_abbrev', '_num')
			if convertedob.Member?(field)
				numField = convertedob[field]
			if Reporter_query.HasColumn?(.query, numField)
				return numField
			return field
			}
		return field
		}

	calc_param(field, prompt)
		{
		idx = .data.formulas.FindIf({ it.key is field.Extract('\d\d*')})
		type = .data.formulas[idx].type
		ctrl = false
		for c in CustomFieldTypes(reporter:)
			if type is c.name
				{
				ctrl = CustomFieldTypes.GetControl(c)
				break
				}
		if not Object?(ctrl)
			return false

		ob = ctrl.Copy()
		ob.name = field
		ob.prompt = prompt
		return ob
		}
	}