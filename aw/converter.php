<?php
include('file-functions.php');

function call_convert($str){
  $converted = convert_file($str, "US-ORU");
  return $converted;
}

$arr = call_convert($argv[1]);
echo $arr['ead'];
?>
