### Проблема, с которой я столкнулась при использовании PostgreSQL 14.4 c pgbouncer (готовое решение с контейнером edoburu/pgbouncer)

В версиях новее 14ой по умолчанию используется для пароля не MD5, а SCRAM-SHA-256. И в готовом контейнере edoburu/pgbouncer, который я использую в проекте пароль хранится в MD5. Пароль MD5 в SCRAM-SHA-256 никак не преобразовать. И получается, что PostgreSQL 14.4 ждёт, что к нему придут с паролем, зашифрованным SCRAM-SHA-256, а приходят с MD5, и он не может его преобразовать и проверить.

Есть два пути решения проблемы, оставив последнюю версию PostgreSQL и готовый контейнер edoburu/pgbouncer:
1. Можно задать переменную окружения контейнера pgbouncer `AUTH_TYPE=plain`, это будет значить, что пароль хранится в контейнере в открытом виде, да это менее безопасно (хотя, на мой взгляд, риски минимальные). Так как это plain (оригинальный пароль), то pgbouncer сможет его зашифровать в SCRAM-SHA-256 и передать в PostgreSQL. (auth_file может содержать как зашифрованные с помощью MD5, SCRAM-SHA-256, так и текстовые пароли (только edoburu/pgbouncer не умеет создавать SCRAM-SHA-256 из MD5). Проверка этой переменной в скрипте инициализации контейнера https://github.com/edoburu/docker-pgbouncer/blob/master/entrypoint.sh#L51
2. Для контейнера с БД PostgreSQL добавить две переменные окружения:
- `POSTGRES_HOST_AUTH_METHOD: md5` (Postgresql будем ожидать пароль в MD5)
- `POSTGRES_INITDB_ARGS: "--auth-host=md5"`  (Postgresql при инициализации задаёт пароль в MD5)
Нужны обе переменные, с только одной из них не работает.
Нашла на стековерфлоу фактически ту же проблему https://stackoverflow.com/questions/62415752/how-to-get-postgres-docker-container-to-initialize-with-scram-sha-256-on-any-con

edoburu/pgbouncer
https://github.com/edoburu/docker-pgbouncer/blob/master/entrypoint.sh
```sh
  if [ "$AUTH_TYPE" != "plain" ]; then
     pass="md5$(echo -n "$DB_PASSWORD$DB_USER" | md5sum | cut -f 1 -d ' ')"
  else
     pass="$DB_PASSWORD"
  fi
```