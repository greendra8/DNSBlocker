<script>
  import { onMount } from 'svelte';
  import { browser } from '$app/environment';
  import axios from 'axios';
  
  export let points;
  
  let rewards = [];
  let blockedSites = [];
  let newReward = { name: '', points: 10 };
  let showNewRewardForm = false;
  
  const API_BASE_URL = 'https://localhost:443'; // Update this to your CoreDNS server address
  
  onMount(async () => {
    if (browser) {
      const savedRewards = localStorage.getItem('rewards');
      if (savedRewards) {
        rewards = JSON.parse(savedRewards);
      }
    }
    await fetchBlockedSites();
  });
  
  async function fetchBlockedSites() {
    try {
      const userId = localStorage.getItem('userId');
      const response = await axios.get(`${API_BASE_URL}/rules/${userId}`);
      blockedSites = response.data.map(site => ({
        ...site,
        points: 50 // Default cost to unblock
      }));
    } catch (error) {
      console.error('Error fetching blocked sites:', error);
    }
  }
  
  function addReward() {
    rewards = [...rewards, { ...newReward, id: Date.now() }];
    newReward = { name: '', points: 10 };
    saveRewards();
    showNewRewardForm = false;
  }
  
  function claimReward(reward) {
    if (points >= reward.points) {
      points -= reward.points;
      // Implement reward claiming logic here
      alert(`Reward claimed: ${reward.name}`);
    } else {
      alert('Not enough points to claim this reward');
    }
  }
  
  async function unblockSite(site) {
    if (points >= site.points) {
      points -= site.points;
      try {
        const userId = localStorage.getItem('userId');
        await axios.post(`${API_BASE_URL}/remove_rule`, {
          user_id: userId,
          rule_id: site.id
        });
        alert(`Site unblocked: ${site.domain}`);
        await fetchBlockedSites();
      } catch (error) {
        console.error('Error unblocking site:', error);
        alert('Error unblocking site');
      }
    } else {
      alert('Not enough points to unblock this site');
    }
  }
  
  function saveRewards() {
    if (browser) {
      localStorage.setItem('rewards', JSON.stringify(rewards));
    }
  }
</script>

<div class="reward-list">
  <div class="list-header">
    <h2>Rewards</h2>
    <button class="add-button" on:click={() => showNewRewardForm = !showNewRewardForm}>+</button>
  </div>
  
  {#if showNewRewardForm}
    <form on:submit|preventDefault={addReward} class="new-reward-form">
      <input bind:value={newReward.name} placeholder="New reward" required>
      <input type="number" bind:value={newReward.points} min="1" required>
      <button type="submit">Add Reward</button>
    </form>
  {/if}
  
  <ul>
    {#each rewards as reward}
      <li>
        {reward.name} (-{reward.points})
        <button on:click={() => claimReward(reward)}>Claim</button>
      </li>
    {/each}
  </ul>
  
  <h3>Blocked Sites</h3>
  <ul>
    {#each blockedSites as site}
      <li>
        {site.domain} (-{site.points})
        <button on:click={() => unblockSite(site)}>Unblock</button>
      </li>
    {/each}
  </ul>
</div>

<style>
  .reward-list {
    background-color: #4a1a4a; /* Purple */
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
  .add-button {
    font-size: 1.5em;
    background: none;
    border: none;
    cursor: pointer;
    color: #FFFFFF;
  }
  h2, h3 {
    font-family: 'Playfair Display', serif;
    margin-bottom: 1rem;
  }
  ul {
    list-style-type: none;
    padding: 0;
  }
  li {
    margin-bottom: 0.5rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .new-reward-form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  .new-reward-form input, .new-reward-form button {
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
</style>