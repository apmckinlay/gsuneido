// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// TODO check "layering"
// i.e. that libraries do not reference names in "later" libraries
class
	{
	CallClass(library)
		{
		return (new this).InternalRun(library)
		}
	InternalRun(library)
		{
		sups = LibrarySuppressions(library)
		results = ""
		// should be .GlobalName() but want to include e.g. Name?__protect
		QueryApply(library $ ' where name =~ `\A[[:upper:]]\w*[!?]?\w*\Z`', group: -1)
			{|x|
			if not CodeTags.Matches(x.lib_current_text) or sups.Has?(x.name)
				continue
			results $= .check1(library, x)
			if results.Size() > 4000 /*= max size */
				return results $ "... TOO MANY RESULTS ..."
			}
		return results
		}
	check1(library, x)
		{
		if .skip?(x)
			return ''
		CheckCode(x.lib_current_text, x.name, library, warnings = [])
		results = warnings.
			Filter({ not it.msg.Prefix?("WARNING") }).
			Map!({ library $ ':' $ x.name $ ' - ' $ it.msg $ '\n' }).
			Join()
		return results
		}
	skip?(x)
		{
		if LibRecordType(x.lib_current_text) in (#dll, #struct, #callback)
			return true
		return .old_update?(x) or .BuiltDate_skip?(x.lib_current_text)
		}
	old_days: 10
	old_update?(x)
		{
		return (false isnt date = x.name.Extract('^Update_(20\d\d\d\d\d\d)')) and
			Date(date) < Date().Minus(days: .old_days)
		}
	BuiltDate_skip?(text)
		{
		if text.Has?('// BuiltDate ') // fast first check
			{
			if false isnt bd = text.Extract("^// BuiltDate > (\d+)")
				return BuiltDate() < Date(bd)
			if false isnt bd = text.Extract("^// BuiltDate < (\d+)")
				return BuiltDate() > Date(bd)
			}
		return false
		}
	getter_stdnames()
		{
		stdnames = Object()
		for name in QueryList('stdlib where group = -1', #name)
			stdnames[name] = true
		for name in BuiltinNames()
			stdnames[name] = true
		return .stdnames = stdnames // once only
		}
	}
