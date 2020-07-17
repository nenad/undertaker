<?php

use Nenad\Undertaker\Preloader;

$autoload = __DIR__ . '/../vendor/autoload.php';
include $autoload;

$undertaker = new Preloader($autoload);
$undertaker->load(__DIR__ . '/../src');
