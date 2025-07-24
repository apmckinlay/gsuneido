// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		mb = MoneyBag()
		Assert(mb.Amounts() is: #())
		mb.Plus(100, 'CAD')
		Assert(mb.Amounts() is: #(CAD: 100))
		mb.Minus(100, 'CAD')
		Assert(mb.Amounts() is: #(CAD: 0))
		mb.Plus(400, 'USD')
		mb.Plus(56.78, 'USD')
		Assert(mb.Amounts() is: #(CAD: 0, USD: 456.78))
		Assert(mb.Currencies() equalsSet: #(CAD, USD))
		}
	Test_Zero?()
		{
		mb = MoneyBag()
		Assert(mb.Zero?())
		mb.Plus(100, 'CAD')
		mb.Plus(100, 'USD')
		Assert(not mb.Zero?())
		mb.Plus(-100, 'CAD')
		Assert(not mb.Zero?())
		mb.Plus(-100, 'USD')
		Assert(mb.Zero?())
		}
	Test_PlusMB()
		{
		mb = MoneyBag(CAD: 100, USD: 200, PSO: 300)
		mb2 = MoneyBag(CAD: 100, USD: 50, YEN: 150)
		expected = MoneyBag(CAD: 200, USD: 250, PSO: 300, YEN: 150)

		mb.PlusMB(mb2)
		Assert(mb is: expected)
		mb.PlusMB(MoneyBag())
		Assert(mb is: expected)

		mb3 = MoneyBag()
		mb3.PlusMB(mb)
		Assert(mb3 is: expected)

		mb4 = MoneyBag()
		mb4.PlusMB(mb)
		Assert(mb4 is: expected)
		}
	Test_RemoveCurrency()
		{
		mb = MoneyBag(CAD: 100, USD: 200, PSO: 300)
		mb.RemoveCurrency('PSO')
		Assert(mb.Amounts() is: #(CAD: 100, USD: 200))
		}
	Test_ToString()
		{
		test = function (s)
			{
			x = s.Eval() // needs to use eval
			Assert(x.ToString() is: s)
			Assert(Display(x) is: s)
			Assert('' $ x is: s)
			}
		test('MoneyBag()')
		test('MoneyBag(USD: 123)')
		test('MoneyBag(CAD: 100, USD: 50, YEN: 150)')
		}
	Test_Display()
		{
		Assert(MoneyBag().Display() is: '')
		Assert(MoneyBag(USD: 123).Display() is: '123.00 USD')
		Assert(MoneyBag(CAD: 100, USD: 50, YEN: 150000).Display()
			is: '100.00 CAD, 50.00 USD, 150,000.00 YEN')
		}
	}
