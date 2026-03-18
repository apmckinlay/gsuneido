// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// BuiltDate > 20260303
Test
	{
	Test_EncryptHead_returns_unchanged_when_encrypt_false()
		{
		merger = PdfMergerEncrypt()
		obj = Object(head: "1 0 obj <</Type /Page>>")
		result = merger.EncryptHead(obj, false)
		Assert(result.head is: "1 0 obj <</Type /Page>>")
		Assert(result is: obj)
		}

	Test_EncryptHead_returns_same_when_no_hex_or_string()
		{
		encrypt = PdfMergerEncrypt().Setup("user", "owner")
		obj = Object(head: "1 0 obj <</Type /Page>>")
		result = PdfMergerEncrypt().EncryptHead(obj, encrypt)
		Assert(result.head is: "1 0 obj <</Type /Page>>")
		}

	Test_EncryptHead_encrypts_hex_string()
		{
		encrypt = PdfMergerEncrypt().Setup("user", "owner")
		obj = Object(head: "1 0 obj <</Length <448B>>>>")
		PdfMergerEncrypt().EncryptHead(obj, encrypt)
		Assert(obj.head matches: `1 0 obj <</Length <[[:xdigit:]]+>>>>`)
		}

	Test_EncryptHead_encrypts_multiple_hex_strings()
		{
		encrypt = PdfMergerEncrypt().Setup("user", "owner")
		obj = Object(head: "1 0 obj <</Length <448B> /Filter <448B>>>>")
		original_head = obj.head
		PdfMergerEncrypt().EncryptHead(obj, encrypt)
		Assert(obj.head isnt: original_head)
		Assert(obj.head matches:
			`1 0 obj <</Length <[[:xdigit:]]+> /Filter <[[:xdigit:]]+>>>>`)
		}

	Test_EncryptHead_encrypts_string_content_in_parentheses()
		{
		encrypt = PdfMergerEncrypt().Setup("user", "owner")
		obj = Object(head: "1 0 obj <</Author (John Doe) /Title (Test)>>")
		PdfMergerEncrypt().EncryptHead(obj, encrypt)
		Assert(obj.head matches:
			`1 0 obj <</Author <[[:xdigit:]]+> /Title <[[:xdigit:]]+>>>`)
		}

	Test_EncryptHead_mixed_hex_and_string_content()
		{
		encrypt = PdfMergerEncrypt().Setup("user", "owner")
		obj = Object(head:
			"1 0 obj <</Length <448B> /Author (Jo\('\)hn) /Title (Test) " $
			"/Test \(Hello\) /Info <ABC1>>>>")
		PdfMergerEncrypt().EncryptHead(obj, encrypt)
		Assert(obj.head matches:
			`1 0 obj <</Length <[[:xdigit:]]+> /Author <[[:xdigit:]]+> ` $
				`/Title <[[:xdigit:]]+> /Test \(Hello\) /Info <[[:xdigit:]]+>>>>`)
		}

	Test_Setup_validates_passwords()
		{
		Assert(PdfMergerEncrypt().Setup(false, "owner") is: false)
		Assert({ PdfMergerEncrypt().Setup(123, "owner") } throws: 'expected a string')
		Assert({ PdfMergerEncrypt().Setup("user", 123) } throws: 'expected a string')
		}

	Test_updateLength()
		{
		obj = Object()
		obj.head = "1 0 obj <</Length 100>>"
		obj.streamSize = 200
		PdfMergerEncrypt.PdfMergerEncrypt_updateLength(obj)
		Assert(obj.head is: "1 0 obj <</Length 200>>")

		obj.head = "1 0 obj <</Length \t   100 0 R /Filter /FlateDecode>>"
		obj.streamSize = 200
		PdfMergerEncrypt.PdfMergerEncrypt_updateLength(obj)
		Assert(obj.head is: "1 0 obj <</Length 200/Filter /FlateDecode>>")

		obj.head = "1 0 obj <</Length 100 /Filter /FlateDecode /Resources <<>>>>"
		obj.streamSize = 200
		PdfMergerEncrypt.PdfMergerEncrypt_updateLength(obj)
		Assert(obj.head is: "1 0 obj <</Length 200/Filter /FlateDecode /Resources <<>>>>")

		obj.head = "1 0 obj <</Filter /FlateDecode /Length 500>>"
		obj.streamSize = 600
		PdfMergerEncrypt.PdfMergerEncrypt_updateLength(obj)
		Assert(obj.head is: "1 0 obj <</Filter /FlateDecode /Length 600>>")

		obj.head = "1 0 obj <</Filter /FlateDecode /Length1 500 /Length 300>>"
		obj.streamSize = 600
		PdfMergerEncrypt.PdfMergerEncrypt_updateLength(obj)
		Assert(obj.head is: "1 0 obj <</Filter /FlateDecode /Length1 500 /Length 600>>")
		}
	}
