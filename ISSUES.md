### The problem I faced when using PostgreSQL 14.4 with pgbouncer (ready-made solution with container edoburu/pgbouncer)

In versions newer than 14, by default SCRAM-SHA-256 is used for the password and not MD5. And in the ready-made container edoburu/pgbouncer that I use in the project, the password is stored in MD5. There is no way to convert MD5 password to SCRAM-SHA-256. So it turns out that PostgreSQL 14.4 expects to receive password encrypted with SCRAM-SHA-256, but receives password encrypted with MD5, that's why PostgreSQL 14.4 can't convert and verify password.

There are two ways to solve the problem, while using the latest version of PostgreSQL and the ready-made container edoburu/pgbouncer:
1. Set the container pgbouncer environment variable `AUTH_TYPE=plain`, it would mean, that the password is stored in the container in plain text, it is less safe. As it is plain (original password), pgbouncer will be able to encrypt it in SCRAM-SHA-256 and transfer it into PostgreSQL. (auth_file can contain both MD5-encrypted, SCRAM-SHA-256, and clear-text passwords (only edoburu/pgbouncer can't create SCRAM-SHA-256 from MD5). This variable is checked in the container initialization script https://github.com/edoburu/docker-pgbouncer/blob/master/entrypoint.sh#L51
2. Add two environment variables for PostgreSQL database container:
- `POSTGRES_HOST_AUTH_METHOD: md5` (PostgreSQL will expect a password in MD5)
- `POSTGRES_INITDB_ARGS: "--auth-host=md5"`  (PostgreSQL sets the password in MD5 during initialization)
Both variables are needed, it doesn't work with only one of them.
I actually found the same problem in stackoverflow: https://stackoverflow.com/questions/62415752/how-to-get-postgres-docker-container-to-initialize-with-scram-sha-256-on-any-con

edoburu/pgbouncer
https://github.com/edoburu/docker-pgbouncer/blob/master/entrypoint.sh
```sh
  if [ "$AUTH_TYPE" != "plain" ]; then
     pass="md5$(echo -n "$DB_PASSWORD$DB_USER" | md5sum | cut -f 1 -d ' ')"
  else
     pass="$DB_PASSWORD"
  fi
```

_In my solution the second variant is used._
\
&nbsp;

---
### Проблема, с которой я столкнулась при использовании PostgreSQL 14.4 c pgbouncer (готовое решение с контейнером edoburu/pgbouncer)

В версиях новее 14ой по умолчанию используется для пароля не MD5, а SCRAM-SHA-256. И в готовом контейнере edoburu/pgbouncer, который я использую в проекте пароль хранится в MD5. Пароль MD5 в SCRAM-SHA-256 никак не преобразовать. И получается, что PostgreSQL 14.4 ждёт, что к нему придут с паролем, зашифрованным SCRAM-SHA-256, а приходят с MD5, и он не может его преобразовать и проверить.

Есть два пути решения проблемы, оставив последнюю версию PostgreSQL и готовый контейнер edoburu/pgbouncer:
1. Можно задать переменную окружения контейнера pgbouncer `AUTH_TYPE=plain`, это будет значить, что пароль хранится в контейнере в открытом виде, да это менее безопасно. Так как это plain (оригинальный пароль), то pgbouncer сможет его зашифровать в SCRAM-SHA-256 и передать в PostgreSQL. (auth_file может содержать как зашифрованные с помощью MD5, SCRAM-SHA-256, так и текстовые пароли (только edoburu/pgbouncer не умеет создавать SCRAM-SHA-256 из MD5). Проверка этой переменной в скрипте инициализации контейнера https://github.com/edoburu/docker-pgbouncer/blob/master/entrypoint.sh#L51
2. Для контейнера с БД PostgreSQL добавить две переменные окружения:
- `POSTGRES_HOST_AUTH_METHOD: md5` (PostgreSQL будет ожидать пароль в MD5)
- `POSTGRES_INITDB_ARGS: "--auth-host=md5"`  (PostgreSQL при инициализации задаёт пароль в MD5)
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

_В моем решении используется второй вариант._
