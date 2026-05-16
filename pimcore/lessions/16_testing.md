# Chapter 16: Testing

## Overview

Testing in Pimcore follows Symfony's testing philosophy: **PHPUnit** for unit tests, **Symfony's WebTestCase** for functional/controller tests, and integration tests when you need to interact with actual Pimcore objects and the database.

---

## 1. PHPUnit Setup

```bash
composer require --dev phpunit/phpunit symfony/phpunit-bridge
```

**`phpunit.xml.dist` at project root:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<phpunit
    bootstrap="vendor/autoload.php"
    colors="true"
    failOnWarning="true"
    stopOnFailure="false"
    verbose="true"
>
    <testsuites>
        <testsuite name="Application Test Suite">
            <directory>tests</directory>
        </testsuite>
    </testsuites>

    <php>
        <env name="APP_ENV" value="test"/>
        <env name="APP_DEBUG" value="1"/>
        <env name="PIMCORE_ENVIRONMENT" value="test"/>
        <env name="DATABASE_URL" value="mysql://root:@127.0.0.1:3306/pimcore_test"/>
    </php>
</phpunit>
```

---

## 2. Unit Tests

Unit tests are fast — they don't need a database or HTTP client.

### The Service

```php
// src/Service/GreetingService.php
namespace App\Service;

class GreetingService
{
    public function greet(string $name): string
    {
        return "Hello, " . ucfirst($name) . "!";
    }
}
```

### The Test

```php
// tests/Service/GreetingServiceTest.php
namespace App\Tests\Service;

use App\Service\GreetingService;
use PHPUnit\Framework\TestCase;

class GreetingServiceTest extends TestCase
{
    public function testGreetReturnsFormattedMessage(): void
    {
        $service = new GreetingService();
        $result = $service->greet('sandeep');

        $this->assertSame("Hello, Sandeep!", $result);
    }

    public function testGreetHandlesEmptyName(): void
    {
        $service = new GreetingService();
        $result = $service->greet('');

        $this->assertStringStartsWith("Hello, ", $result);
    }
}
```

---

## 3. Functional Tests (Controllers & APIs)

Functional tests spin up the Symfony kernel and make real HTTP requests through the test client.

### The Controller

```php
// src/Controller/Api/ProductController.php
namespace App\Controller\Api;

use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\Routing\Annotation\Route;

class ProductController extends AbstractController
{
    #[Route('/api/test', name: 'api_test', methods: ['GET'])]
    public function test(): JsonResponse
    {
        return $this->json(['message' => 'API works']);
    }
}
```

### The Functional Test

```php
// tests/Controller/ProductControllerTest.php
namespace App\Tests\Controller;

use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;

class ProductControllerTest extends WebTestCase
{
    public function testApiEndpointReturnsJson(): void
    {
        $client = static::createClient();
        $client->request('GET', '/api/test');

        $this->assertResponseIsSuccessful();
        $this->assertResponseHeaderSame('content-type', 'application/json');
        $this->assertJsonStringEqualsJsonString(
            json_encode(['message' => 'API works']),
            $client->getResponse()->getContent()
        );
    }

    public function testProductListEndpoint(): void
    {
        $client = static::createClient();
        $client->request('GET', '/api/products?page=1&limit=5');

        $this->assertResponseIsSuccessful();
        $data = json_decode($client->getResponse()->getContent(), true);

        $this->assertArrayHasKey('products', $data);
        $this->assertArrayHasKey('count', $data);
        $this->assertArrayHasKey('page', $data);
    }

    public function testProductNotFound(): void
    {
        $client = static::createClient();
        $client->request('GET', '/product/99999999');

        $this->assertResponseStatusCodeSame(404);
    }
}
```

---

## 4. Running Tests

```bash
# Run all tests
bin/phpunit

# Run a specific test file
bin/phpunit tests/Service/GreetingServiceTest.php

# Run a specific test method
bin/phpunit --filter testGreetReturnsFormattedMessage

# Run with coverage report (requires Xdebug or PCOV)
XDEBUG_MODE=coverage bin/phpunit --coverage-html var/coverage

# Verbose output
bin/phpunit --testdox
```

Open coverage report:
```bash
xdg-open var/coverage/index.html
```

---

## 5. Testing Commands

```php
// tests/Command/ImportProductsCommandTest.php
namespace App\Tests\Command;

use Symfony\Bundle\FrameworkBundle\Console\Application;
use Symfony\Bundle\FrameworkBundle\Test\KernelTestCase;
use Symfony\Component\Console\Tester\CommandTester;

class ImportProductsCommandTest extends KernelTestCase
{
    public function testCommandFailsWithMissingCsv(): void
    {
        $kernel = static::createKernel();
        $application = new Application($kernel);

        $command = $application->find('app:import-products');
        $tester = new CommandTester($command);

        $tester->execute([]);

        $this->assertStringContainsString('CSV not found', $tester->getDisplay());
    }
}
```

---

## 6. CI/CD Integration (GitHub Actions)

**.github/workflows/tests.yml**

```yaml
name: Pimcore Tests

on: [push, pull_request]

jobs:
  tests:
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8
        env:
          MYSQL_DATABASE: pimcore_test
          MYSQL_USER: root
          MYSQL_ALLOW_EMPTY_PASSWORD: yes
        ports:
          - 3306:3306
        options: >-
          --health-cmd="mysqladmin ping --silent"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=3

      redis:
        image: redis:7
        ports:
          - 6379:6379

    steps:
      - uses: actions/checkout@v4

      - name: Setup PHP
        uses: shivammathur/setup-php@v2
        with:
          php-version: '8.3'
          extensions: mbstring, intl, pdo_mysql, redis
          coverage: xdebug

      - name: Install Composer dependencies
        run: composer install --prefer-dist --no-progress

      - name: Setup .env.test
        run: |
          cp .env .env.test.local
          echo "DATABASE_URL=mysql://root:@127.0.0.1:3306/pimcore_test" >> .env.test.local
          echo "APP_ENV=test" >> .env.test.local

      - name: Setup database
        run: |
          bin/console doctrine:database:create --env=test --if-not-exists
          bin/console doctrine:migrations:migrate --env=test --no-interaction

      - name: Run tests
        run: bin/phpunit --testdox

      - name: Upload coverage
        if: always()
        uses: codecov/codecov-action@v4
```

---

## 7. Testing Best Practices

| Type | When to Use | Speed |
|---|---|---|
| **Unit tests** | Services, helper classes, business logic (no DB) | Fast (milliseconds) |
| **Functional tests** | Controllers, API responses, routes | Medium (seconds) |
| **Integration tests** | Full Pimcore stack: create objects, fetch via API | Slow (requires DB) |

### Rules

- Write **unit tests for services** — fast, no DB dependency
- Write **functional tests for controllers/APIs** — via Symfony's test client
- Keep **Pimcore-specific code** (DataObjects, Assets) in integration tests — these need DB + Pimcore bootstrap
- Use a **separate test database** (`pimcore_test`) and reset before test runs
- Run tests in **CI** (GitHub Actions, GitLab CI, Jenkins)
- Mock external dependencies (HTTP calls, email, queues) in unit/functional tests

### Useful Assertions

```php
$this->assertResponseIsSuccessful();          // 2xx status
$this->assertResponseStatusCodeSame(404);     // exact status code
$this->assertResponseHeaderSame('content-type', 'application/json');
$this->assertJsonStringEqualsJsonString($expected, $actual);
$this->assertStringContainsString('text', $response->getContent());
$this->assertCount(5, $data['products']);
```
