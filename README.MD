# Clean Architecture - Go Expert Fullcyle

Este repositório completa o desafio do Go Expert, Clean architecture.

# Passo a passo para execução
- clonar o repositório
- na raiz do repositório, digitar o comando `docker compose up -d`

# Acessos
- os testes estão na pasta api com *.http
- API Rest nesta rota:
    - Port: 8000
    - POST: [http://localhost:8000/order](http://localhost:8000/order)
    - GET:  [http://localhost:8000/order](http://localhost:8000/order)
- GraphQL:
    - Port: 8080
    - http://localhost:8080/
- GRPC:
    - Port: 50051
    - Package: pb
    - Service: OrderService

## GraphQL samples

### Query
```
query orderSelect{
    orders: orders{
        Id,
        Price,
        Tax,
        FinalPrice
    }
}
```

### Mutation
```
mutation orderCreate { 
  createOrder(input: {
    Price: 4.65,
    Tax: 1.99,
  }) {
    Id,
    Price,
    Tax,
    FinalPrice
  }
}
```
