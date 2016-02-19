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
                        . 'index.gzip';

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
                        . 'clientmeta.go';

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
 *
 * !!!!!!!!!!!!!!!!!!!! DO NOT MODIFIY THIS FILE !!!!!!!!!!!!!!!!!!!!
 * !!!!!!!!!!! Regenerate it with ./statusui/generate.php !!!!!!!!!!!
 *
 *
 */

package status

const STATUS_HOME_STATIC_PAGE =
%ZIPPEDDATA%
";

$gzipedContent = '';

foreach (explode("\n", chunk_split(bin2hex($content), 32, "\n")) as $val) {
    if (!$val) {
        continue;
    }

    if ($gzipedContent) {
        $gzipedContent .= " + \n";
    }

    $gzipedContent .= "    \"\\x" . substr(chunk_split($val, 2, "\\x"), 0, -2) . '"';
}

if (!file_put_contents(
    $outputTo,
    str_replace("%ZIPPEDDATA%", $gzipedContent, $template)
)) {
    printf("Can't save data to file %s\r\n", $outputTo);

    exit(1);
}

printf("OK\r\n");

exit(0);