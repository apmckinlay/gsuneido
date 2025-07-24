// Copyright :C 2017 Suneido Software Corp. All rights reserved worldwide.
#(
NONE: 							0 // No value type
SZ: 							1 // Unicode nul terminated string
EXPAND_SZ: 						2 // Unicode nul terminated string
								  //:with environment variable references
BINARY: 						3 // Free form binary
DWORD: 							4 // 32-bit number
DWORD_LITTLE_ENDIAN: 			4 // 32-bit number:same as REG_DWORD
DWORD_BIG_ENDIAN: 				5 // 32-bit number
LINK: 							6 // Symbolic Link:unicode
MULTI_SZ: 						7 // Multiple Unicode strings
RESOURCE_LIST: 					8 // Resource list in the resource map
FULL_RESOURCE_DESCRIPTOR: 		9 // Resource list in the hardware description
RESOURCE_REQUIREMENTS_LIST: 	10
QWORD: 							11 // 64-bit number
QWORD_LITTLE_ENDIAN: 			11 // 64-bit number :same as REG_QWORD
)