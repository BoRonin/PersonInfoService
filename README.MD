# Сервис обогащения данными со сторонних API по имени пользователя

Сервис использует Postgres в качестве DB и Redis в качестве кэша. Запускается через docker compose up. 
Миграции через утилиту migrate или через docker.

Makefile:

Если через докер:

migrate_dock_up:
	docker run -v $(P)/db/migration:/migrations --network emtest_emtest migrate/migrate -path=/migrations/ -database 'postgres://emtest:emtest@postgres:5432/emtest?sslmode=disable' up

migrate_dock_down:
	docker run -v $(P)/db/migration:/migrations --network emtest_emtest migrate/migrate -path=/migrations/ -database 'postgres://emtest:emtest@postgres:5432/emtest?sslmode=disable' down -all

Передавать через переменную P абсолютный путь к папке с миграциями.

## /person
### POST
Отправляем json в следующем виде

```bash
{
    "name": "Dmitriy",
    "surname": "Ushakov",
    "patronymic": "Vasilevich" // необязательно
}
```

Сервис получает дополнительную информацию, складывает ее в Postgres и кешарует в Redis на некоторое время, чтобы не делать повторных запросов по API. При повторном аналогичном запросе (имя) данные будут браться из кэша и ключ будет продлен.

### DELETE

Запрос с URL параметром (/person/{id}), указывающим на ID, по которому будет удален человек из базы данных

### PUT

Запрос с URL параметром (/person/{id}) и телом в следующем виде

```bash
{
    "name": "oleg",
    "surname": "ko",
    "nationality": "hehe"
}
```

Человек с данным ID будет изменен

## /persons

### GET

Запрос с URL параметрами 
#### page 
страница для выдачи

#### per_page 
элементов на странице

#### order_by
сортировка отфильтрованных позиций. Варианты:
______________________________________
age_asc - возраст, возрастающая
age_desc - возраст, убывающая
name_asc - имя, возрастающая
name_desc - имя, убывающая
nationality_asc - национальность, возрастающая
nationality_desc - национальность, убывающая
gender_asc - пол, возрастающая
gender_desc - пол, убывающая
surname_asc - фамилия, возрастающая
surname_desc - фамилия, убывающая
______________________________________

#### name
поиск по именам, можно добавлять больше пар с ключом name

#### surname
поиск по фамилиям, можно добавлять больше пар с ключом surname

#### patronymic
поиск по отчествам, можно добавлять больше пар с ключом patronymic

#### age
поиск по возрасту, можно добавлять больше пар с ключом age

#### nationality
поиск по национальности, можно добавлять больше пар с ключом nationality

#### gender
поиск по полу

#### age_gt
поиск возрасту, больше чем указаном в значении

#### age_lt
поиск возрасту, меньше чем указаном в значении