# Mailmerger Server

A server for Mailmerger

Run server:
```
make run-server
```

# ERD
```mermaid
erDiagram
	Campaign ||--o| Csv: have
	Campaign ||--o| Template: have

	Campaign {
		string id
		string csv_id
		string template_id
	}
	Csv {
		string id
		name string
		path string
	}
	Template {
		id string
		name string
		body string
		subject string
	}

```