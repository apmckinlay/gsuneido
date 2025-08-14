// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	DeviceErr: 'device is not currently connected|device not ready|permission denied' $
		'|unauthenticated guest access'
	IsDrive?(path)
		{
		return path =~ '^[A-Za-z]:'
		}

	IsFullUNC?(path)
		{
		return path.Prefix?(`\\`) or path.Prefix?('//')
		}

	IsLinuxAbsolutePath?(path)
		{
		return path.Prefix?(`/`) and not path.Prefix?(`//`)
		}

	Accessible?(path)
		{
		if path is ''
			return true
		file = .buildPath(path, .file())
		result = false
		if .writeFile(file, s = 'test')
			{
			if GetFile(file) is s
				result = true
			if true isnt DeleteFile(file)
				QueryOutputIfNew(.table, [cd_path: path])
			}
		return result
		}

	file()
		{
		uniqueTxt = Display(Timestamp()).RemovePrefix('#').Replace('\.', '_')
		return .FilePrefix() $ uniqueTxt.RightFill(18 /*= full TS length*/, '0')
		}

	FilePrefix()
		{
		return OptContribution(#CheckDirectoryPrefix, #suneido_dir_test_)
		}

	buildPath(path, file)
		{
		return Paths.ToLocal(Paths.Combine(path, file))
		}

	table: checked_directories
	writeFile(file, s)
		{
		try
			{
			PutFile(file, s)
			return true
			}
		return false
		}

	InvalidChars: `\/*?":<>|`
	ReservedNames: #(CON, PRN, AUX, NUL, COM1, COM2, COM3, COM4, COM5,
		COM6, COM7, COM8, COM9, LPT1, LPT2, LPT3, LPT4, LPT5, LPT6, LPT7, LPT8, LPT9)
	ValidFileName?(path)
		{
		pos = path.FindLast1of("\\/") // : is not valid in file name - want to find these
		filename = pos is false ? path : path[pos + 1..]
		if filename.Blank?() or filename.BeforeFirst('.').Blank?() or
			filename.Find1of(.InvalidChars) isnt filename.Size() or
			.ReservedNames.Has?(filename.BeforeFirst('.'))
			return false
		return true
		}

	ValidFilePath(path)
		{
		if not .ValidFileName?(path)
			return 'File Names may not contain ' $ .InvalidChars $ '\n' $
				'Please enter another File Name.'

		if not .Accessible?(Paths.ParentOf(path))
			return 'You do not have permission to write to the specified destination.\n' $
				'Please enter another location to save the File Name.'
		return ''
		}

	maxStuckDays: 2
	deleteFailure: 'failed to delete'
	ReviewPaths()
		{
		failures = Object()
		failResult = 'WARNING'
		startTime = .startTime()
		filePrefix = .FilePrefix()
		QueryAll(.table).Each()
			{|rec|
			try
				if not .deleteOutstanding(rec.cd_path, startTime, filePrefix)
					throw .deleteFailure
			catch (error)
				if .checkFailed(startTime, failures, rec, error)
					failResult = 'FAILED'
			}
		return failures.Empty?()
			? 'SUCCEEDED'
			: failResult $ '\r\n- ' $ failures.Join('\r\n- ')
		}

	startTime()
		{
		return Timestamp()
		}

	dirFailure: 'path inaccessible'
	deleteOutstanding(path, startTime, filePrefix)
		{
		allDeleted? = true
		if not .dirExists?(path)
			throw .dirFailure
		.dir(dirPath = .buildPath(path, filePrefix)).Each()
			{
			if .deleteFile?(it.name, it.date, startTime, filePrefix) and
				true isnt .deleteFileApi(.buildPath(path, it.name))
				allDeleted? = false
			}
		if allDeleted?
			.deletePathRec(path)
		// Re-output the checked_directories record if there are any new test files.
		// This ensures these files are cleaned up next run and the cd_TS is updated.
		if .dir(dirPath).NotEmpty?()
			QueryOutputIfNew(.table, [cd_path: path])
		return allDeleted?
		}

	dirExists?(path)
		{
		return DirExists?(path)
		}

	dir(path)
		{
		return Dir(path $ '*.*', files:, details:)
		}

	deleteFile?(filename, fileDate, startTime, filePrefix)
		{
		endRegex = '^\d\d\d\d\d\d\d\d_\d\d\d\d\d\d\d\d\d$'
		return fileDate < startTime
			? filename.AfterLast(filePrefix).Match(endRegex) isnt false
			: false
		}

	/* NOTE: .ReviewPaths purposely uses DeleteFileApi in order to avoid the retry.
	The general assumption is that the file is stuck, if DeleteFile fails
	during .Accessible?. .ReviewPaths is designed to be run intermittently
	at a later date. During which, it re-attempts the delete,
	assuming that the file is no longer stuck.
	*/
	deleteFileApi(path)
		{
		return DeleteFileApi(path)
		}

	deletePathRec(path)
		{
		QueryDo('delete ' $ .table $ ' where cd_path is ' $ Display(path))
		}

	checkFailed(startTime, failures, rec, error)
		{
		daysStuck = startTime.MinusDays(rec.cd_TS)
		if error is .deleteFailure
			failures.Add('Path: ' $ rec.cd_path $ ', Days Stuck: ' $ daysStuck)
		else
			.deletePathRec(rec.cd_path)
		return daysStuck > .maxStuckDays
		}
	}
