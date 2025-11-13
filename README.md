# Desafio Fullstack – Mini Kanban de Tarefas (React + Go)

## Autor
Gustavo Henrique Kermes Sousa

## Descrição
Este projeto é um Kanban simples fullstack, com backend em Go e frontend em React, que permite criar, editar, mover e excluir tarefas entre três colunas: A Fazer, Em Progresso e Concluídas. Os dados são persistidos em arquivo JSON e a aplicação pode ser executada facilmente via Docker Compose.

---

## Como rodar o projeto

### 1. Usando Docker Compose (recomendado)

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

### 2. Rodando localmente (sem Docker)

#### Backend
1. Instale o Go (versão 1.22 ou superior).
2. No terminal:
   ```sh
   cd backend
   go run .
   ```
3. O backend estará disponível em http://localhost:8080

#### Frontend
1. Instale o Node.js (versão 20 ou superior) e npm.
2. No terminal:
   ```sh
   cd frontend
   npm install
   npm start
   ```
3. O frontend estará disponível em http://localhost:3000

---

## Decisões técnicas
- **Go**: Escolhido pelo desempenho, simplicidade e fácil manipulação de arquivos para persistência local.
- **React**: Permite uma UI moderna, componentizada e fácil de evoluir.
- **Persistência em JSON**: Ideal para MVPs, facilita testes e não exige banco de dados.
- **Docker**: Garante que o projeto rode igual em qualquer ambiente, facilitando testes e deploy.

---

## Limitações conhecidas
- Não possui autenticação de usuários.
- Não faz deploy automático em nuvem.
- Persistência apenas local (arquivo `tasks.json`).
- Não há testes end-to-end (E2E).

---

## Melhorias futuras
- Adicionar autenticação e autorização de usuários.
- Deploy em nuvem (Heroku, AWS, etc).
- Testes E2E (Cypress, Playwright).
- Melhor responsividade e acessibilidade.
- Filtros, busca e ordenação de tarefas.
- Integração com banco de dados relacional.

---

## User Flow (Fluxo do Usuário)
Veja o arquivo `docs/user-flow.png` para um diagrama visual.

**Resumo do fluxo:**
1. Usuário acessa o sistema.
2. Visualiza as três colunas do Kanban.
3. Adiciona uma nova tarefa.
4. Move tarefas entre colunas (drag and drop ou botões).
5. Edita ou exclui tarefas.
6. Todas as ações são salvas e refletidas na interface.

---

## Testes
- **Backend:**
  - Execute `go test` na pasta backend (cobre casos de sucesso e erro)
- **Frontend:**
  - Execute `npm test` na pasta frontend (exemplo em `TaskForm.test.js`)

---

## Estrutura do Projeto
```
/backend      # Código Go (API REST, persistência JSON, testes, Dockerfile)
/frontend     # Código React (componentes, testes, Dockerfile)
/docs         # Documentação, user-flow.png, data-flow.png
```

---

**Dúvidas? Entre em contato!**
