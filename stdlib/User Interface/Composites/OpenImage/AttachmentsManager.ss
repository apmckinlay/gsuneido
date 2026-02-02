// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	creationNumField: false
	New(.query, .keyFields)
		{
		.oldAttachments = Object()
		.creationNumField = .findCreationNumField()
		}

	findCreationNumField()
		{
		if .query is ''
			return false
		potentialNums = .keyFields.Filter({ it.Has?('_num') })
		if potentialNums.Size() is 1
			return potentialNums[0]
		prefixes = .keyFields.Map({ it.BeforeFirst('_') $ '_num' })
		for key in QueryKeys(.query)
			{
			if key.Has?(",")
				continue
			// if field has _num but not _number
			if prefixes.Any?({ key.Prefix?(it) and not key.Prefix?(it $ 'ber') })
				return key
			}
		return false
		}

	KeyField()
		{
		return .keyFields
		}

	RecordDeleteAction: 'record delete'
	RemoveAction: 'remove'
	ReplaceAction: 'replace'
	QueueDeleteFile(new_file, old_file, rec, fieldName, action)
		{
		if not .normallyLinkCopy?()
			return false
		if .skipToDelete?(old_file, rec, fieldName, action)
			old_file = ''
		if old_file is '' and new_file is ''
			return true
		.oldAttachments.Add(Object(:new_file, :old_file, :rec, :fieldName, :action))
		return true
		}

	normallyLinkCopy?()
		{
		return OpenImageSettings.Normally_linkcopy?()
		}

	skipToDelete?(old_file, rec, fieldName, action)
		{
		if old_file is ""
			return false

		if .checkCreationDate(rec)
			return true

		old_file = Paths.ToStd(old_file)
		folders = Paths.ParentOf(old_file).Split('/')

		linkCopyPath = Paths.ToStd(.copyTo())
		if not old_file.Prefix?(linkCopyPath) or
			not .SubfolderMatches?(Paths.ParentOf(old_file)) or
			folders.Intersects?(.protectFolders())
			return true

		if action is .RecordDeleteAction
			return false

		if not .fileExist?(old_file)
			return true

		if .recordHasAttachment?(old_file, rec, fieldName)
			return true

		return false
		}

	SubfolderMatches?(subFolder)
		{
		subfolderPattern = '\d\d\d\d\d\d$'
		return subFolder =~ subfolderPattern
		}

	// Duplicate attachment links should no longer exist as of this date
	CreationCutOff: #20240416.0200
	checkCreationDate(rec)
		{
		if .creationNumField isnt false
			{
			creationDate = Date(rec[.creationNumField])
			if Date?(creationDate) and creationDate < .CreationCutOff
				return true
			}
		return false
		}

	copyTo()
		{
		return OpenImageSettings.Copyto()
		}

	protectFolders()
		{
		if false is c = OptContribution('AttachmentsManagerProtect', false)
			return Object()
		protectFolders = c()
		return Object?(protectFolders) ? protectFolders : Object()
		}

	recordHasAttachment?(old_file, rec, fieldName)
		{
		if not Object?(rec[fieldName])
			return false
		for row in rec[fieldName]
			{
			if not Object?(row)
				continue
			for att in row
				{
				fullPath = OpenImageWithLabelsControl.SplitFullPath(att)
				if .sameFile?(fullPath, old_file)
					return true
				}
			}
		return false
		}

	sameFile?(fullPath, old_file)
		{
		fullPath = Paths.ToStd(fullPath)
		old_file = Paths.ToStd(old_file)
		if .windows?()
			return fullPath.Lower() is old_file.Lower()
		else
			return fullPath is old_file
		}

	windows?()
		{
		return Sys.Windows?()
		}

	fileExist?(file)
		{
		if '' isnt bucket = AttachmentS3Bucket()
			return AmazonS3.FileExist?(bucket, FormatAttachmentPath(file))

		// Assume file exists even if error is thrown (attachments_cleanup_queue handles
		// cleaning it up later if there is a file access issue)
		result = true
		try	result = FileExists?(file)
		return result
		}

	ProcessQueue(restore? = false)
		{
		if restore?
			.oldAttachments.Each({ it.action = '' })
		.deleteQueuedFiles(restore?)
		}

	deleteQueuedFiles(restore?)
		{
		for fileOb in .oldAttachments
			{
			file = restore? ? fileOb.new_file : fileOb.old_file
			if fileOb.action isnt ''
				.logAction(file, fileOb.rec.Project(.keyFields),
					fileOb.fieldName, fileOb.action)
			.deleteFile(file, fileOb.rec, fileOb.fieldName)
			}
		.oldAttachments = Object()
		}

	deleteFile(file, rec, fieldName)
		{
		if file is ''
			return
		if '' isnt bucket = AttachmentS3Bucket()
			{
			if not AmazonS3.DeleteFile(bucket, FormatAttachmentPath(file))
				.addToDeleteQueue(file, rec, fieldName)
			return
			}
		// using DeleteFileApi to avoid retries on locked files
		// if there are several files locked / can't be deleted, then user can
		// have a long delay while system attempts to delete
		// failed deletes will be re-tried later by CleanupAttachments contribution
		if String?(result = DeleteFileApi(file)) and not result.Has?('does not exist')
			.addToDeleteQueue(file, rec, fieldName)
		}

	addToDeleteQueue(file, rec, fieldName)
		{
		failed = OptContribution('CleanupAttachments', .deleteFailed)
		failed(file, .query, .keyFields, rec.Project(.keyFields), fieldName)
		}

	logAction(file, key, fieldName, action)
		{
		table = QueryGetTable(.query, noThrow:)
		SuneidoLog('Attachment ' $ action, params: [:table, :key, keyFields: .keyFields,
			:fieldName, :file])
		}

	deleteFailed(file, query, keyFields, keys, fieldName)
		{
		params = Object(:file, :query, :keyFields, :keys, :fieldName)
		SuneidoLog('ERRATIC: Attachment cleanup/delete failed', :params)
		}

	// handles restore for a single record in a browse
	RestoreOneByKey(rec)
		{
		results = .oldAttachments.FindAllIf(find = { Same?(it.rec, rec) })
		for idx in results
			{
			ob = .oldAttachments[idx]
			.deleteFile(ob.new_file, rec, ob.fieldName)
			}
		.oldAttachments.RemoveIf(find)
		}

	QueueDeleteRecordFiles(rec)
		{
		if not .normallyLinkCopy?()
			return
		.handleStdAttachments(rec)
		.handleCustomAttachments(rec)
		}

	handleStdAttachments(rec)
		{
		for field in rec.MembersIf({ it.Suffix?('_attachments') })
			{
			if not Object?(rec[field])
				continue
			for fileInfo in rec[field].Flatten()
				.queueDeleteRecordFile(rec, field, fileInfo)
			}
		}

	queueDeleteRecordFile(rec, field, fileInfo)
		{
		if '' is old_file = OpenImageWithLabelsControl.SplitFullPath(fileInfo)
			return
		new_file = ''
		.QueueDeleteFile(new_file, old_file, rec, field, .RecordDeleteAction)
		}

	handleCustomAttachments(rec)
		{
		for customFld in rec.MembersIf({ Customizable.CustomField?(it, true) })
			if .customAttachmentField?(customFld) and rec[customFld] isnt ''
				.queueDeleteRecordFile(rec, customFld, rec[customFld])
		}

	customAttachmentField?(field)
		{
		dd = Datadict(field)
		return dd.Base?(Field_attachment)
		}

	DeleteNewRecordFiles(rec)
		{
		if not .normallyLinkCopy?()
			return
		.RestoreOneByKey(rec)
		}
	}
