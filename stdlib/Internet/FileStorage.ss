// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Exists?(file)
		{
		if false isnt fileInfo = .s3?(file)
			return AmazonS3.FileExist?(fileInfo.bucket, fileInfo.file)
		else
			return FileExists?(file)
		}

	FileSize(file)
		{
		if false isnt fileInfo = .s3?(file)
			return AmazonS3.FileSize(fileInfo.bucket, fileInfo.file)
		else
			return FileSize(file)
		}

	s3?(file)
		{
		if file is '' or file.Prefix?(GetAppTempPath())
			return false
		if '' is bucket = AttachmentS3Bucket()
			return false
		if Paths.ToStd(file) is s3File = FormatAttachmentPath(file)
			return false
		return Object(:bucket, file: s3File)
		}

	GetAccessibleFilePath(file)
		{
		if false is s3Info = .s3?(file)
			return file

		cache = .getCachedLocalFile
		path = cache(s3Info.bucket, s3Info.file)
		if FileExists?(path)
			return path

		cache.ResetCache()
		return cache(s3Info.bucket, s3Info.file)
		}

	getCachedLocalFile: Memoize
		{
		Func(bucket, file)
			{
			// temp files will be cleaned up by scheduled task
			tmpPath = Paths.Combine(GetAppTempPath(), String(Timestamp()).Tr('#.'))
			EnsureDir(tmpPath)
			tmp = Paths.Combine(tmpPath, Paths.Basename(file))
			region = AmazonS3.GetBucketLocationCached(bucket)
			AmazonS3.GetFile(bucket, file, tmp, :region)
			return tmp
			}
		}

	CopyFile(from, to, failIfExists?)
		{
		if false is s3Info = .s3?(from)
			return CopyFile(from, to, failIfExists?)

		// TODO: handle failIfExists?
		toS3 = FormatAttachmentPath(to)
		res = AmazonS3.CopyFile(s3Info.bucket, s3Info.file, s3Info.bucket, toS3)
		if res is false
			return throw 'AmazonS3.CopyFile failed'
		return res
		}

	SaveFile(file, dest = false)
		{
		if '' is bucket = AttachmentS3Bucket()
			return dest is false
				? OpenImageSelect.ResultFile(file, useSubFolder:)
				: true is CopyFile(file, dest, false)
					? dest
					: false

		if dest is false
			{
			copyfolder = OpenImageSelect.SubFolder()
			fileBasename = Paths.Basename(file)
			dest = OpenImageSelect.GetCopyToFilename(copyfolder, fileBasename)
			}

		dest = FormatAttachmentPath(dest)
		if false is s3Info = .s3?(file) // file is local
			{
			region = AmazonS3.GetBucketLocationCached(bucket)
			if true is AmazonS3.PutFile(bucket, file, dest, :region)
				return dest
			return false
			}
		else // file is already on s3
			{
			if true is AmazonS3.CopyFile(bucket, s3Info.file, bucket, dest)
				return dest
			return false
			}
		}

	Dir(path, files = false, details = false)
		{
		if false is s3Info = .s3?(path)
			return Dir(path, :files, :details)

		s3Path = FormatAttachmentPath(path)
		return AmazonS3.Dir(s3Info.bucket, s3Path, :details)
		}
	}
