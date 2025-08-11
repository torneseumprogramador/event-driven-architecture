#!/bin/bash

# Script para criar tópicos do Kafka e suas DLQs
# Executar após o Kafka estar rodando

KAFKA_BROKERS="localhost:9093"
REPLICATION_FACTOR=1
PARTITIONS=3

echo "Aguardando Kafka estar pronto..."
until kafka-topics --bootstrap-server $KAFKA_BROKERS --list > /dev/null 2>&1; do
    echo "Kafka ainda não está pronto, aguardando..."
    sleep 5
done

echo "Criando tópicos do Kafka..."

# Tópicos de usuário
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic user.created --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic user.updated --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic user.created.dlq --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic user.updated.dlq --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists

# Tópicos de produto
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic product.created --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic product.updated --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic stock.reserved --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic stock.released --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic product.created.dlq --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic product.updated.dlq --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic stock.reserved.dlq --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic stock.released.dlq --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists

# Tópicos de pedido
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic order.created --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic order.paid --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic order.canceled --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic order.created.dlq --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic order.paid.dlq --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists
kafka-topics --bootstrap-server $KAFKA_BROKERS --create --topic order.canceled.dlq --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR --if-not-exists

echo "Listando tópicos criados:"
kafka-topics --bootstrap-server $KAFKA_BROKERS --list

echo "Tópicos criados com sucesso!"
