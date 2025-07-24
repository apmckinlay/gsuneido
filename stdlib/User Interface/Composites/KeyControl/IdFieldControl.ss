// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
KeyFieldBaseControl
	{
	Name: "IdField"
	New(query = "", field = "", fillin = false, from = false,
		mandatory = false, allowOther = false, noClear = false,
		width = 10, readonly = false,
		restrictions = false, nameField = false, abbrevField = false,
		restrictionsValid = true, style = 0, status = '',
		invalidRestrictions = false, whereFieldKey = false,
		font = "", size = "", weight = "", addAccessOption = false, tabover = false,
		hidden = '', optionalRestrictions = #())
		{
		super(:width, :readonly, :style, :status, :font, :size, :weight, :tabover,
			:hidden, :addAccessOption, :mandatory, :allowOther, :noClear, :fillin, :from,
			:query, :restrictions, :invalidRestrictions, :optionalRestrictions)

		.base_query = query // used for set
		.numField = field
		.nameField = nameField isnt false
			? nameField
			: field.Replace("(_num|_name|_abbrev)$", "_name")
		.abbrevField = abbrevField isnt false
			? abbrevField
			: field.Replace("(_num|_name|_abbrev)$", "_abbrev")
		.GetLowerFieldName(.nameField)
		.useAbbrevField? = QueryColumns(.GetQuery()).Has?(.abbrevField) and
			.abbrevField isnt .nameField
		.restrictionsValid = restrictionsValid
		.whereFieldKey = whereFieldKey
		}

	InitQuery(query)
		{
		.Setquery(.BuildQuery(query))
		}

	num: false
	Get()
		{
		// if dirty we have to call Process_newvalue method to get the record
		// and set .num to get the new record's timestamp. We do NOT want the
		// Process_newvalue to do the fillins (causes unwanted side-effects)
		if (.Dirty?())
			.Process_newvalue(fillin: false)
		val = .num isnt false ? .num : GetWindowText(.Hwnd)
		return String?(val) ? val.Trim() : val
		}
	set_val: false
	Set(num)
		{
		.num = .set_val = num
		.userTypedMultiMatch = false

		rec = .lookup_record(DatadictEncode(.numField, num),
			.numField, noRestrictions:, excludeWhereField:)

		value = rec is false and not .AllowOther?() and String?(num) and num isnt ""
			? .setInvalidValue(num)
			: .setValidValue(num, rec)

		.SelectAll()
		.SetValid(.valid = true)
		.Send("IdField_Set", num)
		}

	setInvalidValue(val)
		{
		.num = val
		SetWindowText(.Hwnd, val)
		return val
		}

	setValidValue(num, rec)
		{
		if rec is false
			{
			rec = Record(:num)
			rec[.nameField] = rec[.numField] = String?(.num) ? .num : '???'
			}
		return .set(rec)
		}

	lookup_record(val, field = false, noRestrictions = false, excludeWhereField = false)
		{
		if val is ""
			return false

		if (field is false)
			field = .nameField

		// lookup value in table, applying whereField if applicable
		// only use restrictions on lookup if restrictionsValid is false
		q = noRestrictions ? .base_query : .BuildQuery(.base_query, .restrictionsValid)

		// don't apply wherefield if lookup is for set value (saved vals should be valid)
		// - BUT, if the wherefield is required to make the lookup unique (key),
		// still must apply it
		if ((not excludeWhereField or .whereFieldKey) and
			(whereField = .Send("GetWhereField")) isnt false)
			{
			wherevalue = .Send('GetField', whereField)
			q = QueryAddWhere(q, " where " $ whereField $ " = " $ Display(wherevalue))
			}
		return Query1(QueryAddWhere(q, " where " $ field $ " = " $ Display(val)))
		}

	set(rec)
		{
		.num = rec[.numField]
		SetWindowText(.Hwnd, rec[.nameField])
		return rec[.nameField]
		}

	Valid?(forceCheck = false)
		{
		if .GetReadOnly()
			return true
		val = super.Get().Trim()
		if .empty?(.num, val)
			return not .Mandatory?()
		if .masterExist?(.num)
			return true
		if .userTypedMultiMatch isnt false
			return false
		if .AllowOther?()
			return .validAllowOther(.num, val)
		if .valid is val and not forceCheck
			return true
		return .forceValid(val)
		}

	empty?(num_val, val)
		{
		return num_val is "" or val is ""
		}

	masterExist?(num_val)
		{
		return num_val is .set_val and
			false isnt .lookup_record(num_val, .numField,
				noRestrictions:, excludeWhereField:)
		}

	validAllowOther(num_val, val)
		{
		if Date?(num_val) and
			false is .lookup_record(num_val, .numField,
				noRestrictions:, excludeWhereField:)
			return false
		return .ValidLength?(val)
		}

	forceValid(val)
		{
		rec = .lookup_record(val)
		.valid = rec isnt false ? rec[.nameField] : false
		return .valid isnt false
		}

	valid: true
	Process_newvalue(fillin = true, userTyped = false)
		{
		.num = false
		x = .getrec(:userTyped)
		.valid = x isnt false ? x[.nameField] : false
		if fillin
			.Fillin_fields(Record?(x) ? x : Record())
		}

	userTypedMultiMatch: false
	getrec(userTyped = false)
		{
		val = super.Get()
		if val is "" or .userTypedMultiMatch is val
			return false

		if userTyped
			.userTypedMultiMatch = false

		nameMatch = .NameMatchFieldAndValue(.nameField, val)
		nameMatchRec = .match_prefix(nameMatch.field, nameMatch.value)
		abbrevMatchRec = .getAbbrevMatch(val)

		// if different match records are found for name and abbrev fields,
		// we don't know which one to use, so return false (also makes field invalid)
		if .differentMatchRecord?(userTyped, nameMatchRec, abbrevMatchRec)
			{
			.userTypedMultiMatch = val
			return false
			}

		matchRec = nameMatchRec isnt false ? nameMatchRec : abbrevMatchRec
		if matchRec isnt false
			.set(matchRec)
		return matchRec
		}
	getAbbrevMatch(val)
		{
		if not .useAbbrevField?
			return false
		return .match_prefix(.abbrevField, val)
		}
	match_prefix(field, val)
		{
		// check for exact match on basequery
		if (false isnt (rec = .lookup_record(val, field)))
			return rec

		return MatchRecordFromPrefix(.GetQuery(), field, val)
		}
	differentMatchRecord?(userTyped, nameMatchRec, abbrevMatchRec)
		{
		return userTyped and nameMatchRec isnt false and abbrevMatchRec isnt false and
			nameMatchRec[.numField] isnt abbrevMatchRec[.numField]
		}

	AllowFillin?(rec)
		{
		// when allow other - if user modifies control and the old value was "other"
		// and new value is "other" - do not clear the fill fields
		return super.AllowFillin?(rec) and
			(.currentValueIsKey?(rec) or .previousValueIsKey?())
		}

	currentValueIsKey?(rec)
		{
		return not (.AllowOther?() and rec.Empty?())
		}

	previousValueIsKey?()
		{
		return Date?(.Send("GetField", .Parent.Name))
		}

	// used by KeyControl, for the "In List" multi-select option
	Key_DisplayField()
		{
		return .nameField
		}

	SetFieldRestrictions(restrictions)
		{
		super.SetFieldRestrictions(restrictions)
		.Setquery(.BuildQuery(.base_query))
		}
	}
