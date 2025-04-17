# LDA API

API для классификации кредитных заявок с использованием Linear Discriminant Analysis

## 📡 API Endpoints

### 1. Получение обучающих данных
**Метод:** `GET`  
**Эндпоинт:** `/selectLearnData`  
**Ответ:**  
```json
{
  "users_data": [
    {
      "loan_id": 1,
      "loan_term": 12,
      "income_annum": 9600000,
      "loan_amount": 29900000,
      "cibil_score": 778,
      "loan_status": "Approved"
    },
    {
      "loan_id": 2,
      "loan_term": 8,
      "income_annum": 4100000,
      "loan_amount": 12200000,
      "cibil_score": 417,
      "loan_status": "Rejected"
    }
  ]
}
```
