<?php declare(strict_types=1);

require_once __DIR__ . '/../Undertaker.php';

class UndertakerTest
{
    /**
     * @var \Composer\Autoload\ClassLoader
     */
    private $loader;
    /**
     * @var string
     */
    private $autoloadFile;

    public function __construct(string $composerAutoload)
    {
        $this->autoloadFile = $composerAutoload;
        $this->loader = include $composerAutoload;
    }

    public function runTests()
    {
        $this->testAllClassesAreLoaded();
        $this->testFQCN();
    }

    public function testFQCN()
    {
        $testCases = [
            [
                'filename' => __DIR__ . '/fqcnFixtures/simple_class.php',
                'fqcn' => 'Simple\HelloWorld',
                'exception' => null,
            ],
            [
                'filename' => __DIR__ . '/fqcnFixtures/abstract_class.php',
                'fqcn' => 'Abstracted\HelloWorld',
                'exception' => null,
            ],
            [
                'filename' => __DIR__ . '/fqcnFixtures/no_namespace.php',
                'fqcn' => 'HelloWorld',
                'exception' => null,
            ],
            [
                'filename' => __DIR__ . '/fqcnFixtures/empty_file.php',
                'fqcn' => null,
                'exception' => sprintf('could not find file type for: %s', __DIR__ . '/fqcnFixtures/empty_file.php'),
            ],
        ];

        foreach ($testCases as $case) {
            try {
                $fqcn = Undertaker::extractFQCN($case['filename']);
                if ($fqcn !== $case['fqcn']) {
                    throw new Exception(sprintf('want "%s", got "%s"', $case['fqcn'], $fqcn));
                }
            } catch (\RuntimeException $e) {
                if ($e->getMessage() !== $case['exception']) {
                    throw new Exception(sprintf('want "%s", got "%s"', $case['exception'], $e->getMessage()), 0, $e);
                }
            }
        }
    }

    public function testAllClassesAreLoaded()
    {
        $expected = [
            'Undertaker\Dummy\Calculator\PriceCalculator',
            'Undertaker\Dummy\Model\House',
            'Undertaker\Dummy\Model\Building',
            'Undertaker\Dummy\Model\InhabitableInterface',
        ];

        $actual = Undertaker::preload($this->autoloadFile, __DIR__ . '/fixture/src');

        if ($expected !== $actual) {
            throw new Exception('test failed, mismatch in expected preloaded classes');
        }
    }
}

$test = new UndertakerTest(__DIR__ . '/fixture/vendor/autoload.php');
$test->runTests();
