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
### ????????????????, ?? ?????????????? ?? ?????????????????????? ?????? ?????????????????????????? PostgreSQL 14.4 c pgbouncer (?????????????? ?????????????? ?? ?????????????????????? edoburu/pgbouncer)

?? ?????????????? ?????????? 14???? ???? ?????????????????? ???????????????????????? ?????? ???????????? ???? MD5, ?? SCRAM-SHA-256. ?? ?? ?????????????? ???????????????????? edoburu/pgbouncer, ?????????????? ?? ?????????????????? ?? ?????????????? ???????????? ???????????????? ?? MD5. ???????????? MD5 ?? SCRAM-SHA-256 ?????????? ???? ??????????????????????????. ?? ????????????????????, ?????? PostgreSQL 14.4 ????????, ?????? ?? ???????? ???????????? ?? ??????????????, ?????????????????????????? SCRAM-SHA-256, ?? ???????????????? ?? MD5, ?? ???? ???? ?????????? ?????? ?????????????????????????? ?? ??????????????????.

???????? ?????? ???????? ?????????????? ????????????????, ?????????????? ?????????????????? ???????????? PostgreSQL ?? ?????????????? ?????????????????? edoburu/pgbouncer:
1. ?????????? ???????????? ???????????????????? ?????????????????? ???????????????????? pgbouncer `AUTH_TYPE=plain`, ?????? ?????????? ??????????????, ?????? ???????????? ???????????????? ?? ???????????????????? ?? ???????????????? ????????, ???? ?????? ?????????? ??????????????????. ?????? ?????? ?????? plain (???????????????????????? ????????????), ???? pgbouncer ???????????? ?????? ?????????????????????? ?? SCRAM-SHA-256 ?? ???????????????? ?? PostgreSQL. (auth_file ?????????? ?????????????????? ?????? ?????????????????????????? ?? ?????????????? MD5, SCRAM-SHA-256, ?????? ?? ?????????????????? ???????????? (???????????? edoburu/pgbouncer ???? ?????????? ?????????????????? SCRAM-SHA-256 ???? MD5). ???????????????? ???????? ???????????????????? ?? ?????????????? ?????????????????????????? ???????????????????? https://github.com/edoburu/docker-pgbouncer/blob/master/entrypoint.sh#L51
2. ?????? ???????????????????? ?? ???? PostgreSQL ???????????????? ?????? ???????????????????? ??????????????????:
- `POSTGRES_HOST_AUTH_METHOD: md5` (PostgreSQL ?????????? ?????????????? ???????????? ?? MD5)
- `POSTGRES_INITDB_ARGS: "--auth-host=md5"`  (PostgreSQL ?????? ?????????????????????????? ???????????? ???????????? ?? MD5)
?????????? ?????? ????????????????????, ?? ???????????? ?????????? ???? ?????? ???? ????????????????.
?????????? ???? ???????????????????????? ???????????????????? ???? ???? ???????????????? https://stackoverflow.com/questions/62415752/how-to-get-postgres-docker-container-to-initialize-with-scram-sha-256-on-any-con

edoburu/pgbouncer
https://github.com/edoburu/docker-pgbouncer/blob/master/entrypoint.sh
```sh
  if [ "$AUTH_TYPE" != "plain" ]; then
     pass="md5$(echo -n "$DB_PASSWORD$DB_USER" | md5sum | cut -f 1 -d ' ')"
  else
     pass="$DB_PASSWORD"
  fi
```

_?? ???????? ?????????????? ???????????????????????? ???????????? ??????????????._
