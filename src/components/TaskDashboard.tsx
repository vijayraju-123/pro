'use client';

import { useEffect, useState } from 'react';

export default function TaskDashboard() {
  const [tasks, setTasks] = useState<string[]>([]);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/ws');
    ws.onmessage = (event) => {
      console.log('Received:', event.data);
    };
    fetch('/api/tasks')
      .then((res) => res.json())
      .then((data) => setTasks(data.tasks));

    return () => ws.close();
  }, []);

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Task Dashboard</h1>
      <ul>
        {tasks.map((task, index) => (
          <li key={index} className="p-2 bg-white rounded shadow mb-2">
            {task}
          </li>
        ))}
      </ul>
    </div>
  );
}