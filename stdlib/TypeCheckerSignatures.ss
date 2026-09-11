// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// BuiltDate > 20260824
function(force = false)
	{
	if force is false and true is Suneido.GetDefault(#TypeCheckerSignaturesLoaded, false)
		return true
	try
		TypeChecker.RegisterSignatures(TypeSignatures())
	catch (e)
		{
		if e.Prefix?("method not found") or e.Prefix?("can't find")
			return false
		throw e
		}
	Suneido.TypeCheckerSignaturesLoaded = true
	return true
	}
