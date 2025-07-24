// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Generator
	{
	Export: false
	grandTotalRecord: 999
	New(@args)
		{
		for m in args.Members()
			if m is 0
				.Query = args[0]
			else
				this[m.Capitalize()] = args[m]
		_report.RegisterForClose(this)
		.Params = _report.Params

		query = .Val_or_func('Query')
		if args.Member?('bookLogQuery') and args.bookLogQuery is true
			BookLog(query)

		if .Member?('Inputparams')
			query = QueryAddWhere(query, _report.GetParamsWhere(@.Inputparams))
		if not _report.Member?('MainQueryFormat')
			_report.MainQueryFormat = this
		.query = query
		.q = .Open(query)

		.output = .Val_or_func('Output')
		.outputformat = .Val_or_func('OutputFormat')
		if .output is ''
			.output = .outputformat is false ? Object("Row").Add(@.q.Columns()) : false
		if .output isnt false
			.output = _report.Construct(.output)
		if .outputformat is false
			.outputformat = .output
		else
			.outputformat = _report.Construct(.outputformat)

		.order = .q.Order()

		.prev = false
		.tryQNext()
			{
			if false is .next = .nextRecord()
				.Close()
			}


		.totals = Object()
		.counts = Object()
		for o in .order.Members()
			{
			.totals[o] = Object()
			.counts[o] = Object()
			}
		.totals[.grandTotalRecord] = Object() // grand totals
		.totalFields = .Val_or_func('Total')
		.counts[.grandTotalRecord] = Object() // grand count
		.countFields = .Val_or_func('Count')
		}

	tryQNext(block)
		{
		try
			block()
		catch (unused, 'summarize too large')
			{
			Alert(.getSummarizeMsgWarning(this.Base?(CrosstabFormat)),
				'Unable to open report')
			.next = false
			.Close()
			}
		}

	summarizeWarnMsg: "Too many records to summarize.  Please use #REPLACE# to reduce " $
		"the number of records"
	getSummarizeMsgWarning(baseCrosstab?)
		{
		msg = .summarizeWarnMsg
		replacement = baseCrosstab?
				? 'the "Select..." button'
				: 'filters'
		return 'Unable to preform action.\r\n' $ msg.Replace('#REPLACE#', replacement)
		}

	Open(query) // overridden by ObjectFormat
		{
		return _tran.Query(query)
		}
	More()
		{
		do
			{
			.list = Object()
			if .next is false
				{
				.next = "alldone"
				.add("AfterAll", #{})
				if not .list.Empty?()
					{
					for fmt in .list
						super.Output(fmt) // Generator Output
					.list.Delete(all:)
					return
					}
				}
			if .next is "alldone"
				{
				_report.UnregisterForClose(this)
				.Close()
				return
				}
			data = .next
			.next = .nextRecord()

			// before's
			brk = false
			if .prev is false
				{
				.reset_summary(.grandTotalRecord)
				for o in .order.Members()
					.reset_summary(o)
				brk = true
				.add("Before", data)
				}
			for i in .order.Members()
				{
				fld = .order[i]
				if not brk and data[fld] isnt .prev[fld]
					brk = true
				if brk
					{
					.reset_summary(i)
					.add("Before_" $ fld, data)
					}
				}

			.add("BeforeOutput", data)

			.accum_summary(data)

			if .output isnt false
				{
				out = .output.Copy()
				out.Data = data
				.list.Add(out)
				}

			.add("AfterOutput", data)

			if .next is false
				n = 0
			else
				{
				for (n = 0; .order.Member?(n); ++n)
					{
					fld = .order[n]
					if data[fld] isnt .next[fld]
						break
					}
				}
			for (i = .order.Size() - 1; i >= n; --i)
				.add("After_" $ .order[i], .merge_sort_summary(data.Copy(), i))
			if .next is false
				.add("After", .merge_sort_summary(data, .grandTotalRecord))
			.prev = data
			}
			while (.list.Empty?())

		for (fmt in .list)
			super.Output(fmt); // Generator Output
			// so lack of output isnt taken as eof
		.list.Delete(all:)
		}

	alreadyDetectedSlow: false
	nextRecord()
		{
		if .skipSlowQueryChecking()
			return .q.Next()
		t = Timer()
			{
			rec = .q.Next()
			}
		SlowQuery.LogIfTooSlow(t, .query)
			{ |hash|
			.alreadyDetectedSlow = true
			SuneidoLog('ERRATIC: (' $ hash $ ')' $
				' Report query takes more than 5 seconds',
				params: [
					queryHead: .query[..99], queryTail: .query[-99..]], /*= to not trim */
				switch_prefix_limit: 3)
			}
		return rec
		}

	skipSlowQueryChecking()
		{
		if not String?(.query)
			return true
		if .alreadyDetectedSlow
			return true
		rptArgs = _report.GetReportArgs()
		if false is rptArgs.GetDefault('slowQueryFilter', false)
			return true
		if rptArgs.GetDefault('suppressSlowQuery', false)
			return true
		if _report.Member?('MainQueryFormat')
			return _report.MainQueryFormat isnt this
		return false
		}

	reset_summary(o)
		{
		totals = .totals[o]
		counts = .counts[o]

		for t in .totalFields.Members()
			totals[t] = Object?(.totalFields[t]) ? new MoneyBag : 0

		for c in .countFields.Members()
			counts[c] = 0
		}

	// totalling
	accum_totals(data)
		{
		.accum_totals1(data, .grandTotalRecord)
		for o in .order.Members()
			.accum_totals1(data, o)
		}
	accum_totals1(data, o)
		{
		for t in .totalFields.Members()
			if Object?(ob = .totalFields[t])
				if Instance?(data[ob[0]])
					.totals[o][t].PlusMB(data[ob[0]])
				else
					.totals[o][t].Plus(data[ob[0]], data[ob[1]])
			else if Numberable?(data[.totalFields[t]])
				.totals[o][t] += Number(data[.totalFields[t]])
		}

	accum_counts1(data /*unused*/, o)
		{
		for c in .countFields.Members()
			.counts[o][c] += 1
		}

	accum_summary(data)
		{
		.accum_totals1(data, .grandTotalRecord)
		.accum_counts1(data, .grandTotalRecord)
		for orderMember in .order.Members()
			{
			.accum_totals1(data, orderMember)
			.accum_counts1(data, orderMember)
			}
		}

	merge_sort_summary1(data, type, synop, fields)
		{
		for field in synop.Members()
			{
			if Object?(fld = fields[field])
				fld = fld[0]
			data[type $ fld] = synop[field]
			}
		return data
		}
	merge_sort_summary(data, o)
		{
		if .totalFields isnt #()
			data = .merge_sort_summary1(data, "total_", .totals[o], .totalFields)
		if .countFields isnt #()
			data = .merge_sort_summary1(data, "count_", .counts[o], .countFields)
		return data
		}

	AddToTotals(data)
		{
		.accum_totals(data)
		}

	subtract_totals(data)
		{
		.subtract_totals1(data, .grandTotalRecord)
		for o in .order.Members()
			.subtract_totals1(data, o)
		}
	subtract_totals1(data, o)
		{
		for t in .totalFields.Members()
			if Object?(ob = .totalFields[t])
				if Instance?(data[ob[0]])
					.totals[o][t].MinusMB(data[ob[0]])
				else
					.totals[o][t].Minus(data[ob[0]], data[ob[1]])
			else
				.totals[o][t] -= data[.totalFields[t]]
		}
	SubtractFromTotals(data)
		{
		.subtract_totals(data)
		}
	GetTotals()
		{
		return .totals
		}
	GetCounts()
		{
		return .counts
		}
	add(brk, data)
		{
		.append_data = data
		fmt = false
		if .Method?(brk)
			fmt = this[brk](:data)
		else if .Member?(brk)
			fmt = this[brk]
		else if brk.Prefix?("Before_") or brk.Prefix?("After_")
			fmt = this[brk.BeforeFirst('_') $ '_'](brk.AfterFirst('_'), :data)
		if fmt is false
			return
		.Append(fmt)
		}
	Append(@args)
		{
		data = args.Member?('data') ? args.data : .append_data
		for fmt in args.Values(list:)
			{
			if Object?(fmt) and fmt.Member?(0) and String?(fmt[0]) and fmt[0][0] is "_"
				{
				name = fmt[0][1 ..].Capitalize()
				fmt = .outputformat[name](fmt, data)
				}
			else
				{
				fmt = _report.Construct(fmt)
				if Instance?(fmt) and fmt.Member?("Data") and fmt.Data is false
					fmt.Data = (fmt.Member?("Field") ? data[fmt.Field] : data)
				}
			.list.Add(fmt)
			}
		}
	// default methods
	Header()
		{
		return .outputformat is false ? false : .outputformat.Header()
		}
	Output: ''
	OutputFormat: false
	Total: ()
	Count: ()
	Before_(field /*unused*/)
		{ return false }
	After_(field /*unused*/)
		{ return false }
	Close() // overridden by ObjectFormat
		{
		// Issue 35186 check if member exists to prevent errors closing when query error
		// is caught during initialization
		if not .Member?('q') or .q is false
			return
		.q.Close()
		.q = false
		}
	}
