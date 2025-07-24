// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
/*
e.g.
	b = XmlBuilder()
	b.p(align: 'center')
		{
		.S('hello world')
		.span(style: 'margin: 10px;') { 'bye' }
		.img(src: 'x.jpg')
		.S('end')
		}
	b
=> <p align="center">hello world<span style="margin: 10px;">bye</span>
	<img src="x.jpg" />end</p>

indent is the number of spaces for each indent level.
Specifying an indent will also insert newlines.
Margin is the number of spaces for the initial margin/indent.
*/
class
	{
	New(.indent = 0, margin = 0)
		{
		.margin = ' '.Repeat(margin)
		}
	s: ""
	Default(@args)
		{
		tag = args[0]
		if tag[0].Upper?() and tag.Has?('_')
			tag = tag.AfterLast('_') // inside class methods
		.newline()
		.s $= .margin $ '<' $ tag
		for m in args.Members(named:)
			if m isnt #block and m isnt '_'
				.s $= ' ' $ m $ '="' $ XmlEntityEncode(args[m]) $ '"'
		if args.Member?(#block)
			{
			.s $= '>'
			i = .s.Size()
			.newline()
			.margin $= ' '.Repeat(.indent) // increase indent
			result = .Eval2(args.block)
			.margin = .margin[.indent ..] // decrease indent
			if result.Size() is 1 and result[0] isnt this
				.s = .s[.. i] $ XmlEntityEncode(result[0])
			else
				{
				.newline()
				.s $= .margin
				}
			.s $= '</' $ tag $ '>'
			}
		else if args.Member?(1)
			.s $= '>' $ XmlEntityEncode(args[1]) $ '</' $ tag $ '>'
		else
			.s $= ' />'
		.newline()
		return this
		}
	S(s)
		{
		.s $= .margin $ XmlEntityEncode(s)
		return this
		}
	Instruct(version = "1.0", encoding = "us-ascii")
		{
		.newline()
		.s $= .margin $ '<?xml' $
			' version="' $ version $ '"' $
			' encoding="' $ encoding $ '"' $
			'?>'
		.newline()
		return this
		}
	Declare(@args)
		{
		if args.Empty?()
			args = ['DOCTYPE html PUBLIC',
				"-//W3C//DTD XHTML 1.0 Transitional//EN"
				"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"]
		.newline()
		.s $= .margin $ '<!' $ args.PopFirst()
		args.Each
			{ .s $= ' "' $ it $ '"' }
		.s $= '>'
		.newline()
		return this
		}
	Comment(s)
		{
		.newline()
		.s $= .margin $ '<!-- ' $ s $ ' -->'
		.newline()
		return this
		}
	newline()
		{
		if .indent isnt 0 and .s isnt '' and not .s.Suffix?('\n')
			.s $= '\n'
		}
	ToString()
		{ return .s }
	}