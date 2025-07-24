// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_fitText()
		{
		mock = Mock(TabCalcs)
		mock.When.ImageWidth(true).Return(6)
		mock.When.ImageWidth(false).Return(0)
		mock.When.fitText([anyArgs:]).CallThrough()
		mock.TabCalcs_trimChar = '.'
		mock.TabCalcs_trimChars = 3
		mock.PaddingSide = 7

		tab = Object(tabName: 'Tab name with text', image:)
		tab.renderWidth = tab.width = tab.tabName.Size() * charWidth = 3
		tab.trimCharSize = charWidth
		tab.baseCharSizes = tab.tabName.Size().Of({ charWidth })
		Assert(mock.fitText(tab, false) is: 'Tab name with text')
		mock.Verify.Never().ImageWidth([anyArgs:])

		tab.renderWidth -= 1
		Assert(mock.fitText(tab, false) is: 'Tab name w...')
		mock.Verify.ImageWidth([anyArgs:])

		tab.renderWidth -= 3
		Assert(mock.fitText(tab, false) is: 'Tab name ...')

		tab.renderWidth -= 6
		Assert(mock.fitText(tab, false) is: 'Tab nam...')

		tab.renderWidth -= 9
		Assert(mock.fitText(tab, false) is: 'Tab ...')

		tab.renderWidth -= 3
		Assert(mock.fitText(tab, false) is: 'Tab...')

		tab.renderWidth -= 3
		Assert(mock.fitText(tab, false) is: 'Ta...')

		tab.renderWidth -= 3
		Assert(mock.fitText(tab, false) is: 'T...')

		tab.renderWidth -= 3
		Assert(mock.fitText(tab, false) is: 'T..')

		tab.renderWidth -= 3
		Assert(mock.fitText(tab, false) is: 'T.')

		tab.renderWidth -= 3
		Assert(mock.fitText(tab, false) is: 'T')

		tab.renderWidth -= 3
		Assert(mock.fitText(tab, false) is: '')
		}
	}
