/*
Copyright 2014 Alex Crawford

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

package ishmail

import (
	"testing"
	"time"
	"net/mail"
)

func TestHeadersToString(t *testing.T) {
	tests := []struct {
		input  map[string]string
		output string
	}{
		{
			input:  nil,
			output: "",
		},
		{
			input: map[string]string{
				"From":         "test1@example.com",
				"To":           "test2@example.com",
				"Subject":      "Testing",
				"MIME-Version": "1.0",
				"Content-Type": "text/html",
				"Date":         time.Date(2014, time.March, 29, 2, 32, 17, 1232, time.UTC).Format(time.RFC1123Z),
			},
			output: "From: test1@example.com\r\nTo: test2@example.com\r\nSubject: Testing\r\nMIME-Version: 1.0\r\nContent-Type: text/html\r\nDate: Sat, 29 Mar 2014 02:32:17 +0000\r\n",
		},
	}
	for _, test := range tests {
		got := headersToString(test.input)
		if got != test.output {
			t.Errorf("Expected:\n'''%s'''\n\nbut got:\n'''%s'''", test.output, got)
		}
	}
}

func TestAddrsToAddresses(t *testing.T) {
	tests := []struct {
		input  []*mail.Address
		output []string
	}{
		{
			input:  []*mail.Address{},
			output: []string{},
		},
		{
			input: []*mail.Address{
				&mail.Address{
					Name:    "First Last",
					Address: "test1@example.com",
				},
			},
			output: []string{
				"test1@example.com",
			},
		},
		{
			input: []*mail.Address{
				&mail.Address{
					Name:    "First Last",
					Address: "test1@example.com",
				},
				&mail.Address{
					Name:    "A B",
					Address: "test2@example.com",
				},
			},
			output: []string{
				"test1@example.com",
				"test2@example.com",
			},
		},
	}
	for _, test := range tests {
		got := addrsToAddresses(test.input)
		if len(got) != len(test.output) {
			t.Errorf("Expected %d addresses but got %d addresses", len(test.output), len(got))
		}
		for i, result := range got {
			if result != test.output[i] {
				t.Errorf("Expected '%s' but got '%s'", test.output[i], result)
			}
		}
	}
}

func TestAddrsToString(t *testing.T) {
	tests := []struct {
		input  []*mail.Address
		output string
	}{
		{
			input:  []*mail.Address{},
			output: "",
		},
		{
			input: []*mail.Address{
				&mail.Address{
					Name:    "First Last",
					Address: "test1@example.com",
				},
			},
			output: "=?utf-8?q?First_Last?= <test1@example.com>",
		},
		{
			input: []*mail.Address{
				&mail.Address{
					Name:    "First Last",
					Address: "test1@example.com",
				},
				&mail.Address{
					Name:    "A B",
					Address: "test2@example.com",
				},
			},
			output: "=?utf-8?q?First_Last?= <test1@example.com>, =?utf-8?q?A_B?= <test2@example.com>",
		},
	}
	for _, test := range tests {
		got := addrsToString(test.input)
		if got != test.output {
			t.Errorf("Expected '%s' but got '%s'", test.output, got)
		}
	}
}
