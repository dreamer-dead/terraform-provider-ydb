terraform {
  required_providers {
    ydb = {
       source = "terraform.storage.ydb.tech/provider/ydb"
    }
  }
}

# To be discussed.
# Do we really any configuration parameters except token?
provider "ydb" {
  token = "my_token"
}

resource "ydb_table" "table1" {
  path              = "path/to/table"              # Will create table at /path/to/my/table
  database_endpoint = "grpcs://ydb.serverless.cloud-preprod.yandex.net:2135?database=/pre-prod_ydb_public/aoedo0ji1lgce9l91har/cc8pfiaj0ab96vmvp5v8"

  // ТОЛЬКО ДОБАВЛЯЕМ КОЛОНКИ, НЕ УДАЛЯЕМ!!!
  column {
    name   = "a"
    type   = "Uint64"
    not_null = true // default = false
  }
  column {
    name   = "b"
    type   = "Uint8"
    not_null = true
  }
  column {
    name   = "c"
    type   = "Text"
    # not_null = true
  }
  column {
    // Сравнение колонок по именам. Создание -- смотрим на порядок. Потом -- нет.
    name = "d"
    type = "Timestamp" // YQL types
    # not_null = true
  }
  column {
    name = "e"
    type = "Bytes"
    not_null = false
  }
  /*
    PRIMARY KEY ( column, ... ),
    FAMILY column_family ( family_options, ... )
  */

  # family {
  #   name = "family_name"
  #   data = "ssd"
  #   compression = "off"
  # }
  # family {
  #   name = "name2"
  #   data = "hdd"
  #   compression = "lz4"
  # }

  primary_key = [
    "a", "b"
  ] // Can not be changed or altered: error on modification ops.


  // TODO(shmel1k@): А мы ждём вообще создание индекса в терраформе? Он же может день идти. Другие операции могут ждать и понимать, что БД готова к обновлению приложения.
  // operations list + watch till created.
  // XXX wait_async_operations = true // Ждать, пока все операции применятся успешно.
  // XXX: дожидаемся. Флажок не пилим.

  // MODIFY INDEX ONLY THROUGH DROP + CREATE
  index {
      name    = "index_1_name"
      columns = ["b", "a", "c"]
      type    = "global_sync" // global_async
      cover   = ["d", "e", "f"]
  }
  index {
      name    = "index_2_name"
      columns = ["a", "c", "b"]
      type    = "global_sync" // global_async
      cover   = ["d", "e", "f"]
  }
  // TODO: А не подождать ли команды для атомарной модификации индекса?..
  // Инструмент миграции!

  ttl { // Can be dropped, modified, created, etc.
    column_name          = "d" # Колонка должна присутствовать в списке колонок. // modifiable. Меняется через RESET + CREATE.
#    mode                 = "date_type" // mode = "since_unix_epoch" // modifiable. Меняется через RESET + CREATE.
    expire_interval = "PT0S" // modifiable. Меняется через RESET + CREATE.
    // https://ydb.tech/en/docs/concepts/ttl - change to ISO 8601
  }

  // changefeed { // Делается через ALTER
  //   mode = "KEYS_ONLY" // https://ydb.tech/en/docs/yql/reference/syntax/alter_table#changefeed-options
  //   format = "JSON"
  // }

  partitioning_settings { // https://ydb.tech/en/docs/concepts/datamodel/table
    auto_partitioning_by_size_enabled = true
    auto_partitioning_by_load = true
    auto_partitioning_partition_size_mb = 1024
    auto_partitioning_min_partitions_count = 1
    auto_partitioning_max_partitions_count = 2
    uniform_partitions = 2
    partition_at_keys {
      keys = [100, "1000"]
    }
    // can be set only on create
    // PARTITION_AT_KEYS - ONLY ON CREATE
    // UNIFORM_PARTITIONS - ONLY ON CREATE
    // Остальное -- изменяем, как нам скажут.
  }

  read_replicas_settings = "PER_AZ:3"

  key_bloom_filter = true # Дефолт -- false

  // terraform specific
  lifecycle {
    ignore_changes = [
      column, // disables alter
      partitioning_settings, // disables partitioning_settings changes
    ]
  }
}
