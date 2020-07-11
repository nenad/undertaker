<?php

namespace Undertaker\Dummy\Model;

abstract class Building implements InhabitableInterface
{
    abstract function rooms(): int;
}
