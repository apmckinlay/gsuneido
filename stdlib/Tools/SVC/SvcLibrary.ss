// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
/*
SvcLibrary is used to get information from and update local library tables. SvcBook and
SvcLibrary handle the differences between libraries and books.

NOTE: This should never be called directly. Use SvcTable as a interface for this class
*/
class
	{
	Type: 'lib'
	library: false
	New(table)
		{
		SvcLibraryMonitor()
		.library = table
		}

	SvcEnabled?()
		{
		return not SvcDisabledLibraries().Has?(.library)
		}

	Query(deleted = false)
		{
		return .library $ ' where group is ' $ Display(deleted ? -2 : -1 )
		}

	Where()
		{
		return ' where group in (-1, -2)'
		}

	Table()
		{
		return .library
		}

	NameQuery(name, deleted = false)
		{
		return .Query(:deleted) $ ' where name is ' $ Display(name)
		}

	GetData(rec, t = false)
		{
		rec.path = .GetPath(rec, t)
		}

	MakeName(rec)
		{
		return rec.name
		}

	GetPath(rec, t = false)
		{
		path = ''
		DoWithTran(t)
			{|t|
			while rec.parent isnt 0 and rec.parent isnt ''
				{
				if false is rec = t.Query1(.library, num: rec.parent)
					break
				path = rec.name $ '/' $ path
				}
			}
		return path[.. -1]
		}

	rootPlaceholder: '<root>'
	FormatData(rec, t, deleted = false)
		{
		rec.group = deleted ? -2 : -1
		if rec.path is .rootPlaceholder
			rec.path = ''
		if rec.parent is '' or (rec.Member?(#path) and rec.path is '')
			rec.parent = 0
		if rec.path isnt ''
			{
			.ensurePath(rec, t)
			rec.path = .GetPath(rec, t)
			}
		}

	ensurePath(rec, t)
		{
		rec.parent = 0
		for name in rec.path.Split("/")
			{
			if name is .library
				continue
			folder = t.Query1(.library $ " where parent = " $ Display(rec.parent) $
				" and name = " $ Display(name) $ ' and group isnt -1')
			if folder is false
				folder = .make_folder(rec.parent, name, t)
			rec.parent = folder.num
			}
		}

	make_folder(parent, name, t)
		{
		rec = Record(:parent, group: parent, num: NextTableNum(.Table(), t), :name)
		t.QueryOutput(.Table(), rec)
		return rec
		}

	FormatText(text)
		{
		return text.Entab()
		}

	Restore(rec)
		{
		rec.lib_invalid_text = ''
		}

	Rename(rec, oldName /*unused*/, t)
		{
		if rec.path is ''
			if '' is rec.path = .GetPath(rec, t)
				rec.path = .rootPlaceholder
		}

	StageDelete(rec, t /*unused*/)
		{
		if rec.lib_before_path is ''
			rec.lib_before_path = .rootPlaceholder
		}

	Deleted?(rec)
		{
		return rec.group is -2
		}

	Move(rec, newParent, t)
		{
		t.QueryApply1(.library, num: rec.num)
			{
			if it.lib_before_path is ''
				it.lib_before_path = rec.lib_before_path
			else if it.lib_before_path is .rootPlaceholder
				it.lib_before_path = ''
			it.parent = rec.parent = newParent
			it.path = .GetPath(it, t)
			if it.path is it.lib_before_path
				it.lib_before_path = ''
			else if it.lib_before_path is ''
				it.lib_before_path = .rootPlaceholder
			svcTable = SvcTable(.library)
			svcTable.Update(it, :t)
			svcTable.Publish(#TreeChange)
			}
		}

	MoveLibrary(rec, newLib, newParent, t)
		{
		DoWithTran(t, update:)
			{|t|
			SvcTable(.library).StageDelete(rec.name, :t)
			rec.Delete(#num, #lib_before_text, #lib_before_path, #lib_committed)
			svcTable = SvcTable(newLib)
			rec.parent = newParent
			TreeModel.EnsureUnique(rec, newLib)
			rec.path = svcTable.GetPath(rec, t)
			t.QueryApply1(svcTable.NameQuery(rec.name, deleted:))
				{|deletedRec|
				rec.lib_committed = deletedRec.lib_committed
				rec.lib_before_text = deletedRec.lib_before_text
				rec.lib_before_path = deletedRec.lib_before_path
				if rec.path is '' and rec.lib_before_path is .rootPlaceholder
					rec.lib_before_path = ''
				deletedRec.Delete()
				}
			svcTable.VerifyModified(rec)
			svcTable.Output(rec, :t, committed: rec.lib_modified is '')
			}
		}

	FormatMaxCommitRecord(rec)
		{
		rec.group = -3
		rec.parent = 0
		}

	MaxCommitQuery()
		{
		return .library $
			' where name is "' $ SvcTable.MaxCommitName $ '" and group is -3'
		}

	Compare(local, master, discrepancies)
		{
		if Svc.Hash(local.lib_current_text) isnt Svc.Hash(master.text)
			discrepancies.Add(#text)
		}

	CheckRecord(rec, type)
		{
		if .skipRecord?(rec, type)
			return #()

		if true is .verifyCode(rec.text, rec.name, .library, results = [])
			return #()

		return results.Map!({ it.msg })
		}

	skipRecord?(rec, type)
		{
		return type is '-' or .library is 'Contrib' or
			CheckLibrary.BuiltDate_skip?(rec.text)
		}

	verifyCode(code, name, lib, results = false)
		{
		if true isnt msg = .checkCode(code, name, lib, results)
			return msg
		return .checkName(name, results, lib)
		}

	checkCode(code, name, lib, results = false) // split out for testing
		{
		return CheckCode(code, name, lib, results)
		}

	checkName(name, results, lib) //, type, text)
		{
		// check does not work if record name has ?__protect, .css, .js, skipping for now
		if name.Has?('?__protect') or CheckCode.IsWeb?(name)
			return true
		if '' isnt msg = .CheckOnServer(name, lib)
			{
			results.Add(Object(:msg))
			return false
			}
		return true
		}

	CheckOnServer(name, lib)
		{
		if not Libraries().Has?(lib)
			return ""
		if Sys.Client?()
			return ServerEval('SvcLibrary.CheckOnServer', name, lib)
		if false isnt (rec = Query1(lib, :name, group: -1)) and
			not CodeTags.Matches(rec.text)
			return ""
		try
			Global(LibraryTags.RemoveTagFromName(name))
		catch (err)
			{
			if err.Prefix?("can't find")
				return Libraries().Has?(lib)
					? "ERROR: can't find: " $ name : ""
			return err
			}
		return ""
		}

	ResetSvcDisabledCache()
		{
		SvcDisabledLibraries.ResetCache()
		}

	// Whenever a library record is restored to its original state, clear out
	// the test runner success table to ensure the tests are run again.
	// Otherwise, due to the restored record having lib_modified: ""
	// the reversion would not be seen via SvcCommitChecker.need_to_run_tests?
	PostRestore(restored? = false)
		{
		if restored?
			TestRunner.RequireRun()
		}

	// Whenever a new (not committed) library record is deleted, clear out
	// the test runner success table to ensure the tests are run again.
	// As non-committed records are immediately deleted (not staged), tests need to
	// be forced to run to ensure no tests / records depended on the deleted record
	PostStageDelete(name, skipPublish = false, committed? = false)
		{
		if skipPublish
			LibUnload(name)
		if not committed?
			TestRunner.RequireRun()
		}
	}
