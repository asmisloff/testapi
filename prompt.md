\s code_go
\m claude
Мы пишем утилиту на Go. Утилита будет иметь конфигурационный файл в формате JSON.
Напиши модуль для чтения json-файла и отображения его в структуру. Формат json-файла ниже.
```json
{
    "base-url": "http://localhost:8080/api/v1/request",
    "auth": "user",
    "set": [
        {
            "id": "01", 
            "method": "GET",
            "params": {
                "filter": [
                    {"name": "RECORD_NUMBER", "values": [1, 2, 3], "type": "IN"}
                ],
                "login": "apushkin"
            }
        }
    ]
}
```
Модуль чтения конфигурации будет называться config.