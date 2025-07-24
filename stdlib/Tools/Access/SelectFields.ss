// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(fields = #(), excludeFields = #(), joins = true, .headerSelectPrompt = false,
		includeInternal = false, includeMasterNum = false, convertCustomNumFields = false)
		{
		.orig_fields = fields
		fields = fields.Copy()
		.exclude_fields(fields, excludeFields, includeInternal)

		.field_ob = Object()
		.table_fieldprefix = Object()
		.table_fieldnum = Object()
		.field_suffixes = Object()
		.other = Object()
		.converted = Object()

		.add_abbrev_and_name_for_num_fields(fields, joins, includeMasterNum,
			convertCustomNumFields)

		.handlePrompts(fields, excludeFields, includeInternal)
		}

	exclude_fields(fields, excludeFields, includeInternal)
		{
		fields.RemoveIf({ excludeFields.Has?(it) })

		fields.RemoveIf(
			{
			dd = Datadict(it, #(ExcludeSelect))
			not includeInternal and	dd.Member?(#ExcludeSelect) and
				dd.ExcludeSelect is true
			})
		}

	add_abbrev_and_name_for_num_fields(fields, joins, includeMasterNum,
		convertCustomNumFields)
		{
		if joins isnt true or .headerSelectPrompt is 'no_prompts'
			return

		// add _abbrev, _name for num fields
		for fld in fields.Copy() // need copy because of removes
			{
			if false isnt .HandleSpecialJoinFields(fld, fields, .other,
				.field_suffixes, .table_fieldprefix, .table_fieldnum, includeMasterNum,
				:convertCustomNumFields, converted: .converted)
				continue

			if not fld.Suffix?("_num") and not fld.Has?("_num_")
				continue

			// get num field suffix if applicable
			suffix = Opt('_', fld.AfterLast('_num_'))
			// strip off everything after and including '_num'
			f = fld.BeforeFirst('_num')
			if .nameAndAbbrevExist(fields, f, fld)
				continue
			.addForeignTableFields(f, fld, fields, suffix, includeMasterNum)
			}
		}
	GetConverted()
		{
		return .converted
		}

	// f is the prefix (fld.BeforeFirst('_num'))
	// fld is a num field (either ending with _num or with _num_ followed by a suffix
	nameAndAbbrevExist(fields, f, fld)
		{
		if fields.Has?(f $ "_name") and fields.Has?(f $ "_abbrev")
			{
			fields.Remove(fld)
			return true
			}

		if fld.Has?('_num_') and fields.Has?(fld.Replace('_num_', '_abbrev_')) and
			fields.Has?(fld.Replace('_num_', '_name_'))
			{
			fields.Remove(fld)
			return true
			}

		return false
		}

	addForeignTableFields(f, fld, fields, suffix, includeMasterNum)
		{
		if false is tablename = .foreign_key_table(f)
			return
		prompt = .GetFieldPrompt(fld, .field_ob.Members())
		fields.Add(Object(f $ "_name" $ suffix, prompt $ " Name"))
		fields.Add(Object(f $ "_abbrev" $ suffix, prompt $ " Abbrev"))
		.other[f $ "_name" $ suffix] = tablename
		.other[f $ "_abbrev" $ suffix] = tablename
		.field_suffixes[f $ "_name" $ suffix] = suffix
		.field_suffixes[f $ "_abbrev" $ suffix] = suffix
		if false is includeMasterNum
			fields.Remove(fld)
		.table_fieldprefix[tablename] = f
		.table_fieldnum[tablename $ "," $ suffix] = fld
		}

	GetJoinNumField(field)
		{
		if not .NameAbbrev?(field)
			return false
		if false is tablename = .other.GetDefault(field, false)
			return false
		suffix = .field_suffixes[field]
		return .table_fieldnum[tablename $ "," $ suffix]
		}

	handlePrompts(fields, excludeFields, includeInternal)
		{
		// handle prompts, if duplicate, heading or field name is used
		for f in fields
			{
			if Object?(f) // num fields
				{
				prompt = f[1].Trim()
				f = f[0]
				}
			else if .headerSelectPrompt is 'no_prompts'
				prompt = f
			else
				prompt = .GetFieldPrompt(f, .field_ob.Members()).Trim()

			// check excludeFields again in case it included a joined field
			if excludeFields.Has?(f)
				continue

			.warnIfNoPrompt(prompt, f, includeInternal)

			.AddField(f, prompt)
			}
		}

	// The purpose of GetFieldPrompt is to handle duplicate custom field prompts.
	GetFieldPrompt(field, promptList = #())
		{
		if field.Prefix?('custom_')
			{
			prompt = Prompt(field)
			if not promptList.Has?(prompt)
				return prompt
			}
		return SelectPrompt(field)
		}

	warnIfNoPrompt(prompt, f, includeInternal)
		{
		if ((prompt is '' or prompt is f) and
			.headerSelectPrompt isnt 'no_prompts' and includeInternal isnt true)
			ProgrammerError('no prompt for: ' $ f)
		}

	logError(s)
		{
		if not TestRunner.RunningTests?() and not s.Has?('custom_') and
			not s.Has?('_reference') and not s.Has?('calc')
			ProgrammerError(s)
		}

	cacheSize: 500
	foreign_key_table(field)
		{
		cache = Suneido.GetInit(#ForeignKeyTables,
			{ LruCache(.foreign_key_table_getter, .cacheSize) })
		return cache.Get(field)
		}

	foreign_key_table_getter(field)
		{
		QueryApply("indexes", key: true, columns: field $ "_num")
			{ |x|
			cols = TableModelGetColumns(x.table)
			if cols isnt false
				{
				if cols.Has?(field $ "_name") and cols.Has?(field $ "_abbrev")
					return x.table
				}
			else
				if not QueryEmpty?("columns", table: x.table, column: field $ "_name") and
					not QueryEmpty?("columns", table: x.table, column: field $ "_abbrev")
					return x.table
			}
		return false
		}

	HandleSpecialJoinFields(@args)
		{
		for fn in GetContributions('SelectFieldsHandleSpecialJoinFields')
			if fn(@args) is true
				return true
		return false
		}

	AddField(field, prompt)
		{
		if Customizable.DeletedField?(field)
			return

		// reporter handles num/name/abbrev, but if the datasource already has _name
		// _name can't be excluded from datasource or it won't be in the list at all
		if .field_ob.Member?(prompt) and field isnt .field_ob[prompt]
			{
			.logError('duplicate select prompt: ' $ field $ ' & ' $
				.field_ob[prompt])

			if field.Prefix?('custom_') and .field_ob[prompt].Prefix?('custom_')
				{
				// get the original custopm field mapped to this prompt
				origField = .field_ob[prompt]
				// assign "new" custom field to this prompt
				.field_ob[prompt] = field
				// now the "prompt" is already mapped to another field, GetFieldPrompt
				// might return the SelectPrompt for the origField
				newPrompt = .GetFieldPrompt(origField, .field_ob.Members())
				// if newPrompt is same as prompt, it will set the map back to what it
				// was before, otherwise it will have entries for both fields
				.field_ob[newPrompt] = origField
				}
			else if .field_ob.Member?(heading = Heading(field).Trim())
				.field_ob[field] = field
			else
				.field_ob[heading] = field
			}
		else
			.field_ob[prompt] = field
		}

	Prompts()
		{
		return .field_ob.Members()
		}

	PromptsToFields(prompts)
		{
		if String?(prompts)
			prompts = prompts.Split(',')
		return prompts.Map(.PromptToField).Instantiate()
		}

	PromptToField(prompt)
		{
		prompt = prompt.Trim()
		return .HasPrompt?(prompt)
			? .field_ob[prompt]
			: false
		}

	HasPrompt?(prompt)
		{
		return .field_ob.Member?(prompt)
		}

	FieldSuffix(fld)
		{
		return .field_suffixes.GetDefault(fld, "")
		}

	FieldToPrompt(field)
		{
		return .field_ob.Find(field)
		}

	FieldsToPrompts(fieldList)
		{
		if String?(fieldList)
			fieldList = fieldList.Split(',')
		return fieldList.Map(.FieldToPrompt).Instantiate()
		}

	OrigFields()
		{
		return .orig_fields
		}

	Getter_Fields()
		{
		return .field_ob // object with prompt: field
		}

	AllJoins()
		{
		return .Joins(.other.Members())
		}

	Joins(fields)
		{
		if String?(fields)
			fields = fields.Split(',')

		joins = .buildJoinsOb(fields)
		return joins.Join('')
		}
	// add joins so select can be done using name and abbrev fields from
	// master tables
	// keep track of tables joined to prevent joining to same master table twice
	// for the same field suffix (in case they used name and abbrev)
	buildJoinsOb(fields, withDetails? = false)
		{
		join = .getJoinItems(fields)
		joins = Object()
		join_tables = Object()
		for item in join.Members()
			{
			str = ''
			field = join[item]
			table = item[.. item.Find(",")]
			suffix = .field_suffixes[field]
			if join_tables.Has?(table $ suffix)
				continue
			join_tables.Add(table $ suffix)
			prefix = .table_fieldprefix[table]
			numField = .table_fieldnum[table $ "," $ suffix]
			str $= " leftjoin by(" $ numField $ ") (" $
				table $ " project " $ prefix $ "_num, " $
				prefix $ "_name, " $ prefix $ "_abbrev"
			if prefix $ "_num" isnt numField
				str $= " rename " $ prefix $ "_num to " $ numField
			if suffix isnt ""
				{
				str $= " rename " $ prefix $ "_name to " $ prefix $ "_name" $ suffix
				str $= " rename " $ prefix $ "_abbrev to " $ prefix $ "_abbrev" $ suffix
				}
			str $= ") "
			if withDetails?
				{
				fields = Object()
				if suffix is ''
					fields.Add(prefix $ "_name").Add(prefix $ "_abbrev")
				else
					fields.Add(prefix $ "_name" $ suffix).Add(prefix $ "_abbrev" $ suffix)
				joins.Add(Object(:str, :fields))
				}
			else
				joins.Add(str)
			}
		return joins
		}
	JoinsOb(fields, withDetails? = false)
		{
		if String?(fields)
			fields = fields.Split(',')
		return .buildJoinsOb(fields, :withDetails?)
		}

	getJoinItems(fields)
		{
		join = Object()
		for fld in fields
			{
			fld = fld.Trim()
			if .field_ob.Member?(fld)
				fld = .field_ob[fld] // convert to field name
			if .other.Member?(fld)
				join[.other[fld] $ "," $ fld] = fld
			}
		return join
		}

	FormulaFields(src)
		{
		fields = Object()
		.ScanFormula(src, { |f| fields.Add(f) })
		return fields
		}

	FormulaPromptsToFields(src, fields)
		{
		dst = ""
		.ScanFormula(src,
			{ |f| dst $= f; fields.AddUnique(f) },
			{ |s| dst $= s })
		return dst
		}

	ScanFields(src, reserved = false)
		{
		ob = Object()
		pos = 0
		while pos < src.Size()
			{
			whitespace = src[pos..].Extract('^(\s*?)\S')
			pos += whitespace.Size()
			if false is first = .first(src[pos..])
				return ob
			if first is '('
				{
				++pos
				continue
				}
			matched = false
			for match in Object(.best_match, .text_match, .number_match, .reserved_match)
				{
				if '' isnt text = match(first, src[pos..], :reserved)
					{
					ob.Add(Object(:pos, end: pos + text.Size()))
					pos += text.Size()
					matched = true
					break
					}
				}
			if not matched
				pos += first.Size()
			}
		return ob
		}

	reserved_match(unused, src, reserved)
		{
		if reserved is false
			return ''

		for item in reserved
			if src.Prefix?(item)
				return item

		return ''
		}

	text_match(first, src)
		{
		if first is '`'
			{
			pos = src.Find('`', 1)
			return pos is src.Size() ? '' : src[..pos+1]
			}
		else if first in ('"', "'")
			{
			pos = 1
			while pos < src.Size()
				{
				pos = Min(src.Find(`\`, pos), src.Find(first, pos))
				if src[pos] is first
					return src[..pos+1]
				if pos is src.Size()
					return ''
				pos += 2
				}
			}
		return ''
		}

	number_match(first /*unused*/, src)
		{
		if false is numberText = src.Extract('^(\.[0-9]+|[0-9]+\.[0-9]+|[0-9]+)')
			return ''
		return src[..numberText.Size()]
		}

	ScanFormula(src, field_block, other_block = function (unused) { })
		{
		src = src.Trim()
		while src isnt ''
			{
			whitespace = src.Extract('^(\s*?)\S')
			other_block(whitespace)
			src = src[whitespace.Size() ..]
			next = .first(src)
			if '' isnt best = .best_match(next, src)
				field_block(.PromptToField(best))
			else if '' isnt best = .text_match(next, src)
				other_block(best)
			else
				other_block(next)

			if best isnt ''
				next = best
			src = src[next.Size() ..]
			}
		}

	best_match(next, src)
		{
		if src[next.Size()] is '('
			return ""
		possible = .lookup[next]
		best = ''
		for p in possible
			if src.Prefix?(p) and p.Size() > best.Size()
				best = p
		return best
		}

	getter_lookup()
		{
		tbl = Object().Set_default(#())
		for p in .Prompts()
			tbl[.first(p)].Add(p)
		return .lookup = tbl // once only
		}

	notWordChars: '^_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'
	first(s) // string should not start with space charcters
		{
		start = s.Find1of(.notWordChars)
		if start is 0 // when the first character is operator
			start = 1
		return s[..start]
		}

	NameAbbrev?(field)
		{
		return .field_suffixes.Member?(field)
		}
	}
