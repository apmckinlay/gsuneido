// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Recs: (
		(name: 'Name1', lib: 'Lib1', lib_committed: #20241015, lib_modified: #20241015,
			text: 'function () { return Name3 + NameBuiltin }'),
		(name: 'Name2', lib: 'Lib1', lib_committed: #20241015, lib_modified: #20241016,
			text: 'Name4 { CallClass() { return Name1() + Name3 + NameSkip } }'),
		(name: 'Name3', lib: 'Lib1', lib_committed: #20241015, lib_modified: #20241017,
			text: '123'),
		(name: 'Name4', lib: 'Lib1', lib_committed: #20241015, lib_modified: #20241018,
			text: '1+1 /*INVALID*/'),
		(name: 'Name4_Override', lib: 'Lib1', lib_committed: #20241015,
			lib_modified: #20241018, text: 'class { }'),
		(name: 'NameSkip', lib: 'Lib1', text: 'function () { Name5 }',
			lib_committed: #20241015, lib_modified: '')
		)
	Test_buildBundle()
		{
		cl = SuCode
			{
			CollectDependencies()
				{
				return [recs: _newRecs]
				}
			SuCode_getBundleRecs()
				{
				return _oldRecs
				}
			SuCode_buildCode(rec)
				{
				return rec.name $ ': ' $ rec.lib $ '\r\n'
				}
			SuCode_update(@args)
				{
				_update.Add(args)
				}
			SuCode_getRec(name, lib)
				{
				return SuCode_Test.Recs.FindOne({
					it.name is name and it.lib is lib })
				}
			}
		fn = cl.SuCode_buildBundle

		_update = Object()
		_newRecs = Object()
		_oldRecs = Object()
		fn()
		.check(0, Date.Begin(), Date.Begin(), #(), #())

		_newRecs = [[name: 'Name1', lib: 'Lib1'], [name: 'Name2', lib: 'Lib1']]
		_oldRecs = [[name: 'Name1', lib: 'Lib1'], [name: 'Name2', lib: 'Lib1']]
		fn()
		.check(1, #20241016, #20241015, #(), #())

		_newRecs = [[name: 'Name1', lib: 'Lib1'], [name: 'Name2', lib: 'Lib1']]
		_oldRecs = []
		fn()
		.check(2, #20241016, #20241015, #(),
			#([name: 'Name1', lib: 'Lib1'], [name: 'Name2', lib: 'Lib1']))

		_newRecs = []
		_oldRecs = [[name: 'Name1', lib: 'Lib1'], [name: 'Name2', lib: 'Lib1']]
		fn()
		.check(3, Date.Begin(), Date.Begin(),
			#([name: 'Name1', lib: 'Lib1'], [name: 'Name2', lib: 'Lib1']), #())

		_newRecs = [[name: 'Name1', lib: 'Lib1'], [name: 'Name2', lib: 'Lib1'],
			[name: 'Name3', lib: 'Lib1'], [name: 'Name4', lib: 'Lib1']]
		_oldRecs = [[name: 'Name1', lib: 'Lib1'], [name: 'Name3', lib: 'Lib1']]
		fn()
		.check(4, #20241018, #20241015, #(),
			#([name: 'Name2', lib: 'Lib1'], [name: 'Name4', lib: 'Lib1']))

		_newRecs = [[name: 'Name1', lib: 'Lib1'], [name: 'Name3', lib: 'Lib1']]
		_oldRecs = [[name: 'Name1', lib: 'Lib1'], [name: 'Name2', lib: 'Lib1'],
			[name: 'Name3', lib: 'Lib1'], [name: 'Name4', lib: 'Lib1']]
		fn()
		.check(5, #20241017, #20241015,
			#([name: 'Name2', lib: 'Lib1'], [name: 'Name4', lib: 'Lib1']), #())

		_newRecs = [[name: 'Name1', lib: 'Lib1'], [name: 'Name2', lib: 'Lib1'],
			[name: 'Name3', lib: 'Lib1'], ]
		_oldRecs = [[name: 'Name2', lib: 'Lib1'], [name: 'Name3', lib: 'Lib1'],
			[name: 'Name4', lib: 'Lib1']]
		fn()
		.check(6, #20241017, #20241015,
			#([name: 'Name4', lib: 'Lib1']), #([name: 'Name1', lib: 'Lib1']))
		}

	check(i, modified, committed, toDelete, toAdd)
		{
		Assert(_update[i][1].lib_modified is: modified)
		Assert(_update[i][1].lib_committed is: committed)
		Assert(_update[i][2] is: toDelete)
		Assert(_update[i][3] is: toAdd)
		}

	Test_CollectDependencies()
		{
		cl = SuCode
			{
			Libraries: #(Lib1)
			SuCode_overrides: #(Name4: Name4_Override)
			SuCode_skips: #(NameSkip)
			SuCode_initSeeds(@unused)
				{
				return SuCode_Test.Recs.Take(2)
				}
			SuCode_queryRec(lib, name)
				{
				if false is rec = SuCode_Test.Recs.
					FindOne({ it.lib is lib and it.name is name })
					return false
				return rec.Copy()
				}
			}
		result = cl.CollectDependencies()
		Assert(result.recs.Map({ it.Project(#lib, #name) })
			equalsSet: #((lib: "Lib1", name: "Name1"),
				(lib: "Lib1", name: "Name2"),
				(lib: "Lib1", name: "Name4"),
				(lib: "Lib1", name: "NameSkip"),
				(lib: "Lib1", name: "Name3")))
		Assert(result.deps is: #(Name4: ("Name2"), Name3: ("Name2"), Name1: (),
			NameSkip: ("Name2"), Name2: (), NameBuiltin: ("Name1")))
		Assert(result.builtins is: #(NameBuiltin))
		}

	Test_checkUpdate()
		{
		cl = SuCode
			{
			CollectDependencies()
				{
				return _deps
				}
			SuCode_getBundleRecs(block)
				{
				#((lib: "Lib1", name: "Name1"),
					(lib: "Lib1", name: "Name2"),
					(lib: "Lib1", name: "Name3"),
					(lib: "Lib1", name: "Name4"),
					(lib: "Lib1", name: "NameSkip")).Each(block)
				}
			SuCode_getCurrentBundleBuildInfo()
				{
				return _buildInfo
				}
			SuCode_clear() { }
			}
		fn = cl.SuCode_checkUpdate

		_buildInfo = false
		Assert(fn() is: false)

		_buildInfo = [lib_committed: #20241015, lib_modified: #20241018]
		_deps = [recs: Object(
			#(name: "Name1", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241015),
			#(name: "Name2", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241016),
			#(name: "Name3", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241017),
			#(name: "Name4", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241018),
			#(name: "NameSkip", lib: "Lib1", lib_committed: #20241015,
				lib_modified: ""))]
		Assert(fn() is: false)

		// name not match
		_deps = [recs: Object(
			#(name: "Name1", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241015),
			#(name: "Name2", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241016),
			#(name: "NameWRONG", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241017),
			#(name: "Name4", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241018),
			#(name: "NameSkip", lib: "Lib1", lib_committed: #20241015,
				lib_modified: ""))]
		Assert(fn())

		// lib not match
		_deps = [recs: Object(
			#(name: "Name1", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241015),
			#(name: "Name2", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241016),
			#(name: "Name3", lib: "Lib2", lib_committed: #20241015,
				lib_modified: #20241017),
			#(name: "Name4", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241018),
			#(name: "NameSkip", lib: "Lib1", lib_committed: #20241015,
				lib_modified: ""))]
		Assert(fn())

		// lib_modified changed
		_deps = [recs: Object(
			#(name: "Name1", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241015),
			#(name: "Name2", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241016),
			#(name: "Name3", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241019),
			#(name: "Name4", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241018),
			#(name: "NameSkip", lib: "Lib1", lib_committed: #20241015,
				lib_modified: ""))]
		Assert(fn())

		// lib_committed changed
		_deps = [recs: Object(
			#(name: "Name1", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241015),
			#(name: "Name2", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241016),
			#(name: "Name3", lib: "Lib1", lib_committed: #20241018,
				lib_modified: #20241017),
			#(name: "Name4", lib: "Lib1", lib_committed: #20241015,
				lib_modified: #20241018),
			#(name: "NameSkip", lib: "Lib1", lib_committed: #20241015,
				lib_modified: ""))]
		Assert(fn())

		// changes restored
		_deps = [recs: Object(
			#(name: "Name1", lib: "Lib1", lib_committed: #20241015,
				lib_modified: ""),
			#(name: "Name2", lib: "Lib1", lib_committed: #20241015,
				lib_modified: ""),
			#(name: "NameWRONG", lib: "Lib1", lib_committed: #20241015,
				lib_modified: ""),
			#(name: "Name4", lib: "Lib1", lib_committed: #20241015,
				lib_modified: ""),
			#(name: "NameSkip", lib: "Lib1", lib_committed: #20241015,
				lib_modified: ""))]
		Assert(fn())
		}
	}