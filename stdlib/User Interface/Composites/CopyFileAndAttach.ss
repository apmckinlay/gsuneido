// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Title: 'Copy File And Attach'
	CallClass(file, copyto = false, deletesource? = false, quiet? = false,
		attachOb = false)
		{
		if file is ""
			return

		if copyto is false
			copyto = OpenImageSettings.Copyto()

		fileOb = .copyFile(copyto, file, quiet?)
		if String?(fileOb) and Object?(attachOb)
			{
			.BuildAttachOb(attachOb, fileOb.Tr('\n'))
			return
			}

		if deletesource?
			.deleteSource(file, fileOb.dest, quiet?)

		if attachOb isnt false and Object?(fileOb)
			.BuildAttachOb(attachOb, fileOb.fileBaseName)
		}

	copyFileMsg: 'Could not copy the file.\n\n' $
		'There was a problem accessing the source file'
	copyFile(copyto, file, quiet?)
		{
		fileInfoOb = .buildFileOb(file, copyto, quiet?)

		if fileInfoOb.GetDefault(#error, false) isnt false
			{
			if not quiet?
				Alert(fileInfoOb.error, .Title, flags: MB.ICONWARNING)
			return fileInfoOb.error
			}

		if fileInfoOb.copyto is "" or true isnt .copyFileToStorage(file, fileInfoOb.dest)
			{
			if not quiet?
				Alert(.copyFileMsg,
					.Title, flags: MB.ICONWARNING)
			return .copyFileMsg
			}
		return fileInfoOb
		}

	copyFileToStorage(file, fileTo)
		{
		if '' isnt bucket = AttachmentS3Bucket()
			{
			region = AmazonS3.GetBucketLocationCached(bucket)
			if not AmazonS3.FileExist?(bucket, file)
				return AmazonS3.PutFile(
					bucket, file, FormatAttachmentPath(fileTo), :region)
			else
				return AmazonS3.CopyFile(
					bucket, file, bucket, FormatAttachmentPath(fileTo))
			}

		return CopyFile(file, fileTo, true)
		}

	buildFileOb(file, copyto, quiet?)
		{
		destFileName = file
		if copyto is ""
			return Object(fileBaseName: file).Set_default('')

		// this is when copyto is a designated file path (a copyto from a company
		// attachment path always has a trailling slash)
		if copyto[-1] not in (`\`, `/`)
			{
			destFileName = Paths.Basename(copyto)
			copyto = Paths.EnsureTrailingSlash(Paths.ParentOf(copyto))
			}

		fileBasename = Paths.Basename(destFileName)
		if '' isnt msg = CheckFileName.WithErrorMsg(fileBasename)
			return Object(error: msg)
		copyToFolder = .buildCopyToFolder(copyto, quiet?)

		if Object?(copyToFolder)
			return Object(error: copyToFolder.msg)

		file = MakeUniqueFileName(copyToFolder, fileBasename)
		// NEED to ensure that fileBaseName has the renamed file, not the original
		fileBasename = Paths.Basename(file.dest)

		return Object(:copyto, dest: file.dest, base: file.base, ext: file.ext,
			fileBaseName: .SubFolder() $ fileBasename)
		}

	SubFolder()
		{
		return Paths.EnsureTrailingSlash(Date().Format('yyyyMM'))
		}

	buildCopyToFolder(copyto, quiet?)
		{
		copyfolder = copyto $ .SubFolder()
		return .EnsureDirExists(copyfolder, copyto, quiet?)
		}

	EnsureDirExists(copyfolder, copyto, quiet?)
		{
		if AttachmentS3Bucket() isnt ''
			return copyfolder

		if String?(exists = CheckDirExists(copyfolder))
			return Object(msg: exists)
		if not exists
			if not .createDirectory(copyfolder)
				{
				msg = 'System could not access or create folder: ' $ copyfolder $
					' to place copy of current file'
				.copyWarning(msg, quiet?)
				copyfolder = copyto
				}
		return copyfolder
		}

	createDirectory(dir)
		{
		dir = dir.RightTrim('\\/')
		if CheckDirectory.IsFullUNC?(dir)
			return true is CreateDir(dir)

		return true is ServerEval('CreateDir', dir)
		}

	deleteSource(file, dest, quiet?)
		{
		if true is DeleteFile(file)
			return dest
		DeleteFile(dest)
		msg = 'Can not delete source file: ' $ file $
			(quiet? is true ? ' err: delete retry failed' : '')
		.copyWarning(msg, quiet?)
		}

	copyWarning(msg, quiet?)
		{
		if quiet?
			SuneidoLog('INFO: CopyFileAndAttach: ' $ msg,
				caughtMsg: 'quiet? is true, no alert')
		else
			Alert(msg, .Title, flags: MB.ICONWARNING)
		}

	BuildAttachOb(ob, fileName, noPath? = false)
		{
		if not Object?(ob)
			ob = Object()
		next = UnlimitedAttachments.GetNextAvailPos(ob)
		if not ob.Member?(next.row)
			ob[next.row] = Record()
		gpg? = fileName.AfterLast('.') is 'gpg'
		ob[next.row]['attachment' $ next.pos] = noPath?
			? Paths.Basename(fileName)
			: gpg? is true
				? fileName.BeforeLast('.')
				: fileName
		}
	}
