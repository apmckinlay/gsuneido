// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.noIndent? = false)
		{
		}
	indent: 0
	s: ''
	DoWithIndent(block)
		{
		if .noIndent? is false
			.indent++
		block()
		if .noIndent? is false
			.indent--
		}

	Add(tag, s, extra = '')
		{
		if s is false
			{
			.s $= '\n' $ '\t'.Repeat(.indent) $ '<' $ tag $ Opt(' ', extra) $ ' />'
			return
			}

		.s $= '\n' $ '\t'.Repeat(.indent) $ '<' $ tag $ Opt(' ', extra) $ '>' $
			s $ '</' $ tag $ '>'
		}

	AddPure(s)
		{
		.s $= s
		}

	AddWithBlock(tag, block, attr = #(), noNewline? = false)
		{
		.AddOpen(tag, attr)
		.DoWithIndent(block)
		.AddClose(tag, noNewline?)
		}

	AddOpen(tag, attr = #())
		{
		attr = attr.Map2({ |m, v| m $ '="' $ v $ '"' }).Join(' ')
		.s $= '\n' $ '\t'.Repeat(.indent) $ '<' $ tag $ Opt(' ', attr) $ '>'
		}

	AddClose(tag, noNewline? = false)
		{
		.s $= (noNewline? ? '' : '\n') $ '\t'.Repeat(.indent) $ '</' $ tag $ '>'
		}

	Suffix?(suffix)
		{
		return .s.Suffix?(suffix)
		}

	Get()
		{
		return Opt(.s.RemovePrefix('\n'), '\n')
		}
	}