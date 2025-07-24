// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// Based on ideas from the book Building Maintainable Software
class
	{
	CallClass(@args)
		{
		Qc_ContinuousChecks(@args)
		}

	SlowChecks(lib, name, code)
		{
		if lib is 0 or name is 0 or lib is "" or name is ""
			return ""
		recordData = Record(:lib, recordName: name, :code)
		return Qc_ContinuousChecks.RunChecksAsContributions('slow', [:recordData])
		}

	CheckWithExtra(lib, name, text, minimizeOutput? = false, missingTest = "check_local")
		{
		warnings = Qc_ContinuousChecks(lib, name, text, minimizeOutput?, :missingTest)
		extraWarnings = Qc_ContinuousChecks(lib, name, text, minimizeOutput?,
			extraChecks:, :missingTest)
		warnings.lineWarnings.Append(extraWarnings.Extract('lineWarnings'))
		warnings.Append(extraWarnings)
		return warnings
		}

	CalcRatings(warnings)
		{
		count = 0
		sum = 0
		maxRating = 10 //Each whole number represents half a star
		warnings.Each({
			if it.Member?("rating")
				{
				if it.rating is false
					return false
				sum += it.rating
				count++
				}
			if it.Member?("maxRating")
				maxRating = Min(it.maxRating * 2, maxRating)
			})
		return Min(maxRating, (sum * 2 / count).RoundDown(0))
		}

	// Used to compare ratings when sending or receiving
	CompareRatings(lib, name, revision, old, missingTestOld = "check_local",
		missingTestRevision = "check_revision")
		{
		qcErrs = Object().Set_default(Object())
		// Use missingTestRevision when svc holds the change,
		// Record is type incompatible with quality checker
		if false is revisionRating = .getCodeQualityRating(
			lib, name, revision, missingTestRevision)
			return ''
		revision = revisionRating / 2
		if old is false
			{
			if revisionRating < .minNewQuality
				qcErrs.newRecordErrors.Add(
					Object(:lib, :name, :revision, min: .minNewQuality / 2))
			}
		else
			{
			oldRating = .getCodeQualityRating(lib, name, old, :missingTestOld)
			//Record used to be type incompatible with quality checker
			if oldRating is false and revisionRating < .minNewQuality
				qcErrs.existingRecordErrors.Add(
					Object(:lib, :name, :revision, old: .minNewQuality / 2))
			else if revisionRating < oldRating
				qcErrs.existingRecordErrors.Add(
					Object(:lib, :name, :revision, old: oldRating / 2))
			}
		return qcErrs.Empty?() ? '' : .outputCompareResult(qcErrs) $ '\n'
		}

	getCodeQualityRating(lib, change, code, missingTest = "check_local")
		{
		allWarnings = .CheckWithExtra(lib, change, code, minimizeOutput?:, :missingTest)
		return .CalcRatings(allWarnings)
		}

	minNewQuality: 8
	outputCompareResult(qcErrors)
		{
		qualityErrors = ""
		for existRec in qcErrors.existingRecordErrors
			qualityErrors $= ' ' $ existRec.lib $ ':' $ existRec.name $ " rating is: " $
				String(existRec.revision) $ ",  maintain/exceed: " $
				String(existRec.old) $ '\n'
		for newRec in qcErrors.newRecordErrors
			qualityErrors $= '+' $ newRec.lib $ ':' $ newRec.name $ ' rating is: ' $
				String(newRec.revision) $ ',  min: ' $ String(newRec.min) $ '\n'
		return qualityErrors.Trim()
		}

	RunQCWith(libs, afterFn)
		{
		qcErrors = .findQcErrors(libs)
		if qcErrors isnt ""
			{
			qcErrors $= "\n\n" $ "Continue running tests?"
			if YesNo(qcErrors,
				"Quality Checker Failed - NO effect on sending to SVC currently")
				afterFn()
			return
			}
		afterFn()
		}

	RunQC(libs)
		{
		qcErrors = .findQcErrors(libs)
		if qcErrors isnt ""
			Alert(msg: qcErrors,
				title: "Quality Checker Failed - NO effect on sending to SVC")
		else
			Alert(msg: "All records have acceptable quality ratings",
				title: "Quality Checker Passed")
		}

	findQcErrors(libs)
		{
		if false isnt results = .getCachedWarning('svc_all_changes')
			return results
		if false is settings = SvcSettings()
			return "Svc settings are not set. Quality Checker unable to run"
		results = ''
		Working("Checking Code Quality...")
			{
			results = .checkAllLibs(libs, settings)
			}
		return results
		}
	getCachedWarning(lib)
		{
		if not Suneido.Member?('SvcCommit_Warnings')
			return false

		return SvcRunChecks.GetPreCheckResults(lib)
		}
	checkAllLibs(libs, settings)
		{
		results = ""
		try
			svc = Svc(server: settings.svc_server, local?: settings.svc_local?)
		catch(e)
			return 'Cannot connect to version control for quality checker: ' $ e
		for lib in libs
			{
			if false isnt cached = .getCachedWarning(lib)
				{
				results $= cached
				continue
				}

			QueryApply(lib $ ' where group is -1 and lib_modified isnt ""')
				{
				svcGet = svc.Get(lib, it.name)
				results $= Qc_Main.CompareRatings(:lib, name: it.name,
					revision: it.lib_current_text,
					old: svcGet is false ? false : svcGet.text,
					missingTestOld: svc.MissingTest?(lib, it.name))
				}
			}
		return results
		}
	}
