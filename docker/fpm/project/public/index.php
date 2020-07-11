<?php

use Undertaker\Dummy\Calculator\PriceCalculator;
use Undertaker\Dummy\Model\House;

require_once __DIR__ . '/../vendor/autoload.php';

$calc = new PriceCalculator();

$b = new House(4, 'concrete');

$calc->allRoomsPrice($b);
