// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	InvalidChars: '*?"<>|\t\n\r\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c' $
		'\x0d\x0e\x0f\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x7f'
	InvalidFileChars: '\\/:'
	ReservedNames: #(CON, PRN, AUX, NUL, COM1, COM2, COM3, COM4, COM5,
		COM6, COM7, COM8, COM9, LPT1, LPT2, LPT3, LPT4, LPT5, LPT6, LPT7, LPT8, LPT9)
	BlankNameDisplay: 'File Name cannot be blank or just an extension'
	InvalidCharsDisplay: `File Name may not contain newlines, tabs, control characters ` $
		`or any of the following: /*\?":<>|`
	ReservedNameDisplay: 'File Name cannot be a reserved Windows File Name'
	MaxAllowedFileNameChars: 255 /*Common between linux and windows */
	MaxAllowedCharsMsg: 'File Name is too long'

	ValidWithPath?(path)
		{
		if path.Has1of?(.InvalidChars)
			return false
		return '' is .WithErrorMsg(path, withPath?:)
		}

	Valid?(filename)
		{
		return '' is .WithErrorMsg(filename)
		}

	WithErrorMsg(filename, withPath? = false, isFolder? = false)
		{
		if withPath?
			filename = Paths.Basename(filename)
		if filename.Blank?() or (filename.BeforeFirst('.').Blank?() and not isFolder?)
			return .BlankNameDisplay
		if 	filename.Has1of?(.InvalidChars $ .InvalidFileChars)
			return '' $ .InvalidCharsDisplay
		if .ReservedNames.Has?(filename.BeforeFirst('.').Upper())
			return .ReservedNameDisplay
		if filename.Size() > .MaxAllowedFileNameChars
			return .MaxAllowedCharsMsg
		return ''
		}
	}