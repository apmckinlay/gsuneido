// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(restrictions, build_callable = false)
		{
		where = ''
		for r in restrictions
			{
			if build_callable
				r[0] = .convertField(r[0])
			if r.GetDefault('built', false) is true // restriction is pre-built
				where $= " " $ r[0]
			else if r[1] in ('in list', 'not in list')
				where $= .build_in_list(r)
			else if r[1] is 'not in range'
				where $= .notInRange(r)
			else // normal case with field, operator and value
				where $= ' and ' $ r[0] $ ' ' $ r[1] $ ' ' $ Display(r[2])
			}
		where = where.Replace('and', not build_callable ? 'where' : '', 1)
		return not build_callable
			? where
			: .callable(where)
		}

	convertField(field)
		{
		idx = field.Find('(')
		return idx is field.Size()
			? 'rec.' $ field
			: field[..idx+1] $ 'rec.' $ field[idx+1..] // for fields with SelectFunction
		}
	notInRange(r)
		{
		return ' and (' $ r[0] $ ' < ' $ Display(r[2]) $ ' or ' $
			r[0] $ ' > ' $ Display(r[3 /* = second range val */]) $ ')'
		}

	callable(where)
		{
		try
			return .compile(where)
		catch (err)
			throw .processError(err)
		}

	compile(where)
		{
		return ('function (rec)
			{
			return ' $ (where.Blank?() ? 'true' : where) $ '
			}').Compile()
		}

	processError(err)
		{
		if not err.Has?('compile error')
			return err
		caughtMsg = .client?()
			? 'user notified of possible invalid filter'
			: 'unattended'
		SuneidoLog('ERROR: (CAUGHT) ' $ err, calls:, :caughtMsg)
		return 'SHOW: There was a problem with the filter.\n' $
			'This could be caused by invalid filter options or too many In List values'
		}

	client?()
		{
		return Sys.Client?()
		}

	build_in_list(r)
		{
		if not Object?(r[2]) or r[2].Size() is 0
			return ""

		field = r[0].RemovePrefix('rec.')
		list = r[2].Map({ Display(DatadictEncode(field, it)) })
		operator =  r[1].RemoveSuffix(' list')
		return " and " $ r[0] $ " " $ operator $ " (" $ list.Join(', ') $ ')'
		}
	}
