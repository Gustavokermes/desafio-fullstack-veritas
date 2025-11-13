import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import TaskForm from './TaskForm';

test('exibe erro se tentar salvar sem título', () => {
  const onSubmit = jest.fn();
  render(<TaskForm onSubmit={onSubmit} />);
  fireEvent.click(screen.getByText(/salvar/i));
  expect(onSubmit).not.toHaveBeenCalled();
  expect(screen.getByText(/título é obrigatório/i)).toBeInTheDocument();
});

test('envia dados corretos ao salvar', () => {
  const onSubmit = jest.fn();
  render(<TaskForm onSubmit={onSubmit} />);
  fireEvent.change(screen.getByLabelText(/título/i), { target: { value: 'Nova tarefa' } });
  fireEvent.click(screen.getByText(/salvar/i));
  expect(onSubmit).toHaveBeenCalledWith(expect.objectContaining({ title: 'Nova tarefa' }));
});