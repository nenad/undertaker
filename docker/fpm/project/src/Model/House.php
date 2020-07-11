<?php

namespace Undertaker\Dummy\Model;

class House extends Building
{
    /**
     * @var int
     */
    private $rooms;
    /**
     * @var string
     */
    private $exteriorType;

    public function __construct(int $rooms, string $exteriorType)
    {
        $this->rooms = $rooms;
        $this->exteriorType = $exteriorType;
    }

    public function rooms(): int
    {
        return $this->rooms;
    }

    public function exteriorType(): string
    {
        return $this->exteriorType;
    }
}
