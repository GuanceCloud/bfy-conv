package parse

var ServiceMaps = map[int16][2]string{
	// Pinpoint Internal (0 ~ 999)
	-1:  {"UNDEFINED", ""},
	1:   {"UNKNOWN", ""},
	2:   {"USER", ""},
	3:   {"UNKNOWN_GROUP", ""},
	5:   {"TEST", "test"},
	7:   {"COLLECTOR", ""},
	100: {"ASYNC", ""},
	500: {"SDK", ""},
	510: {"SDK_ASYNC", ""},

	// Server (1000 ~ 1899)
	1000: {"STAND_ALONE", "custom"},
	1005: {"TEST_STAND_ALONE", "test"},
	1010: {"TOMCAT", "web"},
	1011: {"TOMCAT_METHOD", "web"},
	1020: {"RESERVED", ""},
	1021: {"RESERVED", ""},
	1030: {"JETTY", "web"},
	1031: {"JETTY_METHOD", "web"},
	1040: {"JBOSS", "web"},
	1041: {"JBOSS_METHOD", "web"},
	1050: {"VERTX", "web"},
	1051: {"VERTX_INTERNAL", "web"},
	1052: {"VERTX_HTTP_SERVER", "web"},
	1053: {"VERTX_HTTP_SERVER_INTERNAL", "web"},
	1060: {"WEBSPHERE", "web"},
	1061: {"WEBSPHERE_METHOD", "web"},
	1070: {"WEBLOGIC", "web"},
	1071: {"WEBLOGIC_METHOD", "web"},
	1080: {"RESERVED", ""},
	1081: {"RESERVED", ""},
	1100: {"THRIFT_SERVER", "rpc"},
	1101: {"THRIFT_SERVER_INTERNAL", "rpc"},
	1110: {"DUBBO_PROVIDER", "rpc"},
	1120: {"UNDERTOW", "web"},
	1121: {"UNDERTOW_METHOD", "web"},
	1126: {"UNDERTOW_SERVLET_METHOD", "web"},
	1130: {"GRPC_SERVER", "rpc"},
	1140: {"REACTOR_NETTY", "web"},
	1141: {"REACTOR_NETTY_INTERNAL", "web"},
	1300: {"C_CPP", "custom"},
	1301: {"C_CPP_METHOD", "custom"},
	1400: {"NODE", "web"},
	1401: {"NODE_METHOD", "web"},
	1500: {"PHP", "web"},
	1501: {"PHP_METHOD", "web"},
	1550: {"ENVOY", "http"},
	1620: {"OPENWHISK_INTERNAL", "http"},
	1621: {"OPENWHISK_CONTROLLER", "http"},
	1622: {"OPENWHISK_INVOKER", "http"},
	1700: {"PYTHON", "web"},
	1701: {"PYTHON_METHOD", "web"},
	1702: {"CELERY", "queue"},
	1703: {"CELERY-WORKER", "worker"},
	1800: {"GO", "web"},
	1801: {"GO_FUNCTION", "web"},
	// Server Sandbox (1900 ~ 1999)
	// Database (2000 ~ 2899)
	2050: {"UNKNOWN_DB", "db"},
	2051: {"UNKNOWN_DB_EXECUTE_QUERY", "sql"},
	2100: {"MYSQL", "mysql"},
	2101: {"MYSQL_EXECUTE_QUERY", "sql"},
	2102: {"R2DBC_MYSQL", "mysql"},
	2103: {"R2DBC_MYSQL_EXECUTE_QUERY", "sql"},
	2150: {"MARIADB", "mysql"},
	2151: {"MARIADB_EXECUTE_QUERY", "sql"},
	2152: {"R2DBC_MARIADB", "mysql"},
	2153: {"R2DBC_MARIADB_EXECUTE_QUERY", ""},
	2200: {"MSSQL", "db"},
	2201: {"MSSQL_EXECUTE_QUERY", "sql"},
	2250: {"MSSQL_JDBC", "db"},
	2251: {"MSSQL_JDBC_QUERY", "sql"},
	2252: {"R2DBC_MSSQL_JDBC", "db"},
	2253: {"R2DBC_MSSQL_JDBC_QUERY", "sql"},
	2300: {"ORACLE", "db"},
	2301: {"ORACLE_EXECUTE_QUERY", "sql"},
	2302: {"R2DBC_ORACLE", "sql"},
	2303: {"R2DBC_ORACLE_EXECUTE_QUERY", "sql"},
	2400: {"CUBRID", "db"},
	2401: {"CUBRID_EXECUTE_QUERY", "sql"},
	2410: {"NBASET", "db"},
	2411: {"NBASET_EXECUTE_QUERY", "sql"},
	2412: {"NBASET_INTERNAL", "sql"},
	2450: {"INFORMIX", "db"},
	2451: {"INFORMIX_EXECUTE_QUERY", "sql"},
	2500: {"POSTGRESQL", "db"},
	2501: {"POSTGRESQL_EXECUTE_QUERY", "sql"},
	2502: {"R2DBC_POSTGRESQL", "db"},
	2503: {"R2DBC_POSTGRESQL_EXECUTE_QUERY", "sql"},
	2600: {"CASSANDRA", "cassandra"},
	2601: {"CASSANDRA_EXECUTE_QUERY", "sql"},
	2650: {"MONGO", "mongodb"},
	2651: {"MONGO_EXECUTE_QUERY", "sql"},
	2652: {"MONGO_REACTIVE", "mongodb"},
	2700: {"COUCHDB", "db"},
	2701: {"COUCHDB_EXECUTE_QUERY", "sql"},
	2750: {"H2", "db"},
	2751: {"H2_EXECUTE_QUERY", "sql"},
	2752: {"R2DBC_H2", "db"},
	2753: {"R2DBC_H2_EXECUTE_QUERY", "sql"},
	// Database Sandbox (2900 ~ 2999)
	// RESERVED (3000 ~ 4999)
	// Library (5000 ~ 7499)
	5000: {"INTERNAL_METHOD", "custom"},
	5005: {"JSP", "custom"},
	5010: {"GSON", "custom"},
	5011: {"JACKSON", "custom"},
	5012: {"JSON-LIB", "custom"},
	5013: {"FASTJSON", "custom"},
	5020: {"JDK_FUTURE", "custom"},
	5050: {"SPRING", "spring"},
	5051: {"SPRING_MVC", "spring"},
	5052: {"SPRING_ASYNC", "spring"},
	5053: {"SPRING_WEBFLUX", "spring"},
	5061: {"RESERVED", ""},
	5071: {"SPRING_BEAN", "spring"},
	5500: {"IBATIS", "custom"},
	5501: {"IBATIS-SPRING", "spring"},
	5510: {"MYBATIS", "custom"},
	6001: {"THREAD_ASYNC", "custom"},
	6005: {"PROCESS", "custom"},
	6050: {"DBCP", "custom"},
	6052: {"DBCP2", "custom"},
	6060: {"HIKARICP", "custom"},
	6062: {"DRUID", "custom"},
	6500: {"RXJAVA", "custom"},
	6510: {"REACTOR", "custom"},
	6600: {"EXPRESS", "express"},
	6610: {"KOA", "express"},
	6620: {"HAPI", "custom"},
	6630: {"RESTIFY", "custom"},
	6640: {"SPRING_DATA_R2DBC", "spring"},
	7010: {"USER_INCLUDE", "custom"},
	// Library Sandbox (7500 ~ 7999)
	// Cache & File Library (8000 ~ 8899) Fast Histogram
	8050: {"MEMCACHED", "memcached"},
	8051: {"MEMCACHED_FUTURE_GET", "memcached"},
	8100: {"ARCUS", "custom"},
	8101: {"ARCUS_FUTURE_GET", "custom"},
	8102: {"ARCUS_EHCACHE_FUTURE_GET", "custom"},
	8103: {"ARCUS_INTERNAL", "custom"},
	8200: {"REDIS", "redis"},
	8201: {"REDIS_LETTUCE", "redis"},
	8202: {"IOREDIS", "redis"},
	8203: {"REDIS_REDISSON", "redis"},
	8204: {"REDIS_REDISSON_INTERNAL", "redis"},
	8250: {"RESERVED", ""},
	8251: {"RESERVED", ""},
	8260: {"RESERVED", ""},
	8280: {"ETCD", "cache"},
	8300: {"RABBITMQ", "rabbitmq"},
	8310: {"ACTIVEMQ_CLIENT", "queue"},
	8311: {"ACTIVEMQ_CLIENT_INTERNAL", "queue"},
	8660: {"KAFKA_CLIENT", "kafka"},
	8661: {"KAFKA_CLIENT_INTERNAL", "kafka"},
	8800: {"HBASE_CLIENT", "db"},
	8801: {"HBASE_CLIENT_ADMIN", "db"},
	8802: {"HBASE_CLIENT_TABLE", "db"},
	8803: {"HBASE_ASYNC_CLIENT", "db"},
	// Cache Library Sandbox (8900 ~ 8999) Histogram type: Fast
	// RPC (9000 ~ 9899)
	9050: {"HTTP_CLIENT_3", "http"},
	9051: {"HTTP_CLIENT_3_INTERNAL", "http"},
	9052: {"HTTP_CLIENT_4", "http"},
	9053: {"HTTP_CLIENT_4_INTERNAL", "http"},
	9054: {"GOOGLE_HTTP_CLIENT_INTERNAL", "http"},
	9055: {"JDK_HTTPURLCONNECTOR", "http"},
	9056: {"ASYNC_HTTP_CLIENT", "http"},
	9057: {"ASYNC_HTTP_CLIENT_INTERNAL", "http"},
	9058: {"OK_HTTP_CLIENT", "http"},
	9059: {"OK_HTTP_CLIENT_INTERNAL", "http"},
	9060: {"RESERVED", ""},
	9070: {"RESERVED", ""},
	9080: {"APACHE_CXF_CLIENT", "soap"},
	9081: {"APACHE_CXF_SERVICE_INVOKER", "soap"},
	9082: {"APACHE_CXF_MESSAGE_SENDER", "soap"},
	9083: {"APACHE_CXF_LOGGING_IN", "soap"},
	9084: {"APACHE_CXF_LOGGING_OUT", "soap"},
	9100: {"THRIFT_CLIENT", "rpc"},
	9101: {"THRIFT_CLIENT_INTERNAL", "rpc"},
	9110: {"DUBBO_CONSUMER", "rpc"},
	9120: {"HYSTRIX_COMMAND", "rpc"},
	9130: {"VERTX_HTTP_CLIENT", "http"},
	9131: {"VERTX_HTTP_CLIENT_INTERNAL", "http"},
	9140: {"REST_TEMPLATE", "rpc"},
	9150: {"NETTY", "http"},
	9151: {"NETTY_INTERNAL", "http"},
	9152: {"NETTY_HTTP", "http"},
	9153: {"SPRING_WEBFLUX_CLIENT", "http"},
	9154: {"REACTOR_NETTY_CLIENT", "http"},
	9155: {"REACTOR_NETTY_CLIENT_INTERNAL", "http"},
	9160: {"GRPC", "grpc"},
	9161: {"GRPC_INTERNAL", "grpc"},
	9162: {"GRPC_SERVER_INTERNAL", "grpc"},
	9201: {"ElasticsearchBBoss @Deprecated", "elasticsearch"},
	9202: {"ElasticsearchBBossExecutor @Deprecated", "elasticsearch"},
	9203: {"ELASTICSEARCH", "elasticsearch"},
	9204: {"ELASTICSEARCH_HIGHLEVEL_CLIENT", "elasticsearch"},
	9205: {"ELASTICSEARCH8", "elasticsearch"},
	9206: {"ELASTICSEARCH8_CLIENT", "elasticsearch"},
	9301: {"ENVOY_INGRESS", "http"},
	9302: {"ENVOY_EGRESS", "http"},
	9401: {"GO_HTTP_CLIENT", "http"},
	9622: {"OPENWHISK_CLIENT", "http"},
	9700: {"PHP_REMOTE_METHOD", "rpc"},
	9800: {"C_CPP_REMOTE_METHOD", "rpc"},
	9900: {"PYTHON_REMOTE_METHOD", "rpc"},
	// RPC Sandbox (9900 ~ 9999)
}
