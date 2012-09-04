// Copyright 2012 Aaron Jacobs. All Rights Reserved.
// Author: aaronjjacobs@gmail.com (Aaron Jacobs)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/gob"
	"flag"
	"github.com/jacobsa/aes/testing"
	"log"
	"math/rand"
	"os"
)

var function = flag.String("func", "", "Function for which to generate cases.")
var randSrc = rand.New(rand.NewSource(0xdeadbeef))

func randBytes(n uint32) []byte {
	b := make([]byte, n)
	for i, _ := range b {
		b[i] = byte(rand.Intn(256))
	}

	return b
}

func doGenerateSubkey() []testing.GenerateSubkeyTestCase {
	numCases := (1 << 13)
	cases := make([]testing.GenerateSubkeyTestCase, numCases)

	for i, _ := range cases {
		c := &cases[i]
		c.Key = randBytes(16)
		c.K1, c.K2 = generateSubkey(c.Key)
	}

	return cases
}

func doCmac() []testing.CmacTestCase {
	numCases := (1 << 13)
	cases := make([]testing.CmacTestCase, numCases)

	for i, _ := range cases {
		c := &cases[i]
		c.Key = randBytes(16)
		c.Msg = randBytes(uint32(i % 256))
		c.Mac = generateCmac(c.Key, c.Msg)
	}

	return cases
}

func doDbl() []testing.DblTestCase {
	numCases := (1 << 10)
	cases := make([]testing.DblTestCase, numCases)

	for i, _ := range cases {
		c := &cases[i]
		c.Input = randBytes(16)
		c.Output = dbl(c.Input)
	}

	return cases
}

func main() {
	flag.Parse()

	var cases interface{}
	switch *function {
	case "":
		log.Fatalf("You must set -func.")
	case "generateSubkey":
		cases = doGenerateSubkey()
	case "AES-CMAC":
		cases = doCmac()
	case "dbl":
		cases = doDbl()
	default:
		log.Fatalf("Unrecognized function: %s", *function)
	}

	encoder := gob.NewEncoder(os.Stdout)
	encoder.Encode(cases)
}
