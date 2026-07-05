function ()
	{
	op  = .confidenceOp
	pct = .confidencePct
	if op is "" or pct is ""
		return ""
	return op $ (Number(pct) / 100).Format("0.00")
	}