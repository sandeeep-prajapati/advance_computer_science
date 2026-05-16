# Chapter 12: Data Import & Export

## Overview

Pimcore projects frequently need to import data from CSVs, XMLs, or external APIs, and export data for syndication, backups, or reporting. This chapter covers a complete import pipeline with error handling.

---

## CSV Import: Step by Step

### Example CSV Format

```csv
sku,name,price,description,image_url
P001,Phone X,999,Latest smartphone,https://example.com/images/phone.jpg
P002,Laptop Z,1499,High-end laptop,https://example.com/images/laptop.jpg
P003,Wireless Mouse,29,Ergonomic mouse,https://example.com/images/mouse.jpg
```

---

## Full Import Command

**File:** `src/Command/ImportProductsCommand.php`

```php
<?php
namespace App\Command;

use Pimcore\Model\DataObject\Product;
use Pimcore\Model\Asset;
use Pimcore\Model\Asset\Image;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Style\SymfonyStyle;

class ImportProductsCommand extends Command
{
    protected static $defaultName = 'app:import-products';

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $io = new SymfonyStyle($input, $output);

        $csvPath = PIMCORE_PROJECT_ROOT . '/var/import/products.csv';

        if (!file_exists($csvPath)) {
            $io->error("CSV not found: $csvPath");
            return Command::FAILURE;
        }

        $handle = fopen($csvPath, 'r');
        $header = fgetcsv($handle, 1000, ','); // first row = headers
        $count = 0;

        while (($row = fgetcsv($handle, 1000, ',')) !== false) {
            $data = array_combine($header, $row);

            try {
                $sku = $data['sku'] ?? null;
                if (!$sku) {
                    throw new \Exception("Missing SKU in row");
                }

                // Find existing product by SKU or create new
                $product = Product::getBySku($sku, 1);
                if (!$product) {
                    $product = new Product();
                    $product->setKey(\Pimcore\File::getValidFilename($sku));
                    $product->setParentId(1); // root folder
                    $product->setPublished(true);
                    $product->setSku($sku);
                }

                $product->setName($data['name'] ?? '');
                $product->setPrice((float)($data['price'] ?? 0));
                $product->setDescription($data['description'] ?? '');

                // Handle image download
                if (!empty($data['image_url'])) {
                    $asset = $this->downloadImage($data['image_url'], $sku);
                    if ($asset instanceof Image) {
                        $product->setImages([$asset]);
                    }
                }

                $product->save();
                $count++;
                $io->success("Imported/Updated: {$sku}");

            } catch (\Throwable $e) {
                $io->error("Error importing row [{$sku}]: " . $e->getMessage());
            }
        }

        fclose($handle);
        $io->success("Import finished. Total products: $count");

        return Command::SUCCESS;
    }

    private function downloadImage(string $url, string $sku): ?Image
    {
        $contents = @file_get_contents($url);
        if ($contents === false) {
            throw new \Exception("Failed to download image from $url");
        }

        $filename = $sku . '-' . basename(parse_url($url, PHP_URL_PATH));

        // Get or create the /imported folder
        $folder = Asset::getByPath('/imported');
        if (!$folder) {
            $folder = new Asset\Folder();
            $folder->setKey('imported');
            $folder->setParentId(1);
            $folder->save();
        }

        // Return existing if already imported
        $existing = Asset::getByPath('/imported/' . $filename);
        if ($existing instanceof Image) {
            return $existing;
        }

        // Create new asset
        $image = new Image();
        $image->setParent($folder);
        $image->setFilename($filename);
        $image->setData($contents);
        $image->save();

        return $image;
    }
}
```

### Run the Import

```bash
# Place your CSV file
mkdir -p var/import
cp your-products.csv var/import/products.csv

# Run the import command
bin/console app:import-products
```

---

## Handling Categories from CSV

Extend the CSV format:
```csv
sku,name,price,description,category_path,image_url
P001,Phone X,999,Latest smartphone,/categories/electronics,https://...
```

Add to the import:
```php
// Find/create category
if (!empty($data['category_path'])) {
    $category = \Pimcore\Model\DataObject\Category::getByPath($data['category_path']);
    if ($category) {
        $product->setCategory($category);
    }
}
```

---

## Export Products to CSV

**File:** `src/Command/ExportProductsCommand.php`

```php
<?php
namespace App\Command;

use Pimcore\Model\DataObject\Product\Listing;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Style\SymfonyStyle;

class ExportProductsCommand extends Command
{
    protected static $defaultName = 'app:export-products';

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $io = new SymfonyStyle($input, $output);
        $exportPath = PIMCORE_PROJECT_ROOT . '/var/export/products_export.csv';

        // Ensure export directory exists
        if (!is_dir(dirname($exportPath))) {
            mkdir(dirname($exportPath), 0755, true);
        }

        $handle = fopen($exportPath, 'w');
        fputcsv($handle, ['id', 'sku', 'name', 'price', 'description', 'category']);

        $listing = new Listing();
        $listing->setOrderKey('name');
        $listing->setOrder('asc');

        foreach ($listing as $product) {
            fputcsv($handle, [
                $product->getId(),
                $product->getSku(),
                $product->getName(),
                $product->getPrice(),
                strip_tags($product->getDescription()),
                $product->getCategory()?->getName() ?? '',
            ]);
        }

        fclose($handle);

        $io->success("Exported to: $exportPath");
        return Command::SUCCESS;
    }
}
```

```bash
bin/console app:export-products
```

---

## Using League\CSV for Better CSV Handling

For large files or complex CSV requirements:

```bash
composer require league/csv
```

```php
use League\Csv\Reader;
use League\Csv\Writer;

// Reading
$reader = Reader::createFromPath('/var/import/products.csv', 'r');
$reader->setHeaderOffset(0);

foreach ($reader->getRecords() as $record) {
    // $record['sku'], $record['name'], etc.
}

// Writing
$writer = Writer::createFromPath('/var/export/products.csv', 'w');
$writer->insertOne(['sku', 'name', 'price']);
$writer->insertAll($rows);
```

---

## Import via Pimcore Admin UI (Built-in)

Pimcore also has a **CSV Import** feature built into the Admin:
1. Go to **Data Objects** → select a class
2. Click **Import CSV** button in the toolbar
3. Map CSV columns to class fields
4. Preview & execute import

This is useful for non-technical users but has less flexibility than a custom command.

---

## Error Handling Best Practices

| Strategy | Implementation |
|---|---|
| Per-row try/catch | Don't let one bad row stop the entire import |
| Log errors to file | `\Pimcore\Log\Simple::log('import', $message)` |
| Dry-run mode | Add `--dry-run` flag to simulate import without saving |
| SKU deduplication | Always check `getBySku()` before creating new objects |
| Image download retry | Wrap HTTP calls in retry logic for flaky URLs |

### Adding a Dry-Run Flag

```php
protected function configure(): void
{
    $this->addOption('dry-run', null, InputOption::VALUE_NONE, 'Simulate without saving');
}

// In execute():
$dryRun = $input->getOption('dry-run');
if (!$dryRun) {
    $product->save();
}
```

```bash
bin/console app:import-products --dry-run
```
