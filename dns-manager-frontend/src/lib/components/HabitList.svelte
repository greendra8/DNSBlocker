<script>
  import { onMount } from 'svelte';
  import { browser } from '$app/environment';
  
  export let points;
  
  let habits = [];
  let newHabit = { name: '', goodPoints: 1, badPoints: 1 };
  let showNewHabitForm = false;
  let showRemoveButtons = false;
  
  onMount(() => {
    if (browser) {
      const savedHabits = localStorage.getItem('habits');
      if (savedHabits) {
        habits = JSON.parse(savedHabits);
      }
    }
  });
  
  function addHabit() {
    habits = [...habits, { ...newHabit, id: Date.now() }];
    newHabit = { name: '', goodPoints: 1, badPoints: 1 };
    saveHabits();
    showNewHabitForm = false;
  }
  
  function removeHabit(id) {
    habits = habits.filter(habit => habit.id !== id);
    saveHabits();
  }
  
  function toggleGoodHabit(habit) {
    points += habit.goodPoints;
    saveHabits();
  }
  
  function toggleBadHabit(habit) {
    points -= habit.badPoints;
    saveHabits();
  }
  
  function saveHabits() {
    if (browser) {
      localStorage.setItem('habits', JSON.stringify(habits));
    }
  }
</script>

<div class="habit-list">
  <div class="list-header">
    <h2>Habits</h2>
    <div>
      <button class="icon-button" on:click={() => showNewHabitForm = !showNewHabitForm}>
        {showNewHabitForm ? '−' : '+'}
      </button>
      <button class="icon-button" on:click={() => showRemoveButtons = !showRemoveButtons}>
        {showRemoveButtons ? '✓' : '−'}
      </button>
    </div>
  </div>
  
  {#if showNewHabitForm}
    <form on:submit|preventDefault={addHabit} class="new-habit-form">
      <input bind:value={newHabit.name} placeholder="New habit" required>
      <input type="number" bind:value={newHabit.goodPoints} min="1" placeholder="Good points" required>
      <input type="number" bind:value={newHabit.badPoints} min="1" placeholder="Bad points" required>
      <button type="submit">Add Habit</button>
    </form>
  {/if}
  
  <ul>
    {#each habits as habit}
      <li>
        <div class="habit-item">
          <button class="icon-button" on:click={() => toggleGoodHabit(habit)}>+</button>
          <span>{habit.name} (+{habit.goodPoints}/-{habit.badPoints})</span>
          <button class="icon-button" on:click={() => toggleBadHabit(habit)}>−</button>
        </div>
        {#if showRemoveButtons}
          <button class="remove-button" on:click={() => removeHabit(habit.id)}>Remove</button>
        {/if}
      </li>
    {/each}
  </ul>
</div>

<style>
  .habit-list {
    background-color: #1a4a1a; /* Green */
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
  .icon-button {
    font-size: 1.2em;
    background: none;
    border: none;
    cursor: pointer;
    color: #FFFFFF;
    padding: 0 5px;
  }
  ul {
    list-style-type: none;
    padding: 0;
  }
  li {
    margin-bottom: 0.5rem;
  }
  .habit-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
  }
  .habit-item span {
    flex-grow: 1;
    text-align: center;
  }
  .new-habit-form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  .new-habit-form input, .new-habit-form button {
    width: 100%;
    padding: 0.5rem;
    box-sizing: border-box;
  }
  button {
    background-color: #FFFFFF;
    color: #000033;
    border: none;
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
    cursor: pointer;
    transition: background-color 0.3s ease;
  }
  button:hover {
    background-color: #CCCCCC;
  }
  .remove-button {
    background-color: #ff4444;
    color: white;
    border: none;
    padding: 2px 5px;
    border-radius: 3px;
    cursor: pointer;
    margin-top: 5px;
  }
</style>