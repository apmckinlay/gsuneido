// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
/*
SvcCommitChecker contains all pre-commit checks for when changes are about to be sent. It
is called from SvcControl.
*/
class
	{
	// Send mandatory checks ===================================================
	MandatoryChecks(changes, model, table)
		{
		for check in [.table_selected?,
			.have_all_master_changes?,
			.Local_changes_up_to_date?,
			.have_checkmarked_local_changes?,
			.errors_in_local_changes?,
			.whitespace_in_name?,
			.checkLineEnds,
			.need_to_run_tests?,
			.additional_svc_checks,
			]
			if '' isnt msg = check(:changes, :model, :table)
				return msg
		return ''
		}

	table_selected?(table)
		{
		return table is false
			? "Please select a library or book"
			: ''
		}

	have_all_master_changes?(changes, model, table)
		{
		if table is "svc_all_changes"
			return false is model.NoMasterChangesInLibs?(changes)
				? "Please get the master changes from selected libraries before sending"
				: ''
		else
			return model.MasterChanges.NotEmpty?()
				? "Please get the master changes before sending"
				: ''
		}

	Local_changes_up_to_date?(model, table)
		{
		return model.LocalChangeNeedsUpdate?(table, log:)
			? 'Local Changes are not up to date, please refresh'
			:  ''
		}

	have_checkmarked_local_changes?(changes)
		{
		return changes.Empty?()
			? 'Please checkmark the local changes to send'
			: ''
		}

	errors_in_local_changes?(changes, table)
		{
		if table is 'Contrib'
			return ''
		errors = Object()
		changes.Each({ .collectErrors(it, errors) })
		return Opt('Unable to send selected changes:\r\n\t- ', errors.Join('\r\n\t- '))
		}

	maxErrLength: 100
	collectErrors(change, errors)
		{
		if change.type is '-'
			return

		results = .checkRecord(change)
		if results.errors.Empty?()
			return

		error = results.errors.Join('\r\n').Ellipsis(.maxErrLength, atEnd:)
		errors.Add(change.lib $ ':' $ change.name $ ' (' $ error $ ')')
		}

	checkRecord(change)
		{
		if Record?(CodeState.InvalidRec(change.lib, change.name))
			return #(errors: ('invalid'), warnings: ())

		results = Object(errors: Object(), warnings: Object())
		for line in SvcTable(change.lib).Check(change.name, change.type)
			if line.Prefix?('WARNING: ')
				results.warnings.Add(line.RemovePrefix('WARNING: '))
			else
				results.errors.Add(line)
		return results
		}

	whitespace_in_name?(model, changes, table)
		{
		if not model.Library?(table)
			return ''

		recordsWithSpace = Object()
		for change in changes
			if change.type is '+' and change.name =~ '\s'
				recordsWithSpace.Add(change.name)

		return recordsWithSpace.Empty?() ? '' :
			'Record "' $ recordsWithSpace.Join('", "') $
				'" contains whitespace in the name'
		}

	maxAllowedRecs: 10
	checkLineEnds(model, changes, table)
		{
		if not .libraryCheck?(model, table)
			return ''
		invalidRecords = Object()
		for change in changes
			{
			if invalidRecords.Size() is .maxAllowedRecs
				{
				invalidRecords.Add('Too many record to display')
				break
				}
			if change.type isnt '-'
				{
				name = change.name
				lib = change.lib
				if CheckCode.HasInvalidLineEnd?(:lib, :name)
					invalidRecords.Add(lib $ ":" $ name)
				}
			}
		return invalidRecords.Empty?()
			? ''
			: 'Following record(s) uses non-standard line ending characters:\n\t' $
				invalidRecords.Join('\n\t')
		}
	// extracted for Test override
	libraryCheck?(model, table)
		{
		return model.Library?(table)
		}

	need_to_run_tests?(model, changes, table)
		{
		if not model.Library?(table)
			return ''

		forceTestsFn = OptContribution('SvcForceTests?', function (@unused) {return true})
		if forceTestsFn(changes.Map({ it.name })) is false
			return ''

		lastLocalChange = [what: 'N/A', when: Date.Begin(), how: 'SVC']
		lastTestRun = TestRunner.LastSuccess()
		return .last_library_change(lastLocalChange) > lastTestRun
			? .libraryChangeMessage(lastLocalChange, lastTestRun)
			: ''
		}

	last_library_change(lastLocalChange)
		{
		for lib in LibraryTables().Difference(SvcTable.ExcludedTables)
			if false isnt x = QueryLast(lib $ ' sort lib_modified')
				.checkLastGet(x, lib, lastLocalChange)
		return lastLocalChange.when
		}

	checkLastGet(x, lib, llc)
		{
		if Date?(x.lib_modified) and x.lib_modified > llc.when
			{
			llc.how = x.group is -2 ? 'deleted' : 'modified'
			llc.what = lib $ ': ' $ x.name
			llc.when = x.lib_modified
			}
		}

	libraryChangeMessage(lastLocalChange, lastTestRun)
		{
		how = lastLocalChange.how
		details = how isnt 'SVC'
			? '\n\nLast record ' $ how $ ':\n\t- ' $ lastLocalChange.what $
				'\n\t- Date: ' $ lastLocalChange.when.ShortDateTime()
			: ''
		testRun = lastTestRun isnt Date.Begin()
			? '\n\nLast successful test: ' $ lastTestRun.ShortDateTime()
			: ''
		return 'You must run all the tests successfully before sending changes' $
			testRun $ details
		}

	additional_svc_checks(model)
		{
		return model.CheckSvcStatus()
		}

	ClearPreCheck()
		{
		Suneido.SvcCommit_Warnings = Object().Set_default(Object().Set_default(Object()))
		Suneido.SvcPreCheck_ForceStop = Object()
		}

	// Send warnings - return message (or "") ==================================
	PreCommitChecks(local_list, id, change, masterRec, index)
		{
		if .stopCurrent(id, index, local_list)
			return
		try
			row = local_list.GetRow(index)
		catch (unused, 'member not found')
			return

		library = id.BeforeFirst('__#')
		cacheMember = library $ '_' $ change.lib $ '_' $ change.name
		if Suneido.SvcCommit_Warnings.errMap[cacheMember] isnt #()
			{
			row.svc_warning = Suneido.SvcCommit_Warnings.errMap[cacheMember]
			.repaintList(local_list)
			return
			}
		if false is local = SvcTable(change.lib).Get(change.name)
			{
			.handleDeleted(change, library, row, local_list, cacheMember)
			return
			}
		large = .checkLargeRecord(change, local, library)
		if .stopCurrent(id, index, local_list)
			return
		quality = .checkCodeQuality(change, local, masterRec, library)
		if .stopCurrent(id, index, local_list)
			return

		results = .checkRecord(change)
		.setStatus(row, cacheMember, .getStatus(large, quality, results))
		if '' isnt msg = .formatMsg(large, quality, results)
			Suneido.SvcCommit_Warnings.msgMap[change.lib $ '_' $ change.name] = msg
		.repaintList(local_list)
		}

	checkLargeRecord(change, local, library)
		{
		if large = change.type isnt '-' and local.text.Size() > 100.Kb() /*= max text */
			Suneido.SvcCommit_Warnings[library].largeCheck.
				AddUnique(Object(lib: change.lib, name: change.name))
		return large
		}

	checkCodeQuality(change, local, masterRec, library)
		{
		// could take a while to check code quality
		if '' isnt quality = .verifyCodeQuality(change, local, masterRec)
			{
			qualityCheckOb = Object(lib: change.lib, name: change.name, results: quality)
			Suneido.SvcCommit_Warnings[library].qualityCheck.AddUnique(qualityCheckOb)
			}
		return quality
		}

	getStatus(large, quality, checkRecord)
		{
		return checkRecord.errors.NotEmpty?()
			? .error
			: large or quality isnt '' or checkRecord.warnings.NotEmpty?()
				? .warning
				: .pass
		}

	stopCurrent(id, index, local_list)
		{
		if .stopRunning?(id, index, local_list)
			{
			Suneido.SvcPreCheck_ForceStop.Delete(id)
			return true
			}
		else
			return false
		}

	handleDeleted(change, library, row, local_list, cacheMember)
		{
		if change.type is '-'
			{
			.checkDeleted(change, library, row, cacheMember)
			.repaintList(local_list)
			}
		}

	stopRunning?(id, index, local_list)
		{
		return not local_list.Member?('Hwnd') or
			.forceStop?(id) or (index > local_list.GetNumRows() - 1)
		}

	checkDeleted(change, library, row, cacheMember)
		{
		if .delete_and_referenced?(change.name, change.type, change.lib) is false
			.setStatus(row, cacheMember, .pass)
		else
			{
			Suneido.SvcCommit_Warnings[library].referencedDeleted.AddUnique(change)
			Suneido.SvcCommit_Warnings.msgMap[change.lib $ '_' $ change.name] =
				'Deleted but referenced!'
			.setStatus(row, cacheMember, .warning)
			}
		}

	repaintList(local_list)
		{
		if local_list.Member?('Hwnd') and
			local_list.CompareAndSet(#Repainter, false, true) // temp set to true
			local_list.Repainter = RunOnGui({
				if local_list.Member?('Hwnd')
					local_list.Repaint()
				local_list.Repainter = false
				})
		}

	forceStop?(id)
		{
		return Suneido.GetDefault('SvcPreCheck_ForceStop', #()).GetDefault(id, false)
		}

	getter_warning()
		{
		return CLR.WarnColor
		}

	getter_pass()
		{
		return CLR.GREEN
		}

	getter_error()
		{
		return CLR.ErrorColor
		}

	setStatus(row, cacheMember, status)
		{
		row.svc_warning = status
		Suneido.SvcCommit_Warnings.errMap[cacheMember] = status
		}

	formatMsg(large, quality, checkRecord)
		{
		msgOb = Object()
		if large
			msgOb.Add('Record is over size limit (100 kb)')
		msgOb.Add(quality)

		join = '\r\n- '
		if checkRecord.errors.NotEmpty?()
			msgOb.Add('Record has syntax error(s)' $
				(checkRecord.errors[0] isnt 'invalid'
					? ':' $ join $ checkRecord.errors.Join(join)
					: ''))
		msgOb.Add(Opt('Record has the following warning(s):' $ join,
			checkRecord.warnings.Join(join)))
		return msgOb.Remove('').Join('\r\n')
		}

	delete_and_referenced?(name, type, lib)
		{
		if type isnt '-' or lib is 'Contrib'
			return false
		name = LibraryTags.RemoveTagFromName(name)
		if not name.GlobalName?()
			return false
		if not Uninit?(name)
			return false
		if FindReferences.DefinitionExists?(UnusedStandardLibraries(), name)
			return false
		return FindReferences.ReferenceExists?(name, #(Contrib))
		}

	//Qc additions
	verifyCodeQuality(change, local, masterRec)
		{
		return Qc_Main.CompareRatings(lib: change.lib, name: change.name,
			revision: local.text, missingTestOld: masterRec.missingTestOld,
			old: masterRec.master is false ? false : masterRec.master.text)
		}

	extra_check(changes)
		{
		msgs = Object()
		for contrib in Contributions('Svc_ExtraChecks')
			msgs.Add(contrib(:changes))
		return msgs.Join('\n')
		}

	ProcessWarnings(library, changes = false)
		{
		if not Suneido.Member?('SvcCommit_Warnings')
			return ''

		largeCheck = Suneido.SvcCommit_Warnings[library].largeCheck
		qualityCheck = Suneido.SvcCommit_Warnings[library].qualityCheck
		referencedDeleted = Suneido.SvcCommit_Warnings[library].referencedDeleted

		if changes isnt false
			{
			checkedChanges = changes.Copy().Map({ it.lib $ ':' $ it.name })
			largeCheck = largeCheck.Copy().RemoveIf({
				not checkedChanges.Has?(it.lib $ ':' $ it.name) })
			qualityCheck = qualityCheck.Copy().RemoveIf({
				not checkedChanges.Has?(it.lib $ ':' $ it.name) })
			referencedDeleted = referencedDeleted.Copy().RemoveIf({
				not checkedChanges.Has?(it.lib $ ':' $ it.name) })
			}

		referencedDeleted = referencedDeleted.Empty?()
			? ''
			: "REFERENCES FOUND TO:\r\n" $ referencedDeleted.
				Map({ it.type $ it.lib $ ':' $ it.name }).Join('\r\n')

		return Opt('Record(s) larger than 100K: \n',
				largeCheck.Map({ it.lib $ ':' $ it.name }).Join('\n').Trim(), '\n') $
			Opt('Record(s) with quality issues: \n',
				qualityCheck.Map({ it.results }).Join().Trim(), '\n') $
			Opt(referencedDeleted, '\n') $
			.extra_check(changes)
		}
	}
