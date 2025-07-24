// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
// used by Select2Control and ReporterModel
// called like: Select2(sf).Where(data)
class
	{
	New(sf)
		{
		.sf = sf
		.ops_desc = .TranslateOps()
		}
	Ops: (
		('greater than',				'>',	pre: '',			suf: ''),
		('greater than or equal to',	'>=',	pre: '',			suf: ''),
		('less than',					'<',	pre: '',			suf: ''),
		('less than or equal to',		'<=',	pre: '',			suf: ''),
		('equals',						'is',	pre: '',			suf: ''),
		('not equal to', 				'isnt',	pre: '',			suf: ''),
		('empty',						'is',	pre: '',			suf: ''),
		('not empty',					'isnt',	pre: '',			suf: ''),
		('contains',					'=~',	pre: '(?i)(?q)',	suf: ''),
		('does not contain',			'!~',	pre: '(?i)(?q)',	suf: ''),
		('starts with',					'=~',	pre: '^(?i)(?q)',	suf: ''),
		('ends with',					'=~',	pre: '(?i)(?q)',	suf: '(?-q)$'),
		('matches',						'=~',	pre: '',			suf: '')
		('does not match',				'!~',	pre: '',			suf: '')
		)
	TranslateOps()
		{
		return .Ops.Map({ TranslateLanguage(it[0]) })
		}

	Where(data, fields = false, except = false, extra_dd = #(),
		exclude_menu_options = false)
		{
		where = errs = ""
		joinflds = Object()
		for (i = 0; i < Select2Control.Numrows; i++)
			{
			if exclude_menu_options and data['menu_option' $ i] is true
				continue
			if false is (op = .operation(data, i))
				continue
			result = .field(data, i, errs, fields, except)
			errs = result.errs
			if false is (fld = result.fld)
				continue
			result = .value(extra_dd, fld, op, data, i)
			if result.errs isnt ''
				{
				errs $= result.errs
				continue
				}
			fld = .Empty_field(fld, op)
			.doIfForeignNameAbbrevField(fld, data, i)
				{ |joinNums|
				op = #('', 'in', pre: '', suf: '')
				result.val = '(' $ joinNums.nums.Map(Display).Join(',') $ ')'
				fld = joinNums.numField
				}
			joinflds.Add(fld)
			if result.selectFunction isnt false
				fld = result.selectFunction $ "(" $ fld $ ")"
			where $= " where " $ fld $ " " $ op[1] $ " " $ result.val
			}
		return Object(:where, :errs, :joinflds)
		}
	operation(data, i)
		{
		op = data['oplist' $ i]
		if (data['checkbox' $ i] isnt true or
			op is "" or data['fieldlist' $ i] is "" or
			false is (pos = .ops_desc.Find(op)))
			return false
		return .Ops[pos]
		}
	field(data, i, errs, fields, except)
		{
		if false is (fld = .sf.PromptToField(data['fieldlist' $ i]))
			{
			errs $= "Can't Find Field: " $ data['fieldlist' $ i] $ "\n"
			return Object(fld: false, :errs)
			}
		if ((fields isnt false and not fields.Has?(fld)) or
			(except isnt false and except.Has?(fld)))
			return Object(fld: false, :errs)
		return Object(:fld, :errs)
		}
	value(extra_dd, fld, op, data, i)
		{
		dd = extra_dd.Member?(fld) ? extra_dd[fld] : Datadict(fld)
		dd = .formatSummarizeCalcField(fld, extra_dd, dd)
		if .Invalid_operator?(op, dd)
			{
			err= "Invalid operator for " $ data['fieldlist' $ i] $ '\n'
			return Object(errs: err, val: false)
			}

		val = Display(not .value_required?(op)
			? ""
			:  .string_op?(op)
				? op.pre $ data['val' $ i] $ op.suf
				: dd.Encode(data['val' $ i]))

		if .InvalidOpValue?(op, data['val' $ i])
			return Object(
				errs: "Invalid value " $ Display(data['val' $ i]) $
					" for operator " $ Display(op[0]),
				val: false)

		// handle selecting on fields where what is displayed to the user
		// is formatted differently than what is stored in the db
		selectFunction = false
		if dd.Member?('SelectFunction')
			selectFunction = dd.SelectFunction

		return Object(errs: '', :val, :selectFunction)
		}
	Invalid_operator?(op, dd)
		{
		if .string_op?(op) and not dd.Base?(Field_string)
			return true
		// can't use operators that require a value for images
		return dd.Base?(Field_image) and .value_required?(op)
		}
	formatSummarizeCalcField(fld, extra_dd, dd)
		{
		// summarize field (from Reporter) should have the same format as field
		if fld =~ '^total_calc\d|^max_calc\d|^min_calc\d|^average_calc\d'
			{
			calc_field = fld.Replace('^total_|^max_|^min_|^average_', '')
			if extra_dd.Member?(calc_field)
				dd = extra_dd[calc_field]
			}
		return dd
		}
	string_op?(op)
		{
		return op[1].Suffix?('~')
		}
	value_required?(op)
		{
		return not op[0].Suffix?('empty')
		}
	InvalidOpValue?(op, val)
		{
		if not .string_op?(op)
			return false

		return not Regex?(op.pre $ val $ op.suf)
		}
	Empty_field(fld, op, sf = false)
		// for empty/notempty on joined _name use _num
		{
		if sf is false
			sf = .sf

		if (.value_required?(op) or
			sf.OrigFields().Has?(fld))
			return fld
		suffix = sf.FieldSuffix(fld)
		if (fld.Suffix?('_name' $ suffix))
			numfld = fld.Replace('_name' $ suffix $ '$', '_num' $ suffix)
		else if (fld.Suffix?('_abbrev' $ suffix))
			numfld = fld.Replace('_abbrev' $ suffix $ '$', '_num' $ suffix)
		else
			return fld
		if (not sf.OrigFields().Has?(numfld))
			return fld
		return numfld
		}
	doIfForeignNameAbbrevField(fld,  data, i, block)
		{
		if false is joinNums = GetForeignNumsFromNameAbbrevFilter(fld, .sf,
			data['oplist' $ i], data['val' $ i])
			return
		block(joinNums)
		}
	}
