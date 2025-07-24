// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
class
	{
	exportKey: "suieo001\n"
	Export(title, data, type, filter, hwnd, name = '')
		{
		DoWithSaveFileName(:filter, :hwnd, :title, file: name, alert: 'Unable to export')
			{ |file|
			data.import_export_type = type
			data.import_export_expiry_date =
				Date().NoTime().Plus(days: 8) /*= after 7 days */
			text = .toText(data)
			.doWithFile(file, "w")
				{|f|
				f.Write(.exportKey)
				f.Write(text)
				}
			}
		}

	toText(data)
		{
		return Zlib.Compress(Pack(data))
		}

	Import(title, type, filter, hwnd)
		{
		file = .openFileName(:filter, :hwnd, :title)
		if file is ""
			return false
		try
			{
			.doWithFile(file, 'r')
				{ |f|
				if .exportKey isnt f.Read(.exportKey.Size())
					{
					.displayMessage(title, 'Unable to import. Invalid file.')
					return false
					}
				text = f.Read()
				}
			x = .fromText(text)
			}
		catch
			{
			.displayMessage(title, "Unable to import from:\n\n\t" $ file)
			return false
			}

		return .valid?(title, type, file, x) ? x : false
		}


	valid?(title, type, file, x)
		{
		expiryDate = x.GetDefault(#import_export_expiry_date, Date().NoTime())
		if expiryDate <= Date().NoTime()
			{
			.displayMessage(title, 'Import not allowed, "' $
				Paths.ToStd(file).AfterLast(`/`) $
				'" has expired (' $ expiryDate.ShortDate() $ ').')
			return false
			}

		if x.GetDefault(#import_export_type, type) isnt type
			{
			.displayMessage(title,
				'Sorry, this is not a valid ' $ type $ ' file.\n\n' $ file)
			return false
			}

		return true
		}

	doWithFile(file, mode, block)
		{
		File(file, mode, block)
		}

	fromText(s)
		{
		return Unpack(Zlib.Uncompress(s))
		}

	displayMessage(title, message)
		{
		Alert(message, title, flags: MB.ICONWARNING)
		}

	openFileName(filter, hwnd, title)
		{
		return OpenFileName(:filter, :hwnd, :title)
		}
	}
