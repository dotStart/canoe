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
package internal

import (
	"encoding/binary"
	"fmt"
	"github.com/dotstart/canoe/internal/metadata"
	"github.com/golang/protobuf/proto"
	"io"
	"os"
)

const expectedMagicNumber = 0xBADC0FEE
const footerSize = 4 + 2 // Magic number + Length

var byteOrder = binary.BigEndian

func ReadExecutableFooter(target string) (*metadata.ApplicationContainer, error) {
	f, err := os.OpenFile(target, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open target: %w", err)
	}

	// size field is 16-bit
	if _, err := f.Seek(-footerSize, 2); err != nil {
		return nil, fmt.Errorf("failed to seek to wrapper footer: %w", err)
	}

	var magicNumber uint32
	if err := binary.Read(f, byteOrder, &magicNumber); err != nil {
		return nil, fmt.Errorf("failed to decode magic number: %w", err)
	}

	if magicNumber != expectedMagicNumber {
		return nil, fmt.Errorf("magic number mismatch")
	}

	var length uint16
	if err := binary.Read(f, byteOrder, &length); err != nil {
		return nil, fmt.Errorf("failed to decode footer size: %w", err)
	}

	if _, err := f.Seek(-(footerSize + int64(length)), 2); err != nil {
		return nil, fmt.Errorf("failed to seek to wrapper configuration: %w", err)
	}

	heap := make([]byte, length)
	if _, err := f.Read(heap); err != nil {
		return nil, fmt.Errorf("failed to read container metadata: %w", err)
	}

	var meta metadata.ApplicationContainer
	if err := proto.Unmarshal(heap, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal container metadata: %w", err)
	}

	return &meta, nil
}

func WriteExecutableFooter(writer io.Writer, meta *metadata.ApplicationContainer) (int, error) {
	length := 0

	encoded, err := proto.Marshal(meta)
	if err != nil {
		return length, fmt.Errorf("failed to encode container metadata: %w", err)
	}

	if _, err := writer.Write(encoded); err != nil {
		return length, err
	}
	length += len(encoded)

	if err := binary.Write(writer, byteOrder, uint32(expectedMagicNumber)); err != nil {
		return length, err
	}

	return length, binary.Write(writer, byteOrder, uint16(length))
}
