**Структура проекта**

Применены некоторые концепции чистой архитекту(так, как я их понял). Некоторые модели из доменного слоя используются в слоях представления/инфраструктуры для упрощения.

**API**

Для пагинации выбран метод keyset вместо offset для избежания оверхеда поиска первой строки при большом offset.<br>
Для изменения статуса задачи помимо id задачи необходимо предоставлять id юзера т.к. если задача еще не существует, то ее необходимо будет привязать к юзеру.

**Сущности**

ID: выбран тип данных KSUID вместо UUID т.к. сортировка по нему с большой вероятностью будет эквивалентна сортировке по времени создания. Это, в частности, подразумевает маленькую вероятность того, что при переборе юзеров с пагинацией новый юзер будет пропущен. 