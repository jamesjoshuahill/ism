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

package ui

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"github.com/lunixbochs/vtclean"
	runewidth "github.com/mattn/go-runewidth"
)

const DefaultTableSpacePadding = 3

type UI struct {
	Out io.Writer
	Err io.Writer
}

func (ui *UI) DisplayText(text string, data ...map[string]interface{}) {
	var keys interface{}
	if len(data) > 0 {
		keys = data[0]
	}

	formattedTemplate := template.Must(template.New("Display Text").Parse(text + "\n"))
	formattedTemplate.Execute(ui.Out, keys)
}

func (ui *UI) DisplayYAML(data interface{}) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprintf(ui.Out, string(bytes))
	return nil
}

func (ui *UI) DisplayError(err error) {
	fmt.Fprintln(ui.Err, err.Error())
}

func (ui *UI) DisplayTable(table [][]string) {
	if len(table) == 0 {
		return
	}

	for i, str := range table[0] {
		style := color.New(color.Bold)
		table[0][i] = style.Sprint(str)
	}

	var columnPadding []int

	rows := len(table)
	columns := len(table[0])
	for col := 0; col < columns; col++ {
		var max int
		for row := 0; row < rows; row++ {
			if strLen := wordSize(table[row][col]); max < strLen {
				max = strLen
			}
		}
		columnPadding = append(columnPadding, max+DefaultTableSpacePadding)
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			data := table[row][col]
			var addedPadding int
			if col+1 != columns {
				addedPadding = columnPadding[col] - wordSize(data)
			}
			fmt.Fprintf(ui.Out, "%s%s", data, strings.Repeat(" ", addedPadding))
		}
		fmt.Fprintf(ui.Out, "\n")
	}
}

func wordSize(str string) int {
	cleanStr := vtclean.Clean(str, false)
	return runewidth.StringWidth(cleanStr)
}
