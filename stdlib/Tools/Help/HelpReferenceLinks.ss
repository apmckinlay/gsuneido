// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(howtos = #(), crossref = #(), trouble = #(), firsts = #(),
		setup = #(), concepts = '', tips = '', side? = false, skipTitle? = false,
		newUserTitle? = false)
		{
		str = ''
		if not skipTitle?
			str $= '<p><b>' $ (newUserTitle? ? 'New to Axon?' : 'See Also:') $
				'</b></p>\n'
		str $= .AddHowDoILinks(howtos, side?)
		str $= .AddCrossrefLinks(crossref, side?)
		str $= .AddTroubleLinks(trouble, side?)
		str $= .addTips(tips, side?)
		str $= .addFirstsLinks(firsts, side?)
		str $= .addSetupLinks(setup, side?)
		str $= .addConcepts(concepts, side?)
		return str
		}
	AddHowDoILinks(howtos, side?)
		{
		if not Object?(howtos)
			return howtos
		linktags = .linktags(howtos, side?, '<p class="howto">', '</p>')
		howtosidestart = '<div class="howto" onClick="showhide(\'howtolinks\')">' $
			'<p>How Do I...?&nbsp;' $ HelpArrow('down') $ '</p></div>\r\n' $
			'<div id="howtolinks" class="showhide">'
		return .formatOb(howtos, side?, howtosidestart, '<ul class="howto">',
			'</div>', '</ul>', linktags.open, linktags.close)
		}
	formatOb(ob, side?, sidestart, notsidestart, sideend, notsideend,
		linkopentag = '', linkclosetag = '')
		{
		str = ''
		if ob.Empty?()
			return str
		if side?
			str $= sidestart
		else
			str $= notsidestart
		for link in ob
			{
			if not link.Prefix?('<') and not link.Has?('<a href')
				link = '<$' $ link $ '$>'
			if link.Prefix?(linkopentag.BeforeLast('>'))
				str $= link $ '\r\n'
			else
				str $= linkopentag $ link $ linkclosetag $ '\r\n'
			}
		if side?
			str $= sideend
		else
			str $= notsideend
		return str
		}
	AddCrossrefLinks(crossref, side?)
		{
		linktags = .linktags(crossref, side?, '<p class="crossref">', '</p>')
		return .formatOb(crossref, side?, '', '<ul class="crossref">', '', '</ul>',
			linktags.open, linktags.close)
		}
	AddTroubleLinks(trouble, side?)
		{
		if not Object?(trouble)
			return trouble

		linktags = .linktags(trouble, side?, '<p class="troubleshooting">', '</p>')
		troublesidestart = '<div class="troubleshooting" ' $
			'onClick="showhide(\'troublelinks\')">' $
			'<p>Troubleshooting&nbsp;' $ HelpArrow('down') $ '</p></div>\r\n' $
			'<div id="troublelinks" class="showhide">'
		return .formatOb(trouble, side?, troublesidestart,
			'<ul class="troubleshooting">', '</div>', '</ul>',
			linktags.open, linktags.close)
		}
	addFirstsLinks(firsts, side?)
		{
		linktags = .linktags(firsts, side?, '<p class="settingupaxon">', '</p>')
		return .formatOb(firsts, side?, '', '<ul class="settingupaxon">', '', '</ul>',
			linktags.open, linktags.close)
		}
	addSetupLinks(setup, side?)
		{
		linktags = .linktags(setup, side?, '<p class="settingupaxon">', '</p>')
		return .formatOb(setup, side?, '', '<ul class="setup">', '', '</ul>',
			linktags.open, linktags.close)
		}
	addConcepts(concepts, side?)
		{
		if concepts is ''
			return ''
		return side?
			? '<p class="concepts">' $ concepts $ '</p>'
			: ''
		}
	addTips(tips, side?)
		{
		if tips is ''
			return ''
		if not Object?(tips)
			return side?
				? '<p class="tips">' $ tips $ '</p>'
				: '<ul class="tips"><li>' $ tips $ '</li></ul>'

		linktags = .linktags(tips, side?, '<p class="tips">', '</p>')
		return .formatOb(tips, side?, '', '<ul class="tips">', '', '</ul>',
			linktags.open, linktags.close)
		}
	linktags(ob, side?, open, close)
		{
		if not ob.Empty?()
			if ob[0].Has?('<a href')
				return side?
					? Object(:open, :close)
					: Object(open: '<li>', close: '</li>')
		return Object(open: '', close: '')
		}
	}
