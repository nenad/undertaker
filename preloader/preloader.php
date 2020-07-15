<?php declare(strict_types=1);

/** @var Undertaker $undertaker */
include __DIR__ . '/Undertaker.php';
Undertaker::preload(getenv('AUTOLOAD_FILE'), getenv('PHP_DIR'));
