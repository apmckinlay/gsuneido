### EncryptControl

A FieldControl that uses string.Xor to "encrypt" the stored value.

The key comes from EncryptControlKey. The stdlib version simply returns "encrypt" but you can override this in your own library.

**Note:** This is <u>not</u> strong security, but at least it means that the data is not stored in clear text in the database and therefore will not be found by simply scanning the database file.

See also: [EncryptFormat](<../../Reports/Reference/EncryptFormat.md>)