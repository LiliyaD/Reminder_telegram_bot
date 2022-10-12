# Telegram Bot Reminder

## Telegram bot for recording daily activities

The application is divided into 2 services: interface service and data service.
The interface part validates the data and passes it on.
The data service saves and retrieves data.
Communication between services occurs via gRPC. ActivityListStream is a stream method for getting a list of activities.
For simplicity, further in the text, the service interface will be called the client, and the service for working with data will be called the server.

### _Available commands:_
- /help - list commands
- /examples - show examples
- /list, /list_stream - list all activities
- /get <name> - get information about activity
- /today - list today activities
- /add <name> <begin_date> <end_date> <times_per_day> <quantity_per_time> - add new activity
- /delete <name> - delete activity
- /update <name> <begin_date> <end_date> <times_per_day> <quantity_per_time> - update activity, if you don't want to change all fields write _ for parameters, which you won't change

The server has functions for working with storage.
The default storage is the PostgreSQL database, but you can use only the local cache (if you call the server with the --local key).
Standard CRUD operations and return of a list of saved entities are implemented for the storage.
The data warehouse can work in a multi-threaded environment, the number of simultaneous accesses to the storage is limited to 10.
The storage has a context and 2 seconds timeouts so waiting for free workers is limited in time.

In addition to using a telegram bot, you can simply make gRPC calls or http calls from the client side. There is not only a telegram bot as a client, but also gRPC, http interfaces. The project also added a gRPC gateway. The project also provides swagger for API documentation.
The Telegram bot, in turn, uses gRPC to get its data from the server.

What has been done additionally for PostgreSQL:
- Added idempotent migrations with table initialization.
- Added parameters for pagination (limit, offset, order).
- PostgreSQL is accessed via a connection pool (pgbouncer).

Also, service communication is implemented through the Kafka message broker, instead of gRPC for some operations (this is not necessary, I just wanted to get to know Kafka). The Producer/Consumer pattern was used. Producer sends data from client service to Kafka. On the server side the consumer is listening to see if something is transfered to Kafka and performs one of the CRUD operations, depending on the topic of the message. This implementation can be seen in commands_kafka.go, before CRUD operations were changed for Kafka, they used gRPC call to data service part, this implementation can be seen in commands_grpc.go.

Added cache work with Redis for faster data retrieval on serious data work. Before attempting to read and add an entry to the database, the cache is checked first. If it contains such data, the result is returned, otherwise the database is checked. When a record is successfully read from the database, it is added to the cache. When updating or deleting an entry in the database, the entry is removed from the cache too.
To receive responses from CRUD operations, messages about which were recorded in Kafka, a notification mechanism is implemented using the Publish–subscribe Redis pattern. On the server side, Consumer, having received and processed the operation, sends the result by the same id to Redis, and methods that sent to Kafka listens (implemented for deletion)/periodically checks (implemented for creating and updating) Redis and return the result.

Unit tests and integration tests were written for part of the code, mocks were made for external calls (PostgreSQL, gRPC)

Added logging, log is printed to the console. The mechanism to write log to the file is provided too. There is also a division by logging levels.
Added trace and counters for input, output, success, failed requests and errors, as well as hit and miss counters to work with the cache.

_This project was done in stages in terms of learning the language and new technologies, so some moments in the code are strange._
\
&nbsp;

---
# Телеграм бот для напоминаний

## Телеграм бот для записи ежедневных активностей

Приложение разделено на 2 сервиса: сервис интерфейс и сервис по работе с данными.
Интерфейсная часть осуществляет валидацию данных и передает их дальше.
Сервис работы с данными осуществляет непосредственное сохранение и получение данных.
Общение между сервисами происходит по gRPC. ActivityListStream - потоковый метод для получения списка активностей.
Для упрощение далее в тексте сервис интерфейс будем называть клиентом, а сервис по работе с данными сервером.

### _Доступные команды:_
- /help - перечисление всех команд
- /examples - примеры
- /list, /list_stream - вывести все активности
- /get <name> - получить информацию об активности
- /today - вывести сегодняшнии активности
- /add <name> <begin_date> <end_date> <times_per_day> <quantity_per_time> - добавить новую активность
- /delete <name> - удалить активность
- /update <name> <begin_date> <end_date> <times_per_day> <quantity_per_time> - обновить активность, если какие-то поля не требуется менять, напишите "_" на месте таких параметров

Сервер имеет функции по работе с хранилищем.
В качестве хранилища по умолчанию используется БД PostgreSQL, но можно использовать и только локальный кеш (если вызвать сервер с ключом --local).
Для хранилища реализованы стандартные CRUD операций + возврат списка сохраненных сущностей.
Хранилище данных может работать в многопоточной среде, количество одновременных обращений к хранилищу ограничено до 10.
У хранилища задан контекст и таймауты в 2 секунды, чтобы ожидание свободных воркеров было ограниченно по времени.

Помимо использования телеграм бота, можно просто делать gRPC вызовы или http вызовы со стороны клиента. То есть в качестве клиента выступает не только телеграм бот, но но и gRPC, http интерфейсы. В проект также добавлен gRPC gateway. В проекте также предусмотрен swagger для документации API.
Телеграм бот в свою очередь использует gRPC, чтобы получить данные от сервера.

Что сделано дополнительно для PostgreSQL:
- Добавлены идемпотентные миграции с инициализацией таблиц.
- Добавлены параметры для пагинации (limit, offset, order).
- Доступ к PostgreSQL осуществляется через пул соединений (pgbouncer).

Так же общение сервисов реализовано через брокер сообщений Kafka, вместо gRPC для части операций (в этом нет необходимости, просто хотелось познакомиться с Kafka). Для реализации использовался паттерн Producer/Consumer. Producer отправляет данные из сервиса-клиента в Kafka. Consumer на стороне сервера слушает, не пришло ли что-то в Kafka и исполяет одну из CRUD операций, в зависимости от топика сообщения. Эту реализацию можно увидеть в commands_kafka.go, до того как CRUD операции были изменены для Kafka, для них использовался gRPC вызов сервисной части работы с данными, эту реализацию можно увидеть в commands_grpc.go.

Добавлена работа с кешем с помощью Redis для более быстрого получения данных на серевере работы с данными. Перед попыткой считать и добавить запись в БД, сначала проверяется кеш. Если в нём есть такие данные, выдается результат, иначе проверяется БД. При успешном чтении записи из БД, она добавляется в кеш. При обновлении и удалении записи в БД, запись также удаляется из кеша.
Для получения ответов по CRUD операциям, сообщения о которых записывались в Kafka реализован механизм оповещения с помощью паттерна Publish–subscribe Redis, на стороне сервера Consumer получив и обработав операцию отправляет результат по ней по такому же id в Redis, а методы, которые делали отправку в Kafka слушают(реализовано для удаления)/переодически проверяют(реализовано для создания и обновления) Redis и выдают результат.

Для части кода написаны unit-тесты и интеграционные тесты, сделаны моки для для внешних вызовов (PostgreSQL, gRPC)

Добавлено логгирование, вывод лога в консоль, в логе предусмотрен механизм, чтобы писать его в файл. Есть разделение по уровням логгирования.
Добавлен trace и счетчики для входящих, исходящих, успешных, неудачных запросов и ошибок, а также счетчики hit и miss для работы с кешем.

_Этот проект делался поэтапно в плане изучения языка и новых технологий, поэтому некоторые моменты в коде странные._
