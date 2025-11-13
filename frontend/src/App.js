import React, { useState, useEffect } from 'react';
import KanbanBoard from './components/KanbanBoard';
import TaskForm from './components/TaskForm';
import { api } from './services/api';
import './styles/App.css';

function App() {
    const [tasks, setTasks] = useState([]); // Garantir que começa como array vazio
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [editingTask, setEditingTask] = useState(null);
    const [showForm, setShowForm] = useState(false);

    useEffect(() => {
        fetchTasks();
    }, []);

    const fetchTasks = async () => {
        setLoading(true);
        setError('');
        try {
            const data = await api.getTasks();
            // Garantir que sempre seja um array, mesmo se a API retornar null
            setTasks(Array.isArray(data) ? data : []);
        } catch (err) {
            setError('Erro ao carregar tarefas: ' + err.message);
            setTasks([]); // Garantir array vazio em caso de erro
        } finally {
            setLoading(false);
        }
    };

    const handleCreateTask = async (taskData) => {
        setError('');
        try {
            const newTask = await api.createTask(taskData);
            setTasks(prevTasks => [...prevTasks, newTask]);
            setShowForm(false);
        } catch (err) {
            setError('Erro ao criar tarefa: ' + err.message);
        }
    };

    const handleUpdateTask = async (taskData) => {
        setError('');
        try {
            const updatedTask = await api.updateTask(taskData.id, taskData);
            setTasks(prevTasks => prevTasks.map(task => task.id === taskData.id ? updatedTask : task));
            setEditingTask(null);
            setShowForm(false);
        } catch (err) {
            setError('Erro ao atualizar tarefa: ' + err.message);
        }
    };

    const handleDeleteTask = async (taskId) => {
        setError('');
        try {
            await api.deleteTask(taskId);
            setTasks(prevTasks => prevTasks.filter(task => task.id !== taskId));
        } catch (err) {
            setError('Erro ao excluir tarefa: ' + err.message);
        }
    };

    const handleMoveTask = async (taskId, newStatus) => {
        const task = tasks.find(t => t.id === taskId);
        if (task) {
            await handleUpdateTask({ ...task, status: newStatus });
        }
    };

    return (
        <div className="app">
            <header className="app-header">
                <h1>Mini Kanban</h1>
                <button 
                    className="btn-primary"
                    onClick={() => {
                        setEditingTask(null);
                        setShowForm(true);
                    }}
                >
                    Nova Tarefa
                </button>
            </header>

            {error && <div className="error-message">{error}</div>}
            {loading && <div className="loading">Carregando...</div>}

            {showForm && (
                <TaskForm
                    task={editingTask}
                    onSubmit={editingTask ? handleUpdateTask : handleCreateTask}
                    onCancel={() => {
                        setShowForm(false);
                        setEditingTask(null);
                    }}
                />
            )}

            <KanbanBoard
                tasks={tasks || []} // Garantir que nunca seja null
                onEdit={(task) => {
                    setEditingTask(task);
                    setShowForm(true);
                }}
                onDelete={handleDeleteTask}
                onMove={handleMoveTask}
            />
        </div>
    );
}

export default App;