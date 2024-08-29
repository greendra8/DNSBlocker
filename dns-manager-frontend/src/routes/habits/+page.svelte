<script>
  import { onMount } from 'svelte';
  import { browser } from '$app/environment';
  import Nav from '$lib/components/Nav.svelte';
  import HabitList from '$lib/components/HabitList.svelte';
  import DailyList from '$lib/components/DailyList.svelte';
  import TodoList from '$lib/components/TodoList.svelte';
  import RewardList from '$lib/components/RewardList.svelte';
  
  let points = 0;
  
  onMount(() => {
    if (browser) {
      const savedPoints = localStorage.getItem('points');
      if (savedPoints) {
        points = parseInt(savedPoints, 10);
      }
    }
  });
  
  $: {
    if (browser) {
      localStorage.setItem('points', points.toString());
    }
  }
</script>

<Nav />

<main>
  <h1>Habit Tracker</h1>
  <p class="points">Current Points: {points}</p>
  
  <div class="lists-container">
    <HabitList bind:points />
    <DailyList bind:points />
    <TodoList bind:points />
    <RewardList bind:points />
  </div>
</main>

<style>
  @import url('https://fonts.googleapis.com/css2?family=Playfair+Display:wght@700&family=Roboto:wght@300;400&display=swap');

  :global(body) {
    font-family: 'Roboto', sans-serif;
    color: #FFFFFF;
    background-color: #000033;
    margin: 0;
    padding: 0;
  }

  main {
    width: 100%;
    max-width: none;
    padding: 2rem;
    box-sizing: border-box;
  }

  h1 {
    font-family: 'Playfair Display', serif;
    font-size: 2.5rem;
    color: #FFFFFF;
    margin-bottom: 1rem;
    text-align: center;
  }

  .points {
    font-size: 1.2rem;
    text-align: center;
    margin-bottom: 2rem;
  }

  .lists-container {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
  }

  :global(.lists-container > *) {
    background-color: rgba(255, 255, 255, 0.1);
    padding: 1.5rem;
    border-radius: 8px;
  }

  :global(.lists-container h2) {
    font-family: 'Playfair Display', serif;
    font-size: 1.8rem;
    margin-bottom: 1rem;
  }

  :global(.lists-container button) {
    background-color: #FFFFFF;
    color: #000033;
    border: none;
    padding: 0.8rem 1.5rem;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.3s ease;
  }

  :global(.lists-container button:hover) {
    background-color: #CCCCCC;
  }
</style>