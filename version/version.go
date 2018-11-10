// Copyright 2018 flowq Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package version

import (
	"encoding/json"
	"fmt"

	"github.com/mgutz/ansi"
)

const Version = "0.1.2"

const Vendor = "flowq.io"

const Server = "flowqlet"

func VersionInfo() map[string]string {

	apiInfo := make(map[string]string)
	apiInfo["version"] = Version
	apiInfo["server"] = Server
	apiInfo["vendor"] = Vendor
	return apiInfo
}

func PrintVersionInfo() {

	data, _ := json.MarshalIndent(VersionInfo(), "", " ")
	fmt.Println(string(data))
}

//PrintBanner print Flowq banner information
func PrintBanner() {

	//print welcome message use ansi/color Green
	phosphorize := ansi.ColorFunc("green+h")

	fmt.Print(phosphorize(` ________  __                           ______  ` + "\n\r"))
	fmt.Print(phosphorize(`/        |/  |                         /      \ ` + "\n\r"))
	fmt.Print(phosphorize(`########/ ## |  ______   __   __   __ /######  |` + "\n\r"))
	fmt.Print(phosphorize(`## |__    ## | /      \ /  | /  | /  |## |  ## |` + "\n\r"))
	fmt.Print(phosphorize(`##    |   ## |/######  |## | ## | ## |## |  ## |` + "\n\r"))
	fmt.Print(phosphorize(`#####/    ## |## |  ## |## | ## | ## |## |_ ## |` + "\n\r"))
	fmt.Print(phosphorize(`## |      ## |## \__## |## \_## \_## |## / \## |` + "\n\r"))
	fmt.Print(phosphorize(`## |      ## |##    ##/ ##   ##   ##/ ## ## ## |` + "\n\r"))
	fmt.Print(phosphorize(`##/       ##/  ######/   #####/####/   ######  |` + "\n\r"))
	fmt.Print(phosphorize(`                                           ###/ ` + "\n\r"))
	fmt.Print(phosphorize(`                              flowqlet ver ` + Version + " \n\r"))

}
