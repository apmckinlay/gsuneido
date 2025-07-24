// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
/*
SvcBook is used to get information from and update local book tables. SvcLibrary and
SvcBook handle the differences between libraries and books.

NOTE: This should never be called directly. Use SvcTable as a interface for this class
*/
class
	{
	Type: 'book'
	book: false
	New(table)
		{
		.book = table
		}

	SvcEnabled?()
		{
		return not SvcDisabledBooks().Has?(.book)
		}

	deletedPrefix: '<deleted>'
	Query(deleted = false)
		{
		criteria = ' path.Prefix?("' $ .deletedPrefix $ '")'
		deletedWhere = ' where ' $ (deleted ? '' : 'not') $ criteria
		return .book $ .Where() $ deletedWhere
		}

	Where()
		{
		where = ' where not path.Prefix?("commitRec")'
		if .book is OptContribution('Svc_IncludeReporterReports', '') or
			not QueryColumns(.book).Has?('plugin')
			return where
		return where $ ' and plugin isnt true and
			not path.Has?("Reporter Reports") and
			not path.Has?("Reporter Forms")'
		}

	Table()
		{
		return .book
		}

	NameQuery(pathname, deleted = false)
		{
		name = Display(.splitName(pathname))
		path = Display(.deletedPath(.splitPath(pathname), :deleted))
		return .Query(:deleted) $ ' where name is ' $ name $ ' and path is ' $ path
		}

	GetData(rec, t /*unused*/ = false)
		{
		if not BookResource?(rec.name) and rec.path !~ '^/res\>'
			rec.text = 'Order: ' $ rec.order $ '\r\n\r\n' $ rec.text
		}

	MakeName(rec, clean = false)
		{
		if clean
			rec.path = .clean(rec.path)
		return rec.path $ '/' $ rec.name
		}

	clean(string)
		{
		return string.RemovePrefix(.deletedPrefix)
		}

	splitName(pathname)
		{
		i = pathname.FindLast('/')
		return i is false ? pathname : pathname[i + 1 ..]
		}

	splitPath(pathname)
		{
		return pathname[.. pathname.FindLast('/')]
		}

	GetPath(rec)
		{
		return rec.path
		}

	FormatData(rec, t /*unused*/, deleted = false)
		{
		rec.path = .deletedPath(rec.path, deleted)
		rec.name = .splitName(rec.name)
		.SplitText(rec)
		}

	FormatText(text, name)
		{
		return BookResource?(name)
			? text
			: text.Trim()
		}

	deletedPath(path, deleted)
		{
		if deleted
			{
			if not .Deleted?([:path])
				path = .deletedPrefix $ path
			}
		else
			path = .clean(path)
		return path
		}

	SplitText(rec)
		{
		if BookResource?(.MakeName(rec)) or not rec.text.Prefix?('Order:')
			return
		order = rec.text.BeforeFirst('\n')
		rec.text = rec.text.AfterFirst(order.Suffix?('\r') ? '\r\n\r\n' : '\n\n')
		rec.order = order.AfterFirst('Order:').Trim()
		if rec.order isnt ''
			rec.order = Number(rec.order)
		}

	Move(movedRec, newParent, t)
		{
		newPath = newParent is -1 ? '' : .MakeName(t.Query1(.book, num: newParent))
		.move(movedRec, newPath, t)
		}

	move(movedRec, newPath, t)
		{
		// Look for potential restores prior to move
		svcTable = SvcTable(.book)
		restore = .potentialRestore(movedRec.name, newPath, svcTable, t)
		childrenQuery = .childrenQuery(delete = .MakeName(movedRec))

		// Delete ONLY this record, children records will be moved / deleted as
		// they are processed during their own moves.
		_skipChildren = true
		svcTable.StageDelete(delete, t, newNum:)

		// Output a new record to be the parent of the children records
		.processRecForMove(movedRec, restore, newPath)
		svcTable.VerifyModified(movedRec)
		svcTable.Output(movedRec, t, committed: movedRec.lib_modified is '')

		// Move all the children record to the new parent
		newPath = .MakeName(movedRec)
		t.QueryApply(childrenQuery)
			{
			.move(it, newPath, t)
			}
		}

	potentialRestore(name, path, svcTable, t)
		{
		potentialRestore = .MakeName([:path, :name])
		return svcTable.Get(potentialRestore, :t, deleted:)
		}

	childrenQuery(name)
		{
		return .book $ ' where path is ' $ Display(.clean(name))
		}

	processRecForMove(movedRec, restoreRec, newPath)
		{
		if restoreRec isnt false
			{
			movedRec.lib_committed = restoreRec.lib_committed
			movedRec.Delete(#lib_before_path)
			restoreRec.Delete()
			}
		else
			{
			movedRec.lib_before_path = movedRec.path
			movedRec.Delete(#lib_committed)
			}
		movedRec.path = newPath
		}

	Rename(rec, oldName, t)
		{
		t.QueryApply(.book $ ' where path is ' $ Display(oldName))
			{
			.move(it, rec.name, t)
			}
		}

	StageDelete(rec, t, newNum = false, _skipChildren = false)
		{
		if skipChildren
			return
		svcTable = SvcTable(.book)
		t.QueryApply(.childrenQuery(.MakeName(rec)))
			{
			svcTable.StageDelete(.MakeName(it, clean:), :t, :newNum)
			}
		}

	Deleted?(rec)
		{
		return rec.path.Prefix?(.deletedPrefix)
		}

	FormatMaxCommitRecord(rec)
		{
		rec.path = 'commitRec'
		}

	MaxCommitQuery()
		{
		return .book $
			' where name is "' $ SvcTable.MaxCommitName $ '" and path is "commitRec"'
		}

	Compare(local, master, discrepancies)
		{
		.SplitText(local)
		.SplitText(master)
		if Svc.Hash(local.text) isnt Svc.Hash(master.text)
			discrepancies.Add(#text)
		if local.order isnt master.order
			discrepancies.Add(#order)
		}

	CheckRecord(rec, type)
		{
		if type is '-' or BookResource?(.MakeName(rec))
			return #()
		results = Object()
		if rec.text =~ '\r[^\n]' or rec.text.Suffix?('\r')
			results.Add('Please ensure that this record does ' $
				'not have any invalid newlines.')
		if not rec.text.Prefix?('<!-- skipXmlCheck -->')
			try
				XmlCheckNest(rec.text)
			catch (e)
			results.Add(e)
		return results
		}

	Dir(num) // not sure about this, since folders have text
		{
		rec = Query1(.book, :num)
		query = .book $
			' where path = ' $ Display(rec.path $ '/' $ rec.name) $
			' sort path, order, name'
		return QueryAll(query).Map!({ it.name }).Join('\r\n')
		}

	ResetSvcDisabledCache()
		{
		SvcDisabledBooks.ResetCache()
		}

	// Whenever a image book record is modified or restored, we need to clear the
	// applicable caches. Otherwise, the image displayed via the Book controls will not
	// match the actual image stored in the table
	PostRestore(name)
		{
		if BookResource?(name, imagesOnly?:)
			ClearBookImageCache(.book, clearQuery1Cache?:)
		}
	}
