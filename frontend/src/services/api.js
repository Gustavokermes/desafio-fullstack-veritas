const API_BASE_URL = 'http://localhost:8080';

export const api = {
    async getTasks() {
        const response = await fetch(`${API_BASE_URL}/tasks`);
        if (!response.ok) throw new Error('Erro ao buscar tarefas');
        return response.json();
    },

    async createTask(task) {
        const response = await fetch(`${API_BASE_URL}/tasks`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(task)
        });
        if (!response.ok) throw new Error('Erro ao criar tarefa');
        return response.json();
    },

    async updateTask(id, task) {
        const response = await fetch(`${API_BASE_URL}/tasks/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(task)
        });
        if (!response.ok) throw new Error('Erro ao atualizar tarefa');
        return response.json();
    },

    async deleteTask(id) {
        const response = await fetch(`${API_BASE_URL}/tasks/${id}`, {
            method: 'DELETE'
        });
        if (!response.ok) throw new Error('Erro ao excluir tarefa');
    }
};
