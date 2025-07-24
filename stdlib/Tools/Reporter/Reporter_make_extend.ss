// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// adds to formula_fields and where objects
	CallClass(data, sf, foreach_calc, formula_fields, where,
		summarize_by, summarize_fields, checkOnly = false)
		{
		uf = UserFuncs('Reporter', :checkOnly)
		extends = Object().Set_default("")

		formula_deps = Object()
		foreach_calc()
			{|key /*unused*/, prompt, formula|
			field = sf.PromptToField(prompt)
			sf.FormulaPromptsToFields(formula, formula_deps[field] = Object())
			}
		which = .which_extend(formula_deps,
			sf.PromptsToFields(data.reporter_summarizeby_cols),
			summarize_by, summarize_fields)

		foreach_calc()
			{ |i, prompt, formula, row|
			field = sf.PromptToField(prompt)
			if not field.Prefix?('calc')
				throw "Reporter: Formula field name " $ prompt $
					" is in use. Please rename the formula."
			newFn = FormulaEditor.ConstructFormula(sf, prompt, row.type, formula, field,
				skipReturnTypeCheck?:, :checkOnly)
			if newFn.formulaCode is false or newFn.formulaCode is "" or
				(Object?(newFn.formulaCode) and newFn.formulaCode.Member?('err'))
				{
				formulaCode = sf.FormulaPromptsToFields(formula, fields = Object())
				fn = .buildFormulaFunction(fields, formulaCode, prompt)
				}
			else
				{
				fn = newFn.formulaCode
				fields = newFn.fields.Split(',')
				}
			if not .compilable?(fn)
				throw "Reporter: Invalid Formula for " $ prompt
			name = uf.NeedFunc(fn)
			extends[which[field]] $=
				'extend calc' $ i $ ' = ' $ name $ '(' $ fields.Join(',') $ ')\n'
			if which[field] is 1
				where.Add(field)
			for f in fields
				formula_fields.AddUnique(f)
			}
		return extends
		}
	buildFormulaFunction(fields,  fn,  prompt)
		{
		fieldlist = fields.Join(',')
		// using Pack on the return value, so error can be caught by try block
		return 'function (' $ fieldlist $ ')
	{
	try {
		func = function (' $ fieldlist $ ')
			{
			' $ fn $ '
			}
		rtnVal = func(' $ fieldlist $ ')
		Pack(rtnVal)
		return rtnVal
		}
	catch (err)
		return "ERROR: ' $ prompt $ ': " $ err
	}'
		}
	compilable?(fn)
		{
		if not Compilable?(fn)
			return false
		return true
		}

	// NOTE: which_extend only deals with field names, NOT prompts

	// formula_deps are the dependencies for (fields used by) each formula_deps
	//		e.g. #(myformula: (field_one, field_two)
	// before_fields are all the fields available before the summarize
	// summarize_fields are the fields only available after the summarize
	// 		NOT including the summarize by fields (e.g. count, total_amount)
	which_extend(formula_deps, before_fields, summarize_by, summarize_fields)
		{
		which = Object()
		only_before = before_fields.Difference(summarize_by)
		only_after = summarize_fields
		sumfunc_fields = summarize_fields.Copy().Remove('count').
			Map!({ it.AfterFirst('_') }) // strip 'total_' etc.
		which = Object().Set_default(1) // default to after
		for field in formula_deps.Members()
			{
			deps = formula_deps[field]
			must_be_before = deps.Intersects?(only_before) or
				summarize_by.Has?(field) or sumfunc_fields.Has?(field)
			must_be_after = deps.Intersects?(only_after)
			if must_be_before and must_be_after
				.conflict()
			if must_be_before or must_be_after
				.mark(field, must_be_before ? 0 : 1, which, formula_deps, summarize_by)
			}
		return which
		}
	mark(field, w, which, formula_deps, summarize_by)
		{
		if summarize_by.Has?(field)
			{
			which[field] = 0
			return
			}
		if which.Member?(field)
			if which[field] is w
				return // already done
			else
				.conflict()
		if not formula_deps.Member?(field)
			return
		which[field] = w
		for f in formula_deps.GetDefault(field, #())
			.mark(f, w, which, formula_deps, summarize_by) // recursive
		for f in formula_deps.Members()
			{
			if f is field
				continue
			if formula_deps[f].Has?(field)
				.mark(f, w, which, formula_deps, summarize_by) // recursive
			}
		}
	conflict()
		{
		throw "Reporter: Formulas cannot depend on fields " $
			"from both before and after the Summarize"
		}
	}