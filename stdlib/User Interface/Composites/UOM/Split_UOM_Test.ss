// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Split_UOM(10) is: #(value: 10, uom: ''))
		Assert(Split_UOM('10 kg') is: #(value: 10, uom: 'kg'))
		Assert(Split_UOM('10') is: #(value: 10, uom: ''))
		Assert(Split_UOM('kg') is: #(value: '', uom: 'kg'))
		Assert(Split_UOM('10 kgs', strips:) is: #(value: 10, uom: 'kg'))
		Assert(Split_UOM('10 KGS', strips:) is: #(value: 10, uom: 'KG'))
		Assert(Split_UOM('10 kg', strips:) is: #(value: 10, uom: 'kg'))
		}

	Test_main()
		{
		test_uom = Split_UOM(" ")
		Assert(test_uom.value is: "")
		Assert(test_uom.uom is: "")

		test_uom = Split_UOM("")
		Assert(test_uom.value is: "")
		Assert(test_uom.uom is: "")

		test_uom = Split_UOM(false)
		Assert(test_uom.value is: false)
		Assert(test_uom.uom is: "")

		test_uom = Split_UOM(5)
		Assert(test_uom.value is: 5)
		Assert(test_uom.uom is: "")

		test_uom = Split_UOM("10")
		Assert(test_uom.value is: 10)
		Assert(test_uom.uom is: "")

		test_uom = Split_UOM(' miles')
		Assert(test_uom.value is: "")
		Assert(test_uom.uom is: 'miles')

		test_uom = Split_UOM('us gallons')
		Assert(test_uom.value is: '')
		Assert(test_uom.uom is: 'us gallons')

		test_uom = Split_UOM('miles')
		Assert(test_uom.value is: "")
		Assert(test_uom.uom is: 'miles')

		test_uom = Split_UOM("12 lbs")
		Assert(test_uom.value is: 12)
		Assert(test_uom.uom is: "lbs")

		test_uom = Split_UOM("120,000 kms")
		Assert(test_uom.value is: 120000)
		Assert(test_uom.uom is: 'kms')

		test_uom = Split_UOM("120,000")
		Assert(test_uom.value is: 120000)
		Assert(test_uom.uom is: '')

		test_uom = Split_UOM("120,000,000")
		Assert(test_uom.value is: 120000000)
		Assert(test_uom.uom is: "")

		test_uom = Split_UOM("12 lbs", strips:)
		Assert(test_uom.value is: 12)
		Assert(test_uom.uom is: "lb")

		test_uom = Split_UOM("12 l,bs", strips:)
		Assert(test_uom.value is: 12)
		Assert(test_uom.uom is: "l,b")

		test_uom = Split_UOM(",lbs", strips:)
		Assert(test_uom.value is: "")
		Assert(test_uom.uom is: ",lb")

		test_uom = Split_UOM("12,0 lbs", strips:)
		Assert(test_uom.value is: 120)
		Assert(test_uom.uom is: "lb")
		}
	}