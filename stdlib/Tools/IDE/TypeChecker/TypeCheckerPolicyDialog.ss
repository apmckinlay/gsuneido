// TypeCheckerPolicyDialog
// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Type Checker Policy'
	levels: #('off', 'warn', 'error')
	ops: #('<', '<=', '>', '>=')

	CallClass(policy)
		{
		return OkCancel(Object(this, policy), .Title)
		}

	New(policy)
		{
		super(.layout(policy))
		data = Record().Merge(policy)
		if policy.Member?(#confidence)
			{
			// edited as operator + percent, recomposed into `confidence` in OK
			conf = .parsePred(policy.confidence)
			data.confidenceOp = conf.op
			data.confidencePct = conf.pct
			data.Delete(#confidence)
			}
		.Data.Set(data)
		}

	layout(policy)
		{
		vert = Object('Vert')
		for name in policy.Members().Sort!()
			{
			if name is #confidence
				continue   // rendered separately, not as off/warn/error
			vert.Add(Object('Pair',
				Object('Static', name),
				Object('ChooseList', .levels, width: 6, mandatory:, :name)))
			}
		if policy.Member?(#confidence)
			vert.Add(Object('Pair',
				Object('Static', 'confidence'),
				Object('Horz',
					Object('ChooseList', .ops, width: 4, mandatory:,
						name: 'confidenceOp'),
					Object('Skip', small:),
					Object('Number', name: 'confidencePct', width: 5, mask: '###'),
					Object('Static', '%'))))
		return Object('Record', vert)
		}

	maxPercent: 100
	// "=0.70 predicate string" -> Object(op:, pct:).  e.g. ">=0.70" -> (">=", 70)
	parsePred(pred)
		{
		op = '>='
		for cand in #('<=', '>=', '<', '>')   // 2-char first so ">=" beats ">"
			if pred.Prefix?(cand)
				{
				op = cand
				break
				}
		return Object(:op, pct: (Number(pred[op.Size() ..]) * .maxPercent).Round(0))
		}

	OK()
		{
		data = .Data.Get().Copy()
		if data.Member?(#confidenceOp)
			{
			pct = data.GetDefault(#confidencePct, '')
			if not (Number?(pct) and pct >= 0 and pct <= .maxPercent)
				{
				AlertError('Confidence % must be a number between 0 and 100')
				return false   // OkCancel keeps the dialog open when OK() is false
				}
			// Go accepts .7 / 0.7 / 0.70 alike, so the bare fraction is fine
			data.confidence = data.confidenceOp $ (pct / .maxPercent)
			data.Delete(#confidenceOp)
			data.Delete(#confidencePct)
			}
		return data
		}
	}
