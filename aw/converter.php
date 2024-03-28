<?php
include('file-functions.php');

function call_convert($str){
  $converted = convert_file($str, "US-ORU");
  return $converted;
}

$arr = call_convert($argv[1]);
if (sizeof($arr['errors']) > 0)
  echo "errors: " . implode($arr['errors'], "|");
else
  echo html_entity_decode($arr['ead']);
?>
