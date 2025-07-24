// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(srcLib, srcNames, dstLib, dstFolder, print, overwrite? = false)
		{
		if srcNames.Empty?()
			return 0

		if '' isnt msg = .checkIfOverloaded(srcLib, srcNames, dstLib)
			return msg

		Transaction(update:)
			{ |t|
			try
				return .copyRecords(t, srcLib, srcNames, dstLib,
					.createDstFolder(t, dstLib, dstFolder.Trim('/')), print, overwrite?)
			catch (e)
				{
				t.Rollback()
				return e
				}
			}
		}

	createDstFolder(t, dstLib, dstFolder, parent = 0)
		{
		if dstFolder is ''
			return parent

		folder = dstFolder.BeforeFirst(`/`)
		if false is x = t.Query1(dstLib, name: folder, group: parent, :parent)
			t.QueryOutput(dstLib,
				x = [num: .nextNum(t, dstLib), name: folder, group: parent, :parent])
		return .createDstFolder(t, dstLib, dstFolder.AfterFirst('/'), x.num)
		}

	copyRecords(t, srcLib, srcNames, dstLib, folderNum, print, overwrite? = false)
		{
		counter = 0
		svcTable = SvcTable(dstLib)
		srcNames.Each()
			{ |name|
			if false is srcRec = t.Query1(srcLib, :name, group: -1)
				throw name $ " not found in " $ srcLib
			if false is rec = t.Query1(dstLib, :name, group: -1)
				svcTable.Output([
					name: srcRec.name,
					text: srcRec.text,
					lib_invalid_text: srcRec.lib_invalid_text
					parent: folderNum],
					:t)
			else
				{
				if not overwrite?
					rec.parent = folderNum
				rec.lib_invalid_text = srcRec.lib_invalid_text
				svcTable.Update(rec, :t, newText: srcRec.text)
				}
			counter++
			print("Copied", srcLib $ ':' $ name, "=>", dstLib)
			}
		return counter
		}

	checkIfOverloaded(srcLib, srcNames, dstLib)
		{
		libs = .libraries().Remove(dstLib).RemoveIf({ it.Suffix?('webgui') })
		for lib in Contributions('CopyToOverloadCheck')
			libs.Remove(lib)
		// srcLib is dstLib or srcLib is not in use
		if false is libi = libs.Find(srcLib)
			return ''

		msg = Object()
		for (i = libi + 1; i < libs.Size(); ++i)
			srcNames.Each()
				{ |name|
				if not QueryEmpty?(libs[i], :name, group: -1)
					msg.Add(name $ " is overloaded in " $ libs[i])
				}
		return msg.Join('\r\n')
		}

	nextNum(t, lib)
		{
		return t.QueryMax(lib, "num", 0) + 1
		}

	// extract for testing
	libraries()
		{
		return Libraries()
		}
	}
