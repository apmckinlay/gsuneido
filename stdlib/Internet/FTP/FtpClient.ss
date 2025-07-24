// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// A basic FTP-client to access ftp-servers, only supporting passive transfer mode
// Based on RFC 959: File Transfer Protocol (FTP), by J. Postel and J. Reynolds
class
	{
	scCtrl:			false
	scData:			false
	verbose:		false
	supportsMDTM:	true
	remCWD:			"/"
	resText:		""
	resCode:		""
	realPath:		""
	realBase:		""
	portNo:			0

	New(ftp_server, user = "anonymous", pass = "anonymous@", acct = "",
		timeout = 60, .verbose = false)
		{
		Assert(String?(user) and String?(pass) and String?(acct))
		Assert(Number?(timeout) and timeout > 0)
		Assert(Boolean?(verbose) or Function?(verbose))
		.scCtrl = SocketClient(ftp_server, 21, timeout) /*= default port number */
		.getResponse()
		.checkResponse("220", "error connecting to " $ ftp_server)
		.loginAccount(user, pass, acct)
		.sendCommand("TYPE I")	// don't realy need ASCII type, so to avoid switching,
		.checkResponse("200")	// just stick with IMAGE (SIZE command needs this also)
		.ChDir()				// get current directory
		}

	loginAccount(user, pass, acct)
		{
		if user isnt ""
			{
			if .sendCommand("USER " $ user) is "331"
				if .sendCommand("PASS " $ pass) is "332" or acct isnt ""
					.sendCommand("ACCT " $ acct)
			.checkResponse("230", "error logging in")
			}
		}

	ChDir(name = "")
		{
		// Changes the working directory on the server
		// Returns the (new) working directory
		// Use wo name, "" or "." to get the current working directory
		Assert(String?(name))
		if name isnt "" and name isnt "."
			{
			.splitAndSetPath(name)
			.sendCommand("CWD " $ .realPath $ .realBase)
			.checkResponse("250")
			}
		.sendCommand("PWD")
		.checkResponse("257")
		if false is (.remCWD = .resText.Extract('"(.*)"'))
			throw "FTP: malformed PWD response: " $ .resText
		.remCWD $= .remCWD.Suffix?('/') ? '' : '/'
		return .remCWD
		}

	Del(fileName)				// Delete file
		{
		Assert(String?(fileName))
		.splitAndSetPath(fileName)
		.sendCommand("DELE " $ .realPath $ .realBase)
		.checkResponse("250")
		}

	Dir(path = "*.*", files = false, details = false)
		{
		// Works like the standard Dir()
		// The .attr member (if details is true) is mainly supplied to be
		// compatitble with Dir(), only the directory flag is valid.

		Assert(String?(path))
		Assert(Boolean?(files) and Boolean?(details))
		.remDir(path, files, details, caseSense: false)
		}

	Get(remName, locName = "", resume = false, callBack = false, bufSize = 1024)
		{
		// Downloads the file remName from the server.
		// If locName has no basename, the remote filename is used.
		// If resume is true, the local file exists and its size is less than that of
		// the remote, only the remainder is downloaded (if supported by the server).
		// A method or function can be supplied as callBack: before each transmission
		// of bufSize bytes (or the remainder), the number of bytes remaining are
		// passed on to callBack. True should then be returned by callBack to proceed
		// or false to abort further downloading.

		Assert(String?(remName) and String?(locName) and Boolean?(resume))
		Assert(Number?(bufSize) and bufSize > 0)
		Assert(callBack is false or Function?(callBack))

		if false isnt (remName $ locName).Extract("[*?]+")
			throw "FTP: no wildcards allowed in path or filename"
		dir = .remDir(remName, details:, files:)
		if dir.Size() isnt 1
			throw "FTP: no such remote file"
		locName $= Paths.Basename(locName) is "" ? dir[0].name : ""
		remName = .realPath $ dir[0].name
		fileSize = dir[0].size
		dir = Dir1(locName, files:)
		if dir isnt false and dir isnt Paths.Basename(locName)
			throw "FTP: local filename conflict (" $ dir $ ")"

		.download(remName, locName, resume, callBack, bufSize, dir, fileSize)
		}

	download(remName, locName, resume, callBack, bufSize, dir, fileSize)
		{
		offset = 0
		if resume
			{
			if dir.Size() isnt 1
				throw "FTP: no such local file to resume"
			offset = dir[0].size
			if offset >= fileSize
				throw "FTP: local filesize not less remote"
			}
		.createDataSocket()
		if offset isnt 0 and .sendCommand("REST " $ offset) isnt "350"
			offset = 0			// server doesn't support resuming
		File(locName, offset is 0 ? "w" : "a")
			{|f|
			.sendCommand("RETR " $ remName)
			.checkResponse("150|125")
			while (offset < fileSize)
				{
				if (false isnt callBack and not callBack(fileSize - offset))
					{
					.sendCommand("ABOR")
					break
					}
				buffer = .scData.Read(Min(bufSize, fileSize - offset))
				f.Write(buffer)
				offset += buffer.Size()
				}
			}
		.closeDataSocket()
		}

	MkDir(name)					// Make a directory
		{
		Assert(String?(name))
		.splitAndSetPath(name)
		.sendCommand("MKD " $ .realPath $ .realBase)
		.checkResponse("250|257")
		}

	Put(locName, remName = "", resume = false, callBack = false, bufSize = 1024)
		{
		// Uploads the file locName to the server.
		// As Get() with locName <-> remName

		Assert(String?(locName) and String?(remName) and Boolean?(resume))
		Assert(Number?(bufSize) and bufSize > 0)
		Assert(callBack is false or Function?(callBack))

		if (false isnt (remName $ locName).Extract("[*?]+"))
			throw "FTP: no wildcards allowed in path or filename"
		dir = Dir1(locName, details:, files:)
		if dir is false
			throw "FTP: no such local file"
		if Paths.Basename(remName) is ""
			{
			if (dir.name isnt Paths.Basename(locName).RightTrim())
				throw "FTP: local filename conflict (" $ dir.name $ ")"
			remName $= dir.name
			}
		.splitAndSetPath(remName)
		.upload(locName, .realPath $ .realBase, resume, callBack, bufSize, dir.size)
		}

	upload(locName, remName, resume, callBack, bufSize, filesize)
		{
		offset = 0
		if resume
			{
			dir = .remDir(remName, details:, files:)
			if dir.Size() isnt 1
				throw "FTP: no such remote file to resume"
			offset = dir[0].size
			if offset >= filesize
				throw "FTP: remote filesize not less local"
			}
		File(locName, 'r')
			{|f|
			.createDataSocket()
			if (offset isnt 0 and .sendCommand("REST " $ offset) is "350")
				f.seek(offset)		// server supports resuming
			else
				offset = 0
			.sendCommand("STOR " $ remName)
			.checkResponse("150|125")
			while (offset < filesize)
				{
				if (false isnt callBack and not callBack(filesize - offset))
					{
					.sendCommand("ABOR")
					break
					}
				buffer = f.Read( Min(bufSize, filesize - offset))
				.scData.Write(buffer)
				offset += buffer.Size()
				}
			.closeDataSocket()
			}
		}

	Ren(oldName, newName)		// Rename file or directory
		{
		Assert(String?(oldName) and String?(newName))
		newName = newName.RightTrim()
		if newName isnt Paths.Basename(newName)
			throw "FTP: new name cannot have a path"
		dir = .remDir(oldName)
		if dir.Size() isnt 1
			throw "FTP: no such remote file or directory"
		if .realBase is newName
			return
		.sendCommand("RNFR " $ .realPath $ .realBase)
		.checkResponse("350")
		.sendCommand("RNTO " $ newName)
		.checkResponse("250")
		}

	RmDir(name)					// Remove a directory
		{
		Assert(String?(name))
		dir = .remDir(name)
		if dir.Size() isnt 1 or not dir[0].Suffix?('/')
			throw "FTP: no such remote directory"
		.sendCommand("RMD " $ .realPath $ dir[0][.. -1])
		.checkResponse("250")
		}

	checkResponse(prefix, msg = "communication error")
		{
		if false is .resCode.Match("^" $ prefix)
			if .resCode is "553"
				throw "FTP: " $ .resText[4 ..]
			else
				throw "FTP: " $ msg $ " (Expected " $ prefix $ ", got " $ .resCode $ ')'
		}

	closeDataSocket()
		{
		.scData.Close()
		.scData = false
		.getResponse()
		.checkResponse("226|250")
		}

	createDataSocket()
		{
		.sendCommand("PASV")
		.checkResponse("227")

		ipAddress = ""
		portNo = i = 0
		if (false isnt ipSpecs = .resText.Extract("\d+,\d+,\d+,\d+,\d+,\d+"))
			for part in ipSpecs.Split(",")
				if part.Number?()
					{
					val = Number(part)
					if (val >= 0 and val <= 255)
						if (++i < 5)
							ipAddress $= val $ "."
						else
							portNo = (portNo << 8) + val
					}
		if i isnt 6
			throw ("FTP: malformed PASV response: " $ .resText)
		.portNo = portNo
		.scData = SocketClient(ipAddress[.. -1], portNo)
		}

	doVerbose(string)
		{
		if true is .verbose
			Print(string)
		else if Function?(.verbose)
			.verbose(string)
		}

	getResponse()
		{
		.resText = resp = .scCtrl.Readline()
		.resCode = .resText[.. 3]
		if .resCode !~ "^[1-5]\d\d"
			throw "FTP: malformed response: " $ .resText
		while (not resp.Prefix?(.resCode $ ' '))
			{
			resp = .scCtrl.Readline()
			.resText $= '\n' $ resp
			}
		.doVerbose(.resText)
		}

	parseDir(line)
		{
//		Parses one line of a directory listing received from the server.
//		Returns object (name:, size:, date:, attr:) or false if the line is not
//		recognized as a directory entry.
//		Supported formats (as descibed in ftpparse by D. J. Bernstein, djb@cr.yp.to):
//		EPLF ,UNIX ls (with or wo gid), Microsoft FTP Service, Windows NT FTP Server,
//		VMS, WFTPD, NetPresenz (Mac), NetWare, MSDOS.
//		Definitely not covered:
//		Long VMS filenames, with information split across two lines.
//		NCSA Telnet FTP (prior to version 2.1 ?

		if line.Size() < 3
			return false
		dirEntry = Object(name: "", size: -1, date: Date("01/01/1970"), attr: 0)
		// So .size defaults to -1  and .date to 01/01/1970 if not supplied
		if line[0] is '+'		// EPLF
			return .eplf(line, dirEntry)
		else if "bcdlps-".Has?(line[0])	// UNIX-style etc
			return .unix(line, dirEntry)
		else if line =~ "^[01]\d-\d\d-\d\d  [01]\d:\d\d[AP]M"	// MSDOS
			return .mdos(line, dirEntry)
		else if (false isnt part =
			line.Extract("\d.*( [ 123]\d-[A-Z][A-Z][A-Z]-\d\d\d\d \d\d:\d\d(:\d\d)? )"))
			if (false isnt dirEntry.date = Date(part))	// Multinet, non-Multinet Vms
				{
				dirEntry.name = line[.. line.Find("")]
				if dirEntry.name.Suffix?(".DIR")
					{
					dirEntry.attr = FILE_ATTRIBUTE.DIRECTORY
					dirEntry.name = dirEntry.name[.. -4]
					}
				return dirEntry
				}
		return false
		}

	eplf(line, dirEntry)
		{
		if false isnt (dirEntry.name = line.Extract("\x09([\x21-\xff]+)$"))
			{
			if line.Has?(",/,")
				dirEntry.attr = FILE_ATTRIBUTE.DIRECTORY
			if false isnt (part = line.Extract(",s(\d+),"))
				dirEntry.size = Number(part)
			if false isnt (part = line.Extract(",m(\d+),"))
				dirEntry.date = dirEntry.date.Plus(seconds: Number(part))
			return dirEntry
			}
		return false
		}

	unix(line, dirEntry)
		{
		if line[0] is 'l' or false isnt (dirEntry.name = line.Extract(" ([^ ]+) -> "))
			return false
		if false isnt (dirEntry.name = line.Extract("[^ ]+$"))
			{
			dirEntry.attr |= line[0] is 'd' ? FILE_ATTRIBUTE.DIRECTORY : 0
			part = line.Extract(" ([A-Z][a-z][a-z] [ 123]\d ( \d\d\d\d|[012]\d:\d\d)) ")
			if part isnt false
				if (false isnt dirEntry.date = Date(part))
					if (dirEntry.date > Date() and dirEntry.date.MinusMonths(Date()) < 12)
						dirEntry.date = dirEntry.date.Plus(years: -1)
			dirEntry.size = Number(line.Extract(" ([0-9]+) +" $ part $ ".*"))
			return dirEntry
			}
		return false
		}

	mdos(line, dirEntry)
		{
		if (false isnt dirEntry.date = Date(line[.. 17]))
			{
			dirEntry.name = line.Extract("[^ ]+$")
			if (line.Has?("<DIR>"))
				dirEntry.attr = FILE_ATTRIBUTE.DIRECTORY
			else
				dirEntry.size = Number(line.Extract(" \d\d* ").Trim())
			return dirEntry
			}
		return false
		}

	remDir(path = "*.*", files = false, details = false, caseSense = true)
		{
		// The ftp-server LIST command to get a directory listing differs from
		// the Windows FindFirst/FindNextFile functions.
		// So we have to do some real work to make them work the same.
		// Essentialy we extract the pathname and filemask from the path variable,
		// supply the pathname with the LIST command to get the complete directory
		// listing and then filter it ourselves using the filemask.
		.splitAndSetPath(path)
		fmask = .BuildFMask(path, caseSense, .realPath, .realBase, .remCWD)

		.createDataSocket()
		.sendCommand("TYPE A")
		.checkResponse("200")

		command = ("LIST " $ .realPath).RightTrim()
		.doVerbose(">>> " $ (command.Prefix?("PASS") ? "PASS ******" : command))
		.scCtrl.Writeline(command)
		.getResponse()
		.checkResponse("150|125")
		dir = Object()
		forever
			{
			try
				{
				if false is resp = .scData.Readline()
					break
				}
			catch
				break
			dir.Add(resp)
			}
		.closeDataSocket()
		return .dirList(fmask, dir, files, details)
		}

	BuildFMask(path, caseSense, realPath = false, realBase = false, remCWD = '/')
		{
		if realPath is false or realBase is false
			{
			ob = .splitPath(path, remCWD)
			realPath = ob.realPath
			realBase = ob.realBase
			}
		fmask = realBase
		// empty mask and '"', '*' or '?' in path returns empty fmask
		if fmask is ""  or false isnt realPath.Extract("[\"*?]+")
			return ''

		if fmask.Size() > 2 and fmask.Suffix?(".*")
			fmask = fmask[.. -2]
		fmask = "(?q)" $ fmask.Replace("?", "(?-q).(?q)").
			Replace("*", "(?-q).*(?q)") $ "(?-q)$"
		return (caseSense ? "" : "(?i)") $ fmask
		}

	dirList(fmask, dir, files, details)
		{
		list = Object()
		for line in dir
			if false isnt (entry = .parseDir(line))
				if String?(entry.name) and entry.name.Size() > 0
					{
					if entry.name !~ fmask
						continue
					if 0 isnt (entry.attr & FILE_ATTRIBUTE.DIRECTORY)
						{
						if files or entry.name is "." or entry.name is ".."
							continue
						entry.name $= "/"
						}
					else if details
						.dirDetails(entry)

					list.Add(details ? entry : entry.name)
					}
		return list
		}
	dirDetails(entry)
		{
		// try to get filedate with MDTM command if not or
		// (probably) only partial returned by LIST command
		if entry.date is entry.date.NoTime() and .supportsMDTM
			{
			.supportsMDTM = false
			if .sendCommand("MDTM " $ .realPath $ entry.name) is "213"
				{
				if .resText.Size() is 18
					{
					date = .resText[-14 ..]
					if date.Numeric?()
						if false isnt (date = Date('#' $ date / 1000000))
							{
							entry.date = date
							.supportsMDTM = true
							}
					}
				if not .supportsMDTM
					throw "FTP: malformed MDTM response: " $ .resText
				}
			}
		// try to get filesize with SIZE command
		// if not returned by LIST command
		if entry.size is -1
			{
			.sendCommand("SIZE " $ .realPath $ entry.name)
			.checkResponse("213")
			entry.size = Number(.resText[4 ..])
			}
		}

	sendCommand(command)
		{
		command = command.RightTrim()
		.doVerbose(">>> " $ (command.Prefix?("PASS") ? "PASS ******" : command))
		.scCtrl.Writeline(command)
		.getResponse()
		return .resCode
		}

	splitAndSetPath(path)
		{
		ob = .splitPath(path, .remCWD)
		.realBase = ob.realBase
		.realPath = ob.realPath
		}

	splitPath(path, remCWD = '/')
		{
		if path.Has?(":")
			throw "FTP: no ':' allowed in path"
		// translate '\' to '/', remove repeating '/' (except at start),
		// dummy spaces, only spaces as filename (at end) equals "."
		path = path.Tr("\\", "/").Replace("(/[^/]+/)[/]+", "\1").Replace(" +/", "/").
			Replace("/ +$", "/.")
		// prefix with CWD if not absolute
		path = (path.Prefix?('/') ? '' : remCWD) $ path
		if path.Size() > 252
			throw "FTP: path too long"

		realBase = Paths.Basename(path).RightTrim()
		realPath = path.BeforeLast('/') $ '/'
		if realPath.Prefix?(remCWD)
			realPath = realPath.AfterFirst(remCWD)
		if false isnt realPath.Extract("[\"*?<>|]+")
			throw "FTP: path contains invalid characters"

		return Object(:realBase, :realPath)
		}

	Close()
		{
		if (.scCtrl isnt false)
			{
			try
				{
				.scCtrl.Writeline("ABORT")
				if (.scData isnt false)
					{
					.scData.Close()
					.scData = false
					}
				.scCtrl.Writeline("QUIT")
				}
			.scCtrl.Close()
			.scCtrl = false
			}
		}
	}
