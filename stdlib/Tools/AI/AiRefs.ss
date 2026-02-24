function()
	{
	refs = Object().Set_default(0)
	defs = Object()
	QueryApply('stdlib where group = -1 sort name')
		{|x|
		defs[x.name] = true
		scnr = Scanner(x.text)
		for token in scnr
			{
			if token.GlobalName?()
				++refs[token]
			}
		}
	results = Object()
	for name, count in refs
		if defs.Member?(name)
			results.Add([name, count])
	return results.Sort!(By(1))
	}