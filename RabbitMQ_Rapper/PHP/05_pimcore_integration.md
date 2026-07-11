# Chapter 5: Pimcore Integration

## Pimcore Mental Model

Pimcore is Symfony-based. Use Symfony services, events, console commands, and bundles.

RabbitMQ is useful in Pimcore for:

- product indexing
- asset processing
- ERP sync
- PIM data export
- image transformation
- cache invalidation
- webhook delivery

---

## Pimcore Bundle Structure

```text
src/RabbitMqBundle/
├── RabbitMqBundle.php
├── Resources/
│   └── config/
│       └── services.yaml
├── EventSubscriber/
│   └── ProductChangedSubscriber.php
├── Command/
│   └── ConsumeRabbitMqCommand.php
└── Handler/
    └── ProductUpdatedHandler.php
```

---

## Publish When DataObject Changes

```php
use Pimcore\Event\DataObjectEvent;
use Pimcore\Event\DataObjectEvents;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;

final class ProductChangedSubscriber implements EventSubscriberInterface
{
    public function __construct(private RabbitPublisher $publisher) {}

    public static function getSubscribedEvents(): array
    {
        return [
            DataObjectEvents::POST_UPDATE => 'onPostUpdate',
        ];
    }

    public function onPostUpdate(DataObjectEvent $event): void
    {
        $object = $event->getObject();

        if (! $object instanceof Product) {
            return;
        }

        $this->publisher->publish('product.updated', [
            'product_id' => $object->getId(),
        ], [
            'source' => 'pimcore',
            'class' => Product::class,
        ]);
    }
}
```

---

## Pimcore Handler Example

```php
final class ProductUpdatedHandler implements MessageHandler
{
    public function __construct(private SearchIndexer $indexer) {}

    public function handle(MessageEnvelope $message): void
    {
        $product = Product::getById((int) $message->payload['product_id']);

        if (! $product instanceof Product) {
            throw new RuntimeException('Product not found.');
        }

        $this->indexer->index($product);
    }
}
```

---

## Pimcore CLI Worker

Because Pimcore has Symfony console:

```bash
bin/console app:rabbitmq:consume product.search_index.queue
```

Run it under Supervisor or Docker.

---

## Pimcore Design Advice

- Publish small IDs, not full object dumps.
- Load fresh Pimcore objects inside handlers.
- Do heavy sync/index/export work in workers.
- Use DLQ for failed product/asset sync messages.
- Include object class and ID in metadata.
- Do not block admin save actions with long external API calls.

