// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
/* READ ME:
	This class is designed to read in a file's content and output a valid record
	If the imported record would overwrite a preexisting record, it can be overwritten.
	This class works with images, text files, etc.

	NOTE: skipSize should be used sparingly.
		It is only to be used when absolutely necessary
*/
class
	{
	title: Import
	CallClass(filename, table, destination, hwnd = 0, quiet = false, skipSize = false)
		{
		if false is text = .getFile(filename, hwnd, quiet, skipSize)
			return false

		svcTable = SvcTable(table)
		importRec = [name: Paths.Basename(filename), path: destination, :text]
		outputRec? = svcTable.Get(lookupName = svcTable.MakeName(importRec)) is false
		if outputRec?
			.outputNewRecord(svcTable, importRec, lookupName)
		else
			.updateExistingRecord(quiet lookupName, hwnd, svcTable, text)
		return outputRec?
		}

	getFile(filename, hwnd, quiet, skipSize)
		{
		if not skipSize and FileSize(filename) > 256.Kb() /*= size limit */
			err = 'too large (over 256kb)'
		else if false is text = GetFile(filename)
			err = 'inaccessible'
		else
			return text
		if quiet
			SuneidoLog('ERROR: ImportText: ' $ filename $ ' ' $ err)
		else
			Alert(Paths.Basename(filename) $ ' is ' $ err, title: .title,
				:hwnd, flags: MB.ICONERROR)
		return false
		}

	updateExistingRecord(quiet, lookupName, hwnd, svcTable, newText)
		{
		if not quiet and not YesNo(
			lookupName $ ' already exists.\nDo you want to replace it?',
			title: .title, :hwnd, flags: MB.ICONWARNING)
			return
		// Due to the alerts, we need to re-get the record afterwards.
		// SvcTable.Update needs the record lookup to be carried out in a transaction
		Transaction(update:)
			{ |t|
			rec = svcTable.Get(lookupName, :t)
			svcTable.Update(rec, :newText)
			}
		}

	outputNewRecord(svcTable, importRec, lookupName)
		{
		if committed = SvcCheckRestorability(svcTable, importRec)
			svcTable.Remove(lookupName, deleted:)
		svcTable.Output(importRec, :committed, force:)
		}
	}
