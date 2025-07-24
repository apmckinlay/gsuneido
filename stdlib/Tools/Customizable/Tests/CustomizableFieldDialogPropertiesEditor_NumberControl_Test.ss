// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getDigits()
		{
		getDigitMask = CustomizableFieldDialogPropertiesEditor_NumberControl.
			CustomizableFieldDialogPropertiesEditor_NumberControl_getDigitMask
		Assert(getDigitMask(2) is: "-##")
		Assert(getDigitMask(3) is: "-###")
		Assert(getDigitMask(4) is: "-#,###")
		Assert(getDigitMask(5) is: "-##,###")
		Assert(getDigitMask(9) is: "-###,###,###")
		}

	Test_set()
		{
		mock = Mock()
		mock.CustomizableFieldDialogPropertiesEditor_NumberControl_decimals =
			decimals = Mock()
		mock.CustomizableFieldDialogPropertiesEditor_NumberControl_digits =
			digits = Mock()
		mock.CustomizableFieldDialogPropertiesEditor_NumberControl_noFormat =
			noFormat = Mock()
		mock.CustomizableFieldDialogPropertiesEditor_NumberControl_tooltip =
			tooltip = Mock()
		mock.Eval(CustomizableFieldDialogPropertiesEditor_NumberControl.Set,
			Field_number_custom)
		digits.Verify.Set(9)
		decimals.Verify.Set(2)
		noFormat.Verify.Set(false)
		tooltip.Verify.Set('')

		c = Field_number_custom
			{
			Control_mask: "-###.#"
			Format_mask: "-###.#"
			}
		mock.Eval(CustomizableFieldDialogPropertiesEditor_NumberControl.Set, c)
		digits.Verify.Set(3)
		decimals.Verify.Set(1)
		noFormat.Verify.Times(2).Set(false)

		c = Field_number_custom
			{
			Control_mask: false
			Field_mask: false
			}
		mock.Eval(CustomizableFieldDialogPropertiesEditor_NumberControl.Set, c)
		digits.Verify.Times(2).Set([any:])
		decimals.Verify.Times(2).Set([any:])
		noFormat.Verify.Set(true)
		digits.Verify.SetReadOnly(true)
		decimals.Verify.SetReadOnly(true)
		}
	}