// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(recordData)
		{
		valid, function? = .checkable?(recordData)
		if not valid
			return .emptyResult()

		libName = recordData.lib
		recName = recordData.recordName
		qc_warnings = .collectWarnings(libName, recName, function?, recordData.code)
		desc = qc_warnings.NotEmpty?() ? 'Type Check' : ''
		// TC_ERROR implies type checker errors not a type checking error
		nError = qc_warnings.Filter({ it.name.Prefix?('TC_ERROR') }).Size() > 0
			? -1 : qc_warnings.Size()
		return Object(warnings: qc_warnings, :desc, :nError)
		}

	collectWarnings(libName, recName, function?, code)
		{
		skipLineageOrLibName = false
		// if lib not loaded, pass libName so lineage is skipped,
		// we also skip lineage if its a function instead of a class
		if not Libraries().Has?(libName) or function?
			skipLineageOrLibName = libName

		qc_warnings = Object()
		try
			{
			result = TypeCheckHelper.Run(recName, TypeCheckerMethods.Infer,
				:skipLineageOrLibName, src: code)
			errors, unused = TypeCheckHelper.FormatDiagnostics(
				result.diagnostics, library: libName)
			// only keep warnings for this record
			prfx = libName $ ":" $ recName
			for error in errors.Filter({ it.Prefix?(prfx) })
				qc_warnings.Add([name: error])
			}
		catch (e)
			qc_warnings.Add([name: 'TC_ERROR: ' $ String(e)])

		return qc_warnings
		}

	// we cannot check if
	checkable?(recordData)
		{
		// not compilable
		if not Compilable?(recordData.code)
			return false, false

		// no binary
		try
			{
			if not TypeCheckHelper.BinaryExists?()
				return false, false
			}
		catch // sometimes this seems to throw on
			return false, false

		c = recordData.code.Compile()
		func? = Function?(c)
		return Class?(c) or func?, func?
		}

	emptyResult()
		{
		return Object(warnings: #(), desc: '', nError: -1, rating: false)
		}
	}