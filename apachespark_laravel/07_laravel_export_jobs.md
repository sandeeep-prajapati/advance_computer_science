# Lesson 7: Laravel Export Jobs

## Export Source Data

Laravel can prepare clean source data before Spark runs.

```php
final class ExportOrdersForSpark implements ShouldQueue
{
    public function handle(): void
    {
        $path = "spark/input/orders-{$this->sparkJobId}.csv";

        Storage::disk('s3')->put($path, $this->buildCsv());

        SparkJob::find($this->sparkJobId)->update([
            'status' => 'exported',
            'parameters->input_path' => $path,
        ]);
    }
}
```

---

## Chunk Large Exports

Never load millions of rows into memory:

```php
Order::query()
    ->whereBetween('created_at', [$from, $to])
    ->orderBy('id')
    ->chunkById(1000, function ($orders) use ($writer) {
        foreach ($orders as $order) {
            $writer->insertOne([
                $order->id,
                $order->customer_id,
                $order->total,
                $order->created_at->toIso8601String(),
            ]);
        }
    });
```

---

## Prefer Typed Data

CSV is okay for learning. For production analytics, prefer:

- Parquet
- Delta Lake
- database read replica through JDBC
- CDC into a data lake

PHP is not usually the best tool for writing huge Parquet datasets, so Laravel often prepares job metadata while Spark reads source systems directly.
