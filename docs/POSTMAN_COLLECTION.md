# Postman Collection

## Importar la Colección

Puedes importar la especificación OpenAPI directamente en Postman:

1. Abre Postman
2. Haz clic en **Import**
3. Selecciona **File** o **Link**
4. Importa el archivo `docs/openapi.yaml` o usa la URL si está disponible

## Variables de Entorno

Crea un entorno en Postman con las siguientes variables:

```
BASE_URL: http://localhost:8080
PROFILE_ID: 123e4567-e89b-12d3-a456-426614174000
```

## Endpoints Disponibles

### 1. Health Check
- **Method**: GET
- **URL**: `{{BASE_URL}}/health`

### 2. Get Profile
- **Method**: GET
- **URL**: `{{BASE_URL}}/api/v1/profiles/{{PROFILE_ID}}`

### 3. Get Skills
- **Method**: GET
- **URL**: `{{BASE_URL}}/api/v1/profiles/{{PROFILE_ID}}/skills`

### 4. Get Projects
- **Method**: GET
- **URL**: `{{BASE_URL}}/api/v1/profiles/{{PROFILE_ID}}/projects`

### 5. Get Certificates
- **Method**: GET
- **URL**: `{{BASE_URL}}/api/v1/profiles/{{PROFILE_ID}}/certificates`

### 6. Create Contact
- **Method**: POST
- **URL**: `{{BASE_URL}}/api/v1/profiles/{{PROFILE_ID}}/contacts`
- **Body** (JSON):
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "message": "Hello, I would like to discuss a project opportunity with you."
}
```

## Testing

Puedes agregar tests en Postman para cada request:

### Health Check Test
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has status field", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('status');
});
```

### Profile Test
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has id field", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('id');
    pm.expect(jsonData).to.have.property('name');
});
```

### Contact Test
```javascript
pm.test("Status code is 201", function () {
    pm.response.to.have.status(201);
});

pm.test("Response has id field", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('id');
    pm.expect(jsonData).to.have.property('message');
});
```

