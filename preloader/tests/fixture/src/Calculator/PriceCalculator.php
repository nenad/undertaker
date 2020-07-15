<?php

namespace Undertaker\Dummy\Calculator;

use Undertaker\Dummy\Model\Building;

class PriceCalculator
{
    public function allRoomsPrice(Building $building): float
    {
        return 50.3 * $building->rooms();
    }

    public function total(Building $building): float 
    {
        return $this->allRoomsPrice($building) + 300;
    }
}
