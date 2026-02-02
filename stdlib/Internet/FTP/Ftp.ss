// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(ftpServer, user = '', pass = '', acct = "", timeout = 60, verbose = false,
		big = false, type = 'ftp', options = #(), exeSubFolder = '')
		{
		Assert(type.Has?('ftp')) // could be ftp, sftp, ftps
		.ftp = big isnt false
			? new Curl(type, ftpServer, user, pass, timeoutConnect: timeout,
				:options, :exeSubFolder)
			: new FtpClient(ftpServer, user, pass, acct, timeout, verbose)
		}

	Get(remName, locName = '', retries = 0)
		{
		return .ftp.Get(:remName, :locName, :retries)
		}

	// only implemented GetMultiple/DeleteMultiple/DirMultiple for Curl
	GetMultiple(fileList, folderPath, receivingPath)
		{
		return .ftp.GetMultiple(fileList, folderPath, receivingPath)
		}

	DeleteMultiple(fileList, receivingPath, notFromRoot = false)
		{
		return .ftp.DeleteMultiple(fileList, receivingPath, :notFromRoot)
		}

	DirMultiple(directories)
		{
		return .ftp.DirMultiple(directories)
		}

	Put(locName, remName = '', retries = 0)
		{
		return .ftp.Put(:locName, :remName, :retries)
		}

	Del(fileName, path = '', notFromRoot = false)
		{
		return .ftp.Del(fileName, :path, :notFromRoot)
		}

	Ren(oldName, newName)
		{
		return .ftp.Ren(:oldName, :newName)
		}
	RenSFTP(oldName, newName)
		{
		return .ftp.RenSFTP(:oldName, :newName)
		}

	Dir(path = "*.*", files = false, details = false, caseSense = false,
		regExp = '', listOnly = true)
		{
		return .ftp.Dir(:path, :files, :details, :caseSense, :regExp, :listOnly)
		}
	}