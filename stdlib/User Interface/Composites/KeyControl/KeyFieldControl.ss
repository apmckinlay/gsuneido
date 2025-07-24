// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
KeyFieldBaseControl
	{
	Name: "KeyField"
	New(query = "", .field = "", fillin = false, from = false,
		mandatory = false, allowOther = false, .abbrevField = false,
		width = 10, readonly = false, noClear = false, style = 0,
		status = '', font = "", size = "", weight = "", addAccessOption = false,
		restrictions = false, invalidRestrictions = false, tabover = false,
		hidden = false, upper = false, lower = false, optionalRestrictions = #())
		{
		super(:width, :readonly, :style, :status, :font, :size, :weight, :tabover,
			:hidden, :addAccessOption, :mandatory, :allowOther, :noClear, :fillin, :from,
			:query, :restrictions, :invalidRestrictions, :upper, :lower,
			:optionalRestrictions)
		.mode = "infile"
		.GetLowerFieldName(.field)
		}
	Get()
		{
		val = super.Get()
		if (val !~ ',')
			val = DatadictEncode(.field, val)

		if .field.Suffix?('_lower!')
			val = val.Lower()

		return val
		}
	Set(value)
		{
		.valid = true
		super.Set(value)
		}
	lookup_val: false
	lookup_rec: false
	lookup_time: false
	lookup_record(val)
		{
		if .lookup_val is val and Date().MinusSeconds(.lookup_time) < 3 //*= seconds */
			return .lookup_rec
		// lookup value in table, applying whereField if applicable
		q = .GetQuery()
		if((whereField = .Send("GetWhereField")) isnt false)
			{
			wherevalue = .Send('GetField', whereField)
			q = QueryAddWhere(q, " where " $ whereField $ " = " $ Display(wherevalue))
			}
		q = QueryAddWhere(q, " where " $ .field $ " = " $ Display(val))
		.lookup_rec = Query1(q)
		.lookup_val = val
		.lookup_time = Date()
		return .lookup_rec
		}

	Valid?(forceCheck = false)
		{
		if .GetReadOnly()
			return true
		val = .Get()
		if (val is "")
			return not .Mandatory?()
		if .AllowOther?()
			return .ValidLength?(val)

		if .mode is "unique"
			return .validUniqueMode(val)

		// if using wherefield, always do lookup in case wherevalue changed
		if .noChange?(val, forceCheck)
			return true
		rec = .lookup_record(val)
		.valid = rec isnt false ? rec[.field] : false
		return .valid isnt false
		}

	validUniqueMode(val)
		{
		newrec? = .Send("NewRecord?")
		if (newrec?)
			{
			// key value can't be in the table (would be duplicate)
			return not .keyexists?(val)
			}
		else
			{
			// old record,  Need to get original value from controller
			original_rec = .Send("GetOriginal")
			if not Object?(original_rec)
				return true
			if (val is original_rec[.field])
				return true
			return not .keyexists?(val)
			}
		}

	noChange?(val, forceCheck)
		{
		return .valid is val and .Send("GetWhereField") is false and not forceCheck
		}

	keyexists?(key)
		{
		return not QueryEmpty?(QueryAddWhere(.GetQuery(),
			" where " $ .field $ " is " $ Display(key)))
		}

	valid: true
	Process_newvalue()
		{
		x = .getrec()
		.valid = x isnt false ? x[.field] : false
		.Fillin_fields(Record?(x) ? x : Object().Set_default(""))
		}
	getrec()
		{
		val = .Get()
		if (val is "")
			return false
		//check prefix if String
		if (String?(val))
			{
			nameMatch = .NameMatchFieldAndValue(.field, val)
			if (false isnt (rec = .match_prefix(nameMatch.value, nameMatch.field)))
				return rec
			if (.abbrevField isnt false and
				false isnt (rec = .match_prefix(val, .abbrevField)))
				return rec
			}
		else // must be exact match for data types other than string
			{
			if (false isnt rec = Query1(QueryAddWhere(.GetQuery(),
				" where " $ .field $ " is " $ Display(val))))
				return rec
			}
		return false
		}
	match_prefix(val, field)
		{
		displayField = KeyDisplayField(.field)
		if (false isnt (rec = .lookup_record(val)))
			{
			if displayField isnt .field
				.Set(rec[displayField])
			return rec
			}
		record = false
		Transaction(read:)
			{ |t|
			query = .BuildQuery(QueryStripSort(.GetQuery()))
			q = t.Query(QueryAddWhere(query,
				" where " $ field $ " >= " $ Display(val)) $ " sort " $ field)
			if (false isnt (rec = q.Next()) and String(rec[field]).Prefix?(val))
				{
				if (rec[field] is val)
					{
					.Set(rec[displayField])
					record = rec
					}
				else
					{
					next = q.Next()
					if (next is false or not String(next[field]).Prefix?(val))
						{
						.Set(rec[displayField])
						record = rec
						}
					}
				}
			}
		return record
		}
	Setmode(mode)
		{
		.mode = mode
		}
	ChangeKey(newkey)
		{
		.field = newkey
		}
	}
