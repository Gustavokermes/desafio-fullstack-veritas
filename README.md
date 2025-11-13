# Desafio Fullstack – Mini Kanban de Tarefas (React + Go)

## Autor
Gustavo Henrique Kermes Sousa

## Descrição
Este projeto é um Kanban simples fullstack, com backend em Go e frontend em React, que permite criar, editar, mover e excluir tarefas entre três colunas: A Fazer, Em Progresso e Concluídas. Os dados são persistidos em arquivo JSON e a aplicação pode ser executada facilmente via Docker Compose.

---

## Demonstração

![User Flow](docs/user-flow.png)

---

## Como rodar o projeto

### Pré-requisitos
- Docker e Docker Compose instalados
- (Opcional) Go e Node.js instalados para rodar localmente sem Docker

### Passos rápidos (recomendado)
1. Clone o repositório:
   ```sh
   git clone https://github.com/seu-usuario/desafio-fullstack-veritas.git
   cd desafio-fullstack-veritas
   ```
2. Suba tudo com Docker Compose:
   ```sh
   docker compose up --build
   ```
3. Acesse:
   - Frontend: http://localhost:3000
   - Backend: http://localhost:8080

### Rodar localmente (sem Docker)
- Backend:
  ```sh
  cd backend
  go run .
  ```
- Frontend:
  ```sh
  cd frontend
  npm install
  npm start
  ```

---

## Estrutura do Projeto
```
/backend      # Código Go (API REST, persistência JSON, testes, Dockerfile)
/frontend     # Código React (componentes, testes, Dockerfile)
/docs         # Documentação, user-flow.png, data-flow.png
```

---

## Funcionalidades
- Visualização de três colunas fixas: A Fazer, Em Progresso e Concluídas
- Adição de tarefas com título obrigatório e descrição opcional
- Edição, exclusão e movimentação de tarefas entre colunas (drag and drop ou botões)
- Feedback visual de loading e erro
- Persistência dos dados em arquivo JSON
- API RESTful completa (GET, POST, PUT, DELETE)
- Testes automatizados no backend e frontend
- Docker e Docker Compose para facilitar execução

---

## Decisões técnicas
- **Go**: Simplicidade, performance e fácil manipulação de arquivos.
- **React**: Componentização e experiência moderna.
- **Persistência em JSON**: Fácil para MVP e testes locais.
- **Docker**: Facilita rodar em qualquer ambiente.

---

## Limitações conhecidas
- Sem autenticação de usuários.
- Sem deploy em nuvem.
- Persistência apenas local (arquivo JSON).
- Não há testes E2E.

---

## Melhorias futuras
- Autenticação e autorização.
- Deploy em nuvem (ex: Heroku, AWS).
- Testes E2E (Cypress, Playwright).
- Responsividade avançada.
- Filtros e busca de tarefas.

---

## User Flow (Fluxo do Usuário)
Veja o arquivo `docs/user-flow.png`.

**Exemplo de fluxo:**
1. Usuário acessa o sistema.
2. Visualiza as três colunas do Kanban.
3. Adiciona uma nova tarefa.
4. Move tarefas entre colunas (drag and drop ou botões).
5. Edita ou exclui tarefas.
6. Todas as ações são salvas e refletidas na interface.

---

## Data Flow (opcional)
Veja o arquivo `docs/data-flow.png`.

- Frontend faz requisições REST para o backend.
- Backend manipula dados em memória e salva no arquivo `tasks.json`.

---

## Testes
- **Backend:**
  - Execute `go test` na pasta backend (cobre casos de sucesso e erro)
- **Frontend:**
  - Execute `npm test` na pasta frontend (exemplo em `TaskForm.test.js`)

---

## Contato
Dúvidas ou sugestões? Entre em contato pelo GitHub ou e-mail.
