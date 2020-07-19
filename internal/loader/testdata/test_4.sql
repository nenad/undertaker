INSERT INTO test.__undertaker_test (function, first_seen_at)
VALUES ('Undertaker\Dummy\Model\House::__construct', now()),
       ('Undertaker\Dummy\Model\House::__destruct', now()),
       ('Undertaker\Dummy\Model\Motel::__construct', null),
       ('Undertaker\Dummy\Model\Motel::__destruct', null),
       ('Undertaker\Dummy\Model\Hotel::__construct', null),
       ('Undertaker\Dummy\Model\Hotel::__destruct', null),
       ('Undertaker\Dummy\Model\Building::__construct', null)
