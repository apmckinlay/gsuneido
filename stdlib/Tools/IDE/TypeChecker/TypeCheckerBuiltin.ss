// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// BuiltDate > 20991231
// LEGACY-TYPECHECK-BINARY
class
	{
	Check(method, orderedSrc, refs, policy, restartOnError? = true)
		{
		TypeCheckerSignatures()
		if method is TypeCheckerMethods.Infer
			return TypeChecker.Infer(orderedSrc, refs, policy)
		if method is TypeCheckerMethods.Annotate
			return TypeChecker.Annotate(orderedSrc, refs, policy)
		throw "TypeChecker: unknown method: " $ Display(method)
		}

	Available?()
		{
		return true
		}

	Start()
		{
		return false
		}

	Stop() { }

	Path()
		{
		return false
		}

	SetPath(path/*unused*/) { }
	}
