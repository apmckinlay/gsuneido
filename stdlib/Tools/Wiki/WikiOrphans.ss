// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		ob = .collectNames()
		.removeReferenced(ob)
		for name in ob.Members()
			{
			QueryApply1('wiki', :name)
				{
				it.orphaned = true
				it.Update()
				}
			}
		.ResetList()
		}
	// StartPage is expected to be the top.
	// RecentChanges is a special page
	// EtaNotes is ETA Wiki notes
	specialPages: #('^StartPage$', '^RecentChanges$', '^EtaNotes')
	collectNames()
		{
		ob = Object()
		// don't bother checking ones we have already flagged
		QueryApply('wiki where orphaned isnt true')
			{|x|
			if .specialPages.HasIf?({ x.name =~ it})
				continue
			ob[x.name] = true
			}
		return ob
		}
	removeReferenced(ob)
		{
		rx = `\<[A-Z][a-z0-9]+[A-Z][a-zA-Z0-9]*\>`
		// these lists are dynamically created, need to check them seperatly
		for m in GetContributions('WikiOrphans')
			Global(m)().Replace(rx, { ob.Delete(it);; })
		QueryApply('wiki')
			{|x|
			x.text.Replace(rx, { ob.Delete(it);; })
			}
		}

	OrphanedMessage: '(orphan)'

	listOrphans: MemoizeSingle
		{
		Func()
			{
			QueryList('wiki where orphaned is true', 'name')
			}
		}

	ListOrphans()
		{
		if Client?()
			return ServerEval('WikiOrphans.ListOrphans')
		return (.listOrphans)()
		}

	ResetList()
		{
		if Client?()
			return ServerEval("WikiOrphans.ResetList")
		(.listOrphans).ResetCache()
		return 0 // Prevent Can't Pack LruCache error
		}
	}
