Тестовое задание в компанию hezzl
# ЗАДАЧА подробное описание можно найти по [ссылке](https://hezzlcom.atlassian.net/wiki/spaces/PUB/pages/3211332/TASK+5)
1. Развернуть сервис на Golang, Postgres, Clickhouse, Nats (альтернатива kafka), Redis
2. Описать модели данных и миграций
3. В миграциях Postgres
   - Проставить primary-key и индексы на указанные поля
   - При добавлении записи в таблицу устанавливать приоритет как макс приоритет в таблице +1. Приоритеты начинаются с 1
   -  При накатке миграций добавить одну запись в Campaigns таблицу по умолчанию
      - id = serial
      - name = Первая запись
4. Реализовать CRUD методы на GET-POST-PATCH-DELETE данных в таблице ITEMS в Postgres
5. При редактировании данных в Postgres ставить блокировку на чтение записи и оборачивать все в транзакцию. Валидируем поля при редактировании.
6. При редактировании данных в ITEMS инвалидируем данные в REDIS
7. Если записи нет (проверяем на PATCH-DELETE), выдаем ошибку (статус 404)
   - code = 3
   - message = “errors.item.notFound“
   - details = {}
8. При GET запросе данных из Postgres кешировать данные в Redis на минуту. Пытаемся получить данные сперва из Redis, если их нет, идем в БД и кладем их в REDIS
9. При добавлении, редактировании или удалении записи в Postgres писать лог в Clickhouse через очередь Nats (альтернатива kafka). Логи писать пачками в Clickhouse

## Как запустить в docker-compose
```shell
docker-compose up -d --build
```
запускаеться сборка, pulling всех образов.

## Как запустить без докера
1. поднимем окружение
   ```shell
   docker-compose up -d pg_db redis nats clickhouse
   ```
2. В студии запустить как обычно run main
