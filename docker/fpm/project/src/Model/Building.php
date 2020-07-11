<?php

namespace Undertaker\Dummy\Model;

abstract class Building implements InhabitableInterface
{
    public abstract function rooms(): int;
}
