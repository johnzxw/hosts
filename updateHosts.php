<?php
$file        = '/etc/hosts';
$url         = 'https://coding.net/api/user/scaffrey/project/hosts/git/blob/master/hosts-files/hosts';
$fileContent = json_decode(file_get_contents($url), true);
$fileContent = $fileContent['data']['file']['data'];
$fileContent = explode("\n", $fileContent);
$search      = '###################*******************';
$key         = array_search($search, $fileContent);
$oldHosts    = array_slice($fileContent, 0, $key);
$data        = implode("\n", $oldHosts);
$data        .= $search . "\n\n";
$data        .= '####' . date('Y-m-d H:i:s', time()) . "\n\n";
$data        .= implode("\n", array_slice($fileContent, $key, count($fileContent)));

file_put_contents($file, $data);
