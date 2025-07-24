// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(@args)
		{
		data = args.data
		args.Delete('data')
		no_encode = args.GetDefault('no_encode', false)
		args.Delete('no_encode')
		build_callable = args.GetDefault('build_callable', false)
		args.Delete('build_callable')
		fields = args
		restrictions = .buildRestrictions(fields, data, no_encode, build_callable)
		return BuildQueryWhere(restrictions, :build_callable)
		}

	buildRestrictions(fields, data, no_encode, build_callable)
		{
		restrictions = Object() // get fields values from paramsdata
		for field in fields
			{
			result = .isTimestamp?(field)
			timestamp? = result.timestamp?
			field = result.field
			if String?(field) and field.Suffix?('Filters')
				for ob in data[field]
					{
					paramData = .getParamData(ob, ob.condition_field)
					.buildFieldRestriction(paramData, ob.condition_field,
						restrictions, timestamp?, no_encode, build_callable)
					}
			else
				{
				param = .getParamField(data, field)
				if param.Has?('__protect')
					continue

				.buildFieldRestriction(.getParamData(data, param),
					field, restrictions, timestamp?, no_encode, build_callable)
				}
			}
		return restrictions
		}

	getParamData(data, paramField)
		{
		if not data.Member?(paramField)
			return ''
		return data[paramField]
		}

	buildFieldRestriction(data, field, restrictions, timestamp?, no_encode,
		build_callable)
		{
		if not .checkParamData(data, field)
			return

		ob = data.Copy()
		.handleIdAllowOther(ob)
		ob.operation = ob.operation.Trim()
		if .invalidOperation?(ob)
			return

		.ConvertDateCodes(ob)
		.handleEmptyValues(field, ob)
		f = Datadict(field, #(SelectFunction))
		fieldName = f.Member?('SelectFunction')
			? f.SelectFunction $ "(" $ field $ ")"
			: field
		if .handleInList(ob, restrictions, fieldName)
			return

		if timestamp?
			{
			.add_timestamp_restriction(field, ob, restrictions)
			return
			}

		.handleBoolean(field, ob)
		val = no_encode is true ? ob.value : DatadictEncode(field, ob.value)
		val2 = ob.GetDefault('value2', '')
		if no_encode isnt true
			val2 = DatadictEncode(field, val2)

		.addRestrictions(ob, restrictions, fieldName, val, val2, build_callable)
		}

	checkParamData(data, field)
		{
		if Object?(data)
			return true

		if data isnt ""
			SuneidoLog('ERROR: invalid param data', calls:, params: Object(:field, :data))
		return false
		}

	ConvertDateCodes(param)
		{
		if param.GetDefault('dateConversionInfo', false) isnt false
			{
			showTime = param.dateConversionInfo.showTime
			for valMem in #(value, value2)
				if param.Member?(valMem) and param[valMem] isnt ''
					param[valMem] = DateControl.ConvertToDate(
						param[valMem], true, showTime)
			}
		}

	handleInList(ob, restrictions, fieldName)
		{
		// check if "in list" since this is handled differently
		if ob.operation is "in list" or ob.operation is "not in list"
			{
			// idtexts is used for Id field with allowOther option
			list = ob.value.MergeUnion(ob.value.GetDefault("idtexts", #()))
			restrictions.Add(Object(fieldName, ob.operation, list))
			return true
			}
		return false
		}

	invalidOperation?(ob)
		{
		if ob.operation is ""
			return true

		if not .ops.Member?(ob.operation)
			{
			// should be throwing error, but first log to see if it happens
			SuneidoLog("ERROR: GetParamsWhere.CallClass: " $
				"invalid operation: " $ ob.operation)
			return true// ignores it
			}
		return false
		}

	addRestrictions(ob, restrictions, field, val, val2, build_callable)
		{
		if ob.operation.Has?("range")
			{
			.handleRanges(ob, field, val, val2, restrictions)
			return
			}
		if not String?(val) and .ops[ob.operation].string_op
			{
			if Number?(val)
				val = String(val)
			else
				{
				// should be throwing error, but first log to see if it happens
				SuneidoLog("ERROR: GetParamsWhere.addRestrictions: " $
					"string operation on non-string", params: Object(:field, :val, :ob))
				return // ignores it
				}
			}
		if .ops[ob.operation].pre > ""
			val = .ops[ob.operation].pre $ val
		if .ops[ob.operation].suf > ""
			val $= .ops[ob.operation].suf
		if .ops[ob.operation].val isnt false
			val = .ops[ob.operation].val
		op = .ops[ob.operation].op
		if build_callable is true and ob.operation is 'not empty'
			op = 'isnt'
		restrictions.Add(Object(field, op, val))
		}

	handleRanges(ob, field, val, val2, restrictions)
		{
		if ob.operation is "range"
			{
			restrictions.Add(Object(field, ">=", val))
			restrictions.Add(Object(field, "<=", val2))
			}
		else if ob.operation is "not in range"
			restrictions.Add(Object(field, "not in range", val, val2))
		}

	handleIdAllowOther(ob)
		{
		if ob.operation is '' or not ob.Member?('idtext')
			return
		Assert(ob.operation is 'equals' or ob.operation is 'not equal to')
		ob.operation = ob.operation is 'not equal to'
			? 'not in list' : 'in list'
		ob.value = Object(ob.value, ob.idtext)
		// to handle Print from Preview
		ob.Delete('idtext')
		}

	isTimestamp?(field)
		{
		if Object?(field) and field[1] is 'timestamp'
			return Object(timestamp?: true, field: field[0])
		else
			return Object(timestamp?: false, :field)
		}

	handleBoolean(field, ob)
		{
		if DatadictType(field) isnt 'boolean'
			return
		if ob.operation in ('equals', 'not equal to') and ob.value in (false, "")
			{
			opposite = ob.operation is 'equals' ? 'not equal to' : 'equals'
			ob.operation = opposite
			ob.value = true
			}
		}

	handleEmptyValues(field, ob)
		{
		if false is vals = .getSelectEmptyVals(field)
			return

		operatorTranslations = #(
			'equals': 'in list'
			'not equal to': 'not in list'
			'empty': 'in list'
			'not empty': 'not in list'
			)
		if ((ob.operation in ('equals', 'not equal to') and vals.Has?(ob.value)) or
			ob.operation in ('empty', 'not empty'))
			{
			ob.operation = operatorTranslations[ob.operation]
			ob.value = vals
			}
		// handle inlist, not in list:
		//  -> if ANY of the empty values are in the list, add them all
		if ob.operation in ('in list', 'not in list') and
			not ob.value.Intersect(vals).Empty?()
			ob.value = ob.value.Copy().MergeUnion(vals)
		}

	getSelectEmptyVals(field)
		{
		ddOb = Datadict(field, #(Select_EmptyValues))
		if not ddOb.Member?('Select_EmptyValues')
			return false

		return ddOb.Select_EmptyValues
		}

	getParamField(data, field)
		{
		matches = Object()
		.addParamNameMatches(data, field, matches)
		.addPrefixMatches(data, field, matches)
		if matches.Size() > 1
			SuneidoLog("ERROR: (CAUGHT) GetParamsWhere.getParamField " $ field $
				" multiple matches, returning first of " $ Display(matches),
				caughtMsg: 'something needs to be excluded or renamed')
		return matches.Empty?() ? field : matches.First()
		}

	addParamNameMatches(data, field, matches)
		{
		for suffix in #('_param', '_params')
			if data.Member?(field $ suffix)
				matches.AddUnique(field $ suffix)
		}

	addPrefixMatches(data, field, matches)
		{
		if not matches.Empty?() or data.Member?(field) // already has match
			return

		for m in data.Members()
			if m.Prefix?(field) and not m.Suffix?('__protect')
				matches.AddUnique(m)
		}

	ops: (
		'greater than':
			(op: ">", pre: "", suf: "", val: false, string_op: false),
		'less than':
			(op: "<", pre: "", suf: "", val: false, string_op: false),
		'equals':
			(op: "is", pre: "", suf: "", val: false, string_op: false),
		'less than or equal to':
			(op: "<=", pre: "", suf: "", val: false, string_op: false),
		'greater than or equal to':
			(op: ">=", pre: "", suf: "", val: false, string_op: false),
		'not equal to':
			(op: "isnt", pre: "", suf: "", val: false, string_op: false),
		'empty':
			(op: "is", pre: "", suf: "", val: "", string_op: false),
		'not empty':
			(op: "isnt", pre: "", suf: "", val: "", string_op: false),
		'contains':
			(op: "=~", pre: "(?i)(?q)", suf: "", val: false, string_op: true),
		'does not contain':
			(op: "!~", pre: "(?i)(?q)", suf: "", val: false, string_op: true),
		'starts with':
			(op: "=~", pre: "^(?i)(?q)", suf: "", val: false, string_op: true),
		'ends with':
			(op: "=~", pre: "(?i)(?q)", suf: "(?-q)$", val: false, string_op: true),
		'matches':
			(op: "=~", pre: "", suf: "", val: false, string_op: true)
		'does not match':
			(op: "!~", pre: "", suf: "", val: false, string_op: true)
		'range':
			(op: "", pre: "", suf: "", val: false, string_op: false)
		'not in range':
			(op: "", pre: "", suf: "", val: false, string_op: false)
		'in list':
			(op: "", pre: "", suf: "", val: false, string_op: false)
		'not in list':
			(op: "", pre: "", suf: "", val: false, string_op: false)
		)

	add_timestamp_restriction(f, ob, restrictions)
		{
		val = DatadictEncode(f, ob.value)
		val2 = DatadictEncode(f, ob.GetDefault('value2', ''))
		.checktimestampRestrictionValid(ob, val, val2)
		endOfDay = startOfDay = endOfRange = false
		if Date?(val)
			{
			startOfDay = val.StartOfDay()
			endOfDay = val.EndOfDay()
			}
		if Date?(val2)
			endOfRange = val2.EndOfDay()
		restrictions.Add(@.getTSRestrictions(ob, f, startOfDay, endOfDay, endOfRange))
		}

	checktimestampRestrictionValid(ob, val, val2)
		{
		Assert(not (.ops.Member?(ob.operation) and .ops[ob.operation].string_op),
			"string operation on timestamp")
		Assert(not ((not Date?(val) and
			ob.operation isnt 'empty' and ob.operation isnt 'not empty') or
			(not Date?(val2) and
			(ob.operation is "range" or ob.operation is "not in range"))),
			"invalid operation on non-date")
		}

	getTSRestrictions(ob, f, startOfDay, endOfDay, endOfRange)
		{
		restrictions = Object()
		switch (ob.operation)
			{
		case 'greater than', 'less than or equal to':
			restrictions.Add(Object(f, .ops[ob.operation].op, endOfDay))
		case 'less than', 'greater than or equal to':
			restrictions.Add(Object(f, .ops[ob.operation].op, startOfDay))
		case 'equals':
			restrictions.Add(Object(f, ">=", startOfDay))
			restrictions.Add(Object(f, "<=", endOfDay))
		case 'range':
			restrictions.Add(Object(f, ">=", startOfDay))
			restrictions.Add(Object(f, "<=", endOfRange))
		case 'not in range':
			restrictions.Add(Object(f, "not in range", startOfDay, endOfRange))
		case 'not equal to':
			r = " and (" $ f $ " < " $ Display(startOfDay) $
				" or " $ f $ " > " $ Display(endOfDay) $ ")"
			restrictions.Add(Object(r, built: true))
		case 'empty', 'not empty':
			restrictions.Add(Object(f, .ops[ob.operation].op, ""))
			}
		return restrictions
		}
	}
