// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"log"
	"math"
)

// VerifyConsistency checks the internal consistency of the Check data structures.
// It validates the following rules:
// 1. cmtdTran values should be newer than the oldest actvTran
// 2. bytable should only include transactions in actvTran or cmtdTran
// 3. oldest should be either MaxInt or the oldest in actvTran
// 4. map keys should match CkTran.start values in both actvTran and cmtdTran
func (ck *Check) VerifyConsistency() []string {
	var errors []string

	// Rule 4: Verify map keys match CkTran.start values
	for mapKey, tran := range ck.actvTran {
		if mapKey != tran.start {
			errors = append(errors, fmt.Sprintf("actvTran map key %d does not match CkTran.start %d", mapKey, tran.start))
		}
	}

	for mapKey, tran := range ck.cmtdTran {
		if mapKey != tran.start {
			errors = append(errors, fmt.Sprintf("cmtdTran map key %d does not match CkTran.start %d", mapKey, tran.start))
		}
	}

	// Find the actual oldest active transaction
	actualOldest := math.MaxInt
	for _, tran := range ck.actvTran {
		if tran.start < actualOldest {
			actualOldest = tran.start
		}
	}

	// Rule 3: Verify oldest field is correct
	if len(ck.actvTran) == 0 {
		// No active transactions, oldest should be MaxInt
		if ck.oldest != math.MaxInt {
			errors = append(errors, fmt.Sprintf("oldest should be MaxInt (%d) when no active transactions, but is %d", math.MaxInt, ck.oldest))
		}
	} else {
		// Active transactions exist, oldest should match the actual oldest
		if ck.oldest != math.MaxInt && ck.oldest != actualOldest {
			errors = append(errors, fmt.Sprintf("oldest is %d but should be %d (actual oldest active transaction)", ck.oldest, actualOldest))
		}
	}

	// Rule 1: Verify cmtdTran values are newer than oldest actvTran
	if actualOldest != math.MaxInt {
		for tranId, tran := range ck.cmtdTran {
			if tran.end <= actualOldest {
				errors = append(errors, fmt.Sprintf("committed transaction %d (end=%d) is not newer than oldest active transaction (start=%d)", tranId, tran.end, actualOldest))
			}
		}
	}

	// Rule 2: Verify bytable only references transactions in actvTran or cmtdTran
	for tableName, tableActions := range ck.bytable {
		for tranId := range tableActions {
			_, inActv := ck.actvTran[tranId]
			_, inCmtd := ck.cmtdTran[tranId]
			if !inActv && !inCmtd {
				errors = append(errors, fmt.Sprintf("bytable[%s] references transaction %d which is not in actvTran or cmtdTran", tableName, tranId))
			}
		}
	}

	for _, err := range errors {
		fmt.Println(err)
	}
	if len(errors) > 0 {
		// print actvTran and cmtdTran start values
		fmt.Println("actvTran", ck.actvTran)
		fmt.Println("cmtdTran", ck.cmtdTran)
		fmt.Println("bytable", ck.bytable["test_table"])
		log.Fatalln("Check consistency errors")
	}
	return errors
}
