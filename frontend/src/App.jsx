import { useEffect, useMemo, useState } from "react";

const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

const columns = [
  { id: "todo", title: "A Fazer", accent: "var(--todo)" },
  { id: "in_progress", title: "Em Progresso", accent: "var(--progress)" },
  { id: "done", title: "Concluidas", accent: "var(--done)" }
];

const emptyForm = {
  title: "",
  description: "",
  status: "todo"
};

export default function App() {
  const [tasks, setTasks] = useState([]);
  const [form, setForm] = useState(emptyForm);
  const [editingTask, setEditingTask] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [feedback, setFeedback] = useState(null);
  const [draggedId, setDraggedId] = useState(null);

  useEffect(() => {
    loadTasks();
  }, []);

  const counters = useMemo(() => {
    return columns.reduce((acc, column) => {
      acc[column.id] = tasks.filter((task) => task.status === column.id).length;
      return acc;
    }, {});
  }, [tasks]);

  async function request(path, options = {}) {
    const response = await fetch(`${API_URL}${path}`, {
      headers: { "Content-Type": "application/json" },
      ...options
    });

    if (!response.ok) {
      let message = "Nao foi possivel concluir a acao.";
      try {
        const body = await response.json();
        message = body.error ?? message;
      } catch {
        // Keep the generic message when the API sends no JSON body.
      }
      throw new Error(message);
    }

    if (response.status === 204) {
      return null;
    }

    return response.json();
  }

  async function loadTasks() {
    try {
      setLoading(true);
      const data = await request("/tasks");
      setTasks(data);
    } catch (error) {
      showFeedback(error.message, "error");
    } finally {
      setLoading(false);
    }
  }

  function showFeedback(message, type = "success") {
    setFeedback({ message, type });
    window.clearTimeout(showFeedback.timeout);
    showFeedback.timeout = window.setTimeout(() => setFeedback(null), 3200);
  }

  function openCreateModal(status = "todo") {
    setEditingTask(null);
    setForm({ ...emptyForm, status });
    setIsModalOpen(true);
  }

  function openEditModal(task) {
    setEditingTask(task);
    setForm({
      title: task.title,
      description: task.description,
      status: task.status
    });
    setIsModalOpen(true);
  }

  function closeModal() {
    if (saving) return;
    setIsModalOpen(false);
    setEditingTask(null);
    setForm(emptyForm);
  }

  async function handleSubmit(event) {
    event.preventDefault();
    if (!form.title.trim()) {
      showFeedback("Informe um titulo para a tarefa.", "error");
      return;
    }

    try {
      setSaving(true);
      const payload = {
        title: form.title.trim(),
        description: form.description.trim(),
        status: form.status
      };

      if (editingTask) {
        const updated = await request(`/tasks/${editingTask.id}`, {
          method: "PUT",
          body: JSON.stringify(payload)
        });
        setTasks((current) => current.map((task) => (task.id === updated.id ? updated : task)));
        showFeedback("Tarefa atualizada.");
      } else {
        const created = await request("/tasks", {
          method: "POST",
          body: JSON.stringify(payload)
        });
        setTasks((current) => [...current, created]);
        showFeedback("Tarefa criada.");
      }

      closeModal();
    } catch (error) {
      showFeedback(error.message, "error");
    } finally {
      setSaving(false);
    }
  }

  async function updateTaskStatus(task, nextStatus) {
    if (task.status === nextStatus) return;

    const previous = tasks;
    const optimistic = { ...task, status: nextStatus };
    setTasks((current) => current.map((item) => (item.id === task.id ? optimistic : item)));

    try {
      const updated = await request(`/tasks/${task.id}`, {
        method: "PUT",
        body: JSON.stringify({
          title: task.title,
          description: task.description,
          status: nextStatus
        })
      });
      setTasks((current) => current.map((item) => (item.id === task.id ? updated : item)));
      showFeedback("Tarefa movida.");
    } catch (error) {
      setTasks(previous);
      showFeedback(error.message, "error");
    }
  }

  async function deleteTask(task) {
    const confirmed = window.confirm(`Excluir "${task.title}"?`);
    if (!confirmed) return;

    const previous = tasks;
    setTasks((current) => current.filter((item) => item.id !== task.id));

    try {
      await request(`/tasks/${task.id}`, { method: "DELETE" });
      showFeedback("Tarefa excluida.");
    } catch (error) {
      setTasks(previous);
      showFeedback(error.message, "error");
    }
  }

  function moveByStep(task, direction) {
    const index = columns.findIndex((column) => column.id === task.status);
    const nextColumn = columns[index + direction];
    if (!nextColumn) return;
    updateTaskStatus(task, nextColumn.id);
  }

  function handleDrop(status) {
    const task = tasks.find((item) => item.id === draggedId);
    setDraggedId(null);
    if (task) {
      updateTaskStatus(task, status);
    }
  }

  return (
    <main className="app-shell">
      <header className="topbar">
        <div>
          <span className="eyebrow">Desafio Fullstack</span>
          <h1>Mini Kanban de Tarefas</h1>
        </div>

        <button className="primary-button" onClick={() => openCreateModal()}>
          <span aria-hidden="true">+</span>
          Nova tarefa
        </button>
      </header>

      <section className="metrics" aria-label="Resumo das tarefas">
        {columns.map((column) => (
          <div className="metric" key={column.id}>
            <span className="metric-dot" style={{ background: column.accent }} />
            <span>{column.title}</span>
            <strong>{counters[column.id] ?? 0}</strong>
          </div>
        ))}
      </section>

      {feedback && (
        <div className={`toast ${feedback.type}`} role="status">
          {feedback.message}
        </div>
      )}

      {loading ? (
        <section className="loading-state" aria-live="polite">
          Carregando tarefas...
        </section>
      ) : (
        <section className="board" aria-label="Quadro kanban">
          {columns.map((column) => {
            const columnTasks = tasks.filter((task) => task.status === column.id);

            return (
              <article
                className="kanban-column"
                key={column.id}
                onDragOver={(event) => event.preventDefault()}
                onDrop={() => handleDrop(column.id)}
              >
                <div className="column-header" style={{ borderTopColor: column.accent }}>
                  <h2>{column.title}</h2>
                  <button
                    className="icon-button"
                    title={`Adicionar em ${column.title}`}
                    aria-label={`Adicionar em ${column.title}`}
                    onClick={() => openCreateModal(column.id)}
                  >
                    +
                  </button>
                </div>

                <div className="task-list">
                  {columnTasks.length === 0 ? (
                    <div className="empty-column">Sem tarefas</div>
                  ) : (
                    columnTasks.map((task) => (
                      <div
                        className="task-card"
                        key={task.id}
                        draggable
                        onDragStart={() => setDraggedId(task.id)}
                        onDragEnd={() => setDraggedId(null)}
                      >
                        <div className="task-content">
                          <h3>{task.title}</h3>
                          {task.description && <p>{task.description}</p>}
                        </div>

                        <div className="task-actions" aria-label={`Acoes para ${task.title}`}>
                          <button
                            className="icon-button"
                            title="Mover para coluna anterior"
                            aria-label="Mover para coluna anterior"
                            disabled={task.status === "todo"}
                            onClick={() => moveByStep(task, -1)}
                          >
                            &lt;
                          </button>
                          <button
                            className="icon-button"
                            title="Editar"
                            aria-label="Editar"
                            onClick={() => openEditModal(task)}
                          >
                            E
                          </button>
                          <button
                            className="icon-button danger"
                            title="Excluir"
                            aria-label="Excluir"
                            onClick={() => deleteTask(task)}
                          >
                            x
                          </button>
                          <button
                            className="icon-button"
                            title="Mover para proxima coluna"
                            aria-label="Mover para proxima coluna"
                            disabled={task.status === "done"}
                            onClick={() => moveByStep(task, 1)}
                          >
                            &gt;
                          </button>
                        </div>
                      </div>
                    ))
                  )}
                </div>
              </article>
            );
          })}
        </section>
      )}

      {isModalOpen && (
        <div className="modal-backdrop" role="presentation" onMouseDown={closeModal}>
          <form className="task-modal" onSubmit={handleSubmit} onMouseDown={(event) => event.stopPropagation()}>
            <div className="modal-header">
              <h2>{editingTask ? "Editar tarefa" : "Nova tarefa"}</h2>
              <button className="icon-button" type="button" title="Fechar" aria-label="Fechar" onClick={closeModal}>
                x
              </button>
            </div>

            <label>
              Titulo
              <input
                autoFocus
                maxLength={80}
                value={form.title}
                onChange={(event) => setForm((current) => ({ ...current, title: event.target.value }))}
                placeholder="Ex.: Revisar endpoints"
              />
            </label>

            <label>
              Descricao
              <textarea
                maxLength={280}
                value={form.description}
                onChange={(event) => setForm((current) => ({ ...current, description: event.target.value }))}
                placeholder="Detalhe rapido da tarefa"
                rows={4}
              />
            </label>

            <label>
              Coluna
              <select
                value={form.status}
                onChange={(event) => setForm((current) => ({ ...current, status: event.target.value }))}
              >
                {columns.map((column) => (
                  <option key={column.id} value={column.id}>
                    {column.title}
                  </option>
                ))}
              </select>
            </label>

            <footer className="modal-actions">
              <button className="ghost-button" type="button" onClick={closeModal}>
                Cancelar
              </button>
              <button className="primary-button" type="submit" disabled={saving}>
                {saving ? "Salvando..." : "Salvar"}
              </button>
            </footer>
          </form>
        </div>
      )}
    </main>
  );
}
