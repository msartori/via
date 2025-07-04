# VIA Guide Process System

Description

## Development


go install entgo.io/ent/cmd/ent@v0.14.4
ent new Operator
ent new Guide
ent new GuideHistory
# go to api directory
ent generate ./ent/schema

// Run the auto migration tool.
	/*
		if err := entClient.Schema.Create(context.Background()); err != nil {
			logger.Fatal(context.Background(), err, "msg", "failed creating schema resources")
		}*/
	

psql -h localhost -p 5432 -U viauser -d viadb


pkill -f "localtunnel"

export PATH=$PATH:$(go env GOPATH)/bin


fly auth login

fly logs -i [instance]

fly mpg connect --cluster [cluster]




openssl rsa -in private.pem -pubout -out public.pem    

openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048