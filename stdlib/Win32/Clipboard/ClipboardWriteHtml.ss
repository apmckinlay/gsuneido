// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
function (html, add? = false)
	{
	s = 'Version:0.9\n' $
		'StartHTML:000000\n' $
		'EndHTML:000001\n' $
		'StartFragment:000002\n' $
		'EndFragment:000003\n' $
		'<html><body>\n' $
		html $ '\n' $
		'</body></html>'
	len = s.Size()
	starthtml = s.Find('<html>')
	s = s.Replace('000000', starthtml.Pad(6))
	s = s.Replace('000001', len.Pad(6))
	s = s.Replace('000002', (starthtml + 13).Pad(6))
	s = s.Replace('000003', (len - 15).Pad(6))
	Assert(s.Size() is: len)
	fmt = RegisterClipboardFormat('HTML Format')
	ClipboardWriteString(s, fmt, add?)
	return s
	}