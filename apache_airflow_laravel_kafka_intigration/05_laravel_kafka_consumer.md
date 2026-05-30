# Lesson 5: Laravel Kafka Consumer

## Consume Airflow Results

Airflow can publish workflow result events to Kafka. Laravel consumes them and updates application state.

```php
use Junges\Kafka\Facades\Kafka;
use Junges\Kafka\Contracts\ConsumerMessage;
use Junges\Kafka\Contracts\MessageConsumer;

Kafka::consumer(['reports.generated'], 'laravel-report-updater', config('kafka.brokers'))
    ->withHandler(function (ConsumerMessage $message, MessageConsumer $consumer): void {
        $payload = $message->getBody();

        Report::where('id', $payload['report_id'])->update([
            'status' => 'ready',
            'file_path' => $payload['file_path'],
            'generated_at' => $payload['generated_at'],
        ]);
    })
    ->build()
    ->consume();
```

---

## Run as Artisan Command

```php
final class ConsumeReportsCommand extends Command
{
    protected $signature = 'kafka:consume-reports';

    public function handle(): int
    {
        Kafka::consumer(['reports.generated'], 'laravel-report-updater', config('kafka.brokers'))
            ->withHandler(app(ReportGeneratedHandler::class))
            ->build()
            ->consume();

        return self::SUCCESS;
    }
}
```

Run it with Supervisor, systemd, Docker, or Kubernetes.

---

## Idempotency

Kafka consumers may receive duplicates. Store processed event IDs:

```php
ProcessedKafkaMessage::firstOrCreate([
    'event_id' => $payload['event_id'],
    'consumer_group' => 'laravel-report-updater',
]);
```

Make database updates safe to repeat.
