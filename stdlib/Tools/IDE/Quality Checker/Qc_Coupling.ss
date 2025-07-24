// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	warningsThreshold: 20 //With > this many records coupled the rating decreases

	CallClass(recordData, minimizeOutput? = false)
		{
		coupledRecords = Qc_globalRefs(recordData.code, recordData.lib is 'stdlib')
		coupledSize = coupledRecords.Size()

		warnings = .generateWarnings(minimizeOutput?, coupledRecords, coupledSize)
		rating = .getRating(coupledSize)
		desc = .getDescription(minimizeOutput?, coupledSize)
		size = Max(0, coupledSize - .warningsThreshold)
		return Object(:warnings, :desc, :rating, :size)
		}

	generateWarnings(minimizeOutput?, coupledRecords, coupledSize)
		{
		if minimizeOutput? and coupledSize <= .warningsThreshold
			return Object()

		group = Object().Set_default(Object())
		for c in coupledRecords.Members().Sort!()
			{
			if false is lib = Qc_whichLib(c)
				lib = 'Unknown'
			group[lib].Add(Object(name: c, count: coupledRecords[c]))
			}
		warningText = "Depends on:\n"
		for lib in .applicationLibraries().MergeUnion(#(stdlib, Unknown))
			if group.Member?(lib)
				{
				warningText $= '\t' $ lib $ ': '
				for rec in group[lib]
					warningText $= rec.name $ '(' $ rec.count $ "), "
				warningText = warningText.RemoveSuffix(', ') $ '\n'
				}
		name = warningText.RemoveSuffix('\n')
		return Object([:name])
		}

	applicationLibraries()
		{
		return GetContributions('ApplicationLibraries')
		}

	getRating(numDependencies)
		{
		maxRating = 5
		return .warningsThreshold > numDependencies
			? maxRating
			: Max(0, maxRating + .warningsThreshold - numDependencies)
		}

	getDescription(minimizeOutput?, coupledSize)
		{
		if coupledSize > .warningsThreshold or not minimizeOutput?
			return 'This record depends on ' $ coupledSize $ ' records' $
				' - Attempt to limit to ' $ .warningsThreshold

		if not minimizeOutput?
			return 'This record does not depend on any other records'
		return ''
		}
	}