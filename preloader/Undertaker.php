<?php declare(strict_types=1);

use Composer\Autoload\ClassLoader;

class Undertaker
{
    /**
     * Preloads all classes/interfaces/traits found in a directory by running `require_once` on them.
     *
     * @param string $autoloadFile Path to the generated composer autoload file
     * @param string $src Directory where to look for .php files
     * @return string[] List of all preloaded objects
     */
    public static function preload(string $autoloadFile, string $src): array
    {
        /** @var ClassLoader $loader */
        $loader = include $autoloadFile;

        $dir = new RecursiveDirectoryIterator($src);
        $iter = new RecursiveIteratorIterator($dir);
        $files = new RegexIterator($iter, '/^.+\.php$/', RecursiveRegexIterator::GET_MATCH);

        $allClasses = [];
        foreach ($files as $file) {
            try {
                $allClasses[] = self::extractFQCN($file[0]);
            } catch (RuntimeException $e) {
                printf("Error while extracting FQCN: %s\n", $e->getMessage());
            }
        }

        foreach ($allClasses as $class) {
            try {
                require_once $loader->findFile($class);
            } catch (\Throwable $e) {
                printf("Error while loading class %s: %s\n", $class, $e->getMessage());
            }
        }

        return $allClasses;
    }

    /**
     * @param string $filename Path to a .php file
     * @return string Fully qualified class name
     * @throws RuntimeException
     */
    public static function extractFQCN(string $filename): string
    {
        $src = file_get_contents($filename);
        $matches = [];
        $res = preg_match('/^namespace\s+([a-z0-9A-Z\\\]+);/m', $src, $matches);

        $parts = [];
        if ($res) {
            $parts[] = $matches[1];
        }

        $res = preg_match('/^(abstract\sclass|final\sclass|class|trait|interface)\s+([a-zA-Z0-9]+)/m', $src, $matches);
        if (!$res) {
            throw new RuntimeException(sprintf('could not find file type for: %s', $filename));
        }
        $parts[] = $matches[2];

        return implode("\\", $parts);
    }
}
