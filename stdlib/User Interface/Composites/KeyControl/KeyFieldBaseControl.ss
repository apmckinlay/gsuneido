// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	Name: 'KeyFieldBase'

	New(width, readonly = false, style = 0, status = '',
		font = '', size = '', weight = "", tabover = false, hidden = false,
		.addAccessOption = false, .mandatory = false, .allowOther = false,
		.noClear = false, .fillin = false, .from = false, query = '',
		.restrictions = false, .invalidRestrictions = false, upper = false, lower = false,
		.optionalRestrictions = #())
		{
		super(:width, :readonly, :style, :status, :font, :size, :weight, :tabover,
			:hidden, :upper, :lower)

		if addAccessOption is true
			{
			.AddContextMenuItem('', false)
			.AddContextMenuItem('Access', .On_Access)
			}

		.InitQuery(query)
		.fillincustom = .getFillinCustom()

		.SubClass()
		}
	GetLowerFieldName(fieldName)
		{
		lowerName = fieldName $ '_lower!'
		.lowerField = QueryColumns(.GetQuery()).Has?(lowerName)
			? lowerName : false
		}
	NameMatchFieldAndValue(field, val)
		{
		matchField = .lowerField isnt false ? .lowerField : field
		matchVal = matchField.Suffix?('_lower!') ? val.Lower() : val
		return Object(field: matchField, value: matchVal)
		}
	InitQuery(query)
		{
		.Setquery(query)
		}

	BuildQuery(query, restrictionsValid = false)
		{
		restriction = restrictionsValid ? false : .restrictions
		optRestrictions = restrictionsValid ? #() : .optionalRestrictions
		.Send("Key_BuildQuery", query, restriction, .invalidRestrictions,
			:optRestrictions)
		}

	Setquery(query)
		{
		.query = query
		}

	GetQuery()
		{
		return .query
		}

	Last_ellipsis: false
	KillFocus()
		{
		if (.Dirty?() and
			(.Last_ellipsis is false or .Get() isnt .Last_ellipsis))
			.Process_newvalue(userTyped:)
		.Last_ellipsis = false
		}

	SetFieldRestrictions(restrictions)
		{
		.restrictions = restrictions
		}

	Mandatory?()
		{
		return .mandatory is true
		}

	AllowOther?()
		{
		return .allowOther is true
		}

	getFillinCustom()
		{
		toQuery = .Send('GetTransQuery')
		return .getCustomFillin(.query, toQuery)
		}

	getCustomFillin(query, fillinToQuery)
		{
		if not String?(fillinToQuery) or not String?(query)
			return false
		if QueryGetTable(query, nothrow:) is QueryGetTable(fillinToQuery, nothrow:)
			return false
		return CustomizableMap(query, fillinToQuery, .GetDefault("Custom", false),
			.Parent.Name)
		}

	Fillin_fields(rec)
		{
		if .noClear and rec is Object()
			return

		.handleCustomFillin(rec)

		if .AllowFillin?(rec)
			.forFillins(.fillin)
				{ |i|
				.Send("SetField", .fillin[i], rec[.from is false ? .fillin[i] : .from[i]])
				}
		}

	forFillins(fillin, block)
		{
		for (i in fillin.Members())
			block(i)
		}

	ValidData?(@args)
		{
		NameArgs(args, #(val query field mandatory allowOther), #(false false))
		if args.val is ""
			return not args.mandatory
		if args.allowOther
			{
			maxCharacters = args.GetDefault('maxCharacters', .MaxCharacters)
			return .ValidTextLength?(args.val, maxCharacters)
			}

		if args.Member?('whereField') and args.whereField isnt false
			{
			qOb = Object(QueryStripSort(.findQuery(args.query)))
			qOb[args.whereField] = args.record[args.whereField]
			qOb[args.field] = args.val
			return false isnt Query1(@qOb)
			}

		restrictions = ""
		if args.Member?('restrictions') and
			args.GetDefault('restrictionsValid', true) is false
			restrictions = ' where ' $
				(args.restrictions =~ "^[a-zA-Z0-9_]+[?!]?$" is true
					? Record()[args.restrictions]
					: args.restrictions)

		return false isnt Query1(QueryAddWhere(.findQuery(args.query),
			' where ' $ args.field $ ' is ' $ Display(args.val) $
			restrictions))
		}
	findQuery(query)
		{
		if Function?(query)
			return query()
		if String?(query)
			return query
		throw 'Invalid Query Specification'
		}

	AllowFillin?(rec /*unused*/)
		{
		// uses current record, fillin and from objects
		return .fillin isnt false
		}

	handleCustomFillin(rec)
		{
		if .fillincustom is false or (.allowOther is true and rec.Empty?())
			return
		.fillincustom.ForEachField()
			{ |from, to, trim|
			.Send("SetMainRecordField", to, .fillincustom.TrimValue(rec[from], trim))
			}
		}

	On_Access()
		{
		.Send('KeyIdField_Access')
		return 0
		}
	// Helper method so Customizable can use standard fillin logic
	FillinRecord(query, valfield, fillin, from, rec, value, fillincustom = false)
		{
		if false is masterRec = Query1(QueryAddWhere(query,
			" where " $ valfield $ " is " $ Display(value)))
			return

		if fillincustom isnt false
			fillincustom.CopyCustomFields(rec, masterRec)

		if fillin isnt false and from isnt false
			.forFillins(fillin)
				{ |i|
				rec[fillin[i]] = masterRec[from[i]]
				}
		}
	}
