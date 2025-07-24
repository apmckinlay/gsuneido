// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
//TODO: Do something with all nightly task warnings
class
	{
	Qcm_UnusedPublicMethods(recordData)
		{
		warnings = Object()
		methods = .validMethods(recordData.code)
		numPubMethods = methods.Size()
		refsToClass = FindReferences(recordData.recordName)
		.pruneReferencedMethods(methods, refsToClass)
		numPubMethodsUnused = methods.Size()

		if numPubMethods > 0 and numPubMethodsUnused is numPubMethods
			Print("Record unused - No public methods used")
		Print("The following methods are not utilized: " $ Display(methods))
		}

	pruneReferencedMethods(methods, refsToClass)
		{
		for ref in refsToClass.list
			{
			if ref.Member?("Found")
				{
				for method in methods.Copy()
					{
					if ref.Found.Has?('.' $ method)
						methods.Remove(method)
					}
				}
			}
		}

	validMethods(code)
		{
		methods = ClassHelp.Methods(code)
		methods.RemoveIf({ it.Capitalized?() is false or it.Prefix?('On_') or
			it.Prefix?("Getter_") })
		methods.Remove('New', 'CallClass', 'Destroy','Controls', 'Layout')
		return methods
		}
	}






