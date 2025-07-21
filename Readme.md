# Guide withdraw flow 🚀

A lightweight Go service for handling guide to withdraw internal flow.  
It's composed by 3 main use cases:
- Client view to search and guide to withdraw.
- Monitor view to display guide staus
- Operator view to allow guide flow management

---

## 📦 Features

- Stateless front end views.
- Stateless restfull api.
- SSE api.
- Authentication integrated using oauth2.
- JWT authorization.
- Guide provider integration.
- Business logic managed by configuration.
- Custom messages internationalized.

---

## 🛠️ Technologies

- Go
- Vue
- Tailwind
- PostgreSQL
- Redis
- Docker

---

## 🚀 Getting Started

### Installation & Deployment

```bash
git clone https://github.com/msartori/via.git
make rebuild ENV=[ENV]
```

### Development

#### View local start

```
cd web
npm run dev
> web@0.0.0 dev
> vite

  VITE v6.3.5  ready in 727 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
  ➜  press h + enter to show help
```
#### Api local start

vscode launch configuration
```
 "version": "0.2.0",
    "configurations": [

        {
            "name": "via API",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${[PROJECT FOLDER]}/api/cmd",
            "env": {
				... env configuration
			}
		}
	]			
```

#### Database & Data Source local start
```
make rebuild ENV=local
```

#### Local database access
```
psql -h localhost -p 5432 -U viauser -d viadb
```

#### Database schema modification
[Ent library](https://github.com/ent/ent)
```
go install entgo.io/ent/cmd/ent@v0.14.4
cd api/ent
ent new [Model struct that will/manages the table]
ent generate ./ent/schema
```

#### JWT Key generation
```
openssl rsa -in private.pem -pubout -out public.pem    
openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
```



