// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	InvalidChars: '\\/*?":<>|\t\n\r\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c' $
		'\x0d\x0e\x0f\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x7f'
	InvalidCharsDisplay: `newlines, tabs, control characters or any of the following: ` $
		`/*\?":<>|`
	ReservedNames: #(CON, PRN, AUX, NUL, COM1, COM2, COM3, COM4, COM5,
		COM6, COM7, COM8, COM9, LPT1, LPT2, LPT3, LPT4, LPT5, LPT6, LPT7, LPT8, LPT9)

	Basename(path)
		{
		if .windows?()
			// remove volume e.g. c: or \\server\share
			path = path.Replace(`\A(.:|[/\\][/\\].+?[/\\][^/\\]+)`)
		pos = path.FindLast1of(`\/`)
		return pos is false ? path : path[pos + 1 ..]
		}

	windows?() // Extracted so we can safely override it in the tests
		{
		return not Sys.Browser?() and Sys.Windows?()
		}

	ToStd(path)
		{
		return path.Tr('\\', '/')
		}

	ToWindows(path)
		{
		return path.Tr('/', '\\')
		}

	// To Unix is the same as the Std path
	// keeping this method so operating system dependent code can be kept clearer
	ToUnix(path)
		{
		.ToStd(path)
		}

	ToLocal(path)
		{
		return .windows?() ? .ToWindows(path) : .ToUnix(path)
		}

	ParentOf(path)
		{
		if false is pos = path.FindLast1of("\\/:")
			return '.'
		if path[pos] is ':'
			++pos
		return path[.. pos]
		}

	Combine(@paths)
		{
		path = paths[0]
		for (i = 1; i < paths.Size(); i++)
			path = .combine(path, paths[i])
		return path
		}

	combine(base, path)
		{
		return base.RightTrim(`\/`) $ '/' $ path.LeftTrim(`\/`)
		}

	//TODO handle multiple . or ..
	ToAbsolute(currentPath, relativePath)
		{
		if not relativePath.Prefix?('.')
			return .ToLocal(relativePath)
		if relativePath.Prefix?('..')
			currentPath = .ParentOf(currentPath)
		return .ToLocal(.Combine(currentPath, relativePath.Replace('^\.\.?[\\]?[/]?')))
		}

	IsValid?(path)
		{
		return path isnt '' and path.Size() < MAX_PATH and path.Tr(' -~') is ''
		}

	ValidFileName?(fileName)
		{
		if fileName is ''
			return false
		return not fileName.Has1of?(.InvalidChars)
		}

	/* NOTE: UNC paths, Windows VS Linux/Unix
	Windows UNC paths are structured as follows: \\<server>\<share>\<file/path>
	In order to be considered a complete/valid UNC path, the <server> and
	<share> are required. <file/path> is optional and can include multiple
	directory levels. The slash direction is inconsequential to Windows.

	Conversely, Linux/Unix systems do not use UNC paths.
	Instead, it requires mounting network resources, via the "mount" command.
	Once mounted, these "shares" appear as normal file paths.
		IE: /<ip address>/share => mount to /mnt/share_a, is accessed via: /mnt/share_a
	*/
	ParseUNC(path)
		{
		path = .ToStd(path)
		if not path.Prefix?(`//`)
			return false 					// Not a UNC windows path

		path = path.RemovePrefix(`//`) 		// Remove leading slashes
		if '' is server = path.BeforeFirst(`/`)
			return false

		path = path.AfterFirst(`/`) 		// Remove server name
		if '' is share = path.BeforeFirst(`/`)
			return false

		file = path.AfterFirst(`/`) 		// Optional file name/path component
		return Object(:server, :share, :file)
		}

	/* NOTE: Comparing paths, Windows VS Linux/Unix
	When comparing paths, we need to take into consideration the operating system.
	Otherwise, we cannot accurately determine if the paths are truly equal.
	For reference:
		Windows paths are not case sensitive (meaning \\path and \\PATH are equivalent)
		Linux/Unix paths are case sensitive (meaning \\path and \\PATH are not equivalent)
	*/
	Equal?(path1, path2)
		{
		return .windows?()
			? path1.Lower() is path2.Lower()
			: path1 is path2
		}

	/* NOTE: Linux/Unix and trailing slashes
	In Linux, the standard is to not include a trailing slash when referencing
	directories. A path with a trailing slash can behave differently then a
	matching path without a trailing slash.
	Example:
		mv testFile /home/myuser/existingFolder
			If "existingFolder" does not exist in the example above, "mv" will just
			rename testFile and move it to /home/myuser.

		mv testFile /home/myuser/existingFolder/
			Adding the trailing slash, will check if existingFolder really exists in
			/home/myuser. If this is not the case, it will throw an error and no further
			actions will happen.

	As a result, and for consistency, Equivalent? does not trim the trailing slash
	before carrying out the comparison.
	*/
	Equivalent?(path1, path2)
		{
		return .Equal?(.ToLocal(path1), .ToLocal(path2))
		}

	/* NOTE: Windows path prefixes
	As windows allows for a mixture of slashes to be used in the path, we need to
	be flexible enough to identify matching prefixes given inconsistent slashes.

	While most of our paths should be consistent, we do have areas of code which do
	not format the paths.
		IE: Paths.Combine(...) VS Paths.ToLocal(Paths.Combine(...))

	As a result, Prefix? is designed to be flexible enough to handle these inconsistencies
	*/
	Prefix?(path, prefix)
		{
		return .windows?()
			? .ToStd(path.Lower()).Prefix?(.ToStd(prefix.Lower()))
			: path.Prefix?(prefix)
		}

	/*
	This method is temporary, and should be removed once all paths are built using
	Paths.Combine instead of string concatenation. Reference suggestion: 32860
	*/
	EnsureTrailingSlash(path, noSlashWhenEmpty = false)
		{
		if path is '' and noSlashWhenEmpty
			return ''
		return path.RightTrim(`\/`) $ `/`
		}
	}
