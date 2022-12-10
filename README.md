# Mailmerger Server

A server for [Mailmerger](https://github.com/fahmifan/mailmerger)

## Features
- [x] Manage campaigns
- [x] Custom templates
- [x] Preview email templates 
- [x] Blast emails from imported csv files

## Development
Tools to install:
- [rubenv/sql-migrate: SQL schema migration tool for Go.](https://github.com/rubenv/sql-migrate)
- [cortesi/modd: A flexible developer tool that runs processes and responds to filesystem changes](https://github.com/cortesi/modd)
- [Install Docker Engine | Docker Documentation](https://docs.docker.com/engine/install/)

### Start development server
Copy the `example.dbconfig.yml` adjust to your local setup. 

```bash
# run postgres & create the db
make run-dev-postgres
# run schema migration
make migrate-up

# run mailhog the email server
make run-dev-mailhog
# run server using modd
make run-server
```

# ERD
```mermaid
erDiagram
	Campaign ||--o| File: have
	Campaign ||--o| Template: have
	Campaign ||--o{ Event: have

	Campaign {
		id string
		file_id string
		template_id string
		subject string
		body string
	}
	File {
		id string
		name string
		path string
	}
	Template {
		id string
		name string
		html string
	}
	Event {
		id string
		created_at timestamp
		status string
		detail string
	}

```

### Pages
![home](doc/localhost_8080_.png)
![campaigns](doc/localhost_8080_campaigns.png)
![campaigns show](doc/localhost_8080_campaigns_01GAN44ZMFYAHJYKQ3BYDQRJTX.png)
![campaigns edit](doc/localhost_8080_campaigns_01GAN44ZMFYAHJYKQ3BYDQRJTX_edit.png)
![templates](doc/localhost_8080_templates.png)
![templates show](doc/localhost_8080_templates_01GAMXWK6BBZE6W4QHZMKBB7C5.png)
![templates edit](doc/localhost_8080_templates_01GAMXWK6BBZE6W4QHZMKBB7C5_edit.png)