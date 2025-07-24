// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(recordData)
		{
		dupFinder = new DupCodeFinder()
		allDuplicates = dupFinder.FindDuplicates(recordData.lib)
		warnings = Object()

		for hash in allDuplicates.Members()
			{
			//Matches the class name exactly
			positionFound = allDuplicates[hash].Find(":" $ recordData.recordName $ ":")
			if positionFound isnt allDuplicates[hash].Size() //if className is a duplicate
				warnings.Add(Record(name: "Duplicate: " $ allDuplicates[hash]))
			}

		desc = warnings.Empty?()
			? "No duplicates found of code in this class "
			: "Duplicate code from in this class was found "
		desc $= "-> This check does not affect the rating of code"
		return Object(:desc, :warnings)
		}
	}
