# Go-project---arithmetic-expressions
# Примеры (переходим по ссылке: http://127.0.0.1:8080/)
*для того чтобы увидеть результат надо обновить главную страницу страницу
![example](https://github.com/alek-pc/Go-project---arithmetic-expressions/blob/main/src/example.png)

*деление только целочисленное (2/8 = 0), поэтому при подборе примеров учитывайте это
2 + 2 * 2 = 6 \n
5 * 10 / 5 - 60 / 15 + 20 = 46
120 * 5 - 95 * 3 + 65 / 5 = 328
50 * 5 / 2 - 65 / 13 = 120
10 * 10 * 60 / 100 = 60 

10 * 0 = 0 

Схема программы в кратком виде:
![scheme](https://github.com/alek-pc/Go-project---arithmetic-expressions/blob/main/src/Scheme.jpg)

Начнем по порядку:
# Главная страница /
1 пункт - главная страница. С нее на сервер через форму посылаются post-запросы с выражением
Они попадают в функцию server.GettingResponse, здесь проверяется, не подсчитано ли выражение, если да, то отправляется сообщение "expression is already computing", если нет - выражение -> struct Expression отправка выражения в SeparateExpression, где выражение проверяется на валидность и из строки превращается в массив типа "2 + 2 * 2" -> {"2" + "2" * "2"}
Если выражение правильное (только цифры 0123456789 и символы +-*/), то структура Expression добавляется в Storage (Expressions []Expression)
Выражение (уже структура Expression) -> orchestrator.Expression -> server.orchestrator

# Оркестратор
находится в server/orchestrator
Выражение переводится в operations - массив объектов Operation
Operation поля: num1 int - первое число, oper string - операция в выражении, num2 int - второе число, res int - результат выражения, status bool - статус выражения: true - посчитано, false - считается
В бесконечном цикле проходимся по operations, берем число сейчас и следующее, если операция посчитана, то проверяем другие условия и закидываем в Agent (считает только одну операцию)
В цикле массив operations уменьшается и когда его длина становится нулем, то выражение подсчитано - записываем время конца подсчета, результат и обновляем статус на true

В GettingResponse проверяем expression на status, если статус - true, storage.Update() - обновляем сторэдж и файл expressions_logs.csv

#  Страница настроек /settings
Чекаем форму, чекаем ее содержимое, загружаем изменения в settings.csv (settings.Upload()), показываем настройки (SettingsPage)

# Страница воркеров /workers
воркер - условный сервер, на котором выполняются операции - агенты (во всяком случае, я так понял условие)
У каждого воркера кол-во занятых агентов

# функция main:
Подключение пакетов
Инициализирование server.Storage (загрузка (storage.Download) данных из expressions_logs.csv) и закидываем в megaServer (SendStorage)
Инициализирование settings.Settings (загрузка (settings.Upload) данных из csv) и закидываем в settings (SendSettings)
запуск сервера


# Дизайн страниц
![MainPage](https://github.com/alek-pc/Go-project---arithmetic-expressions/blob/main/src/Main%20page%20design.png)
1 - строка ввода выражения (тег input)
2 - сообщение от сервера
3 - список выражений (зеленый - посчитано) строка вида: выражение строкой = ответ start: время начала вычислений end: время конца вычислений
4 - переход на страницу настроек. Клик!
(какой переход!)
![Settings page](https://github.com/alek-pc/Go-project---arithmetic-expressions/blob/main/src/settings%20page%20design.png)

1 - переход на главную страницу
2 - настройка времени выполнения сложения от 1 до 101 (в секундах)
3 - настройка времени вычитания от 1 до 101 (в секундах)
4 - настройка умножения от 1 до 101 (в секундах)
5 - настройка деления от 1 до 101 (в секундах)
6 - кол-во воркеров (агентов) от 0 до 20
7 - отправка формы


![workers page](https://github.com/alek-pc/Go-project---arithmetic-expressions/blob/main/src/workers_page.png)
1 - список воркеров (типа серверы), у каждого подписано кол-во занятых агентов
2 - переход на другие страницы

По вопросам связывайтесь в тг: (https://t.me/alek_asd)



