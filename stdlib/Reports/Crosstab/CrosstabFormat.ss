// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// BUG: column and row averages are averages of averages
QueryFormat
	{
	New(query, .rows = "", .cols = "", .value = "", .func = "", .sortcols = false,
		.selectFields = false)
		{
		super(@.setup(query))
		}

	setup(query)
		{
		.query = Query.StripSort(query)
		.accum = .func is #count ? #total : .func
		.warns = Object()
		return #()
		}

	Query()
		{
		.columns = Object()
		query = QueryHelper.ExtendColumns(.query, .selectFields, .collectFields())
		if .cols isnt ""
			.buildColumns(query)

		.sort_columns()
		if .func is #count
			.value = ""
		else if .value is ""
			throw "CrossTable " $ .func $ " requires Value field"
		.funcfield = .func is #count
			? #count
			: .func $ '_' $ .value

		if .rows isnt ""
			query $= " extend rowfield = " $ .rows
		query $= " summarize "
		sort = ""
		if .rows isnt ""
			{
			query $= " rowfield, "
			sort = " sort rowfield"
			}
		if .cols isnt ""
			query $= .cols $ ", "
		query $= .func $ ' ' $ .value $ sort
		return query
		}

	collectFields()
		{
		fields = Object()
		if .rows isnt ""
			fields.Add(.rows)
		if .cols isnt ""
			fields.Add(.cols)
		if .value isnt ""
			fields.Add(.value)
		return fields
		}

	columnLimit: 1000
	buildColumns(query)
		{
		.columns = _tran.QueryList(query, .cols)
		if .columns.Size() > .columnLimit
			{
			.warns.Add("too many column values (max " $ .columnLimit $ ')')
			.columns = .columns.Take(.columnLimit)
			}
		}

	sort_columns()
		{
		if .sortcols is false
			return
		else if .sortcols is true
			.columns.Sort!()
		else if Object?(.sortcols)
			{
			invalid = .sortcols.Difference(.columns)
			if not invalid.Empty?()
				throw String(invalid) $ " not found in Crosstab column: " $ .cols
			.columns = .sortcols.Copy().MergeUnion(.columns)
			}
		else
			try
				.columns.Sort!(.sortcols) /* user-defined sort fn? */
			catch /* if not a callable block or function */
				throw "Crosstab[sortcols] must be a list, " $
					"boolean, or valid .Sort!() function"
		}

	Before()
		{
		.totals = Object(last: Accumulator(.accum))
		.row = Object(last: Accumulator(.accum)).Set_default("")
		return false
		}

	Header()
		{
		f = Datadict(.cols).Format.Copy()
		f.justify = #center
		fmt = [#Row, [#Text, ""]]
		fmt.font = .font
		for i in .columns.Members()
			{
			f = f.Copy()
			f.data = .columns[i]
			f.abbrev = true // for IdFormat
			fmt.Add([#Vert, f, [#Hline, .widths[i+1], xstretch: 0]])
			}
		fmt.Add([#Vert,
			[#Text, .func.Capitalize(), justify: #center],
			[#Hline, .widths.Last(), xstretch: 0]])
		fmt.widths = .widths
		return fmt
		}

	BeforeOutput(data)
		{
		val = data[.funcfield]

		if .cols isnt ""
			{
			i = .columns.Find(data[.cols])
			col = "col" $ i
			.row[col] = val
			if not .totals.Member?(col)
				.totals[col] = Accumulator(.accum)
			.totals[col].Value(val)
			}
		.row.last.Value(val)
		.totals.last.Value(val)
		return false
		}

	Output()
		{
		fmt = [#RowHead]
		fmt.Add(.rows)
		f = .func is #count
			? #(Number, mask: "###,###")
			: Datadict(.value).Format
		for i in .columns.Members()
			{
			f = f.Copy()
			f.field = "col" $ i
			fmt.Add(f)
			}
		f = f.Copy()
		f.field = #last
		fmt.Add(f)

		cf = _report.Construct(fmt)
		.widths = cf.GetWidths()
		.font = cf.GetFont()
		.rowfmt = fmt.Copy()
		.rowfmt[0] = #Row
		return fmt
		}

	AfterOutput(data)
		{
		return .cols is "" ? .output(data) : false
		}

	After_rowfield(data)
		{
		return .cols isnt "" ? .output(data) : false
		}

	output(data)
		{
		.row.last = .row.last.Result()
		.row[.rows] = data.rowfield
		.rowfmt.data = .row
		.row = Object(last: Accumulator(.accum)).Set_default("")
		return .rowfmt
		}

	After()
		{
		if .rows is "" and .cols is ""
			return false
		fmt = ["_output"]
		fmt[.rows] = [#Vert,
			[#Vskip, 100.TwipsInInch()], /*= 1 hundredth of a twip*/
			[#Text, .func.Capitalize()]]
		f = .func is #count
			? #(Number, mask: "###,###")
			: Datadict(.value).Format
		for i in .columns.Members()
			{
			f = f.Copy()
			f.field = "col" $ i
			fmt["col" $ i] = [#Total, f]
			}
		f = f.Copy()
		f.field = #last
		fmt.last = [#Total, f]
		for col in .totals.Members()
			.totals[col] = .totals[col].Result()
		fmt.data = .totals
		return fmt
		}

	AfterAll()
		{
		if .warns.NotEmpty?()
			.Append([#Horz, #Hfill,
				[#Text, "*** " $ .warns.Join("; ") $ " ***"
					font: #(name: Arial, size: 12, weight: bold)], #Hfill])
		return false
		}
	}
