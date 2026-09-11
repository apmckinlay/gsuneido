// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// sig_<Method> members carry the signatures, and how to interpret them:
//		receiver:	methods on a value of that type, e.g. receiver: 'object'
//		class:		statics on that global class, e.g. class: 'Database'
//		neither:	free functions
// A contribution may instead carry sigs: with a literal list for signatures
// that have no defining record
#(
	ExtensionPoints: ((Signatures)),
	Contributions: (
		(TypeSignatures, Signatures, receiver: object, source: Objects),
		(TypeSignatures, Signatures, receiver: string, source: Strings),
		(TypeSignatures, Signatures, receiver: number, source: Numbers),
		(TypeSignatures, Signatures, receiver: date, source: Dates)
		)
	)
