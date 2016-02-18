#!/usr/bin/php
<?php

/**
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
 *
 */

$currentDir     =   dirname(__FILE__);
$base64File     =   $currentDir
                        . DIRECTORY_SEPARATOR
                        . 'dist'
                        . DIRECTORY_SEPARATOR
                        . 'index.b64';

$outputTo       =   $currentDir
                        . DIRECTORY_SEPARATOR
                        . '..'
                        . DIRECTORY_SEPARATOR
                        . 'trap'
                        . DIRECTORY_SEPARATOR
                        . 'core'
                        . DIRECTORY_SEPARATOR
                        . 'status'
                        . DIRECTORY_SEPARATOR
                        . 'template.go';

if (!file_exists($base64File)) {
    printf("Can't found file %s\r\n", $base64File);

    exit(1);
}

$content = file_get_contents($base64File);

if (!$content) {
    printf("No content in file %s\r\n", $base64File);

    exit(1);
}

$template = "/*
 * Trap
 * An anti-pryer server for better privacy
 *
 * This file is a part of Trap project
 *
 * Copyright 2016 Rain Lee <raincious@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the \"License\");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an \"AS IS\" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *
 * NOTICE:
 *
 * This file may contains third-party components, they are the properties
 * of it's respective holder(s). See README.md for more information.
 *
 */

package status

import (
    \"encoding/base64\"
)

var (
    StaticClientPage, staticClientPageErr =
        base64.StdEncoding.DecodeString(
%Base64Code%)
)";

$base64Code = '';

foreach (explode("\n", chunk_split($content, 64, "\n")) as $val) {
    if (!$val) {
        continue;
    }

    if ($base64Code) {
        $base64Code .= " +\n";
    }

    $base64Code .= '            "' . $val . '"';
}

$base64Code .= "\n";

$base64Code = rtrim($base64Code);

if (!file_put_contents(
    $outputTo,
    str_replace("%Base64Code%", $base64Code, $template)
)) {
    printf("Can't save data to file %s\r\n", $outputTo);

    exit(1);
}

printf("OK\r\n");

exit(0);