// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(initial_selects = #(), .name = '')
		{
		.Reset(initial_selects)
		}

	Name()
		{
		return .name
		}

	// expecting simple format of conditions
	Reset(initial_selects = #())
		{
		.select_vals = Record() // used in RecordControl in Select2Control
		.convertSimpleSelects(initial_selects)
		.initial_selects = .select_vals.DeepCopy()
		}

	PrependInitialSelect(initial_selects)
		{
		convertedSelects = Object()
		.convertSimpleSelects(initial_selects, convertedSelects)
		.initial_selects = convertedSelects.DeepCopy()
		.select_vals = convertedSelects.MergeUnion(.select_vals)
		}

	Select_vals()
		{
		// WARNING: returns actual object - other code may modify
		return .select_vals
		}

	selectChanged?: false
	UsingDefaultFilter?()
		{
		if .selectChanged? is true
			return false
		return .select_vals.Any?({ it.check is true }) and
			.select_vals.Filter({ it.check is true }) is
				.defaultSelects().Filter({ it.check is true })
		}

	SetSelectVals(.select_vals, sf) // used by SelectControl
		{
		if not .UsingDefaultFilter?()
			.selectChanged? = true
		.select_vals.RemoveIf({ not sf.Fields.Has?(it.condition_field) })
		}
	ops: #(
		"=":		"equals",
		"==":		"equals",
		"!=":		"not equal to",
		"<":		"less than",
		"<=":		"less than or equal to",
		">":		"greater than",
		">=":		"greater than or equal to",
		"=~":		"matches",
		"!~":		"does not match",
		"in":		"in list"
		)
	convertSimpleSelects(initial_selects, select_vals = false)
		{
		if select_vals is false
			select_vals = .select_vals
		for r in initial_selects
			{
			op = r.Member?(1) ? .ops[r[1]] : ''
			val = r.GetDefault(2, '')
			condition = .buildCondition(r[0], op, val)
			condition.check = r.Member?(1)
			select_vals.Add(condition)
			}
		}
	Add_initial_selects(initial_selects, select_vals = false) // used by reporter
		{
		i = 0
		if select_vals is false
			select_vals = .select_vals
		for r in initial_selects
			{
			select_vals['fieldlist' $ i] = SelectPrompt(r[0])
			select_vals['oplist' $ i] = .ops[r[1]]
			select_vals['val' $ i] = r[2]
			select_vals['checkbox' $ i] = true
			++i
			}
		}

	SaveSelects()
		{
		.Ensure()
		selects = .filterInitial(.select_vals, .defaultSelects())
		.retrySave()
			{|t|
			.delete_selects(t)
			.output_selects(t, selects)
			}
		}
	Ensure()
		{
		Database('ensure userselects (userselect_user, userselect_title,
			userselect_selects, userselect_TS) key (userselect_user, userselect_title)')
		}
	filterInitial(select_vals, initial_selects)
		{
		// setting check to true so that initial_selects will be found correctly
		// can get away with this ONLY BECAUSE we are already setting check to false for
		// all the saved selects
		return select_vals.Copy().Each({ it.check = true }).Difference(initial_selects).
			Each({ it.check = false})
		}
	// if user has used Save as Default, then filter out that from the same
	// instead of initial_selects
	defaultSelects()
		{
		return .defaultSel.Empty?()
			? .initial_selects
			: .defaultSel
		}
	// broken out so we can test
	retrySave(block)
		{
		RetryTransaction()
			{ |t|
			block(t)
			}
		}
	delete_selects(t)
		{
		t.QueryDo('delete ' $ .userselects_query())
		}
	output_selects(t, selects)
		{
		if selects.Empty?()
			return
		t.QueryOutput('userselects', Record(
			userselect_user: Suneido.User,
			userselect_title: .name,
			userselect_selects: selects))
		}
	defaultSel: #()
	LoadSelects(ctrl, noUserDefaultSelects? = false, _accessGoTo? = false)
		{
		if not TableExists?('userselects') or .name is '' or
			.name is VirtualListColModel.TmpSelectName
			return false
		sf = ctrl.GetSelectFields()
		noDefault? = noUserDefaultSelects? or accessGoTo?
		// in case old bad selects, users can still get into screens and we are notified
		try
			{
			.defaultSel = noDefault? ? Object() : .getSelects(savedDefault?:).Copy()
			hisSel = .getSelects()
			defaultSelUnChecked = .defaultSel.Map({ it.Copy().Merge(#(check: false)) })
			hisSel = hisSel.Map({ it.Copy().Merge(#(check: false)) })
			combinedSel = .defaultSel.Copy().Append(
				hisSel.Difference(defaultSelUnChecked))
			if .HasSavedDefault? // if there is default saved, remove initial selects
				.select_vals = Object()
			.select_vals.Append(combinedSel)
			// clean up invalid ones
			.select_vals.RemoveIf({ not sf.Fields.Has?(it.condition_field) })
			.removeDuplicateEmptyDefaults()
			.select_vals = .select_vals.UniqueValues()[.. .maxSelectRecords()]
			}
		catch (err)
			SuneidoLog('ERROR: (CAUGHT) Cannot load Selects: ' $ err,
				caughtMsg: 'bad select data, needs attention')
		}
	maxSelectRecords()
		{
		return SelectRepeatControl.MaxRecords
		}

	HasSavedDefault?: false
	getSelects(savedDefault? = false)
		{
		if false is defaultSel = .querySelect(:savedDefault?)
			return #()
		if savedDefault?
			.HasSavedDefault? = true
		return defaultSel.userselect_selects
		}
	querySelect(savedDefault? = false)
		{
		return Query1(.userselects_query(:savedDefault?))
		}
	removeDuplicateEmptyDefaults()
		{
		// TODO: treat default same as initial select
		for f in .initial_selects.Filter(.emptyOp?).Map({ it.condition_field })
			if .select_vals.HasIf?({ it.condition_field is f and not .emptyOp?(it) })
				.select_vals.RemoveIf({
					it.condition_field is f and it.check is false and .emptyOp?(it) })
		}
	emptyOp?(x)
		{
		return x[x.condition_field].operation is ''
		}

	emptyOps: #('less than or equal to', 'equals')
	notEmptyOps: #('greater than', 'not equal to')
	buildCondition(field, selOp, selVal)
		{
		condition = Record()
		condition.condition_field = field

		value = ''
		if selOp.Has?('empty')
			operation = selOp
		else if selOp is ''
			operation = ''
		else if selVal is ''
			operation = .emptyOps.Has?(selOp)
				? 'empty'
				: .notEmptyOps.Has?(selOp)
					? 'not empty'
					: ''
		else
			{
			value = Datadict(field).Encode(selVal)
			operation = selOp
			}
		condition[field] = Object(:operation, :value, value2: '')
		return condition
		}

	userselects_query(savedDefault? = false)
		{
		name = savedDefault? ? .name $ '~default' : .name
		return 'userselects where
			userselect_user is ' $ Display(Suneido.User) $
			' and userselect_title is ' $ Display(name)
		}

	DefaultSelectName(ctrl, title, query)
		{
		selectName = ListCustomize.BuildCustomKeyFromQueryTitle(query, title)
		if 0 isnt overrideSelect? = ctrl.Send('OverrideSelectManager?')
			{
			return overrideSelect?
				? ''
				: selectName
			}
		return selectName
		}
	}