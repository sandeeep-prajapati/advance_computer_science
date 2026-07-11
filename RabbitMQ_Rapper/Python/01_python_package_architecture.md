# Chapter 1: Python Package Architecture

## Package Structure

```text
python-rabbit-app/
├── pyproject.toml
├── messaging/
│   ├── __init__.py
│   ├── envelope.py
│   ├── serializer.py
│   ├── publisher.py
│   ├── consumer.py
│   ├── registry.py
│   ├── topology.py
│   └── exceptions.py
├── handlers/
│   ├── __init__.py
│   └── product_updated.py
└── worker.py
```

Install:

```bash
pip install pika
```

For Django:

```bash
pip install pika django
```

---

## Message Envelope

```python
from dataclasses import dataclass, field
from datetime import datetime, timezone
from typing import Any

@dataclass(frozen=True)
class MessageEnvelope:
    id: str
    type: str
    version: int
    occurred_at: str
    payload: dict[str, Any]
    meta: dict[str, Any] = field(default_factory=dict)

    @staticmethod
    def now_iso() -> str:
        return datetime.now(timezone.utc).isoformat()
```

---

## Handler Protocol

```python
from typing import Protocol

class MessageHandler(Protocol):
    def handle(self, message: MessageEnvelope) -> None:
        ...
```

Using a protocol keeps handlers easy to test.

---

## Serializer

```python
import json

class JsonMessageSerializer:
    def encode(self, message: MessageEnvelope) -> bytes:
        return json.dumps({
            "id": message.id,
            "type": message.type,
            "version": message.version,
            "occurred_at": message.occurred_at,
            "payload": message.payload,
            "meta": message.meta,
        }).encode("utf-8")

    def decode(self, body: bytes) -> MessageEnvelope:
        data = json.loads(body.decode("utf-8"))

        return MessageEnvelope(
            id=data["id"],
            type=data["type"],
            version=int(data["version"]),
            occurred_at=data["occurred_at"],
            payload=data.get("payload", {}),
            meta=data.get("meta", {}),
        )
```

---

## Exceptions

```python
class TemporaryMessageError(Exception):
    pass

class PermanentMessageError(Exception):
    pass
```

Consumer behavior:

```text
TemporaryMessageError -> retry
PermanentMessageError -> DLQ
unknown Exception -> limited retry then DLQ
```

---

## Design Rule

Keep the `messaging` package independent from Django.

Good:

```text
messaging package -> pure Python
django app -> imports messaging and wires settings/commands
```

Bad:

```text
messaging package imports Django models/settings everywhere
```

