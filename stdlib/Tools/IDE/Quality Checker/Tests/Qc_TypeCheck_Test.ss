// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cannotRun?: false
	Setup()
		{
		try
			.cannotRun? = not TypeCheckHelper.BinaryExists?()
		catch
			.cannotRun? = true
		}

	Test_skipsNonCompilableCode()
		{
		rd = .recordData(code: 'this is not @#$ valid')
		Assert(Qc_TypeCheck(rd).nError is: -1)
		}

	Test_skipsMissingRecord()
		{
		rd = .recordData(recordName: 'NoSuchRecord_zzz')
		Assert(Qc_TypeCheck(rd).nError is: -1)
		}

	Test_checksFunctions()
		{
		if .cannotRun?
			Assert(true)
		else
			{
			code = `function(s)
						{
						try
							s.Compile()
						catch
							return false
						return true
						}`
			rd = .recordData(:code, recordName: 'Compilable?')
			Assert(Qc_TypeCheck(rd).nError is: 0)
			}
		}

	Test_liveResultIsConsistent()
		{
		if .cannotRun?
			Assert(true)
		else
			{
			result = Qc_TypeCheck(.recordData())
			Assert(result.nError is: result.warnings.Size())
			Assert(result.desc is: (result.warnings.Empty?() ? '' : 'Type Check'))
			}
		}

	Test_returnsWellFormedResult()
		{
		result = Qc_TypeCheck(.recordData())
		Assert(result.Member?(#warnings))
		Assert(result.Member?(#desc))
		Assert(result.Member?(#nError))
		}

	recordData(code = 'class { }', lib = 'stdlib', recordName = 'TypeCheckerControl')
		{
		return [:code, :lib, :recordName]
		}
	}
