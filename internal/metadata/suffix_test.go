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
	"testing"
)

const testKiloByte = 1024
const testMegaByte = testKiloByte * 1024
const testGigaByte = testMegaByte * 1024
const testTeraByte = testGigaByte * 1024
const testPetaByte = testTeraByte * 1024
const testExaByte = testPetaByte * 1024

func TestAppendByteSuffix(t *testing.T) {
	none := AppendByteSuffix(testKiloByte - 1)             // none
	kilo0 := AppendByteSuffix(testKiloByte)                // K
	kilo1 := AppendByteSuffix(testMegaByte - testKiloByte) // K
	mega0 := AppendByteSuffix(testMegaByte)                // M
	mega1 := AppendByteSuffix(testGigaByte - testMegaByte) // M
	giga0 := AppendByteSuffix(testGigaByte)                // G
	giga1 := AppendByteSuffix(testTeraByte - testGigaByte) // G
	tera0 := AppendByteSuffix(testTeraByte)                // T
	tera1 := AppendByteSuffix(testPetaByte - testTeraByte) // T
	peta0 := AppendByteSuffix(testPetaByte)                // P
	peta1 := AppendByteSuffix(testExaByte - testPetaByte)  // P
	exa0 := AppendByteSuffix(testExaByte)                  // E

	if none != "1023" {
		t.Errorf("expected 1023 but got %q", none)
	}
	if kilo0 != "1K" {
		t.Errorf("expected 1K but got %q", kilo0)
	}
	if kilo1 != "1023K" {
		t.Errorf("expected 1023K but got %q", kilo1)
	}
	if mega0 != "1M" {
		t.Errorf("expected 1M but got %q", mega0)
	}
	if mega1 != "1023M" {
		t.Errorf("expected 1023M but got %q", mega1)
	}
	if giga0 != "1G" {
		t.Errorf("expected 1G but got %q", giga0)
	}
	if giga1 != "1023G" {
		t.Errorf("expected 1023G but got %q", giga1)
	}
	if tera0 != "1T" {
		t.Errorf("expected 1T but got %q", tera0)
	}
	if tera1 != "1023T" {
		t.Errorf("expected 1023T but got %q", tera1)
	}
	if peta0 != "1P" {
		t.Errorf("expected 1P but got %q", peta0)
	}
	if peta1 != "1023P" {
		t.Errorf("expected 1023P but got %q", peta1)
	}
	if exa0 != "1E" {
		t.Errorf("expected 1E but got %q", exa0)
	}
}

func TestParseByteSuffix(t *testing.T) {
	none, err := ParseByteSuffix(fmt.Sprintf("%d", testKiloByte-1))
	if err != nil {
		t.Errorf("received error for 1023: %s", err)
	}
	kilo0, err := ParseByteSuffix(fmt.Sprintf("%dk", 1))
	if err != nil {
		t.Errorf("received error for 1K: %s", err)
	}
	kilo1, err := ParseByteSuffix(fmt.Sprintf("%dK", 1023))
	if err != nil {
		t.Errorf("received error for 1023K: %s", err)
	}
	mega0, err := ParseByteSuffix(fmt.Sprintf("%dm", 1))
	if err != nil {
		t.Errorf("received error for 1m: %s", err)
	}
	mega1, err := ParseByteSuffix(fmt.Sprintf("%dM", 1023))
	if err != nil {
		t.Errorf("received error for 1023M: %s", err)
	}
	giga0, err := ParseByteSuffix(fmt.Sprintf("%dg", 1))
	if err != nil {
		t.Errorf("received error for 1g: %s", err)
	}
	giga1, err := ParseByteSuffix(fmt.Sprintf("%dG", 1023))
	if err != nil {
		t.Errorf("received error for 1023G: %s", err)
	}
	tera0, err := ParseByteSuffix(fmt.Sprintf("%dt", 1))
	if err != nil {
		t.Errorf("received error for 1t: %s", err)
	}
	tera1, err := ParseByteSuffix(fmt.Sprintf("%dT", 1023))
	if err != nil {
		t.Errorf("received error for 1023T: %s", err)
	}
	peta0, err := ParseByteSuffix(fmt.Sprintf("%dp", 1))
	if err != nil {
		t.Errorf("received error for 1p: %s", err)
	}
	peta1, err := ParseByteSuffix(fmt.Sprintf("%dP", 1023))
	if err != nil {
		t.Errorf("received error for 1023P: %s", err)
	}
	exa0, err := ParseByteSuffix(fmt.Sprintf("%de", 1))
	if err != nil {
		t.Errorf("received error for 1e: %s", err)
	}
	exa1, err := ParseByteSuffix(fmt.Sprintf("%dE", 1))
	if err != nil {
		t.Errorf("received error for 1E: %s", err)
	}

	if none != testKiloByte-1 {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testKiloByte-1, none)
	}
	if kilo0 != testKiloByte {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testKiloByte, kilo0)
	}
	if kilo1 != testMegaByte-testKiloByte {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testMegaByte-testKiloByte, kilo1)
	}
	if mega0 != testMegaByte {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testMegaByte, mega0)
	}
	if mega1 != testGigaByte-testMegaByte {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testGigaByte-testMegaByte, mega1)
	}
	if giga0 != testGigaByte {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testGigaByte, giga0)
	}
	if giga1 != testTeraByte-testGigaByte {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testTeraByte-testGigaByte, giga1)
	}
	if tera0 != testTeraByte {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testTeraByte, tera0)
	}
	if tera1 != testPetaByte-testTeraByte {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testPetaByte-testTeraByte, tera1)
	}
	if peta0 != testPetaByte {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testPetaByte, peta0)
	}
	if peta1 != testExaByte-testPetaByte {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testExaByte-testPetaByte, peta1)
	}
	if exa0 != testExaByte {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testExaByte, exa0)
	}
	if exa1 != testExaByte {
		t.Errorf("expected %[1]d for input %[1]d but got %[2]d", testExaByte, exa1)
	}
}
