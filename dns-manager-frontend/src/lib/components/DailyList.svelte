<script>
  import { onMount } from 'svelte';
  import { browser } from '$app/environment';
  
  export let points;
  
  let dailies = [];
  let newDaily = { name: '', points: 1 };
  let showNewDailyForm = false;
  let showRemoveButtons = false;
  
  onMount(() => {
    if (browser) {
      const savedDailies = localStorage.getItem('dailies');
      if (savedDailies) {
        dailies = JSON.parse(savedDailies);
      }
      resetDailies();
    }
  });
  
  function addDaily() {
    dailies = [...dailies, { ...newDaily, id: Date.now(), lastCompleted: null }];
    newDaily = { name: '', points: 1 };
    saveDailies();
    showNewDailyForm = false;
  }
  
  function toggleDaily(daily) {
    const today = new Date().toDateString();
    if (daily.lastCompleted === today) {
      daily.lastCompleted = null;
      points -= daily.points;
    } else {
      daily.lastCompleted = today;
      points += daily.points;
    }
    saveDailies();
  }
  
  function resetDailies() {
    const today = new Date().toDateString();
    dailies = dailies.map(daily => {
      if (daily.lastCompleted !== today) {
        daily.lastCompleted = null;
      }
      return daily;
    });
    saveDailies();
  }
  
  function saveDailies() {
    if (browser) {
      localStorage.setItem('dailies', JSON.stringify(dailies));
    }
  }
  
  function removeDaily(id) {
    dailies = dailies.filter(daily => daily.id !== id);
    saveDailies();
  }
</script>

<div class="daily-list">
  <div class="list-header">
    <h2>Dailies</h2>
    <div>
      <button class="icon-button" on:click={() => showNewDailyForm = !showNewDailyForm}>
        {showNewDailyForm ? '−' : '+'}
      </button>
      <button class="icon-button" on:click={() => showRemoveButtons = !showRemoveButtons}>
        {showRemoveButtons ? '✓' : '−'}
      </button>
    </div>
  </div>
  
  {#if showNewDailyForm}
    <form on:submit|preventDefault={addDaily} class="new-daily-form">
      <input bind:value={newDaily.name} placeholder="New daily" required>
      <input type="number" bind:value={newDaily.points} min="1" required>
      <button type="submit">Add Daily</button>
    </form>
  {/if}
  
  <ul>
    {#each dailies as daily}
      <li>
        <div class="daily-item">
          <input 
            type="checkbox" 
            checked={daily.lastCompleted === new Date().toDateString()}
            on:change={() => toggleDaily(daily)}
          >
          <span>{daily.name} ({daily.points})</span>
          {#if showRemoveButtons}
            <button class="remove-button" on:click={() => removeDaily(daily.id)}>Remove</button>
          {/if}
        </div>
      </li>
    {/each}
  </ul>
</div>

<style>
  .daily-list {
    background-color: #4a1a1a; /* Red */
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
  .daily-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
  }
  .daily-item span {
    flex-grow: 1;
    margin-left: 10px;
  }
  .new-daily-form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  .new-daily-form input, .new-daily-form button {
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
  }
</style>