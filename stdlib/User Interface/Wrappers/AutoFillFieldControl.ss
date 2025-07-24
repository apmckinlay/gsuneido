// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	origCandidates: false
	getCandidateFn: false
	New(@args)
		{
		super(@args)
		.origCandidates = args.GetDefault('candidates', false)
		.getCandidateFn = args.GetDefault('getCandidateFn', false)
		}

	getCandidates()
		{
		if .origCandidates isnt false
			return .origCandidates
		if .getCandidateFn isnt false
			return (.getCandidateFn)(.Get())
		return false
		}

	prevValue: ''
	EN_CHANGE(@args)
		{
		super.EN_CHANGE(@args)
		value = .Get()
		candidates = .getCandidates()

		if false is res = .findMatch(value, .prevValue, candidates)
			return 0

		.prevValue = value
		.Set(res)
		.SetSel(value.Size(), -1)
		.Send("NewValue", .Get())
		return 0
		}

	findMatch(value, prevValue, candidates)
		{
		if candidates is false or value is '' or value is prevValue
			return false

		candidatesLower = candidates.Map(#Lower)
		valueLower = value.Lower()
		preValueLower = prevValue.Lower()

		if valueLower is '' or
			preValueLower.Prefix?(valueLower) or
			false is (i = candidatesLower.FindIf({ it.Prefix?(valueLower) })) or
			valueLower is candidatesLower[i]
			return false

		return candidates[i]
		}
	}