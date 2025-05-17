# Dockerfile

# -------------------------
# Etapa de construcción
# -------------------------
    FROM golang:1.24 AS builder

    WORKDIR /app
    
    ARG APPNAME
    ENV GO111MODULE=on

    # Copia los archivos necesarios
    COPY . .
    
    # Instala make (por si no viene en la imagen base)
    RUN apt-get update && apt-get install -y make
    
    # Ejecuta tests, chequea cobertura y construye el binario
    RUN make build
    
    # -------------------------
    # Imagen final mínima
    # -------------------------
    FROM debian:bookworm-slim

    
    WORKDIR /app
    
    ARG APPNAME
    ARG PORT

    # Copiamos el binario desde el build stage
    COPY --from=builder /app/build/${APPNAME} .
    
    # Abrimos el puerto de la API (ajustalo si tu app usa otro)
    EXPOSE ${PORT}
    
    # Comando que se ejecuta al iniciar el contenedor
    #CMD ["sh", "-c", "./${APPNAME}"]
    ENV APPNAME=${APPNAME}
    ENTRYPOINT ["sh", "-c"]
    CMD ["./$APPNAME"]