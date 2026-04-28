
# Modulul Backend — Orders (Go)

Documentație pentru modulul backend al proiectului Orders.


## Descriere scurtă

Backend-ul este implementat în Go folosind framework-ul Gin pentru HTTP și GORM pentru lucrul cu PostgreSQL. Modulul conține modele, repository-uri, strat de servicii și handler-e HTTP care gestionează REST API pentru administrarea utilizatorilor, clienților, contractelor și comenzilor.


## Structură (fișiere importante)

- `cmd/server/main.go` — punctul de intrare: încărcarea configurației, conectarea la BD, migrații și pornirea routerului.
- `internal/config/config.go` — configurația aplicației (DSN, secrete JWT etc.).
- `internal/models/models.go` — modele GORM: `User`, `Client`, `Contract`, `ContractAddress`, `Product`, `Order`, `OrderItem` și altele.
- `internal/repository` — straturi de repository pentru lucrul direct cu GORM (operații CRUD).
- `internal/service` — logică de business: wrapper peste repository-uri, validare, operații suplimentare.
- `internal/api/handlers.go` — înregistrarea rutelor și handler-e HTTP generale (login, signup, rute protejate).
- `internal/api/clients.go` *(recomandat)* — handler-e separate pentru clienți/contracte (dacă sunt adăugate).

> Pentru a găsi puncte de intrare specifice folosiți `grep`/IDE: `SetupRoutes`, `NewService`, `AutoMigrate`.


## Principii arhitecturale

- Separare clară a responsabilităților: handlers → service → repository → models.
- Autorizare prin JWT: `/login` returnează un token care se folosește în header-ul `Authorization: Bearer <token>` pentru rutele protejate.
- Migrațiile se execută automat la pornire (`AutoMigrate`) — verificați schema reală la deploy în producție.


## Modele (pe scurt)

- User — utilizator, stochează email, parolă (hash), rol și legătură la canalele de vânzări. Poate fi nullable.
- Client — client (beneficiar): câmpuri `Name`, `Email`, `Phone`, `Address`, `UserID` (proprietar).
- Contract — contractul clientului: `ClientID`, `Number`, `Date`, `Amount`, `Status` și legături către `Client` și `ContractAddress`.
- ContractAddress, Product, Order, OrderItem — entități auxiliare pentru adrese, produse și comenzi.

Dacă este necesar, deschideți `internal/models/models.go` pentru structura completă și tag-urile GORM.


## Principale endpoint-uri API (exemplu)

Toate exemplele presupun că serverul ascultă la `http://localhost:8080`.

- POST /login — autentificare: primește JSON `{ "email": "...", "password": "..." }`, returnează JSON `{ "token": "..." }`.

- POST /clients — creează client (protejată): header `Authorization: Bearer <token>`; body: `{ "name":"ACME", "email":"acme@example.com", "phone":"...", "address":"..." }`. `UserID` se recomandă să fie preluat din token pe server.

- GET /clients/:id — obține client după id (inclusiv contractele asociate, dacă service folosește Preload).

- POST /clients/:id/contracts — creează contract pentru client: body `{ "number":"CTR-001","date":"2025-11-01","amount":1000.0,"status":"active" }`.

- GET /contracts/:id — obține contract după id.

Verificați `internal/api/handlers.go` pentru harta completă a rutelor și middlewares.


## Exemple curl

Obține token (login):


```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret"}'
```

Creează client (înlocuiți `<token>`):

```bash
curl -X POST http://localhost:8080/api/v1/clients \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"ACME","email":"acme@example.com","phone":"123","address":"addr"}'
```

Creează contract pentru clientul id=1:

```bash
curl -X POST http://localhost:8080/api/v1/clients/1/contracts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"number":"CTR-001","date":"2025-11-01","amount":1000.00,"status":"active"}'
```


## JWT și testare în Postman

1. Efectuați `POST /login` cu date reale, salvați `token` în environment-ul Postman (Tests script):

```javascript
const json = pm.response.json();
pm.environment.set("auth_token", json.token);
```

2. În Authorization pentru următoarele request-uri selectați `Bearer Token` și folosiți `{{auth_token}}`.

3. Dacă trebuie să generați manual JWT în Pre-request, folosiți Secret și algoritmul (HS256). Dar este mai sigur să obțineți token-ul de la server.


## Migrații și deploy

- În `cmd/server/main.go` se execută `db.AutoMigrate(...)` la pornire. Pentru producție este mai bine să folosiți migrații explicite (migrate tool) și procese de rollback.
- Configurația BD se ia din `internal/config` (DSN). Asigurați-vă că variabilele de mediu sunt setate.


## Inițializarea serviciului (exemplu)

În `cmd/server/main.go` se creează repository-ul și serviciul și se transmit către rute:

```go
repo := repository.NewRepository(db)
svc := service.NewService(repo, cfg.JWTSecret)
api.SetupRoutes(r, svc)
```

Service conține metode pentru: înregistrare/login, creare clienți/contracte, obținere entități etc.


## Sfaturi de securitate și validare

- Parolele se stochează doar ca hash (bcrypt).
- Verificați ownership-ul: client.UserID trebuie completat pe server (din token), nu permiteți clientului să trimită user_id străin.
- Validați unicitatea (email, număr contract) — verificați în DB și returnați erori clare.


## Ce se poate îmbunătăți / next steps

- Adăugați o colecție Postman cu exemple de login și CRUD pentru clienți/contracte.
- Adăugați teste unit/integration pentru stratul de servicii și handler-e (httptest + gin).
- Migrați la un instrument de migrații (golang-migrate) pentru migrații controlate.
- Adăugați logare structurată pentru request-uri și erori.


## Unde găsiți codul

- Handler-e API: `internal/api/handlers.go` (+ fișiere suplimentare în `internal/api/`).
- Logică de business: `internal/service/*`.
- Repository/DB: `internal/repository/*`.
- Modele: `internal/models/models.go`.
