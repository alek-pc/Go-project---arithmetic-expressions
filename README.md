# Go-project---arithmetic-expressions

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

# функция main:
Подключение пакетов
Инициализирование server.Storage (загрузка (storage.Download) данных из expressions_logs.csv) и закидываем в megaServer (SendStorage)
Инициализирование settings.Settings (загрузка (settings.Upload) данных из csv) и закидываем в settings (SendSettings)
запуск сервера


# Примеры
![MainPage](https://github.com/alek-pc/Go-project---arithmetic-expressions/blob/main/src/Main%20page%20design.png)
