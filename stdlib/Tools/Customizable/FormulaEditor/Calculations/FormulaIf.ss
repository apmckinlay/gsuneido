// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(@args)
		{
		.checkArgSize(args)
		for (i = 0; i < args.Size() - 1; i += 2)
			{
			condition = .getArgValue(args[i])
			.checkConditionType(Object(condition.type))
			if condition.value is true
				return .getArgValue(args[i + 1])
			}
		return .getArgValue(args.Last())
		}

	getArgValue(arg)
		{
		// Suggestion 24575 changed the args of FormulaIf to block
		// but old formulas will still pass the arg value
		if Function?(arg)
			return arg()
		return arg
		}

	Validate(@args)
		{
		.checkArgSize(args)

		for (i = 0; i < args.Size() - 1; i += 2)
			.checkConditionType(args[i])

		rtnTypes = args.Last().Copy()
		for (i = 1; i < args.Size(); i += 2)
			rtnTypes.MergeUnion(args[i])
		return rtnTypes
		}

	checkArgSize(args)
		{
		if args.Size() < 3/*=minimun args*/
			throw "Formula: IF must have at least 3 arguments"

		if not args.Size().Odd?()
			throw "Formula: IF must have odd number of arguments"
		}

	checkConditionType(conditionType)
		{
		if conditionType isnt Object(FORMULATYPE.BOOLEAN)
			throw "Formula: IF condition must be a <Boolean>"
		}
	}
