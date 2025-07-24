// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
/*
abstract base class
	- handles opening and closing the textfile
	- iterates through the query
	- returns the number of records imported
provides
	From_file - the name of the source text file
	To_query - the name of the destination query (usually just a table)
	Fields - the list of fields
	Tf - the open textfile
	Getline()
	Output(x) - output the records
derived classes define:
	Header() - optional, for getting file headers
	Before() - optional, used for outputing at the beginning of the file
	Import1(x) - mandatory, process a single record
*/
class
	{
	CallClass(@args)
		{
		import = new this(@args)
		try
			{
			import.Import()
			import.Close()
			return import.N
			}
		catch (e)
			{
			import.Close()
			throw e
			}
		}
	New(from_file, to_query, fields = false, header = true, .convertCaps = false)
		{
		.From_file = from_file
		.To_query = to_query
		.Tf = from_file isnt false ? File(from_file, 'r') : false
		.Fields = fields
		.Header? = header is true
		}
	Header?: false
	convertCaps: false
	DateFmt() { return false }
	Import()
		{
		try
			{
			.Before()
			if (.Header?)
				.Header()
			.DoImport()
			}
		catch (err, 'File: Readline: line too long')
			.readlineFailure(err)
		}

	DoImport()
		{
		while (.ImportEach()) {}
		}

	ImportEach()
		{
		if false is line = .Getline()
			return false

		if false is x = .Import1(line)
			return false
		x = .ConvertRecord(x, .DateFmt())
		.BeforeOutputLine(x)
		.Output(x)
		return true
		}

	BeforeOutputLine(x /*unused*/)
		{ /* do nothing */ }
	ConvertRecord(x, dateFmt)
		{
		// WARNING: The new rec doen't have deps inheritted, so it may cause inconsistent
		// 			values if there is a field that is depended by another field and
		//			its value is modified by encode/conversion (see suggestion 28861).
		//			Could apply the conversion on the original x directly if the
		//			original x is a record.
		rec = Record()
		for field in .getRecordFields(x)
			{
			rec[field] = x[field]
			// number values may have been converted from string from Import1
			if String?(rec[field])
				rec[field] = rec[field].Trim()

			dd = Datadict(field)
			// Datadict returns Field_String if it cannot find "field".
			// Therefore IF dd is Field_String and "field" isn't "string",
			// we need to convert the field value manually (done via DefaultConversion)
			if dd is Field_string and field isnt 'string'
				.DefaultConversion(rec, field)
			else
				{
				rec[field] = dd.Encode(rec[field], fmt: dateFmt)
				.EncodeConversion(dd, rec, field)
				}

			if x[field] isnt "" and (field.Has?('_abbrev') or field.Has?('_name'))
				.LookupNumField(rec, field)
			}
		return rec
		}

	// need to sort fields so that num is processed first
	// if name/abbrev is set first it won't do a proper lookup for master rec
	getRecordFields(x)
		{
		fields = x.Members()
		normal = { not (it.Suffix?('_num') or it.Has?('_num_')) }
		return fields.SortWith!(normal)
		}

	EncodeConversion(dd, rec, field)
		{
		if Number?(rec[field])
			.roundNumericValues(rec, field, dd)

		if Date?(rec[field]) and dd.Member?('Control') and
			(dd.Control[0] is 'ChooseDate' or
			dd.Control[0] is 'ChooseDateControl')
			rec[field] = rec[field].NoTime()

		// apply proper case if applicable, can't use GetDefault on dd class
		if .convertCaps and dd.Member?('ProperCase?') and dd.ProperCase?
			.convertToProperCase(rec, field)

		rec[field] = .ConvertAbbrev(field, rec[field])
		}

	roundNumericValues(rec, field, dd)
		{
		mask = dd.Control.Member?('mask') ? dd.Control.mask : false
		if mask is false
			return

		mask = NumberControl.RetrieveMask(mask)
		decimals = mask.AfterLast('.').Size()
		rec[field] = rec[field].Round(decimals)
		}

	convertToProperCase(rec, field)
		{
		if (rec[field] is "")
			return

		orig = rec[field]
		// not CapitalizeWords if string is already mixed case
		if not (orig =~ "[A-Z]" and orig =~"[a-z]")
			rec[field] = rec[field].CapitalizeWords()
		}

	ConvertAbbrev(field, value)
		{
		if not field.Has?('_abbrev') or value is "" or .nonLowerAbbrevField?(field)
			return value

		return value.Lower()
		}

	nonLowerAbbrevField?(field)
		{
		return Datadict(field, #(NonLowerAbbrev?)).GetDefault('NonLowerAbbrev?', false)
		}

	DefaultConversion(rec, field)
		{
		// don't convert name and abbrev fields or the lookup for the num will fail
		if not field.Has?('_abbrev') and not field.Has?('_name')
			rec[field] = ConvertNumeric(rec[field])
		}

	// If importing directly into the table we cannot import invalid data
	// however, if importing from an application where users can preview
	// the data, we want them to be able to correct the data
	AllowInvalidLookupVals: false
	LookupNumField(x, field, allowEmpty? = false, filter = "")
		{
		numfield = field.Replace('_name|_abbrev', '_num')
		lookupVal = x[field]
		if lookupVal is '' and allowEmpty? is true
			{
			x[numfield] = ''
			return
			}

		if false is y = FindForeignRecWithAbbrevNameOrNum(x, field, :filter)
			{
			if Datadict(numfield).Control.GetDefault('allowOther', false) is true or
				.AllowInvalidLookupVals
				x[numfield] = lookupVal
			return
			}
		x[numfield] = y[field.BeforeFirst("_") $ '_num']
		}
	Before()
		{ }
	Header()
		{ }
	Import1(line /*unused*/)
		{
		throw "Import1 not defined"
		}
	readLineLimit: 4000
	Getline()
		{
		bTell = .Tf.Tell()
		line = .Tf.Readline()
		if .readLineLimit < (.Tf.Tell() - bTell)
			throw 'File: Readline: line too long'
		return line
		}

	AlertOnReadLineFailure: false
	readlineFailure(err, type = 'ERROR: (CAUGHT) ')
		{
		if .AlertOnReadLineFailure
			Alert('Invalid file - line length exceeded maximum', title: 'Load',
				flags: MB.ICONERROR)
		}

	N: 0
	GetTransaction()
		{
		return .t
		}
	t: false
	Output(x)
		{
		// start a new transaction every 100 records
		if (.N++ % 100 is 0) /*= new transaction threshold */
			{
			if (.t isnt false)
				.t.Complete()
			.t = Transaction(update:)
			.q = .t.Query(.To_query)
			}
		count = 0
		forever
			{
			try
				{
				.q.Output(x)
				break
				}
			catch (err, "*duplicate key")
				.Make_unique(x, .To_query, err)
			if (++count > 100) /*= output attempt threshold */
				throw "error during import"
			}
		}
	Make_unique(x, query, err)
		{
		key = err.Extract("key: ([a-zA-Z0-9_]+!?)")
		if key.Suffix?('_lower!')
			key = .getActualKey(key, query)

		// if key is a "_name" field and key is empty, try abbrev
		abbrev = x[key[.. -4] $ "abbrev"] /*= trimming name*/
		if (x[key] is "" and key =~ "_name$" and abbrev isnt "")
			{
			x[key] = abbrev
			return
			}
		try
			keyvalue = String(x[key])
		catch
			throw err
		pattern = "\*([0-9]+)$"
		if (keyvalue =~ pattern)
			x[key] = keyvalue.Replace(pattern,
				'*' $ String(Number(keyvalue.Extract(pattern)) + 1))
		else
			x[key] $= "*1"
		}

	getActualKey(lowerKey, query)
		{
		origKey = lowerKey.BeforeFirst('_lower!')
		for key in QueryKeys(query)
			if key.Prefix?(origKey) and key isnt lowerKey
				return key
		return origKey
		}

	Close()
		{
		if .Tf isnt false
			{
			.Tf.Close()
			.Tf = false
			}
		if .t isnt false
			{
			.t.Complete()
			.t = false
			}
		}
	}