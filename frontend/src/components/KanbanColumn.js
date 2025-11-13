import React from 'react';
import TaskCard from './TaskCard';

const KanbanColumn = ({ title, tasks, status, onEdit, onDelete, onMove }) => {
    const otherStatuses = ['A Fazer', 'Em Progresso', 'Concluídas'].filter(s => s !== status);

    const handleDragOver = (e) => {
        e.preventDefault();
        e.currentTarget.classList.add('drag-over');
    };

    const handleDragLeave = (e) => {
        e.preventDefault();
        e.currentTarget.classList.remove('drag-over');
    };

    const handleDrop = (e) => {
        e.preventDefault();
        e.currentTarget.classList.remove('drag-over');
        
        const taskId = e.dataTransfer.getData('taskId');
        const currentStatus = e.dataTransfer.getData('currentStatus');
        
        if (currentStatus !== status) {
            onMove(taskId, status);
        }
    };

    return (
        <div 
            className="kanban-column"
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
        >
            <h2>{title} ({tasks.length})</h2>
            <div className="tasks-list">
                {tasks.map(task => (
                    <TaskCard
                        key={task.id}
                        task={task}
                        onEdit={onEdit}
                        onDelete={onDelete}
                        onMove={onMove}
                        availableStatuses={otherStatuses}
                    />
                ))}
                {tasks.length === 0 && (
                    <div className="empty-state">Nenhuma tarefa</div>
                )}
            </div>
        </div>
    );
};

export default KanbanColumn;
