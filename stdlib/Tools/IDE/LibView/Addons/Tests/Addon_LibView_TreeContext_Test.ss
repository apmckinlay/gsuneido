// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_buildContextMenu()
		{
		cl = new Addon_LibView_TreeContext([], [])
		m = cl.Addon_LibView_TreeContext_buildContextMenu

		result = m('FolderName', recordSelected?: false, rootSelected?: false)
		Assert(result is: [[order: 40, name: "&New"], ["&Folder", "&Item", order: 41]])

		result = m('(TableName)', recordSelected?: false, rootSelected?:)
		Assert(result
			is: [
				[name: '&Delete', order: 30],
					['&Delete Library', order: 31],
				[name: '', order: 90],
				[name: 'Dump', order: 91],
				[name: 'Import Records...', order: 92],
				[name: 'Use', order: 93],
				[name: 'Undelete...', order: 94],
				[name: '', order: 110],
				[name: 'Check Records', order: 111],
				[order: 40, name: "&New"],
					["&Folder", "&Item", order: 41]
			])

		result = m('TableName', recordSelected?: false, rootSelected?:)
		Assert(result
			is: [
				[name: '&Delete', order: 30],
					['&Delete Library', order: 31],
				[name: '', order: 90],
				[name: 'Dump', order: 91],
				[name: 'Import Records...', order: 92],
				[name: 'Unuse', order: 93],
				[name: 'Undelete...', order: 94],
				[name: '', order: 110],
				[name: 'Check Records', order: 111],
				[order: 40, name: "&New"],
					["&Folder", "&Item", order: 41]
			])

		result = m('TableName', recordSelected?:, rootSelected?: false)
		Assert(result
			is: [
				[name: '&Run', def:, order: 1],
				[name: '', order: 2],
				[name: '', order: 80],
				[name: 'Export Record', order: 81],
				[order: 40, name: "&New"],
					["&Folder", "&Item", order: 41]
			])
		}
	}
