FROM alpine

WORKDIR /code

COPY priceservice_linux priceservice
COPY partners.csv partners.csv
COPY products.csv products.csv

EXPOSE 8080

CMD ["./priceservice"]

