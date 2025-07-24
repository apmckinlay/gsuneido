// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
class
	{
	table: 'params'
	Save(name, params, user = false)
		{
		if user is false
			user = Suneido.User

		RetryTransaction()
			{ |t|
			t.QueryDo('delete ' $ .table $ ' where user is ' $ Display(user) $
				' and report is ' $ Display(name))
			t.QueryOutput(.table, Record(:user, report: name, :params))
			}
		}

	Load(name, user = false)
		{
		if user is false
			user = Suneido.User

		if false isnt rec = Query1(.table, :user, report: name)
			return rec.params
		return Record()
		}

	Sync(syncFromName, syncToName, syncMembers, forUser)
		{
		.syncUserParam(syncToName, syncFromName, syncMembers)
		.syncPresets(syncToName, syncFromName, syncMembers, forUser)
		}

	syncPresets(syncToName, syncFromName, syncMembers, forUser)
		{
		validPresets = Object()
		QueryApply(.buildSyncQuery(syncFromName))
			{ |from_rpt|
			Assert(from_rpt.report.Prefix?(syncFromName $ '~presets~'))
			presetName = from_rpt.report.AfterFirst(syncFromName)
			validPresets.Add(presetName)
			.sync(from_rpt, syncToName $ presetName, user: forUser, :syncMembers)
			}

		QueryApply(.buildSyncQuery(syncToName), update:)
			{ |to_rpt|
			Assert(to_rpt.report.Prefix?(syncToName $ '~presets~'))
			presetName = to_rpt.report.AfterFirst(syncToName)
			if not validPresets.Has?(presetName)
				to_rpt.Delete()
			}
		}

	buildSyncQuery(name)
		{
		// builds range on report name in order to use index
		return .table $ ' where report >= ' $  Display(name $ '~presets~') $ ' and ' $
			' report < ' $ Display(name $ '~presets~~')
		}

	sync(other, syncToName, user, syncMembers)
		{
		params = other.params.Project(syncMembers)
		found? = false
		QueryApply(.table $ .syncWhere(syncToName, user), update:)
			{ |x|
			if other.params_TS > x.params_TS
				{
				x.params.Merge(params)
				x.Update()
				}
			found? = true
			}
		if not found?
			QueryOutput(.table, [report: syncToName, :params, :user])
		}

	syncWhere(syncToName, user)
		{
		userWhere = syncToName.Has?('~presets~') ? '' : ' and user is ' $ Display(user)
		return ' where report is ' $ Display(syncToName) $ userWhere
		}

	syncUserParam(syncToName, syncFromName, syncMembers)
		{
		user = Suneido.User
		if false isnt other = Query1(.table, :user, report: syncFromName)
			.sync(other, syncToName, user, syncMembers)
		}
}
