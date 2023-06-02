/*
	Copyright 2023 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package expression

import (
	"strings"

	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/encoding"
)

/**
 * Create an Expression from some wat source.
 *
 */
func ExpressionFromWat(d string) ([]*Expression, error) {
	newex := make([]*Expression, 0)
	lines := strings.Split(d, "\n")
	for _, toline := range lines {
		// Strip any comment from the end
		cptr := strings.Index(toline, ";;")
		if cptr != -1 {
			toline = toline[:cptr]
		}
		toline = strings.Trim(toline, encoding.Whitespace)
		if len(toline) > 0 {
			newe := &Expression{}
			err := newe.DecodeWat(toline, nil)
			if err != nil {
				return newex, err
			}
			newex = append(newex, newe)
		}
	}
	return newex, nil
}

/**
 * Add an expression to the start of some code
 *
 */
func AddExpressionStart(exp []*Expression, to string) ([]*Expression, error) {
	newex, err := ExpressionFromWat(to)
	if err != nil {
		return nil, err
	}

	adjustedExpression := make([]*Expression, 0)
	adjustedExpression = append(adjustedExpression, newex...)
	adjustedExpression = append(adjustedExpression, exp...)
	return adjustedExpression, nil
}

/**
 * Add an expression to the end of some code
 *
 */
func AddExpressionEnd(exp []*Expression, to string) ([]*Expression, error) {
	newex, err := ExpressionFromWat(to)
	if err != nil {
		return nil, err
	}

	adjustedExpression := make([]*Expression, 0)
	adjustedExpression = append(adjustedExpression, exp...)
	adjustedExpression = append(adjustedExpression, newex...)
	return adjustedExpression, nil
}

/**
 * Insert an expression after any instructions that need relocation fixup (offset())
 *
 */
func InsertAfterRelocating(exp []*Expression, to string) ([]*Expression, error) {
	newex, err := ExpressionFromWat(to)
	if err != nil {
		return nil, err
	}

	// Now we need to find where to insert the code
	adjustedExpression := make([]*Expression, 0)
	for _, e := range exp {
		adjustedExpression = append(adjustedExpression, e)
		if e.DataOffsetNeedsAdjusting {
			adjustedExpression = append(adjustedExpression, newex...)
		}
	}
	return adjustedExpression, nil
}
