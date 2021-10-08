/*
 * Copyright 2021 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package metadata

import (
	"fmt"
	"strconv"
	"unicode"
)

var byteSuffixes = []rune{
	'K',
	'M',
	'G',
	'T',
	'P',
	'E',
	// Z and Y exceed the storage capacity of uint64
}

const byteSuffixBase = uint64(1024)
const byteSuffixOffset = uint64(10)

func AppendByteSuffix(size uint64) string {
	for i := len(byteSuffixes) - 1; i >= 0; i-- {
		divisor := byteSuffixBase << (uint64(i) * byteSuffixOffset)

		if size%divisor == 0 {
			return strconv.FormatUint(size/divisor, 10) + string(byteSuffixes[i])
		}
	}

	return strconv.FormatUint(size, 10)
}

func ParseByteSuffix(input string) (uint64, error) {
	suffix := rune(input[len(input)-1])
	if unicode.IsNumber(suffix) {
		number, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("illegal size: %w", err)
		}

		return number, nil
	} else if !unicode.IsLetter(suffix) {
		return 0, fmt.Errorf("illegal suffix: %c", suffix)
	}

	suffix = unicode.ToUpper(suffix)

	number, err := strconv.ParseUint(input[:len(input)-1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("illegal size: %w", err)
	}

	for i, s := range byteSuffixes {
		if s == suffix {
			multiplier := byteSuffixBase << (uint64(i) * byteSuffixOffset)
			return number * multiplier, nil
		}
	}

	return 0, fmt.Errorf("illegal size suffix: %c", suffix)
}
