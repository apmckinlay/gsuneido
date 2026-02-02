// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
/*
abstract base class
	- defaults the text file name to the table name plus an extension
	- handles opening and closing the textfile
	- iterates through the query
	- returns the number of records exported
provides:
	From_query
	To_file
	Fields
	Tf - the open text file
	Putline(line)
derived classes define:
	Ext - optional, default file extension, default is "txt"
	Before() - optional, used for outputing at the beginning of the file
	Header() - optional, for outputing file headers
	Export1(x) - mandatory, output a single record
	After() - optional, for outputing at the end of the file
*/
class
	{
	Ext: "txt"
	Fields: false
	CallClass(@args)
		{
		e = new this(@args)
		n = e.Export()
		e.Close()
		return n
		}
	New(from_query, to_file = false, fields = false, header = "Fields", prompts = false)
		{
		.From_query = from_query
		.To_file = to_file isnt false ? to_file : from_query $ "." $ .Ext
		.OpenFile()
		.Fields = fields isnt false ? fields : QuerySelectColumns(from_query)
		.HeaderType = header
		.GetHead(prompts)
		}

	OpenFile()
		{
		.Tf = File(.To_file, 'w')
		}

	// Extracted for redefinition by derived classes (especially XML)
	GetHead(prompts)
		{
		if not Boolean?(.HeaderType) and .HeaderType.Has?("Prompt")
			{
			.Head = Object()
			for field in .Fields
				if prompts is false
					{
					if .HeaderType is "SelectPrompt"
						.Head.Add(SelectPrompt(field))
					else
						.Head.Add(Prompt(field))
					}
				else .Head.Add(prompts.Find(field))
			}
		else
			.Head = .convertHeaderNumFieldsToName()
		}

	convertHeaderNumFieldsToName()
		{
		headerFields = Object()
		for field in .Fields
			{
			if false isnt .numFieldInfo(field)
				headerFields.Add(field.Replace("_num", "_name"))
			else
				headerFields.Add(field)
			}
		return headerFields
		}

	HeaderType: "None"
	Export()
		{
		.Before()
		if .HeaderType isnt "None"
			.Header()
		n = 0
		QueryApply(.From_query)
			{ |x|
			for field in .Fields // need to kick in rules before formatting values
				x[field]
			.convertNumFieldsToName(x)
			.format_dates(x)
			.format_Encrypt(x)
			.Export1(x)
			++n
			}
		.After()
		return n
		}
	format_dates(x, fields = false)
		{
		if fields is false
			fields = .Fields

		for field in fields
			{
			fmt = .Datadict(field).Format[0]
			// don't format _num fields or they won't re-import properly
			// because they lose the second and millisecond portion
			if field.Has?("_num") and fmt is "Id"
				continue

			value = x[field]
			if Date?(value)
				x.PreSet(field, fmt is 'DateTime'
					? value.ShortDateTime()
					: value.ShortDate())
			}
		return x
		}
	convertNumFieldsToName(x)
		{
		for field in .Fields
			{
			namefield = field.BeforeFirst('_num') $ '_name'
			if false is rec = .GetMasterRecFromNumField(field, x[field])
				continue
			if rec[namefield] isnt ''
				x[field] = rec[namefield]
			}
		}
	GetMasterRecFromNumField(field, num)
		{
		if false is result = .numFieldInfo(field)
			return false
		dd = result.dd
		table = result.table
		numfield = dd.Control.field
		// catch non-existent column errors since we are not certain
		// that the num field exists, although it
		// should if it is being used as the IdControl's field
		rec = false
		try
			rec = Query1(table $ ' where ' $ numfield $ ' is ' $ Display(num))
		catch (err /*unused*/, "query: select: nonexistent columns") {}
		return rec
		}
	numFieldInfo(field)
		{
		if not field.Has?("_num")
			return false
		dd = .Datadict(field)
		if dd.Control[0] isnt 'Id'
			return false
		query = dd.Control[1]
		if Function?(query)
			query = query()
		table = query.BeforeFirst(' ')
		if table.Suffix?('/*')
			table = query.AfterFirst('*/ ').BeforeFirst(' ')
		if not dd.Control.Member?('field') or
			(false is Query1Cached('tables', :table) and
			false is Query1Cached('views', view_name: table))
			return false
		return Object(:dd, :table)
		}
	format_Encrypt(x)
		{
		for field in .Fields
			{
			dd = .Datadict(field)
			control = dd.Control[0]
			if control is 'SSNSIN' or dd.GetDefault('Format', [])[0] is 'SINSSN'
				{
				if x[field].Has?(": ")
					x[field] = x[field].AfterFirst(': ')
				x[field] = Decrypt(x[field])
				}
			else if control is 'Encrypt'
				x[field] = Decrypt(x[field])
			}
		}
	Datadict(field)
		{
		return Datadict(field)
		}
	Before()
		{ }
	Header()
		{ }
	Export1(unused)
		{
		throw "Export1 not defined"
		}
	Putline(line)
		{
		.Tf.Writeline(line)
		}
	After()
		{ }

	Close()
		{
		.Tf.Close()
		}
	}
