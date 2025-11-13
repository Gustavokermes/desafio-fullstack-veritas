import React from 'react';
import KanbanColumn from './KanbanColumn';

const KanbanBoard = ({ tasks, onEdit, onDelete, onMove }) => {
    const columns = [
        { id: 'todo', title: 'A Fazer', status: 'A Fazer' },
        { id: 'inProgress', title: 'Em Progresso', status: 'Em Progresso' },
        { id: 'done', title: 'Concluídas', status: 'Concluídas' }
    ];

    // Garantir que tasks seja sempre um array
    const safeTasks = Array.isArray(tasks) ? tasks : [];

    return (
        <div className="kanban-board">
            {columns.map(column => (
                <KanbanColumn
                    key={column.id}
                    title={column.title}
                    tasks={safeTasks.filter(task => task.status === column.status)}
                    status={column.status}
                    onEdit={onEdit}
                    onDelete={onDelete}
                    onMove={onMove}
                />
            ))}
        </div>
    );
};

export default KanbanBoard;