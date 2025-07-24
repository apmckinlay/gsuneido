// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
// NOTE: usage is: msg.Header and msg.Body
class
	{
	New(.message)
		// message is a "raw" RFC822/2822 message
		{
		}
	Getter_Message()
		{
		return .message
		}
	Getter_Header()
		// returns "unfolded" header (split lines combined)
		{
		.split()
		return .Header // once only
		}
	Getter_Body()
		{
		.split()
		return .Body // once only
		}
	split()
		{
		.Header = ""
		lines = .message.Lines().Iter()
		while lines isnt (line = lines.Next()) and line isnt ""
			{
			if line[0] is ' ' or line[0] is '\t'
				.Header = .Header[.. -2]
			.Header $= line $ '\r\n'
			}
		return .Body = lines.Remainder()
		}
	Getter_Fields()
		{
		flds = Object()
		for line in .Header.Lines()
			{
			ob = line.SplitOnFirst(':')
			flds[ob[0]] = ob[1].Trim()
			}
		return .Fields = flds // once only
		}
	Field(field, default_value = "")
		{
		return .Fields.GetDefault(field, default_value)
		}
	// WARNING: Address and DisplayName don't handle lists
	addrpat: "<(.*)>"
	Address(field)
		{
		s = .Field(field)
		t = s.Extract(.addrpat)
		return (t is false ? s : t).Lower()
		}
	DisplayName(field)
		{
		s = .Field(field)
		return s =~ .addrpat ? s.BeforeFirst('<').Trim() : ""
		}
	// TODO: Date(field)
	Getter_MimeHeader()
		{
		mhdr = ""
		lines = .message.Lines().Iter()
		while lines isnt (line = lines.Next()) and line isnt ""
			if line =~ "(?i)^(Mime-|Content-)"
				{
				mhdr $= line $ "\r\n"
				// handle folded, but don't unfold
				while lines.Remainder()[0] =~ "[ \t]"
					mhdr $= lines.Next() $ "\r\n"
				}
		return .MimeHeader = mhdr
		}

	// handles joining continued lines and skipping 100 continue
	ReadHeader(sc)
		{
		do
			{
			hdr = Object()
			while ((line = sc.Readline()) not in ('', false))
				{
				if hdr.Size() > 0 and (line[0] is ' ' or line[0] is '\t')
					hdr[hdr.Size()-1] $= line
				else
					hdr.Add(line)
				}
			} while not hdr.Empty?() and
				'100' is hdr[0].Extract("^HTTP/[\d.]+ (\d\d\d)")
		return hdr
		}

	HeaderValues(header, env, translateHeaderNames = false)
		{
		for line in header
			{
			name = line.BeforeFirst(':')
			if translateHeaderNames
				name = .TranslateHeaderName(name)
			value = line.AfterFirst(':').Trim()
			if env.Member?(name)
				env[name] $= '\n' $ value
			else
				env[name] = value
			}
		}

	TranslateHeaderName(name)
		{
		return name.Lower().Tr('-', '_')
		}

	// utility methods (don't use instance)

	Split(msg)
		{
		// NOTE: does NOT combine continued lines
		// this is preferable when re-sending
		header = ""
		lines = msg.Lines().Iter()
		while lines isnt (line = lines.Next()) and line isnt ""
			header $= line $ '\r\n'
		return Object(:header, body: lines.Remainder())
		}
	}