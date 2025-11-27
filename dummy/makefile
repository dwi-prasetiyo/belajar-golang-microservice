# Redis cluster
CONF_DIR=./conf
PREFIX=redis-node
PORT=6379
BIND_IP=0.0.0.0
CLUSTER_ENABLED=yes
REQUIREPASS=rahasia
MASTERAUTH=rahasia
CLUSTER_NODE_TIMEOUT=5000
CLUSTER_ANNOUNCE_PORT=6379
CLUSTER_ANNOUNCE_BUS_PORT=16379
APPENDONLY=yes
NODES=1 2 3 4 5 6

# Redis
.PHONY: redis-conf
redis-conf: ${NODES:%=${CONF_DIR}/${PREFIX}-%.conf}

.PHONY: ${CONF_DIR}/${PREFIX}-%.conf
${CONF_DIR}/${PREFIX}-%.conf: ${CONF_DIR}
	@echo "Creating $@"
	@echo "port ${PORT}" > $@
	@echo "bind ${BIND_IP}" >> $@
	@echo "cluster-enabled ${CLUSTER_ENABLED}" >> $@
	@echo "requirepass ${REQUIREPASS}" >> $@
	@echo "masterauth ${MASTERAUTH}" >> $@
	@echo "cluster-config-file node-$*.conf" >> $@
	@echo "cluster-announce-ip 172.48.0.1$*" >> $@
	@echo "cluster-announce-port 6379" >> $@
	@echo "cluster-announce-bus-port 16379" >> $@
	@echo "appendonly ${APPENDONLY}" >> $@

.PHONY: ${CONF_DIR}
${CONF_DIR}: 
	mkdir -p ${CONF_DIR}

.PHONY: clean-redis-conf
clean-redis-conf:
	rm -f ${CONF_DIR}/${PREFIX}-*.conf

.PHONY: all-redis-conf
all-redis-conf: clean-redis-conf redis-conf
