/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms
of the under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

package ui_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	. "github.com/pivotal-cf/ism/ui"
)

var _ = Describe("UI", func() {
	var testUI *UI

	BeforeEach(func() {
		testUI = &UI{
			Out: NewBuffer(),
			Err: NewBuffer(),
		}
	})

	Describe("DisplayText", func() {
		It("prints text with templated values to the out buffer", func() {
			testUI.DisplayText("This is a test for the {{.Struct}} struct", map[string]interface{}{"Struct": "UI"})
			Expect(testUI.Out).To(Say("This is a test for the UI struct\n"))
		})
	})

	Describe("DisplayTable", func() {
		It("prints a table with headers", func() {
			testUI.DisplayTable([][]string{
				{"header1", "header2", "header3"},
				{"data1", "mydata2", "data3"},
				{"data4", "data5", "data6"},
			})
			Expect(testUI.Out).To(Say("header1"))
			Expect(testUI.Out).To(Say("header2"))
			Expect(testUI.Out).To(Say("header3"))
			Expect(testUI.Out).To(Say(`data1\s+mydata2\s+data3`))
			Expect(testUI.Out).To(Say(`data4\s+data5\s+data6`))
		})
	})
})
