// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
/*
*	minimizeOutput? is used to determine formatting of output text for TODO window
*	or for quality check GUI. True is for TODO window.
*
*	lineWarnings is used by Addon_check_code to add grey warnings + squiggles to libview.
*/
Memoize
	{
	HashArgs?: true
	OkForResetAll?: false
	Func(lib, name, code, minimizeOutput?, extraChecks = false,
		missingTest = "check_local")
		{
		if .skipQc?(lib, name, code)
			return Object(#(warnings: #(), rating: false, desc: ""),
				lineWarnings: Object())

		recordData = Record(:lib, recordName: name, :code)
		qcArgs = [:recordData, :minimizeOutput?, :missingTest]
		type = extraChecks ? 'extra' : 'normal'
		return .RunChecksAsContributions(type, qcArgs)
		}

	skipQc?(lib, name, code)
		{
		return lib is 0 or name is 0 or lib is "" or name is "" or .EmptyRecord(code) or
			name.Suffix?(".js") or name.Suffix?(".css") or
			LibRecordType(code) not in ("function", "class")
		}

	EmptyRecord(code)
		{
		return not code.Has?('}') or code.AfterFirst('{').BeforeLast('}').LineCount() <= 2
		}

	RunChecksAsContributions(type, qcArgs, _checkStop? = false)
		{
		warnings = Object(lineWarnings: Object())
		for c in SoleContribution("ContinuousQualityChecks").Filter({ it.type is type })
			{
			if .skipThisContrib?(c, qcArgs.recordData.recordName)
				{
				warnings.Add(Object(warnings: Object(), desc: "", rating: 5))
				continue
				}
			if checkStop? isnt false and checkStop?()
				throw "Outdated"
			cl = Global(c.name)
			result = cl(@qcArgs)
			if result.Member?('lineWarnings')
				warnings.lineWarnings.Append(result.Extract('lineWarnings'))
			warnings.Add(result)
			}
		return warnings
		}

	skipThisContrib?(c, recordName)
		{
		if c.GetDefault('noUpdates', false) and recordName.Prefix?('Update_')
			return true

		if c.GetDefault('noTests', false) and
			LibraryTags.RemoveTagFromName(recordName).Suffix?('Test')
			return true

		return false
		}
	}
