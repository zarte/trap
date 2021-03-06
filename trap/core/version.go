/*
 * Trap
 * An anti-pryer server for better privacy
 *
 * This file is a part of Trap project
 *
 * Copyright 2016 Rain Lee <raincious@gmail.com>
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

package core

const (
	TRAP_LOGO = "                #####                   " +
		"                    \r\n             ###" +
		"#   ##                                  " +
		"    \r\n           ###       ##         " +
		"                            \r\n        " +
		"####  ###     ##                        " +
		"            \r\n     ####  ####        #" +
		"#                                   \r\n" +
		"  ####    ## ###        ##   ######     " +
		"  ###       ######  \r\n###           ##" +
		"#       ### ##    ##      ###       ##  " +
		" ## \r\n##             ###     ###  ##  " +
		"  ##     ## ##      ##    ##\r\n###     " +
		"        ##    ##    ##    ##     ## ##  " +
		"    ##   ## \r\n ###             ##     " +
		"     #######     ##  ##     ######  \r\n" +
		"  ###             ##          ##  ##    " +
		"#######     ##      \r\n   ###          " +
		"   ##        ##   ##    #######     ##  " +
		"    \r\n    ###  #          ##      ##  " +
		"  ##   ##     ##    ##      \r\n     ###" +
		"##           ##     ##    ##   ##     ##" +
		"    ##      \r\n                        " +
		"                                    \r\n"

	TRAP_NAME        = "Trap"
	TRAP_DESCRIPTION = "An anti-pryer server for better privacy"
	TRAP_VERSION     = "0.0-prototype"
	TRAP_PROJECTURL  = "https://www.3ax.org/work/trap"
	TRAP_SOURCEURL   = "https://github.com/raincious/trap"
	TRAP_LICENSE     = "Apache License, Version 2.0"
	TRAP_LICENSEURL  = "https://www.apache.org/licenses/LICENSE-2.0"
	TRAP_AUTHOR      = "Rain Lee <raincious@gmail.com>"
	TRAP_COPYRIGHT   = "(C) 2016 Rain Lee"
)

const (
	TRAP_COMMAND_BANNDER = "\r\n" +
		TRAP_LOGO +
		"\r\n  %s\r\n\r\n" +
		"  %s\r\n\r\n" +
		"----------------------------------" +
		"--------------------------\r\n" +
		"  Version  |   %s\r\n" +
		"  License  |   %s\r\n" +
		"           |   %s\r\n" +
		"  Website  |   %s\r\n" +
		"  Source   |   %s\r\n" +
		"----------------------------------" +
		"--------------------------\r\n\r\n"
)
