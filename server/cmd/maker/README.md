Metadata maker (Magia în acțiune)

Deschide terminalul la rădăcina proiectului și tastează comanda:
Bash

go run cmd/maker/main.go product

Vei vedea acest meniu interactiv. Răspunde-i așa:
Plaintext

🤖 GOrders Maker: Creăm modulul 'Test'
👉 Introdu coloanele (format: nume_coloana:tip:req). Scrie 'exit' când ai terminat.
Column: code:string:req
Column: name:string
Column: price:float64
Column: exit

✅ Generated: internal/models/mod_test.go
✅ Generated: internal/dto/dto_test.go
✅ Generated: internal/service/svc_test.go
✅ Generated: internal/handler/hdl_test.go

🎉 GATA! Tot scheletul a fost generat cu succes.