<?php
// Functions for converting, validating, checking compliance for, and uploading contents of files

// Strip namespaces from root
function strip_namespaces($xml_string) {
  $new_string = preg_replace('/<ead [^>]+>/', '<ead>', $xml_string);
  $new_string = preg_replace('/xlink:type/', 'linktype', $new_string);
  $new_string = preg_replace('/xlink:/', '', $new_string);
  $new_string = preg_replace('/\sxsi:[^"]+"[^"]+"/', '', $new_string);
  return $new_string;
}

// Add DTD declaration, if missing
function add_dtd($xml_string) {
  if (!stristr($xml_string, '<!DOCTYPE')) {
    return str_replace('<ead>','<!DOCTYPE ead PUBLIC "+//ISBN 1-931666-00-8//DTD ead.dtd (Encoded Archival Description (EAD) Version 2002)//EN" "http://archiveswest.orbiscascade.org/ead.dtd">' . "\r\n\r\n" . '<ead>', $xml_string);
  }
  return $xml_string;
}

// Add submission date to publicationstmt
function add_submission_date($xml_string, $time) {
  if ($xml = simplexml_load_string($xml_string)) {
    unset($xml->xpath('//eadheader/filedesc/publicationstmt/date[@type="archiveswest"]')[0][0]);
    $date = $xml->eadheader->filedesc->publicationstmt->addChild('date', date('F j, Y', $time));
    $date->addAttribute('type', 'archiveswest');
    $date->addAttribute('normal', date('Ymd', $time));
    $date->addAttribute('era', 'ce');
    $date->addAttribute('calendar', 'gregorian');
    return $xml->asXML();
  }
  else {
    return false;
  }
}

// Rename c nodes in as2aw conversion
function rename_c($c, $ead) {
  $c_level = 1;
  $parent_node = $c->parentNode;
  if (substr($parent_node->tagName, 0, 2) == 'c0') {
    $parent_level = (int) substr($parent_node->tagName, 2);
    $c_level = $parent_level + 1;
  }
  $new_c = $ead->createElement('c0' . $c_level);
  if ($level = $c->getAttribute('level')) {
    $new_c->setAttribute('level', $level);
  }
  if ($c->hasChildNodes()) {
    for ($i = 0; $i < count($c->childNodes); $i++) {
      $child = $c->childNodes->item($i);
      $new_child = $child->cloneNode(true);
      $new_c->appendChild($new_child);
      if ($new_child->tagName == 'c') {
        rename_c($new_child, $ead);
      }
    }
  }
  $c->parentNode->replaceChild($new_c, $c);
}

// Convert file contents from ArchivesSpace to Archives West EAD
function convert_file($file_contents, $mainagencycode) {
  $errors = array();
  
  // Remove namespaces
  $ead_string = strip_namespaces($file_contents);
  
  // String to DOM
  set_error_handler(function($number, $error){
    if (preg_match('/^DOMDocument::loadXML\(\): (.+)$/', $error, $m) === 1) {
      throw new Exception($m[1]);
      restore_error_handler();
    }
  });
  try {
    $ead = new DOMDocument('1.0', 'utf-8');
    $ead->preserveWhiteSpace = false;
    $ead->formatOutput = true;
    $ead->loadXML($ead_string);
  }
  catch (Exception $e) {
    $errors[] = $e->getMessage();
  }
  restore_error_handler();
  
  if (empty($errors)) {
  
    // Add DOCTYPE with DTD
    if (!stristr($ead_string, '<!DOCTYPE')) {
      $implementation = new DOMImplementation();
      $doctype = $implementation->createDocumentType('ead', '+//ISBN 1-931666-00-8//DTD ead.dtd (Encoded Archival Description (EAD) Version 2002)//EN', 'ead.dtd');
      $ead->insertBefore($doctype, $ead->documentElement);
    }
    
    // Create Xpath
    $xpath = new DOMXpath($ead);
    
    // Check EADID
    if ($eadid = $xpath->query('//eadheader/eadid')->item(0)) {

      // Remove globally unwanted attributes
      foreach (array('audience', 'label', 'id', 'datechar', 'parent') as $attribute) {
        foreach ($xpath->query('//*[@' . $attribute . ']') as $node) {
          $node->removeAttribute($attribute);
        }
      }
      
      // Remove specific unwanted attributes
      $to_remove = array(
        '//eadheader' => 'findaidstatus',
        '//archdesc//physdesc' => 'altrender',
        '//archdesc//physdesc//extent' => 'altrender'
      );
      foreach ($to_remove as $query => $attribute) {
        foreach($xpath->query($query . '[@' . $attribute . ']') as $node) {
          $node->removeAttribute($attribute);
        }
      }
      
      // Rewrite attribute values to lowercase
      $to_convert = array(
        '//container' => 'type',
        '//extref' => 'actuate',
        '//dao' => 'actuate'
      );
      foreach ($to_convert as $query => $attribute) {
        foreach ($xpath->query($query . '[@' . $attribute . ']') as $node) {
          $value = $node->getAttribute($attribute);
          $lower_value = strtolower($value);
          if ($attribute == 'type') {
            $lower_value = preg_replace('/\s+/', '-', $lower_value);
          }
          $node->removeAttribute($attribute);
          if ($lower_value) {
            $node->setAttribute($attribute, $lower_value);
          }
        }
      }
      
      // Add new attributes
      $to_add = array(
        '//eadheader' => array(
          'relatedencoding' => 'dc',
          'scriptencoding' => 'iso15924'
        ),
        '//eadheader/eadid' => array(
          'identifier' => extract_ark($eadid->getAttribute('url')),
          'mainagencycode' => $mainagencycode
        ),
        '//archdesc' => array(
          'relatedencoding' => 'marc21',
          'type' => 'inventory'
        ),
        '//archdesc/did/unitid' => array(
          'countrycode' => $eadid->getAttribute('countrycode'),
          'repositorycode' => $mainagencycode
        ),
        '//archdesc//controlaccess/subject[@source="archiveswest"]' => array(
          'altrender' => 'nodisplay'
        )
      );
      foreach ($to_add as $query => $attributes) {
        foreach ($xpath->query($query) as $node) {
          foreach ($attributes as $name => $value) {
            $node->removeAttribute($name);
            $node->setAttribute($name, $value);
          }
        }
      }
      
      // Dsc: add type attribute and add digits to names of <c> elements
      // Reiterative function to nest <c> is at the bottom of this file
      $levels = array(
          'class' => 'analyticover',
          'collection' => 'analyticover',
          'file' => 'in-depth',
          'fonds' => 'analyticover',
          'item' => 'in-depth',
          'otherlevel' => 'othertype',
          'recordgrp' => 'analyticover',
          'series' => 'analyticover',
          'subgrp' => 'analyticover',
          'subseries' => 'analyticover'
        );
        foreach ($xpath->query('//dsc') as $dsc) {
          if ($dsc->hasChildNodes()) {
            $types = array();
            foreach ($dsc->childNodes as $c) {
              $type = '';
              $level = $c->getAttribute('level');
              if (isset($levels[$level])) {
                $type = $levels[$level];
              }
              if ($type && !in_array($type, $types)) {
                $types[] = $type;
              }
            }
            if (count($types) > 1) {
              $dsc->setAttribute('type', 'combined');
            }
            else {
              if (!empty($types)) {
                $dsc->setAttribute('type', $types[0]);
              }
            }
          }
          else {
            $dsc->parentNode->removeChild($dsc);
          }
        }
        foreach ($xpath->query('//dsc/c') as $c) {
          rename_c($c, $ead);
        }
      
      // Rewrite incorrect role attribute values
      foreach ($xpath->query('//*[@role]') as $node) {
        $role = $node->getAttribute('role');
        if (stristr($role, ' (')) {
          $split_role = explode(' (', $role);
          $role = $split_role[0];
        }
        $lower_role = strtolower($role);
        $node->removeAttribute('role');
        $node->setAttribute('role', $lower_role);
      }
      
      // Change source "naf" to "lcnaf"
      foreach (array('persname', 'corpname') as $type) {
        foreach ($xpath->query('//' . $type) as $node) {
          if ($node->getAttribute('source') == 'naf') {
            $node->removeAttribute('source');
            $node->setAttribute('source', 'lcnaf');
          }
        }
      }
      
      // Address: extract URL only in last line
      if ($last_line = $xpath->query('//eadheader//address/addressline[last()]')->item(0)) {
        if ($extptr = $last_line->getElementsByTagName('extptr')->item(0)) {
          $url = $extptr->getAttribute('href');
          $new_last_line = $ead->createElement('addressline', $url);
          $address = $last_line->parentNode;
          $address->removeChild($last_line);
          $address->appendChild($new_last_line);
        }
      }
      
      // Titlestmt: Reorder children
      foreach ($xpath->query('//eadheader/filedesc/titlestmt') as $titlestmt) {
        $titleproper = $titlestmt->getElementsByTagName('titleproper');
        if (count($titleproper) > 1) {
          $titlestmt->appendChild($titleproper->item(1));
          $titlestmt->appendChild($titleproper->item(0));
        }
        foreach (array('author', 'sponsor') as $child_type) {
          $children = $titlestmt->getElementsByTagName($child_type);
          foreach ($children as $child) {
            $titlestmt->appendChild($child);
          }
        }
        $encodinganalog = $ead->createAttribute('encodinganalog');
        $encodinganalog->value = 'title';
        $titleproper->item(0)->appendChild($encodinganalog);
      }
      
      // Titleproper: remove call number and add altrender
      $titleproper_query = $xpath->query('//titleproper');
      foreach ($titleproper_query as $titleproper) {
        if ($callnumber = $titleproper->getElementsByTagName('num')->item(0)) {
          $titleproper->removeChild($callnumber);
        }
        if ($titleproper->getAttribute('type') == 'filing') {
          $titleproper->setAttribute('altrender', 'nodisplay');
        }
        $titleproper->nodeValue = trim(htmlentities($titleproper->nodeValue));
      }
      
      // First titleproper: copy archdesc/did/unitdate into new date element
      $date_clone = $xpath->query('//archdesc/did/unitdate')->item(0)->cloneNode();
      $titleproper_query->item(0)->appendChild($date_clone);
      $renamed_clone = $ead->createElement('date');
      foreach ($date_clone->attributes as $attribute) {
        $renamed_clone->setAttribute($attribute->nodeName, $attribute->nodeValue);
      }
      while ($date_clone->firstChild) {
        $renamed_clone->appendChild($date_clone->firstChild);
      }
      $date_clone->parentNode->replaceChild($renamed_clone, $date_clone);
      
      // Publicationstmt: Remove <p> and normalize date
      foreach ($xpath->query('//eadheader//publicationstmt') as $publicationstmt) {
        foreach ($xpath->query('//eadheader//publicationstmt/p/date') as $date) {
          $date->setAttribute('encodinganalog', 'date');
          $date->setAttribute('calendar', 'gregorian');
          $date->setAttribute('era', 'ce');
          $year = $date->nodeValue;
          $normalized_year = preg_replace('|\D|', '', $year);
          $date->setAttribute('normal', $normalized_year);
          $publicationstmt->replaceChild($date, $date->parentNode);
        }
      }
      
      // Profiledesc: rewrite date to Y-m-d and add language
      $date = $xpath->query('//eadheader//profiledesc/creation/date')->item(0);
      $date->nodeValue = date('Y-m-d', strtotime($date->nodeValue));
      if ($langusage = $xpath->query('//eadheader/profiledesc/langusage[not(language)]')->item(0)) {
        $language = $ead->createElement('language', $langusage->nodeValue);
        foreach (array('langcode' => 'eng', 'scriptcode' => 'latn', 'encodinganalog' => 'language') as $name => $value) {
          $language_attribute = $ead->createAttribute($name);
          $language_attribute->value = $value;
          $language->appendChild($language_attribute);
        }
        $new_langusage = $ead->createElement('langusage');
        $new_langusage->appendChild($language);
        $profiledesc = $langusage->parentNode;
        $profiledesc->removeChild($langusage);
        $profiledesc->appendChild($profiledesc->getElementsByTagName('creation')->item(0));
        $profiledesc->appendChild($new_langusage);
        $profiledesc->appendChild($profiledesc->getElementsByTagName('descrules')->item(0));
      }
      else if ($language = $xpath->query('//eadheader/profiledesc/langusage/language[not(@encodinganalog)]')->item(0)) {
        $language_attribute = $ead->createAttribute('encodinganalog');
        $language_attribute->value = 'language';
        $language->appendChild($language_attribute);
      }
      
      // Descrules: Change old notice to new notice
      $descrules = $xpath->query('//descrules')->item(0);
      if ($descrules->nodeValue == 'Describing Archives: A Content Standard') {
        $descrules->nodeValue = 'Finding aid based on DACS (Describing Archives: A Content Standard), 2nd Edition.';
      }
      
      // Origination: Collapse multiple nodes into one
      $origination_nodes = $xpath->query('//archdesc/did/origination');
      $origination_count = count($origination_nodes);
      if ($origination_count > 1) {
        $first_node = $origination_nodes->item(0);
        for ($o = 1; $o < $origination_count; $o++) {
          $node = $origination_nodes->item($o);
          foreach ($node->childNodes as $child) {
            $first_node->appendChild($child);
          }
          $node->parentNode->removeChild($node);
        }
      }
      
      // Physdesc: Rewrite extent to lowercase
      foreach ($xpath->query('//physdesc/extent') as $extent) {
        $extent->nodeValue = strtolower($extent->nodeValue);
      }
      
      // Archdesc: Change <dao> to <daogrp>
      foreach ($xpath->query('//archdesc//dao') as $dao) {
        $daogrp = $ead->createElement('daogrp');
        // resource
        $resource = $ead->createElement('resource');
        $resource->setAttribute('label', 'start');
        $daogrp->appendChild($resource);
        // daoloc
        $daoloc = $ead->createElement('daoloc');
        $daoloc->setAttribute('label', 'icon');
        $daoloc->setAttribute('role', 'text/html');
        if ($title = $dao->getAttribute('title')) {
          $daoloc->setAttribute('title', $title);
        }
        if ($href = $dao->getAttribute('href')) {
          $daoloc->setAttribute('href', $href);
        }
        $daogrp->appendChild($daoloc);
        // arc
        $arc = $ead->createElement('arc');
        $arc->setAttribute('from', 'start');
        $arc->setAttribute('to', 'icon');
        if ($show = $dao->getAttribute('show')) {
          $arc->setAttribute('show', $show);
        }
        if ($actuate = $dao->getAttribute('actuate')) {
          $arc->setAttribute('actuate', strtolower($actuate));
        }
        $daogrp->appendChild($arc);
        $dao->parentNode->replaceChild($daogrp, $dao);
      }
      
      // Controlaccess: remove whitespace from headings and separate into nested children
      if ($controlaccess = $xpath->query('//controlaccess')->item(0)) {
        $headings = $controlaccess->childNodes;
        $headings_by_type = array();
        foreach ($headings as $heading) {
          $type = $heading->tagName;
          if ($type == 'subject') {
            if ($heading->getAttribute('source') == 'archiveswest') {
              $type = 'subject-aw';
            }
          }
          $heading->nodeValue = str_replace(' -- ', '--', htmlentities($heading->nodeValue));
          $headings_by_type[$type][] = $heading;
        }
        foreach (array('persname', 'corpname', 'famname', 'geogname', 'subject', 'subject-aw', 'function', 'genreform', 'occupation', 'title') as $type) {
          $new_controlaccess = $ead->createElement('controlaccess');
          if (isset($headings_by_type[$type])) {
            foreach ($headings_by_type[$type] as $heading) {
              $new_controlaccess->appendChild($heading);
            }
            $controlaccess->appendChild($new_controlaccess);
          }
        }
      }
      
      // Add encodinganalog attributes
      $upper = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
      $lower = 'abcdefghijklmnopqrstuvwxyz';
      $encodinganalogs = array(
        '//eadheader/eadid' => 'identifier',
        '//eadheader//author' => 'creator',
        '//eadheader//sponsor' => 'contributor',
        '//eadheader//publisher' => 'publisher',
        '//eadheader//publicationstmt/p/date' => 'date',
        '//archdesc//repository/corpname' => '852$a',
        '//archdesc//repository/subarea' => '852$b',
        '//archdesc//origination/persname' => '100',
        '//archdesc//origination/corpname' => '110',
        '//archdesc//origination/famname' => '100',
        '//archdesc//unittitle' => '245$a',
        '//archdesc//unitdate[@type=\'inclusive\']' => '245$f',
        '//archdesc//abstract' => '5203_',
        '//archdesc//langmaterial[1]/language' => '546',
        '//archdesc//bioghist[substring(translate(head/text(), "' . $upper . '", "' . $lower . '"), 1, 17)="biographical note"]' => '5450_',
        '//archdesc//bioghist[substring(translate(head/text(), "' . $upper . '", "' . $lower . '"), 1, 15)="historical note"]' => '5451_',
        '//archdesc//bioghist[not(@encodinganalog)]' => '545',
        '//archdesc//scopecontent' => '5202_',
        '//archdesc//odd' => '500',
        '//archdesc//arrangement' => '351',
        '//archdesc//altformavail' => '530',
        '//archdesc//accessrestrict' => '506',
        '//archdesc//userestrict' => '540',
        '//archdesc//prefercite' => '524',
        '//archdesc//custodhist' => '561',
        '//archdesc//acqinfo' => '541',
        '//archdesc//accruals' => '584',
        '//archdesc//separatedmaterial' => '5440_',
        '//archdesc//otherfindaid' => '555',
        '//archdesc//relatedmaterial' => '5441_',
        '//archdesc//controlaccess/subject[@source!="archiveswest"]' => '650',
        '//archdesc//controlaccess/subject[@source="archiveswest"]' => '690',
        '//archdesc//controlaccess/persname[not(@role!="")]' => '600',
        '//archdesc//controlaccess/persname[@role!=""]' => '700',
        '//archdesc//controlaccess/corpname[not(@role!="")]' => '610',
        '//archdesc//controlaccess/corpname[@role!=""]' => '710',
        '//archdesc//controlaccess/famname[not(@role!="")]' => '600',
        '//archdesc//controlaccess/famname[@role!=""]' => '700',
        '//archdesc//controlaccess/geogname' => '651',
        '//archdesc//controlaccess/genreform' => '655',
        '//archdesc//controlaccess/occupation' => '656',
        '//archdesc//controlaccess/function' => '657',
        '//unitdate[@type="bulk"]' => '245$g',
        '//physdesc/extent' => '300$a',
        '//archdesc//did/unitid' => '099'
      );
      foreach ($encodinganalogs as $query => $encodinganalog) {
        foreach ($xpath->query($query) as $node) {
          $attribute = $ead->createAttribute('encodinganalog');
          $attribute->value = $encodinganalog;
          $node->appendChild($attribute);
        }
      }
      
      // Convert anchor tags to extref
      foreach ($xpath->query('//a') as $anchor) {
        $extref = $anchor->ownerDocument->createElement('extref');
        foreach ($anchor->childNodes as $child){
          $extref->appendChild($anchor->ownerDocument->importNode($child, true));
        }
        foreach($anchor->attributes as $attrName => $attrNode) {
          $extref->setAttribute($attrNode->nodeName, $attrNode->nodeValue);
        }
        $anchor->parentNode->replaceChild($extref, $anchor);
      }
      
      // Remove <head> elements
      foreach ($xpath->query('//head') as $head) {
        $head->parentNode->removeChild($head);
      }
      
      // Remove empty <prefercite> elements
      foreach ($xpath->query('//prefercite[not(normalize-space())]') as $empty) {
        $empty->parentNode->removeChild($empty);
      }
      
      // DOM to string
      $converted_ead = $ead->saveXML();
    }
    else {
      $errors[] = 'EADID not found.';
    }
  }
  return array('ead'=>$converted_ead, 'errors'=>$errors);
}

// from tools-functions
function extract_ark($string) {
  $ark = '';
  preg_match('|80444\/xv\d{5,6}|', $string, $matches);
  if (isset($matches[0])) {
    $ark = $matches[0];
  }
  return $ark;
}
?>
