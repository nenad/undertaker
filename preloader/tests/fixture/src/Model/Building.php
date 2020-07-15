<?php

namespace Undertaker\Dummy\Model;

abstract class Building implements InhabitableInterface
{
    abstract public function rooms(): int;
}
