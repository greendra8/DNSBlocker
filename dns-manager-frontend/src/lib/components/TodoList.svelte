<script>
  import { onMount } from 'svelte';
  import { browser } from '$app/environment';
  
  export let points;
  
  let todos = [];
  let newTodo = { name: '', points: 1 };
  let showNewTodoForm = false;
  
  onMount(() => {
    if (browser) {
      const savedTodos = localStorage.getItem('todos');
      if (savedTodos) {
        todos = JSON.parse(savedTodos);
      }
    }
  });
  
  function addTodo() {
    todos = [...todos, { ...newTodo, id: Date.now() }];
    newTodo = { name: '', points: 1 };
    saveTodos();
    showNewTodoForm = false;
  }
  
  function completeTodo(todo) {
    points += todo.points;
    todos = todos.filter(t => t.id !== todo.id);
    saveTodos();
  }
  
  function saveTodos() {
    if (browser) {
      localStorage.setItem('todos', JSON.stringify(todos));
    }
  }
</script>

<div class="todo-list">
  <div class="list-header">
    <h2>To-Do's</h2>
    <button class="add-button" on:click={() => showNewTodoForm = !showNewTodoForm}>
      {showNewTodoForm ? '−' : '+'}
    </button>
  </div>
  
  {#if showNewTodoForm}
    <form on:submit|preventDefault={addTodo} class="new-todo-form">
      <input bind:value={newTodo.name} placeholder="New to-do" required>
      <input type="number" bind:value={newTodo.points} min="1" required>
      <button type="submit">Add To-Do</button>
    </form>
  {/if}
  
  <ul>
    {#each todos as todo}
      <li>
        <label>
          <input type="checkbox" on:change={() => completeTodo(todo)}>
          {todo.name} (+{todo.points})
        </label>
      </li>
    {/each}
  </ul>
</div>

<style>
  .todo-list {
    background-color: #1a4a4a; /* Teal */
    padding: 1.5rem;
    border-radius: 8px;
    color: #FFFFFF;
  }
  .list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  h2 {
    font-family: 'Playfair Display', serif;
    margin-bottom: 1rem;
  }
  .add-button {
    font-size: 1.5em;
    background: none;
    border: none;
    cursor: pointer;
    color: #FFFFFF;
  }
  ul {
    list-style-type: none;
    padding: 0;
  }
  li {
    margin-bottom: 0.5rem;
  }
  label {
    display: flex;
    align-items: center;
  }
  input[type="checkbox"] {
    margin-right: 10px;
  }
  .new-todo-form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  .new-todo-form input, .new-todo-form button {
    width: 100%;
    padding: 0.5rem;
    box-sizing: border-box;
  }
  button {
    background-color: #FFFFFF;
    color: #1a4a4a;
    border: none;
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
    cursor: pointer;
    transition: background-color 0.3s ease;
  }
  button:hover {
    background-color: #CCCCCC;
  }
</style>