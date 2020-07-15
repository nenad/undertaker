<?php

namespace Undertaker\Dummy\Model;

use Ramsey\Uuid\Uuid;
use Ramsey\Uuid\UuidInterface;

class House extends Building
{
    /**
     * @var \Ramsey\Uuid\UuidInterface
     */
    private $id;
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
        $this->id = Uuid::uuid4();
        $this->rooms = $rooms;
        $this->exteriorType = $exteriorType;
    }

    public function getId(): UuidInterface
    {
        return $this->id;
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
