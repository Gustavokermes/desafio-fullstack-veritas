import React from 'react';

const TaskCard = ({ task, onEdit, onDelete, onMove, availableStatuses }) => {
    const handleDragStart = (e) => {
        e.dataTransfer.setData('taskId', task.id);
        e.dataTransfer.setData('currentStatus', task.status);
    };

    const handleDelete = () => {
        if (window.confirm('Tem certeza que deseja excluir esta tarefa?')) {
            onDelete(task.id);
        }
    };

    return (
        <div 
            className="task-card"
            draggable
            onDragStart={handleDragStart}
        >
            <div className="task-header">
                <h3>{task.title}</h3>
                <div className="task-actions">
                    <button onClick={() => onEdit(task)}>Editar</button>
                    <button onClick={handleDelete}>Excluir</button>
                </div>
            </div>
            {task.description && <p>{task.description}</p>}
            <div className="move-actions">
                <span>Mover para:</span>
                {availableStatuses.map(status => (
                    <button
                        key={status}
                        className="btn-small"
                        onClick={() => onMove(task.id, status)}
                    >
                        {status}
                    </button>
                ))}
            </div>
        </div>
    );
};

export default TaskCard;
