import React, { useState, useEffect } from 'react';

const TaskForm = ({ task, onSubmit, onCancel }) => {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [status, setStatus] = useState('A Fazer');

    useEffect(() => {
        if (task) {
            setTitle(task.title);
            setDescription(task.description || '');
            setStatus(task.status);
        }
    }, [task]);

    const handleSubmit = (e) => {
        e.preventDefault();
        if (!title.trim()) {
            alert('Título é obrigatório');
            return;
        }
        
        onSubmit({
            id: task?.id,
            title: title.trim(),
            description: description.trim(),
            status
        });
    };

    return (
        <div className="modal-overlay">
            <div className="task-form">
                <h2>{task ? 'Editar Tarefa' : 'Nova Tarefa'}</h2>
                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label>Título *</label>
                        <input
                            type="text"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            required
                        />
                    </div>
                    <div className="form-group">
                        <label>Descrição</label>
                        <textarea
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            rows="3"
                        />
                    </div>
                    <div className="form-group">
                        <label>Status</label>
                        <select value={status} onChange={(e) => setStatus(e.target.value)}>
                            <option value="A Fazer">A Fazer</option>
                            <option value="Em Progresso">Em Progresso</option>
                            <option value="Concluídas">Concluídas</option>
                        </select>
                    </div>
                    <div className="form-actions">
                        <button type="submit" className="btn-primary">
                            {task ? 'Atualizar' : 'Criar'}
                        </button>
                        <button type="button" onClick={onCancel}>
                            Cancelar
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default TaskForm;
